// MFP - Multi-Function Printers and scanners toolkit
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions
//
// mfp-test command definition

package test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"math"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/OpenPrinting/go-mfp/argv"
	"github.com/OpenPrinting/go-mfp/internal/evaluate"
	"github.com/OpenPrinting/go-mfp/internal/testutils"
	"github.com/OpenPrinting/go-mfp/log"
	"github.com/OpenPrinting/go-mfp/modeling"
	"github.com/OpenPrinting/go-mfp/transport"
	"github.com/h2non/bimg"
)

// defaultTCPPort 0 asks the OS to pick a free port automatically,
// avoiding conflicts with other services (e.g. ipp-usb uses port 60000).
const defaultTCPPort = 0


// queueNamePrefix is prepended to the sanitised printer model name
// to form the CUPS queue name (e.g. "mfp-test-xerox-b235").
const queueNamePrefix = "mfp-test"

// Command is the mfp-test command description.
var Command = argv.Command{
	Name: "mfp-test",
	Help: "Print system testing pipeline",
	Options: []argv.Option{
		{
			Name:      "-m",
			Aliases:   []string{"--model"},
			Help:      "printer model file",
			HelpArg:   "file",
			Required:  true,
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		{
			Name:      "-P",
			Aliases:   []string{"--port"},
			Help:      "IPP server TCP port (default: OS-assigned free port)",
			HelpArg:   "port",
			Singleton: true,
			Validate:  argv.ValidateUint16,
		},
		{
			Name:      "-n",
			Aliases:   []string{"--name"},
			Help:      "CUPS queue name (default: mfp-test-<printer-model>)",
			HelpArg:   "name",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name:      "-o",
			Aliases:   []string{"--output"},
			Help:      "write JSON report to file",
			HelpArg:   "file",
			Singleton: true,
			Validate:  argv.ValidateAny,
			Complete:  argv.CompleteOSPath,
		},
		{
			Name:      "--threshold",
			Help:      fmt.Sprintf("minimum similarity score to pass (0.0-1.0, default %.2f)", defaultThreshold),
			HelpArg:   "score",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name:      "--timeout",
			Help:      fmt.Sprintf("document capture timeout (default %s)", defaultTimeout),
			HelpArg:   "duration",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name: "--list",
			Help: "list all test configurations (full batch matrix) and exit",
		},
		{
			Name: "--list-quick",
			Help: "list quick test configurations and exit",
		},
		{
			Name:      "--batch",
			Help:      "run all test configurations",
			Singleton: true,
		},
		{
			Name:      "--single",
			Help:      "run a single test configuration by name (sides/color-mode/format)",
			HelpArg:   "name",
			Singleton: true,
			Validate:  argv.ValidateAny,
		},
		{
			Name: "--quick",
			Help: "run reduced test matrix (all sides × all color modes, first format only)",
		},
		{
			Name: "--keep",
			Help: "save captured document bytes to current directory",
		},
		{
			Name:    "-v",
			Aliases: []string{"--verbose"},
			Help:    "enable verbose output",
		},
		argv.HelpOption,
	},
	Handler: cmdTestHandler,
}

// cmdTestHandler is the top-level handler for the mfp-test command.
func cmdTestHandler(ctx context.Context, inv *argv.Invocation) error {
	level := log.LevelInfo
	if inv.Flag("-v") {
		level = log.LevelTrace
	}
	logger := log.NewLogger(level, log.Console)
	ctx = log.NewContext(ctx, logger)

	// Model file is required: without it, NewIPPServer() returns nil.
	modelfile, ok := inv.Get("-m")
	if !ok {
		return fmt.Errorf("model file required: use -m <file>")
	}

	model, err := modeling.NewModel()
	if err != nil {
		return err
	}
	defer model.Close()

	if err := model.Load(modelfile); err != nil {
		return fmt.Errorf("load model %q: %w", modelfile, err)
	}

	// Parse port number
	port := defaultTCPPort
	if portStr, ok := inv.Get("-P"); ok {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid port %q: %w", portStr, err)
		}
		port = p
	}

	// Create document capture backend
	capture := newDocumentCapture()

	// Create virtual IPP printer from model and hook capture into it.
	ippPrinter := model.NewIPPPrinter()
	if ippPrinter == nil {
		return fmt.Errorf("model has no IPP printer attributes")
	}
	ippPrinter.SetPrintBackend(capture)

	// Register IPP handler on the URL path /ipp/print
	mux := transport.NewPathMux()
	mux.Add("/ipp/print", ippPrinter)

	// Open TCP port and start HTTP server (IPP runs over HTTP)
	addr := fmt.Sprintf("localhost:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// When port 0 was requested the OS assigns a free port; read it back.
	actualPort := ln.Addr().(*net.TCPAddr).Port

	srvr := transport.NewServer(ctx, nil, mux)
	log.Info(ctx, "virtual IPP printer at ipp://localhost:%d/ipp/print", actualPort)
	go srvr.Serve(ln)
	defer srvr.Close()

	// Build the IPP URL with the actual assigned port.
	ippURL := fmt.Sprintf("ipp://localhost:%d/ipp/print", actualPort)

	// Read printer capabilities directly from the model, avoiding an
	// HTTP round-trip to the virtual printer.
	caps, err := capsFromModel(model)
	if err != nil {
		return fmt.Errorf("read printer capabilities: %w", err)
	}

	// Determine CUPS queue name: use -n override or derive from printer model.
	queueName := queueNamePrefix + "-" + safeModelName(caps.ModelName)
	if name, ok := inv.Get("-n"); ok {
		queueName = name
	}

	// Register virtual printer with CUPS
	if err := CreateCUPSQueue(ctx, queueName, ippURL); err != nil {
		return err
	}
	// WithoutCancel preserves logging and values from ctx but
	// prevents cancellation from stopping the cleanup operation.
	defer RemoveCUPSQueue(context.WithoutCancel(ctx), queueName)

	log.Info(ctx, "CUPS queue %q ready at %s", queueName, ippURL)

	// --list: print all configurations (full batch matrix) and exit.
	if inv.Flag("--list") {
		for _, cfg := range batchMatrix(caps) {
			fmt.Println(cfg.Name)
		}
		return nil
	}

	// --list-quick: print quick test configurations and exit.
	if inv.Flag("--list-quick") {
		for _, cfg := range quickMatrix(caps) {
			fmt.Println(cfg.Name)
		}
		return nil
	}

	// Determine which test configurations to run.
	// One of --batch, --quick, or --single must be explicitly set.
	var configs []testConfig
	switch {
	case inv.Flag("--batch"):
		configs = batchMatrix(caps)
	case inv.Flag("--quick"):
		configs = quickMatrix(caps)
	default:
		if spec, ok := inv.Get("--single"); ok {
			cfg, err := singleConfig(spec, caps)
			if err != nil {
				return err
			}
			configs = []testConfig{*cfg}
		} else {
			return fmt.Errorf("no test mode specified: use --batch, --quick, or --single")
		}
	}

	// Parse similarity threshold.
	threshold := defaultThreshold
	if ts, ok := inv.Get("--threshold"); ok {
		var t float64
		if _, err := fmt.Sscanf(ts, "%f", &t); err != nil {
			return fmt.Errorf("invalid threshold %q: %w", ts, err)
		}
		threshold = t
	}

	// Parse capture timeout.
	timeout := defaultTimeout
	if ts, ok := inv.Get("--timeout"); ok {
		d, err := time.ParseDuration(ts)
		if err != nil {
			return fmt.Errorf("invalid timeout %q: %w", ts, err)
		}
		timeout = d
	}

	// Set up image evaluator using the embedded enhanced_comparison.py.
	var eval *evaluate.Evaluator
	e, err := evaluate.NewDefaultEvaluator(defaultComparatorScript)
	if err != nil {
		log.Info(ctx, "image evaluation disabled (embedded comparator unavailable): %v", err)
	} else {
		defer e.Close()
		eval = e
	}

	keep := inv.Flag("--keep")
	verbose := inv.Flag("-v")

	// Run each test configuration and collect results.
	var results []*testResult
	for _, cfg := range configs {
		log.Info(ctx, "running test: %s", cfg.Name)
		result, err := runTest(ctx, cfg, queueName, capture, eval, threshold, timeout, keep, verbose)
		if err != nil {
			log.Info(ctx, "FAIL %s: %v", cfg.Name, err)
			results = append(results, &testResult{Config: cfg, Score: 0, Passed: false})
			continue
		}
		if result.Passed {
			log.Info(ctx, "PASS %s (score=%.4f)", cfg.Name, result.Score)
		} else {
			log.Info(ctx, "FAIL %s (score=%.4f < threshold=%.4f)", cfg.Name, result.Score, threshold)
		}
		results = append(results, result)
	}

	// Print summary table and per-metric details.
	printSummaryTable(results)
	for _, res := range results {
		if verbose || !res.Passed {
			printVerboseDetails(res)
		}
	}

	// Write JSON report if --output is set.
	if outPath, ok := inv.Get("-o"); ok {
		report := buildReport(results)
		if err := writeReport(outPath, report); err != nil {
			log.Info(ctx, "report: %v", err)
		} else {
			log.Info(ctx, "report written to %s", outPath)
		}
	}

	return nil
}

// generateTestPNG writes the UEIT RGB test image to a temp file and returns
// its path. UEIT images have varied content (colour patches, gradients) that
// provides enough features for the full evaluator metric suite.
// The caller is responsible for removing the file after use.
func generateTestPNG() (string, error) {
	return writeTempPNG(testutils.Images.PNG100x75rgb8)
}

// generateGrayscaleTestPNG writes the UEIT greyscale test image to a temp
// file and returns its path. Used when the printer operates in monochrome
// mode; the native greyscale image avoids dependence on CUPS colour-to-grey
// conversion formula.
// The caller is responsible for removing the file after use.
func generateGrayscaleTestPNG() (string, error) {
	return writeTempPNG(testutils.Images.PNG100x75gray8)
}

// writeTempPNG builds a 300×300 test PNG from the given UEIT source and
// writes it to a temp file. Three post-processing steps maximise evaluator
// coverage:
//
//  1. bimg force-resize to 300×300 gives enough pixels for ORB and texture
//     metrics (the raw 100×75 UEIT is too small).
//  2. addTestBorder draws a 1-pixel white margin followed by a 5-pixel black
//     frame. The white margin guarantees a Canny gradient ≥255 on all four
//     sides of the frame (white→black, not boundary→black), forming a clean
//     closed rectangular contour that RETR_EXTERNAL picks up cleanly for
//     avg_rectangularity. UEIT content is kept for the inner 288×288 (92%).
//  3. injectPNGDPI embeds 300 DPI so CUPS prints the image at its native
//     printer resolution (1 inch × 1 inch → 300×300 captured pixels). This
//     eliminates the 1.5× CUPS upscale / 0.67× evaluator downscale that
//     otherwise degrades PSNR and edge-similarity scores.
func writeTempPNG(data []byte) (string, error) {
	scaled, err := bimg.NewImage(data).Process(bimg.Options{
		Width: 300, Height: 300, Force: true, Type: bimg.PNG,
	})
	if err != nil {
		return "", fmt.Errorf("mfp-test: scale test image: %w", err)
	}
	gridded, err := addTestBorder(scaled)
	if err != nil {
		return "", fmt.Errorf("mfp-test: add border: %w", err)
	}
	final, err := injectPNGDPI(gridded, 300)
	if err != nil {
		return "", fmt.Errorf("mfp-test: inject DPI: %w", err)
	}
	f, err := os.CreateTemp("", "mfp-test-*.png")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(final); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), f.Close()
}

// addTestBorder draws an inset black frame on the 300×300 image, with a
// 1-pixel white outer margin before the frame.
//
// Why the white margin matters for avg_rectangularity:
// Canny edge detection uses a 3×3 Sobel kernel. For the frame's outer edge to
// form a CLOSED, STRONGLY-DETECTED rectangular contour, we need the pixels
// immediately above/left of the frame to be white (255). With no margin, the
// frame sits at the image boundary where Canny uses reflection padding — the
// reflected pixel equals the frame pixel (both black), giving gradient ≈ 0 and
// no edge. With a 1-pixel white margin, the gradient at the frame's outer edge
// is √2×255 ≈ 360 at corners and 255 along sides, well above the Canny upper
// threshold of 150 — creating a definite edge on all four sides simultaneously.
//
// The UEIT content starts 6 pixels inside the frame edge (1 margin + 5 border),
// so UEIT texture edges can never connect to the frame contour. RETR_EXTERNAL
// returns the frame rectangle cleanly as a single external contour with
// area ≈ (298)² ≈ 88 800 px >> 1000 px threshold, approximated as exactly 4
// sides, giving avg_rectangularity ≈ 1.0.
//
// Compared with thin internal grid lines (previous approach), the thick outer
// frame avoids the MSE spike caused by printer rasterization of 3-pixel-wide
// lines. The white margin also reduces content coverage only marginally
// (288×288 = 92% of image area keeps UEIT texture).
func addTestBorder(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("border: decode: %w", err)
	}
	b := img.Bounds()
	dst := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	const (
		size   = 300
		margin = 1 // 1-pixel white strip; ensures Canny sees gradient ≥255 at the frame edge
		bw     = 5 // frame width in pixels; thick enough for reliable printer reproduction
	)
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if x < margin || x >= size-margin || y < margin || y >= size-margin {
				dst.Set(x, y, white)
			} else if x < margin+bw || x >= size-margin-bw ||
				y < margin+bw || y >= size-margin-bw {
				dst.Set(x, y, black)
			}
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, fmt.Errorf("border: encode: %w", err)
	}
	return buf.Bytes(), nil
}

// injectPNGDPI inserts a pHYs chunk with the specified DPI immediately after
// the IHDR chunk of a PNG byte slice. Because CUPS uses the pHYs chunk to
// determine the printed physical size, embedding 300 DPI causes a 300×300
// image to print at exactly 1×1 inch on A4 at 300 DPI, matching the capture
// resolution and eliminating the evaluator's lossy resize step.
func injectPNGDPI(data []byte, dpi int) ([]byte, error) {
	// PNG structure: 8-byte signature + IHDR (4+4+13+4 = 25 bytes) = 33 bytes.
	const ihdrEnd = 33
	if len(data) < ihdrEnd {
		return data, nil
	}
	// Compute pixels-per-metre: 1 inch = 0.0254 metres.
	ppm := uint32(math.Round(float64(dpi) / 0.0254))
	var physData [9]byte
	binary.BigEndian.PutUint32(physData[0:4], ppm) // X density
	binary.BigEndian.PutUint32(physData[4:8], ppm) // Y density
	physData[8] = 1                                  // unit = metre

	h := crc32.NewIEEE()
	h.Write([]byte("pHYs"))
	h.Write(physData[:])

	chunk := make([]byte, 21) // 4 len + 4 type + 9 data + 4 CRC
	binary.BigEndian.PutUint32(chunk[0:4], 9)
	copy(chunk[4:8], "pHYs")
	copy(chunk[8:17], physData[:])
	binary.BigEndian.PutUint32(chunk[17:21], h.Sum32())

	result := make([]byte, 0, len(data)+21)
	result = append(result, data[:ihdrEnd]...)
	result = append(result, chunk...)
	result = append(result, data[ihdrEnd:]...)
	return result, nil
}

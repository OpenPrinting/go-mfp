// MFP - Multi-Function Printers and scanners toolkit
// The "mfp-test" command
//
// Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
// See LICENSE for license terms and conditions

// # mfp-test — Print System Testing Pipeline
//
// mfp-test is an automated end-to-end testing tool for driverless printing.
// It simulates a virtual IPP printer, sends print jobs through the real CUPS
// printing stack, captures the output, and evaluates image quality to verify
// that the full pipeline works correctly — without requiring any physical
// hardware.
//
// # Motivation
//
// Modern driverless printing relies on a complex software stack: the
// application formats the document, CUPS selects a conversion filter based on
// the printer's advertised capabilities, the filter produces the printer's
// native format (PWG Raster, PDF, PostScript, …), and the result is delivered
// over IPP. A defect anywhere in this chain — a wrong colour space, a broken
// raster encoder, a filter that silently drops pages — produces a bad printout
// that is only discovered when someone looks at the paper.
//
// mfp-test exercises the entire chain automatically: it generates a known test
// image, sends it through CUPS exactly as a real application would, and
// compares what comes out of CUPS against the original. Failures are reported
// with a numeric quality score and per-metric details so the root cause can be
// diagnosed quickly.
//
// # Architecture
//
// The pipeline is built from five components that run in sequence for each
// test configuration:
//
//  1. Virtual IPP printer (modeling + proto/ipp)
//
//     A printer model file (a Python script) defines the printer's IPP
//     attribute set: supported paper sizes, colour modes, document formats,
//     duplex capability, and so on. go-mfp's modeling package loads the file
//     and starts a real HTTP/IPP server on a local TCP port. The server speaks
//     IPP/2.x and responds correctly to Get-Printer-Attributes, Create-Job,
//     and Send-Document operations.
//
//  2. CUPS queue (cupsctl.go)
//
//     mfp-test registers the virtual printer as a real CUPS queue using
//     lpadmin. From CUPS's perspective there is no difference between this
//     queue and one that points at a physical device on the network. CUPS
//     applies its full filter pipeline — format negotiation, colour
//     conversion, rasterisation — before delivering the job.
//     The queue is removed automatically when mfp-test exits.
//
//  3. Document capture (capture.go)
//
//     The virtual printer's print backend is wired to a DocumentCapture
//     object instead of a real device driver. When CUPS delivers a job,
//     DocumentCapture stores the raw bytes and their MIME type in memory and
//     signals the test runner that the document has arrived.
//
//  4. Raster conversion (raster.go)
//
//     The captured bytes are converted to PNG so that image comparison is
//     format-independent. The converter dispatches on MIME type:
//
//       - image/pwg-raster, image/urf:
//         Decoded by the vendored pure-Go cupsraster package
//         (derived from github.com/rusq/thermoprint/cupsraster, MIT licence;
//         see CREDITS.txt). No external dependency required.
//
//       - application/pdf, image/jpeg, image/tiff, image/webp, image/gif:
//         Converted by bimg (github.com/h2non/bimg), a Go binding for
//         libvips. Requires libvips-dev to be installed.
//
//       - application/postscript, application/vnd.cups-postscript:
//         Rendered by a Ghostscript subprocess (gs). The first page is
//         rendered at the printer's negotiated resolution (printer-resolution-
//         supported[0]), defaulting to 300 DPI. Requires ghostscript to be
//         installed.
//
//  5. Image evaluation (internal/evaluate)
//
//     The captured PNG is compared against the original test image using
//     ImageComparator from Sanskar Yaduka's OpenPrinting-Image-Evaluation
//     framework (2-clause BSD licence). The comparator runs as a Python 3
//     subprocess and analyses the images across 14 metrics including SSIM,
//     PSNR, edge quality, colour accuracy, histogram similarity, and feature
//     matching. It returns an overall_quality score between 0.0 (completely
//     different) and 1.0 (pixel-perfect match). A score ≥ 0.95 is a PASS.
//     The Python script is embedded into the mfp-test binary at compile time
//     using //go:embed, so no external file path is required.
//
// # Test Matrix
//
// mfp-test queries the virtual printer's IPP attributes to discover all
// supported values for sides, print-color-mode, and document-format. It
// then generates test configurations from their combinations.
//
// Three run modes are available:
//
//   - --batch  runs every combination of sides × colour mode × format.
//     This is the exhaustive matrix and may take several minutes on a
//     printer that advertises many formats.
//
//   - --quick  runs all sides × all colour modes but only the first
//     document format. This covers duplex and colour correctness without
//     the overhead of testing every format variant.
//
//   - --single NAME  runs exactly one configuration identified by its
//     name (sides/colour-mode/format, e.g.
//     "one-sided/color/application/pdf"). Use this to reproduce a specific
//     failure after a batch run.
//
// Use --list or --list-quick to print all configuration names without
// running any tests.
//
// # Usage
//
//	mfp-test -m MODEL [options] --batch | --quick | --single NAME
//
// Required flag:
//
//	-m, --model FILE      Printer model file (Python script defining
//	                      ipp.attrs, and optionally escl.caps / wsd.caps)
//
// Run mode (exactly one required):
//
//	    --batch           Run the full test matrix
//	    --quick           Run the reduced (quick) test matrix
//	    --single NAME     Run one configuration by name
//	    --list            List all batch configurations and exit
//	    --list-quick      List all quick configurations and exit
//
// Options:
//
//	-P, --port PORT       IPP server TCP port (default: OS-assigned free port)
//	-n, --name NAME       CUPS queue name (default: mfp-test-<printer-model>)
//	-o, --output FILE     Write JSON report to FILE after all tests finish
//	    --threshold SCORE Minimum overall_quality score to pass (default 0.95)
//	    --timeout DUR     Document capture timeout per test (default 30s)
//	    --keep            Save the raw captured document to the current
//	                      directory (useful for debugging format issues)
//	-v, --verbose         Print per-metric similarity scores for every test,
//	                      not only for failures
//
// # System Dependencies
//
// The following packages must be installed on the test machine:
//
//   - CUPS           — queue management and job processing
//                      (sudo apt install cups)
//
// mfp-test calls lpadmin to create and remove the temporary CUPS queue.
// lpadmin requires CUPS administrative privileges. On most systems the
// invoking user must be a member of the "lpadmin" group (Debian/Ubuntu) or
// the "sys" group (Fedora/RHEL):
//
//	sudo usermod -aG lpadmin $USER   # Debian/Ubuntu
//	sudo usermod -aG sys    $USER   # Fedora/RHEL
//
// Alternatively, run mfp-test as root (not recommended in production).
//   - Ghostscript    — PostScript → PNG conversion
//                      (sudo apt install ghostscript)
//   - libvips        — PDF / JPEG / TIFF / WEBP / GIF conversion
//                      (sudo apt install libvips-dev)
//   - Python 3       — image evaluation subprocess
//                      with numpy, opencv-python, scikit-image, scipy
//                      (pip3 install numpy opencv-python scikit-image scipy)
//
// # Output
//
// After all tests complete, mfp-test prints a summary table to stdout:
//
//	TEST                                                SCORE  RESULT
//	----------------------------------------------------------------------
//	one-sided/color/application/pdf                   0.9821  PASS
//	one-sided/monochrome/image/pwg-raster              0.9134  PASS
//	two-sided-long-edge/color/application/pdf         0.8741  FAIL
//	----------------------------------------------------------------------
//	TOTAL                                              2/3 passed
//
// When a test fails (or --verbose is set), per-metric scores are printed
// below the table so the specific cause of the quality loss can be identified:
//
//	details for two-sided-long-edge/color/application/pdf:
//	  ssim_score                               0.8120
//	  edge_quality_score                       0.7340
//	  color_accuracy_score                     0.9210
//	  ...
//
// A machine-readable JSON report can be written with --output report.json.
// The report contains the same information as the table plus the generation
// timestamp and is suitable for processing in CI pipelines.
//
// JSON report schema:
//
//	{
//	  "generated_at": "<RFC3339 timestamp>",
//	  "total":        <integer — number of test configurations run>,
//	  "passed":       <integer — number that scored >= threshold>,
//	  "failed":       <integer — number that scored < threshold>,
//	  "results": [
//	    {
//	      "name":    "<sides/color-mode/format>",
//	      "passed":  <bool>,
//	      "score":   <float64 0.0–1.0>,
//	      "details": {            // omitted when evaluator is disabled
//	        "<metric>": <float64>,
//	        ...
//	      }
//	    },
//	    ...
//	  ]
//	}
//
// # Printer Model Files
//
// A printer model file is a Python script that sets the ipp.attrs dictionary
// to a dict of IPP printer attribute values. Example model files are provided
// in the modeling/examples/ directory of the go-mfp repository, including
// Xerox-B235.py (a real monochrome laser printer) and
// OpenPrinting-IPP-scan.py (a scanner-capable device).
//
// # Limitations
//
//   - image/vnd.cups-raster (legacy CUPS raster) and image/jpeg+gzip are not
//     supported for image evaluation. The captured bytes are saved with --keep
//     but the test is reported as a conversion error rather than a quality score.
//
//   - The image evaluator requires Python 3 and several scientific computing
//     packages. If Python 3 is unavailable, evaluation is skipped and all
//     tests are reported as passed with score 1.0 (capture-only mode).
package main

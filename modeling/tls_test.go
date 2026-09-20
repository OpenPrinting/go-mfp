// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// TLS certificates tests

package modeling

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/OpenPrinting/go-mfp/util/generic"
	"github.com/OpenPrinting/go-mfp/util/uuid"
)

// TestTLSPredictableCert tests the tlsPredictableCert function
func TestTLSPredictableCert(t *testing.T) {
	const rolls = 5

	// Test that different seeds yields different certs
	fingerprints := generic.NewSet[string]()
	for i := 0; i < rolls; i++ {
		seed := ([]byte)(uuid.Random().String())
		cert := tlsPredictableCert(seed)

		if len(cert.Certificate) != 1 {
			t.Fatalf("self-sinned cert can't contain %d certs in chain",
				len(cert.Certificate))
		}

		fingerpint := fmt.Sprintf("%02X",
			sha256.Sum256(cert.Certificate[0]))
		fingerprints.Add(fingerpint)
	}

	if fingerprints.Count() != rolls {
		t.Errorf("%d different seeds yields only %d different certs",
			rolls, fingerprints.Count())
	}

	// Test that the same seed yields the same cert
	for i := 0; i < rolls; i++ {
		seed := ([]byte)(uuid.Random().String())
		cert1 := tlsPredictableCert(seed)
		cert2 := tlsPredictableCert(seed)

		if len(cert1.Certificate) != 1 {
			t.Fatalf("self-sinned cert can't contain %d certs in chain",
				len(cert1.Certificate))
		}

		if len(cert2.Certificate) != 1 {
			t.Fatalf("self-sinned cert can't contain %d certs in chain",
				len(cert2.Certificate))
		}

		if !bytes.Equal(cert1.Certificate[0], cert2.Certificate[0]) {
			t.Fatalf("certificates are not predictable")
		}
	}
}

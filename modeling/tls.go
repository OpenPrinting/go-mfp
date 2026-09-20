// MFP - Miulti-Function Printers and scanners toolkit
// Printer and scanner modeling.
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// TLS certificates

package modeling

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// GetTLSCertificate returns TLS certificate associated with
// the model.
func (model *Model) GetTLSCertificate() tls.Certificate {
	// CUPS uses "Trust On the First Use" approach to the
	// printer certificates, so our certificate must be
	// predictable.
	//
	// Currently we seed the certificate generator with some
	// arbitrary random UUID. So all our certificates are
	// the same.
	//
	// TODO: consider binding certificates to the model
	// (for example, derive random seed from the IPP UUID).
	return tlsPredictableCert([]byte("a6007f1b-be59-404d-a849-34046e8bae0c"))
}

// tlsPredictableCert creates a deterministic self-signed TLS
// certificate based on the input seed.  For the same input bytes, it
// will always return the exact same tls.Certificate.
//
// WARNING: SECURITY RISK!
//
// This function is NOT secure and must ONLY be used for integration
// testing purposes. Because the certificate and private
// key are fully deterministic and bound to the input seed, they lack
// any cryptographic randomness. NEVER use this in production or any
// context requiring actual cryptographic protection, secrecy, or trust.
func tlsPredictableCert(seed []byte) tls.Certificate {
	// --- 1. Deterministic Private Key Generation (ECDSA P-256) ---
	curve := elliptic.P256()
	params := curve.Params()

	// Hash the seed to ensure uniform byte distribution for
	// the private scalar D and sufficient seed size (32 bytes).
	hashed := sha256.Sum256(seed[:])
	d := new(big.Int).SetBytes(hashed[:])
	d.Mod(d, params.N)
	if d.Sign() == 0 {
		d.SetInt64(1)
	}

	// Compute the public key point
	x, y := curve.ScalarBaseMult(d.Bytes())
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y},
		D:         d,
	}

	// --- 2. Create X.509 Template (Hardware Profile) ---
	// Generate a unique, yet fixed Serial Number based on the seed.
	// Make sure the serial number is positive.
	serialNumber := new(big.Int).SetBytes(hashed[:16])
	serialNumber.Add(serialNumber, big.NewInt(1))

	// Hardcode fixed lifespans. 100 years should be enough
	notBefore := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	notAfter := notBefore.Add(100 * 365 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "go-mfp test certificate",
			Organization: []string{"OpenPrinting"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// --- 3. Sign the Certificate ---
	derBytes, err := x509.CreateCertificate(nil,
		&template, &template, &privKey.PublicKey, privKey)
	assert.NoError(err)

	// --- 4. Return Ready tls.Certificate ---
	return tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privKey,
	}
}

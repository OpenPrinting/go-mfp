// MFP - Miulti-Function Printers and scanners toolkit
// UUID mini library
//
// Copyright (C) 2024 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// UUID type and functions.

package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"

	"github.com/OpenPrinting/go-mfp/internal/assert"
)

// Predefined UUID values:
var (
	NilUUID = UUID{}
	MaxUUID = MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
)

// Well-known namespaces:
var (
	NameSpaceDNS  = MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // [RFC9499]
	NameSpaceURL  = MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8") // [RFC1738]
	NameSpaceOID  = MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8") // [X660]
	NameSpaceX500 = MustParse("6ba7b814-9dad-11d1-80b4-00c04fd430c8") // [X500]
)

// UUID represents a parsed UUID. This type is comparable and can be
// used as the map key.
type UUID [16]byte

// Parse is the universal UUID parser.
//
// It recognizes at least the following UUID formats:
//   - urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//   - uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//   - xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//   - {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
//   - xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//
// It is very forgiving and should not be used as a validating
// parser.
func Parse(s string) (UUID, error) {
	// Strip decrations
	switch {
	case strings.HasPrefix(s, "urn:uuid:"):
		s = s[9:]
	case strings.HasPrefix(s, "uuid:"):
		s = s[5:]
	case strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}"):
		s = s[1 : len(s)-1]
	}

	// Now parse
	var uuid UUID
	cnt := 0
	for _, c := range s {
		var v rune

		// Parse next character. Ignore dashes ('-')
		switch {
		case c == '-':
			// Ignore dash characters
			continue

		case '0' <= c && c <= '9':
			v = c - '0'

		case 'a' <= c && c <= 'f':
			v = c - 'a' + 10

		case 'A' <= c && c <= 'F':
			v = c - 'A' + 10

		default:
			err := fmt.Errorf(
				"UUID contains invalid character: %q",
				string(c))
			return NilUUID, err
		}

		if cnt < 32 {
			if cnt&1 == 0 {
				uuid[cnt/2] = byte(v << 4)
			} else {
				uuid[cnt/2] |= byte(v)
			}
		}

		cnt++
	}

	// We must have exactly 32 digits
	switch {
	case cnt < 32:
		err := fmt.Errorf("UUID is too short (%d digits)", cnt)
		return NilUUID, err
	case cnt > 32:
		err := fmt.Errorf("UUID is too long (%d digits)", cnt)
		return NilUUID, err
	}

	return uuid, nil
}

// Random generates a random UUID.
// It uses [rand.Reader] as the source of entropy.
func Random() UUID {
	uuid, err := RandomFrom(rand.Reader)
	assert.NoError(err)
	return uuid
}

// RandomFrom generates a random UUID.
// It uses provided [io.Reader] as a source of entropy.
func RandomFrom(reader io.Reader) (UUID, error) {
	var uuid UUID

	_, err := io.ReadFull(reader, uuid[:])
	if err != nil {
		return NilUUID, err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4 (VersionRandom)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // VariantRFC4122

	return uuid, nil
}

// MD5 generates a new Name-Based [UUID], using MD5 crypto hash
// function.
//
// Name-Based UUIDs are, essentially, the crypto hash of supplied
// parameters, encoded as UUID. Parameters are:
//   - The namespace, supplied as UUID
//   - Some "name"
//
// It has the following properties:
//   - For the same combination of the space and name, it always
//     yield the same output
//   - If two different UUIDs are generated, using different names
//     and/or namespaces, these two UUIDs will be different with
//     very high probability.
//
// There is a number of well-known namespaces exists ([NameSpaceDNS] etc).
//
// See [RFC 4122, 4.3.] for details.
//
// [RFC 4122, 4.3.]: https://www.rfc-editor.org/rfc/rfc4122.html#section-4.3
func MD5(space UUID, name string) UUID {
	// Merge space and name
	data := make([]byte, len(space)+len(name))
	copy(data, space[:])
	copy(data[len(space):], ([]byte)(name))

	// Compute a sum
	sum := md5.Sum(data)

	// Make UUID
	var uuid UUID
	copy(uuid[:], sum[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x30 // Version 3 (VersionNameBasedMD5)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // VariantRFC4122

	return uuid
}

// SHA1 generates a new Name-Based [UUID], using SHA1 crypto hash
// function.
//
// Name-Based UUIDs are, essentially, the crypto hash of supplied
// parameters, encoded as UUID. Parameters are:
//   - The namespace, supplied as UUID
//   - Some "name"
//
// It has the following properties:
//   - For the same combination of the space and name, it always
//     yield the same output
//   - If two different UUIDs are generated, using different names
//     and/or namespaces, these two UUIDs will be different with
//     very high probability.
//
// There is a number of well-known namespaces exists ([NameSpaceDNS] etc).
//
// See [RFC 4122, 4.3.] for details.
//
// [RFC 4122, 4.3.]: https://www.rfc-editor.org/rfc/rfc4122.html#section-4.3
func SHA1(space UUID, name string) UUID {
	// Merge space and name
	data := make([]byte, len(space)+len(name))
	copy(data, space[:])
	copy(data[len(space):], ([]byte)(name))

	// Compute a sum
	sum := sha1.Sum(data)

	// Make UUID
	var uuid UUID
	copy(uuid[:], sum[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x50 // Version 5 (VersionNameBasedSHA1)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // VariantRFC4122

	return uuid
}

// MustParse calls [Parse] and panics in a case of any error.
func MustParse(s string) UUID {
	uuid, err := Parse(s)
	assert.NoError(err)
	return uuid
}

// Compare returns an integer comparing two UUIDs lexicographically.
//
// The result will be 0 id uuid == uuid2, -1 if uuid < uuid2,
// and +1 if uuid > uuid2
func (uuid UUID) Compare(uuid2 UUID) int {
	return bytes.Compare(uuid[:], uuid2[:])
}

// Less returns true if uuid is less that uuid2 in the lexicographical
// order.
func (uuid UUID) Less(uuid2 UUID) bool {
	return uuid.Compare(uuid2) < 0
}

// String returns the string form of UUID:
//
//	xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (uuid UUID) String() string {
	const format = "" +
		"%.2x%.2x%.2x%.2x-" +
		"%.2x%.2x-" +
		"%.2x%.2x-" +
		"%.2x%.2x-" +
		"%.2x%.2x%.2x%.2x%.2x%.2x"

	return fmt.Sprintf(
		format,
		uuid[0], uuid[1], uuid[2], uuid[3],
		uuid[4], uuid[5], uuid[6], uuid[7],
		uuid[8], uuid[9], uuid[10], uuid[11],
		uuid[12], uuid[13], uuid[14], uuid[15])
}

// GoString returns the Go source representation for UUID
func (uuid UUID) GoString() string {
	return fmt.Sprintf("uuid.MustParse(%q)", uuid.String())
}

// URN returns the URN form of UUID, per [RFC 2141]:
//
//	urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
//
// [RFC 2141]: https://www.rfc-editor.org/rfc/rfc2141.html
func (uuid UUID) URN() string {
	return "urn:uuid:" + uuid.String()
}

// MarshalText implements [encoding.TextMarshaler] interface for UUID.
func (uuid UUID) MarshalText() ([]byte, error) {
	return []byte(uuid.String()), nil
}

// Microsoft returns the Microsoft style form of UUID:
//
//	{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
func (uuid UUID) Microsoft() string {
	return "{" + uuid.String() + "}"
}

// Variant returns the [Variant] encoded in uuid.
func (uuid UUID) Variant() Variant {
	switch uuid[8] & 0b_111_00000 {
	case 0b_100_00000, 0b_101_00000:
		return VariantRFC4122
	case 0b_110_00000:
		return VariantMicrosoft
	case 0b_111_00000:
		return VariantFuture
	default: // 0b_0xx_00000
		return VariantNCS
	}
}

// Version returns the [Version] of uuid.
func (uuid UUID) Version() Version {
	return Version(uuid[6] >> 4)
}

// This package provides immutable UUID structs and the functions
// NewV3, NewV4, NewV5 and Parse() for generating versions 3, 4
// and 5 UUIDs as specified in RFC 4122.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"
)

// The UUID reserved variants. 
const (
	ReservedNCS       byte = 0x80
	ReservedRFC4122   byte = 0x40
	ReservedMicrosoft byte = 0x20
	ReservedFuture    byte = 0x00
	urnUuidPrefix	string = "urn:uuid:"
	dash			byte = byte('-')
)

// The following standard UUIDs are for use with NewV3() or NewV5().
var (
	NamespaceDNS, _  = ParseHex("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = ParseHex("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = ParseHex("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

// A UUID representation copmliant with specification in
// RFC 4122 document.
type UUID [16]byte

// ParseHex creates a UUID object from given hex string
// representation. Function accepts UUID string in following
// formats:
//
//     uuid.ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//     uuid.ParseHex("{6ba7b814-9dad-11d1-80b4-00c04fd430c8}")
//     uuid.ParseHex("urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//
func ParseHex(s string) (u *UUID, err error) {
	if strings.HasPrefix(s, urnUuidPrefix) {
		s = s[len(urnUuidPrefix):]
	}
	if strings.HasPrefix(s, "{") {
		if !strings.HasSuffix(s, "}") {
			err = errors.New("Invalid UUID string has an opening '{' bracket but no closing '}' bracket")
			return
		} else if len(s) != 38 {
			err = errors.New("Invalid UUID string")
			return
		}
		s = s[1:37]
	} else if strings.HasSuffix(s, "}") {
		err = errors.New("Invalid UUID string has no opening '{' bracket but does have a closing '}' bracket")
		return
	} else if (len(s) != 36) {
		err = errors.New("Invalid UUID string")
		return
	}
	var a, b byte
	var half bool
	u = new(UUID)
	v := u[0:0]
	for i, c := range []byte(s) {
		if (i == 8 || i == 13 || i == 18 || i == 23) {
			if (c != dash) {
				return nil, errors.New("Invalid UUID string had an improper character where '-' expected")
			}
			continue
		}
		switch {
		case '0' <= c && c <= '9':
			b = c - '0'
		case 'a' <= c && c <= 'f':
			b = c - 'a' + 10
		default:
			return nil, hex.InvalidByteError(i)
		}
		if half {
			v = append(v, (a << 4) | b)
		} else {
			a = b
		}
		half = !half
	}
	return
}

// Parse creates a UUID object from given bytes slice.
func Parse(b []byte) (u *UUID, err error) {
	if len(b) != 16 {
		err = errors.New("Given slice is not valid UUID sequence")
		return
	}
	u = new(UUID)
	copy(u[:], b)
	return
}

// Generate a UUID based on the MD5 hash of a namespace identifier
// and a name.
func NewV3(ns *UUID, name []byte) (u *UUID, err error) {
	if ns == nil {
		err = errors.New("Invalid namespace UUID")
		return
	}
	u = new(UUID)
	// Set all bits to MD5 hash generated from namespace and name.
	u.setBytesFromHash(md5.New(), ns[:], name)
	u.setVariant(ReservedRFC4122)
	u.setVersion(3)
	return
}

// Generate a random UUID.
func NewV4() (u *UUID, err error) {
	u = new(UUID)
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	_, err = rand.Read(u[:])
	if err != nil {
		return
	}
	u.setVariant(ReservedRFC4122)
	u.setVersion(4)
	return
}

// Generate a UUID based on the SHA-1 hash of a namespace identifier
// and a name.
func NewV5(ns *UUID, name []byte) (u *UUID, err error) {
	u = new(UUID)
	// Set all bits to truncated SHA1 hash generated from namespace
	// and name.
	u.setBytesFromHash(sha1.New(), ns[:], name)
	u.setVariant(ReservedRFC4122)
	u.setVersion(5)
	return
}

// Generate a MD5 hash of a namespace and a name, and copy it to the
// UUID slice.
func (u *UUID) setBytesFromHash(hash hash.Hash, ns, name []byte) {
	hash.Write(ns[:])
	hash.Write(name)
	copy(u[:], hash.Sum([]byte{})[:16])
}

// Set the two most significant bits (bits 6 and 7) of the
// clock_seq_hi_and_reserved to zero and one, respectively.
func (u *UUID) setVariant(v byte) {
	switch v {
	case ReservedNCS:
		u[8] = (u[8] | ReservedNCS) & 0xBF
	case ReservedRFC4122:
		u[8] = (u[8] | ReservedRFC4122) & 0x7F
	case ReservedMicrosoft:
		u[8] = (u[8] | ReservedMicrosoft) & 0x3F
	}
}

// Variant returns the UUID Variant, which determines the internal
// layout of the UUID. This will be one of the constants: RESERVED_NCS,
// RFC_4122, RESERVED_MICROSOFT, RESERVED_FUTURE.
func (u *UUID) Variant() byte {
	if u[8]&ReservedNCS == ReservedNCS {
		return ReservedNCS
	} else if u[8]&ReservedRFC4122 == ReservedRFC4122 {
		return ReservedRFC4122
	} else if u[8]&ReservedMicrosoft == ReservedMicrosoft {
		return ReservedMicrosoft
	}
	return ReservedFuture
}

// Set the four most significant bits (bits 12 through 15) of the
// time_hi_and_version field to the 4-bit version number.
func (u *UUID) setVersion(v byte) {
	u[6] = (u[6] & 0xF) | (v << 4)
}

// Version returns a version number of the algorithm used to
// generate the UUID sequence.
func (u *UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Returns unparsed version of the generated UUID sequence.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

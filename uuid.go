// This package provides immutable UUID structs and the functions
// NewV3, NewV4, NewV5 and Parse() for generating versions 3, 4
// and 5 UUIDs as specified in RFC 4122.
//
// If all you want is a unique ID, you should probably call NewV4
// which creates a random UUID.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

import (
	"crypto/rand"
	"fmt"
)

// The UUID reserved variants. 
const (
	ReservedNCS       byte = 0x80
	ReservedRFC4122   byte = 0x40
	ReservedMicrosoft byte = 0x20
	ReservedFuture    byte = 0x00
)

// The following standard UUIDs are for use with NewV3() or NewV5().
var (
	NamespaceDNS, _  = ParseString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = ParseString("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = ParseString("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = ParseString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")	
)

type UUID [16]byte

// ParseString creates a UUID object from given string.
func ParseString(s string) (u *UUID, err error) {
	return
}

// Parse creates a UUID object from given bytes slice.
func Parse(b []byte) (u *UUID, err error) {
	return
}

// Generate a UUID based on the MD5 hash of a namespace identifier
// and a name.
func NewV3(ns *UUID, name string) (u *UUID, err error) {
	u = new(UUID)
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
func NewV5(ns *UUID, name string) (u *UUID, err error) {
	u = new(UUID)
	u.setVariant(ReservedRFC4122)
	u.setVersion(5)
	return
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
	if u[8] & ReservedNCS == ReservedNCS {
		return ReservedNCS
	} else if u[8] & ReservedRFC4122 == ReservedRFC4122 {
		return ReservedRFC4122
	} else if u[8] & ReservedMicrosoft == ReservedMicrosoft {
		return ReservedMicrosoft
	} 
	return ReservedFuture
}

// Set the four most significant bits (bits 12 through 15) of the
// time_hi_and_version field to the 4-bit version number.
func (u *UUID) setVersion(v byte) {
	u[6] = (u[6] & 0xF) | 0x40
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

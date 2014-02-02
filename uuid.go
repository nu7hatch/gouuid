// This package provides immutable UUID structs and the functions
// NewV3, NewV4, NewV5, New([]byte) and Parse(string) for generating versions 3, 4
// and 5 UUIDs as specified in RFC 4122.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package gouuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"regexp"
	"io"
)

// The UUID reserved variants.
const (
	ReservedNCS       byte = 0x00
	ReservedRFC4122   byte = 0x80
	ReservedMicrosoft byte = 0xC0
	ReservedFuture    byte = 0xE0
)

const (
	variantIndex = 8
	versionIndex = 6
	length       = 16
)

// The following standard UUIDs are for use with NewV3() or NewV5().
var (
	NamespaceDNS, NamespaceURL, NamespaceOID, NamespaceX500 UUID
	reg *regexp.Regexp
)

// Pattern used to parse hex string representation of the UUID.
// FIXME: do something to consider both brackets at one time,
// current one allows to parse string with only one opening
// or closing bracket.
const hexPattern = `^(urn\:uuid\:)?\{?(\w{8})-(\w{4})-([1-5]\w{3})-(\w{4})-(\w{12})\}?$`

func init() {
	NamespaceDNS = newHex("6ba7b8109dad11d180b400c04fd430c8")
	NamespaceURL = newHex("6ba7b8119dad11d180b400c04fd430c8")
	NamespaceOID = newHex("6ba7b8129dad11d180b400c04fd430c8")
	NamespaceX500 = newHex("6ba7b8149dad11d180b400c04fd430c8")
	reg = regexp.MustCompile(hexPattern)
}

// A UUID representation compliant with specification in
// RFC 4122 document.
type UUID [length]byte

// Set the three most significant bits (bits 0, 1 and 2) of the
// clock_seq_hi_and_reserved to variant mask v.
func (u *UUID) setVariant(v byte) {
	switch v {
	case ReservedRFC4122:
		u[variantIndex] &= 0x3F
	case ReservedFuture, ReservedMicrosoft:
		u[variantIndex] &= 0x1F
	case ReservedNCS:
		u[variantIndex] &= 0x7F
	default:
		panic(errors.New("UUID.setVariant: invalid variant mask"))
	}
	u[variantIndex] |= v
}

// Variant returns the UUID Variant, which determines the internal
// layout of the UUID. This will be one of the constants: RESERVED_NCS,
// RFC_4122, RESERVED_MICROSOFT, RESERVED_FUTURE.
func (u *UUID) Variant() byte {
	switch u[variantIndex] & 0xE0 {
	case ReservedRFC4122, 0xA0:
		return ReservedRFC4122
	case ReservedMicrosoft:
		return ReservedMicrosoft
	case ReservedFuture:
		return ReservedFuture
	}
	return ReservedNCS
}

// Set the four most significant bits (bits 0 through 3) of the
// time_hi_and_version field to the 4-bit version number.
func (u *UUID) setVersion(v byte) {
	u[versionIndex] &= 0x0F
	u[versionIndex] |= v<<4
}

// Version returns a version number of the algorithm used to
// generate the UUID sequence.
func (u *UUID) Version() int {
	return int(u[versionIndex]>>4)
}

// Returns unparsed version of the generated UUID sequence.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// UnmarshalBinary copies data into an existing UUID.
// This conforms to the BinaryUnmarshaller interface
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != length {
		err = errors.New("UUID.UnmarshalBinary: data invalid length")
		return
	}
	copy(u[:], data)
	return
}

// MarshalBinary returns a slice of UUID array
// This conforms to the BinaryMarshaller interface
func (u *UUID) MarshalBinary() (data []byte, err error) {
	return u[:], nil
}

// Parse creates a UUID object from given hex string
// representation. Function accepts UUID string in following
// formats:
//
//     gouuid.Parse("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//     gouuid.Parse("{6ba7b814-9dad-11d1-80b4-00c04fd430c8}")
//     gouuid.Parse("urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//
func Parse(s string) (u *UUID, err error) {
	md := reg.FindStringSubmatch(s)
	if md == nil {
		err = errors.New("Parse: invalid UUID string")
		return
	}
	return NewHex(md[2] + md[3] + md[4] + md[5] + md[6])
}

// Unsafe does not check string format
func NewHex(uuid string) (u *UUID, err error) {
	bytes, err := hex.DecodeString(uuid)
	if err != nil {
		return
	}
	return New(bytes)
}

// Unsafe does not check string format
func newHex(uuid string) UUID {
	u, err := NewHex(uuid)
	if err != nil {
		panic(err)
	}
	return *u
}

// New creates a UUID object from data byte slice.
func New(data []byte) (u *UUID, err error) {
	u = new(UUID)
	err = u.UnmarshalBinary(data)
	return
}

// Generate a UUID based on the MD5 hash of a namespace identifier
// and a name.
func NewV3(ns UUID, name string) (u *UUID, err error) {
	// Set all bits to MD5 hash generated from namespace and name.
	u, err = New(sum(md5.New(), ns, name))
	if err != nil {
		return
	}
	u.setVariant(ReservedRFC4122)
	u.setVersion(3)
	return
}

// Generate a random UUID.
func NewV4() (u *UUID, err error) {
	u = new(UUID)
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	_, err = rand.Read(u[:length]) // explicit length to ensure proper slice
	if err != nil {
		panic(err)
	}
	u.setVariant(ReservedRFC4122)
	u.setVersion(4)
	return
}

// Generate a UUID based on the SHA-1 hash of a namespace identifier
// and a name.
func NewV5(ns UUID, name string) (u *UUID, err error) {
	// Set all bits to truncated SHA1 hash generated from namespace
	// and name.
	u, err = New(sum(sha1.New(), ns, name))
	if err != nil {
		return
	}
	u.setVariant(ReservedRFC4122)
	u.setVersion(5)
	return
}

// Generate a check sum hash of a namespace and a name, and copy it to the
// UUID slice.
func sum(h hash.Hash, ns UUID, name string) []byte {
	h.Write(ns[:])
	io.WriteString(h, name)
	return h.Sum(nil)[:length]
}

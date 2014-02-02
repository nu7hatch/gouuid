// This package provides immutable UUID structs and the functions
// NewV3, NewV4, NewV5 and Parse() for generating versions 3, 4
// and 5 UUIDs as specified in RFC 4122.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package gouuid

import (
	"regexp"
	"testing"
	"fmt"
)

const format = "^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$"

var (
	uuidBytes    = []byte{
	0xAA, 0xCF, 0xEE, 0x12,
	0xD4, 0x00,
	0x67, 0x23,
	0x00,
	0xD3,
	0x23, 0x12, 0x4A, 0x11, 0x89, 0xFF,
}
	uuidVariants = []byte{
	ReservedNCS, ReservedRFC4122, ReservedMicrosoft, ReservedFuture,
}
	printer      = false
)

func TestVarientBits(t *testing.T) {

	u := new(UUID)
	for _, v := range uuidVariants {
		for i := 0; i <= 255; i ++ {
			uuidBytes[variantIndex] = byte(i)
			copy(u[:], uuidBytes)
			u.setVariant(v)
			b := u[variantIndex]>>4
			t_VariantConstraint(v, b, u, t)
			if u.Variant() != v {
				t.Errorf("%d does not resolve to %x", i, v)
			}
		}
	}
}

func t_VariantConstraint(v, b byte, u *UUID, t *testing.T) {
	output(u)
	switch v {
	case ReservedNCS:
		switch b {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
			outputf(": %X ", b)
			break
		default: t.Errorf("%X most high bits do not resolve to 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07", b)
		}
	case ReservedRFC4122:
		switch b {
		case 0x08, 0x09, 0x0A, 0x0B:
			outputf(": %X ", b)
			break
		default: t.Errorf("%X most high bits do not resolve to 0x08, 0x09, 0x0A, 0x0B", b)
		}
	case ReservedMicrosoft:
		switch b {
		case 0x0C, 0x0D:
			outputf(": %X ", b)
			break
		default: t.Errorf("%X most high bits do not resolve to 0x0C, 0x0D", b)
		}
	case ReservedFuture:
		switch b {
		case 0x0E, 0x0F:
			outputf(": %X ", b)
			break
		default: t.Errorf("%X most high bits do not resolve to 0x0E, 0x0F", b)
		}
	}
	output("\n")
}

func TestUnmarshalBinary(t *testing.T) {
	u := new(UUID)
	err := u.UnmarshalBinary([]byte{1, 2, 3, 4, 5})
	if err == nil {
		t.Errorf("Expected error due to invalid bte length")
	}

	err = u.UnmarshalBinary(uuidBytes)
	if err != nil {
		t.Errorf("Expected to parse UUID sequence without problems")
		return
	}
	if u.String() != u.String() {
		t.Errorf("Expected parsed UUID to be the same as base, %s != %s", u.String(), u.String())
	}
}

func TestUUID_VersionBits(t *testing.T) {
	u := new(UUID)
	for v := 0; v < 16; v++ {
		for i := 0; i <= 255; i ++ {
			uuidBytes[versionIndex] = byte(i)
			t_VersionConstraints(v, u, t)
		}
	}
}

func t_VersionConstraints(v int, u *UUID, t *testing.T) {
	copy(u[:], uuidBytes)
	u.setVersion(byte(v))
	output(u)
	if u.Version() != v {
		t.Errorf("%x does not resolve to %x", byte(u.Version()), v)
	}
	output("\n")
}

func TestUUID_ParseString(t *testing.T) {
	_, err := ParseHex("foo")
	if err == nil {
		t.Errorf("Expected error due to invalid UUID string")
	}
	base, _ := NewV4()
	u, err := ParseHex(base.String())
	if err != nil {
		t.Errorf("Expected to parse UUID sequence without problems")
		return
	}
	if u.String() != base.String() {
		t.Errorf("Expected parsed UUID to be the same as base, %s != %s", u.String(), base.String())
	}
}

var (
	invalidHexStrings = [...]string{
	"foo",
	"6ba7b814-9dad-11d1-80b4-",
	"6ba7b8149dad-11d1-80b4-00c04fd430c8",
	"6ba7b814-9dad11d1-80b4-00c04fd430c8",
	"6ba7b814-9dad-11d180b4-00c04fd430c8",
	"6ba7b814-9dad-11d1-80b400c04fd430c8",
	"6ba7b8149dad11d180b400c04fd430c8",
	"6ba7b8147-9dad11d1-80b400c04fd430c8",
	"6ba7b814--9dad-11d1-80b4--00c04fd430c8",
	"6ba7b814-9dad7-11d1-80b4-00c04fd430c8999",
	"{6ba7b814-9dad-1180b4-00c04fd430c8",
	"{6ba7b814--11d1-80b4-00c04fd430c8}",
	"6ba7b8149dad-11d1-80b4-00c04fd430c8}",
	"{6ba7b8149dad-11d1-80b400c04fd430c8}",
	"{6ba7b814-9dad11d180b400c04fd430c8}",
	"urn:uuid:6ba7b814-9dad-1666666680b4-00c04fd430c8",
}
	validHexStrings    = [...]string{
	"6ba7b814-9dad-11d1-80b4-00c04fd430c8",
	"{6ba7b814-9dad-11d1-80b4-00c04fd430c8}",
	"{6ba7b814-9dad-11d1-80b4-00c04fd430c8",
	"6ba7b814-9dad-11d1-80b4-00c04fd430c8}",
	"urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8",
}
)

func TestUUID_ParseHex(t *testing.T) {
	for _, v := range invalidHexStrings {
		_, err := ParseHex(v)
		if err == nil {
			t.Errorf("Expected error due to invalid UUID string %s", v)
		}
	}
	for _, v := range validHexStrings {
		_, err := ParseHex(v)
		if err != nil {
			t.Errorf("Expected valid UUID string %s but got error", v)
		}
	}
}

func TestNewV3(t *testing.T) {
	u, err := NewV3(NamespaceURL, "golang.org")
	if err != nil {
		t.Errorf("Expected to generate UUID without problems, error thrown: %d", err.Error())
		return
	}
	if u.Version() != 3 {
		t.Errorf("Expected to generate UUIDv3, given %d", u.Version())
	}
	if u.Variant() != ReservedRFC4122 {
		t.Errorf("Expected to generate UUIDv3 RFC4122 variant, given %x", u.Variant())
	}
	re := regexp.MustCompile(format)
	if !re.MatchString(u.String()) {
		t.Errorf("Expected string representation to be valid, given %s", u.String())
	}
	u2, _ := NewV3(NamespaceURL, "golang.org")
	if u2.String() != u.String() {
		t.Errorf("Expected UUIDs generated of the same namespace and name to be the same")
	}
	u3, _ := NewV3(NamespaceDNS, "golang.org")
	if u3.String() == u.String() {
		t.Errorf("Expected UUIDs generated of different namespace and the same name to be different")
	}
	u4, _ := NewV3(NamespaceURL, "code.google.com")
	if u4.String() == u.String() {
		t.Errorf("Expected UUIDs generated of the same namespace and different names to be different")
	}
}

func TestNewV4(t *testing.T) {
	u, err := NewV4()
	if err != nil {
		t.Errorf("Expected to generate UUID without problems, error thrown: %s", err.Error())
		return
	}
	if u.Version() != 4 {
		t.Errorf("Expected to generate UUIDv4, given %d", u.Version())
	}
	if u.Variant() != ReservedRFC4122 {
		t.Errorf("Expected to generate UUIDv4 RFC4122 variant, given %x", u.Variant())
	}
	re := regexp.MustCompile(format)
	if !re.MatchString(u.String()) {
		t.Errorf("Expected string representation to be valid, given %s", u.String())
	}
}

func TestNewV5(t *testing.T) {
	u, err := NewV5(NamespaceURL, "golang.org")
	if err != nil {
		t.Errorf("Expected to generate UUID without problems, error thrown: %d", err.Error())
		return
	}
	if u.Version() != 5 {
		t.Errorf("Expected to generate UUIDv5, given %d", u.Version())
	}
	if u.Variant() != ReservedRFC4122 {
		t.Errorf("Expected to generate UUIDv5 RFC4122 variant, given %x", u.Variant())
	}
	re := regexp.MustCompile(format)
	if !re.MatchString(u.String()) {
		t.Errorf("Expected string representation to be valid, given %s", u.String())
	}
	u2, _ := NewV5(NamespaceURL, "golang.org")
	if u2.String() != u.String() {
		t.Errorf("Expected UUIDs generated of the same namespace and name to be the same")
	}
	u3, _ := NewV5(NamespaceDNS, "golang.org")
	if u3.String() == u.String() {
		t.Errorf("Expected UUIDs generated of different namespace and the same name to be different")
	}
	u4, _ := NewV5(NamespaceURL, "code.google.com")
	if u4.String() == u.String() {
		t.Errorf("Expected UUIDs generated of the same namespace and different names to be different")
	}
}

func BenchmarkParseHex(b *testing.B) {
	s := "f3593cff-ee92-40df-4086-87825b523f13"
	for i := 0; i < b.N; i++ {
		_, err := ParseHex(s)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.ReportAllocs()
}

func output(a interface{}) {
	if printer {
		fmt.Print(a)
	}
}

func outputf(format string, a... interface{}) {
	if printer {
		fmt.Printf(format, a)
	}
}

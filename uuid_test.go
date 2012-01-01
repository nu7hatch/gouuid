// This package contains a binding for the `libuuid` - 
// http://linux.die.net/man/3/libuuid.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

import (
	"regexp"
	"testing"
)

const v4Format = "^[a-z0-9]{8}-[a-z0-9]{4}-4[a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$"

func TestNewV3(t *testing.T) {
	//uuid := NewV3()
	//if uuid.Version() != 3 {
	//	t.Errorf("Expected to generate UUIDv3, given %d", uuid.Version())
	//}
}

func TestNewV4(t *testing.T) {
	uuid, err := NewV4()
	if err != nil {
		t.Errorf("Expected to generate UUID without problems, error thrown: %s", err.Error())
	}
	if uuid.Version() != 4 {
		t.Errorf("Expected to generate UUIDv4, given %d", uuid.Version())
	}
	if uuid.Variant() != ReservedRFC4122 {
		t.Errorf("Expected to generate UUIDv4 RFC4122 variant, given %x", uuid.Variant())
	}
	re := regexp.MustCompile(v4Format)
	if !re.MatchString(uuid.String()) {
		t.Errorf("Expected string representation to be valid, given %s", uuid.String())
	}
}

func TestNewV5(t *testing.T) {
	//uuid := NewV5()
	//if uuid.Version() != 5 {
	//	t.Errorf("Expected to generate UUIDv5, given %d", uuid.Version())
	//
}

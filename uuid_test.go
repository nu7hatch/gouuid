// This package contains a binding for the `libuuid` - 
// http://linux.die.net/man/3/libuuid.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

import (
	"testing"
	"regexp"
)

const pattern = "^[a-z0-9]{8}-([a-z0-9]{4}-){3}[a-z0-9]{12}$"

func TestGenerateRand(t *testing.T) {
	r, _ := regexp.Compile(pattern)
	uuid := GenerateRand()
	if !r.MatchString(uuid) {
		t.Errorf("Expected to generate correct UUID, given '%s'", uuid)
	}
}

func TestGenerateTime(t *testing.T) {
	r, _ := regexp.Compile(pattern)
	uuid := GenerateTime()
	if !r.MatchString(uuid) {
		t.Errorf("Expected to generate correct UUID, given '%s'", uuid)
	}
}
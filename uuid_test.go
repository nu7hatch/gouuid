// This package provides a hybrid of MQ and WebSockets server with
// support for horizontal scalability.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package webrocket

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
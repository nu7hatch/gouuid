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

//#include <stdlib.h>
//#include <uuid.h>
import "C"

const (
	uuidSize         = 16
	uuidUnparsedSize = 36
)

// GenerateTime generates an universally unique identifier based on
// current time. 
func GenerateTime() string {
	buf := C.malloc(C.size_t(uuidUnparsedSize))
	uuid := C.malloc(C.size_t(uuidSize))
	C.uuid_generate_time((*C.uchar)(uuid))
	C.uuid_unparse((*C.uchar)(uuid), (*C.char)(buf))
	return C.GoString((*C.char)(buf))
}

// GenerateRand generate an universally unique identifier based on
// randomly generated numbers.
func GenerateRand() (res string) {
	buf := C.malloc(C.size_t(uuidSize))
	uuid := C.malloc(C.size_t(uuidUnparsedSize))
	C.uuid_generate_random((*C.uchar)(uuid))
	C.uuid_unparse((*C.uchar)(uuid), (*C.char)(buf))
	return C.GoString((*C.char)(buf))
}

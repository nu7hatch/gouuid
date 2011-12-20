// This package contains a binding for the `libuuid` - 
// http://linux.die.net/man/3/libuuid.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

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
	defer C.free(buf)
	uuid := C.malloc(C.size_t(uuidSize))
	defer C.free(uuid)
	C.uuid_generate_time((*C.uchar)(uuid))
	C.uuid_unparse((*C.uchar)(uuid), (*C.char)(buf))
	return C.GoString((*C.char)(buf))
}

// GenerateRand generate an universally unique identifier based on
// randomly generated numbers.
func GenerateRand() (res string) {
	buf := C.malloc(C.size_t(uuidUnparsedSize))
	defer C.free(buf)
	uuid := C.malloc(C.size_t(uuidSize))
	defer C.free(uuid)
	C.uuid_generate_random((*C.uchar)(uuid))
	C.uuid_unparse((*C.uchar)(uuid), (*C.char)(buf))
	return C.GoString((*C.char)(buf))
}

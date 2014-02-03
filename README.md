# Pure Go UUID implementation

This package provides immutable UUID structs and the functions
NewV3, NewV4, NewV5, New and Parse() for generating versions 3, 4
and 5 UUIDs as specified in [RFC 4122](http://www.ietf.org/rfc/rfc4122.txt).

# Recent Changes
* varient bits and type is now set correctly
* varient bits and varient type can now be retrieved more efficiently
* new tests added for variant setting
* new tests added to confirm proper version setting
* type UUID now conforms to the BinaryMarshaller and BinaryUnmarshaller interfaces
* New was added to create a base UUID from a []byte
* ParseHex was renamed to simply Parse
* Parse creates a UUID and properly checks the string format
* NewHex now performs unsafe creation of UUID from a hex string
* NewV3 and NewV5 now take strings as a namespace name

## Installation

Use the `go` tool:

	$ go get github.com/nu7hatch/gouuid

## Usage

See [documentation and examples](http://godoc.org/github.com/nu7hatch/gouuid)
for more information.

## Copyright

Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>. See [COPYING](https://github.com/nu7hatch/gouuid/tree/master/COPYING)
file for details.

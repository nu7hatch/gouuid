include $(GOROOT)/src/Make.inc

CFLAGS=$(shell pkg-config --cflags uuid) -I.
LDFLAGS=$(shell pkg-config --libs uuid)

TARG=uuid
CGOFILES=uuid.go
CGO_CFLAGS=$(CFLAGS)
CGO_LDFLAGS=$(LDFLAGS)

include $(GOROOT)/src/Make.pkg
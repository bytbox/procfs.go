include ${GOROOT}/src/Make.inc

TARG = procd
GOFILES = procd.go

include ${GOROOT}/src/Make.cmd

fmt:
	gofmt -w *.go


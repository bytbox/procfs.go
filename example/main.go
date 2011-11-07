package main

import (
	"github.com/bytbox/procfs.go/procfs"
)

func main() {
	var pfs procfs.ProcFS
	pfs.Fill()
	println(pfs.Self)
}

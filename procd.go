package main

import (
	"github.com/bytbox/procfs.go/procfs"
	"json"
	"log"
)

func main() {
	var pfs procfs.ProcFS
	pfs.Fill()
	str, err := json.MarshalIndent(pfs, "", "  ")
	if err != nil {
		log.Fatal("ERR: ", err)
	}
	println(string(str))
}

package main

import (
	"flag"
	"log"
	"time"
	"sync"

	"github.com/bytbox/procfs.go/procfs"
)

var pfs *procfs.ProcFS
var pfsMutex sync.Mutex
var interval = flag.Int64("interval", 10, "update interval")

func main() {
	flag.Parse()

	go updater()

/*
	var pfs procfs.ProcFS
	pfs.Fill()
	str, err := json.MarshalIndent(pfs, "", "  ")
	if err != nil {
		log.Fatal("ERR: ", err)
	}
	println(string(str))
*/

	<-make(chan int)
}

// Maintain an updated version of ProcFS
func updater() {
	ticker := time.Tick(1e9 * *interval)
	for true {
		log.Print("Updating...")
		sn := time.Nanoseconds()
		var pfs2 procfs.ProcFS
		pfs2.Fill()
		pfsMutex.Lock()
		pfs = &pfs2
		pfsMutex.Unlock()
		en := time.Nanoseconds()
		log.Print("Done in ", (en-sn)/1000, "μs")
		<-ticker
	}
}

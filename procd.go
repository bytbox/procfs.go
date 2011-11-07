package main

import (
	"flag"
	"log"
	"net"
	"rpc"
	"rpc/jsonrpc"
	"sync"
	"time"

	"github.com/bytbox/procfs.go/procfs"
)

var pfs *procfs.ProcFS
var pfsMutex sync.Mutex
var interval = flag.Int64("interval", 10, "update interval")
var port = flag.String("port", ":16070", "JSON-RPC service port")

func main() {
	flag.Parse()

	go updater()
	go serve()

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

type ProcFSServer struct {}

func (ProcFSServer) Get(req string, reply *procfs.ProcFS) error {
	reply = pfs
	return nil
}

// Runs the JSON-RPC server
func serve() {
	server := ProcFSServer{}
	rpc.Register(server)

	l, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Print("WARN: ", err)
			continue
		}

		log.Print("Serving ", c.RemoteAddr())
		go jsonrpc.ServeConn(c)
	}
}

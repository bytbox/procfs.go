package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	"github.com/bytbox/procfs.go/procfs"
)

var pfs *procfs.ProcFS
var pfsMutex sync.Mutex
var interval = flag.Int64("interval", 10, "update interval")
var http_port = flag.String("http", ":6070", "HTTP service port")

func main() {
	flag.Parse()

	go updater()
	go serveHTTP()

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

func (ProcFSServer) Get(req string, reply **procfs.ProcFS) error {
	*reply = pfs
	return nil
}

func serveHTTP() {
	server := ProcFSServer{}
	rpc.RegisterName("ProcFS", server)
	http.HandleFunc("/", HTMLServer)
	http.HandleFunc("/rpc", RPCServer)

	err := http.ListenAndServe(*http_port, http.DefaultServeMux)
	if err != nil {
		log.Print("ERR: ", err)
	}
}

func HTMLServer(w http.ResponseWriter, req *http.Request) {
	c, err := ioutil.ReadFile("proc.html")
	if err != nil {
		log.Fatal("ERR: html/index.html not openable")
		return
	}
	w.Write(c)
}

type combo struct{
	a net.Conn
	b io.ReadWriter
}

func (c combo) Close() error {
	return c.a.Close()
}

func (c combo) Read(b []byte) (int, error) {
	return c.b.Read(b)
}

func (c combo) Write(b []byte) (int, error) {
	return c.a.Write(b)
}

func RPCServer(w http.ResponseWriter, req *http.Request) {
	h, buf, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("ERR: ", err)
	}
	connected := "200 Connected to JSON-RPC"
	io.WriteString(h, "HTTP/1.0 "+connected+"\n\n")
	log.Print("Serving")
	codec := jsonrpc.NewServerCodec(combo{h, buf})
	rpc.ServeRequest(codec)
	h.Close()
}


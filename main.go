package main

import (
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
)

//	constants
const ssePath = "/.hak/fs/sse"

var (
	watchDir *string
	portPtr  *int
	//go:embed localhost.pem
	pubKeyMaterial []byte
	//go:embed localhost-key.pem
	privKeyMaterial []byte
)

var events = make(chan []byte)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func main() {

	//	start watcher
	if err := watchRecursively(*watchDir, events); err != nil {
		log.Fatal(err)
	}

	//	despatch events to SSE broker
	broker := NewServer()
	go func() {
		for b := range events {
			broker.Notifier <- b
		}
	}()

	//	start web server
	cert, err := tls.X509KeyPair(pubKeyMaterial, privKeyMaterial)
	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(*watchDir))
	mux.Handle("/", injectHeadersForStaticFiles(fs))
	mux.Handle(ssePath, broker)
	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	fmt.Printf("listening on port %d\n", *portPtr)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)
	if err != nil {
		log.Fatalln(err)
	}

}

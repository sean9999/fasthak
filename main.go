package main

import (
	"crypto/tls"
	"embed"
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
	//go:embed frontend
	frontend embed.FS
)

var niceEvents = make(chan NiceEvent)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func main() {

	//	start watcher
	if err := watchRecursively(*watchDir, niceEvents); err != nil {
		log.Fatal(err)
	}

	//	despatch events to SSE sseBroker
	sseBroker := NewBroker()
	go func() {
		for b := range niceEvents {
			sseBroker.Notifier <- b
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

	//	.hak/**/*

	//	.hak/js/
	mux.Handle("./hak/js/", http.FileServer(frontend))

	//	./hak/fs/sse
	mux.Handle(ssePath, sseBroker)

	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	fmt.Printf("listening on port %d\n", *portPtr)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)
	if err != nil {
		log.Fatalln(err)
	}

}

package main

import (
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/rjeczalik/notify"
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
	c               = make(chan notify.EventInfo)
)

func init() {
	//	parse options and arguments
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func main() {

	//	start watcher

	if err := notify.Watch(*watchDir+"/...", c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)
	//defer close(c)

	//	start web server
	cert, err := tls.X509KeyPair(pubKeyMaterial, privKeyMaterial)
	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(*watchDir))
	mux.Handle("/", injectHeadersForStaticFiles(fs))
	mux.Handle(ssePath, handler(fsEventHandler))
	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	fmt.Printf("listening on port %d\n", *portPtr)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)
	if err != nil {
		log.Fatalln(err)
	}

}

func fsEventHandler(w http.ResponseWriter, r *http.Request) {

	//	consume fileSystem events and push them to the HTTP response one by one
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	for ei := range c {
		pushEvent(ei, w)
	}

}

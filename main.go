package main

import (
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
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

func init() {
	//	parse options and arguments
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func main() {

	//	start watcher
	if watcherBootstrapError != nil {
		log.Fatal(watcherBootstrapError)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(watcher)
	watcherAddFolderError := watcher.Add(*watchDir)
	if watcherAddFolderError != nil {
		log.Fatal(watcherAddFolderError)
	}
	fmt.Printf("watching folder %s\n", *watchDir)

	//	start web server
	//cert, err := GenX509KeyPair()
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

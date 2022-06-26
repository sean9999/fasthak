package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const hakPrefix = "/.hak"
const ssePath = "fs/sse"

var (
	watchDir *string
	portPtr  *int
	//go:embed frontend/*
	frontend embed.FS
)

type NiceEvent struct {
	Event string
	File  string
}

var niceEvents = make(chan NiceEvent)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func hakHandler() http.Handler {
	fsys := fs.FS(frontend)
	hakFiles, _ := fs.Sub(fsys, "frontend")
	return http.StripPrefix(hakPrefix+"/js/", http.FileServer(http.FS(hakFiles)))
}

func main() {

	//	start watcher
	if err := watchRecursively(*watchDir, niceEvents); err != nil {
		log.Fatal(err)
	}

	//	dispatch events to SSE sseBroker
	sseBroker := NewBroker()
	go func() {
		for b := range niceEvents {
			sseBroker.Notifier <- b
		}
	}()

	//	get TLS keys. panic if they don't exist
	pubKeyMaterial, err := ioutil.ReadFile("localhost.pem")
	if err != nil {
		panic("Could not find localhost.pem")
	}
	privKeyMaterial, err := ioutil.ReadFile("localhost-key.pem")
	if err != nil {
		panic("Could not find localhost-key.pem")
	}

	//	start web server
	cert, err := tls.X509KeyPair(pubKeyMaterial, privKeyMaterial)
	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()

	//	static files
	staticFileServer := http.FileServer(http.Dir(*watchDir))
	mux.Handle("/", injectHeadersForStaticFiles(staticFileServer))

	//	.hak/js/*
	mux.Handle(hakPrefix+"/js/", hakHandler())

	//	./hak/fs/sse
	mux.Handle(path.Join(hakPrefix, ssePath), sseBroker)

	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	fmt.Printf("listening on port %d\n", *portPtr)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)
	if err != nil {
		log.Fatalln(err)
	}

}

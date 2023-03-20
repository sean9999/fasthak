package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/sean9999/rebouncer"
)

// constants
const hakPrefix = "/.hak"
const ssePath = "fs/sse"

var (
	watchDir *string
	portPtr  *int
	privkey  *string
	pubkey   *string
	//go:embed frontend/*
	frontend embed.FS
)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	privkey = flag.String("privkey", "./localhost-key.pem", "private key in PEM format")
	pubkey = flag.String("pubkey", "./localhost.pem", "public key in PEM format")
	flag.Parse()
}

func hakHandler() http.Handler {
	fsys := fs.FS(frontend)
	hakFiles, _ := fs.Sub(fsys, "frontend")
	return http.StripPrefix(hakPrefix+"/js/", http.FileServer(http.FS(hakFiles)))
}

func main() {

	//	load up privkey and pubkey
	cert, err := tls.LoadX509KeyPair(*pubkey, *privkey)
	if err != nil {
		log.Fatal(err)
	}

	stateMachine := rebouncer.NewInotify(*watchDir, 1000)
	niceEvents := stateMachine.Subscribe()

	//	dispatch events to SSE sseBroker
	sseBroker := NewBroker()
	go func() {
		for b := range niceEvents {
			sseBroker.Notifier <- b
		}
	}()

	//	start web server
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

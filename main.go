package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gokyle/fswatch"
)

// constants
const hakPrefix = "/.hak"
const ssePath = "fs/sse"

var (
	watchDir     *string
	subDomainPtr *string
	portPtr      *int

	//go:embed certs/*
	secrets embed.FS

	//go:embed frontend/*
	frontend embed.FS
)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	subDomainPtr = flag.String("subdomain", "fasthak", "what subomain on rec.la to use")
	portPtr = flag.Int("port", 0, "what port to listen on")
	flag.Parse()
}

func barfOn(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	//	load up privkey and pubkey
	pubKeyMaterial, err := secrets.ReadFile("certs/rec.la-cert.crt")
	barfOn(err)
	privKeyMaterial, err := secrets.ReadFile("certs/rec.la-key.pem")
	barfOn(err)
	cert, err := tls.X509KeyPair(pubKeyMaterial, privKeyMaterial)
	barfOn(err)

	//	watch directory
	watcher := fswatch.NewAutoWatcher(*watchDir)
	fsEvents := watcher.Start()

	//	dispatch events to SSE sseBroker
	sseBroker := NewBroker()
	go func() {
		for b := range fsEvents {
			sseBroker.Notifier <- *b
		}
	}()

	//	start web server
	mux := http.NewServeMux()

	//	static files
	staticFileServer := http.FileServer(http.Dir(*watchDir))
	mux.Handle("/", injectHeadersForStaticFiles(staticFileServer))

	//	.hak/js/*
	mux.Handle(hakPrefix+"/js/", hakHandler())

	//	./hak/fs/sse
	mux.Handle(path.Join(hakPrefix, ssePath), sseBroker)

	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	//fmt.Printf("running on https://fasthak.rec.la:%d\n\n", *portPtr)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)

	if err != nil {
		log.Fatalln(err)
	}

}

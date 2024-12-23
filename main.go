package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/sean9999/fasthak/certs"
	gorph "github.com/sean9999/go-fsnotify-recursively"
)

// constants
const hakPrefix = "/.hak"
const ssePath = "fs/sse"
const domain = "backloop.dev"

var (
	dir  *string
	port *int

	//go:embed frontend/*
	frontend embed.FS
)

func init() {
	//	parse options and args
	dir = flag.String("dir", ".", "what directory to watch")
	port = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func barfOn(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	//	watch directory
	watcher, _ := gorph.New(fmt.Sprintf("%s/**", *dir))
	fsEvents, _ := watcher.Listen()

	//	braodcast to all cients
	sseBroker := NewBroadcaster[gorph.GorphEvent]()

	//	debounce noisy events
	debouncer := debounce(fsEvents)

	//	broadcsast debounced events
	go func() {
		for ev := range debouncer.Subscribe() {

			if ev.NotifyEvent.Name != "CHMOD" {
				fmt.Println(ev)
			}

			sseBroker.Broadcast(ev)
		}
	}()

	mux := http.NewServeMux()

	//	static files
	staticFileServer := http.FileServer(http.Dir(*dir))
	mux.Handle("/", injectHeaders(staticFileServer))

	//	.hak/js/*
	mux.Handle(hakPrefix+"/js/", hakHandler())

	//	./hak/fs/sse
	mux.Handle(path.Join(hakPrefix, ssePath), sseBroker)

	cert, err := tls.X509KeyPair(certs.Cert, certs.Key)
	barfOn(err)

	//	start server
	url := fmt.Sprintf("%s:%d", domain, *port)
	fmt.Printf("running on https://%s\n\n", url)
	err = ListenAndServeTLSKeyPair(fmt.Sprintf(":%d", *port), cert, mux)
	barfOn(err)

}

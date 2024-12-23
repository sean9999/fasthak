package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/sean9999/fasthak/certs"
	gorph "github.com/sean9999/go-fsnotify-recursively"
)

const (
	hakPrefix = "/.hak"
	ssePath   = "fs/sse"
	domain    = "backloop.dev"
)

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

	cert, err := certs.KeyPair()
	barfOn(err)

	//	start server
	fmt.Printf("running on %s:%d", domain, *port)
	err = ListenAndServeTLSKeyPair(fmt.Sprintf(":%d", *port), cert, mux)
	barfOn(err)

}

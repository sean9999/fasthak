package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"path"

	"github.com/abc-inc/browser"
	"github.com/sean9999/fasthak/certs"
	gorph "github.com/sean9999/go-fsnotify-recursively"
)

func serve() error {

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
	if err != nil {
		return err
	}

	//	start server
	url := fmt.Sprintf("https://%s.%s:%d", randomSlug(), domain, *port)
	fmt.Printf("running on %s", url)
	browser.Open(url)
	ListenAndServeTLSKeyPair(fmt.Sprintf(":%d", *port), cert, mux)
	// if err != nil {
	// 	return nil, err
	// }

	return nil

}

func randomSlug() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(65536))
	return fmt.Sprintf("%x", n)
}

func barfOn(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

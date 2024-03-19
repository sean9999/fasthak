package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	gorph "github.com/sean9999/go-fsnotify-recursively"
)

// constants
const hakPrefix = "/.hak"
const ssePath = "fs/sse"

var (
	watchDir *string
	portPtr  *int

	//go:embed certs/*
	secrets embed.FS

	//go:embed frontend/*
	frontend embed.FS
)

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	flag.Parse()
}

func barfOn(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	//	load up TLS certs
	pubKeyMaterial, err := secrets.ReadFile("certs/rec.la-cert.crt")
	barfOn(err)
	privKeyMaterial, err := secrets.ReadFile("certs/rec.la-key.pem")
	barfOn(err)
	cert, err := tls.X509KeyPair(pubKeyMaterial, privKeyMaterial)
	barfOn(err)

	//	watch directory
	watcher, _ := gorph.New(fmt.Sprintf("%s/**", *watchDir))
	fsEvents, _ := watcher.Listen()

	//	braodcast to all cients
	sseBroker := NewBroadcaster[gorph.GorphEvent]()

	//	debounce noisy events
	debouncer := debounce(fsEvents)

	//	broadcsast debounced events
	go func() {
		for ev := range debouncer.Subscribe() {
			fmt.Println(ev)
			sseBroker.Broadcast(ev)
		}
	}()

	// go func() {
	// 	for ev := range fsEvents {
	// 		fmt.Println("fsEvent", ev)
	// 		sseBroker.Broadcast(ev)
	// 	}
	// }()

	mux := http.NewServeMux()

	//	static files
	staticFileServer := http.FileServer(http.Dir(*watchDir))
	mux.Handle("/", injectHeadersForStaticFiles(staticFileServer))

	//	.hak/js/*
	mux.Handle(hakPrefix+"/js/", hakHandler())

	//	./hak/fs/sse
	mux.Handle(path.Join(hakPrefix, ssePath), sseBroker)

	//	start server
	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	fmt.Printf("running on https://fasthak.rec.la:%d\n\n", *portPtr)
	err = ListenAndServeTLSKeyPair(portString, cert, mux)
	if err != nil {
		log.Fatalln(err)
	}

}

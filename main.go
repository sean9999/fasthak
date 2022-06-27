package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"path"
)

const hakPrefix = "/.hak"
const ssePath = "fs/sse"

var (
	watchDir *string
	portPtr  *int
	//go:embed frontend/*
	frontend embed.FS
	privKey  *string
	pubKey   *string
)

type NiceEvent struct {
	Event string
	File  string
}

func init() {
	//	parse options and arguments
	//	@todo: sanity checking
	watchDir = flag.String("dir", ".", "what directory to watch")
	portPtr = flag.Int("port", 9443, "what port to listen on")
	privKey = flag.String("privkey", "localhost-key.pem", "location of private key")
	pubKey = flag.String("pubkey", "localhost.pem", "location of public key")
	flag.Parse()
}

func main() {

	//	start watcher
	var niceEvents = make(chan NiceEvent)
	if err := watchRecursively(*watchDir, niceEvents); err != nil {
		log.Fatal(err)
	}

	//	dispatch watcher events to SSE sseBroker
	sseBroker := NewBroker()
	go func() {
		for b := range niceEvents {
			sseBroker.Notifier <- b
		}
	}()

	fmt.Println(path.Join(hakPrefix, ssePath))

	//	start server
	err := serve(*watchDir, *portPtr, *privKey, *pubKey, sseBroker)
	if err != nil {
		panic("server could not start")
	}

}

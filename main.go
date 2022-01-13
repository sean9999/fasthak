package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/rjeczalik/notify"
	"log"
	"net/http"
)

//	constants
const ssePath = "/.hak/fs/sse"

func main() {
	//	parse options and arguments
	watchDir := flag.String("dir", ".", "what directory to watch")
	portPtr := flag.Int("port", 9443, "what port to listen on")
	flag.Parse()

	//	start watcher
	eventsChannel := make(chan notify.EventInfo, 1)
	if err := notify.Watch(*watchDir, eventsChannel, notify.Remove); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(eventsChannel)
	fmt.Printf("watching folder %s\n", *watchDir)

	//	Configure web server
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(*watchDir))
	mux.Handle("/", fileServer)

	//	SSE Events
	mux.Handle(ssePath, &fsEventHandler{})

	//	Start Web Server
	portString := fmt.Sprintf("%s%d", ":", *portPtr)
	fmt.Printf("listening on port %d\n", *portPtr)
	err := http.ListenAndServeTLS(portString, "./localhost.pem", "localhost-key.pem", mux)
	if err != nil {
		log.Fatal(err)
		return
	}

}

type fsEventHandler struct {
	msg notify.Event
}

func (ev *fsEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for {
		select {
		case event := <-eventsChannel.Events:
			pushEvent(event, w)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			//log.Panic("watcher.Events was Context().Done(). What does this mean?")
			return
		}
	}
}

func pushEvent(msg notify.EventInfo, w http.ResponseWriter) {
	/**
	 *	push a fileSystem Event through SSE
	 */

	// @todo: maybe find a way to only call this once per connection
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	e := FsEvent{
		File:  msg.Path(),
		Event: msg.Event().String(),
	}
	log.Println(msg)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		return
	}
	fprintf, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", "fs", buf.String())
	if err != nil {
		log.Panic(fprintf, err)
		return
	}
}

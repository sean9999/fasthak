package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
)

const hakPrefix = ".hak"
const hakEventNamespace = "fs"

//	watcher needs to be global (for now until I figure out how Go works)
var watcher, watcherBootstrapError = fsnotify.NewWatcher()

type FsEvent struct {
	Event string
	File  string
}

func main() {

	//	parse options and arguments
	watchDir := flag.String("dir", ".", "what directory to watch")
	portPtr := flag.Int("port", 9001, "what port to listen on")
	flag.Parse()

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

	//	start web server
	fs := http.FileServer(http.Dir(*watchDir))
	http.Handle("/", injectHeaders(fs))
	http.Handle("/"+hakPrefix+"/"+hakEventNamespace+"/sse", handler(fsEventHandler))
	portString := fmt.Sprintf("%s%d", ":", *portPtr)

	err := http.ListenAndServeTLS(portString, "./localhost.pem", "localhost-key.pem", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func injectHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, max-age=1")
		fs.ServeHTTP(w, r)
	}
}

func handler(f http.HandlerFunc) http.Handler {
	return f
}

func pushEvent(msg fsnotify.Event, w http.ResponseWriter) {

	e := FsEvent{
		File:  msg.Name,
		Event: msg.Op.String(),
	}
	log.Println(msg)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", "fs", buf.String())

}

func fsEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher := w.(http.Flusher)

	//	consume fileSystem events
	for {
		select {
		case event := <-watcher.Events:
			pushEvent(event, w)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}

}

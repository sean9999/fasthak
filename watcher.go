package main

import (
	"bytes"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

func notifyEventInfoToNiceEvent(ei notify.EventInfo) NiceEvent {
	abs, _ := filepath.Abs(*watchDir)
	return NiceEvent{
		File:  strings.TrimPrefix(ei.Path(), abs+"/"),
		Event: ei.Event().String(),
	}
}

func niceEventToBuffer(ne NiceEvent) (bytes.Buffer, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(ne)
	return buf, err
}

func toBytes(ei notify.EventInfo) []byte {
	ne := notifyEventInfoToNiceEvent(ei)
	buf, _ := niceEventToBuffer(ne)
	b := buf.Bytes()
	return b
}

//	watchRecursively emits event info to the "niceEvents" channel
func watchRecursively(path string, niceEvents chan NiceEvent) error {

	var c = make(chan notify.EventInfo)
	err := notify.Watch(path+"/...", c, notify.All)

	//	massage the event to the format we want
	go func() {
		for eventInfo := range c {
			niceEvent := notifyEventInfoToNiceEvent(eventInfo)
			log.Printf("%s - %s", niceEvent.Event, niceEvent.File)
			niceEvents <- niceEvent
		}
	}()

	return err
}

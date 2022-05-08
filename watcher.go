package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

type NiceEvent struct {
	Event string
	File  string
}

func toNiceEvent(ei notify.EventInfo) NiceEvent {

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

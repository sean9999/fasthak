package main

import "github.com/fsnotify/fsnotify"

//	watcher needs to be global (for now until I figure out how Go works)
var watcher, watcherBootstrapError = fsnotify.NewWatcher()

type FsEvent struct {
	Event string
	File  string
}

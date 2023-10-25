package main

import (
	"github.com/gokyle/fswatch"
)

func Xwatch(dir string) {
	var _ = fswatch.NewAutoWatcher(dir)

}

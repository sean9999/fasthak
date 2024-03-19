package main

import (
	"time"

	gorph "github.com/sean9999/go-fsnotify-recursively"
	"github.com/sean9999/rebouncer"
)

type gev = gorph.GorphEvent

func debounce(messyEvents chan gev) rebouncer.Rebouncer[gev] {

	//	simply pass in all events
	var ingestFn rebouncer.Ingester[gev] = func(fsEvents chan<- gev) {
		for ev := range messyEvents {
			fsEvents <- ev
		}
	}

	//	simply accumulate all events except ones involving .DS_Store
	var reduceFn rebouncer.Reducer[gev] = func(evs []gev) []gev {
		out := make([]gev, 0)
		for _, ev := range evs {
			//	@todo: filter out tmp files
			if ev.Path != ".DS_Store" {
				out = append(out, ev)
			}
		}
		return out
	}

	//	flush every second
	var quantizeFn rebouncer.Quantizer[gev] = func(queue []gev) bool {
		time.Sleep(1 * time.Second)
		return len(queue) > 0
	}

	rebby := rebouncer.NewRebouncer(ingestFn, reduceFn, quantizeFn, 1024)

	return rebby

}

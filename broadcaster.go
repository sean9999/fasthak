package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Broadcaster[T fmt.Stringer] interface {
	http.Handler
	AddClient(chan T)
	RemoveClient(chan T)
	Clients() []chan T
	Incoming() chan T
	Broadcast(T)
}

type broadcaster[T fmt.Stringer] struct {
	clients  map[chan T]bool
	incoming chan T
	sync.RWMutex
}

func (b *broadcaster[T]) AddClient(msgChan chan T) {
	b.Lock()
	//fmt.Println("adding")
	defer b.Unlock()
	b.clients[msgChan] = true
}

func (b *broadcaster[T]) RemoveClient(msgChan chan T) {
	b.Lock()
	defer b.Unlock()
	delete(b.clients, msgChan)
}

func (b *broadcaster[T]) Clients() []chan T {
	b.RLock()
	defer b.RUnlock()
	chans := make([]chan T, len(b.clients), 0)
	for k, ok := range b.clients {
		if ok {
			chans = append(chans, k)
		}
	}
	return chans
}

func (b *broadcaster[T]) Incoming() chan T {
	return b.incoming
}

func (b *broadcaster[T]) Broadcast(msg T) {
	fmt.Println("broadcast", len(b.clients))
	for ch := range b.clients {
		fmt.Println(ch)
		ch <- msg
	}
}

func (b *broadcaster[T]) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// ResponseWriter must support flushing for SSE to be effective
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Flushing unsupported", http.StatusInternalServerError)
		return
	}

	//	HTTP headers
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// a dedicated channel for each client connection
	messageChan := make(chan T)
	b.AddClient(messageChan)
	defer func() {
		b.RemoveClient(messageChan)
	}()

	//	request cancellation
	notify := req.Context().Done()

	for {
		select {
		case <-notify:
			return
		default:
			// Write Server Sent Event to the response
			msg := <-messageChan
			fmt.Fprintf(rw, "%s", msg)
			// Flush data immediately to user agent
			flusher.Flush()
		}
	}

}

func NewBroadcaster[T fmt.Stringer]() Broadcaster[T] {
	b := broadcaster[T]{
		clients:  make(map[chan T]bool),
		incoming: make(chan T),
	}
	return &b
}

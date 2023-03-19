package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sean9999/rebouncer"
)

const patience time.Duration = time.Second * 1

func StringifyEvent(ne rebouncer.NiceEvent) string {
	dataToken := ne.Operation + "\n" + ne.File
	data := base64.StdEncoding.EncodeToString([]byte(dataToken))
	return fmt.Sprintf("id: %d\nevent: fs\ndata: %s\nretry: 3001\n\n", ne.Id, data)
}

type Broker struct {
	Notifier       chan rebouncer.NiceEvent
	newClients     chan chan rebouncer.NiceEvent
	closingClients chan chan rebouncer.NiceEvent
	clients        map[chan rebouncer.NiceEvent]bool
}

func NewBroker() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan rebouncer.NiceEvent, 1),
		newClients:     make(chan chan rebouncer.NiceEvent),
		closingClients: make(chan chan rebouncer.NiceEvent),
		clients:        make(map[chan rebouncer.NiceEvent]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()
	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan rebouncer.NiceEvent)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := req.Context().Done()

	for {
		select {
		case <-notify:
			return
		default:

			// Write to the ResponseWriter
			// Server Sent Events compatible

			stringifiedEvent := StringifyEvent(<-messageChan)
			fmt.Fprintf(rw, "%s", stringifiedEvent)

			// Flush data immediately
			flusher.Flush()
		}
	}

}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has detached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan := range broker.clients {
				select {
				case clientMessageChan <- event:
				case <-time.After(patience):
					log.Print("Skipping client.")
				}
			}
		}
	}
}

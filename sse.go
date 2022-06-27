package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

const patience time.Duration = time.Second * 1

//	serializer for NiceEvent
func (ne NiceEvent) String() string {
	dataToken := ne.Event + "\n" + ne.File
	data := base64.StdEncoding.EncodeToString([]byte(dataToken))
	return fmt.Sprintf("event: fs\ndata: %s\nretry: 3001\n\n", data)
}

type Broker struct {
	Notifier       chan NiceEvent
	newClients     chan chan NiceEvent
	closingClients chan chan NiceEvent
	clients        map[chan NiceEvent]bool
}

func NewBroker() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan NiceEvent, 1),
		newClients:     make(chan chan NiceEvent),
		closingClients: make(chan chan NiceEvent),
		clients:        make(map[chan NiceEvent]bool),
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
	messageChan := make(chan NiceEvent)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	ctx := req.Context()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// send SSE data
			fmt.Fprintf(rw, "%s", <-messageChan)
			// Flush to user-agent
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

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan := range broker.clients {
				select {
				case clientMessageChan <- event:
					log.Print("sse event")
				case <-time.After(patience):
					log.Print("Skipping client.")
				}
			}
		}
	}
}

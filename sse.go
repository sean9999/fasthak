package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gokyle/fswatch"
)

const patience time.Duration = time.Second * 1

// const (
// 	NONE     = iota // No event, initial state.
// 	CREATED         // File was created.
// 	DELETED         // File was deleted.
// 	MODIFIED        // File was modified.
// 	PERM            // Changed permissions
// 	NOEXIST         // File does not exist.
// 	NOPERM          // No permissions for the file (see const block comment).
// 	INVALID         // Any type of error not represented above.
// )

func EventToString(ev fswatch.Notification) string {
	lookup := map[int]string{
		0: "none",
		1: "created",
		2: "deleted",
		3: "modified",
		4: "changed_permissions",
		5: "file_does_not_exist",
		6: "no_permissions_for_file",
		7: "invalid_event",
	}
	return fmt.Sprintf("%s\n%s", lookup[ev.Event], ev.Path)
}

func StringifyEvent(ev fswatch.Notification) string {
	dataToken := EventToString(ev)
	data := base64.StdEncoding.EncodeToString([]byte(dataToken))
	return fmt.Sprintf("id: %d\nevent: fs\ndata: %s\nretry: 3001\n\n", time.Now().Unix(), data)
}

type Broker struct {
	Notifier       chan fswatch.Notification
	newClients     chan chan fswatch.Notification
	closingClients chan chan fswatch.Notification
	clients        map[chan fswatch.Notification]bool
}

func NewBroker() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan fswatch.Notification, 1),
		newClients:     make(chan chan fswatch.Notification),
		closingClients: make(chan chan fswatch.Notification),
		clients:        make(map[chan fswatch.Notification]bool),
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
	messageChan := make(chan fswatch.Notification)

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

			// Flush data immediately to user agent
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

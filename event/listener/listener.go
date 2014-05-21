// listener package implement an event listener server with resequencing
// support.
// It listen for incoming tcp connection and start decoding events sent
// unordered through that connection.
// Reading is performed using buffered io.

// Once an event is decoded a resequencing strategy is applied,
// ensuring that outgoing events are sent in the correct order regarding
// in respect to their sequence ID.

// Once an EventSource disconnect, EventSourceCloseChan is sent a value,
// which other routines can use to handle event source disconnection.

package listener

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/andreadipersio/efr/event"
)

type Listener struct {
	Port int

	// Resequenced events are sent through this channel
	DispatchChan chan *event.Event

	// Use to inform other routines of EventSource disconnection
	EventSourceCloseChan chan interface{}

	ResequencerType string

	// Max size of incoming events queue.
	// The bigger the value, the bigger the memory consumption.
	ResequencerCapacity int
}

// Listen for incoming connection from EventSource
func (l *Listener) Listen() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", l.Port))

	defer ln.Close()

	if err != nil {
		log.Fatalf("Cannot start Event Listener: %v", err)
	}

	log.Printf("=== Event Listener waiting for connection on %v", l.Port)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Printf("Cannot read from socket: %v", err)
			continue
		}

		go l.handleEventSourceConnection(conn)
	}
}

// handleEventSourceConnection handle a tcp connection sending
// batch of events.
func (l *Listener) handleEventSourceConnection(conn net.Conn) {
	defer conn.Close()

	resequencer := NewResequencer(l.ResequencerType, l.ResequencerCapacity)
	scanner := bufio.NewScanner(conn)

	log.Printf("  = EventSource connected, %s enabled", resequencer)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Printf("Cannot read payload: %v", err)
			continue
		}

		payload := scanner.Text()
		evt, err := event.FromString(payload)

		if err != nil {
			log.Printf("Cannot decode payload: %v", err)
			continue
		}

		resequencer.Resequence(evt, l.DispatchChan)
	}

	log.Println("EventSource disconnected")
	// send all events in resequencer buffer (guaranted to be sorted)
	resequencer.Flush(l.DispatchChan)
	// notify other routines of event source disconnection
	l.EventSourceCloseChan <- nil

	return
}

func New(
	port int,
	dspChan chan *event.Event,
	ctrlChan chan interface{},
	resequencerType string,
	resequencerCap int,
) *Listener {
	return &Listener{
		Port:                 port,
		DispatchChan:         dspChan,
		EventSourceCloseChan: ctrlChan,
		ResequencerCapacity:  resequencerCap,
		ResequencerType:      resequencerType,
	}
}

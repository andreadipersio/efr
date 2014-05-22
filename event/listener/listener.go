// listener package implement an event listener server with resequencing
// support.
// It listen for incoming tcp connection and start decoding events sent
// unordered through that connection.
// Reading is performed using buffered io.
//
// Once an event is decoded a resequencing strategy is applied,
// ensuring that outgoing events are sent in the correct order regarding
// in respect to their sequence ID.
//
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
	DispatchChan chan event.Event

	// Use to inform other routines of EventSource disconnection
	EventSourceCloseChan chan interface{}

	ResequencerConfig *ResequencerConfig

	EventFactory event.EventFactoryType
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
	defer func() {
		// notify other routines of event source disconnection
		l.EventSourceCloseChan <- nil

		// terminate connection with event source
		conn.Close()
	}()

	resequencer := NewResequencer(l.ResequencerConfig)
	scanner := bufio.NewScanner(conn)

	log.Printf("  = EventSource connected, %s enabled", resequencer)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Printf("Cannot read payload: %v", err)
			continue
		}

		payload := scanner.Text()
		e, err := l.EventFactory(payload)

		if err != nil {
			log.Printf("Cannot create event: %v", err)
			continue
		}

		resequencer.Resequence(e, l.DispatchChan)
	}

	log.Println("  = EventSource disconnected")

	// send all events in resequencer buffer (guaranted to be sorted)
	resequencer.Flush(l.DispatchChan)

	return
}

func New(
	port int,
	dspChan chan event.Event,
	ctrlChan chan interface{},
	resequencerConfig *ResequencerConfig,
	eventFactory event.EventFactoryType,
) *Listener {
	return &Listener{
		Port:                 port,
		DispatchChan:         dspChan,
		EventSourceCloseChan: ctrlChan,
		ResequencerConfig:    resequencerConfig,
		EventFactory:         eventFactory,
	}
}

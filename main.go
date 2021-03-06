// Run event forwarder using example event and subscriber implementation
package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/andreadipersio/efr/event"
	"github.com/andreadipersio/efr/event/dispatcher"
	"github.com/andreadipersio/efr/event/listener"
	"github.com/andreadipersio/efr/event/subscription"
	"github.com/andreadipersio/efr/example"
)

func main() {
	var (
		resequencerType = flag.String("resequencerType", "stream",
			"Resequencer type, can be 'batch' or 'stream'")

		resequencerCap = flag.Int("resequencerCapacity", 100,
			"Resequencer capacity.")

		maxProcs = flag.Int("maxProcs", 1,
			"Max number of OS Thread that can run simultaneously")

		eventSourcePort = flag.Int("eventSourcePort", 9090, "EventSource connection port")
		subPort         = flag.Int("clientPort", 9099, "Clients will subscribe using to this port")

		sequenceIndex = flag.Int("sequenceIndex", 0,
			"Last know sequence number. Stream resequencer "+
				"will start resequencing from sequenceIndex+1")
	)

	flag.Parse()

	resequencerConfig := &listener.ResequencerConfig{
		*resequencerType,
		*resequencerCap,
		*sequenceIndex,
	}

	runtime.GOMAXPROCS(*maxProcs)
	log.Printf("Maximum number of concurrent threads set to %v", *maxProcs)

	// Acknowledge event source disconnection
	ctrlChan := make(chan interface{})

	// Dispatcher will wait for ordered events on that channel
	eventChan := make(chan event.Event)

	// Client connections
	subChan := make(chan *subscription.SubscriptionRequest)

	subscriptionServer := subscription.New(*subPort, subChan)

	dispatcher := dispatcher.New(
		eventChan,
		subChan,
		ctrlChan,
		example.NewUser,
	)

	listener := listener.New(
		*eventSourcePort,
		eventChan,
		ctrlChan,
		resequencerConfig,
		example.NewEvent,
	)

	// Listen for event source connection.
	// Provide resequenceing of events
	go listener.Listen()

	// Listen for new client connection
	go subscriptionServer.Listen()

	// Dispatch event between connected client
	go dispatcher.Dispatch()

	select {}
}

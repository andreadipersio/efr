// Run an event forwarder, using maze.NewUser as subscriber.
package main

import (
	"flag"
	"runtime"

	"github.com/andreadipersio/efr/event"
	"github.com/andreadipersio/efr/event/dispatcher"
	"github.com/andreadipersio/efr/event/listener"
	"github.com/andreadipersio/efr/event/subscription"
	"github.com/andreadipersio/efr/followersmaze"
)

func main() {
	var (
		resequencerType = flag.String("resequencerType", "batch",
			"Resequencer type, can be 'batch' or 'stream'")

		resequencerCap = flag.Int("resequencerCapacity", 100,
			"Resequencer capacity.")

		maxProcs = flag.Int("maxProcs", 1,
			"Max number of OS Thread that can run simultaneously")

		eventSourcePort = flag.Int("eventSourcePort", 9090, "EventSource connection port")
		subPort         = flag.Int("clientPort", 9099, "Clients will subscribe using to this port")
	)

	flag.Parse()

	runtime.GOMAXPROCS(*maxProcs)

	// Acknowledge event source disconnection
	ctrlChan := make(chan interface{})

	// Dispatcher will wait for ordered events on that channel
	eventChan := make(chan *event.Event)

	// Client connections
	subChan := make(chan *subscription.SubscriptionRequest)

	subscriptionServer := subscription.New(*subPort, subChan)
	dispatcher := dispatcher.New(eventChan, subChan, ctrlChan, followersmaze.NewUser)
	listener := listener.New(
		*eventSourcePort,
		eventChan,
		ctrlChan,
		*resequencerType,
		*resequencerCap,
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

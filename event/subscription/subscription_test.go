package subscription

import (
	"fmt"
	"net"
	"testing"

	"github.com/andreadipersio/efr/event/subscription"
)

// TestSubscription prove that subscription server can accept connection
// on a port and create a subscriptionRequest which wil be routed through
// a channel
func TestSubscription(t *testing.T) {
	// setup server
	port := 11111
	addr := fmt.Sprintf("localhost:%v", port)

	subChan := make(chan *subscription.SubscriptionRequest)
	s := subscription.New(port, subChan)

	go s.Listen()

	// setup client
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		t.Fatalf("Cannot connect to %v: %v", addr, err)
	}

	testSubscriberID := "123"

	// send request
	fmt.Fprintf(conn, fmt.Sprintf("%v\n", testSubscriberID))

	// verify that SubscriptionRequest has been correctly created
	subReq := <-subChan

	if subReq.SubscriberID != testSubscriberID {
		t.Fatalf("Expected ID %v got '%v'", testSubscriberID, subReq.SubscriberID)
	}
}

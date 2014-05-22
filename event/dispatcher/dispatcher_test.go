package dispatcher

import (
	"bytes"
	"testing"
	"time"

	"github.com/andreadipersio/efr/event"
	"github.com/andreadipersio/efr/event/dispatcher"
	"github.com/andreadipersio/efr/event/subscription"
	"github.com/andreadipersio/efr/followersmaze"
)

// Need to implement io.WriteCloser on top of a bytes buffer
type testBuffer struct {
	bytes.Buffer

	// We mock Write so everything will be writted there
	Content string
}

// Mock write in a way is easy to read what we send through
// Subscriber.Conn
func (b *testBuffer) Write(p []byte) (n int, err error) {
	b.Content = string(p)

	return len(b.Content), nil
}

// Required by io.WriteCloser
func (b *testBuffer) Close() error {
	return nil
}

var (
	subscriberFactory = followersmaze.NewUser
	eventFactory      = followersmaze.NewEvent
)

func createSubscribtionRequest(
	subscriberID string,
) *subscription.SubscriptionRequest {
	return &subscription.SubscriptionRequest{
		SubscriberID: subscriberID,

		// Our mocked io.WriteCloser
		Conn: &testBuffer{},
	}
}

// TestDispatcher prove that dispatcher is able to dispatch
// broadcast events to two subscribers
func TestDispatcher(t *testing.T) {
	dspChan := make(chan event.Event)
	subChan := make(chan *subscription.SubscriptionRequest)
	ctrlChan := make(chan interface{})

	// our test dispatcher
	dsp := dispatcher.New(dspChan, subChan, ctrlChan, subscriberFactory)

	// start dispatching
	go dsp.Dispatch()

	sr1 := createSubscribtionRequest("1")
	sr2 := createSubscribtionRequest("2")

	// subscribe our fake subscribers
	subChan <- sr1
	subChan <- sr2

	// Broadcast
	broadcastEvent, _ := eventFactory("1|B")

	dspChan <- broadcastEvent

	ctrlChan <- nil

	// Check if buffer got data
	writeComplete := func(buff *testBuffer) bool {
		return buff.Content != ""
	}

	// Check if buffer has been written correctly
	check := func(buff *testBuffer) {
		if buff.Content != "1|B\n" {
			t.Fatalf("Expected '%v', got %v", broadcastEvent, buff.Content)
		}
	}

	// Wait max 1 second for buffer to being written
	// by dispatcher go routine
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout while waiting for goroutine" +
			" to send event to subscriber")
	default:
		if writeComplete(sr1.Conn.(*testBuffer)) &&
			writeComplete(sr2.Conn.(*testBuffer)) {
			break
		}
	}

	check(sr1.Conn.(*testBuffer))
	check(sr2.Conn.(*testBuffer))
}

package dispatcher

import (
	"testing"

	"github.com/andreadipersio/efr/event/dispatcher"
	"github.com/andreadipersio/efr/followersmaze"
)

var subscriberFactory = followersmaze.NewUser

func TestGetOrcreate(t *testing.T) {
	testSubscriberID := "foo"

	d := dispatcher.NewDirectory(subscriberFactory)
	s := d.GetOrCreate(testSubscriberID)

	if s.GetID() != testSubscriberID {
		t.Fatalf("Subscriber was created with wrong ID! Expected %v, got %v",
			testSubscriberID, s.GetID())
	}
}

func TestSubscribe(t *testing.T) {
	testSubscriberID := "foo"
	d := dispatcher.NewDirectory(subscriberFactory)

	subscriber := subscriberFactory(testSubscriberID)

	d.Subscribe(subscriber)

	if _, ok := d.GetByID(testSubscriberID); !ok {
		t.Fatalf("Expected subscriber with ID %v "+
			"to be subscribed. Is not!", testSubscriberID)
	}
}

func TestNew(t *testing.T) {
	testSubscriberID := "foo"

	d := dispatcher.NewDirectory(subscriberFactory)
	d.New(testSubscriberID)

	if _, ok := d.GetByID(testSubscriberID); !ok {
		t.Fatalf("Expected subscriber with ID %v "+
			"to be subscribed. Is not!", testSubscriberID)
	}
}

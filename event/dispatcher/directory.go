package dispatcher

import "github.com/andreadipersio/efr/event"

// dispatchDirectory provide storing of subscribers by their ids
type dispatchDirectory struct {
	storage           map[string]event.Subscriber
	subscriberFactory event.SubscriberFactoryType
}

// GetOrCreate try to get a subscriber from directory, if it does not exist,
// create a disconnected user
func (d *dispatchDirectory) GetOrCreate(subscriberID string) event.Subscriber {
	subscriber, exist := d.storage[subscriberID]

	if exist {
		return subscriber
	}

	d.New(subscriberID)

	return d.storage[subscriberID]
}

// Create a new disconnected directory subscriber
func (d *dispatchDirectory) New(subscriberID string) {
	s := d.subscriberFactory(subscriberID)

	d.storage[subscriberID] = s
}

// Subscribe perform subscriber subscription to the subscribers directory
func (d *dispatchDirectory) Subscribe(s event.Subscriber) {
	d.storage[s.GetID()] = s
}

// UnsubscribeAll unsubscribe all subscribers by deleting them from
// the subscriber directory and, if they are connected, disconnect them
func (d *dispatchDirectory) UnsubscribeAll() {
	for subscriberID, s := range d.storage {
		if s.IsConnected() {
			s.Disconnect()
		}

		delete(d.storage, subscriberID)
	}
}

// Broadcast send event e to all subscribers in the directory
func (d *dispatchDirectory) Broadcast(e *event.Event) {
	for _, s := range d.storage {
		s.SendEvent(e)
	}
}

// SenderAndRecipientFromEvent return event Sender and event Receiver.
// If they are not registered in the directory, create them.
func (d *dispatchDirectory) SenderAndRecipientFromEvent(e *event.Event) (event.Subscriber, event.Subscriber) {
	return d.GetOrCreate(e.SenderID), d.GetOrCreate(e.RecipientID)
}

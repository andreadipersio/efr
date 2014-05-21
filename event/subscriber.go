package event

import (
	"io"
)

// Subscriber interface define an event subscriber,
// which should be able to Connnect and receive events
// directed to him, and also to handle/notify followers.
type Subscriber interface {
    // Provide a connection to the subscriber,
    // so events can be sent to it
	Connect(io.WriteCloser)

    // Close connection with subscriber
	Disconnect()

    // Return whenever a subscriber is connected
	IsConnected() bool

    // Get the identity value for this subscriber
	GetID() string
    // Set the identity value for this subscriber
	SetID(string)

    // Handle an incoming event
	HandleEvent(*Event, Subscriber) error
    // Send an event to the subscriber
	SendEvent(*Event)

    // A subscriber can have followers, which depending on message
    // type have to be notified
	GetFollowers() []Subscriber
	NewFollower(Subscriber)
	RemoveFollower(string)

    // Used for initialization of subscriber implementation
	Init()
}

// Given an ID return a subscriber.
type SubscriberFactoryType func(ID string) Subscriber

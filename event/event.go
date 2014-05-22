// event package implement Event and Subscriber data type.
package event

import "fmt"

type Event interface {
	SequenceNum() int
	SenderID() string
	RecipientID() string
	EventType() string

	Parse(string) error
	fmt.Stringer
}

// BySequence provide sorting of events by Sequence
type BySequence []Event

func (e BySequence) Len() int           { return len(e) }
func (e BySequence) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e BySequence) Less(i, j int) bool { return e[i].SequenceNum() < e[j].SequenceNum() }

// EventFactory represent a function that taken a string
// return an event concrete value or an error
type EventFactoryType func(string) (Event, error)

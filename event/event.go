// event package implement Event and Subscriber data type.
package event

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	fieldDelimiter = "|"
)

type Event struct {
	Sequence              int
	Type                  string
	SenderID, RecipientID string
}

// fromString read a string in the format
// 123|S|56
// and return an Event
func FromString(s string) (*Event, error) {
	parts := strings.Split(s, fieldDelimiter)

	validateFieldNumber := func(fields []string) error {
		if len(parts) < 2 {
			return fmt.Errorf("Event is incomplete, should contains at least Sequence and Type")
		}

		return nil
	}

	if err := validateFieldNumber(parts); err != nil {
		return nil, err
	}

	sequenceNum, err := strconv.Atoi(parts[0])

	if err != nil {
		return nil, fmt.Errorf("Invalid sequence")
	}

	e := &Event{
		Sequence: sequenceNum,
		Type:     parts[1],
	}

	if len(parts) > 2 {
		e.SenderID = parts[2]
	}

	if len(parts) == 4 {
		e.RecipientID = parts[3]
	}

	return e, nil
}

// String return the same string from which the Event has been
// constructed
func (e *Event) String() string {
	parts := []string{fmt.Sprintf("%v", e.Sequence), e.Type}

	if e.SenderID != "" {
		parts = append(parts, e.SenderID)
	}

	if e.RecipientID != "" {
		parts = append(parts, e.RecipientID)
	}

	return strings.Join(parts, fieldDelimiter)
}

// BySequence provide sorting of events by Sequence
type BySequence []*Event

func (e BySequence) Len() int           { return len(e) }
func (e BySequence) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e BySequence) Less(i, j int) bool { return e[i].Sequence < e[j].Sequence }

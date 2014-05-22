package example

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andreadipersio/efr/event"
)

const (
	fieldDelimiter = "|"
)

type Event struct {
	sequence int
	// Event Type
	eType                 string
	senderID, recipientID string
}

// fromString read a string in the format
//     123|S|56
// and return an Event
func (e *Event) Parse(s string) error {
	parts := strings.Split(s, fieldDelimiter)

	validateFieldNumber := func(fields []string) error {
		if len(parts) < 2 {
			return fmt.Errorf("Event is incomplete, should contains at least Sequence and Type")
		}

		return nil
	}

	if err := validateFieldNumber(parts); err != nil {
		return err
	}

	sequenceNum, err := strconv.Atoi(parts[0])

	if err != nil {
		return fmt.Errorf("Invalid sequence: %v", err)
	}

	e.sequence = sequenceNum
	e.eType = parts[1]

	if len(parts) > 2 {
		e.senderID = parts[2]
	}

	if len(parts) == 4 {
		e.recipientID = parts[3]
	}

	return nil
}

func (e *Event) SequenceNum() int {
	return e.sequence
}

func (e *Event) SenderID() string {
	return e.senderID
}

func (e *Event) RecipientID() string {
	return e.recipientID
}

func (e *Event) EventType() string {
	return e.eType
}

// String return the same string from which the Event has been
// constructed
func (e *Event) String() string {
	parts := []string{fmt.Sprintf("%v", e.sequence), e.eType}

	if e.senderID != "" {
		parts = append(parts, e.senderID)
	}

	if e.recipientID != "" {
		parts = append(parts, e.recipientID)
	}

	return strings.Join(parts, fieldDelimiter)
}

func NewEvent(s string) (event.Event, error) {
	e := &Event{}

	if err := e.Parse(s); err == nil {
		return e, nil
	} else {
		return nil, err
	}
}

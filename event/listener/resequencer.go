package listener

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andreadipersio/efr/event"
)

type Resequencer interface {
	// Resequence append an event to a buffer and based on the
	// resequencing strategy it check if an ordered sequence of
	// events can be streamed to the output channel.
	Resequence(e *event.Event, outChan chan *event.Event)

	// Send all the remaining events in buffer to outChan.
	// Sent item are always ordered but sequence may be incomplete.
	// Always empty the buffer.
	Flush(outChan chan *event.Event)

	fmt.Stringer
}

// NewResequencer return the correct resequencer for the choosen type
// or BatchResequencer if type is wrong.
func NewResequencer(rType string, capacity int) Resequencer {
	switch strings.ToLower(rType) {
	case "batch":
		return NewBatchResequencer(capacity)
	case "stream":
		return NewStreamResequencer(capacity)
	default:
		return NewBatchResequencer(capacity)
	}
}

// A Batch Resequencer which keep a buffer of events and
// when the buffer is on full Capacity it flush it.
// Effectiveness is determined by the relation between randomness
// and batch size of the event source in respect to Capacity,
// which also is directly related to memory consumption.
type BatchResequencer struct {
	Capacity int
	buffer   []*event.Event
}

func (r *BatchResequencer) String() string {
	return fmt.Sprintf("Batch Resequencer(cap %v)", r.Capacity)
}

func (r *BatchResequencer) Append(e *event.Event) {
	r.buffer = append(r.buffer, e)
}

func (r *BatchResequencer) BufferIsFull() bool {
	return len(r.buffer) == r.Capacity
}

// Flush send all the remaining events in the events buffer to
// outChan. Buffer is sorted before sending.
func (r *BatchResequencer) Flush(outChan chan *event.Event) {
	if len(r.buffer) == 0 {
		return
	}

	sort.Sort(event.BySequence(r.buffer))

	for _, e := range r.buffer {
		outChan <- e
	}

	r.buffer = r.buffer[:0]
}

func NewBatchResequencer(capacity int) *BatchResequencer {
	return &BatchResequencer{
		Capacity: capacity,
		buffer:   []*event.Event{},
	}
}

func (r *BatchResequencer) Resequence(e *event.Event, dspChan chan *event.Event) {
	r.Append(e)

	if r.BufferIsFull() {
		r.Flush(dspChan)
	}
}

// A Stream resequencer which send messages as soon
// as they represent a linear
type StreamResequencer struct {
	Capacity  int
	buffer    []*event.Event
	lastIndex int
}

func (r *StreamResequencer) String() string {
	return "Stream Resequencer"
}

func (r *StreamResequencer) Flush(outChan chan *event.Event) {
	if len(r.buffer) == 0 {
		return
	}

	sort.Sort(event.BySequence(r.buffer))

	for _, e := range r.buffer {
		outChan <- e
	}

	r.buffer = r.buffer[:0]
}

func NewStreamResequencer(capacity int) *StreamResequencer {
	return &StreamResequencer{
		Capacity: capacity,
		buffer:   []*event.Event{},
	}
}

func (r *StreamResequencer) Resequence(e *event.Event, dspChan chan *event.Event) {
	r.buffer = append(r.buffer, e)
	sort.Sort(event.BySequence(r.buffer))

	buffer := []*event.Event{}
	copy(buffer, r.buffer)

	for i, e := range buffer {
		index := i + r.lastIndex

		if index == e.Sequence-1 {
			dspChan <- e
			r.buffer = append(r.buffer[:index], r.buffer[index+1])
		} else {
			r.lastIndex = index
		}
	}
}

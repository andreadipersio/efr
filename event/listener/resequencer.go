package listener

import (
	"fmt"
	"sort"
	"strings"

	"github.com/andreadipersio/efr/event"
)

type ResequencerConfig struct {
	// 'stream' or 'batch'
	Type string

	// Max size of incoming events queue.
	// The bigger the value, the bigger the memory consumption.
	Capacity int

	// Start resequencing from SequenceIndex+1
	SequenceIndex int
}

type Resequencer interface {
	// Resequence append an event to a buffer and based on the
	// resequencing strategy it check if an ordered sequence of
	// events can be streamed to the output channel.
	Resequence(e event.Event, outChan chan event.Event)

	// Send all the remaining events in buffer to outChan.
	// Sent item are always ordered but sequence may be incomplete.
	// Always empty the buffer.
	Flush(outChan chan event.Event)

	fmt.Stringer
}

// NewResequencer return the correct resequencer for the choosen type
// or BatchResequencer if type is wrong.
func NewResequencer(config *ResequencerConfig) Resequencer {
	switch strings.ToLower(config.Type) {
	case "batch":
		return NewBatchResequencer(config)
	case "stream":
		return NewStreamResequencer(config)
	default:
		return NewStreamResequencer(config)
	}
}

// A Batch Resequencer which keep a buffer of events and
// when the buffer is on full Capacity it flush it.
// Effectiveness is determined by the relation between randomness
// and batch size of the event source in respect to Capacity,
// which also is directly related to memory consumption.
type BatchResequencer struct {
	Capacity int
	buffer   []event.Event
}

func (r *BatchResequencer) String() string {
	return fmt.Sprintf("Batch Resequencer(cap %v)", r.Capacity)
}

func (r *BatchResequencer) Append(e event.Event) {
	r.buffer = append(r.buffer, e)
}

func (r *BatchResequencer) BufferIsFull() bool {
	return len(r.buffer) == r.Capacity
}

// Flush send all the remaining events in the events buffer to
// outChan. Buffer is sorted before sending.
func (r *BatchResequencer) Flush(outChan chan event.Event) {
	if len(r.buffer) == 0 {
		return
	}

	sort.Sort(event.BySequence(r.buffer))

	for _, e := range r.buffer {
		outChan <- e
	}

	r.buffer = r.buffer[:0]
}

func NewBatchResequencer(config *ResequencerConfig) *BatchResequencer {
	return &BatchResequencer{
		Capacity: config.Capacity,
		buffer:   []event.Event{},
	}
}

func (r *BatchResequencer) Resequence(e event.Event, dspChan chan event.Event) {
	r.Append(e)

	if r.BufferIsFull() {
		r.Flush(dspChan)
	}
}

// A Stream resequencer which send messages as soon
// as they represent a linear
type StreamResequencer struct {
	buffer    map[int]event.Event
	lastIndex int
}

func (r *StreamResequencer) String() string {
	return "Stream Resequencer"
}

func (r *StreamResequencer) Flush(outChan chan event.Event) {
	if len(r.buffer) == 0 {
		return
	}

	// create a temporary buffer slice to sort events
	buff := []event.Event{}

	for _, e := range r.buffer {
		buff = append(buff, e)
	}

	// reset buffer
	r.buffer = map[int]event.Event{}

	// order slice
	sort.Sort(event.BySequence(buff))

	// send item in order
	for _, e := range buff {
		outChan <- e
	}

	// clear our temporary buffer slice
	buff = buff[:0]
}

func NewStreamResequencer(config *ResequencerConfig) *StreamResequencer {
	return &StreamResequencer{
		buffer:    map[int]event.Event{},
		lastIndex: config.SequenceIndex,
	}
}

func (r *StreamResequencer) Resequence(e event.Event, dspChan chan event.Event) {
	r.buffer[e.SequenceNum()] = e

	// Check if we have a valid sequence
	for {
		nextSeqNum := r.lastIndex + 1
		if s, ok := r.buffer[nextSeqNum]; ok {
			dspChan <- s
			r.lastIndex++
			delete(r.buffer, nextSeqNum)
		} else {
			break
		}
	}
}

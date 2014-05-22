package listener

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/andreadipersio/efr/event"
	"github.com/andreadipersio/efr/event/listener"
	"github.com/andreadipersio/efr/followersmaze"
)

// testResequencer prove that resequencer is able to resequence
// unordered events
func testResequencer(t *testing.T, r listener.Resequencer, batchSize int) {
	eventFactory := followersmaze.NewEvent

	rand.Seed(42)

	testSequence := rand.Perm(batchSize)

	// ensure that we are feeding random ordered sequence
	if sort.IntsAreSorted(testSequence) {
		t.Fatalf("testSequence is ordered. Should be random :)")
	}

	// receive ordered event on this channel
	dspChan := make(chan event.Event)

	resequenced := []event.Event{}

	// generate random sequence
	go func() {
		for e := range dspChan {
			resequenced = append(resequenced, e)
		}
	}()

	for _, seq := range testSequence {
		e, err := eventFactory(fmt.Sprintf("%v|B", seq))

		if err != nil {
			t.Fatalf("Cannot create event with sequence %v: %v", seq, err)
		}

		r.Resequence(e, dspChan)
	}

	// if sequence generator does not return all sequence before deadline
	// fail test
	select {
	case <-time.After(1 * time.Second):
		if len(resequenced) != len(testSequence) {
			t.Fatalf("Expected a sequence %v items long, got %v", len(testSequence), len(resequenced))
		}
	default:
		if len(resequenced) == len(testSequence) {
			break
		}
	}

	if !sort.IsSorted(event.BySequence(resequenced)) {
		t.Fatalf("Sequence is not ordered: %v", resequenced)
	}
}

func TestBatchResequencer(t *testing.T) {
	batchSize := 100

	config := &listener.ResequencerConfig{"batch", batchSize, 0}
	r := listener.NewBatchResequencer(config)

	testResequencer(t, r, batchSize)
}

func TestStreamResequencer(t *testing.T) {
	batchSize := 100

	config := &listener.ResequencerConfig{"batch", batchSize, 0}
	r := listener.NewStreamResequencer(config)

	testResequencer(t, r, batchSize)
}

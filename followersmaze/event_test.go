package followersmaze

import (
	"log"
	"testing"

	"github.com/andreadipersio/efr/followersmaze"
)

// TestParse prove that Event.Parse method is able
// to correctly parse strings and generate error when
// strings contains wrong format.
func TestParse(t *testing.T) {
	type testDataType struct {
		payload string
		isValid bool
	}

	testEvents := []testDataType{
		testDataType{"1|B", true},
		testDataType{"150|F|12|13", true},
		testDataType{"22|S|12", true},

		testDataType{"A|S|15", false},
		testDataType{"", false},

		testDataType{"2|S|123|13|53", true},
	}

	for _, testEvent := range testEvents {
		evt, err := followersmaze.NewEvent(testEvent.payload)

		// Test should fail, but it hasn't!
		if !testEvent.isValid && err == nil {
			log.Fatalf("Parsing %v should fail! Got event %v", testEvent.payload, evt)
		}

		// Test should succeed but it hasn't!
		if testEvent.isValid && err != nil {
			log.Fatalf("Parsing %v should succeed! Got error: %v", testEvent.payload, err)
		}
	}
}

// TestString prove that Event.String return the same
// string it got when parsing when invoked.
func TestString(t *testing.T) {
	type testDataType struct {
		payload string
		isValid bool
	}

	testEvents := []testDataType{
		testDataType{"1|B", true},
		testDataType{"150|F|12|13", true},
		testDataType{"22|S|12", true},
	}

	for _, testEvent := range testEvents {
		e, err := followersmaze.NewEvent(testEvent.payload)

		if e.String() != testEvent.payload {
			log.Fatalf("Got '%v', expected '%v'", e, testEvent.payload)
		}

		if err != nil {
			log.Fatalf("Should not have failed! %v -> %v", testEvent.payload, err)
		}
	}
}

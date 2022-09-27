package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/protocol"
)

func TestDownstreamMessage(t *testing.T) {
	eui := protocol.EUIFromInt64(0x0abcdef0)
	msg := NewDownstreamMessage(eui, 100)
	if msg.CreatedTime == 0 {
		t.Fatal("Expected created time to be set")
	}

	if msg.IsComplete() || msg.State() != UnsentState {
		t.Fatal("Message should not be completed and in unsent state")
	}

	// No ack and message is sent: not pending
	msg.SentTime = time.Now().Unix()
	if !msg.IsComplete() || msg.State() != SentState {
		t.Fatal("Expected message to be completed and in sent state")
	}

	// Set ack flag. The state should become pending
	msg.Ack = true
	if msg.IsComplete() {
		t.Fatal("Expected message not to be completed")
	}

	msg.AckTime = time.Now().Unix()
	if !msg.IsComplete() || msg.State() != AcknowledgedState {
		t.Fatal("Expected message to be completed and acknowledged state")
	}

	msg.Data = "010203040506070809"
	if !reflect.DeepEqual(msg.Payload(), []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}) {
		t.Fatal("Not the payload I expected")
	}

	msg.Data = "Random characters"
	if len(msg.Payload()) != 0 {
		t.Fatal("Expected empty payload")
	}
}

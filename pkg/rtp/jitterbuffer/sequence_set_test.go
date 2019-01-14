package jitterbuffer

import (
	"testing"

	"github.com/pions/webrtc/pkg/rtp"
)

func TestSequentialPush(t *testing.T) {
	ss := sequenceSet{}

	ss.Push(&rtp.Packet{SequenceNumber: 0})
	ss.Push(&rtp.Packet{SequenceNumber: 1})
	ss.Push(&rtp.Packet{SequenceNumber: 2})

	if len(ss) != 3 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(ss), 3)
	}
	if ss[0].SequenceNumber != 0 {
		t.Errorf("Element 0 incorrect, got: %d, want: %d.", ss[0].SequenceNumber, 0)
	}
	if ss[1].SequenceNumber != 1 {
		t.Errorf("Element 1 incorrect, got: %d, want: %d.", ss[1].SequenceNumber, 1)
	}
	if ss[2].SequenceNumber != 2 {
		t.Errorf("Element 2 incorrect, got: %d, want: %d.", ss[2].SequenceNumber, 2)
	}
}

func TestNonSequentialPush(t *testing.T) {
	ss := sequenceSet{}

	ss.Push(&rtp.Packet{SequenceNumber: 2})
	ss.Push(&rtp.Packet{SequenceNumber: 0})
	ss.Push(&rtp.Packet{SequenceNumber: 1})

	if len(ss) != 3 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(ss), 3)
	}
	if ss[0].SequenceNumber != 0 {
		t.Errorf("Element 0 incorrect, got: %d, want: %d.", ss[0].SequenceNumber, 0)
	}
	if ss[1].SequenceNumber != 1 {
		t.Errorf("Element 1 incorrect, got: %d, want: %d.", ss[1].SequenceNumber, 1)
	}
	if ss[2].SequenceNumber != 2 {
		t.Errorf("Element 2 incorrect, got: %d, want: %d.", ss[2].SequenceNumber, 2)
	}
}

func TestSequentialWrapPush(t *testing.T) {
	ss := sequenceSet{}

	ss.Push(&rtp.Packet{SequenceNumber: 65534})
	ss.Push(&rtp.Packet{SequenceNumber: 65535})
	ss.Push(&rtp.Packet{SequenceNumber: 0})

	if len(ss) != 3 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(ss), 3)
	}
	if ss[0].SequenceNumber != 65534 {
		t.Errorf("Element 0 incorrect, got: %d, want: %d.", ss[0].SequenceNumber, 65534)
	}
	if ss[1].SequenceNumber != 65535 {
		t.Errorf("Element 1 incorrect, got: %d, want: %d.", ss[1].SequenceNumber, 65535)
	}
	if ss[2].SequenceNumber != 0 {
		t.Errorf("Element 2 incorrect, got: %d, want: %d.", ss[2].SequenceNumber, 0)
	}
}

func TestNonSequentialWrapPush(t *testing.T) {
	ss := sequenceSet{}

	ss.Push(&rtp.Packet{SequenceNumber: 65534})
	ss.Push(&rtp.Packet{SequenceNumber: 0})
	ss.Push(&rtp.Packet{SequenceNumber: 65535})

	if len(ss) != 3 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(ss), 3)
	}
	if ss[0].SequenceNumber != 65534 {
		t.Errorf("Element 0 incorrect, got: %d, want: %d.", ss[0].SequenceNumber, 65534)
	}
	if ss[1].SequenceNumber != 65535 {
		t.Errorf("Element 1 incorrect, got: %d, want: %d.", ss[1].SequenceNumber, 65535)
	}
	if ss[2].SequenceNumber != 0 {
		t.Errorf("Element 2 incorrect, got: %d, want: %d.", ss[2].SequenceNumber, 0)
	}
}

func TestInsertDuplicatePush(t *testing.T) {
	ss := sequenceSet{}

	if ok := ss.Push(&rtp.Packet{SequenceNumber: 65534}); !ok {
		t.Errorf("Push should have succeeded")
	}

	if ok := ss.Push(&rtp.Packet{SequenceNumber: 0}); !ok {
		t.Errorf("Push should have succeeded")
	}

	if ok := ss.Push(&rtp.Packet{SequenceNumber: 65534}); ok {
		t.Errorf("Push should have failed")
	}

	if len(ss) != 2 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(ss), 3)
	}
	if ss[0].SequenceNumber != 65534 {
		t.Errorf("Element 0 incorrect, got: %d, want: %d.", ss[0].SequenceNumber, 65534)
	}
	if ss[1].SequenceNumber != 0 {
		t.Errorf("Element 1 incorrect, got: %d, want: %d.", ss[1].SequenceNumber, 0)
	}
}

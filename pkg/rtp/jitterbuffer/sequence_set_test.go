package jitterbuffer

import (
	"testing"
)

func TestSequentialPush(t *testing.T) {
	ss := sequenceSet{}

	ss.Push(newRTPTestPacket(0, false))
	ss.Push(newRTPTestPacket(1, false))
	ss.Push(newRTPTestPacket(2, false))

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

	ss.Push(newRTPTestPacket(2, false))
	ss.Push(newRTPTestPacket(0, false))
	ss.Push(newRTPTestPacket(1, false))

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

	ss.Push(newRTPTestPacket(65534, false))
	ss.Push(newRTPTestPacket(65535, false))
	ss.Push(newRTPTestPacket(0, false))

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

	ss.Push(newRTPTestPacket(65534, false))
	ss.Push(newRTPTestPacket(0, false))
	ss.Push(newRTPTestPacket(65535, false))

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

	if ok := ss.Push(newRTPTestPacket(65534, false)); !ok {
		t.Errorf("Push should have succeeded")
	}

	if ok := ss.Push(newRTPTestPacket(0, false)); !ok {
		t.Errorf("Push should have succeeded")
	}

	if ok := ss.Push(newRTPTestPacket(65534, false)); ok {
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

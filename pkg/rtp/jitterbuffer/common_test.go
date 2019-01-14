package jitterbuffer

import (
	"fmt"
	"testing"

	"github.com/pions/webrtc/pkg/rtp"
)

func newRTPTestPacket(sequenceNumber uint16, marker bool) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			SequenceNumber: sequenceNumber,
			Marker:         marker,
		},
	}
}

func TestSequenceCompare(t *testing.T) {
	fmt.Println()
	if cmp := compareSeqNum(65535, 0); cmp != -1 {
		t.Errorf("Compare failed, got: %d, want: %d.", cmp, -1)
	}
	if cmp := compareSeqNum(65535, 1); cmp != -2 {
		t.Errorf("Compare failed, got: %d, want: %d.", cmp, -2)
	}
	if cmp := compareSeqNum(0, 2); cmp != -2 {
		t.Errorf("Compare failed, got: %d, want: %d.", cmp, -2)
	}
	if cmp := compareSeqNum(2, 0); cmp != 2 {
		t.Errorf("Compare failed, got: %d, want: %d.", cmp, 2)
	}
	if cmp := compareSeqNum(65535, 65534); cmp != 1 {
		t.Errorf("Compare failed, got: %d, want: %d.", cmp, 1)
	}
}

func TestSequenceNumber(t *testing.T) {
	var s sequenceNumber

	if s.isValid {
		t.Error("Sequence number should be invalid at initialization")
	}

	s.setValue(123)
	if !s.isValid {
		t.Error("Sequence number should be valid after set")
	}
	if s.value != 123 {
		t.Errorf("Value incorrect, got: %d, want: %d.", s.value, 123)
	}

	s.reset()
	if s.isValid {
		t.Error("Sequence number should be invalid after reset")
	}
}

package jitterbuffer

import (
	"testing"
)

func TestPushNoFirstPacket(t *testing.T) {
	sample := newFragmentedSample()
	sample.pushPacket(newRTPTestPacket(1, false))

	if sample.hasFirstPacket() {
		t.Error("Sample should not have first packet")
	}

	if sample.hasLastPacket() {
		t.Error("Sample should not have last packet")
	}

	if sample.isComplete() {
		t.Error("Sample should not be complete")
	}
}

func TestPushFirstPacket(t *testing.T) {
	sample := newFragmentedSample()
	sample.pushPacket(newRTPTestPacket(1, false))

	if sample.hasFirstPacket() {
		t.Error("Sample should not have first packet yet, not set")
	}

	sample.setFirstPacketSeqNum(1)
	if !sample.hasFirstPacket() {
		t.Error("Sample should have first packet after setting first packet sequence number")
	}

	if sample.hasLastPacket() {
		t.Error("Sample should not have last packet")
	}

	if sample.isComplete() {
		t.Error("Sample should not be complete")
	}
}

func TestPushLastPacket(t *testing.T) {
	sample := newFragmentedSample()

	if sample.hasLastPacket() {
		t.Error("Sample should not have last packet yet")
	}

	sample.pushPacket(newRTPTestPacket(1, true))

	if sample.hasFirstPacket() {
		t.Error("Sample should not have first packet")
	}

	if !sample.hasLastPacket() {
		t.Error("Sample should have last packet")
	}

	if sample.isComplete() {
		t.Error("Sample should not be complete")
	}
}

func TestPushCompleteWithSinglePacket(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, true))

	if !sample.isComplete() {
		t.Error("Sample should be complete")
	}
}

func TestPushCompleteWithMultiPackets(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(3, false))
	sample.pushPacket(newRTPTestPacket(4, true))

	if !sample.isComplete() {
		t.Error("Sample should be complete")
	}
}

func TestPushCompleteWithMultiWrappingPackets(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	sample.setFirstPacketSeqNum(65534)
	sample.pushPacket(newRTPTestPacket(65534, false))
	sample.pushPacket(newRTPTestPacket(65535, false))
	sample.pushPacket(newRTPTestPacket(0, false))
	sample.pushPacket(newRTPTestPacket(1, true))

	if !sample.isComplete() {
		t.Error("Sample should be complete")
	}
}

func TestPushWithMissingMiddlePacket(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(5, true))

	if sample.isComplete() {
		t.Error("Sample should not be complete as it is missing 1 packet")
	}
}

func TestIsCompleteCaching(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	if sample.isCompleteCached {
		t.Error("Complete should not be cached as it is still incomplete")
	}

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(5, true))

	if sample.isComplete() {
		t.Error("Sample should not be complete as it is missing 1 packet")
	}
	if sample.isCompleteCached {
		t.Error("Complete should still not be cached as it is still incomplete (missing 1 packet)")
	}

	sample.pushPacket(newRTPTestPacket(3, false))
	if !sample.isComplete() {
		t.Error("Sample should be complete")
	}
	if !sample.isCompleteCached {
		t.Error("Complete should be cached as match the return value of isComplete()")
	}
	// Excercises the early return
	if !sample.isComplete() {
		t.Error("Sample should be complete with early return from cache being set")
	}
}

func TestGenerateNacksWithCompleteSinglePacketSample(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, true))

	nacks, err := sample.generateNacks(sequenceNumber{value: 2, isValid: true})
	if err != nil {
		t.Errorf("Error, got: %v, want: %v.", err, nil)
	}

	if len(nacks) > 0 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(nacks), 0)
	}
}

func TestGenerateNacksWithCompleteMultiPacketSample(t *testing.T) {
	sample := newFragmentedSample()

	if sample.isComplete() {
		t.Error("Sample should not be complete yet")
	}

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(3, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(5, true))

	nacks, err := sample.generateNacks(sequenceNumber{value: 6, isValid: true})
	if err != nil {
		t.Errorf("Error, got: %v, want: %v.", err, nil)
	}

	if len(nacks) > 0 {
		t.Errorf("Length incorrect, got: %d, want: %d.", len(nacks), 0)
	}
}

func TestGenerateNacksWithOverlapNextSeqNum(t *testing.T) {
	sample := newFragmentedSample()

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(3, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(5, true))

	_, err := sample.generateNacks(sequenceNumber{value: 3, isValid: true})
	if err != errNextSeqNumOverlaps {
		t.Errorf("Error incorrect, got: %v, want: %v.", err, errNextSeqNumOverlaps)
	}
}

func TestGenerateNacksWithTooSmallSeqNum(t *testing.T) {
	sample := newFragmentedSample()

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(3, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(5, true))

	_, err := sample.generateNacks(sequenceNumber{value: 50000, isValid: true})
	if err != errNextSeqNumTooSmall {
		t.Errorf("Error incorrect, got: %v, want: %v.", err, errNextSeqNumTooSmall)
	}
}

func TestGenerateNacksWithOneMissingFragment(t *testing.T) {
	sample := newFragmentedSample()

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(5, true))

	nacks, err := sample.generateNacks(sequenceNumber{value: 6, isValid: true})
	if err != nil {
		t.Errorf("Error incorrect, got: %v, want: %v.", err, nil)
	}

	if len(nacks) != 1 {
		t.Errorf("Nacks length incorrect, got: %d, want: %d.", len(nacks), 1)
	}

	if nacks[0] != 3 {
		t.Errorf("Nack incorrect, got: %d, want: %d.", nacks[0], 3)
	}
}

func TestGenerateNacksWithSparseMissingFragments(t *testing.T) {
	sample := newFragmentedSample()

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, false))
	sample.pushPacket(newRTPTestPacket(2, false))
	sample.pushPacket(newRTPTestPacket(4, false))
	sample.pushPacket(newRTPTestPacket(6, true))

	nacks, err := sample.generateNacks(sequenceNumber{value: 7, isValid: true})
	if err != nil {
		t.Errorf("Error incorrect, got: %v, want: %v.", err, nil)
	}

	if len(nacks) != 2 {
		t.Errorf("Nacks length incorrect, got: %d, want: %d.", len(nacks), 2)
	}

	if nacks[0] != 3 {
		t.Errorf("Nack incorrect, got: %d, want: %d.", nacks[0], 3)
	}

	if nacks[1] != 5 {
		t.Errorf("Nack incorrect, got: %d, want: %d.", nacks[1], 5)
	}
}

func TestGenerateNacksWithNextSeqNumAfterLastPacket(t *testing.T) {
	sample := newFragmentedSample()

	sample.setFirstPacketSeqNum(1)
	sample.pushPacket(newRTPTestPacket(1, true))

	nacks, err := sample.generateNacks(sequenceNumber{value: 3, isValid: true})
	if err != nil {
		t.Errorf("Error incorrect, got: %v, want: %v.", err, nil)
	}

	if len(nacks) != 1 {
		t.Errorf("Nacks length incorrect, got: %d, want: %d.", len(nacks), 1)
	}

	if nacks[0] != 2 {
		t.Errorf("Nack incorrect, got: %d, want: %d.", nacks[0], 2)
	}
}

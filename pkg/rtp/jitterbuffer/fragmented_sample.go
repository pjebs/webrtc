package jitterbuffer

import (
	"errors"

	"github.com/pions/webrtc/pkg/rtp"
)

const invalidSeqNum = -1

var (
	errNextSeqNumTooSmall = errors.New("next sequence number less than first sequence number in the sample")
	errNextSeqNumOverlaps = errors.New("next sequence number overlaps between first and last sequence number in the sample")
)

type fragmentedSample struct {
	isCompleteCached  bool
	firstPacketSeqNum sequenceNumber
	packets           sequenceSet
}

func newFragmentedSample() *fragmentedSample {
	return &fragmentedSample{}
}

func (f *fragmentedSample) hasFirstPacket() bool {
	if f.firstPacketSeqNum.isValid && len(f.packets) > 0 {
		return f.packets[0].SequenceNumber == f.firstPacketSeqNum.value
	}

	return false
}

func (f *fragmentedSample) hasLastPacket() bool {
	if len(f.packets) > 0 {
		return f.packets[len(f.packets)-1].Marker
	}

	return false
}

func (f *fragmentedSample) setFirstPacketSeqNum(seqNum uint16) {
	f.firstPacketSeqNum.setValue(seqNum)
}

func (f *fragmentedSample) isComplete() bool {
	if f.isCompleteCached {
		return true
	}

	if f.hasFirstPacket() && f.hasLastPacket() {
		// First packet == Last packet (and has mark)
		if len(f.packets) == 1 {
			f.isCompleteCached = true
			return true
		}

		// We know the first packet is valid due to hasFirstPacket check
		previousSeqNum := f.firstPacketSeqNum.value
		for _, p := range f.packets[1:] {
			if previousSeqNum+1 != p.SequenceNumber {
				return false
			}
			previousSeqNum++
		}
		f.isCompleteCached = true
		return true
	}

	return false
}

func (f *fragmentedSample) pushPacket(packet *rtp.Packet) {
	if f.isCompleteCached {
		// Trying to push to completed packet
		// TODO: Add some sort of very stern warning
		return
	}
	if ok := f.packets.Push(packet); !ok {
		// This packet already exists in this fragmented sample
		// TODO: Record the number of duplicate packets
		return
	}
}

func (f *fragmentedSample) generateNacks(nextSampleStartingSeqNum sequenceNumber) (nacks []uint16, err error) {
	// If we have a fragmented sample, we have at least 1 packet
	previousSeqNum := f.packets[0].SequenceNumber

	if nextSampleStartingSeqNum.isValid {
		// We are larger or equal to the next samples sequence number
		// This should never happen unless we are extremely out of date
		if compareSeqNum(previousSeqNum, nextSampleStartingSeqNum.value) > 0 {
			return nil, errNextSeqNumTooSmall
		}

		// Get the last sequence number
		lastSeqNum := f.packets[len(f.packets)-1].SequenceNumber
		// Only check if it not the same as the one we already checked
		// Check if there is overlap between nextSampleStartingSeqNum and first to last packet in this sample
		if lastSeqNum != previousSeqNum {
			if compareSeqNum(lastSeqNum, nextSampleStartingSeqNum.value) >= 0 {
				return nil, errNextSeqNumOverlaps
			}
		}

	}

	for _, p := range f.packets[1:] {
		for i := previousSeqNum + 1; i != p.SequenceNumber; i++ {
			nacks = append(nacks, i)
		}

		previousSeqNum = p.SequenceNumber
	}

	if nextSampleStartingSeqNum.isValid {
		for i := previousSeqNum + 1; i != nextSampleStartingSeqNum.value; i++ {
			nacks = append(nacks, i)
		}
	}

	return nacks, nil
}

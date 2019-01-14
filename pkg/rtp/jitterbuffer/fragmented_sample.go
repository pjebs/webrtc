package jitterbuffer

import "github.com/pions/webrtc/pkg/rtp"

const invalidSeqNum = -1

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

}

func (f *fragmentedSample) generateNacks(nextSampleStartingSeqNum sequenceNumber) (nacks []uint16) {
	// If we have a fragmented sample, we have at least 1 packet
	previousSeqNum := f.packets[0].SequenceNumber

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

	return nacks
}

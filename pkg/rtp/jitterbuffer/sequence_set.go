package jitterbuffer

import (
	"sort"

	"github.com/pions/webrtc/pkg/rtp"
)

type sequenceSet []*rtp.Packet

func (s *sequenceSet) Push(packet *rtp.Packet) (ok bool) {
	index := sort.Search(len(*s), func(i int) bool { return compareSeqNum((*s)[i].SequenceNumber, packet.SequenceNumber) > 0 })
	if index > 0 {
		// Check if this is a duplicate
		if (*s)[index-1].SequenceNumber == packet.SequenceNumber {
			ok = false
			return
		}
	}
	*s = append(*s, packet)
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = packet
	ok = true
	return
}

package jitterbuffer

import "math"

type sequenceNumber struct {
	value   uint16
	isValid bool
}

func (s *sequenceNumber) setValue(value uint16) {
	s.value = value
	s.isValid = true
}

func (s *sequenceNumber) reset() {
	s.isValid = false
}

func seqNumDistance(left uint16, right uint16) uint16 {
	ileft := int(left)
	iright := int(right)
	idist := ileft - iright
	if idist < 0 {
		return uint16(-idist)
	}
	return uint16(idist)
}

func compareSeqNum(left uint16, right uint16) int16 {
	distance := seqNumDistance(left, right)
	adjustedLeft := int(left)
	adjustedRight := int(right)
	if distance > 32767 {
		if left > 32768 {
			adjustedLeft -= 65536
		}
		if right > math.MaxInt16+1 {
			adjustedRight -= 65536
		}
	}
	return int16(adjustedLeft - adjustedRight)
}

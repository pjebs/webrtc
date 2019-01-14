package jitterbuffer

import (
	"fmt"
	"testing"
)

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

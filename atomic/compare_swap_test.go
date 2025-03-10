package main

import "testing"

func TestWaitCompareSwap(t *testing.T) {
	t.Run("wait 3000 milis", func(t *testing.T) {
		if got := WaitCompareSwap(); got < 3000 {
			t.Errorf("WaitCompareSwap() = %v, want > 3000", got)
		}
	})
}

package slack

import (
	"sync"
	"testing"
)

func TestRTMGetInfoZeroValue(t *testing.T) {
	var rtm RTM
	if info := rtm.GetInfo(); info != nil {
		t.Fatalf("GetInfo() = %v, want nil", info)
	}
}

func TestRTMGetInfoConcurrent(t *testing.T) {
	rtm := New("test-token").NewRTM()
	defer rtm.pingDeadman.Stop()
	var wg sync.WaitGroup
	wg.Go(func() {
		for range 1000 {
			// The connection manager uses this lock when replacing the info.
			rtm.mu.Lock()
			rtm.info = &Info{Team: &Team{ID: "T1"}}
			rtm.mu.Unlock()
		}
	})
	for range 1000 {
		_ = rtm.GetInfo()
	}
	wg.Wait()
}

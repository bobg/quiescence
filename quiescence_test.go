package quiescence

import (
	"testing"
	"time"
)

func TestWaiter(t *testing.T) {
	var (
		w         = NewWaiter()
		doneCh    = make(chan struct{})
		notBefore = time.Now().Add(time.Second)
	)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			w.Ping()
		}
	}()

	go func() {
		w.Wait(time.Second)
		close(doneCh)
	}()

	timeout := time.NewTimer(3 * time.Second)

	select {
	case <-doneCh:
		if time.Now().Before(notBefore) {
			t.Errorf("Wait returned too soon")
		}

	case <-timeout.C:
		t.Error("timeout")
	}
}

// Package quiescence defines a "quiescence waiter" type.
// It allows a caller to wait until some resource has gone quiet for a given amount of time.
package quiescence

import (
	"sync"
	"time"
)

// Waiter is a "quiescence waiter."
// It can be pinged,
// and it can wait until no ping has happened for a certain interval.
type Waiter struct {
	mu sync.Mutex
	c  *sync.Cond
	t  time.Time
}

// New produces a new quiescence waiter.
func NewWaiter() *Waiter {
	w := new(Waiter)
	w.c = sync.NewCond(&w.mu)
	return w
}

// Ping resets the timer in the quiescence waiter.
func (w *Waiter) Ping() {
	w.mu.Lock()
	w.t = time.Now()
	w.c.Broadcast()
	w.mu.Unlock()
}

// Wait returns after interval dur elapses with no pings.
func (w *Waiter) Wait(dur time.Duration) {
	var (
		timeCh = make(chan time.Time)
		doneCh = make(chan struct{})
	)

	go func() {
		w.mu.Lock()
		defer w.mu.Unlock()

		for {
			select {
			case <-doneCh:
				return
			default:
			}

			w.c.Wait()
			timeCh <- w.t
		}
	}()

	t := time.NewTimer(dur)

	for {
		select {
		case <-t.C:
			close(doneCh)
			return

		case timestamp := <-timeCh:
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Until(timestamp.Add(dur)))
		}
	}
}

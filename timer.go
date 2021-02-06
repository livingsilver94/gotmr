package gotmr

import "time"

// Timer is a thread-safe wrapper of time.Time that exposes easy Stop and Reset commands,
// without the need to drain the notification channel.
//
// Timer is not a drop-in replacement of time.Timer because Stop and Reset work differently.
// Namely, a Stop is not necessary to Reset the Timer, and there is no boolean returned.
type Timer struct {
	C     <-chan time.Time
	timer timer
}

// Timer works by controlling the internal time.Timer in a single goroutine.
// The commands in fact are passed via a command channel, which the goroutine is receiving from.
// Even the ticks are copied from the Timer.C channel to another exposed channel, so that users
// of Timer have no access to time.Timer.C.

type timer struct {
	c   chan time.Time
	cmd chan interface{}
	t   *time.Timer
}

// New creates a new Timer that will send the current time on its channel after at least duration d.
func New(d time.Duration) Timer {
	tmr := timer{
		c:   make(chan time.Time, 1),
		cmd: make(chan interface{}),
		t:   time.NewTimer(d),
	}
	go func() {
		for {
			select {
			case tick := <-tmr.t.C:
				tmr.c <- tick
			case cmd := <-tmr.cmd:
				if !tmr.t.Stop() {
					<-tmr.t.C
				}
				switch t := cmd.(type) {
				case bool:
					return
				case time.Duration:
					tmr.t.Reset(t)
				}
			}
		}
	}()
	return Timer{C: tmr.c, timer: tmr}
}

// Stop aborts the timer, garbage-collecting its resources.
// No Timer methods must be called after Stop, including Stop itself.
func (t Timer) Stop() {
	t.timer.cmd <- true
}

// Reset changes the timer to expire after duration d.
func (t Timer) Reset(d time.Duration) {
	t.timer.cmd <- d
}

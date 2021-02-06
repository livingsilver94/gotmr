package gotmr_test

import (
	"context"
	"testing"
	"time"

	"github.com/livingsilver94/gotmr"
)

const (
	wait     = time.Millisecond
	longWait = wait + time.Millisecond
)

func TestTimerTicks(t *testing.T) {
	tmr := gotmr.New(wait)
	select {
	case <-time.After(longWait):
		t.Error("A tick was not received from the Timer")
	case <-tmr.C:
	}
}

func TestTimerResets(t *testing.T) {
	t1 := time.Now()
	tmr := gotmr.New(longWait)
	tmr.Reset(wait)
	<-tmr.C
	if time.Since(t1) >= longWait {
		t.Error("Timer did not reset")
	}
}

func TestTimerStops(t *testing.T) {
	tmr := gotmr.New(wait)
	tmr.Stop()
	select {
	case <-tmr.C:
		t.Error("A tick was received: Timer did not stop")
	case <-time.After(longWait):
	}
}

func TestTimerResetsThread(t *testing.T) {
	ctx, stop := context.WithTimeout(context.Background(), longWait*2)
	defer stop()
	tick := make(chan time.Time)
	t1 := time.Now()
	tmr := gotmr.New(longWait)
	go func() {
		select {
		case t := <-tmr.C:
			tick <- t
		case <-ctx.Done():
			tick <- time.Now()
		}
	}()
	tmr.Reset(wait)
	t2 := <-tick
	if t2.Sub(t1) >= longWait {
		t.Error("Timer did not reset")
	}
}

func TestTimerStopsThread(t *testing.T) {
	ctx, stop := context.WithTimeout(context.Background(), longWait)
	defer stop()
	tick := make(chan time.Time)
	t1 := time.Now()
	tmr := gotmr.New(wait)
	go func() {
		select {
		case t := <-tmr.C:
			tick <- t
		case <-ctx.Done():
			tick <- time.Now()
		}
	}()
	tmr.Stop()
	t2 := <-tick
	if t2.Sub(t1) <= wait {
		t.Error("Timer did not stop")
	}
}

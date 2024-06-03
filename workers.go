package main

import (
	"runtime"
	"sync"
	"time"
)

type Worker struct {
	ch      chan func()
	wg      sync.WaitGroup
	dog     chan struct{}
	dogWait sync.WaitGroup
}

// NewWorker creates a worker pool with an optional
// watchdog timeout. If the workers stall for the timeout period,
// the watchdog triggers.
func NewWorker(d time.Duration) *Worker {
	// Create a set of workers that listen on a channel and
	// call a function.
	w := &Worker{}
	count := runtime.NumCPU()
	w.ch = make(chan func(), count)
	if d != 0 {
		w.dog = make(chan struct{}, 10)
		w.dogWait.Add(1)
		go w.watchdog(d)
	}
	w.wg.Add(count)
	for i := 0; i < count; i++ {
		go w.worker()
	}
	return w
}

// Wait shuts down the worker pool by closing the channel and
// waits for the workers to finish.
func (w *Worker) Wait() {
	close(w.ch)
	w.wg.Wait()
	if w.dog != nil {
		close(w.dog)
		w.dogWait.Wait()
	}
}

// Execute a function on one of the workers.
// If a watchdog is enabled, send a keepalive.
func (w *Worker) Run(f func()) {
	if w.dog != nil {
		w.dog <- struct{}{}
	}
	w.ch <- f
}

// worker listens for a function to dispatch and then calls it.
// When the channel closes, exit.
func (w *Worker) worker() {
	defer w.wg.Done()
	for {
		f, ok := <-w.ch
		if !ok {
			return
		}
		f()
	}
}

// watchdog starts a timer and watches for
// keepalives and the timer expiry (which logs a fatal message).
func (w *Worker) watchdog(t time.Duration) {
	defer w.dogWait.Done()
	ticker := time.NewTicker(t)
	for {
		select {
		case <-ticker.C:
			panic("Watchdog timeout!")
		case _, ok := <-w.dog:
			if !ok {
				// Watchdog shutdown.
				ticker.Stop()
				return
			}
			ticker.Reset(t)
		}
	}
}

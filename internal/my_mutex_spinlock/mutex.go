package mymutex

import (
	"runtime"
	"sync/atomic"
)

const (
	locked         = true
	unlocked       = false
	deferredRetrys = 5
)

type Mutex struct {
	state atomic.Bool
}

func (mu *Mutex) Lock() {
	/*
		если CompareAndSwap смог поменять значение
		state, то вернется true, иначе жгем поток for
	*/
	counter := deferredRetrys
	for !mu.state.CompareAndSwap(unlocked, locked) {
		counter--
		if counter == 0 {
			runtime.Gosched()
			counter = deferredRetrys
		}
	}
}

// func (mu *Mutex) TryLock() bool {}
func (mu *Mutex) Unlock() {
	mu.state.Store(unlocked)
}

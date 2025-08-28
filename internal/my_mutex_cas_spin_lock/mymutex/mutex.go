package mymutex

import (
	"runtime"
	"sync/atomic"
)

const (
	locked    = true
	unlocked  = false
	spinCount = 5
)

type Mutex struct {
	state atomic.Bool
}

func (mu *Mutex) Lock() {
	/*
		если CompareAndSwap смог поменять значение
		state, то вернется true, иначе греем поток
		в количестве spinCount раз и после
		переводим горутину в runabler,
		т.е. переводится в конец очереди горутин.
	*/
	counter := spinCount
	for !mu.state.CompareAndSwap(unlocked, locked) {
		counter--
		if counter == 0 {
			runtime.Gosched()
			counter = spinCount
		}
	}
}

// func (mu *Mutex) TryLock() bool {}
func (mu *Mutex) Unlock() {
	mu.state.Store(unlocked)
}

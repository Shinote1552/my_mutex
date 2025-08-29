package mymutexcas

import (
	"runtime"
	"sync/atomic"
)

const (
	locked           = true
	unlocked         = false
	spinCountLock    = 80
	spinCountTryLock = 10
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
	counter := spinCountLock
	for !mu.state.CompareAndSwap(unlocked, locked) {
		counter--
		if counter == 0 {
			runtime.Gosched()
			counter = spinCountLock
		}
	}
}

func (mu *Mutex) TryLock() bool {
	counter := spinCountTryLock
	for !mu.state.CompareAndSwap(unlocked, locked) {
		counter--
		if counter == 0 {
			counter = spinCountTryLock
			return false
		}
	}
	return true
}
func (mu *Mutex) Unlock() {
	mu.state.Store(unlocked)
}

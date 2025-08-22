package mymutexticket

import "sync/atomic"

// Реализовать мьютекс с талонами (аналогично sync.Mutex)
type MyMutex struct {
	nextTicket atomic.Int64
	current    atomic.Int64
}

func (t *Mutex) Lock() int64 {
	// TODO: получить талон и ждать своего номера
}

func (t *Mutex) Unlock(ticket int64) {
	// TODO: передать ход следующему
}

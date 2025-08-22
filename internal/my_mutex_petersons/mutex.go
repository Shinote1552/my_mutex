package mymutexpetersons

import "sync/atomic"

// Реализовать мьютекс для двух горутин по алгоритму Петерсона
type Mutex struct {
	flags [2]atomic.Bool
	turn  atomic.Int32
}

func (p *Mutex) Lock(threadID int) {
	// TODO: алгоритм Петерсона для двух потоков
}

func (p *Mutex) Unlock(threadID int) {
	// TODO: освобождение
}

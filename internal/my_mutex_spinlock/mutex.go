package mymutexspinlock

import "sync/atomic"

// Реализовать спинлок на основе атомарных операций
type Mutex struct {
	locked atomic.Bool
}

func (s *Mutex) Lock() {
	// TODO: реализовать захват через цикл ожидания
}

func (s *Mutex) Unlock() {
	// TODO: реализовать освобождение
}

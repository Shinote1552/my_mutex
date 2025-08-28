package mymutex

import (
	"runtime"
	"sync/atomic"
)

const (
	locked    = true
	unlocked  = false
	spinCount = 80
)

type Mutex struct {
	ownerTicket atomic.Int64
	nextTicket  atomic.Int64
}

/*
по сути ownerTicket потсоянно будет	пытаться
догнать nextTicket, догоняет он его лишь тогда,
когда кто то Unlock вызовет, потом срабатывает Lock

пока ownerTicket будет пытаться он каждый раз переводится
в runable, т.е. переводится в конец очереди горутин
*/

func (mu *Mutex) Lock() {
	// получаем текущий билет и задаем в очереди следующий
	ticket := mu.nextTicket.Add(1) - 1

	for i := 0; i < spinCount; i++ {
		if mu.ownerTicket.Load() == ticket {
			// мьютекс захвачен!
			return
		}
	}

	// Если не получилось за spinCount попыток, переводим в runable
	for mu.ownerTicket.Load() != ticket {
		runtime.Gosched()
	}

}

/*
в оргинале он менее конкуретный, здесь это реализовать тяжело так как
имеет много оптимизаций внутри, поэтому мой TryLock это почти самый обычный Lock
*/
func (mu *Mutex) TryLock() bool {
	if mu.nextTicket.Load() != mu.ownerTicket.Load() {
		return false
	}

	ticket := mu.nextTicket.Add(1) - 1
	return mu.ownerTicket.Load() == ticket
}

func (mu *Mutex) Unlock() {
	mu.ownerTicket.Add(1)
}

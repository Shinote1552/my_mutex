# my_mutex

## The most important thing to understand:
The machine does not execute the code you wrote...


## Ready now Context:
package mycontext `internal/my_context/mycontext`
- WithCancel
- WithoutCancel
- WithDeadline
- WithTimeout

## Ready mutexes:
package mymutex CAS `internal/my_mutex_cas_spin_lock/mymutex`
package mymutex CAS `internal/my_mutex_ticket_spin_lock/mymutex`
- CAS lock(with spin_lock) used 1 atomic.bool
- Ticket lock(with spin_lock) used 2 atomic.Uint32
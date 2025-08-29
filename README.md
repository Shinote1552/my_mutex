# my_concurency

## The most important thing to understand:
> The machine does not execute the code you wrote...
## Available Mutex Implementations

### 1. CAS Spin Lock
**Package**: `mymutex` (`internal/my_mutex_cas_spin_lock/mymutex`)

**Implementation**: Compare-And-Swap with spin lock
**Atomic Variables**: 1 `atomic.Bool`

### 2. Ticket Spin Lock  
**Package**: `mymutex` (`internal/my_mutex_ticket_spin_lock/mymutex`)

**Implementation**: Ticket-based spin lock
**Atomic Variables**: 2 `atomic.Uint32`

## Context Implementation

**Package**: `mycontext` (`internal/my_context/mycontext`)

**Available Functions**:
- `WithCancel`
- `WithoutCancel` 
- `WithDeadline`
- `WithTimeout`

## Performance Testing Results

Tested(not clean Benchmark) `mycontext` with different mutex implementations:

| Mutex Type | Execution Time |
|------------|----------------|
| `sync` (standard library) | 0.416s |
| `mymutextic` (ticket spin lock) | 0.558s |
| `mymutexcas` (CAS spin lock) | 0.418s |


package mymutexcas

import (
	"sync"
	"testing"
)

func TestMutexCAS_LockUnlock(t *testing.T) {
	var mu Mutex
	mu.Lock()
	mu.Unlock()
	// Should not panic
}

func TestMutexCAS_ConcurrentAccess(t *testing.T) {
	var mu Mutex
	var counter int
	var wg sync.WaitGroup
	iterations := 1000

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()

	if counter != iterations {
		t.Errorf("Expected %d, got %d", iterations, counter)
	}
}

func TestMutexCAS_TryLockSuccess(t *testing.T) {
	var mu Mutex
	if !mu.TryLock() {
		t.Error("TryLock should succeed on unlocked mutex")
	}
	mu.Unlock()
}

func TestMutexCAS_TryLockFailure(t *testing.T) {
	var mu Mutex
	mu.Lock()

	// TryLock should fail when mutex is already locked
	if mu.TryLock() {
		t.Error("TryLock should fail on locked mutex")
	}

	mu.Unlock()
}

func TestMutexCAS_TryLockAfterUnlock(t *testing.T) {
	var mu Mutex
	mu.Lock()
	mu.Unlock()

	if !mu.TryLock() {
		t.Error("TryLock should succeed after unlock")
	}
	mu.Unlock()
}

func TestMutexCAS_Reentrancy(t *testing.T) {
	var mu Mutex
	mu.Lock()

	// TryLock should fail - mutex is not reentrant
	if mu.TryLock() {
		t.Error("Mutex should not be reentrant")
	}

	mu.Unlock()
}

package mymutextic

import (
	"sync"
	"testing"
)

func TestMutexTIC_LockUnlock(t *testing.T) {
	var mu Mutex
	mu.Lock()
	mu.Unlock()
	// Should not panic
}

func TestMutexTIC_ConcurrentAccess(t *testing.T) {
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

func TestMutexTIC_TryLockSuccess(t *testing.T) {
	var mu Mutex
	if !mu.TryLock() {
		t.Error("TryLock should succeed on unlocked mutex")
	}
	mu.Unlock()
}

func TestMutexTIC_TryLockFailure(t *testing.T) {
	var mu Mutex
	mu.Lock()

	// TryLock should fail when mutex is already locked
	if mu.TryLock() {
		t.Error("TryLock should fail on locked mutex")
	}

	mu.Unlock()
}

func TestMutexTIC_TryLockAfterUnlock(t *testing.T) {
	var mu Mutex
	mu.Lock()
	mu.Unlock()

	if !mu.TryLock() {
		t.Error("TryLock should succeed after unlock")
	}
	mu.Unlock()
}

func TestMutexTIC_Fairness(t *testing.T) {
	var mu Mutex
	var wg sync.WaitGroup
	order := make(chan int, 3)

	// First goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		order <- 1
		mu.Unlock()
	}()

	// Second goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		order <- 2
		mu.Unlock()
	}()

	// Third goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		order <- 3
		mu.Unlock()
	}()

	wg.Wait()
	close(order)

	// Check that all goroutines executed
	results := make([]int, 0)
	for result := range order {
		results = append(results, result)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}

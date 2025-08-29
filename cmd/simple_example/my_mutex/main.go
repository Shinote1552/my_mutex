// main.go
package main

import (
	"fmt"
	"my_mutex/internal/mymutexcas"
	"my_mutex/internal/mymutextic"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Сравнение реализации мьютексов ===")
	fmt.Println()

	// Тест 1: Базовая функциональность
	fmt.Println("1. Тест базовой функциональности:")
	testBasicFunctionality()

	// Тест 2: Конкурентный доступ
	fmt.Println("\n2. Тест конкурентного доступа:")
	testConcurrentAccess()

	// Тест 3: TryLock
	fmt.Println("\n3. Тест TryLock:")
	testTryLock()

	// Тест 4: Производительность
	fmt.Println("\n4. Тест производительности:")
	testPerformance()

	// Тест 5: Честность (fairness)
	fmt.Println("\n5. Тест честности:")
	testFairness()

	fmt.Println("\n=== Тесты с примерами на Mymutex завершены ===")
}

func testBasicFunctionality() {
	fmt.Println("   CAS Mutex:")
	casMu := mymutexcas.Mutex{}
	casMu.Lock()
	fmt.Println("   - Lock успешен")
	casMu.Unlock()
	fmt.Println("   - Unlock успешен")

	fmt.Println("   Ticket Mutex:")
	ticMu := mymutextic.Mutex{}
	ticMu.Lock()
	fmt.Println("   - Lock успешен")
	ticMu.Unlock()
	fmt.Println("   - Unlock успешен")
}

func testConcurrentAccess() {
	const iterations = 10000
	fmt.Printf("   Тестирование с %d итераций\n", iterations)

	// CAS Mutex
	casCounter := 0
	casMu := mymutexcas.Mutex{}
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			casMu.Lock()
			casCounter++
			casMu.Unlock()
		}()
	}
	wg.Wait()
	casTime := time.Since(start)

	// Ticket Mutex
	ticCounter := 0
	ticMu := mymutextic.Mutex{}

	start = time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticMu.Lock()
			ticCounter++
			ticMu.Unlock()
		}()
	}
	wg.Wait()
	ticTime := time.Since(start)

	fmt.Printf("   CAS Mutex:  счетчик=%d, время=%v\n", casCounter, casTime)
	fmt.Printf("   Ticket Mutex: счетчик=%d, время=%v\n", ticCounter, ticTime)

	if casCounter != iterations || ticCounter != iterations {
		fmt.Println("   ОШИБКА: Не все инкременты выполнены!")
	}
}

func testTryLock() {
	fmt.Println("   Тестирование TryLock:")

	// CAS Mutex
	casMu := mymutexcas.Mutex{}
	if !casMu.TryLock() {
		fmt.Println("   CAS Mutex: TryLock на свободном мьютексе не удался")
	} else {
		fmt.Println("   CAS Mutex: TryLock на свободном мьютексе успешен")
	}

	// Попытка повторного TryLock на занятом мьютексе
	if casMu.TryLock() {
		fmt.Println("   CAS Mutex: TryLock на занятом мьютексе неожиданно удался")
	} else {
		fmt.Println("   CAS Mutex: TryLock на занятом мьютексе правильно失敗")
	}
	casMu.Unlock()

	// Ticket Mutex
	ticMu := mymutextic.Mutex{}
	if !ticMu.TryLock() {
		fmt.Println("   Ticket Mutex: TryLock на свободном мьютексе не удался")
	} else {
		fmt.Println("   Ticket Mutex: TryLock на свободном мьютексе успешен")
	}

	// Попытка повторного TryLock на занятом мьютексе
	if ticMu.TryLock() {
		fmt.Println("   Ticket Mutex: TryLock на занятом мьютексе неожиданно удался")
	} else {
		fmt.Println("   Ticket Mutex: TryLock на занятом мьютексе правильно失敗")
	}
	ticMu.Unlock()
}

func testPerformance() {
	const operations = 100000
	fmt.Printf("   Тест производительности (%d операций):\n", operations)

	// CAS Mutex производительность
	casMu := mymutexcas.Mutex{}
	start := time.Now()
	for i := 0; i < operations; i++ {
		casMu.Lock()
		casMu.Unlock()
	}
	casTime := time.Since(start)

	// Ticket Mutex производительность
	ticMu := mymutextic.Mutex{}
	start = time.Now()
	for i := 0; i < operations; i++ {
		ticMu.Lock()
		ticMu.Unlock()
	}
	ticTime := time.Since(start)

	fmt.Printf("   CAS Mutex:  %v\n", casTime)
	fmt.Printf("   Ticket Mutex: %v\n", ticTime)
	fmt.Printf("   Отношение: %.2f\n", float64(ticTime.Nanoseconds())/float64(casTime.Nanoseconds()))
}

func testFairness() {
	fmt.Println("   Тест честности (порядок захвата):")

	// CAS Mutex fairness
	fmt.Println("   CAS Mutex:")
	testMutexFairness(&mymutexcas.Mutex{})

	// Ticket Mutex fairness (должен быть более fair)
	fmt.Println("   Ticket Mutex:")
	testMutexFairness(&mymutextic.Mutex{})
}

func testMutexFairness(mu interface{}) {
	var wg sync.WaitGroup
	order := make(chan int, 5)

	lockFunc := func(id int) {
		defer wg.Done()
		switch m := mu.(type) {
		case *mymutexcas.Mutex:
			m.Lock()
			order <- id
			m.Unlock()
		case *mymutextic.Mutex:
			m.Lock()
			order <- id
			m.Unlock()
		}
	}

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go lockFunc(i)
		// Небольшая задержка для создания определенного порядка запуска
		time.Sleep(time.Microsecond * 100)
	}

	wg.Wait()
	close(order)

	fmt.Printf("   Порядок захвата: ")
	for id := range order {
		fmt.Printf("%d ", id)
	}
	fmt.Println()
}

// Дополнительный тест: смешанное использование Lock и TryLock
func testMixedUsage() {
	fmt.Println("\n6. Тест смешанного использования Lock/TryLock:")

	casMu := mymutexcas.Mutex{}
	successfulTries := 0
	totalTries := 100

	for i := 0; i < totalTries; i++ {
		if casMu.TryLock() {
			successfulTries++
			casMu.Unlock()
		}
		time.Sleep(time.Microsecond * 10)
	}

	fmt.Printf("   CAS Mutex: успешных TryLock: %d/%d\n", successfulTries, totalTries)

	ticMu := mymutextic.Mutex{}
	successfulTries = 0

	for i := 0; i < totalTries; i++ {
		if ticMu.TryLock() {
			successfulTries++
			ticMu.Unlock()
		}
		time.Sleep(time.Microsecond * 10)
	}

	fmt.Printf("   Ticket Mutex: успешных TryLock: %d/%d\n", successfulTries, totalTries)
}

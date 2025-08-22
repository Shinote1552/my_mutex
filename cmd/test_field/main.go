package main

import (
	"fmt"
	"sync"
)

func testMutex(mutex sync.Locker, name string) {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			counter++
			mutex.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("%s: counter = %d (expected: 1000)\n", name, counter)
}

func main() {
	// Протестировать все реализации
	testMutex(&Spinlock{}, "Spinlock")
	// testMutex(NewChannelMutex(), "ChannelMutex")
	// ... другие тесты
}

/*
	# BONUS LEVEL
*/
/*
 - Добавить TryLock() - неблокирующая попытка захвата
 - Реентерабельный мьютекс - с поддержкой повторного захвата
 - RW-мьютекс - с разделением на читателей и писателей
 - Таймауты - LockWithTimeout(time.Duration) bool
*/

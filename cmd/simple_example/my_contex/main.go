package main

import (
	"fmt"
	"my_mutex/internal/mycontext"
	"sync"
	"time"
)

func testWorker(ctx *mycontext.Context, wg *sync.WaitGroup, i int, resultChan chan string) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		resultChan <- fmt.Sprintf("worker%d: cancelled", i)
	case <-time.After(500 * time.Millisecond):
		resultChan <- fmt.Sprintf("worker%d: completed", i)
	}
}

func collectResults(resultChan chan string, count int) []string {
	results := make([]string, 0, count)
	for i := 0; i < count; i++ {
		results = append(results, <-resultChan)
	}
	return results
}

func main() {
	fmt.Println("=== ТЕСТИРОВАНИЕ MyContext ===")

	// Тест 1: WithCancel
	fmt.Println("\n1. Тест WithCancel:")
	var wg sync.WaitGroup
	resultChan := make(chan string, 3)

	backgroundCtx := mycontext.Background()
	ctx1, cancel1 := mycontext.WithCancel(backgroundCtx)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go testWorker(ctx1, &wg, i, resultChan)
	}

	time.Sleep(100 * time.Millisecond)
	cancel1()
	wg.Wait()
	close(resultChan)

	results1 := collectResults(resultChan, 3)
	fmt.Println("Результаты:", results1)

	// Тест 2: WithTimeout
	fmt.Println("\n2. Тест WithTimeout:")
	wg = sync.WaitGroup{}
	resultChan = make(chan string, 3)

	ctx2, cancel2 := mycontext.WithTimeout(mycontext.Background(), 200*time.Millisecond)
	defer cancel2()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go testWorker(ctx2, &wg, i, resultChan)
	}

	wg.Wait()
	close(resultChan)

	results2 := collectResults(resultChan, 3)
	fmt.Println("Результаты:", results2)

	// Тест 3: WithDeadline
	fmt.Println("\n3. Тест WithDeadline:")
	wg = sync.WaitGroup{}
	resultChan = make(chan string, 3)

	deadline := time.Now().Add(300 * time.Millisecond)
	ctx3, cancel3 := mycontext.WithDeadline(mycontext.Background(), deadline)
	defer cancel3()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go testWorker(ctx3, &wg, i, resultChan)
	}

	wg.Wait()
	close(resultChan)

	results3 := collectResults(resultChan, 3)
	fmt.Println("Результаты:", results3)

	// Тест 4: WithoutCancel
	fmt.Println("\n4. Тест WithoutCancel:")
	wg = sync.WaitGroup{}
	resultChan = make(chan string, 3)

	parentCtx, parentCancel := mycontext.WithCancel(mycontext.Background())
	isolatedCtx := mycontext.WithoutCancel(parentCtx)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go testWorker(isolatedCtx, &wg, i, resultChan)
	}

	// Отменяем родительский контекст, но изолированный должен продолжать работать
	time.Sleep(100 * time.Millisecond)
	parentCancel()
	wg.Wait()
	close(resultChan)

	results4 := collectResults(resultChan, 3)
	fmt.Println("Результаты:", results4)

	// Тест 5: Родительская отмена
	fmt.Println("\n5. Тест родительской отмены:")
	wg = sync.WaitGroup{}
	resultChan = make(chan string, 3)

	parentCtx2, parentCancel2 := mycontext.WithCancel(mycontext.Background())
	childCtx, childCancel := mycontext.WithCancel(parentCtx2)
	defer childCancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go testWorker(childCtx, &wg, i, resultChan)
	}

	time.Sleep(100 * time.Millisecond)
	parentCancel2() // Отменяем родителя - ребенок тоже должен отмениться
	wg.Wait()
	close(resultChan)

	results5 := collectResults(resultChan, 3)
	fmt.Println("Результаты:", results5)

	// Тест 6: Многократный вызов Done()
	fmt.Println("\n6. Тест многократного вызова Done():")
	ctx6, cancel6 := mycontext.WithCancel(mycontext.Background())

	// Многократно вызываем Done() - не должно быть паники
	var doneWg sync.WaitGroup
	for i := 0; i < 10; i++ {
		doneWg.Add(1)
		go func() {
			defer doneWg.Done()
			<-ctx6.Done()
		}()
	}

	time.Sleep(50 * time.Millisecond)
	cancel6()
	doneWg.Wait()

	// Пытаемся вызвать Done() после отмены
	for i := 0; i < 5; i++ {
		<-ctx6.Done()
	}
	fmt.Println("Многократный вызов Done() прошел без паники")

	// Тест 7: Nil timer безопасность
	fmt.Println("\n7. Тест безопасности с nil timer:")
	ctx7, cancel7 := mycontext.WithCancel(mycontext.Background())

	// Сначала отменяем, потом читаем
	cancel7()

	// Многократно вызываем Done() - не должно быть паники из-за nil timer
	for i := 0; i < 5; i++ {
		<-ctx7.Done()
	}
	fmt.Println("Работа с nil timer прошла без паники")

	fmt.Println("\n=== ВСЕ ТЕСТЫ ЗАВЕРШЕНЫ ===")
}

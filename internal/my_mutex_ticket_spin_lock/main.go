package main

import (
	"fmt"
	"my_mutex/internal/my_mutex_ticket_spin_lock/mymutex"
	"sync"
	"time"
)

type Counter struct {
	mu    mymutex.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type BankAccount struct {
	mu      mymutex.Mutex
	balance int
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.balance += amount
}

func (b *BankAccount) Withdraw(amount int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.balance >= amount {
		b.balance -= amount
		return true
	}
	return false
}

func (b *BankAccount) GetBalance() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.balance
}

func main() {
	fmt.Println("=== Тест 1: Базовый счетчик ===")
	testCounter()

	fmt.Println("\n=== Тест 2: Банковский счет ===")
	testBankAccount()

	fmt.Println("\n=== Тест 3: Конкурентные операции ===")
	testConcurrentOperations()
}

func testCounter() {
	var wg sync.WaitGroup
	counter := Counter{}

	// Запускаем 1000 горутин для инкремента
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Ожидаемое значение: 1000\n")
	fmt.Printf("Полученное значение: %d\n", counter.GetValue())
}

func testBankAccount() {
	account := BankAccount{balance: 1000}
	var wg sync.WaitGroup

	// Конкурентные пополнения
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(amount int) {
			defer wg.Done()
			account.Deposit(amount)
			fmt.Printf("Пополнение на %d руб.\n", amount)
		}(i * 100)
	}

	// Конкурентные снятия
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(amount int) {
			defer wg.Done()
			success := account.Withdraw(amount)
			if success {
				fmt.Printf("Снятие %d руб. успешно\n", amount)
			} else {
				fmt.Printf("Снятие %d руб. отклонено (недостаточно средств)\n", amount)
			}
		}(200)
	}

	wg.Wait()
	fmt.Printf("Итоговый баланс: %d руб.\n", account.GetBalance())
}

func testConcurrentOperations() {
	var (
		sharedData int
		mu         sync.Mutex
		wg         sync.WaitGroup
	)

	// Писатели
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				mu.Lock()
				sharedData++
				fmt.Printf("Писатель %d: записал значение %d\n", id, sharedData)
				mu.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Читатели
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 4; j++ {
				mu.Lock()
				value := sharedData
				mu.Unlock()
				fmt.Printf("Читатель %d: прочитал значение %d\n", id, value)
				time.Sleep(15 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Финальное значение sharedData: %d\n", sharedData)
}

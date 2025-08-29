package mycontext

import (
	"sync"
	"testing"
	"time"
)

func TestBackground(t *testing.T) {
	ctx := Background() // ctx is *Context

	select {
	case <-ctx.Done(): // This works with pointer receiver
		t.Error("Background context should not be done")
	default:
		// Expected
	}

	if ctx.Err() != nil { // This works with pointer receiver
		t.Errorf("Expected nil error, got %v", ctx.Err())
	}
}

func TestWithCancel(t *testing.T) {
	parent := Background()
	ctx, cancel := WithCancel(parent)

	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
	}

	cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancellation")
	}

	if ctx.Err() != Canceled {
		t.Errorf("Expected Canceled error, got %v", ctx.Err())
	}
}

func TestParentChildCancellation(t *testing.T) {
	parent, parentCancel := WithCancel(Background())
	child, childCancel := WithCancel(parent)
	defer childCancel()

	parentCancel()

	select {
	case <-parent.Done():
		// Expected
	default:
		t.Error("Parent context should be done")
	}

	select {
	case <-child.Done():
		// Expected
	default:
		t.Error("Child context should be done when parent is cancelled")
	}

	if parent.Err() != Canceled {
		t.Errorf("Expected Canceled error for parent, got %v", parent.Err())
	}

	if child.Err() != Canceled {
		t.Errorf("Expected Canceled error for child, got %v", child.Err())
	}
}

func TestWithCancel_Concurrent(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	var wg sync.WaitGroup
	const goroutines = 10

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
		}()
	}

	time.Sleep(10 * time.Millisecond)
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Error("Not all goroutines completed after cancellation")
	}
}

func TestWithDeadline_AlreadyPassed(t *testing.T) {
	deadline := time.Now().Add(-100 * time.Millisecond)
	ctx, cancel := WithDeadline(Background(), deadline)
	defer cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done immediately for past deadline")
	}

	if ctx.Err() != DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
	}
}

func TestWithoutCancel(t *testing.T) {
	parent, parentCancel := WithCancel(Background())
	isolated := WithoutCancel(parent)

	parentCancel()

	select {
	case <-parent.Done():
		// Expected
	default:
		t.Error("Parent context should be done")
	}

	select {
	case <-isolated.Done():
		t.Error("Isolated context should not be affected by parent cancellation")
		// error: Isolated context should not be affected by parent cancellation

	default:
		// Expected
	}

	if isolated.Err() != nil {
		t.Errorf("Expected nil error for isolated context, got %v", isolated.Err())
	}
}

func TestMultipleCancels(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	cancel()
	cancel()
	cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancellation")
	}
}

func TestRaceConditions(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	var wg sync.WaitGroup
	const readers = 20

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					// Continue
				}
			}
		}()
	}

	time.Sleep(5 * time.Millisecond)
	cancel()

	wg.Wait()
}

func TestDoneChannelReuse(t *testing.T) {
	ctx, cancel := WithCancel(Background())
	done := ctx.Done()

	select {
	case <-done:
		t.Error("Should not receive before cancellation")
	default:
	}

	cancel()

	for i := 0; i < 5; i++ {
		select {
		case <-done:
			// Expected
		default:
			t.Error("Should receive after cancellation")
		}
	}

	if ctx.Done() != done {
		t.Error("Multiple calls to Done() should return the same channel")
	}
}

func TestWithTimeout(t *testing.T) {
	start := time.Now()
	ctx, cancel := WithTimeout(Background(), 100*time.Millisecond)
	defer cancel()

	<-ctx.Done()
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Timeout occurred too early: %v", elapsed)
	}

	if elapsed > 500*time.Millisecond {
		t.Errorf("Timeout occurred too late: %v", elapsed)
	}

	if ctx.Err() != DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
	}
}

func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := WithDeadline(Background(), deadline)
	defer cancel()

	start := time.Now()
	<-ctx.Done()
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("Deadline occurred too early: %v", elapsed)
	}

	if elapsed > 500*time.Millisecond {
		t.Errorf("Deadline occurred too late: %v", elapsed)
	}

	if ctx.Err() != DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
	}
}

func TestWithCancel_NilParent(t *testing.T) {
	ctx, cancel := WithCancel(nil)
	defer cancel()

	if ctx == nil {
		t.Fatal("WithCancel should return non-nil context even with nil parent")
	}

	// Проверяем, что контекст работает нормально
	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Expected
	}

	cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancellation")
	}
}

func TestWithoutCancel_Isolation(t *testing.T) {
	parent, parentCancel := WithCancel(Background())
	isolated := WithoutCancel(parent)

	// Проверяем, что isolated не имеет родителя
	if isolated.parent != nil {
		t.Error("WithoutCancel should create context without parent")
	}

	parentCancel()

	// Родитель должен быть отменен
	select {
	case <-parent.Done():
		// Expected
	default:
		t.Error("Parent context should be done")
	}

	// Изолированный контекст не должен быть отменен
	select {
	case <-isolated.Done():
		t.Error("Isolated context should not be affected by parent cancellation")
	default:
		// Expected
	}

	if isolated.Err() != nil {
		t.Errorf("Expected nil error for isolated context, got %v", isolated.Err())
	}
}

func TestWithDeadline_CancellationBeforeDeadline(t *testing.T) {
	deadline := time.Now().Add(200 * time.Millisecond)
	ctx, cancel := WithDeadline(Background(), deadline)

	// Отменяем до истечения дедлайна
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancellation")
	}

	if ctx.Err() != Canceled {
		t.Errorf("Expected Canceled error, got %v", ctx.Err())
	}
}

func TestWithDeadline_MultipleCancels(t *testing.T) {
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := WithDeadline(Background(), deadline)

	cancel()
	cancel() // Многократный вызов не должен паниковать
	cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancellation")
	}
}

func TestContext_Err_RaceCondition(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	var wg sync.WaitGroup
	const readers = 20

	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = ctx.Err()
			}
		}()
	}

	time.Sleep(5 * time.Millisecond)
	cancel()

	wg.Wait()

	// После отмены Err должен возвращать Canceled
	if ctx.Err() != Canceled {
		t.Errorf("Expected Canceled error after cancellation, got %v", ctx.Err())
	}
}

func TestDone_ChannelConsistency(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	// Multiple calls to Done() should return the same channel
	done1 := ctx.Done()
	done2 := ctx.Done()

	if done1 != done2 {
		t.Error("Multiple calls to Done() should return the same channel")
	}

	cancel()

	// После отмены канал должен оставаться тем же
	done3 := ctx.Done()
	if done1 != done3 {
		t.Error("Done() should return the same channel after cancellation")
	}
}

func TestWithTimeout_ImmediateCancellation(t *testing.T) {
	ctx, cancel := WithTimeout(Background(), 0)
	defer cancel()

	// Контекст с нулевым таймаутом должен быть сразу отменен
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context with zero timeout should be done immediately")
	}

	if ctx.Err() != DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
	}
}

func TestParentCancellationPropagation(t *testing.T) {
	parent, parentCancel := WithCancel(Background())
	child, childCancel := WithCancel(parent)
	grandchild, grandchildCancel := WithCancel(child)
	defer func() {
		childCancel()
		grandchildCancel()
	}()

	parentCancel()

	// Все контексты должны быть отменены
	contexts := []*Context{parent, child, grandchild}
	for i, ctx := range contexts {
		select {
		case <-ctx.Done():
			// Expected
		default:
			t.Errorf("Context at level %d should be done after parent cancellation", i)
		}

		if ctx.Err() != Canceled {
			t.Errorf("Expected Canceled error for context at level %d, got %v", i, ctx.Err())
		}
	}
}

func TestTimerCleanup(t *testing.T) {
	deadline := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := WithDeadline(Background(), deadline)

	// Ждем истечения таймера
	time.Sleep(100 * time.Millisecond)

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after deadline")
	}

	if ctx.Err() != DeadlineExceeded {
		t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
	}

	cancel() // Не должно паниковать после автоматической отмены
}

func TestConcurrentCancellationAndDone(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	var wg sync.WaitGroup
	const goroutines = 10

	// Запускаем горутины, которые читают из Done()
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
		}()
	}

	// И горутины, которые вызывают Err()
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = ctx.Err()
				time.Sleep(time.Microsecond)
			}
		}()
	}

	time.Sleep(5 * time.Millisecond)
	cancel()

	wg.Wait()
}

func TestBackgroundContextProperties(t *testing.T) {
	ctx := Background()

	// Background context should have no parent
	if ctx.parent != nil {
		t.Error("Background context should not have a parent")
	}

	// Background context should have no timer
	if ctx.timer != nil {
		t.Error("Background context should not have a timer")
	}

	// Background context should not be done
	select {
	case <-ctx.Done():
		t.Error("Background context should not be done")
	default:
		// Expected
	}

	if ctx.Err() != nil {
		t.Errorf("Background context should have nil error, got %v", ctx.Err())
	}
}

func TestWithCancel_NoDataRace(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	// Concurrent access to Done() and Err()
	var wg sync.WaitGroup
	const accesses = 1000

	for i := 0; i < accesses; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = ctx.Err()
		}()
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
			default:
			}
		}()
	}

	cancel()
	wg.Wait()
}

func TestWithDeadline_NilParent(t *testing.T) {
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := WithDeadline(nil, deadline)
	defer cancel()

	if ctx == nil {
		t.Fatal("WithDeadline should return non-nil context even with nil parent")
	}

	// Контекст должен нормально работать
	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Expected
	}
}

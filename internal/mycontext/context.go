package mycontext

import (
	"errors"
	"sync"
	"time"
)

var (
	Canceled         = errors.New("context canceled")
	DeadlineExceeded = errors.New("context deadline exceeded")
)

/*
Это простая реализация контекста в учебных целях, здесь 100% есть ошибки.
Будет работать и моими реализациями Мютекса
*/
type Context struct {
	done   chan struct{}
	parent *Context
	timer  *time.Timer
	mu     sync.Mutex
	err    error
}

// Общая функция для безопасной отмены контекста
func (mc *Context) safeCancel(err error) {

	mc.mu.Lock()
	defer mc.mu.Unlock()
	select {
	case <-mc.done:
		return
	default:
		mc.err = err
		close(mc.done)
	}
}

// В оригинале принимает интерфейс, но я пока одной структурой обошелся
func WithCancel(parent *Context) (*Context, func()) {
	child := Context{
		done:   make(chan struct{}),
		parent: parent,
	}

	cancel := func() {
		child.safeCancel(Canceled)
	}
	return &child, cancel
}

func WithoutCancel(parent *Context) *Context {
	child := Context{
		done: make(chan struct{}),
	}

	return &child
}

func WithDeadline(parent *Context, ddl time.Time) (*Context, func()) {
	now := time.Now()
	if now.After(ddl) {
		child := &Context{
			done: make(chan struct{}),
			err:  DeadlineExceeded,
		}
		close(child.done)
		return child, func() {}
	}

	child := &Context{
		done:   make(chan struct{}),
		parent: parent,
		timer:  time.NewTimer(time.Until(ddl)),
	}

	cancel := func() {
		child.mu.Lock()
		defer child.mu.Unlock()
		select {
		case <-child.done:
		default:
			if child.timer != nil {
				child.timer.Stop()
			}
			child.err = Canceled
			close(child.done)
		}
	}

	/*
		Запускаем мониторинг таймера в отдельной горутине,
		как минимум для асинхронной рабоыт таймера
	*/
	go func() {
		select {
		case <-child.timer.C:
			child.mu.Lock()
			select {
			case <-child.done:
			default:
				child.err = DeadlineExceeded
				close(child.done)
			}
			child.mu.Unlock()
		case <-child.done:
		}
	}()

	return child, cancel
}

func WithTimeout(parent *Context, duration time.Duration) (*Context, func()) {
	return WithDeadline(parent, time.Now().Add(duration))
}

func Background() *Context {
	return &Context{
		done: make(chan struct{}),
	}
}

// Проверка отменен ли контекст, если да то есть ли ошибка
func (mc *Context) Err() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	select {
	case <-mc.done:
		return mc.err
	default:
	}

	if mc.parent != nil {
		if err := mc.parent.Err(); err != nil {
			return err
		}
	}

	return nil
}

// Общая функция для обработки родительской отмены
func (mc *Context) handleParentCancellation() bool {
	if mc.parent != nil {
		select {
		case <-mc.parent.Done():
			mc.mu.Lock()
			defer mc.mu.Unlock()
			select {
			case <-mc.done:
				return true
			default:
				if parentErr := mc.parent.Err(); parentErr != nil && mc.err == nil {
					mc.err = parentErr
				}
				close(mc.done)
				return true
			}
		default:
		}
	}
	return false
}

// Общая функция для обработки таймера
func (mc *Context) handleTimerCancellation() bool {
	if mc.timer != nil {
		select {
		case <-mc.timer.C:
			mc.mu.Lock()
			defer mc.mu.Unlock()
			select {
			case <-mc.done:
				return true
			default:
				if mc.err == nil {
					mc.err = DeadlineExceeded
				}
				close(mc.done)
				return true
			}
		default:
		}
	}
	return false
}

// Жертва единой структуры для разных контекстов) тут datarace, livelock можно найти
func (mc *Context) Done() <-chan struct{} {
	mc.mu.Lock()
	select {
	case <-mc.done:
		mc.mu.Unlock()
		return mc.done
	default:
		mc.mu.Unlock()
	}

	if mc.handleParentCancellation() {
		return mc.done
	}

	if mc.handleTimerCancellation() {
		return mc.done
	}

	return mc.done
}

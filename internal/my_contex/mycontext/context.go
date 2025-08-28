package mycontext

import (
	"time"
)

type Context struct {
	done   chan struct{}
	parent *Context
	timer  *time.Timer
}

// в оригинале принимает интерфейс, но я пока одной структурой обошелся
func WithCancel(parent Context) (Context, func()) {
	child := Context{
		done:   make(chan struct{}),
		parent: &parent,
		timer:  nil}

	cancel := func() {
		close(child.done)
	}
	return child, cancel
}

func WithoutCancel(parent Context) Context {
	child := Context{
		done:   make(chan struct{}),
		parent: nil,
		timer:  nil}

	return child
}

func WithDeadline(parent Context, ddl time.Time) (Context, func()) {
	timer := time.NewTimer(time.Until(ddl))

	child := Context{
		done:   make(chan struct{}),
		parent: &parent,
		timer:  timer}

	cancel := func() {
		close(child.done)
	}
	return child, cancel
}

func WithTimeout(parent Context, duration time.Duration) (Context, func()) {
	timer := time.NewTimer(duration)

	child := Context{
		done:   make(chan struct{}),
		parent: &parent,
		timer:  timer}

	cancel := func() {
		close(child.done)
	}
	return child, cancel
}

func Background() Context {
	return Context{
		done:   make(chan struct{}),
		parent: nil,
	}
}

func (mc *Context) Done() <-chan struct{} {
	select {
	case <-mc.done:
		return mc.done
	default:
		if mc.parent != nil {
			select {
			case <-mc.parent.Done():
				close(mc.done)
			default:
				if mc.timer != nil {
					select {
					case <-mc.timer.C:
						close(mc.done)
					default:
					}
				}
			}
		}
	}

	return mc.done
}

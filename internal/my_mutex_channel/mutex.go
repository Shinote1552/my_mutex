package mymutexchannel

// Реализовать мьютекс с FIFO очередью на каналах
type Mutex struct {
	chanQueue chan struct{}
}

func NewMutex() *Mutex {
	// TODO: инициализация
}

func (m *Mutex) Lock() {
	// TODO: встать в очередь
}

func (m *Mutex) Unlock() {
	// TODO: разбудить следующего
}

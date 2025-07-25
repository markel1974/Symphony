package gl_render

import (
	"sync"
	"time"
)

type Ticker struct {
	queue    chan int
	lock     sync.Mutex
	emulate  func()
	chunk    int
	interval time.Duration
	cycles   int
}

func NewTicker(emulate func()) *Ticker {
	const mhz = 1000
	const interval = 20
	t := &Ticker{
		emulate:  emulate,
		interval: interval * time.Millisecond,
		cycles:   interval * mhz,
		queue:    make(chan int, 4096),
		chunk:    20, //runtime.NumCPU(),
	}
	return t
}

func (t *Ticker) loop() {
	for {
		select {
		case cycles := <-t.queue:
			for step := 0; step < cycles; step++ {
				t.emulate()
			}
		}
	}
}

func (t *Ticker) Start() {
	go t.loop()

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	for range ticker.C {
		cycles := t.cycles / t.chunk
		for x := 0; x < t.chunk; x++ {
			t.queue <- cycles
		}
	}
}

package benchmark

/*
#include <unistd.h>
//#include "_cgo_export.h"
void cticker() {
    //   usleep(20000); // 20 ms
	//   Gotask();
}
*/
//
import (
	"sync"
	"time"
)

type Ticker struct {
	queue     chan int
	lock      sync.Mutex
	emulateFn func()
	interval  time.Duration
	blocks    int
	cycles    int
}

func NewTicker(hz int, intervalMs int, blocks int, emulate func()) *Ticker {
	t := &Ticker{
		emulateFn: emulate,
		interval:  time.Duration(intervalMs) * time.Millisecond,
		blocks:    blocks,
		queue:     make(chan int, 4096),
		cycles:    (hz / intervalMs) / blocks, //runtime.NumCPU(),
	}
	return t
}

func (t *Ticker) loop() {
	for {
		select {
		case cycles := <-t.queue:
			for step := 0; step < cycles; step++ {
				t.emulateFn()
			}
		}
	}
}

//go:noinline
//go:nosplit
func (t *Ticker) Emulate2() {

}

//go:noinline
//go:nosplit
func Emulate() {
	//t.emulate()
}

func (t *Ticker) Start() {
	//debug.SetGCPercent(-1)
	//emulate := Emulate //t.emulate
	//emulate := t.Emulate
	//emulate := Emulate
	//emulate := (*Ticker).Emulate2
	for {
		for block := 0; block < t.blocks; block++ {
			for cycle := 0; cycle < t.cycles; cycle++ {
				//emulate(t)
				t.emulateFn()
			}
			//runtime.Gosched()
		}
		//C.cticker()
		time.Sleep(time.Duration(20) * time.Millisecond)
	}
	/*
		go t.loop()

		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()
		for range ticker.C {
			for x := 0; x < t.blocks; x++ {
				t.queue <- t.cycles
			}
		}

	*/
}

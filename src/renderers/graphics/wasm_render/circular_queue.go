// file: wasm_render/circular_queue_wasm.go

package wasm_render

import (
	"fmt"
	"sync/atomic"
)

const (
	// Definiamo gli indici nel blocco di controllo (un array di uint32)
	headControlIndex = 0 // Indice di scrittura (head)
	tailControlIndex = 1 // Indice di lettura (tail)
)

// CircularQueueWasm gestisce la logica del ring buffer sulla memoria condivisa.
type CircularQueueWasm struct {
	control   *[]uint32  // Puntatore alla sezione di controllo (head/tail) della SharedArrayBuffer
	data      *[]float32 // Puntatore alla sezione dati della SharedArrayBuffer
	capacity  uint32     // Capacità totale del buffer in chunk
	chunkSize int        // Dimensione di un singolo chunk in campioni
}

func NewCircularQueueWasm(control *[]uint32, data *[]float32, capacity uint32, chunkSize int) *CircularQueueWasm {
	return &CircularQueueWasm{
		control:   control,
		data:      data,
		capacity:  capacity,
		chunkSize: chunkSize,
	}
}

// Push aggiunge un chunk nel buffer condiviso. È chiamata dal produttore (Go).
func (q *CircularQueueWasm) Push(chunk *[]float32) bool {
	head := atomic.LoadUint32(&(*q.control)[headControlIndex])
	tail := atomic.LoadUint32(&(*q.control)[tailControlIndex])
	nextHead := (head + 1) % q.capacity
	if nextHead == tail {
		fmt.Println("WASM Audio: Ring buffer full, dropping frame!")
		return false // Buffer pieno
	}
	offset := int(head) * q.chunkSize
	copy((*q.data)[offset:], *chunk)
	// Aggiorna l'indice di testa in modo atomico, rendendolo visibile a JS
	atomic.StoreUint32(&(*q.control)[headControlIndex], nextHead)
	return true
}

package oto_render

import (
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"sync"
	"time"
)

const (
	bufferStateEmpty   = 0 // La coda è completamente vuota.
	bufferStateTooFast = 1 // Il consumatore è troppo veloce, la coda si sta svuotando.
	bufferStateNice    = 2 // Stato di equilibrio ideale.
	bufferStateGood    = 3
	bufferStateStable  = 4 // Stato ancora stabile, agisce come "dead zone" per evitare oscillazioni.
)

// ContinuousReader manages continuous audio streaming to an 'oto' player.
// It uses a buffered circular queue and a stateful resampling algorithm
// to dynamically synchronize a free-running producer with a real-time
// consumer (the sound card), ensuring stable and high-quality audio flow.
type ContinuousReader struct {
	player       oto.Player
	writeFn      writeFn
	ring         *CircularQueue
	interpolator *LinearInterpolation
	chunkSize    int
	doubleBuffer *[]float32
	states       [0xf]func([]byte) (int, error)
	lock         sync.RWMutex
}

// NewContinuousReader creates and initializes a new instance of the ContinuousReader for managing continuous audio streams.
func NewContinuousReader() *ContinuousReader {
	return &ContinuousReader{}
}

// Setup initializes the ContinuousReader instance with the specified sample rate, chunks per second, channels, and format.
func (r *ContinuousReader) Setup(sampleRate int, chunkPerSecond int, channels int, fo string) error {
	format, ok := _formats[fo]
	if !ok {
		return fmt.Errorf("audio format not found")
	}
	r.chunkSize = sampleRate / chunkPerSecond
	bufferSize := r.chunkSize * format.Bytes
	options := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channels,
		Format:       format.Format,
		BufferSize:   time.Second / time.Duration(chunkPerSecond),
	}
	ctx, ready, err := oto.NewContextWithOptions(options)
	if err != nil {
		return fmt.Errorf("failed to create oto context: %w", err)
	}
	<-ready
	r.ring = NewCircularQueue(r.chunkSize)
	r.interpolator = NewLinearInterpolation(r.chunkSize * 2)
	doubleBuffer := make([]float32, r.chunkSize*2)
	r.doubleBuffer = &doubleBuffer
	r.writeFn = format.Func
	r.player = ctx.NewPlayer(r)
	r.player.SetVolume(1.0)
	r.player.(oto.BufferSizeSetter).SetBufferSize(bufferSize)

	for x := range r.states {
		r.states[x] = r.handleTooSlow
	}
	r.states[bufferStateEmpty] = r.handleEmpty
	r.states[bufferStateTooFast] = r.handleTooFast
	r.states[bufferStateNice] = r.handleGood
	r.states[bufferStateGood] = r.handleGood
	r.states[bufferStateStable] = r.handleGood
	return nil
}

// Play starts the playback of the audio stream using the underlying oto player.
func (r *ContinuousReader) Play() {
	r.player.Play()
}

// Err returns the current error state of the underlying oto player, if any.
func (r *ContinuousReader) Err() error {
	return r.player.Err()
}

// AddChunk adds a chunk of audio data to the circular queue for processing and playback, locking the queue during the operation.
func (r *ContinuousReader) AddChunk(chunk *[]float32, samples int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.ring.Push(chunk)
}

// Read is the core of the renderer. It is called by 'oto' when it needs data.
// Implements a state machine to decide how to process audio data
// based on the buffer fill level.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.states[r.ring.Counter()&0xf](buf)
}

// handleEmpty is called when the buffer is empty. It fills the 'oto' buffer
// with silence to keep the audio stream active and prevent interruptions.
func (r *ContinuousReader) handleEmpty(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}

// handleGood is called when the buffer is in a balanced state (2 or 3 chunks).
// It takes a chunk and plays it as-is, without resampling.
func (r *ContinuousReader) handleGood(buf []byte) (int, error) {
	chunkToPlay, _ := r.ring.Pop()
	written := 0
	for _, sample := range *chunkToPlay {
		written += r.writeFn(buf, sample, written)
	}
	return written, nil
}

// handleTooFast is called when the buffer is running low (1 chunk remaining).
// This means the consumer (oto) is faster than the producer.
// To "buy time", it takes a chunk, "stretches" it to double its length,
// plays the first half, and puts the second half back in the queue.
func (r *ContinuousReader) handleTooFast(buf []byte) (int, error) {
	//fmt.Printf("[%s] CONSUMER TOO FAST: Stretching...\n", time.Now().Format(time.RFC3339Nano))
	chunkToStretch, ok := r.ring.Pop()
	if !ok {
		return r.handleEmpty(buf)
	}
	stretchedChunk, _ := r.interpolator.LinearInterpolationF32(chunkToStretch, r.chunkSize*2)
	firstHalf := (*stretchedChunk)[:r.chunkSize]
	secondHalf := (*stretchedChunk)[r.chunkSize:]
	r.ring.Push(&secondHalf)
	written := 0
	for _, sample := range firstHalf {
		written += r.writeFn(buf, sample, written)
	}
	return written, nil
}

// handleTooSlow is called when the buffer is filling up (more than 3 chunks).
// This means the producer (emulator) is faster than the consumer.
// To "catch up", it takes two chunks, "compresses" them into one
// and plays it, consuming data from the buffer at twice the speed.
func (r *ContinuousReader) handleTooSlow(buf []byte) (int, error) {
	//fmt.Printf("[%s] CONSUMER TOO SLOW (Lag > %d chunks): Squishing...\n", time.Now().Format(time.RFC3339Nano), r.ring.Counter())
	chunk1, chunk2, ok := r.ring.DoublePop()
	if !ok {
		return r.handleEmpty(buf)
	}
	copy(*r.doubleBuffer, *chunk1)
	copy((*r.doubleBuffer)[len(*chunk1):], *chunk2)
	squishedChunk, squishedLen := r.interpolator.LinearInterpolationF32(r.doubleBuffer, r.chunkSize)
	written := 0
	for i := 0; i < squishedLen; i++ {
		written += r.writeFn(buf, (*squishedChunk)[i], written)
	}
	return written, nil
}

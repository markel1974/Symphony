package oto_render

import (
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"math"
	"time"
)

// ContinuousReader is a struct for managing continuous audio playback using a circular buffer and linear interpolation.
type ContinuousReader struct {
	player        oto.Player
	bytes         int
	writeFn       writeFn
	ring          *CircularQueue
	interpolator  *LinearInterpolation
	maxCorrection float64
	remainder     *[]float32
	remainderLen  int
}

// NewContinuousReader creates and initializes a new instance of ContinuousReader for managing continuous audio playback.
func NewContinuousReader() *ContinuousReader {
	r := &ContinuousReader{}
	return r
}

// Setup initializes the ContinuousReader with the specified sample rate, chunk rate, channels, and audio format.
// It configures audio playback via oto.Player, prepares internal buffers, and initializes interpolation and queue systems.
// Returns an error if the audio format is invalid or the oto context setup fails.
func (r *ContinuousReader) Setup(sampleRate int, chunkPerSecond int, channels int, fo string) error {
	format, ok := _formats[fo]
	if !ok {
		return fmt.Errorf("audio format not found")
	}
	chunkSize := sampleRate / chunkPerSecond
	options := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channels,
		Format:       format.Format,
		BufferSize:   time.Second / time.Duration(chunkPerSecond),
	}
	ctx, ready, err := oto.NewContextWithOptions(options)

	//ctx, ready, err := oto.NewContext(sampleRate, channels, format.Format)
	if err != nil {
		return fmt.Errorf("failed to create oto context: %w", err)
	}
	<-ready
	readyTarget := (chunkPerSecond * 10) / 100
	r.ring = NewCircularQueue(chunkPerSecond, chunkSize, readyTarget)
	r.interpolator = NewLinearInterpolation(chunkSize + 1)
	r.bytes = format.Bytes
	r.writeFn = format.Func
	r.player = ctx.NewPlayer(r)
	r.player.SetVolume(1.0)
	r.maxCorrection = (float64(chunkSize) * 10) / 100
	r.remainderLen = 0
	remainder := make([]float32, chunkPerSecond)
	r.remainder = &remainder
	//bufferSize := chunkLen * r.bytes
	//r.player.(oto.BufferSizeSetter).SetBufferSize(bufferSize)
	return nil
}

// Play starts audio playback using the underlying oto.Player associated with the ContinuousReader instance.
func (r *ContinuousReader) Play() {
	r.player.Play()
}

// Err returns the last error encountered during the continuous audio playback. If no error occurred, it returns nil.
func (r *ContinuousReader) Err() error {
	return r.player.Err()
}

// AddChunk enqueues a chunk of audio samples into the circular buffer for playback and processing.
func (r *ContinuousReader) AddChunk(chunk *[]float32, samples int) {
	r.ring.Push(chunk)
}

// remainCreate updates the remainder buffer with leftover samples from finalChunk that were not written.
func (r *ContinuousReader) remainCreate(finalChunk *[]float32, finalChunkSize int, samplesToWrite int) {
	r.remainderLen = finalChunkSize - samplesToWrite
	if r.remainderLen > len(*r.remainder) {
		r.remainderLen = len(*r.remainder)
	}
	copy((*r.remainder)[:r.remainderLen], (*finalChunk)[samplesToWrite:])
}

// remainFlush transfers the remaining samples in the buffer using the remainder data and updates the remainder state.
func (r *ContinuousReader) remainFlush(buf []uint8) int {
	samplesThatFitInBuf := len(buf) / r.bytes
	samplesToWrite := r.remainderLen
	if samplesToWrite > samplesThatFitInBuf {
		samplesToWrite = samplesThatFitInBuf
	}
	written := 0
	for i := 0; i < samplesToWrite; i++ {
		written += r.writeFn(buf, (*r.remainder)[i], written)
	}
	remaining := r.remainderLen - samplesToWrite
	if remaining > 0 {
		copy((*r.remainder)[:remaining], (*r.remainder)[samplesToWrite:r.remainderLen])
	}
	r.remainderLen = remaining
	return written
}

// emptyFlush clears the provided buffer by setting all its elements to 0 and returns the number of bytes cleared.
func (r *ContinuousReader) emptyFlush(buf []uint8) int {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf)
}

// resampling adjusts the audio chunk size based on the buffer fill level, ensuring consistent audio playback.
// Returns the adjusted audio chunk and its size.
func (r *ContinuousReader) resampling(originalChunkPtr *[]float32) (*[]float32, int) {
	const targetRatio = 0.5 //50%
	originalChunk := originalChunkPtr
	originalChunkSize := len(*originalChunk)
	finalChunk := originalChunk
	finalChunkSize := originalChunkSize
	fillRatio := r.ring.FillRatio()
	// Calcoliamo l'errore: un valore positivo se il buffer è troppo pieno, negativo se troppo vuoto.
	errorRatio := fillRatio - targetRatio
	// Calcoliamo la correzione in modo proporzionale all'errore.
	// L'errore massimo è 0.5 (100% - 50%), quindi moltiplichiamo per 2 per normalizzare.
	correction := int(math.Round(errorRatio * 2 * r.maxCorrection))
	targetSize := originalChunkSize - correction
	if targetSize < originalChunkSize-int(r.maxCorrection) {
		targetSize = originalChunkSize - int(r.maxCorrection)
	}
	if targetSize > originalChunkSize+int(r.maxCorrection) {
		targetSize = originalChunkSize + int(r.maxCorrection)
	}
	if targetSize != originalChunkSize && targetSize > 0 {
		finalChunk, finalChunkSize = r.interpolator.LinearInterpolationF32(originalChunk, targetSize)
	}
	return finalChunk, finalChunkSize
}

// Read reads data into the provided buffer, processes audio chunks, and handles remainder audio samples for playback buffer.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	if r.remainderLen > 0 {
		return r.remainFlush(buf[:]), nil
	}
	samples, ok := r.ring.Pop()
	if !ok {
		return r.emptyFlush(buf[:]), nil
	}
	finalSamples, finalSampleSize := r.resampling(samples)
	sampleSize := finalSampleSize
	if requiredSampleLen := len(buf) / r.bytes; sampleSize > requiredSampleLen {
		sampleSize = requiredSampleLen
		if finalSampleSize > sampleSize {
			r.remainCreate(finalSamples, finalSampleSize, sampleSize)
		}
	}
	written := 0
	for i := 0; i < sampleSize; i++ {
		written += r.writeFn(buf, (*finalSamples)[i], written)
	}
	return written, nil
}

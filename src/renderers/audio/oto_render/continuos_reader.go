package oto_render

import (
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"math"
	"time"
)

// offset16 represents a constant offset value of 16 bits.
//
// offset32 represents a constant offset value of 32 bits.
const (
	offset16 = 2
	offset32 = 4
)

// writeFn is a function type that writes data from a buffer, processes it with a uint32 value and an index, and returns an int.
type writeFn func(buf []byte, data float32, idx int) int

// Format represents a structure that defines an audio format with specific encoding, byte size, and write function.
type Format struct {
	Format int
	Bytes  int
	Func   writeFn
}

// _formats is a map that associates audio format names with their respective Format structures, defining format details.
var _formats = map[string]Format{
	"PCM16":        {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LE},
	"PCM16Mod":     {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LEMod},
	"FLOAT32LE":    {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LE},
	"FLOAT32LEMOD": {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LEMod},
}

// ContinuousReader is a type for continuous audio data streaming and playback management.
// It uses a mutex for thread-safe operations on audio chunks and handles data writing via a custom write function.
// The type integrates with an audio player to support playback functionality.
type ContinuousReader struct {
	player        oto.Player
	bytes         int
	writeFn       writeFn
	ring          *CircularQueue
	interpolator  *LinearInterpolation
	maxCorrection float64
}

// NewContinuousReader initializes a new ContinuousReader for streaming audio based on the specified sample rate.
// Returns a pointer to the ContinuousReader instance and an error if the initialization fails.
func NewContinuousReader() *ContinuousReader {
	r := &ContinuousReader{}
	return r
}

// Setup initializes the ContinuousReader with a given sample rate, channel count, and audio format. Returns an error if failed.
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
	//bufferSize := chunkLen * r.bytes
	//r.player.(oto.BufferSizeSetter).SetBufferSize(bufferSize)
	return nil
}

// Play starts or resumes audio playback using the underlying oto.Player.
func (r *ContinuousReader) Play() {
	r.player.Play()
}

// Err returns the current error state of the audio player associated with the ContinuousReader.
func (r *ContinuousReader) Err() error {
	return r.player.Err()
}

// AddChunk appends a new chunk of audio data to the buffer and updates the sample count for playback synchronization.
func (r *ContinuousReader) AddChunk(chunk *[]float32, samples int) {
	r.ring.Push(chunk)
}

// Read processes audio data from the CircularQueue, dynamically adjusts its size using interpolation, and writes to the buffer.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	const targetRatio = 0.5 //50%
	originalChunkPtr, ok := r.ring.Pop()
	if !ok {
		for i := range buf {
			buf[i] = 0
		}
		return len(buf), nil
	}
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
		resampledBufferPtr, validLen := r.interpolator.LinearInterpolationF32(originalChunk, targetSize)
		finalChunk = resampledBufferPtr
		finalChunkSize = validLen
	}

	written := 0
	for idx := 0; idx < finalChunkSize; idx++ {
		if written+r.bytes > len(buf) {
			break
		}
		written += r.writeFn(buf, (*finalChunk)[idx], written)
	}
	return written, nil
}

// writeSignedInt16LE writes a signed 16-bit integer in little-endian format to the provided buffer at the specified index.
// It returns the constant offset16, representing the number of bytes written.
func writeSignedInt16LE(buf []byte, data float32, idx int) int {
	v := Float32ToPCM16(data)
	slice := buf[idx : idx+offset16]
	slice[0] = byte(v)
	slice[1] = byte(v >> 8)
	return offset16
}

// writeSignedInt16LEMod writes a modified signed 16-bit integer in little-endian format to the specified buffer at given index.
// It scales the input based on a divisor, clamps the result, converts it to PCM16, and sets it into the buffer.
// Returns the byte offset increment for the next write operation (`offset16`).
func writeSignedInt16LEMod(buf []byte, data float32, idx int) int {
	v := clamp(data * 10)
	t := Float32ToPCM16(v)
	//binary.LittleEndian.PutUint16(buf[start:end], t)
	slice := buf[idx : idx+offset16]
	slice[0] = byte(t)
	slice[1] = byte(t >> 8)
	return offset16
}

// writeFloat32LE writes a float32 value in little-endian format derived from the provided uint32 data to the buffer.
// buf is the target byte slice, data is the uint32 input, and idx is the offset index in the buffer.
// Returns the number of bytes written, which is fixed at 4 (offset32).
func writeFloat32LE(buf []byte, data float32, idx int) int {
	//const divisor = float32(1 << 15)
	//v := float32(int32(data)) / divisor
	//t := math.Float32bits(v)
	//if idx+offset32 >= len(buf)+4 {
	//	return 0
	//}

	t := math.Float32bits(data)
	slice := buf[idx : idx+offset32]

	slice[0] = byte(t)
	slice[1] = byte(t >> 8)
	slice[2] = byte(t >> 16)
	slice[3] = byte(t >> 24)
	return offset32
}

// writeFloat32LEMod writes a float32 value, in little-endian format, modified by a divisor and clamped to [-1, 1], to the buffer.
// buf is the byte buffer where data is written.
// data is a uint32 input that gets converted into a float32 value.
// idx specifies the starting index in the buffer where data will be written.
// Returns the number of bytes written, determined by offset32 (4 bytes).
func writeFloat32LEMod(buf []byte, data float32, idx int) int {
	v := clamp(data * 10)
	t := math.Float32bits(v)
	//binary.LittleEndian.PutUint32(buf[start:end], math.Float32bits(v))
	slice := buf[idx : idx+offset32]
	slice[0] = byte(t)
	slice[1] = byte(t >> 8)
	slice[2] = byte(t >> 16)
	slice[3] = byte(t >> 24)
	return offset32
}

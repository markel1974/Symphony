package oto_render

import (
	"encoding/binary"
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"math"
	"sync"
)

const (
	offset16 = 2
	offset32 = 4
)

type writeFn func(buf []byte, data uint32, idx int) (int, bool)

type Format struct {
	Format int
	Bytes  int
	Func   writeFn
}

var _formats = map[string]Format{
	"PCM16":        {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LE},
	"PCM16Mod":     {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LEMod},
	"FLOAT32LE":    {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LE},
	"FLOAT32LEMOD": {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LEMod},
}

type ContinuousReader struct {
	lock             sync.Mutex
	lastChunk        []uint32
	lastChunkSamples int
	player           oto.Player
	bytes            int
	writeFn          writeFn
}

func NewContinuousReader(audioSampleRate int) (*ContinuousReader, error) {
	format, ok := _formats["PCM16"]
	if !ok {
		return nil, fmt.Errorf("audio format not found")
	}
	ctx, ready, err := oto.NewContext(audioSampleRate, 1, format.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to create oto context: %w", err)
	}
	<-ready
	r := &ContinuousReader{
		bytes:   format.Bytes,
		writeFn: format.Func,
	}
	r.player = ctx.NewPlayer(r)
	r.player.SetVolume(1.0)
	return r, nil
}

func (r *ContinuousReader) Play() {
	r.player.Play()
}

func (r *ContinuousReader) Err() error {
	return r.player.Err()
}

func (r *ContinuousReader) AddChunk(chunk []uint32, samples int) {
	chunkLen := len(chunk)
	if chunkLen == 0 {
		return
	}
	lastChunkSamples := samples / 2
	r.lock.Lock()
	defer r.lock.Unlock()
	if chunkLen != len(r.lastChunk) {
		r.player.(oto.BufferSizeSetter).SetBufferSize((chunkLen / 2) * r.bytes)
		r.lastChunk = make([]uint32, chunkLen)
	}
	copy(r.lastChunk, chunk)
	r.lastChunkSamples = lastChunkSamples
}

func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.lastChunk == nil || r.lastChunkSamples == 0 {
		return 0, nil
	}
	written := 0
	for x := 0; x < r.lastChunkSamples; x++ {
		v, ok := r.writeFn(buf, r.lastChunk[x], written)
		if !ok {
			break
		}
		written += v
	}
	return written, nil
}

func writeSignedInt16LE(buf []byte, data uint32, idx int) (int, bool) {
	start := idx
	end := start + offset16
	if end > len(buf) {
		return idx, false
	}
	v := uint16(int32(data))
	binary.LittleEndian.PutUint16(buf[start:end], v)
	return offset16, true
}

func writeSignedInt16LEMod(buf []byte, data uint32, idx int) (int, bool) {
	const divisor = float32(1 << 21)
	start := idx
	end := start + offset16
	if end > len(buf) {
		return idx, false
	}
	v := float32(int32(data)) / divisor
	v = clamp(v * 10)
	binary.LittleEndian.PutUint16(buf[start:end], Float32ToPCM16(v))
	return offset16, true
}

func writeFloat32LE(buf []byte, data uint32, idx int) (int, bool) {
	const divisor = float32(1 << 15)
	start := idx
	end := start + offset32
	if end > len(buf) {
		return 0, false
	}
	v := float32(int32(data)) / divisor
	binary.LittleEndian.PutUint32(buf[start:end], math.Float32bits(v))
	return offset32, true
}

func writeFloat32LEMod(buf []byte, data uint32, idx int) (int, bool) {
	const divisor = float32(1 << 21)
	start := idx
	end := start + offset32
	if end > len(buf) {
		return 0, false
	}
	v := float32(int32(data)) / divisor
	v = clamp(v * 10)
	binary.LittleEndian.PutUint32(buf[start:end], math.Float32bits(v))
	return offset32, true
}

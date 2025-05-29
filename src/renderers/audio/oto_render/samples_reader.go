package oto_render

import (
	"encoding/binary"
	"fmt"
	"math"
)

// SamplesReader implements io.Reader for oto to consume
// uint32 samples that represent float32 bit patterns.
type SamplesReader struct {
	samples []uint32
	pos     int
}

func NewSamplesReader(samples []uint32) *SamplesReader {
	return &SamplesReader{
		samples: samples,
		pos:     0,
	}
}
func (r *SamplesReader) Read(buf []byte) (n int, err error) {
	if r.pos >= len(r.samples) {
		return 0, fmt.Errorf("end of samples")
	}

	numFloatsToRead := len(buf) / 4 // Each float32 is 4 bytes
	if r.pos+numFloatsToRead > len(r.samples) {
		numFloatsToRead = len(r.samples) - r.pos
	}

	for i := 0; i < numFloatsToRead; i++ {
		uint32bits := r.samples[r.pos]
		float32Val := math.Float32frombits(uint32bits)
		binary.LittleEndian.PutUint32(buf[i*4:(i+1)*4], math.Float32bits(float32Val))
		r.pos++
	}

	return numFloatsToRead * 4, nil
}

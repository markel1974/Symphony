package oto_render

import (
	"github.com/hajimehoshi/oto/v2"
	"math"
)

// offset16 represents a constant value of 2, potentially used as an offset for 16-bit operations.
// offset32 represents a constant value of 4, potentially used as an offset for 32-bit operations.
const (
	offset16 = 2
	offset32 = 4
)

// writeFn is a function type that writes data of a float32 value into a byte buffer at a specific index and returns the number of bytes written.
//type writeFn func(buf []byte, data float32, idx int) int

// writeFn defines a function type used to write data from a float32 buffer to a byte buffer with a specified length.
type writeFn func(buf *[]byte, data *[]float32, l int) int

// Format represents audio data processing parameters.
// It includes the format type, byte size, and a function for writing data.
type Format struct {
	Format int
	Bytes  int
	Func   writeFn
}

// _formats maps format names to their respective Format structs, which define audio format details and processing functions.
var _formats = map[string]Format{
	"PCM16":     {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LE},
	"FLOAT32LE": {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LE},
}

// writeFloat32LE writes a slice of float32 values into the buffer as little-endian 32-bit floats and returns the written size.
func writeFloat32LE(buf *[]byte, data *[]float32, l int) int {
	offset := 0
	realen := l * offset32
	if realen > len(*buf) {
		l = len(*buf)
	}
	for idx := 0; idx < l; idx++ {
		t := math.Float32bits((*data)[idx])
		(*buf)[offset] = byte(t)
		offset++
		(*buf)[offset] = byte(t >> 8)
		offset++
		(*buf)[offset] = byte(t >> 16)
		offset++
		(*buf)[offset] = byte(t >> 24)
		offset++
	}
	return offset
}

// writeSignedInt16LE converts a slice of float32 to PCM16 format and writes it to the provided byte buffer in little-endian order.
// buf is the destination byte buffer, data is the source float32 slice, and l is the number of samples to convert and write.
// Returns the total number of bytes written.
func writeSignedInt16LE(buf *[]byte, data *[]float32, l int) int {
	offset := 0
	realen := l * offset16
	if realen > len(*buf) {
		l = len(*buf)
	}
	for idx := 0; idx < l; idx++ {
		t := Float32ToPCM16((*data)[idx])
		(*buf)[offset] = byte(t)
		offset++
		(*buf)[offset] = byte(t >> 8)
		offset++
	}
	return offset
}

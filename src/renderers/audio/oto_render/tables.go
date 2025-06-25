package oto_render

import (
	"github.com/hajimehoshi/oto/v2"
	"math"
)

// offset16 represents a constant offset value of 16 bits.
//
// offset32 represents a constant offset value of 32 bits.
const (
	offset16 = 2
	offset32 = 4
)

// writeFn is a function type that writes data of a float32 value into a byte buffer at a specific index and returns the number of bytes written.
type writeFn func(buf []byte, data float32, idx int) int

// Format represents a data format with its associated properties for audio processing.
// Format specifies the format code, byte size, and a function to write data in the given format.
type Format struct {
	Format int
	Bytes  int
	Func   writeFn
}

// _formats is a mapping of audio format names to their corresponding Format specifications, including format type, byte size, and write function.
var _formats = map[string]Format{
	"PCM16":        {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LE},
	"PCM16Mod":     {Format: oto.FormatSignedInt16LE, Bytes: offset16, Func: writeSignedInt16LEMod},
	"FLOAT32LE":    {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LE},
	"FLOAT32LEMOD": {Format: oto.FormatFloat32LE, Bytes: offset32, Func: writeFloat32LEMod},
}

// writeSignedInt16LE writes a 16-bit signed integer in little-endian format derived from the float32 input to the buffer.
// Converts the float32 value to PCM 16-bit format, storing the result in the specified buffer location.
// Returns the number of bytes written, which is constant and defined by offset16 (2 bytes).
func writeSignedInt16LE(buf []byte, data float32, idx int) int {
	v := Float32ToPCM16(data)
	slice := buf[idx : idx+offset16]
	slice[0] = byte(v)
	slice[1] = byte(v >> 8)
	return offset16
}

// writeSignedInt16LEMod writes a float32 value as a scaled and clamped signed 16-bit little-endian integer to the buffer.
// The float32 value is scaled by a factor of 10, clamped within the range [-1.0, 1.0], and then converted to PCM16 format.
// Writes 2 bytes starting at the specified index in the buffer and returns the number of bytes written (2).
func writeSignedInt16LEMod(buf []byte, data float32, idx int) int {
	v := clamp(data * 10)
	t := Float32ToPCM16(v)
	//binary.LittleEndian.PutUint16(buf[start:end], t)
	slice := buf[idx : idx+offset16]
	slice[0] = byte(t)
	slice[1] = byte(t >> 8)
	return offset16
}

// writeFloat32LE writes a 32-bit floating-point value in little-endian format to the buffer at the specified index.
// buf is the destination byte slice where data will be written.
// data is the input float32 value to be converted into bytes.
// idx specifies the starting index in the buffer for writing.
// Returns the number of bytes written, typically 4 for a 32-bit float.
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

// writeFloat32LEMod writes a float32 value in little-endian format to the specified byte slice starting at the given index.
// The float32 value is clamped to the range [-1.0, 1.0], scaled by 10, and then converted to its IEEE 754 representation.
// It returns the constant offset32, which indicates the number of bytes written to the buffer.
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

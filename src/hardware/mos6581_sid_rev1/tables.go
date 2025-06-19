package mos6581

// SampleFreq represents the standard audio sample frequency in Hz.
// Frequency denotes the clock frequency used for audio processing.
// Cycles calculates the number of clock cycles per audio sample.
// SampleBufSize defines the size of the audio sample buffer in bytes.
// RegisterCount specifies the number of registers available.
// triTableSize indicates the size of the triangular waveform lookup table.
const (
	SampleFreq        = 44100
	Frequency         = 985248
	Cycles            = Frequency / SampleFreq
	SampleBufHalfSize = 0x138
	SampleBufSize     = SampleBufHalfSize * 2
	RegisterCount     = 32
	triTableSize      = 8192 //1 << 13 //0x1fff
	sawTableSize      = 8192
)

// _eGTable16 is a lookup table with precomputed values derived from Cycles for specific scaling factors.
var _eGTable16 = [16]uint32{
	(Cycles << 16) / 9, (Cycles << 16) / 32,
	(Cycles << 16) / 63, (Cycles << 16) / 95,
	(Cycles << 16) / 149, (Cycles << 16) / 220,
	(Cycles << 16) / 267, (Cycles << 16) / 313,
	(Cycles << 16) / 392, (Cycles << 16) / 977,
	(Cycles << 16) / 1954, (Cycles << 16) / 3126,
	(Cycles << 16) / 3906, (Cycles << 16) / 11720,
	(Cycles << 16) / 19531, (Cycles << 16) / 31251,
}

// _eGDRShiftTable256 is a precomputed lookup table containing shift values for specific input ranges, optimized for efficiency.
var _eGDRShiftTable256 [256]uint8

// _triTable8192 is an array that holds precomputed 12-bit triangle wave values for efficient lookup.
var _triTable8192 [triTableSize]uint16

// _sawTable8192 is a precomputed lookup table holding non-linear sawtooth wave values for efficient waveform synthesis.
var _sawTable8192 [sawTableSize]uint16

func computeShift(level uint8) uint8 {
	switch {
	case level < 8:
		return 5
	case level < 16:
		return 4
	case level < 28:
		return 3
	case level < 56:
		return 2
	case level < 96:
		return 1
	default:
		return 0
	}
}

// computeSawtooth generates a sawtooth waveform by combining a triangle wave and a duty cycle-based square wave.
// The function uses a lookup table accessed via triLutFn to compute the triangle wave component.
func computeSawtooth(i int, triLutFn func(uint32) uint16) uint16 {
	phaseAccumulator := uint32(i << 11)
	triOutput := triLutFn(phaseAccumulator)
	// Calculate the output of a 50% duty cycle square wave.
	// It's 0xFFFF for the first half of the accumulator cycle, 0x0000 for the second.
	rectOutput := uint16(0x0000)
	if (phaseAccumulator >> 12) < 2048 {
		// Threshold is half of the 12-bit range (4096 / 2)
		rectOutput = 0xffff
	}
	// Combine the two signals with a logical AND to get the final sawtooth.
	return triOutput & rectOutput
}

// computeTriangle calculates a 12-bit triangle waveform value based on the input integer and returns it as a uint16.
// The function shifts the input, generates an ascending or descending ramp based on the most significant bit, and scales the result.
func computeTriangle(i int) uint16 {
	accumulator := uint16(i >> 1) // 12 bit
	msb := accumulator >> 11
	var triVal uint16
	if msb&1 == 1 {
		triVal = ^accumulator //descending ramp
	} else {
		triVal = accumulator //ascending ramp
	}
	return (triVal << 4) & 0xfff0
}

// egLut calculates a value from a lookup table using the lower 4 bits of the input parameter.
func egLut(d uint8) uint32 {
	idx := d & 0xf
	return _eGTable16[idx]
}

// eGDRShiftLut retrieves a value from a pre-defined table based on the higher 8 bits of the input level.
func eGDRShiftLut(level uint32) uint8 {
	idx := (level >> 16) & 0xff
	return _eGDRShiftTable256[idx]
}

// triLut computes a value by indexing into a precomputed lookup table (_triTable8192) based on the input count.
// The input count is shifted and masked to derive the index, ensuring a bounded access within the table range.
func triLut(count uint32) uint16 {
	idx := (count >> 11) & 0x1fff
	return _triTable8192[idx]
}

// sawLut retrieves a precomputed sawtooth wave value from the lookup table based on the provided count value.
func sawLut(count uint32) uint16 {
	idx := (count >> 11) & 0x1fff
	return _sawTable8192[idx]
}

// init initializes precomputed lookup tables used for waveform synthesis and efficient operations in the system.
func init() {
	for idx := range _eGDRShiftTable256 {
		_eGDRShiftTable256[idx] = computeShift(uint8(idx))
	}
	for idx := range _triTable8192 {
		_triTable8192[idx] = computeTriangle(idx)
	}
	for idx := range _sawTable8192 {
		_sawTable8192[idx] = computeSawtooth(idx, triLut) //computeNonLinearSaw(idx)
	}
}

//triTableLen := uint16(len(_triTable8192))
//triTableHalLen := triTableLen / 2
//for i := uint16(0); i < triTableHalLen; i++ {
//val := computeTriangle(i)
//	_triTable8192[i] = val                 // Ascending ramp, indices 0..4095
//	_triTable8192[(triTableLen-1)-i] = val // Symmetric descending ramp
//}

/*
// Function that calculates the non-linear value for a given 12-bit input
func computeNonLinearSaw(i int) uint16 {
	const valueMax = 8192
	const valueHalf = valueMax / 2
	const rangeMax = valueMax - 1
	const divisor = 1048576
	const outputMask = 0xffff
	// Linear base scales a 13-bit value to a 16-bit range
	output := i << 3
	// Non-linear correction is adapted to the new range
	if i < valueHalf {
		// Divisor is increased to maintain the same "curvature"
		output -= (i * i) / divisor
	} else {
		output += ((i-rangeMax)*(i-rangeMax))/divisor - 32
	}
	if output < 0 {
		return 0
	}
	output &= outputMask
	return uint16(output)
}
*/

/*
// computeTri applies bitwise operations to the input to compute and return a modified 16-bit value.
func computeTriangle(i uint16) uint16 {
	val := (i << 4) | (i >> 8)
	return val
}
*/

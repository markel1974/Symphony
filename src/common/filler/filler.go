package filler

import (
	"math"
)

// InitRandomChanceMax defines the maximum chance value, representing 100% probability.
// InitRandomChanceHalf represents 50% probability, derived as half of InitRandomChanceMax.
// InitRandomChanceDefault defines the default chance value, representing 0.1% probability.
const (
	InitRandomChanceMax     = 10000 // 100%
	InitRandomChanceHalf    = InitRandomChanceMax / 2
	InitRandomChanceDefault = 10 // 0.1%
)

// RandomMethodNone does not flip any bits.
// RandomMethodGeom generates bit intervals between flips geometrically.
// RandomMethodUniform generates discrete uniform values for each bit.
const (
	RandomMethodNone    = 0 // flip no (or all) bits
	RandomMethodGeom    = 1 // generate bit intervals between flips
	RandomMethodUniform = 2 // generate discrete uniform per bit
)

// Filler represents a configurable structure used to generate patterns and randomness for memory initialization.
type Filler struct {
	startValue         uint // first value of the base pattern (byte value)
	valueInvert        uint // number of bytes until start value is inverted
	valueOffset        uint // offset of first pattern in bytes
	patternInvert      uint // invert base pattern after this many bytes
	patternInvertValue uint // invert base pattern with this byte
	randomStart        uint // length of random pattern in bytes
	randomRepeat       uint // repeat random pattern after this many bytes
	randomChance       uint // global random chance
}

// New creates and returns a new instance of Filler with the specified parameters for pattern and randomization configuration.
func New(startValue uint, valueInvert uint, valueOffset uint, patternInvert uint, patternInvertValue uint, randomStart uint, randomRepeat uint, randomChance uint) *Filler {
	return &Filler{
		startValue:         startValue,
		valueInvert:        valueInvert,
		valueOffset:        valueOffset,
		patternInvert:      patternInvert,
		patternInvertValue: patternInvertValue,
		randomStart:        randomStart,
		randomRepeat:       randomRepeat,
		randomChance:       randomChance,
	}
}

// InitWithPattern initializes the memory with a patterned sequence based on the properties and randomness configuration.
func (rp *Filler) InitWithPattern(memRam []uint8, ramSize uint) {
	randomMethod := RandomMethodNone
	randomMaskInitial := uint(0)
	log1mp := math.Inf(-1)
	randomNext := UintMax

	if rp.randomChance == 0 {
		// flipping no bits
		randomMethod = RandomMethodNone
		randomMaskInitial = 0x00
	} else if rp.randomChance >= InitRandomChanceMax {
		// flipping all bits; same as no bits, but with the opposite mask
		randomMethod = RandomMethodNone
		randomMaskInitial = 0xff
	} else if rp.randomChance == InitRandomChanceHalf {
		// flipping bits or not with equal probability; worst-case for the geometric spacing method, so handle separately
		randomMethod = RandomMethodUniform
	} else if rp.randomChance < InitRandomChanceHalf {
		// some other probability less than 0.5; generate the number of bits un-flipped between each flipped bit.
		randomMethod = RandomMethodGeom
		randomMaskInitial = 0x00
		log1mp = math.Log1p(float64(rp.randomChance) / InitRandomChanceMax)
		randomNext = RandomMethodGeomNext(log1mp)
	} else {
		// some other probability greater than 0.5; generate the number of bits flipped between each un-flipped bit.
		randomMethod = RandomMethodGeom
		randomMaskInitial = 0xff
		log1mp = math.Log(float64(rp.randomChance) / InitRandomChanceMax)
		randomNext = RandomMethodGeomNext(log1mp)
	}
	for offset := uint(0); offset < ramSize; offset++ {
		j := uint(0)
		k := uint(0)
		if rp.valueInvert != 0 {
			j = 0x00
			if p := ((offset + rp.valueOffset) / rp.valueInvert) & 1; p != 0 {
				j = 0xff
			}
		}
		if rp.patternInvert != 0 {
			j = 0x00
			if p := (offset / rp.patternInvert) & 1; p != 0 {
				j = rp.patternInvertValue
			}
		}
		value := uint8(rp.startValue ^ j ^ k)
		k = 0
		if rp.randomStart != 0 && rp.randomRepeat != 0 {
			k = 0
			if (offset % rp.randomRepeat) < rp.randomStart {
				k = RandomUint(0, 0xff)
			}
		}
		j = 0
		switch randomMethod {
		case RandomMethodNone:
			j = randomMaskInitial
		case RandomMethodUniform:
			j = RandomUint(0x00, 0xff)
		case RandomMethodGeom:
			j = randomMaskInitial
			for randomNext < 8 {
				j ^= 1 << randomNext
				randomNext += 1 + RandomMethodGeomNext(log1mp)
			}
			randomNext -= 8
		}
		value ^= uint8(k ^ j)
		memRam[offset] = value
	}
}

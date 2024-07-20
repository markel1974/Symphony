package filler

import (
	"math"
)

const (
	InitRandomChanceMax     = 10000 // 100%
	InitRandomChanceHalf    = InitRandomChanceMax / 2
	InitRandomChanceDefault = 10 // 0.1%
)

const (
	RandomMethodNone    = 0 // flip no (or all) bits
	RandomMethodGeom    = 1 // generate bit intervals between flips
	RandomMethodUniform = 2 // generate discrete uniform per bit
)

type Filler struct {
	startValue         int // first value of the base pattern (byte value)
	valueInvert        int // number of bytes until start value is inverted
	valueOffset        int // offset of first pattern in bytes
	patternInvert      int // invert base pattern after this many bytes
	patternInvertValue int // invert base pattern with this byte
	randomStart        int // length of random pattern in bytes
	randomRepeat       int // repeat random pattern after this many bytes
	randomChance       int // global random chance
}

/*
func NewDefaultInitiator() *Filler {
	return &Filler{
		startValue:         255,
		valueInvert:        128,
		valueOffset:        0,
		patternInvert:      0,
		patternInvertValue: 0,
		randomStart:        0,
		randomRepeat:       0,
		randomChance:       0,
	}
}
*/

func New(startValue, valueInvert, valueOffset, patternInvert, patternInvertValue, randomStart, randomRepeat, randomChance int) *Filler {
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

func (rp *Filler) InitWithPattern(memRam []uint8, ramSize uint) {
	randomMethod := RandomMethodNone
	randomMaskInitial := uint(0)
	log1mp := _log1mp //math.Inf(-1)
	randomNext := UintMax

	if rp.randomChance <= 0 {
		// flipping no bits
		randomMethod = RandomMethodNone
		randomMaskInitial = 0x00
	} else if rp.randomChance >= InitRandomChanceMax {
		// flipping all bits; same as no bits, but with the opposite mask
		randomMethod = RandomMethodNone
		randomMaskInitial = 0xff
	} else if rp.randomChance == (InitRandomChanceHalf) {
		// flipping bits or not with equal probability; worst-case for the geometric spacing method, so handle separately
		randomMethod = RandomMethodUniform
	} else if rp.randomChance < (InitRandomChanceHalf) {
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
			if p := ((int(offset) + rp.valueOffset) / rp.valueInvert) & 1; p != 0 {
				j = 0xff
			}
		}
		if rp.patternInvert != 0 {
			j = 0x00
			if p := (int(offset) / rp.patternInvert) & 1; p != 0 {
				j = uint(rp.patternInvertValue)
			}
		}
		value := uint8(uint(rp.startValue) ^ j ^ k)
		k = 0
		if rp.randomStart != 0 && rp.randomRepeat != 0 {
			k = 0
			if (int(offset) % rp.randomRepeat) < rp.randomStart {
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

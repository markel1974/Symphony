package filler

import (
	"math"
	"math/rand"
)

// LcgIncrement is the additive constant used in a linear congruential generator (LCG) for pseudo-random number generation.
const LcgIncrement = 1

// LcgMultiplier is the constant multiplier used in the Linear Congruential Generator (LCG) algorithm for random number generation.
const LcgMultiplier = 6364136223846793005

// UintMax represents the maximum value for the unsigned integer type, derived using bitwise negation.
const UintMax = ^uint(0)

// _randState is the internal state for the random number generator, initialized using LcgMultiplier and LcgIncrement.
var _randState uint64 = LcgMultiplier + LcgIncrement

// RandUint32 generates and returns a pseudo-random 32-bit unsigned integer using a linear congruential generator algorithm.
func RandUint32() uint32 {
	prevState := _randState
	_randState = (LcgMultiplier * prevState) + LcgIncrement
	outputBase := (uint32)((prevState ^ (prevState >> 18)) >> 27)
	outputRot := uint(prevState) >> 59
	return (outputBase >> outputRot) | (outputBase << (-outputRot & 31))
}

// RandUnitFloat64 generates a pseudo-random float64 value uniformly distributed in the range [0, 1).
func RandUnitFloat64() float64 {
	return 0x1.0p-32 * float64(RandUint32())
}

// RandomUint generates a random unsigned integer in the range [min, max).
func RandomUint(min int, max int) uint {
	return uint(rand.Intn(max-min) + min)
}

// RandomMethodGeomNext generates a random geometric distribution based on the logarithm provided (log1mp) and clamps the result.
func RandomMethodGeomNext(log1mp float64) uint {
	u := RandUnitFloat64()
	g := math.Floor(math.Log1p(-u) / log1mp)
	if g > float64(UintMax) {
		return UintMax
	}
	return uint(g)
}

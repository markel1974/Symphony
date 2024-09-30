package filler

import (
	"math"
	"math/rand"
)

const LcgIncrement = 1
const LcgMultiplier = 6364136223846793005
const UintMax = ^uint(0)

var _randState uint64 = LcgMultiplier + LcgIncrement

func RandUint32() uint32 {
	prevState := _randState
	_randState = (LcgMultiplier * prevState) + LcgIncrement
	outputBase := (uint32)((prevState ^ (prevState >> 18)) >> 27)
	outputRot := uint(prevState) >> 59
	return (outputBase >> outputRot) | (outputBase << (-outputRot & 31))
}

func RandUnitFloat64() float64 {
	return 0x1.0p-32 * float64(RandUint32())
}

func RandomUint(min int, max int) uint {
	return uint(rand.Intn(max-min) + min)
}

func RandomMethodGeomNext(log1mp float64) uint {
	u := RandUnitFloat64()
	g := math.Floor(math.Log1p(-u) / log1mp)
	if g > float64(UintMax) {
		return UintMax
	}
	return uint(g)
}

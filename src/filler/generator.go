package filler

// LcgIncrement is the additive constant used in the Linear Congruential Generator (LCG) formula for random number generation.
const LcgIncrement = 1

// LcgMultiplier is the multiplier constant used in the linear congruential generator (LCG) for state transitions.
const LcgMultiplier = 6364136223846793005

// Generator represents a linear congruential generator for producing pseudo-random numbers.
type Generator struct {
	// randState is the internal state for the random number generator, initialized using LcgMultiplier and LcgIncrement.
	randState uint64
}

// NewGenerator initializes and returns a new instance of Generator with a default internal random state.
func NewGenerator() *Generator {
	return &Generator{randState: LcgMultiplier + LcgIncrement}
}

// Generate produces the next pseudo-random 64-bit unsigned integer based on a linear congruential generator algorithm.
func (g *Generator) Generate() uint64 {
	prevState := g.randState
	g.randState = (LcgMultiplier * prevState) + LcgIncrement
	return prevState
}

package mos6581

import "math"

// LinearInterpolation generates a resized version of the input table using linear interpolation.
// originalTable is the input slice of uint16 representing the original data.
// newSize specifies the desired length of the resulting table.
// Returns a new slice of uint16 with interpolated values to fit the specified size.
func LinearInterpolation(originalTable []uint16, newSize int) []uint16 {
	originalSize := len(originalTable)
	if originalSize == 0 {
		return nil
	}
	newTable := make([]uint16, newSize)
	scaleFactor := float64(newSize) / float64(originalSize)
	for i := 0; i < newSize; i++ {
		// Calculate the corresponding position (as float) in the original table
		posInOriginal := float64(i) / scaleFactor

		// Find the integer indices in the original table that surround this position
		idx0 := int(math.Floor(posInOriginal))
		idx1 := idx0 + 1

		var interpolatedValue uint16

		if idx0 >= originalSize-1 {
			// If we're at (or beyond) the last point in the original table,
			// simply use the last value from the original table.
			// This handles the case where posInOriginal is very close to (originalSize - 1).
			interpolatedValue = originalTable[originalSize-1]
		} else {
			// Otherwise, perform linear interpolation
			val0 := float64(originalTable[idx0])
			val1 := float64(originalTable[idx1])

			// 'fraction' represents how far we are beyond idx0, toward idx1
			fraction := posInOriginal - float64(idx0)

			// Linear interpolation formula
			calcValue := val0*(1.0-fraction) + val1*fraction
			interpolatedValue = uint16(math.Round(calcValue))
		}
		newTable[i] = interpolatedValue
	}
	return newTable
}

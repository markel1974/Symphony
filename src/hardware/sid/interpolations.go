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
		// Calcola la posizione corrispondente (come float) nella tabella originale
		posInOriginal := float64(i) / scaleFactor

		// Trova gli indici interi della tabella originale che circondano questa posizione
		idx0 := int(math.Floor(posInOriginal))
		idx1 := idx0 + 1

		var interpolatedValue uint16

		if idx0 >= originalSize-1 {
			// Se siamo all'ultimo punto della tabella originale (o oltre),
			// usa semplicemente l'ultimo valore della tabella originale.
			// Questo gestisce il caso in cui posInOriginal è molto vicino a (originalSize - 1).
			interpolatedValue = originalTable[originalSize-1]
		} else {
			// Altrimenti, esegui l'interpolazione lineare
			val0 := float64(originalTable[idx0])
			val1 := float64(originalTable[idx1])

			// 'fraction' è quanto siamo "oltre" idx0, verso idx1
			fraction := posInOriginal - float64(idx0)

			// Formula dell'interpolazione lineare
			calcValue := val0*(1.0-fraction) + val1*fraction
			interpolatedValue = uint16(math.Round(calcValue))
		}
		newTable[i] = interpolatedValue
	}
	return newTable
}

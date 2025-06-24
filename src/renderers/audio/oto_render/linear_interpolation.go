package oto_render

import "math"

// LinearInterpolation provides a structure for performing linear interpolation on a given dataset.
// It contains a maximum size for the interpolation table and a pointer to the data table.
type LinearInterpolation struct {
	maxSize int
	table   *[]float32
}

// NewLinearInterpolation initializes a LinearInterpolation instance with a pre-allocated table of the specified maximum size.
func NewLinearInterpolation(maxSize int) *LinearInterpolation {
	v := make([]float32, maxSize)
	return &LinearInterpolation{
		table: &v,
	}
}

// LinearInterpolationF32 performs linear interpolation on a given float32 array and resizes it to the specified new size.
// It returns a pointer to the interpolated array and its new size.
// If the input array is empty, it returns nil and a size of 0.
// Allocates memory if the target new size exceeds the current length of the internal table.
func (li *LinearInterpolation) LinearInterpolationF32(originalTable *[]float32, newSize int) (*[]float32, int) {
	originalSize := len(*originalTable)
	if originalSize == 0 {
		return nil, 0
	}
	if newSize > len(*li.table) {
		v := make([]float32, newSize)
		li.table = &v
	}
	scaleFactor := float64(newSize) / float64(originalSize)
	for i := 0; i < newSize; i++ {
		posInOriginal := float64(i) / scaleFactor
		idx0 := int(math.Floor(posInOriginal))
		idx1 := idx0 + 1
		var interpolatedValue float32
		if idx0 >= originalSize-1 {
			interpolatedValue = (*originalTable)[originalSize-1]
		} else {
			val0 := (*originalTable)[idx0]
			val1 := (*originalTable)[idx1]
			fraction := float32(posInOriginal - float64(idx0))
			interpolatedValue = val0*(1.0-fraction) + val1*fraction
		}
		(*li.table)[i] = interpolatedValue
	}
	return li.table, newSize
}

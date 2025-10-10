package pixels

import (
	"fmt"
)

// TriangleData represents a 2D triangle with positional, color, texture, clipping, and intensity-related properties.
type TriangleData struct {
	Position  Vector
	Color     RGBA
	Picture   Vector
	Intensity float64
	ClipRect  Rect
	IsClipped bool
}

// _zeroValueTriangleData is a predefined TriangleData instance with default values, used for initialization and reuse.
var _zeroValueTriangleData = TriangleData{Color: RGBA{1, 1, 1, 1}}

// TrianglesData represents a slice of TriangleData, used to store data for multiple triangles.
type TrianglesData []TriangleData

// NewTrianglesData creates and returns a new instance of TrianglesData populated with a specified number of default elements.
func NewTrianglesData(len int) *TrianglesData {
	td := make(TrianglesData, len)
	for i := 0; i < len; i++ {
		td[i] = _zeroValueTriangleData
	}
	return &td
}

// Len returns the number of elements in the TrianglesData. It is equivalent to the length of the underlying slice.
func (td *TrianglesData) Len() int {
	return len(*td)
}

// SetLen adjusts the length of the TrianglesData slice to the specified value, appending or truncating elements as needed.
func (td *TrianglesData) SetLen(len int) {
	if len > td.Len() {
		needAppend := len - td.Len()
		for i := 0; i < needAppend; i++ {
			*td = append(*td, _zeroValueTriangleData)
		}
	}
	if len < td.Len() {
		*td = (*td)[:len]
	}
}

// Slice extracts a subrange of the TrianglesData from index i to j and returns it as a new ITriangles implementation.
func (td *TrianglesData) Slice(i, j int) ITriangles {
	s := (*td)[i:j]
	return &s
}

// updateData updates the data of the TrianglesData object by copying values from the provided ITriangles input.
// It optimizes for cases where the input is also a TrianglesData, otherwise performs manual property-wise copying.
func (td *TrianglesData) updateData(t ITriangles) {
	// fast path optimization
	if t, ok := t.(*TrianglesData); ok {
		copy(*td, *t)
		return
	}

	// slow path manual copy
	if t, ok := t.(ITrianglesPosition); ok {
		for i := range *td {
			(*td)[i].Position = t.Position(i)
		}
	}
	if t, ok := t.(ITrianglesColor); ok {
		for i := range *td {
			(*td)[i].Color = t.Color(i)
		}
	}
	if t, ok := t.(ITrianglesPicture); ok {
		for i := range *td {
			(*td)[i].Picture, (*td)[i].Intensity = t.Picture(i)
		}
	}
	if t, ok := t.(ITrianglesClipped); ok {
		for i := range *td {
			(*td)[i].ClipRect, (*td)[i].IsClipped = t.ClipRect(i)
		}
	}
}

// Update copies data from the given ITriangles object into the current TrianglesData instance. Panics if lengths differ.
func (td *TrianglesData) Update(t ITriangles) {
	if td.Len() != t.Len() {
		panic(fmt.Errorf("(%T).Update: invalid triangles length", td))
	}
	td.updateData(t)
}

// Copy creates and returns a new ITriangles instance, duplicating the data of the current TrianglesData object.
func (td *TrianglesData) Copy() ITriangles {
	copyTd := NewTrianglesData(td.Len())
	copyTd.Update(td)
	return copyTd
}

// Position retrieves the position vector of the triangle at the specified index in the TrianglesData collection.
func (td *TrianglesData) Position(i int) Vector {
	return (*td)[i].Position
}

// Color returns the RGBA color of the triangle at the specified index i in the TrianglesData slice.
func (td *TrianglesData) Color(i int) RGBA {
	return (*td)[i].Color
}

// Picture retrieves the picture vector and intensity value for the triangle at the specified index.
func (td *TrianglesData) Picture(i int) (pic Vector, intensity float64) {
	return (*td)[i].Picture, (*td)[i].Intensity
}

// ClipRect returns the clipping rectangle and a boolean indicating whether the triangle at index i is clipped.
func (td *TrianglesData) ClipRect(i int) (rect Rect, has bool) {
	return (*td)[i].ClipRect, (*td)[i].IsClipped
}

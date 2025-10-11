package executor

import (
	"fmt"
)

// VertexSlice represents a slice of vertices managed by a VertexArray, with a defined starting and ending index.
type VertexSlice struct {
	va *VertexArray
	i  int
	j  int
}

// NewVertexSlice creates a new VertexSlice with a specified length and capacity tied to the provided Shader.
// Panics if the given length exceeds the capacity. Returns a pointer to the newly created VertexSlice.
func NewVertexSlice(shader *Shader, len int, cap int) *VertexSlice {
	if len > cap {
		panic("failed to make vertex slice: len > cap")
	}
	return &VertexSlice{
		va: NewVertexArray(shader, cap),
		i:  0,
		j:  len,
	}
}

// VertexFormat returns the attribute format of the vertices used in the VertexSlice.
func (vs *VertexSlice) VertexFormat() AttrFormat {
	return vs.va.format
}

// Stride returns the number of float32 elements per vertex, calculated as the vertex array stride divided by 4.
func (vs *VertexSlice) Stride() int {
	return vs.va.stride / 4
}

// Len returns the length of the VertexSlice, calculated as the difference between its end index (j) and start index (i).
func (vs *VertexSlice) Len() int {
	return vs.j - vs.i
}

// Cap returns the remaining capacity of the VertexSlice, calculated as the difference between va.cap and the starting index.
func (vs *VertexSlice) Cap() int {
	return vs.va.cap - vs.i
}

// SetLen adjusts the length of the VertexSlice, resizing it and updating the underlying VertexArray as necessary.
func (vs *VertexSlice) SetLen(len int) {
	vs.End()
	*vs = vs.grow(len)
	vs.Begin()
}

// grow adjusts the size of the VertexSlice to accommodate at least the specified length.
// If the current capacity suffices, it returns a resized slice. Otherwise, it creates a new slice with increased capacity.
func (vs *VertexSlice) grow(len int) VertexSlice {
	if len <= vs.Cap() {
		return VertexSlice{
			va: vs.va,
			i:  vs.i,
			j:  vs.i + len,
		}
	}

	newCap := vs.Cap()
	if newCap < 1024 {
		newCap += newCap
	} else {
		newCap += newCap / 4
	}
	if newCap < len {
		newCap = len
	}
	newVs := VertexSlice{
		va: NewVertexArray(vs.va.shader, newCap),
		i:  0,
		j:  len,
	}

	newVs.Begin()
	newVs.Slice(0, vs.Len()).SetVertexData(vs.VertexData())
	newVs.End()
	return newVs
}

// Slice creates a new VertexSlice from the current slice in the range [i, j). Panics if indices are out of range.
func (vs *VertexSlice) Slice(i, j int) *VertexSlice {
	if i < 0 || j < i || j > vs.va.cap {
		panic("failed to slice vertex slice: index out of range")
	}
	return &VertexSlice{
		va: vs.va,
		i:  vs.i + i,
		j:  vs.i + j,
	}
}

// SetVertexData sets vertex data for the slice, ensuring the data length matches the slice's vertex count and stride.
// Panics if the data length is incompatible with the slice.
func (vs *VertexSlice) SetVertexData(data []float32) {
	if len(data)/vs.Stride() != vs.Len() {
		fmt.Println(len(data)/vs.Stride(), vs.Len())
		panic("set vertex data: wrong length of vertices")
	}
	vs.va.SetVertexData(vs.i, vs.j, data)
}

// VertexData retrieves vertex data as a slice of float32 values between the slice's start and end indices.
func (vs *VertexSlice) VertexData() []float32 {
	return vs.va.VertexData(vs.i, vs.j)
}

// Draw renders the current range of vertices from the VertexSlice using the associated VertexArray's Draw method.
func (vs *VertexSlice) Draw() {
	vs.va.Draw(vs.i, vs.j)
}

// Begin prepares the underlying VertexArray for subsequent OpenGL operations by binding its VAO and VBO.
func (vs *VertexSlice) Begin() {
	vs.va.Begin()
}

// End finalizes operations on the VertexSlice by calling End on the underlying VertexArray.
func (vs *VertexSlice) End() {
	vs.va.End()
}

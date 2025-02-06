package executor

import (
	"fmt"
)

type VertexSlice struct {
	va *VertexArray
	i  int
	j  int
}

func NewVertexSlice(shader *Shader, len, cap int) *VertexSlice {
	if len > cap {
		panic("failed to make vertex slice: len > cap")
	}
	return &VertexSlice{
		va: NewVertexArray(shader, cap),
		i:  0,
		j:  len,
	}
}

func (vs *VertexSlice) VertexFormat() AttrFormat {
	return vs.va.format
}

func (vs *VertexSlice) Stride() int {
	return vs.va.stride / 4
}

func (vs *VertexSlice) Len() int {
	return vs.j - vs.i
}

func (vs *VertexSlice) Cap() int {
	return vs.va.cap - vs.i
}

func (vs *VertexSlice) SetLen(len int) {
	vs.End()
	*vs = vs.grow(len)
	vs.Begin()
}

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

func (vs *VertexSlice) SetVertexData(data []float32) {
	if len(data)/vs.Stride() != vs.Len() {
		fmt.Println(len(data)/vs.Stride(), vs.Len())
		panic("set vertex data: wrong length of vertices")
	}
	vs.va.SetVertexData(vs.i, vs.j, data)
}

func (vs *VertexSlice) VertexData() []float32 {
	return vs.va.VertexData(vs.i, vs.j)
}

func (vs *VertexSlice) Draw() {
	vs.va.Draw(vs.i, vs.j)
}

func (vs *VertexSlice) Begin() {
	vs.va.Begin()
}

func (vs *VertexSlice) End() {
	vs.va.End()
}

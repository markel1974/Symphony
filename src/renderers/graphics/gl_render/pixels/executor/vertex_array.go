package executor

import (
	"errors"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// VertexArray represents an OpenGL Vertex Array Object (VAO) that manages vertex data and its binding state.
// It encapsulates a VAO, Vertex Buffer Object (VBO), data capacity, attribute format, stride, offsets, and shader reference.
// Used for efficient rendering of vertex data with OpenGL.
type VertexArray struct {
	vao, vbo Binder
	cap      int
	format   AttrFormat
	stride   int
	offset   []int
	shader   *Shader
}

// vertexArrayMinCap defines the minimum capacity for a vertex array, ensuring sufficient memory allocation for attributes.
const vertexArrayMinCap = 4

// NewVertexArray creates and initializes a new VertexArray with the provided shader and capacity.
// Ensures a minimum capacity is respected and sets up attribute pointers based on the shader's vertex format.
// Returns a pointer to the created VertexArray.
func NewVertexArray(shader *Shader, cap int) *VertexArray {
	if cap < vertexArrayMinCap {
		cap = vertexArrayMinCap
	}

	va := &VertexArray{
		vao: Binder{
			restoreLoc: gl.VERTEX_ARRAY_BINDING,
			bindFunc: func(obj uint32) {
				gl.BindVertexArray(obj)
			},
		},
		vbo: Binder{
			restoreLoc: gl.ARRAY_BUFFER_BINDING,
			bindFunc: func(obj uint32) {
				gl.BindBuffer(gl.ARRAY_BUFFER, obj)
			},
		},
		cap:    cap,
		format: shader.VertexFormat(),
		stride: shader.VertexFormat().Size(),
		offset: make([]int, len(shader.VertexFormat())),
		shader: shader,
	}

	offset := 0
	for i, attr := range va.format {
		switch attr.Type {
		case Float, Vec2, Vec3, Vec4:
		default:
			panic(errors.New("failed to create vertex array: invalid attribute type"))
		}
		va.offset[i] = offset
		offset += attr.Type.Size()
	}

	gl.GenVertexArrays(1, &va.vao.obj)

	va.vao.Bind()

	gl.GenBuffers(1, &va.vbo.obj)
	defer va.vbo.Bind().Restore()

	emptyData := make([]byte, cap*va.stride)
	gl.BufferData(gl.ARRAY_BUFFER, len(emptyData), gl.Ptr(emptyData), gl.DYNAMIC_DRAW)

	for i, attr := range va.format {
		loc := gl.GetAttribLocation(shader.program.obj, gl.Str(attr.Name+"\x00"))

		var size int32
		switch attr.Type {
		case Float:
			size = 1
		case Vec2:
			size = 2
		case Vec3:
			size = 3
		case Vec4:
			size = 4
		default:
			size = 0
		}

		gl.VertexAttribPointerWithOffset(
			uint32(loc),
			size,
			gl.FLOAT,
			false,
			int32(va.stride),
			uintptr(va.offset[i]),
		)
		gl.EnableVertexAttribArray(uint32(loc))
	}

	va.vao.Restore()

	runtime.SetFinalizer(va, (*VertexArray).Delete)

	return va
}

// Delete releases the GPU resources associated with the VertexArray, including vertex array and buffer objects.
func (va *VertexArray) Delete() {
	GraphicThread.Post(func() {
		gl.DeleteVertexArrays(1, &va.vao.obj)
		gl.DeleteBuffers(1, &va.vbo.obj)
	})
}

// Begin binds the vertex array object (VAO) and vertex buffer object (VBO) for subsequent OpenGL operations.
func (va *VertexArray) Begin() {
	va.vao.Bind()
	va.vbo.Bind()
}

// End restores the previous state of the vertex array and vertex buffer by calling their respective Restore methods.
func (va *VertexArray) End() {
	va.vbo.Restore()
	va.vao.Restore()
}

// Draw renders a portion of the vertex array as triangles, using indices i and j to define the range of vertices.
func (va *VertexArray) Draw(i int, j int) {
	gl.DrawArrays(gl.TRIANGLES, int32(i), int32(j-i))
}

// SetVertexData updates a subset of vertex buffer data between indices i and j with the provided float32 data slice.
func (va *VertexArray) SetVertexData(i int, j int, data []float32) {
	if j-i == 0 {
		// avoid setting 0 bytes of buffer data
		return
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, i*va.stride, len(data)*4, gl.Ptr(data))
}

// VertexData retrieves a slice of vertex data from the buffer between indices i and j. Returns nil if the range is invalid.
func (va *VertexArray) VertexData(i int, j int) []float32 {
	if j-i == 0 {
		// avoid getting 0 bytes of buffer data
		return nil
	}
	data := make([]float32, (j-i)*va.stride/4)
	gl.GetBufferSubData(gl.ARRAY_BUFFER, i*va.stride, len(data)*4, gl.Ptr(data))
	return data
}

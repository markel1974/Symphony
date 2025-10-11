package executor

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// Init initializes the OpenGL context, enabling blending and scissor testing, and sets the blend equation.
// Panics if OpenGL initialization fails.
func Init() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.BlendEquation(gl.FUNC_ADD)
}

// Clear sets the clear color and clears the color buffer with the specified RGBA values.
func Clear(r float32, g float32, b float32, a float32) {
	gl.ClearColor(r, g, b, a)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

// Bounds sets the OpenGL viewport and scissor rectangle to the specified position and size.
func Bounds(x int, y int, w int, h int) {
	gl.Viewport(int32(x), int32(y), int32(w), int32(h))
	gl.Scissor(int32(x), int32(y), int32(w), int32(h))
}

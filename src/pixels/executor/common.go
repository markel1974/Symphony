package executor

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func Init() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.BlendEquation(gl.FUNC_ADD)
}

func Clear(r float32, g float32, b float32, a float32) {
	gl.ClearColor(r, g, b, a)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func Bounds(x int, y int, w int, h int) {
	gl.Viewport(int32(x), int32(y), int32(w), int32(h))
	gl.Scissor(int32(x), int32(y), int32(w), int32(h))
}

package executor

import "github.com/go-gl/gl/v3.3-core/gl"

type BlendFactor int

const (
	One              = BlendFactor(gl.ONE)
	Zero             = BlendFactor(gl.ZERO)
	SrcAlpha         = BlendFactor(gl.SRC_ALPHA)
	DstAlpha         = BlendFactor(gl.DST_ALPHA)
	OneMinusSrcAlpha = BlendFactor(gl.ONE_MINUS_SRC_ALPHA)
	OneMinusDstAlpha = BlendFactor(gl.ONE_MINUS_DST_ALPHA)
)

func BlendFunc(src BlendFactor, dst BlendFactor) {
	gl.BlendFunc(uint32(src), uint32(dst))
}

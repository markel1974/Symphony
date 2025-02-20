package executor

import "github.com/go-gl/gl/v3.3-core/gl"

// BlendFactor represents constants used to specify blending factors in rendering operations.
type BlendFactor int

// One represents the GL constant for a blend factor of one.
// Zero represents the GL constant for a blend factor of zero.
// SrcAlpha represents the GL constant for a source alpha blend factor.
// DstAlpha represents the GL constant for a destination alpha blend factor.
// OneMinusSrcAlpha represents the GL constant for one minus source alpha blending.
// OneMinusDstAlpha represents the GL constant for one minus destination alpha blending.
const (
	One              = BlendFactor(gl.ONE)
	Zero             = BlendFactor(gl.ZERO)
	SrcAlpha         = BlendFactor(gl.SRC_ALPHA)
	DstAlpha         = BlendFactor(gl.DST_ALPHA)
	OneMinusSrcAlpha = BlendFactor(gl.ONE_MINUS_SRC_ALPHA)
	OneMinusDstAlpha = BlendFactor(gl.ONE_MINUS_DST_ALPHA)
)

// BlendFunc sets the RGB and alpha blend factors for OpenGL rendering using the specified source and destination factors.
func BlendFunc(src BlendFactor, dst BlendFactor) {
	gl.BlendFunc(uint32(src), uint32(dst))
}

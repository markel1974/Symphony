package executor

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Texture represents an OpenGL texture object with dimensions, smoothing settings, and utility methods for manipulation.
type Texture struct {
	tex    Binder
	width  int
	height int
	smooth bool
}

// NewTexture creates a new Texture object with the specified width, height, smoothing mode, and pixel data.
// It initializes the OpenGL texture, sets its parameters, and associates the given pixel data with it.
func NewTexture(width int, height int, smooth bool, pixels []uint8) *Texture {
	tex := &Texture{
		tex: Binder{
			restoreLoc: gl.TEXTURE_BINDING_2D,
			bindFunc: func(obj uint32) {
				gl.BindTexture(gl.TEXTURE_2D, obj)
			},
		},
		width:  width,
		height: height,
	}

	gl.GenTextures(1, &tex.tex.obj)

	tex.Begin()

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	//borderColor := mgl32.Vec4{0, 0, 0, 0}
	//gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)

	tex.SetSmooth(smooth)

	runtime.SetFinalizer(tex, (*Texture).delete)

	tex.End()

	return tex
}

// delete releases the GPU texture associated with the Texture instance.
// The operation is executed on the graphical thread using GraphicThread.
func (t *Texture) delete() {
	GraphicThread.Post(func() {
		gl.DeleteTextures(1, &t.tex.obj)
	})
}

// ID returns the OpenGL texture ID associated with the Texture instance.
func (t *Texture) ID() uint32 {
	return t.tex.obj
}

// Width returns the width of the texture in pixels.
func (t *Texture) Width() int {
	return t.width
}

// Height returns the height of the texture in pixels.
func (t *Texture) Height() int {
	return t.height
}

// SetPixels updates a rectangular region of the texture with the provided pixel data, which must match the region size.
func (t *Texture) SetPixels(x int, y int, w int, h int, pixels []uint8) {
	//if len(pixels) != w*h*4 {
	//	return
	//	//panic("set pixels: wrong number of pixels")
	//}
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, int32(x), int32(y), int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
}

// Pixels retrieves the RGBA pixel data from the specified subregion (x, y, w, h) of the texture.
// The pixel data is returned as a slice of uint8 values in alpha-premultiplied RGBA format.
func (t *Texture) Pixels(x int, y int, w int, h int) []uint8 {
	pixels := make([]uint8, t.width*t.height*4)
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	subPixels := make([]uint8, w*h*4)
	for i := 0; i < h; i++ {
		row := pixels[(i+y)*t.width*4+x*4 : (i+y)*t.width*4+(x+w)*4]
		subRow := subPixels[i*w*4 : (i+1)*w*4]
		copy(subRow, row)
	}
	return subPixels
}

// SetSmooth configures the texture sampling method for scaling, toggling between smooth (linear) and pixelated (nearest).
func (t *Texture) SetSmooth(smooth bool) {
	t.smooth = smooth
	if smooth {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}
}

// Smooth returns the current smoothing state of the texture.
func (t *Texture) Smooth() bool {
	return t.smooth
}

// Begin binds the texture object to the current OpenGL context, making it active for subsequent operations.
func (t *Texture) Begin() {
	t.tex.Bind()
}

// End restores the previous OpenGL texture binding state using the underlying binder.
func (t *Texture) End() {
	t.tex.Restore()
}

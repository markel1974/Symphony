package pixels

import (
	"github.com/markel1974/symphony/src/renderers/graphics/gl_render/pixels/executor"
)

// GLFrame represents a graphical frame used for rendering operations and managing pixel data efficiently.
// It contains a frame object, boundary definitions, pixel data, and a flag indicating if the frame is dirty.
type GLFrame struct {
	frame  *executor.Frame
	bounds Rect
	pixels []uint8
	dirty  bool
}

// NewGLFrame creates a new GLFrame instance and sets its bounds to the specified Rect.
// Returns a pointer to the newly created GLFrame.
func NewGLFrame(bounds Rect) *GLFrame {
	gf := new(GLFrame)
	gf.SetBounds(bounds)
	return gf
}

// SetBounds updates the frame's bounds to the specified rectangle and reallocates resources if necessary.
func (gf *GLFrame) SetBounds(bounds Rect) {
	if bounds == gf.Bounds() {
		return
	}

	executor.GraphicThread.Call(func() {
		oldF := gf.frame

		_, _, w, h := intBounds(bounds)
		if w <= 0 {
			w = 1
		}
		if h <= 0 {
			h = 1
		}
		gf.frame = executor.NewFrame(w, h, true)

		if oldF != nil {
			ox, oy, ow, oh := intBounds(bounds)
			oldF.Blit(
				gf.frame,
				ox, oy, ox+ow, oy+oh,
				ox, oy, ox+ow, oy+oh,
			)
		}
	})

	gf.bounds = bounds
	gf.pixels = nil
	gf.dirty = true
}

// Bounds returns the rectangle defining the boundaries of the GLFrame instance.
func (gf *GLFrame) Bounds() Rect {
	return gf.bounds
}

// Color returns the RGBA color at the specified Vector position in the GLFrame.
// If the position is out of bounds, it returns an RGBA with zero alpha.
// Updates pixel data if the GLFrame is marked as dirty.
func (gf *GLFrame) Color(at Vector) RGBA {
	if gf.dirty {
		executor.GraphicThread.Call(func() {
			tex := gf.frame.Texture()
			tex.Begin()
			gf.pixels = tex.Pixels(0, 0, tex.Width(), tex.Height())
			tex.End()
		})
		gf.dirty = false
	}
	if !gf.bounds.Contains(at) {
		return Alpha(0)
	}
	bx, by, bw, _ := intBounds(gf.bounds)
	x, y := int(at.X)-bx, int(at.Y)-by
	off := y*bw + x
	return RGBA{
		R: float64(gf.pixels[off*4+0]) / 255,
		G: float64(gf.pixels[off*4+1]) / 255,
		B: float64(gf.pixels[off*4+2]) / 255,
		A: float64(gf.pixels[off*4+3]) / 255,
	}
}

// Frame returns the underlying executor.Frame instance associated with this GLFrame.
func (gf *GLFrame) Frame() *executor.Frame {
	return gf.frame
}

// Texture retrieves the underlying executor.Texture associated with the GLFrame instance.
func (gf *GLFrame) Texture() *executor.Texture {
	return gf.frame.Texture()
}

// Dirty marks the GLFrame as needing to be redrawn by setting its dirty flag to true.
func (gf *GLFrame) Dirty() {
	gf.dirty = true
}

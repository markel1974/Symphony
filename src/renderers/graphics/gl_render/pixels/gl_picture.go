package pixels

import (
	"math"

	"github.com/markel1974/symphony/src/renderers/graphics/gl_render/pixels/executor"
)

// GLPicture represents a graphical picture with texture, pixel data, and its bounding rectangle.
type GLPicture struct {
	bounds  Rect
	texture *executor.Texture
	pixels  []uint8
}

// NewGLPicture creates a new GLPicture based on the provided IPicture. It initializes texture and pixel data.
func NewGLPicture(p IPicture) *GLPicture {
	var pixels []uint8
	var bw, bh int

	if pr, ok := p.(IPicturePixels); ok {
		_, _, bw, bh, pixels = pr.Pixels()
	} else if pc, ok := p.(IPictureColor); ok {
		_, _, bw, bh, pixels = extractPixelsFromPictureColor(pc)
	} else {
		_, _, bw, bh, pixels = extractPixelsFromPicture(p)
	}

	var texture *executor.Texture

	executor.GraphicThread.Call(func() {
		texture = executor.NewTexture(bw, bh, true, pixels)
	})

	gp := &GLPicture{
		bounds:  p.Bounds(),
		texture: texture,
		pixels:  pixels,
	}

	return gp
}

// Bounds returns the rectangular bounds of the GLPicture.
func (gp *GLPicture) Bounds() Rect {
	return gp.bounds
}

// Texture returns the associated texture of the GLPicture instance.
func (gp *GLPicture) Texture() *executor.Texture {
	return gp.texture
}

// Color returns the RGBA color at the specified Vector position within the GLPicture bounds.
// If the position is outside bounds, it returns fully transparent (Alpha 0).
func (gp *GLPicture) Color(at Vector) RGBA {
	if !gp.bounds.Contains(at) {
		return Alpha(0)
	}
	bx, by, bw, _ := intBounds(gp.bounds)
	x, y := int(at.X)-bx, int(at.Y)-by
	off := y*bw + x
	return RGBA{
		R: float64(gp.pixels[off*4+0]) / 255,
		G: float64(gp.pixels[off*4+1]) / 255,
		B: float64(gp.pixels[off*4+2]) / 255,
		A: float64(gp.pixels[off*4+3]) / 255,
	}
}

// Update updates the GLPicture texture with the pixel data from the specified IPicture implementation.
func (gp *GLPicture) Update(p IPicture) {
	var pixels []uint8
	var bx, by, bw, bh int

	if pr, ok := p.(IPicturePixels); ok {
		bx, by, bw, bh, pixels = pr.Pixels()
	} else if pc, ok := p.(IPictureColor); ok {
		bx, by, bw, bh, pixels = extractPixelsFromPictureColor(pc)
	} else {
		bx, by, bw, bh, pixels = extractPixelsFromPicture(p)
	}

	executor.GraphicThread.Call(func() {
		gp.texture.Begin()
		gp.texture.SetPixels(bx, by, bw, bh, pixels)
		gp.texture.End()
	})
}

// extractPixelsFromPictureColor extracts pixel data from an IPictureColor into a slice and returns dimensions & position.
func extractPixelsFromPictureColor(pc IPictureColor) (int, int, int, int, []uint8) {
	bounds := pc.Bounds()
	bx, by, bw, bh := intBounds(bounds)
	pixels := make([]uint8, 4*bw*bh)
	for y := 0; y < bh; y++ {
		for x := 0; x < bw; x++ {
			at := NewVec(
				math.Max(float64(bx+x), bounds.Min.X),
				math.Max(float64(by+y), bounds.Min.Y),
			)
			color := pc.Color(at)
			off := ((y * bw) + x) * 4
			pixels[off+0] = uint8(color.R * 255)
			pixels[off+1] = uint8(color.G * 255)
			pixels[off+2] = uint8(color.B * 255)
			pixels[off+3] = uint8(color.A * 255)
		}
	}
	return bx, by, bw, bh, pixels
}

// extractPixelsFromPicture extracts pixel data and bounding box dimensions from an IPicture instance.
func extractPixelsFromPicture(pp IPicture) (int, int, int, int, []uint8) {
	bounds := pp.Bounds()
	bx, by, bw, bh := intBounds(bounds)
	pixels := make([]uint8, 4*bw*bh)
	return bx, by, bw, bh, pixels
}

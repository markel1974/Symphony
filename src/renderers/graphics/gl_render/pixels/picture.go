package pixels

import (
	"image"
	"image/draw"
	"math"
)

// Picture represents an image with defined width, height, and a rectangular boundary.
// It holds information about the image's pixel data, stride, and additional dimensions for processing.
type Picture struct {
	width  int
	height int
	rect   Rect
	stride int
	pixels []uint8
	length int
	lastY  int
	bx     int
	by     int
	bw     int
	bh     int
}

// NewPictureFromPicture creates a new Picture from an existing IPicture, copying its bounds and optionally its colors.
func NewPictureFromPicture(picture IPicture) *Picture {
	if pd, ok := picture.(*Picture); ok {
		return pd
	}
	bounds := picture.Bounds()
	pd := NewPicture(bounds)
	if pc, ok := picture.(IPictureColor); ok {
		for y := math.Floor(bounds.Min.Y); y < bounds.Max.Y; y++ {
			for x := math.Floor(bounds.Min.X); x < bounds.Max.X; x++ {
				at := NewVec(math.Max(x, bounds.Min.X), math.Max(y, bounds.Min.Y))
				col := pc.Color(at)
				pd.SetRGBA(int(x), int(y), uint8(col.R*255), uint8(col.G*255), uint8(col.B*255), uint8(col.A*255))
			}
		}
	}
	return pd
}

// NewPictureFromImage creates a new Picture instance from the provided image.Image by copying and vertically flipping it.
func NewPictureFromImage(img image.Image) *Picture {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	verticalFlip(rgba)
	pd := NewPicture(NewRect(
		float64(rgba.Bounds().Min.X),
		float64(rgba.Bounds().Min.Y),
		float64(rgba.Bounds().Max.X),
		float64(rgba.Bounds().Max.Y),
	))
	copy(pd.pixels, rgba.Pix)
	return pd
}

// NewPicture initializes a new Picture object with the given Rect and computes its width, height, stride, and pixel data.
func NewPicture(rect Rect) *Picture {
	w := int(math.Ceil(rect.Max.X)) - int(math.Floor(rect.Min.X))
	h := int(math.Ceil(rect.Max.Y)) - int(math.Floor(rect.Min.Y))
	l := 4 * w * h
	bx, by, bw, bh := intBounds(rect)
	s := &Picture{
		width:  w,
		height: h,
		stride: 4 * w,
		rect:   rect,
		pixels: make([]uint8, l),
		length: l - 4,
		lastY:  int(rect.Max.Y) - 1,
		bx:     bx,
		by:     by,
		bw:     bw,
		bh:     bh,
	}
	return s
}

// ComputeIndex calculates the flat array index for given x and y coordinates in the Picture's pixel data.
// The calculation considers the stride and rectangular bounds of the Picture.
func (s *Picture) ComputeIndex(x int, y int) int {
	y = s.lastY - y
	i := (y-int(s.rect.Min.Y))*s.stride + (x-int(s.rect.Min.X))*4
	return i
}

// SetRGBAArray sets the RGBA values at the specified x, y coordinates using an array of 4 uint8 values (red, green, blue, alpha).
func (s *Picture) SetRGBAArray(x int, y int, rgba []uint8) {
	//flip
	//y = (int(s.rect.Max.Y) -1) - y
	y = s.lastY - y
	i := (y-int(s.rect.Min.Y))*s.stride + (x-int(s.rect.Min.X))*4
	if i >= 0 && i < s.length {
		copy(s.pixels[i:], rgba)
	}
}

// SetRGBADirectArray sets the RGBA color values starting at the given pixel index directly in the pixel data array.
//func (s *Picture) SetRGBADirectArray(i int, rgba []uint8) {
//	copy(s.pixels[i:], rgba)
//}

// SetRGBA4DirectArrayPtr sets a 4-byte RGBA color at the specified index in the pixel buffer using a pointer to a 4-byte array.
func (s *Picture) SetRGBA4DirectArrayPtr(i int, rgba *[4]uint8) {
	copy(s.pixels[i:], rgba[:])
}

// SetRGBA8DirectArrayPtr updates the pixel data at the specified index with values from a pointer to an 8-element RGBA array.
func (s *Picture) SetRGBA8DirectArrayPtr(i int, rgba *[8]uint8) {
	copy(s.pixels[i:], rgba[:])
}

// SetRGBA32DirectArrayPtr copies a 32-length RGBA uint8 array to the pixel data at the specified index i in the Picture.
func (s *Picture) SetRGBA32DirectArrayPtr(i int, rgba *[32]uint8) {
	copy(s.pixels[i:], rgba[:])
}

// SetRGBADirectArrayPtr sets a portion of the pixel buffer starting at the specified index with the given RGBA data slice.
func (s *Picture) SetRGBADirectArrayPtr(i int, rgba *[]uint8, width int) {
	copy(s.pixels[i:], (*rgba)[:width])
}

// SetRGBA sets the RGBA color values at the specified (x, y) pixel position in the Picture.
// The color channels are represented by r, g, b, and a, where a denotes the alpha (transparency).
// The operation updates the corresponding pixel in the Picture's pixel buffer, if the index is within bounds.
func (s *Picture) SetRGBA(x int, y int, r uint8, g uint8, b uint8, a uint8) {
	//flip
	//y = (int(s.rect.Max.Y) -1) - y
	y = s.lastY - y
	i := (y-int(s.rect.Min.Y))*s.stride + (x-int(s.rect.Min.X))*4
	if i >= 0 && i < s.length {
		s.pixels[i] = r
		s.pixels[i+1] = g
		s.pixels[i+2] = b
		s.pixels[i+3] = a
	}
}

// SetRGBASize sets a rectangle of pixels centered at (x, y) with the specified width and RGBA color values.
func (s *Picture) SetRGBASize(x int, y int, r uint8, g uint8, b uint8, a uint8, size int) {
	k := size / 2
	for offsetX := -k; offsetX < k; offsetX++ {
		for offsetY := -k; offsetY < k; offsetY++ {
			s.SetRGBA(x+offsetX, y+offsetY, r, g, b, a)
		}
	}
}

// Width returns the width of the Picture as an integer.
func (s *Picture) Width() int {
	return s.width
}

// Height returns the height dimension of the Picture instance.
func (s *Picture) Height() int {
	return s.height
}

// Length returns the length of the `Picture`, which represents the total number of elements in the underlying pixel data.
func (s *Picture) Length() int {
	return s.length
}

// Stride returns the number of bytes between adjacent rows of the image in memory.
func (s *Picture) Stride() int {
	return s.stride
}

// Bounds returns the rectangular boundaries of the Picture as a Rect.
func (s *Picture) Bounds() Rect {
	return s.rect
}

// Pixels returns the bounding coordinates (bx, by, bw, bh) and the pixel data as a slice of uint8 from the Picture struct.
func (s *Picture) Pixels() (int, int, int, int, []uint8) {
	return s.bx, s.by, s.bw, s.bh, s.pixels
}

// Image converts the Picture's pixel data into an image.RGBA instance and returns it.
func (s *Picture) Image() *image.RGBA {
	bounds := image.Rect(
		int(math.Floor(s.rect.Min.X)),
		int(math.Floor(s.rect.Min.Y)),
		int(math.Ceil(s.rect.Max.X)),
		int(math.Ceil(s.rect.Max.Y)),
	)
	rgba := image.NewRGBA(bounds)
	copy(rgba.Pix, s.pixels)
	return rgba
}

// verticalFlip flips the pixel data of an image.RGBA vertically, modifying the input image in place.
func verticalFlip(rgba *image.RGBA) {
	bounds := rgba.Bounds()
	width := bounds.Dx()
	tmpRow := make([]uint8, width*4)
	for i, j := 0, bounds.Dy()-1; i < j; i, j = i+1, j-1 {
		iRow := rgba.Pix[i*rgba.Stride : i*rgba.Stride+width*4]
		jRow := rgba.Pix[j*rgba.Stride : j*rgba.Stride+width*4]
		copy(tmpRow, iRow)
		copy(iRow, jRow)
		copy(jRow, tmpRow)
	}
}

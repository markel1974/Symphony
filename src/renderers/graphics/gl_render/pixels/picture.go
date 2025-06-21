package pixels

import (
	"image"
	"image/draw"
	"math"
)

// verticalFlip vertically inverts the pixels of the provided RGBA image by swapping rows in place.
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

// Picture represents an image, defining its dimensions, boundary rectangle, pixel data, and other related properties.
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

// NewPictureFromPicture creates a new Picture from an existing IPicture, copying pixel data if applicable.
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

// NewPictureFromImage creates a new Picture from the given image.Image, with a vertically flipped RGBA pixel buffer.
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

// NewPicture creates a new Picture instance based on the specified Rect. It initializes pixel data and associated properties.
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

// ComputeIndex calculates the linear index in the pixel array for the given (x, y) coordinates.
// It adjusts the y-coordinate based on the lastY property and computes the index using stride and rectangle bounds.
func (s *Picture) ComputeIndex(x int, y int) int {
	y = s.lastY - y
	i := (y-int(s.rect.Min.Y))*s.stride + (x-int(s.rect.Min.X))*4
	return i
}

// SetRGBAArray sets the RGBA color data at the specified (x, y) coordinates using a given array of uint8 values.
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
func (s *Picture) SetRGBADirectArray(i int, rgba []uint8) {
	copy(s.pixels[i:], rgba)
}

// SetRGBA4DirectArray sets the RGBA color values starting at the given pixel index directly in the pixel data array.
func (s *Picture) SetRGBA4DirectArray(i int, rgba [4]uint8) {
	copy(s.pixels[i:], rgba[:])
}

// SetRGBA8DirectArray sets the RGBA color values starting at the given pixel index directly in the pixel data array.
func (s *Picture) SetRGBA8DirectArray(i int, rgba [8]uint8) {
	copy(s.pixels[i:], rgba[:])
}

func (s *Picture) SetRGBA32DirectArray(i int, rgba [32]uint8) {
	copy(s.pixels[i:], rgba[:])
}

// SetRGBADirectArrayPtr sets the RGBA color values starting at the given pixel index directly in the pixel data array.
func (s *Picture) SetRGBADirectArrayPtr(i int, rgba *[]uint8) {
	copy(s.pixels[i:], *rgba)
}

// SetRGBA sets the RGBA color value for the pixel at the specified x, y coordinates within the image bounds.
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

// SetRGBASize sets a square block of pixels centered at (x, y) to the specified RGBA color with the given size.
func (s *Picture) SetRGBASize(x int, y int, r uint8, g uint8, b uint8, a uint8, size int) {
	k := size / 2
	for offsetX := -k; offsetX < k; offsetX++ {
		for offsetY := -k; offsetY < k; offsetY++ {
			s.SetRGBA(x+offsetX, y+offsetY, r, g, b, a)
		}
	}
}

// Width returns the width of the Picture in pixels.
func (s *Picture) Width() int {
	return s.width
}

// Height returns the height of the Picture in pixels.
func (s *Picture) Height() int {
	return s.height
}

// Length returns the total length in pixels of the Picture.
func (s *Picture) Length() int {
	return s.length
}

// Stride returns the number of bytes between successive rows of pixels in the Picture.
func (s *Picture) Stride() int {
	return s.stride
}

// Bounds returns the rectangular domain of the Picture as a Rect structure.
func (s *Picture) Bounds() Rect {
	return s.rect
}

// Pixels returns the bounding box coordinates (bx, by, bw, bh) and the pixel data of the Picture.
func (s *Picture) Pixels() (int, int, int, int, []uint8) {
	return s.bx, s.by, s.bw, s.bh, s.pixels
}

// Image converts the Picture object into an *image.RGBA, using its internal pixel data and defined bounds.
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

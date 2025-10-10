package pixels

import (
	"image/color"
	"math"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font/basicfont"
)

// ASCII is a slice of runes representing printable ASCII characters starting from 32 to unicode.MaxASCII.
var ASCII []rune

// Atlas7x13 is a predefined *Atlas instance created with a fixed font face and a set of ASCII characters.
var Atlas7x13 *Atlas

// init initializes the ASCII slice with printable ASCII characters and creates the Atlas7x13 Atlas using a basic font.
func init() {
	ASCII = make([]rune, unicode.MaxASCII-32)
	for i := range ASCII {
		ASCII[i] = rune(32 + i)
	}
	Atlas7x13 = NewAtlas(basicfont.Face7x13, ASCII)
}

// RangeTable converts a given unicode.RangeTable into a slice of runes representing all characters in the table.
func RangeTable(table *unicode.RangeTable) []rune {
	var runes []rune
	for _, rng := range table.R16 {
		for r := rng.Lo; r <= rng.Hi; r += rng.Stride {
			runes = append(runes, rune(r))
		}
	}
	for _, rng := range table.R32 {
		for r := rng.Lo; r <= rng.Hi; r += rng.Stride {
			runes = append(runes, rune(r))
		}
	}
	return runes
}

// Text is a struct that represents a formatted text object for manipulation and rendering.
type Text struct {
	// Orig specifies the text origin, usually the top-left dot position. Dot is always aligned
	// to Orig when writing newlines.
	Orig Vector

	// Dot is the position where the next character will be written. Dot is automatically moved
	// when writing to a Text object, but you can also manipulate it manually
	Dot Vector

	// Color is the color of the text that is to be written. Defaults to white.
	Color color.Color

	// LineHeight is the vertical distance between two lines of text.
	//
	// Example:
	//   txt.LineHeight = 1.5 * txt.Atlas().LineHeight()
	LineHeight float64

	// TabWidth is the horizontal tab width. Tab characters will align to the multiples of this
	// width.
	//
	// Example:
	//   txt.TabWidth = 8 * txt.Atlas().Glyph(' ').Advance
	TabWidth float64

	atlas *Atlas

	buf    []byte
	prevR  rune
	bounds Rect
	glyph  TrianglesData
	tris   TrianglesData

	mat    Matrix
	col    RGBA
	trans  TrianglesData
	transD Drawer
	dirty  bool
	anchor Anchor
}

// NewText creates a new Text object initialized with the provided origin vector and font atlas.
func NewText(orig Vector, atlas *Atlas) *Text {
	txt := &Text{
		Orig:       orig,
		Dot:        orig,
		Color:      Alpha(1),
		LineHeight: atlas.LineHeight(),
		TabWidth:   atlas.Glyph(' ').Advance * 4,
		atlas:      atlas,
		mat:        IM,
		col:        Alpha(1),
	}

	txt.glyph.SetLen(6)
	for i := range txt.glyph {
		txt.glyph[i].Color = Alpha(1)
		txt.glyph[i].Intensity = 1
	}

	txt.transD.SetPicture(txt.atlas.pic)
	txt.transD.SetTriangles(&txt.trans)
	txt.transD.SetCacheMode(CacheModePicture)

	txt.Clear()

	return txt
}

// Atlas returns the *Atlas associated with the Text object, which contains pre-drawn glyphs for efficient text rendering.
func (txt *Text) Atlas() *Atlas {
	return txt.atlas
}

// Bounds returns the rectangular boundary of the text drawn so far.
func (txt *Text) Bounds() Rect {
	return txt.bounds
}

// BoundsOf calculates and returns the bounding rectangle that encompasses the given string when drawn with the current text settings.
func (txt *Text) BoundsOf(s string) Rect {
	dot := txt.Dot
	prevR := txt.prevR
	bounds := Rect{}

	for _, r := range s {
		var control bool
		dot, control = txt.controlRune(r, dot)
		if control {
			continue
		}

		var b Rect
		_, _, b, dot = txt.Atlas().DrawRune(prevR, r, dot)

		if bounds.W()*bounds.H() == 0 {
			bounds = b
		} else {
			bounds = bounds.Union(b)
		}

		prevR = r
	}

	return bounds
}

// AlignedTo sets the alignment anchor for the text and returns the modified Text object.
func (txt *Text) AlignedTo(anchor Anchor) *Text {
	txt.anchor = anchor
	return txt
}

// Clear resets the Text object by clearing its bounds, glyphs, and previous rune, marking it as dirty, and aligning Dot to Orig.
func (txt *Text) Clear() {
	txt.prevR = -1
	txt.bounds = Rect{}
	txt.tris.SetLen(0)
	txt.dirty = true
	txt.Dot = txt.Orig
}

// Write appends the provided byte slice to the internal buffer and updates the rendered text. Returns the number of bytes written.
func (txt *Text) Write(p []byte) (n int, err error) {
	txt.buf = append(txt.buf, p...)
	txt.drawBuf()
	return len(p), nil
}

// WriteString appends the provided string to the internal buffer and updates the text display.
// It returns the number of bytes written and an error, if any.
func (txt *Text) WriteString(s string) (n int, err error) {
	txt.buf = append(txt.buf, s...)
	txt.drawBuf()
	return len(s), nil
}

// WriteByte appends a single byte to the Text buffer and updates the rendering. Returns an error if writing fails.
func (txt *Text) WriteByte(c byte) error {
	txt.buf = append(txt.buf, c)
	txt.drawBuf()
	return nil
}

// WriteRune appends the UTF-8 encoding of the given rune to the text buffer and updates the drawing state.
// It returns the number of bytes written and an error (always nil).
func (txt *Text) WriteRune(r rune) (n int, err error) {
	var b [4]byte
	n = utf8.EncodeRune(b[:], r)
	txt.buf = append(txt.buf, b[:n]...)
	txt.drawBuf()
	return n, nil
}

// Draw renders the text object onto the specified target using the given transformation matrix.
func (txt *Text) Draw(t ITarget, matrix Matrix) {
	txt.DrawColorMask(t, matrix, nil)
}

// DrawColorMask renders the text onto the specified target using the given transformation matrix and color mask.
func (txt *Text) DrawColorMask(t ITarget, matrix Matrix, mask color.Color) {
	if matrix != txt.mat {
		txt.mat = matrix
		txt.dirty = true
	}

	offset := txt.Orig.Sub(txt.Bounds().Max.Add(txt.Bounds().AnchorPos(txt.anchor.Opposite())))
	txt.mat = IM.Moved(offset).Chained(txt.mat)

	if mask == nil {
		mask = Alpha(1)
	}
	rgba := ToRGBA(mask)
	if rgba != txt.col {
		txt.col = rgba
		txt.dirty = true
	}

	if txt.dirty {
		txt.trans.SetLen(txt.tris.Len())
		txt.trans.Update(&txt.tris)

		for i := range txt.trans {
			txt.trans[i].Position = txt.mat.Project(txt.trans[i].Position)
			txt.trans[i].Color = txt.trans[i].Color.Mul(txt.col)
		}

		txt.transD.Dirty()
		txt.dirty = false
	}

	txt.transD.Draw(t)
}

// controlRune processes special control runes like newline, carriage return, or tab, adjusting the text position accordingly.
// It returns the updated dot position and a boolean indicating if the rune was a control character.
func (txt *Text) controlRune(r rune, dot Vector) (newDot Vector, control bool) {
	switch r {
	case '\n':
		dot.X = txt.Orig.X
		dot.Y -= txt.LineHeight
	case '\r':
		dot.X = txt.Orig.X
	case '\t':
		rem := math.Mod(dot.X-txt.Orig.X, txt.TabWidth)
		rem = math.Mod(rem, rem+txt.TabWidth)
		if rem == 0 {
			rem = txt.TabWidth
		}
		dot.X += rem
	default:
		return dot, false
	}
	return dot, true
}

// drawBuf processes the text buffer to decode runes, handle control characters, and generate drawing data for glyphs.
func (txt *Text) drawBuf() {
	if !utf8.FullRune(txt.buf) {
		return
	}

	rgba := ToRGBA(txt.Color)
	for i := range txt.glyph {
		txt.glyph[i].Color = rgba
	}

	for utf8.FullRune(txt.buf) {
		r, size := utf8.DecodeRune(txt.buf)
		txt.buf = txt.buf[size:]

		var control bool
		txt.Dot, control = txt.controlRune(r, txt.Dot)
		if control {
			continue
		}

		var rect, frame, bounds Rect
		rect, frame, bounds, txt.Dot = txt.Atlas().DrawRune(txt.prevR, r, txt.Dot)

		txt.prevR = r

		rv := [...]Vector{
			{X: rect.Min.X, Y: rect.Min.Y},
			{X: rect.Max.X, Y: rect.Min.Y},
			{X: rect.Max.X, Y: rect.Max.Y},
			{X: rect.Min.X, Y: rect.Max.Y},
		}

		fv := [...]Vector{
			{X: frame.Min.X, Y: frame.Min.Y},
			{X: frame.Max.X, Y: frame.Min.Y},
			{X: frame.Max.X, Y: frame.Max.Y},
			{X: frame.Min.X, Y: frame.Max.Y},
		}

		for i, j := range [...]int{0, 1, 2, 0, 2, 3} {
			txt.glyph[i].Position = rv[j]
			txt.glyph[i].Picture = fv[j]
		}

		txt.tris = append(txt.tris, txt.glyph...)
		txt.dirty = true

		if txt.bounds.W()*txt.bounds.H() == 0 {
			txt.bounds = bounds
		} else {
			txt.bounds = txt.bounds.Union(bounds)
		}
	}
}

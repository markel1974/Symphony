package pixels

import (
	"image"
	"image/draw"
	"sort"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Glyph represents a graphical character or symbol with its position, bounds, and spacing in 2D space.
type Glyph struct {
	Dot     Vector
	Frame   Rect
	Advance float64
}

// Atlas represents a font atlas that maps runes to glyphs with their metrics and images.
type Atlas struct {
	face       font.Face
	pic        IPicture
	mapping    map[rune]Glyph
	ascent     float64
	descent    float64
	lineHeight float64
}

// NewAtlas creates and returns a new Atlas using the provided font face and rune sets for glyph mapping and rendering.
func NewAtlas(face font.Face, runeSets ...[]rune) *Atlas {
	seen := make(map[rune]bool)
	runes := []rune{unicode.ReplacementChar}
	for _, set := range runeSets {
		for _, r := range set {
			if !seen[r] {
				runes = append(runes, r)
				seen[r] = true
			}
		}
	}

	fixedMapping, fixedBounds := makeSquareMapping(face, runes, fixed.I(2))

	atlasImg := image.NewRGBA(image.Rect(
		fixedBounds.Min.X.Floor(),
		fixedBounds.Min.Y.Floor(),
		fixedBounds.Max.X.Ceil(),
		fixedBounds.Max.Y.Ceil(),
	))

	for r, fg := range fixedMapping {
		if dr, mask, maskP, _, ok := face.Glyph(fg.dot, r); ok {
			draw.Draw(atlasImg, dr, mask, maskP, draw.Src)
		}
	}

	bounds := NewRect(
		i2f(fixedBounds.Min.X),
		i2f(fixedBounds.Min.Y),
		i2f(fixedBounds.Max.X),
		i2f(fixedBounds.Max.Y),
	)

	mapping := make(map[rune]Glyph)
	for r, fg := range fixedMapping {
		mapping[r] = Glyph{
			Dot: NewVec(
				i2f(fg.dot.X),
				bounds.Max.Y-(i2f(fg.dot.Y)-bounds.Min.Y),
			),
			Frame: NewRect(
				i2f(fg.frame.Min.X),
				bounds.Max.Y-(i2f(fg.frame.Min.Y)-bounds.Min.Y),
				i2f(fg.frame.Max.X),
				bounds.Max.Y-(i2f(fg.frame.Max.Y)-bounds.Min.Y),
			).Norm(),
			Advance: i2f(fg.advance),
		}
	}

	return &Atlas{
		face:       face,
		pic:        NewPictureFromImage(atlasImg),
		mapping:    mapping,
		ascent:     i2f(face.Metrics().Ascent),
		descent:    i2f(face.Metrics().Descent),
		lineHeight: i2f(face.Metrics().Height),
	}
}

// Picture returns the IPicture associated with the Atlas.
func (a *Atlas) Picture() IPicture {
	return a.pic
}

// Contains checks if the atlas contains a glyph mapping for the specified rune. Returns true if the rune is present.
func (a *Atlas) Contains(r rune) bool {
	_, ok := a.mapping[r]
	return ok
}

// Glyph returns the Glyph associated with the given rune from the Atlas's mapping.
func (a *Atlas) Glyph(r rune) Glyph {
	return a.mapping[r]
}

// Kern calculates the kerning adjustment between two runes r0 and r1 in the font, returning the spacing offset as a float64.
func (a *Atlas) Kern(r0, r1 rune) float64 {
	return i2f(a.face.Kern(r0, r1))
}

// Ascent returns the ascent of the font represented by the Atlas as a float64 value.
func (a *Atlas) Ascent() float64 {
	return a.ascent
}

// Descent returns the descent value of the font in the atlas, representing the vertical distance below the baseline.
func (a *Atlas) Descent() float64 {
	return a.descent
}

// LineHeight returns the line height of the font in the atlas as a float64.
func (a *Atlas) LineHeight() float64 {
	return a.lineHeight
}

// DrawRune draws the rune `r` at the specified position `dot`, applying kerning relative to `prevR`.
// It returns the rectangle where the rune is drawn, the glyph's frame, its bounding box, and the updated position `newDot`.
// If the rune is not found in the atlas, the replacement character is used.
func (a *Atlas) DrawRune(prevR, r rune, dot Vector) (rect, frame, bounds Rect, newDot Vector) {
	if !a.Contains(r) {
		r = unicode.ReplacementChar
	}
	if !a.Contains(unicode.ReplacementChar) {
		return Rect{}, Rect{}, Rect{}, dot
	}
	if !a.Contains(prevR) {
		prevR = unicode.ReplacementChar
	}

	if prevR >= 0 {
		dot.X += a.Kern(prevR, r)
	}

	glyph := a.Glyph(r)

	rect = glyph.Frame.Moved(dot.Sub(glyph.Dot))
	bounds = rect

	if bounds.W()*bounds.H() != 0 {
		bounds = NewRect(
			bounds.Min.X,
			dot.Y-a.Descent(),
			bounds.Max.X,
			dot.Y+a.Ascent(),
		)
	}

	dot.X += glyph.Advance

	return rect, glyph.Frame, bounds, dot
}

// fixedGlyph represents a single glyph with its associated properties for rendering in fixed-point coordinates.
type fixedGlyph struct {
	dot     fixed.Point26_6
	frame   fixed.Rectangle26_6
	advance fixed.Int26_6
}

// makeSquareMapping generates a mapping of runes to their fixedGlyph representation ensuring square bounds.
// It adjusts the width dynamically to maintain a square aspect ratio of the overall bounding rectangle.
// Accepts a font.Face, a slice of runes, and padding as parameters.
// Returns a map of runes to fixedGlyphs and the overall bounding rectangle as fixed.Rectangle26_6.
func makeSquareMapping(face font.Face, runes []rune, padding fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {
	width := sort.Search(int(fixed.I(1024*1024)), func(i int) bool {
		width := fixed.Int26_6(i)
		_, bounds := makeMapping(face, runes, padding, width)
		return bounds.Max.X-bounds.Min.X >= bounds.Max.Y-bounds.Min.Y
	})
	return makeMapping(face, runes, padding, fixed.Int26_6(width))
}

// makeMapping generates a mapping of runes to glyph data and computes the bounding rectangle for all glyphs.
// face specifies the font face used for glyph metrics.
// runes is the set of unicode characters to include in the mapping.
// padding adds spacing between glyphs.
// width sets the maximum line width before wrapping to a new row.
// Returns a map where each rune links to its fixedGlyph and the overall bounding rectangle.
func makeMapping(face font.Face, runes []rune, padding, width fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {
	mapping := make(map[rune]fixedGlyph)
	bounds := fixed.Rectangle26_6{}

	dot := fixed.P(0, 0)

	for _, r := range runes {
		b, advance, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}

		// this is important for drawing, artifacts arise otherwise
		frame := fixed.Rectangle26_6{
			Min: fixed.P(b.Min.X.Floor(), b.Min.Y.Floor()),
			Max: fixed.P(b.Max.X.Ceil(), b.Max.Y.Ceil()),
		}

		dot.X -= frame.Min.X
		frame = frame.Add(dot)

		mapping[r] = fixedGlyph{
			dot:     dot,
			frame:   frame,
			advance: advance,
		}
		bounds = bounds.Union(frame)

		dot.X = frame.Max.X

		// padding + align to integer
		dot.X += padding
		dot.X = fixed.I(dot.X.Ceil())

		// width exceeded, new row
		if frame.Max.X >= width {
			dot.X = 0
			dot.Y += face.Metrics().Ascent + face.Metrics().Descent

			// padding + align to integer
			dot.Y += padding
			dot.Y = fixed.I(dot.Y.Ceil())
		}
	}

	return mapping, bounds
}

// i2f converts a fixed.Int26_6 value to a float64 by dividing it by 64.
func i2f(i fixed.Int26_6) float64 {
	return float64(i) / (1 << 6)
}

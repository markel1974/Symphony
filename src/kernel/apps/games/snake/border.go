package snake

import "github.com/markel1974/symphony/src/kernel/interfaces"

// Border represents a rectangular boundary for a two-dimensional plane, defined by its width, height, and coordinate map.
type Border struct {
	width  int
	height int
	coords map[Point]Data
}

// NewBorder initializes and returns a new Border with specified origin coordinates, width, and height.
// It constructs the border using '█' characters with green foreground and no background color.
// The resulting border dimensions are adjusted by (-1, -2) for width and height respectively.
// The function populates the Border's coordinates map with relevant points and attributes.
func NewBorder(origX int, origY int, width, height int) *Border {
	b := new(Border)
	b.width, b.height = width-1, height-2
	b.coords = make(map[Point]Data)

	for x := origX; x < b.width; x++ {
		p1 := Point{X: x, Y: origY}
		p2 := Point{X: x, Y: b.height}

		b.coords[p1] = Data{Point: p1, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
		b.coords[p2] = Data{Point: p2, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
	}

	for y := origY; y < b.height+1; y++ {
		p1 := Point{X: origX, Y: y}
		p2 := Point{X: b.width, Y: y}

		b.coords[p1] = Data{Point: p1, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
		b.coords[p2] = Data{Point: p2, Attributes: Attributes{C: '█', Fg: interfaces.ColorGreenDef, Bg: interfaces.ColorNoneDef}}
	}

	return b
}

// Contains checks if the specified point is within the border's coordinates. Returns true if the point is found.
func (b *Border) contains(point Point) bool {
	_, exists := b.coords[point]
	return exists
}

// Draw renders the border on the specified ISurface by iterating through its coordinates and applying colors to each point.
func (b *Border) Draw(s interfaces.ISurface) {
	for _, c := range b.coords {
		s.DrawColor(c.Y, c.X, c.C, c.Fg, c.Bg, interfaces.ModeNormal)
	}
}

package invaders

import (
	"strings"

	"github.com/markel1974/symphony/src/kernel/interfaces"
)

// barricadeSymbol defines the character used to represent a barricade.
// numBarricades specifies the total number of barricades.
// fgBarricade sets the foreground color of barricades.
// bgBarricade sets the background color of barricades.
// barricadeSpriteWidth defines the width of the barricade sprite.
// barricadeSpriteHeight defines the height of the barricade sprite.
// barricadeSprite specifies the visual representation of the barricade.
const (
	barricadeSymbol       = '▓'
	numBarricades         = 4
	fgBarricade           = interfaces.ColorWhiteDef
	bgBarricade           = interfaces.ColorGreenDef
	barricadeSpriteWidth  = 11
	barricadeSpriteHeight = 5
	barricadeSprite       = `    xxx
  xxxxxxx
xxxxxxxxxxx
xxx     xxx
xxx     xxx`
)

// barricadeYPos calculates the vertical position for placing barricades based on screen height and sprite dimensions.
func barricadeYPos(h int) int {
	//TODO pass playerSpriteHeight
	var playerSpriteHeight = 3
	return (h - playerSpriteBottomOffset - playerSpriteHeight) - barricadeSpriteHeight - 2
}

// Barricade represents a defensive structure in the game, used for shielding and interacting with game entities.
// It maintains width, height, and a data grid to hold the state of individual barricade positions.
type Barricade struct {
	w    int
	h    int
	data [][]int
}

// NewBarricade creates and initializes a new Barricade instance and returns its pointer.
func NewBarricade() *Barricade {
	b := &Barricade{}
	return b
}

// Setup initializes the barricade structure with specified dimensions and creates barricade sprites on the playfield.
func (b *Barricade) Setup(w int, h int) *Barricade {
	b.w = w
	b.h = h

	b.data = make([][]int, b.w)
	for i := range b.data {
		b.data[i] = make([]int, b.h)
		for j := range b.data[i] {
			b.data[i][j] = nonIndex
		}
	}

	gap := (w - 4*barricadeSpriteWidth) / 5
	x := gap
	y := barricadeYPos(h)
	initY := y
	lines := strings.Split(barricadeSprite, "\n")
	for i := 0; i < numBarricades; i++ {
		initX := gap*(i+1) + barricadeSpriteWidth*i
		x = initX
		for _, l := range lines {
			for _, c := range l {
				if c != ' ' {
					b.Set(x, y, i)
				}
				x++
			}
			y++
			x = initX
		}
		y = initY
	}
	return b
}

// Unset clears the data at the specified (x, y) coordinates, setting it to nonIndex. Returns true if data was altered.
func (b *Barricade) Unset(x int, y int) bool {
	var ret = false
	if b.Get(x, y) != nonIndex {
		b.Set(x, y, nonIndex)
		ret = true
	}

	return ret
}

// Draw renders the barricade on the provided surface using specified colors and symbols for non-empty cells.
func (b *Barricade) Draw(surface interfaces.ISurface) {
	for i := range b.data {
		for j := range b.data[i] {
			if b.data[i][j] != nonIndex {
				surface.DrawColor(j, i, barricadeSymbol, fgBarricade, bgBarricade, interfaces.ModeNormal)
			}
		}
	}
}

// Set updates the value at the specified (x, y) position with the given data if the position is within bounds.
func (b *Barricade) Set(x int, y int, data int) {
	if len(b.data) > 0 {
		if x >= 0 && y >= 0 && x < b.w && y < b.h {
			b.data[x][y] = data
		}
	}
}

// Get retrieves the value at the specified x and y coordinates within the barricade data grid. Returns nonIndex if out of bounds.
func (b *Barricade) Get(x int, y int) int {
	data := nonIndex
	if len(b.data) > 0 {
		if x >= 0 && y >= 0 && x < b.w && y < b.h {
			data = b.data[x][y]
		}
	}
	return data
}

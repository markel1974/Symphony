package invaders

import (
	"math/rand"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/servers/render/matrix"
)

// starSymbol represents the character used as a symbol for a star.
// numStars defines the total number of stars.
const (
	starSymbol = '.'
	numStars   = 50
)

// NewStar creates a new star entity with specified position, symbol, color, and speed properties.
func NewStar(x int, y int, symbol rune, color string, speed float64) *matrix.Entity {
	s := matrix.NewEntity(symbol, x, y, 0.5, speed, speed, []string{color}, -1)
	return s
}

// Stars represents a collection of entities and their spatial hierarchy within an AABB tree for efficient operations.
type Stars struct {
	data []*matrix.Entity
	tree *matrix.AABBTree
}

// NewStars creates and initializes a new Stars instance with a pre-allocated AABBTree and an entity slice.
func NewStars() *Stars {
	s := &Stars{
		tree: matrix.NewAABBTree(numStars),
		data: make([]*matrix.Entity, 0, numStars),
	}
	return s
}

// Update adds, updates, or replaces stars in the collection, ensuring they remain within boundaries and handle collisions.
func (s *Stars) Update(w int, h int, fc uint8) {
	if len(s.data) != cap(s.data) && fc%3 == 0 {
		n := len(s.data)
		s.data = s.data[0 : n+1]
		var symbol, speed, color = s.randSymbol()
		star := NewStar(rand.Intn(w), 0, symbol, color, speed)
		s.data[n] = star
	}

	for i, obj1 := range s.data {
		obj1.AddTo(obj1.Vx, obj1.Vy)

		s.tree.RemoveObject(obj1)

		if obj1.GetX() < 0 || int(obj1.GetX()) > w || obj1.GetY() < 0 || int(obj1.GetY()) > h {
			var symbol, speed, color = s.randSymbol()
			star := NewStar(rand.Intn(w), rand.Intn(h), symbol, color, speed)
			s.data[i] = star
			s.tree.InsertObject(star)
		} else {
			s.tree.InsertObject(obj1)
		}

		k := s.tree.QueryOverlaps(obj1)
		if len(k) > 0 {
			obj1.DoPhysics(obj1)
			for _, z := range k {
				obj1.DoPhysics(z.(*matrix.Entity))
			}
		}
	}
}

// Draw iterates over all entities in the Stars data and renders them onto the provided ISurface instance.
func (s *Stars) Draw(surface interfaces.ISurface) {
	for _, s := range s.data {
		s.Draw(surface)
		//surface.DrawEntity(s)
	}
}

// randColor generates a random color by selecting one of six predefined color constants from the ColorDef type.
func (s *Stars) randColor() interfaces.ColorDef {
	z := rand.Intn(6)
	color := interfaces.ColorWhiteDef

	switch z {
	case 0:
		color = interfaces.ColorRedDef
	case 1:
		color = interfaces.ColorYellowDef
	case 2:
		color = interfaces.ColorBlueDef
	case 3:
		color = interfaces.ColorMagentaDef
	case 4:
		color = interfaces.ColorCyanDef
	case 5:
		color = interfaces.ColorWhiteDef
	}

	return color
}

// randSymbol generates a random symbol, velocity, and color combination for a star object.
// Returns a rune representing the symbol, a float64 velocity value, and a string for the color.
func (s *Stars) randSymbol() (rune, float64, string) {
	var z = rand.Intn(9)
	var symbol = starSymbol
	var vx = -0.1
	var color = "W"

	switch z {
	case 1:
		symbol = '.'
		vx = -0.1
		color = "W"
	case 2:
		symbol = '.'
		vx = 0.1
		color = "W"
	case 3:
		symbol = '*'
		vx = -0.5
		color = "B"
	case 5:
		symbol = '*'
		vx = 0.5
		color = "B"
	case 7:
		symbol = '+'
		vx = -0.3
		color = "C"
	case 8:
		symbol = '+'
		vx = 0.3
		color = "C"
	}

	return symbol, vx, color
}

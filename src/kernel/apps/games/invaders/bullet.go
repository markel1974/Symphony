package invaders

import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/servers/render/matrix"
)

// bulletSprite defines the visual representation of a bullet as a slice of strings, where each string corresponds to a sprite frame.
var (
	bulletSprite = []string{"R"}
)

// Bullet represents a specific type of Entity used to model projectiles in the game.
// It embeds the matrix.Entity type to inherit its physical and visual properties.
type Bullet struct {
	*matrix.Entity
}

// NewBullet creates a new Bullet instance at position (x, y) with a vertical velocity vy initialized using the bulletSprite.
func NewBullet(x int, y int, vy int) *Bullet {
	b := &Bullet{Entity: matrix.NewEntity(-1, x, y, 2, 0, float64(vy), bulletSprite, 0)}
	return b
}

// Bullets represents a collection of Bullet objects, typically used to manage and control multiple bullets in a game scenario.
type Bullets struct {
	data []*Bullet
}

// NewBullets initializes and returns a new instance of the Bullets struct.
func NewBullets() *Bullets {
	return &Bullets{}
}

// Setup initializes the Bullets instance by allocating memory for the specified number of Bullet elements.
func (ab *Bullets) Setup(num int) {
	ab.data = make([]*Bullet, num)
}

// Draw iterates over all non-nil bullets and renders them onto the specified surface, advancing their state afterward.
func (ab *Bullets) Draw(surface interfaces.ISurface) {
	for _, k := range ab.data {
		if k != nil {
			k.Draw(surface)
			k.Next()
		}
	}
}

// Update assigns a new bullet at the specified position if an empty slot is available in the bullets array.
func (ab *Bullets) Update(x int, y int) {
	for j := range ab.data {
		if ab.data[j] == nil {
			ab.data[j] = NewBullet(x, y, alienBulletSpeed)
			break
		}
	}
}

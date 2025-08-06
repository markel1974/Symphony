package invaders

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/render/matrix"
)

// alienStartX defines the starting X coordinate for the alien.
// alienStartY defines the starting Y coordinate for the alien.
// alienBulletSpeed specifies the speed of the alien's bullet.
// alienShootValMax sets the maximum randomized value for alien shooting logic.
// alienPadVertical determines the vertical padding around the alien.
// alienPadHorizontal determines the horizontal padding around the alien.
// alienMoveEvery specifies the number of steps after which the alien moves.
const (
	alienStartX        = 10
	alienStartY        = 7
	alienBulletSpeed   = 1
	alienShootValMax   = 100
	alienPadVertical   = 1
	alienPadHorizontal = 3
	alienMoveEvery     = 4
)

// ufoMoveEvery defines the interval at which the UFO moves.
// ufoReward specifies the reward points for destroying a UFO.
const (
	ufoMoveEvery = 3
	ufoReward    = 100
)

// ufoSprite represents a slice of ASCII art strings used for animating a UFO sprite in the game.
var (
	ufoSprite = []string{
		`
 YYYYYY
GMMMMMMM
 BB  BB`,
		`
 YYYYYY
MMGMMMMM
 BB  BB`,
		`
 YYYYYY
MMMMGMMM
 BB  BB`,
		`
 YYYYYY
MMMMMMGM
 BB  BB`,
		`
 YYYYYY
MMMMMMMM
 BB  BB`,
		`
 YYYYYY
MMMMMMMM
 BB  BB`,
		`
 YYYYYY
MMMMMMMM
 BB  BB`,
	}
)

// smAlienSprite represents a collection of small alien sprite ASCII art variations.
// mdAlienSprite represents a collection of medium alien sprite ASCII art variations.
// lgAlienSprite represents a collection of large alien sprite ASCII art variations.
var (
	smAlienSprite = []string{
		`
  YMMM
 GYYYYY
   CC
`,
		`
  MYMM
 YGYYYY
  C  C
`,
		`
  MMYM
 YYGYYY
 C    C
`,
		`
  MMMY
 YYYGYY
  C  C
`,
		`
  MMMM
 YYYYGY
   CC
`,
		`
  MMMM
 YYYYYG
  C  C
`,
		`
  MMMM
 YYYYYY
 C    C
`,
		`
  MMMM
 YYYYYY
  C  C
`,
	}

	mdAlienSprite = []string{
		`
G  BB  G
 RYYYYY
  MMMM
  B  B
`,
		`
   BB
GYRYYYYG
  MMMM 
  B  B
`,
		`
   BB
 YYRYYY
G MMMM G
  B  B
`,
		`
   BB
GYYYRYYG
  MMMM 
  B  B
`,
		`
G  BB  G
 YYYYRY
  MMMM 
  B  B
`,
		`
   BB  
GYYYYYRG
  MMMM 
  B  B
`,
		`
   BB  
 YYYYYY
G MMMM G
  B  B
`,
		`
   BB  
GYYYYYYG
  MMMM 
  B  B
`,
	}

	lgAlienSprite = []string{
		`
  CWWW
 GYYYYY
 GYYYYY
   WW
`,
		`
  WCWW
 YYGYYY
 YYYGYY
  W  W
`,
		`
  WWCW
 YYYYGY
 YYYYYG
 W    W
`,
		`
  WWWC
 YYYYYY
 YYYYYY
  W  W
`,
		`
  WWWW
 YYYYYY
 YYYYYY
   WW
`,
		`
  WWWW
 YYYYYY
 YYYYYY
  W  W
`,
		`
  WWWW
 YYYYYY
 YYYYYY
 W    W
`,
		`
  WWWW
 YYYYYY
 YYYYYY
  W  W
`,
	}
)

// Alien represents an alien entity in the game, inheriting properties and methods from matrix.Entity.
// It includes attributes for indexing, rewards, and a health counter to manage its state and interactions.
type Alien struct {
	*matrix.Entity
	idx     int
	reward  int
	counter int
}

// NewAlien creates and returns a new Alien instance with the specified index, position, sprite, and reward.
func NewAlien(idx int, x int, y int, sprite []string, reward int) *Alien {
	return &Alien{
		Entity:  matrix.NewEntity(-1, x, y, 1, 1, 0, sprite, -1),
		idx:     idx,
		reward:  reward,
		counter: 3,
	}
}

// newUfo creates a new UFO entity with predefined properties and positions it off-screen at the top-left.
func newUfo() *matrix.Entity {
	u := matrix.NewEntity(-1, 0, 0, 1, 0.5, 0.5, ufoSprite, -1)
	u.MoveTo(0-u.GetWidth(), u.GetHeight())
	return u
}

// Aliens represents a collection of alien entities, their movement parameters, state, and related functionalities.
// It includes properties for positioning, size, velocity, and interaction with bullets, as well as spatial management.
type Aliens struct {
	container []*Alien

	columns int
	rows    int
	rowsSm  int
	rowsMd  int
	rowsLg  int
	alienV  [2]int
	altered bool

	rightMove    [2]int
	leftMove     [2]int
	downMove     [2]int
	alienBullets *Bullets

	count int
	minY  int
	minX  int

	w    int
	h    int
	tree *matrix.AABBTree
}

// NewAliens initializes and returns a new instance of the Aliens struct with default values and configuration.
func NewAliens() *Aliens {
	a := &Aliens{}

	a.rightMove = [2]int{1, 0}
	a.leftMove = [2]int{-1, 0}
	a.downMove = [2]int{0, 1}

	a.rowsSm = 2
	a.rowsMd = 2
	a.rowsLg = 1
	a.columns = 1
	a.rows = a.rowsSm + a.rowsMd + a.rowsLg
	a.altered = false
	a.alienV = a.rightMove
	a.alienBullets = NewBullets()

	a.count = 0

	return a
}

// Setup initializes the Aliens instance by determining grid dimensions, rows per size category, and setting up containers.
func (a *Aliens) Setup(w int, h int) {
	var i = 0
	var maxAlienWidth = 8
	var maxAlienHeight = 4

	a.columns = (int(float64(w) * 0.7)) / (maxAlienWidth + alienPadHorizontal)

	i = barricadeYPos(h) / (maxAlienHeight + alienPadVertical)

	switch i {
	case 0:
		a.rowsSm = 0
		a.rowsMd = 0
		a.rowsLg = 1
	case 1:
		a.rowsSm = 0
		a.rowsMd = 0
		a.rowsLg = 1
	case 2:
		a.rowsSm = 0
		a.rowsMd = 0
		a.rowsLg = 1
	case 3:
		a.rowsSm = 0
		a.rowsMd = 1
		a.rowsLg = 1
	case 4:
		a.rowsSm = 1
		a.rowsMd = 1
		a.rowsLg = 1
	case 5:
		a.rowsSm = 1
		a.rowsMd = 2
		a.rowsLg = 1
	default:
		a.rowsSm = 2
		a.rowsMd = 2
		a.rowsLg = 1
	}

	a.rows = a.rowsSm + a.rowsMd + a.rowsLg

	a.container = make([]*Alien, a.columns*a.rows)

	a.tree = matrix.NewAABBTree(uint(a.columns * a.rows))

	a.alienBullets.Setup(a.columns * a.rows / 10)

	a.alienV = a.rightMove
}

// Create generates aliens in three rows (large, medium, small) with specified starting positions and updates offsets.
func (a *Aliens) Create(x int, y int, offset int) {
	y, offset = a.MakeAliens(a.columns, x, y, a.rowsLg, a.columns, rwdLg, offset, lgAlienSprite)
	y, offset = a.MakeAliens(a.columns, x, y, a.rowsMd, a.columns, rwdMd, offset, mdAlienSprite)
	a.MakeAliens(a.columns, x, y, a.rowsSm, a.columns, rwdSm, offset, smAlienSprite)
}

// Draw renders all aliens and their bullets onto the provided surface, advancing each alien to its next state.
func (a *Aliens) Draw(surface interfaces.ISurface) {
	for _, k := range a.container {
		if k != nil {
			k.Draw(surface)
			k.Next()
		}
	}
	a.alienBullets.Draw(surface)
}

// MakeAliens populates the alien container by creating and positioning aliens in a grid layout with specified parameters.
// Returns the updated vertical position and the next offset for the created aliens.
func (a *Aliens) MakeAliens(aliensHorizontal int, x int, y int, rows int, cols int, reward int, offset int, sprite []string) (int, int) {
	a.count = 0
	startX := x
	for i := 0; i < rows; i++ {
		var maxH = 0
		for j := 0; j < cols; j++ {
			idx := offset + i*aliensHorizontal + j
			if len(a.container) > 0 {
				if idx >= 0 && idx <= len(a.container) {
					alien := NewAlien(idx, x, y, sprite, reward)
					a.container[idx] = alien
					a.tree.InsertObject(alien)
					a.count++
					x += int(a.container[idx].GetWidth() + alienPadHorizontal)
					if int(a.container[idx].GetHeight()) > maxH {
						maxH = int(a.container[idx].GetHeight())
					}
				}
			}
		}
		x = startX
		y += maxH + alienPadVertical
	}
	return y, offset + rows*cols
}

// DoBulletCollision checks for collision between a bullet and aliens, computes rewards, and removes destroyed aliens.
func (a *Aliens) DoBulletCollision(bullet *matrix.Entity) (bool, int) {
	var reward = 0
	var found = false

	q := a.tree.QueryOverlaps(bullet)
	if len(q) > 0 {
		found = true
		for _, z := range q {
			var alien = z.(*Alien)
			alien.AddTo(bullet.Vx, bullet.Vy)
			bullet.DoPhysics(alien.Entity)
			alien.counter--
			if alien.counter <= 0 {
				reward = alien.reward
				a.tree.RemoveObject(alien)
				if alien.idx >= 0 && alien.idx < len(a.container) {
					a.container[alien.idx] = nil
				}
				a.count--
			}
		}
	}

	return found, reward
}

// DoMove updates the positions of all aliens and adjusts their movement direction if necessary.
func (a *Aliens) DoMove(w int, _ int) {
	var downFlag = false
	var alienCount = 0

	a.minX = 999999999
	a.minY = 999999999

	for _, alien := range a.container {
		if alien != nil {
			alienCount++
			if a.altered {
				alien.Vx = float64(a.alienV[0])
				alien.Vy = float64(a.alienV[1])
			}

			alien.AddTo(alien.Vx, alien.Vy)

			a.tree.RemoveObject(alien)
			a.tree.InsertObject(alien)

			q := a.tree.QueryOverlaps(alien)
			if len(q) > 0 {
				for _, z := range q {
					var target = z.(*Alien)
					alien.DoPhysics(target.Entity)
				}
			}

			alienX := int(alien.GetX())
			alienY := int(alien.GetY())
			if alienX <= 0 || alienX+int(alien.GetWidth()) >= w {
				downFlag = true
				if alienX < a.minX {
					a.minX = alienX
				}
			}

			if alienY < a.minY {
				a.minY = alienY
			}

			if rand.Intn(alienShootValMax) == 6 {
				a.alienBullets.Update(int(alien.GetX())+int(alien.GetWidth())/2, int(alien.GetY()))
			}
		}
	}

	a.altered = false

	switch {
	case a.alienV == a.downMove:
		if a.minX == 0 {
			a.alienV = a.rightMove
		} else {
			a.alienV = a.leftMove
		}
		a.altered = true
	case downFlag:
		a.alienV = a.downMove
		a.altered = true
	}
}

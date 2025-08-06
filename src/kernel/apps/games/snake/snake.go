package snake

import (
	"math/rand"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// direction represents a type used to specify movement or orientation in a discrete set of directions.
type direction int

// Up represents the upward direction.
// Down represents the downward direction.
// Left represents the leftward direction.
// Right represents the rightward direction.
const (
	Up direction = iota
	Down
	Left
	Right
)

// Point represents a coordinate in a 2D space with X and Y integer values.
type Point struct {
	X int
	Y int
}

// Attributes represents visual properties for a character, including its symbol and foreground/background colors.
type Attributes struct {
	C  rune
	Fg interfaces.ColorDef
	Bg interfaces.ColorDef
}

// Data represents a combination of positional information and visual attributes.
type Data struct {
	Point
	Attributes
}

// Snake represents the state and behavior of the snake in the game, including its position, direction, length, and components.
type Snake struct {
	Position   Point
	Direction  direction
	Length     int
	Body       []Data
	Speed      int
	Rows       int
	Columns    int
	Food       Data
	Score      int
	GameOver   bool
	borders    []*Border
	SpriteHead Attributes
	SpriteBody Attributes
	SpriteTail Attributes
}

// New creates and initializes a new Snake instance with default sprite attributes for its head, body, and tail.
func New() *Snake {
	snake := &Snake{
		SpriteTail: Attributes{C: '░', Fg: interfaces.ColorYellowDef, Bg: interfaces.ColorNoneDef},
		SpriteBody: Attributes{C: '░', Fg: interfaces.ColorMagentaDef, Bg: interfaces.ColorNoneDef},
		SpriteHead: Attributes{C: '@', Fg: interfaces.ColorRedDef, Bg: interfaces.ColorNoneDef},
	}
	return snake
}

// SetSize sets the dimensions of the game grid by setting the number of rows and columns and resets game borders and state.
func (snake *Snake) SetSize(rows int, column int) {
	snake.Rows = rows
	snake.Columns = column
	snake.borders = nil

	snake.Start()
}

// Start initializes the snake's game state, resets parameters, and places food and borders on the game surface.
func (snake *Snake) Start() {
	snake.GameOver = false
	snake.Speed = 1
	snake.Direction = Right
	snake.Body = []Data{
		{Point: Point{X: 1, Y: 6}, Attributes: snake.SpriteTail}, // Tail
		{Point: Point{X: 2, Y: 6}, Attributes: snake.SpriteBody}, // Body
		{Point: Point{X: 3, Y: 6}, Attributes: snake.SpriteHead}, // Head
	}
	snake.moveFood()
	snake.borders = nil
	snake.borders = append(snake.borders, NewBorder(0, 1, snake.Columns, snake.Rows))
	//snake.borders = append(snake.borders, NewBorder(Coordinates{7, 10}, 3, 5))
	//snake.borders = append(snake.borders, NewBorder(Coordinates{32, 5}, 3, 20))
}

// moveFood randomly positions the food inside the game boundaries and sets its appearance attributes.
func (snake *Snake) moveFood() {
	minRows := 1
	maxRows := snake.Rows - 2
	minColumns := 1
	maxColumns := snake.Columns - 2

	mr := maxRows - minRows
	if mr <= 0 {
		mr = 1
	}

	mc := maxColumns - minColumns
	if mc <= 0 {
		mc = 1
	}

	snake.Food.Y = rand.Intn(mr) + minRows
	snake.Food.X = rand.Intn(mc) + minColumns
	snake.Food.C = '*'
	snake.Food.Fg = interfaces.ColorRedDef
	snake.Food.Bg = interfaces.ColorNoneDef
}

// createSurface initializes and returns a 2D slice of strings representing the game surface with empty spaces based on Rows and Columns.
func (snake *Snake) createSurface() [][]string {
	var surface [][]string
	for i := 0; i < snake.Rows; i++ {
		var line []string
		for j := 0; j < snake.Columns; j++ {
			line = append(line, " ")
		}
		surface = append(surface, line)
	}
	return surface
}

// head retrieves the current position of the snake's head as a Point. Returns a default Point if the snake has no body.
func (snake *Snake) head() Point {
	if len(snake.Body) > 0 {
		var idx = len(snake.Body) - 1
		if idx >= 0 && idx < len(snake.Body) {
			return snake.Body[idx].Point
		}
	}
	return Point{}
}

// borderCollision checks if the snake's head has collided with any defined borders and returns true if a collision occurs.
func (snake *Snake) borderCollision() bool {
	k := snake.head()
	found := false
	for _, border := range snake.borders {
		if border.Contains(k) {
			found = true
			break
		}
	}
	return found
}

// foodCollision checks if the snake's head is positioned at the same coordinates as the food. Returns true if they match.
func (snake *Snake) foodCollision() bool {
	k := snake.head()
	if k.X == snake.Food.X && k.Y == snake.Food.Y {
		return true
	}
	return false
}

// snakeCollision checks if the snake's head has collided with its own body and returns true if a collision is detected.
func (snake *Snake) snakeCollision() bool {
	return snake.Contains()
}

// Advance progresses the game state by updating the snake's position and handling collisions or interactions.
func (snake *Snake) Advance() {
	//for x := 0; x <= snake.Score; x ++ {
	snake.update()
	//}
}

// update updates the snake's position, handles game logic, and checks for collisions or food consumption.
func (snake *Snake) update() {
	if !snake.GameOver {
		nHead := snake.head()
		switch snake.Direction {
		case Up:
			nHead.Y--
		case Down:
			nHead.Y++
		case Left:
			nHead.X--
		case Right:
			nHead.X++
		}

		if snake.borderCollision() || snake.snakeCollision() {
			snake.GameOver = true
		} else {
			if snake.foodCollision() {
				snake.Score++
				snake.rebuildSnake(nHead)
				snake.moveFood()
			} else {
				if len(snake.Body) > 0 {
					snake.Body = snake.Body[1:]
				}
				snake.rebuildSnake(nHead)
			}

			snake.Position.X = nHead.X
			snake.Position.Y = nHead.Y
		}
	}
}

// rebuildSnake updates the snake's tail, appends the new head, and sets the appropriate attributes for body parts.
func (snake *Snake) rebuildSnake(nHead Point) {
	if len(snake.Body) > 0 {
		snake.Body[len(snake.Body)-1].Attributes = snake.SpriteBody
	}
	snake.Body = append(snake.Body, Data{Point: nHead, Attributes: snake.SpriteHead})
	if len(snake.Body) > 1 {
		snake.Body[0].Attributes = snake.SpriteTail
	}
}

// Draw renders the snake, food, borders, and the score on the provided surface. Displays "Game Over" if the game ends.
func (snake *Snake) Draw(surface interfaces.ISurface) {
	surface.DrawColor(snake.Food.Y, snake.Food.X, snake.Food.C, snake.Food.Fg, snake.Food.Bg, interfaces.ModeNormal)

	for _, c := range snake.Body {
		surface.DrawColor(c.Y, c.X, c.C, c.Fg, c.Bg, interfaces.ModeNormal)
	}

	for _, border := range snake.borders {
		border.Draw(surface)
	}

	score := "Score: " + strconv.Itoa(snake.Score)

	surface.DrawTextColor(0, 0, score, interfaces.ColorBlueDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	if snake.GameOver {
		rows, column := surface.GetSize()
		gameOver := "Game Over"
		surface.DrawTextColor(rows/2, (column/2)-(len(gameOver)/2), gameOver, interfaces.ColorBlueDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	}
}

// SetPosition sets the snake's position by updating the X and Y coordinates of the snake's position to the given X value.
func (snake *Snake) SetPosition(x int, _ int) {
	snake.Position.X = x
	snake.Position.Y = x
}

// Contains checks if the snake's head overlaps with any part of its body, excluding the head itself. Returns true if it does.
func (snake *Snake) Contains() bool {
	for i := 0; i < len(snake.Body)-1; i++ {
		if snake.head() == snake.Body[i].Point {
			return true
		}
	}
	return false
}

// Border defines a structure representing a rectangular boundary with specific dimensions and coordinates.
// It includes the width, height, and a map tracking points and associated data within the border.
type Border struct {
	width  int
	height int
	coords map[Point]Data
}

// NewBorder creates a new Border struct with specified origin coordinates, width, and height.
// It initializes the border with points having default attributes and stores them in a map.
// The origin defines the starting point, and the width/height determine the boundary dimensions.
// Returns a pointer to the created Border instance.
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

// Contains checks whether the given point exists within the border's defined coordinates. Returns true if it exists.
func (b *Border) Contains(point Point) bool {
	_, exists := b.coords[point]
	return exists
}

// Draw renders the border onto the provided surface by iterating through its coordinates and applying color settings.
func (b *Border) Draw(s interfaces.ISurface) {
	for _, c := range b.coords {
		s.DrawColor(c.Y, c.X, c.C, c.Fg, c.Bg, interfaces.ModeNormal)
	}
}

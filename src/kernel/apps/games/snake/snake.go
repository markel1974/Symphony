package snake

import (
	"math/rand"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// direction represents an enumerated type used to define possible movement or orientation values.
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

// Point represents a coordinate with X and Y integer values.
type Point struct {
	X int
	Y int
}

// Attributes defines visual characteristics including a character, foreground color, and background color.
type Attributes struct {
	C  rune
	Fg interfaces.ColorDef
	Bg interfaces.ColorDef
}

// Data represents a structure that combines a Point and associated Attributes for rendering or processing purposes.
type Data struct {
	Point
	Attributes
}

// Snake represents a snake in the game, including its attributes, state, position, and game-related properties.
type Snake struct {
	process    interfaces.IUserProcess
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

// New creates and returns a new instance of Snake with default attributes and initial setup.
func New() *Snake {
	snake := &Snake{
		SpriteTail: Attributes{C: '░', Fg: interfaces.ColorYellowDef, Bg: interfaces.ColorNoneDef},
		SpriteBody: Attributes{C: '░', Fg: interfaces.ColorMagentaDef, Bg: interfaces.ColorNoneDef},
		SpriteHead: Attributes{C: '@', Fg: interfaces.ColorRedDef, Bg: interfaces.ColorNoneDef},
	}
	return snake
}

// Setup configures the snake's dependencies and binds event handlers for reading, timing, and painting operations.
func (snake *Snake) Setup(process interfaces.IUserProcess) {
	snake.process = process
	process.SetOnKey(snake.onKey)
	process.SetOnTimer(snake.onTimer)
	process.SetOnPaint(snake.onPaint)
}

// Start initializes the game board dimensions and creates the main game timer.
func (snake *Snake) Start() {
	w, h := snake.process.GetScreenSize()
	snake.setSize(h, w)
	snake.process.CreateTimer(0, 200, -1)
}

// onRead handles user input to change the direction of the snake or restart the game when specific keys are pressed.
func (snake *Snake) onKey(_ int, key rune) {
	switch key {
	case 'a':
		if snake.Direction != Left {
			snake.Direction = Left
		}
	case 'd':
		if snake.Direction != Right {
			snake.Direction = Right
		}
	case 'w':
		if snake.Direction != Up {
			snake.Direction = Up
		}
	case 's':
		if snake.Direction != Down {
			snake.Direction = Down
		}
	case '1':
		snake.initialize()
	}
}

// onTimer is triggered periodically by a timer, advancing the snake's state and requesting a repaint of the game surface.
func (snake *Snake) onTimer(_ int, _ int) {
	snake.Advance()
	snake.process.PaintRequest()
}

// onPaint updates the snake's appearance on the surface and adjusts dimensions if the terminal size changes.
func (snake *Snake) onPaint(surface interfaces.ISurface) {
	rows, columns := surface.GetSize()
	if snake.Rows != rows || snake.Columns != columns {
		snake.setSize(rows, columns)
	}
	snake.draw(surface)
}

// setSize updates the Snake's grid size (rows and columns) and reinitializes the game state.
func (snake *Snake) setSize(rows int, column int) {
	snake.Rows = rows
	snake.Columns = column
	snake.borders = nil

	snake.initialize()
}

// initialize resets the snake's state, including direction, speed, body, game status, and initializes borders and food.
func (snake *Snake) initialize() {
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

// moveFood relocates the food to a random position within the playable area and updates its appearance properties.
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

// createSurface initializes and returns a 2D grid representing the game's surface of empty strings with defined dimensions.
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

// head returns the current position of the snake's head as a Point. Returns an empty Point if the Body is empty.
func (snake *Snake) head() Point {
	if len(snake.Body) > 0 {
		var idx = len(snake.Body) - 1
		if idx >= 0 && idx < len(snake.Body) {
			return snake.Body[idx].Point
		}
	}
	return Point{}
}

// borderCollision checks whether the snake's head collides with any of the defined borders and returns true if a collision occurs.
func (snake *Snake) borderCollision() bool {
	k := snake.head()
	found := false
	for _, border := range snake.borders {
		if border.contains(k) {
			found = true
			break
		}
	}
	return found
}

// foodCollision checks if the snake's head position matches the food's position and returns true if they collide.
func (snake *Snake) foodCollision() bool {
	k := snake.head()
	if k.X == snake.Food.X && k.Y == snake.Food.Y {
		return true
	}
	return false
}

// snakeCollision checks if the snake's head collides with its own body and returns true if there is a collision.
func (snake *Snake) snakeCollision() bool {
	return snake.contains()
}

// Advance progresses the state of the snake in the game by handling movement, collisions, and food consumption.
func (snake *Snake) Advance() {
	//for x := 0; x <= snake.Score; x ++ {
	snake.update()
	//}
}

// update updates the snake's position on the grid and handles collisions or interactions like food consumption.
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

// rebuildSnake updates the snake's body by adding a new head and adjusting tail and body attributes accordingly.
func (snake *Snake) rebuildSnake(nHead Point) {
	if len(snake.Body) > 0 {
		snake.Body[len(snake.Body)-1].Attributes = snake.SpriteBody
	}
	snake.Body = append(snake.Body, Data{Point: nHead, Attributes: snake.SpriteHead})
	if len(snake.Body) > 1 {
		snake.Body[0].Attributes = snake.SpriteTail
	}
}

// Draw renders the snake, food, borders, score, and game-over message on the given surface.
func (snake *Snake) draw(surface interfaces.ISurface) {
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

// SetPosition updates the snake's position based on the provided x-coordinate, setting both X and Y to the same value.
func (snake *Snake) setPosition(x int, _ int) {
	snake.Position.X = x
	snake.Position.Y = x
}

// contains checks if the snake's head collides with any part of its body, returning true if a collision is detected.
func (snake *Snake) contains() bool {
	for i := 0; i < len(snake.Body)-1; i++ {
		if snake.head() == snake.Body[i].Point {
			return true
		}
	}
	return false
}

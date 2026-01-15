package invaders

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/servers/render/matrix"
)

// fgPlayText specifies the foreground color for playback text, using the cyan default color.
// bgPlayText specifies the background color for playback text, using the red default color.
// playerSpriteBottomOffset denotes the bottom offset for player sprite positioning.
// scoreX defines the x-coordinate position for the score display.
// scoreY defines the y-coordinate position for the score display.
// livesRightOffset indicates the offset for the lives display from the right edge.
// rwdSm represents the small reward value.
// rwdMd represents the medium reward value.
// rwdLg represents the large reward value.
// nonIndex represents an invalid or non-referenced index value.
const (
	fgPlayText               = interfaces.ColorCyanDef
	bgPlayText               = interfaces.ColorRedDef
	playerSpriteBottomOffset = 1

	scoreX = 10
	scoreY = 1

	livesRightOffset = 0

	rwdSm = 10
	rwdMd = 20
	rwdLg = 30

	nonIndex = -50
)

// GameState represents the current state of the game, defined as an unsigned 8-bit integer.
type GameState uint8

// MenuState represents the game state where the main menu is displayed.
// PlayState represents the game state where the gameplay occurs.
// HighScoresState represents the game state showing the high scores screen.
const (
	MenuState GameState = iota
	PlayState
	HighScoresState
)

// Invaders represents the main game entity controlling game logic, state, and rendering, managing various game components.
type Invaders struct {
	process    interfaces.IUserProcess
	highScores []*HighScore
	state      GameState
	fc         uint8
	hmi        int
	w          int
	h          int

	startX int
	startY int
	player *Player

	ufo            *matrix.Entity
	ufoTimer       *time.Time
	aliens         *Aliens
	alienMoveEvery uint8

	fragments  *Fragments
	barricade  *Barricade
	lvl        int
	freeze     int
	freezeText string

	stars *Stars
}

// New creates and initializes a new Invaders instance with the given width and height.
func New() *Invaders {
	g := &Invaders{
		stars: NewStars(),
	}
	g.init()
	return g
}

// Setup initializes the Invaders game by setting up the user process with the necessary event handlers.
func (g *Invaders) Setup(process interfaces.IUserProcess) {
	g.process = process
	g.process.SetOnKey(g.handleKey)
	g.process.SetOnTimer(g.handleTimer)
	g.process.SetOnPaint(g.handlePaint)
}

// Start initializes the game by setting the menu state and creating a repeating timer for game updates.
func (g *Invaders) Start() {
	g.w, g.h = g.process.GetScreenSize()
	g.setMenuState()
	g.process.CreateTimer(0, 100, -1)
}

// init initializes the Invaders game's core components, including high scores, barricade, fragments, and aliens.
func (g *Invaders) init() {
	g.highScores = make([]*HighScore, 0)
	g.fc = 1

	g.barricade = NewBarricade()
	g.fragments = NewFragments()
	g.aliens = NewAliens()
}

// SetSize updates the width and height of the game and reinitializes its components.
func (g *Invaders) SetSize(w int, h int) {
	g.w = w
	g.h = h
	g.init()
}

// GetSize returns the current width and height of the Invaders game field.
func (g *Invaders) GetSize() (int, int) {
	return g.w, g.h
}

// drawMenu renders the game's menu interface, including title, menu items, and animated stars, onto the provided surface.
func (g *Invaders) drawMenu(surface interfaces.ISurface) {
	x := g.w/2 - logoLineLength/2
	y := logoY

	for _, line := range logoLines {
		drawColor(surface, x, y, fgMenu, bgMenu, line)
		y++
	}

	length := 0
	i := 0
	for _, v := range menuItems {
		length += len(v)
		if i+1 != len(menuItems) {
			length += menuPad
		}
		i++
	}

	x = g.w/2 - length/2
	y += logoHeight + 5
	for i := FirstMenuItem; i < NumMenuItems; i++ {
		v := menuItems[i]
		if i == g.hmi {
			drawColor(surface, x, y, fgMenuHighlight, bgMenuHighlight, v)
		} else {
			drawColor(surface, x, y, fgMenu, bgMenu, v)
		}
		x += len(v) + menuPad
	}

	g.stars.Draw(surface)
}

// drawFlash displays a flashing message at the center of the surface when the freeze timer is active, decrementing it each frame.
func (g *Invaders) drawFlash(surface interfaces.ISurface) {
	if g.freeze > 0 {
		drawColor(surface, g.w/2-len(g.freezeText)/2, g.h/2, fgPlayText, bgPlayText, g.freezeText)
		g.freeze--
	}
}

// drawPlay renders the main play area of the game, including stars, barricades, UFO, aliens, player, score, lives, and fragments.
func (g *Invaders) drawPlay(surface interfaces.ISurface) {
	g.stars.Draw(surface)
	g.barricade.Draw(surface)
	if g.ufo != nil {
		g.ufo.Draw(surface)
		g.ufo.Next()
	}
	g.aliens.Draw(surface)
	g.player.Draw(surface)
	drawColor(surface, scoreX, scoreY, fgPlayText, bgPlayText, strconv.Itoa(g.player.score))
	drawColor(surface, g.w-livesRightOffset-len(g.player.GetLivesStr()), scoreY, fgPlayText, bgPlayText, g.player.GetLivesStr())
	g.fragments.Draw(surface)
	g.drawFlash(surface)
}

// draw renders the appropriate game state UI elements on the provided surface.
func (g *Invaders) draw(surface interfaces.ISurface) {
	switch g.state {
	case MenuState:
		g.drawMenu(surface)
	case PlayState:
		g.drawPlay(surface)
	case HighScoresState:
		//g.DrawHighscores()
	}
}

// update updates the game state based on the current state (MenuState, PlayState, HighScoresState) and increments the frame count.
func (g *Invaders) update() {
	g.fc++
	switch g.state {
	case MenuState:
		g.updateMenu()
	case PlayState:
		g.updatePlay()
	case HighScoresState:
		//g.UpdateHighscores()
	}
}

// updateMenu updates the menu by refreshing the stars' position based on the current game dimensions and frame count.
func (g *Invaders) updateMenu() {
	g.stars.Update(g.w, g.h, g.fc)
}

// updatePlay updates the game state during the play phase, handling collisions, movements, scoring, and game logic.
func (g *Invaders) updatePlay() {
	g.stars.Update(g.w, g.h, g.fc)

	for b := range g.aliens.alienBullets.data {
		if g.aliens.alienBullets.data[b] != nil {
			g.aliens.alienBullets.data[b].AddToY(alienBulletSpeed)

			if int(g.aliens.alienBullets.data[b].GetY()) >= g.h {
				g.aliens.alienBullets.data[b] = nil
			} else {
				x, y := g.aliens.alienBullets.data[b].GetX(), g.aliens.alienBullets.data[b].GetY()

				if g.player.HasCollision(g.aliens.alienBullets.data[b].Entity) {
					if !g.player.DecLives() {
						g.gameOver()
						return
					}

					g.wipeBullets()

					g.freeze = 20
					g.freezeText = fmt.Sprintf("Level %d", g.lvl)
					return
				} else if g.barricade.Unset(int(x), int(y)) {
					g.aliens.alienBullets.data[b] = nil
				}
			}
		}
	}

	if g.player.bullet != nil {
		g.player.bullet.AddToY(g.player.bullet.Vy)
		if g.player.bullet.GetY() < 0 {
			g.player.bullet = nil
		} else {
			found, reward := g.aliens.DoBulletCollision(g.player.bullet.Entity)
			if !found {
				if g.ufo != nil {
					if g.ufo.HasCollision(g.player.bullet.Entity) {
						reward = ufoReward
						g.ufo = nil
						found = true
					}
				}
			}

			if found {
				g.explode(int(g.player.bullet.GetX()), int(g.player.bullet.GetY()))
				g.player.score += reward
				g.player.WipeBullet()
			} else {
				if g.barricade.Unset(int(g.player.bullet.GetX()), int(g.player.bullet.GetY())) {
					g.player.WipeBullet()
				}
			}
		}
	}

	if g.ufo != nil && g.fc%ufoMoveEvery == 0 {
		g.ufo.AddToX(1)
		if int(g.ufo.GetX()) > g.w {
			g.ufo = nil
			t := time.Now()
			t = t.Add(time.Duration(rand.Intn(20)+15) * time.Second)
			g.ufoTimer = &t
		}
	}

	if g.fc%g.alienMoveEvery == 0 {
		g.aliens.DoMove(g.w, g.h)

		if g.aliens.minY >= int(g.player.sprite.GetY()-g.player.sprite.GetHeight()) {
			g.gameOver()
			return
		}

		if g.aliens.count <= 0 && g.ufo == nil {
			g.lvl += 1
			if g.alienMoveEvery != 2 {
				g.alienMoveEvery--
			}
			g.beginNextLevel()
			g.wipeBullets()
		}
	}

	g.fragments.Update()

	if g.ufoTimer != nil {
		if (*g.ufoTimer).Before(time.Now()) {
			g.ufo = newUfo()
			g.ufoTimer = nil
		}
	}
}

// handleKey routes key input to appropriate handlers based on the current game state.
func (g *Invaders) handleKey(_ int, k rune) {
	switch g.state {
	case MenuState:
		g.handleKeyMenu(k)
	case PlayState:
		g.handleKeyPlay(k)
	case HighScoresState:
		//g.HandleKeyHighscores(k)
	}
}

// handlePaint adjusts the game dimensions and redraws the surface based on the current game state.
func (g *Invaders) handlePaint(surface interfaces.ISurface) {
	rows, columns := surface.GetScreenSize()
	w, h := g.GetSize()
	if h != rows || w != columns {
		g.SetSize(columns, rows)
		g.setMenuState()
	}
	g.draw(surface)
}

// handleTimer is triggered by a timer event, updates game state, and requests a repaint of the game surface.
func (g *Invaders) handleTimer(_ int, _ int) {
	g.update()
	g.process.PaintRequest()
}

// handleKeyMenu processes key inputs for navigating or selecting options in the menu state of the game.
func (g *Invaders) handleKeyMenu(k rune) {
	switch k {
	case 'a':
		g.hmi = (g.hmi - 1 + NumMenuItems) % NumMenuItems
	case 'd':
		g.hmi = (g.hmi + 1) % NumMenuItems
	case ' ':
		switch g.hmi {
		case HighScores:
			//g.GoHighscores()
		case Howto:
			//g.GoHowto()
		case Play:
			g.SetPlayState()
		default:
			//Nothing to do
		}
	}
}

// handleKeyPlay processes key inputs during the game's play state and executes corresponding player actions.
func (g *Invaders) handleKeyPlay(k rune) {
	switch k {
	case 'd':
		g.player.MoveRight(g.w)
	case 'a':
		g.player.MoveLeft(0)
	case ' ':
		g.player.NewBullet()
	}
}

// setMenuState sets the game state to the menu and initializes the highlighted menu item to the first option.
func (g *Invaders) setMenuState() {
	g.state = MenuState
	g.hmi = FirstMenuItem
}

// SetPlayState transitions the game to the play state, initializes gameplay elements, and starts the next level.
func (g *Invaders) SetPlayState() {
	g.state = PlayState
	g.initPlay()
	g.beginNextLevel()
}

// beginNextLevel initializes the next level by creating a new set of aliens at the starting coordinates.
func (g *Invaders) beginNextLevel() {
	g.aliens.Create(alienStartX, alienStartY, 0)
}

// explode triggers the addition of a new fragment at the specified coordinates (x, y).
func (g *Invaders) explode(x, y int) {
	g.fragments.AddFragment(x, y)
}

// wipeBullets clears all bullets from the player and resets alien bullets based on the number of aliens present.
func (g *Invaders) wipeBullets() {
	g.player.WipeBullet()
	g.aliens.alienBullets.Setup(g.aliens.columns * g.aliens.rows / 10)
}

// initPlay initializes the game state for the play mode, setting up player, aliens, barricades, fragments, and UFO timer.
func (g *Invaders) initPlay() {
	g.fc = 1
	g.lvl = 1
	g.freeze = 0
	g.freezeText = ""

	g.player = NewPlayer(0, 0)

	g.startX = g.w/2 - int(g.player.sprite.GetWidth()/2)
	g.startY = g.h - playerSpriteBottomOffset - int(g.player.sprite.GetHeight())

	g.player.sprite.MoveTo(float64(g.startX), float64(g.startY))

	g.alienMoveEvery = uint8(alienMoveEvery)

	g.aliens.Setup(g.w, g.h)

	g.fragments.Setup()

	g.barricade.Setup(g.w, g.h)

	t := time.Now()
	t = t.Add(time.Duration(rand.Intn(20)+15) * time.Second)
	g.ufoTimer = &t
}

// gameOver resets the game to the menu state, displays "GAME OVER" text, and freezes the game temporarily.
func (g *Invaders) gameOver() {
	g.freeze = 200
	g.freezeText = "GAME OVER"
	g.setMenuState()
}

package invaders

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/render/matrix"
)

// fgPlayText defines the foreground color for play text.
// bgPlayText defines the background color for play text.
// playerSpriteBottomOffset specifies the offset for the bottom of the player sprite.
// scoreX denotes the X-coordinate position for the score display.
// scoreY denotes the Y-coordinate position for the score display.
// livesRightOffset specifies the right offset for the lives display.
// rwdSm defines the reward for a small achievement.
// rwdMd defines the reward for a medium achievement.
// rwdLg defines the reward for a large achievement.
// nonIndex represents a non-existent or invalid index value.
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

// GameState represents the state of the game, typically using enumerated constants such as menu, playing, or paused.
type GameState uint8

// MenuState represents the game state where the main menu is displayed.
// PlayState represents the game state where the game is actively being played.
// HighScoresState represents the game state where the high scores are displayed.
const (
	MenuState GameState = iota
	PlayState
	HighScoresState
)

// Invaders represents the core structure for the game, managing entities, states, and interactions.
// It maintains player data, alien configurations, UFOs, stars, barricades, and various gameplay states.
type Invaders struct {
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

// NewGame initializes a new game instance with the given width and height, setting up game components and state.
func NewGame(w int, h int) *Invaders {
	g := &Invaders{
		h:     h,
		w:     w,
		stars: NewStars(),
	}
	g.init()
	return g
}

// init initializes the Invaders instance by setting up high scores, frame count, barricades, fragments, and aliens.// init initializes the Invaders game state, setting up highScores, barricades, fragments, and alien entities.
func (g *Invaders) init() {
	g.highScores = make([]*HighScore, 0)
	g.fc = 1

	g.barricade = NewBarricade()
	g.fragments = NewFragments()
	g.aliens = NewAliens()
}

// SetSize sets the width and height of the game area and reinitializes the game state.
func (g *Invaders) SetSize(w int, h int) {
	g.w = w
	g.h = h
	g.init()
}

// GetSize returns the width and height of the Invaders as two integer values.// GetSize returns the width and height of the Invaders game field as two integers.
func (g *Invaders) GetSize() (int, int) {
	return g.w, g.h
}

// drawMenu renders the game menu on the provided surface, including logo, menu items, and starfield background.
// It calculates positions and highlights the currently selected menu item.
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

// drawFlash renders a visual effect or text on the surface during a specific freeze state, then decreases the freeze counter.
func (g *Invaders) drawFlash(surface interfaces.ISurface) {
	if g.freeze > 0 {
		drawColor(surface, g.w/2-len(g.freezeText)/2, g.h/2, fgPlayText, bgPlayText, g.freezeText)
		g.freeze--
	}
}

// drawPlay renders all game components on the provided surface, including stars, barricade, UFO, aliens, player, and fragments.
// It also displays the player's score and remaining lives, and manages visual effects such as screen flashes.
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

// Draw renders game elements to the provided surface based on the current game state (MenuState, PlayState, HighScoresState).
func (g *Invaders) Draw(surface interfaces.ISurface) {
	switch g.state {
	case MenuState:
		g.drawMenu(surface)
	case PlayState:
		g.drawPlay(surface)
	case HighScoresState:
		//g.DrawHighscores()
	}
}

// Update evaluates and acts upon the current game state by incrementing the frame count and updating specific state logic.
func (g *Invaders) Update() {
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

// updateMenu updates the starfield animation for the menu screen based on the current window size and frame count.
func (g *Invaders) updateMenu() {
	g.stars.Update(g.w, g.h, g.fc)
}

// updatePlay controls the gameplay logic, including updates to enemies, player, projectiles, and level progression.
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

// HandleKey processes a key input and delegates handling based on the current game state.
func (g *Invaders) HandleKey(k rune) {
	switch g.state {
	case MenuState:
		g.handleKeyMenu(k)
	case PlayState:
		g.handleKeyPlay(k)
	case HighScoresState:
		//g.HandleKeyHighscores(k)
	}
}

// handleKeyMenu processes key input to navigate and select items in the menu of the Invaders game.
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

// handleKeyPlay processes the player's key input during the game-play state and executes corresponding actions.
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

// SetMenuState sets the game state to MenuState and initializes the highlighted menu item to the first menu item.
func (g *Invaders) SetMenuState() {
	g.state = MenuState
	g.hmi = FirstMenuItem
}

// SetPlayState transitions the game to the play state, initializes the game settings, and starts the first level.
func (g *Invaders) SetPlayState() {
	g.state = PlayState
	g.initPlay()
	g.beginNextLevel()
}

// beginNextLevel initializes the next level by creating a new set of aliens at the starting coordinates.
func (g *Invaders) beginNextLevel() {
	g.aliens.Create(alienStartX, alienStartY, 0)
}

// explode triggers an explosion effect at the specified (x, y) coordinates by adding a new fragment to the fragments system.
func (g *Invaders) explode(x, y int) {
	g.fragments.AddFragment(x, y)
}

// wipeBullets resets the player's bullet and reinitializes alien bullets based on the number of alien columns and rows.
func (g *Invaders) wipeBullets() {
	g.player.WipeBullet()
	g.aliens.alienBullets.Setup(g.aliens.columns * g.aliens.rows / 10)
}

// initPlay initializes the game state for a new play session, setting up player, aliens, barricades, and timers.
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

// gameOver sets the game state to "GAME OVER", displays the freeze text, and transitions to the menu state.
func (g *Invaders) gameOver() {
	g.freeze = 200
	g.freezeText = "GAME OVER"
	g.SetMenuState()
}

package invaders

import (
	"strings"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/render/matrix"
)

// initLives represents the initial number of lives a player starts with in the game.
// livesSprite defines the string representation of a single life in the UI.
// playerMoveSpeed determines the speed at which the player can move.
// playerBulletSpeed sets the speed of bullets fired by the player.
const (
	initLives         = 5
	livesSprite       = `| `
	playerMoveSpeed   = 2
	playerBulletSpeed = -1
)

// playerSprite defines the visual representation of the player as a collection of string-based sprites.
var playerSprite = []string{
	`
  YY
GcGGGG
GGGGGG
`,
}

// Player represents a game character with a visual entity, a bullet, score, lives, and a string representation of lives.
type Player struct {
	sprite   *matrix.Entity
	bullet   *Bullet
	score    int
	lives    int
	livesStr string
}

// NewPlayer creates and returns a new Player instance initialized at the specified position (x, y).
func NewPlayer(x int, y int) *Player {
	player := &Player{
		sprite: matrix.NewEntity(-1, x, y, 1, 1.0, 1.0, playerSprite, 0),
		score:  0,
		lives:  initLives,
		bullet: nil,
	}
	player.syncLives()
	return player
}

// HasCollision checks if the player's sprite collides with the specified entity based on their rectangular bounds.
func (player *Player) HasCollision(e *matrix.Entity) bool {
	return player.sprite.HasCollision(e)
}

// Draw renders the player's sprite and optionally their bullet on the given surface, advancing their frames if present.
func (player *Player) Draw(surface interfaces.ISurface) {
	player.sprite.Draw(surface)
	player.sprite.Next()
	if player.bullet != nil {
		player.bullet.Draw(surface)
		player.bullet.Next()
	}
}

// NewBullet creates a new bullet for the player if one does not already exist, positioned at the player's sprite midpoint.
func (player *Player) NewBullet() {
	if player.bullet == nil {
		player.bullet = NewBullet(int(player.sprite.GetX()+player.sprite.GetWidth()/2), int(player.sprite.GetY()), playerBulletSpeed)
	}
}

// WipeBullet clears the player's currently assigned bullet by setting it to nil.
func (player *Player) WipeBullet() {
	player.bullet = nil
}

// MoveLeft moves the player's sprite to the left by decreasing its x-coordinate, bounded by the specified minimum width.
func (player *Player) MoveLeft(minW int) {
	player.sprite.AddToX(-playerMoveSpeed)
	if int(player.sprite.GetX()) < minW {
		player.sprite.MoveToX(float64(minW))
	}
}

// MoveRight moves the player to the right by a predefined speed and ensures the movement stays within the specified maximum width.
func (player *Player) MoveRight(maxW int) {
	player.sprite.AddToX(playerMoveSpeed)
	if int(player.sprite.GetX()+player.sprite.GetWidth()) > maxW {
		player.sprite.MoveToX(float64(maxW) - player.sprite.GetWidth())
	}
}

// DecLives decreases the player's life count by 1 and updates the lives string representation. Returns false if lives are 0 or less.
func (player *Player) DecLives() bool {
	player.lives -= 1
	player.syncLives()
	if player.lives <= 0 {
		return false
	}
	return true
}

// GetLivesStr returns the string representation of the player's remaining lives.
func (player *Player) GetLivesStr() string {
	return player.livesStr
}

// syncLives updates the player's lives visual representation by repeating the livesSprite string based on the current lives count.
func (player *Player) syncLives() {
	player.livesStr = strings.Repeat(livesSprite, player.lives)
}

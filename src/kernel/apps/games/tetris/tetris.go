package tetris

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// levelMax defines the maximum level a player can reach in the game.
// scoreMax specifies the maximum score a player can achieve in the game.
const (
	levelMax = 20
	scoreMax = 999999
)

// Tetris represents the main structure for managing the Tetris game, including the game board, score, level, and gameplay state.
type Tetris struct {
	process         interfaces.IUserProcess
	board           *Board
	currentMino     *Tetromino
	nextMino        *Tetromino
	score           int
	level           int
	initLevel       int
	deleteLines     int
	gameOver        bool
	backgroundLines []string

	animationGameOverCount int
	animationDeleteCount   int
	animationDelete        []int
}

// New initializes and returns a new instance of Tetris.
func New() *Tetris {
	t := &Tetris{}
	return t
}

// Setup initializes the Tetris instance with the given user process and sets up timer, read, and paint event handlers.
func (t *Tetris) Setup(process interfaces.IUserProcess) {
	t.process = process
	t.process.SetOnTimer(t.onTimer)
	t.process.SetOnKey(t.onKey)
	t.process.SetOnPaint(t.onPaint)
}

// Start initializes the game state and sets up the game screen and timer.
func (t *Tetris) Start() {
	t.initialize()
	t.process.CreateTimer(0, 300, -1)
}

// onKey processes user input represented by the key parameter and performs corresponding game actions.
func (t *Tetris) onKey(_ int, key rune) {
	switch key {
	case 'a':
		t.MoveLeft()
	case 'd':
		t.MoveRight()
	case 'w':
		t.RotateRight()
	case 's':
		t.MoveDown()
	case ' ':
		t.Drop()
	case '1':
		//r := cmd.GetRootContext()
		//w, h := r.GetScreenSize()
		t.initialize()
	}
}

// onTimer is triggered by a timer and applies gravity to the game state and requests a repaint of the game field.
func (t *Tetris) onTimer(_ int, _ int) {
	t.ApplyGravity()
	t.process.PaintRequest()
}

// onPaint is responsible for rendering the current state of the Tetris game onto the provided drawing surface.
func (t *Tetris) onPaint(surface interfaces.ISurface) {
	t.Draw(surface)
}

// initialize sets up the initial state of the Tetris game, configuring the board, score, level, and other game variables.
func (t *Tetris) initialize() {
	w := 10
	h := 18
	//h, w := t.process.GetScreenSize()
	rand.Seed(time.Now().UnixNano())
	t.board = NewBoard(w, h)
	t.level = t.initLevel
	t.score = 0
	t.deleteLines = 0
	t.currentMino = nil
	t.nextMino = nil
	t.gameOver = false
	t.animationGameOverCount = 0
	t.animationDeleteCount = 0
	t.animationDelete = nil
	t.backgroundLines = background

	t.initMino()
}

// GetSize returns the width and height of the Tetris board as two integer values.
//func (t *Tetris) getSize() (int, int) {
//	return t.board.w, t.board.h
//}

// initMino initializes the current and next tetromino for the Tetris game and ensures two tetrominoes are queued.
func (t *Tetris) initMino() {
	t.currentMino = nil
	t.nextMino = nil
	t.pushMino()
	t.pushMino()
}

// deleteCheck evaluates the game board for completed lines, removes them, updates the score, and adjusts the game level accordingly.
func (t *Tetris) deleteCheck() {
	if !t.board.hasFullLine() {
		return
	}

	lines := t.board.fullLines()

	t.animationDeleteCount = 4
	t.animationDelete = lines

	for _, line := range lines {
		t.board.deleteLine(line)
	}
	t.deleteLines += len(lines)
	switch len(lines) {
	case 1:
		t.addScore(40 * (t.level + 1))
	case 2:
		t.addScore(100 * (t.level + 1))
	case 3:
		t.addScore(300 * (t.level + 1))
	case 4:
		t.addScore(1200 * (t.level + 1))
	}
	t.levelUpdate()
}

// levelUpdate updates the game level based on the number of lines deleted, capped at the maximum allowed level.
func (t *Tetris) levelUpdate() {
	if t.level == levelMax {
		return
	}

	targetLevel := t.deleteLines / 10
	if t.level < targetLevel {
		t.level = targetLevel
	}
}

// addScore increases the current score by the specified amount and caps it at a maximum predefined value.
func (t *Tetris) addScore(add int) {
	t.score += add
	if t.score > scoreMax {
		t.score = scoreMax
	}
}

// pushMino manages the transition of the current tetromino to the board and initializes the next tetromino for the game.
func (t *Tetris) pushMino() {
	t.deleteCheck()

	t.currentMino = t.nextMino
	if t.currentMino != nil {
		t.currentMino.x, t.currentMino.y = defaultTetrominoX, defaultTetrominoY
		if t.currentMino.conflicts(t.board) {
			ranking := NewRanking()
			ranking.insertScore(t.score)
			//ranking.save()
			t.gameOver = true
			return
		}
	}
	t.nextMino = NewMino()
}

// ApplyGravity moves the current Tetris piece down by one row, simulating gravity, unless the game is over.
func (t *Tetris) ApplyGravity() {
	if t.gameOver {
		return
	}
	t.MoveDown()
}

// Drop moves the current Tetris piece directly to the lowest possible position on the board.
func (t *Tetris) Drop() {
	if t.gameOver {
		return
	}
	t.addScore(t.currentMino.putBottom(t.board))
	t.board.setCells(t.currentMino.cells())
	t.pushMino()
}

// MoveDown attempts to move the current tetromino one step down; places it if it collides or pushes a new tetromino if needed.
func (t *Tetris) MoveDown() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.forceMoveDown()

	if dstMino.conflicts(t.board) {
		t.board.setCells(t.currentMino.cells())
		t.pushMino()
	} else {
		t.currentMino.forceMoveDown()
	}
}

// MoveLeft moves the current tetromino one unit to the left if it does not cause a conflict with the board or boundaries.
func (t *Tetris) MoveLeft() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.x--
	if !dstMino.conflicts(t.board) {
		t.currentMino.x--
	}
}

// MoveRight moves the current Tetrimino one step to the right if there are no conflicts and the game is not over.
func (t *Tetris) MoveRight() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.x++
	if !dstMino.conflicts(t.board) {
		t.currentMino.x++
	}
}

// RotateRight attempts to rotate the current Tetris piece clockwise if the game is not over and no conflicts would occur.
func (t *Tetris) RotateRight() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.forceRotateRight()
	if !dstMino.conflicts(t.board) {
		t.currentMino.forceRotateRight()
	}
}

// RotateLeft attempts to rotate the current mino to the left. Rotation is applied only if it does not cause a conflict.
func (t *Tetris) RotateLeft() {
	if t.gameOver {
		return
	}
	dstMino := *t.currentMino
	dstMino.forceRotateLeft()
	if !dstMino.conflicts(t.board) {
		t.currentMino.forceRotateLeft()
	}
}

// Draw renders the entire Tetris game state onto the provided surface, including background, board, pieces, and animations.
func (t *Tetris) Draw(surface interfaces.ISurface) {
	t.drawBackGround(surface, 0, 0)
	t.drawBoard(surface, boardXOffset, boardYOffset)
	t.drawMino(surface, t.nextMino, nextMinoXOffset-t.nextMino.x, nextMinoYOffset-t.nextMino.y)
	t.drawTexts(surface)
	t.drawDropMarker(surface)
	t.drawMino(surface, t.currentMino, boardXOffset, boardYOffset)
	t.drawAnimationDelete(surface)
	t.drawGameOver(surface)
}

// drawTexts renders text information such as score, level, and lines on the provided surface with specified styles.
func (t *Tetris) drawTexts(surface interfaces.ISurface) {
	surface.DrawTextColor(9, 32, "SCORE", interfaces.ColorWhiteDef, interfaces.ColorBlueDef, interfaces.ModeNormal)
	surface.DrawTextColor(10, 32, fmt.Sprintf("%7d", t.score), interfaces.ColorBlackDef, interfaces.ColorWhiteDef, interfaces.ModeNormal)
	surface.DrawTextColor(13, 32, "LEVEL", interfaces.ColorWhiteDef, interfaces.ColorBlueDef, interfaces.ModeNormal)
	surface.DrawTextColor(14, 32, fmt.Sprintf("%5d", t.level), interfaces.ColorBlackDef, interfaces.ColorWhiteDef, interfaces.ModeNormal)
	surface.DrawTextColor(16, 32, "LINES", interfaces.ColorWhiteDef, interfaces.ColorBlueDef, interfaces.ModeNormal)
	surface.DrawTextColor(17, 32, fmt.Sprintf("%5d", t.deleteLines), interfaces.ColorBlackDef, interfaces.ColorWhiteDef, interfaces.ModeNormal)
}

// drawDropMarker renders a visual indicator on the board showing where the current Tetris piece will land.
func (t *Tetris) drawDropMarker(surface interfaces.ISurface) {
	marker := *t.currentMino
	marker.putBottom(t.board)

	for y, line := range marker.lines() {
		for x, char := range line {
			if t.board.isOnBoard(x+marker.x, y+marker.y) && colorByChar(char) != blankColor &&
				colorByChar(char) != interfaces.ColorNoneDef {
				t.drawCell(surface, x+marker.x+boardXOffset, y+marker.y+boardYOffset, colorByChar('K'))
			}
		}
	}
}

// drawMino renders the given Tetromino onto the specified surface with applied x and y offsets for positioning.
func (t *Tetris) drawMino(surface interfaces.ISurface, mino *Tetromino, xOffset, yOffset int) {
	for y, line := range mino.lines() {
		for x, char := range line {
			if t.board.isOnBoard(x+mino.x, y+mino.y) {
				color := colorByChar(char)
				t.drawCell(surface, x+mino.x+xOffset, y+mino.y+yOffset, color)
			}
		}
	}
}

// drawBoard renders the Tetris game board onto the specified surface at the given top-left coordinates.
func (t *Tetris) drawBoard(surface interfaces.ISurface, left int, top int) {
	for j := 0; j < t.board.h; j++ {
		for i := 0; i < t.board.w; i++ {
			t.drawCell(surface, left+i, top+j, t.board.colors[i][j])
		}
	}
}

// drawCell renders a single cell on the given surface at the specified coordinates using the provided color definition.
func (t *Tetris) drawCell(surface interfaces.ISurface, x, y int, color interfaces.ColorDef) {
	if color != interfaces.ColorNoneDef && color != blankColor {
		if color == colorByChar('K') {
			surface.DrawColor(y, 2*x-1, '▓', color, interfaces.ColorWhiteDef, interfaces.ModeNormal)
			surface.DrawColor(y, 2*x, ' ', color, interfaces.ColorWhiteDef, interfaces.ModeNormal)
		} else {
			var bg interfaces.ColorDef
			switch color {
			case interfaces.ColorRedDef:
				bg = interfaces.ColorBrightRedDef
			case interfaces.ColorGreenDef:
				bg = interfaces.ColorBrightGreenDef
			case interfaces.ColorYellowDef:
				bg = interfaces.ColorBrightYellowDef
			case interfaces.ColorBlueDef:
				bg = interfaces.ColorBrightBlueDef
			case interfaces.ColorMagentaDef:
				bg = interfaces.ColorBrightMagentaDef
			case interfaces.ColorCyanDef:
				bg = interfaces.ColorBrightCyanDef
			case interfaces.ColorWhiteDef:
				bg = interfaces.ColorBrightWhiteDef
			default:
				bg = color
			}
			surface.DrawColor(y, 2*x-1, '▓', color, bg, interfaces.ModeNormal)
			surface.DrawColor(y, 2*x, ' ', color, bg, interfaces.ModeNormal)
		}
	}
}

// drawBackGround renders the background on the given surface using the provided coordinates as the starting position.
func (t *Tetris) drawBackGround(surface interfaces.ISurface, left int, top int) {
	for y, line := range t.backgroundLines {
		for x, char := range line {
			t.drawBack(surface, left+x, top+y, colorByChar(char))
		}
	}
}

// drawBack renders the background of a Tetris block at the specified coordinates with the given color on the surface.
func (t *Tetris) drawBack(surface interfaces.ISurface, x, y int, color interfaces.ColorDef) {
	surface.DrawColor(y, 2*x-1, ' ', interfaces.ColorNoneDef, color, interfaces.ModeNormal)
	surface.DrawColor(y, 2*x, ' ', interfaces.ColorNoneDef, color, interfaces.ModeNormal)
}

// drawAnimationDelete animates the deletion of lines by alternating their color until the animation completes.
func (t *Tetris) drawAnimationDelete(surface interfaces.ISurface) {
	if t.animationDeleteCount > 0 {
		for _, line := range t.animationDelete {
			color := interfaces.ColorCyanDef
			if t.animationDeleteCount%2 == 0 {
				color = interfaces.ColorMagentaDef
			}
			t.colorizeLine(surface, line, color)
		}
		t.animationDeleteCount--
	}
}

// drawGameOver renders the game over sequence onto the provided surface, including animations and ranking display.
func (t *Tetris) drawGameOver(surface interfaces.ISurface) {
	if t.gameOver {
		if t.animationGameOverCount < t.board.h {
			for y := t.board.h - 1; y > t.board.h-1-t.animationGameOverCount; y -= 1 {
				t.colorizeLine(surface, y, interfaces.ColorBlackDef)
			}
			t.animationGameOverCount++
			return
		}

		for j := 0; j < t.board.h; j++ {
			t.colorizeLine(surface, j, interfaces.ColorBlackDef)
		}
		surface.DrawTextColor(4, 10, "GAME OVER", interfaces.ColorWhiteDef, interfaces.ColorBlackDef, interfaces.ModeNormal)

		ranking := NewRanking()
		for idx, sc := range ranking.scores {
			fg := availableColors[rand.Intn(len(availableColors))]
			surface.DrawTextColor(8+idx, 9, fmt.Sprintf("%2d: %6d", idx+1, sc), fg, interfaces.ColorBlackDef, interfaces.ModeNormal)
		}
	}
}

// colorizeLine applies the specified color to a horizontal line on the game board using the provided surface.
func (t *Tetris) colorizeLine(surface interfaces.ISurface, line int, color interfaces.ColorDef) {
	for i := 0; i < t.board.w; i++ {
		t.drawBack(surface, i+boardXOffset, line+boardYOffset, color)
	}
}

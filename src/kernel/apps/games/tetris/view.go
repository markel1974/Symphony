package tetris

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strings"
)

// srcBackground defines the initial visual layout of the board in string format.
// boardXOffset specifies the horizontal offset for the game board.
// boardYOffset specifies the vertical offset for the game board.
// nextMinoXOffset specifies the horizontal offset for displaying the next piece.
// nextMinoYOffset specifies the vertical offset for displaying the next piece.
// blankColor defines the color to represent blank areas on the board.
const (
	srcBackground = `
		WWWWWWWWWWWW WWWWWW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WkkkkW
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW BBBBBB
		WkkkkkkkkkkW WWWWWW
		WkkkkkkkkkkW
		WkkkkkkkkkkW
		WWWWWWWWWWWW
	`

	boardXOffset    = 3
	boardYOffset    = 2
	nextMinoXOffset = 16
	nextMinoYOffset = 2
	blankColor      = interfaces.ColorBlackDef
)

// background holds a slice of strings split by newline from the srcBackground constant.
//
// colorMapping maps rune values to their corresponding ColorDef values from the interfaces package.
var (
	background = strings.Split(srcBackground, "\n")

	colorMapping = map[rune]interfaces.ColorDef{
		'k': interfaces.ColorBlackDef,
		'K': interfaces.ColorBrightBlackDef,
		'r': interfaces.ColorRedDef,
		'R': interfaces.ColorBrightRedDef,
		'g': interfaces.ColorGreenDef,
		'G': interfaces.ColorBrightGreenDef,
		'y': interfaces.ColorYellowDef,
		'Y': interfaces.ColorBrightYellowDef,
		'b': interfaces.ColorBlueDef,
		'B': interfaces.ColorBrightBlueDef,
		'm': interfaces.ColorMagentaDef,
		'M': interfaces.ColorBrightMagentaDef,
		'c': interfaces.ColorCyanDef,
		'C': interfaces.ColorBrightCyanDef,
		'w': interfaces.ColorWhiteDef,
		'W': interfaces.ColorWhiteDef,
	}
)

// availableColors defines a list of predefined color definitions used in the application.
var (
	availableColors = []interfaces.ColorDef{
		interfaces.ColorRedDef,
		interfaces.ColorGreenDef,
		interfaces.ColorYellowDef,
		interfaces.ColorBlueDef,
		interfaces.ColorMagentaDef,
		interfaces.ColorCyanDef,
		interfaces.ColorWhiteDef,
		interfaces.ColorBrightRedDef,
		interfaces.ColorBrightGreenDef,
		interfaces.ColorBrightYellowDef,
		interfaces.ColorBrightBlueDef,
		interfaces.ColorBrightMagentaDef,
		interfaces.ColorBrightCyanDef,
		interfaces.ColorBrightWhiteDef,
		interfaces.ColorBrightBlackDef,
	}
)

// colorByChar maps a given rune to a corresponding color definition from the predefined colorMapping.
// Returns the matching interfaces.ColorDef value based on the input rune.
func colorByChar(ch rune) interfaces.ColorDef {
	return colorMapping[ch]
}

// _ returns the rune corresponding to the provided ColorDef from the colorMapping. If no match is found, it returns '.'.
func _(color interfaces.ColorDef) rune {
	for ch, currentColor := range colorMapping {
		if currentColor == color {
			return ch
		}
	}
	return '.'
}

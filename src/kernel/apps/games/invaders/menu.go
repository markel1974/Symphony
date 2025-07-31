package invaders

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"strings"
)

// menuPad defines the padding value for the menu layout.
// fgMenu represents the foreground color for the menu, set to red.
// bgMenu represents the background color for the menu, set to black.
// fgMenuHighlight represents the foreground color for highlighted menu items, set to white.
// bgMenuHighlight represents the background color for highlighted menu items, set to cyan.
// logoY defines the vertical positioning of the logo in the UI.
// logo provides an ASCII art logo for the application or interface.
const (
	menuPad         = 10
	fgMenu          = interfaces.ColorRedDef
	bgMenu          = interfaces.ColorBlackDef
	fgMenuHighlight = interfaces.ColorWhiteDef
	bgMenuHighlight = interfaces.ColorCyanDef
	logoY           = 2
	logo            = `
                          _                     _
                         (_)                   | |              
 ___ _ __   __ _  ___ ___ _ _ ____   ____ _  __| | ___ _ __ ___ 
/ __| '_ \ / _  |/ __/ _ \ | '_ \ \ / /  | |/ _  |/ _ \ '__/ __|
\__ \ |_) | (_| | (_|  __/ | | | \ V / (_| | (_| |  __/ |  \__ \
|___/ .__/ \__,_|\___\___|_|_| |_|\_/ \__,_|\__,_|\___|_|  |___/
| |
|_|
`
)

// FirstMenuItem represents the initial menu item index.
// Play represents the play menu option, assigned a value of -1.
// HighScores represent the high-scores menu option, sequentially incremented from Play.
// Howto represent the how-to menu option, sequentially incremented from HighScores.
// NumMenuItems represents the total number of menu items, calculated as the next value in the sequence.
const (
	FirstMenuItem     = 0
	Play          int = iota - 1
	HighScores
	Howto
	NumMenuItems
)

// menuItems defines a map of menu item IDs to their corresponding display strings, currently only including "PLAY".
// logoLines represent the individual lines of the logo, split by newline characters.
// logoLineLength stores the length of the first line of the logo.
// logoHeight stores the total number of lines in the logo.
var (
	//menuItems = map[int]string{Play: "PLAY", HighScores: "HIGHSCORES", Howto: "HOWTO"}
	menuItems      = map[int]string{Play: "PLAY"}
	logoLines      = strings.Split(logo, "\n")
	logoLineLength = len(logoLines[0])
	logoHeight     = len(logoLines)
)

package invaders

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// drawColor renders a string onto the given surface at specified coordinates with foreground and background colors.
// Each character is drawn sequentially, incrementing the x-coordinate.
// The rendering uses the default normal color mode for characters.
func drawColor(surface interfaces.ISurface, x, y int, fg, bg interfaces.ColorDef, data string) {
	for _, c := range data {
		surface.DrawColor(y, x, c, fg, bg, interfaces.ModeNormal)
		x++
	}
}

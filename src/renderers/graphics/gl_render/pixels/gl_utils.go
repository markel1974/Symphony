package pixels

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
)

// intBounds computes integer bounds of a Rect by flooring Min points and ceiling Max points to nearest integers.
// Returns x, y as coordinates of the bottom-left corner, and w, h as the width and height of the rectangle.
func intBounds(bounds Rect) (int, int, int, int) {
	x0 := int(math.Floor(bounds.Min.X))
	y0 := int(math.Floor(bounds.Min.Y))
	x1 := int(math.Ceil(bounds.Max.X))
	y1 := int(math.Ceil(bounds.Max.Y))
	return x0, y0, x1 - x0, y1 - y0
}

func GLGetTime() float64 {
	return glfw.GetTime()
}

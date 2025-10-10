package pixels

import "image/color"

// RGBA represents a color using red, green, blue, and alpha components, where each component is a float64.
type RGBA struct {
	R, G, B, A float64
}

// RGB creates an RGBA color with the specified red, green, and blue components, setting alpha to full opacity (1.0).
func RGB(r, g, b float64) RGBA {
	return RGBA{r, g, b, 1}
}

// Alpha returns an RGBA color with all channels set to the specified alpha value.
func Alpha(a float64) RGBA {
	return RGBA{a, a, a, a}
}

// Add returns a new RGBA color whose components are the sum of the respective components of the two input colors.
func (c RGBA) Add(d RGBA) RGBA {
	return RGBA{
		R: c.R + d.R,
		G: c.G + d.G,
		B: c.B + d.B,
		A: c.A + d.A,
	}
}

// Sub subtracts the RGBA values of another RGBA struct from the current struct and returns the resulting RGBA.
func (c RGBA) Sub(d RGBA) RGBA {
	return RGBA{
		R: c.R - d.R,
		G: c.G - d.G,
		B: c.B - d.B,
		A: c.A - d.A,
	}
}

// Mul multiplies the R, G, B, and A components of the receiver RGBA by the corresponding components of another RGBA and returns the result.
func (c RGBA) Mul(d RGBA) RGBA {
	return RGBA{
		R: c.R * d.R,
		G: c.G * d.G,
		B: c.B * d.B,
		A: c.A * d.A,
	}
}

// Scaled returns a new RGBA where each component is scaled by the specified factor.
func (c RGBA) Scaled(scale float64) RGBA {
	return RGBA{
		R: c.R * scale,
		G: c.G * scale,
		B: c.B * scale,
		A: c.A * scale,
	}
}

// RGBA converts the RGBA struct values into 16-bit unsigned integers for r, g, b, and a, scaled by 0xffff multiplier.
func (c RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(0xffff * c.R)
	g = uint32(0xffff * c.G)
	b = uint32(0xffff * c.B)
	a = uint32(0xffff * c.A)
	return
}

// ToRGBA converts a color.Color to an RGBA instance, normalizing its components to the range [0.0, 1.0].
func ToRGBA(c color.Color) RGBA {
	if c, ok := c.(RGBA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA{
		float64(r) / 0xffff,
		float64(g) / 0xffff,
		float64(b) / 0xffff,
		float64(a) / 0xffff,
	}
}

// RGBAModel is the color model for RGBA colors, converting any color.Color to an RGBA representation.
var RGBAModel = color.ModelFunc(rgbaModel)

// rgbaModel converts a color.Color to the RGBA format using the ToRGBA function.
func rgbaModel(c color.Color) color.Color {
	return ToRGBA(c)
}

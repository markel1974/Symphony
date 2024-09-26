package asciirender

import (
	"fmt"
	"math"
	"strings"
)

const (
	AsFg uint8 = iota
	AsBg
)

const (
	TplFgRGB = "38;2;%d;%d;%d"
	TplBgRGB = "48;2;%d;%d;%d"
	FgRGBPfx = "38;2;"
	BgRGBPfx = "48;2;"
)

var _increments = []uint8{0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}

type RGBColor [4]uint8

func RGB(r, g, b uint8, isBg ...bool) RGBColor {
	rgb := RGBColor{r, g, b}
	if len(isBg) > 0 && isBg[0] {
		rgb[3] = AsBg
	}
	return rgb
}

func (c RGBColor) Sprintf(format string, a ...any) string {
	return c.renderString(c.String(), fmt.Sprintf(format, a...))
}

func (c RGBColor) String() string {
	if c[3] == AsFg {
		return fmt.Sprintf(TplFgRGB, c[0], c[1], c[2])
	}
	if c[3] == AsBg {
		return fmt.Sprintf(TplBgRGB, c[0], c[1], c[2])
	}
	return ""
}

func (c RGBColor) renderString(code string, str string) string {
	if len(code) == 0 || str == "" {
		return str
	}
	if !_enable || !(_colorLevel > ColorLevelNone) {
		if !strings.Contains(str, CodeSuffix) {
			return str
		}
		return _codeRegex.ReplaceAllString(str, "")
	}
	return StartSet + code + "m" + str + ResetSet
}

func (c RGBColor) C256() Color256 {
	return C256(c.rgbTo256(c[0], c[1], c[2]), c[3] == AsBg)
}

func (c RGBColor) rgbTo256(r uint8, g uint8, b uint8) uint8 {
	res := make([]uint8, 3)
	for partI, part := range [3]uint8{r, g, b} {
		for i := 0; i < len(_increments)-1; i++ {
			smaller := _increments[i]
			bigger := _increments[i+1]
			if smaller <= part && part <= bigger {
				s1 := math.Abs(float64(smaller) - float64(part))
				b1 := math.Abs(float64(bigger) - float64(part))
				var closest uint8
				if s1 < b1 {
					closest = smaller
				} else {
					closest = bigger
				}
				res[partI] = closest
				break
			}

		}
	}
	hex := fmt.Sprintf("%02x%02x%02x", res[0], res[1], res[2])
	equiv := _hexTo256Table[hex]
	return equiv
}

package asciirender

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TplFg256 = "38;5;%d"
	TplBg256 = "48;5;%d"
	Fg256Pfx = "38;5;"
	Bg256Pfx = "48;5;"
)

type Color256 [2]uint8

func C256(val uint8, isBg ...bool) Color256 {
	bc := Color256{val}
	if len(isBg) > 0 && isBg[0] {
		bc[1] = AsBg
	}
	return bc
}

func (c Color256) Sprintf(format string, a ...any) string {
	return c.renderString(c.String(), fmt.Sprintf(format, a...))
}

func (c Color256) String() string {
	if c[1] == AsFg { // 0 is Fg
		return Fg256Pfx + strconv.Itoa(int(c[0]))
	}
	if c[1] == AsBg { // 1 is Bg
		return Bg256Pfx + strconv.Itoa(int(c[0]))
	}
	return ""
}

func (c Color256) renderString(code string, str string) string {
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

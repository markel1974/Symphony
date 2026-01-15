package vt100

import (
	"strconv"

	"github.com/markel1974/symphony/src/kernel/interfaces"
)

//var controlSequenceIntroducer = []byte{ 27, 91 }

var (
	_escMoveCursorLeftDef    = []byte{27, 91, 68}
	_escMoveCursorRightDef   = []byte{27, 91, 67}
	_escMoveCursorTopLeftDef = []byte{27, 91, 72}
	_escSaveCursorDef        = []byte{27, 55}
	_escRestoreCursorDef     = []byte{27, 56}
	_clearLine               = []byte{27, 91, 50, 75, 13}
	_clearScreen             = []byte{27, 91, 50, 74, 13}
	//	escClearLineDef =   	  []byte{ 27, 91, '2', 'K' }
)

// VT100 represents a structure for handling VT100 terminal operations and settings.
// It encapsulates functionality like cursor movement, screen clearing, and colorizing text.
// The debug field enables or disables debug logs for terminal actions.
// The enterKey field specifies the character used as the enter key.
type VT100 struct {
	enterKey rune
}

// NewVt100 initializes and returns a new VT100 instance with the specified enter key. If enterKey is invalid, it defaults to 13.
func NewVt100(enterKey rune) *VT100 {
	if enterKey < 0 {
		enterKey = 13
	}
	t := &VT100{
		enterKey: enterKey,
	}
	return t
}

// CreateEmpty returns a single space character as a string, representing an empty value for the VT100 terminal.
func (l *VT100) CreateEmpty() string {
	return " "
}

// CreateColorize applies foreground and background colors to a given text string based on the specified color mode.
func (l *VT100) CreateColorize(text string, f int, b int, mode interfaces.ColorMode) string {
	switch mode {
	case interfaces.ModeNormal:
		fg := interfaces.ColorDef(f)
		bg := interfaces.ColorDef(b)
		if fg == interfaces.ColorNoneDef && bg == interfaces.ColorNoneDef {
			return text
		}
		fgColor := l.colorFgAdapter(fg)
		if bg == interfaces.ColorNoneDef {
			return Colorize(text, fgColor)
		}
		bgColor := l.colorBgAdapter(bg)
		return ColorizeWithBackground(text, fgColor, bgColor)

	case interfaces.Mode8bit:
		if f == -1 && b == -1 {
			return text
		}
		if b == -1 {
			return Colorize8(text, f)
		}
		return Colorize8WithBackground(text, f, b)

	default:
		return text
	}
}

func (l *VT100) CreateMoveCursor(row int, col int) []byte {
	//ESC[<L>;<C>H
	//v := fmt.Sprintf("\x1b[%d;%dH", row, colum)
	//return []byte(v)
	// Pre-alloca uno slice con una capacità ragionevole per evitare riallocazioni
	// \x1b [ r o w ; c o l H -> 7 caratteri + cifre
	buf := make([]byte, 0, 16)
	buf = append(buf, 27, 91)
	buf = strconv.AppendInt(buf, int64(row+1), 10) // +1 perché VT100 è 1-based
	buf = append(buf, 59)
	buf = strconv.AppendInt(buf, int64(col+1), 10) // +1 perché VT100 è 1-based
	buf = append(buf, 72)
	return buf
}

// CreateSaveCursor returns the VT100 escape sequence to save the current cursor position.
func (l *VT100) CreateSaveCursor() []byte {
	return _escSaveCursorDef
}

// CreateRestoreCursor returns the VT100 escape sequence for restoring the cursor to the last saved position.
func (l *VT100) CreateRestoreCursor() []byte {
	return _escRestoreCursorDef
}

// CreateMoveCursorLeft generates the VT100 escape sequence for moving the cursor one position to the left.
func (l *VT100) CreateMoveCursorLeft() []byte {
	return _escMoveCursorLeftDef
}

// CreateMoveCursorRight returns a byte sequence to move the cursor one position to the right according to VT100 specifications.
func (l *VT100) CreateMoveCursorRight() []byte {
	return _escMoveCursorRightDef
}

// CreateMoveCursorTopLeft generates a byte sequence to move the cursor to the top-left position of the terminal.
func (l *VT100) CreateMoveCursorTopLeft() []byte {
	return _escMoveCursorTopLeftDef
}

// CreateClearLine sends the VT100 escape code to erase the current line and resets the cursor to the start of the line.
func (l *VT100) CreateClearLine(_ string) []byte {
	//l.ResetBuffer()
	//l.current = []rune(line)
	//l.pos = len(l.current)

	// CreateClearLine sends the VT100 code for erasing the line followed by a carriage return
	// to move the cursor back to the beginning of the line
	return _clearLine
}

// CreateClearScreen returns the VT100 escape sequence for clearing the screen and moving the cursor to the top-left position.
func (l *VT100) CreateClearScreen() []byte {
	return _clearScreen
}

// CreateSetEnterKey updates the rune value to be used as the Enter (Return) key input handler in the VT100 terminal.
func (l *VT100) CreateSetEnterKey(key rune) {
	l.enterKey = key
}

// colorFgAdapter maps a ColorDef to its corresponding foreground colorCode for VT100 terminal coloring.
func (l *VT100) colorFgAdapter(c interfaces.ColorDef) colorCode {
	var color colorCode

	switch c {
	case interfaces.ColorNoneDef:
		color = normal
	case interfaces.ColorRedDef:
		color = fgRed
	case interfaces.ColorGreenDef:
		color = fgGreen
	case interfaces.ColorYellowDef:
		color = fgYellow
	case interfaces.ColorBlueDef:
		color = fgBlue
	case interfaces.ColorMagentaDef:
		color = fgMagenta
	case interfaces.ColorCyanDef:
		color = fgCyan
	case interfaces.ColorWhiteDef:
		color = fgWhite
	case interfaces.ColorBrightRedDef:
		color = fgBrightRed
	case interfaces.ColorBrightGreenDef:
		color = fgBrightGreen
	case interfaces.ColorBrightYellowDef:
		color = fgBrightYellow
	case interfaces.ColorBrightBlueDef:
		color = fgBrightBlue
	case interfaces.ColorBrightMagentaDef:
		color = fgBrightMagenta
	case interfaces.ColorBrightCyanDef:
		color = fgBrightCyan
	case interfaces.ColorBrightWhiteDef:
		color = fgBrightWhite
	case interfaces.ColorBrightBlackDef:
		color = fgBrightBlack
	case interfaces.ColorBlackDef:
		color = fgBlack
	case interfaces.ColorGrayDef:
		color = fgBrightBlack
	case interfaces.ColorNormalDef:
		color = normal
	default:
		color = normal
	}

	return color
}

// colorBgAdapter maps a color definition to its corresponding background color code for VT100 terminal escape sequences.
func (l *VT100) colorBgAdapter(c interfaces.ColorDef) colorCode {
	var color colorCode

	switch c {
	case interfaces.ColorNoneDef:
		color = normal
	case interfaces.ColorRedDef:
		color = bgRed
	case interfaces.ColorGreenDef:
		color = bgGreen
	case interfaces.ColorYellowDef:
		color = bgYellow
	case interfaces.ColorBlueDef:
		color = bgBlue
	case interfaces.ColorMagentaDef:
		color = bgMagenta
	case interfaces.ColorCyanDef:
		color = bgCyan
	case interfaces.ColorWhiteDef:
		color = bgWhite
	case interfaces.ColorBrightRedDef:
		color = bgBrightRed
	case interfaces.ColorBrightGreenDef:
		color = bgBrightGreen
	case interfaces.ColorBrightYellowDef:
		color = bgBrightYellow
	case interfaces.ColorBrightBlueDef:
		color = bgBrightBlue
	case interfaces.ColorBrightMagentaDef:
		color = bgBrightMagenta
	case interfaces.ColorBrightCyanDef:
		color = bgBrightCyan
	case interfaces.ColorBrightWhiteDef:
		color = bgBrightWhite
	case interfaces.ColorBrightBlackDef:
		color = bgBrightBlack
	case interfaces.ColorBlackDef:
		color = bgBlack
	case interfaces.ColorGrayDef:
		color = bgBrightBlack
	case interfaces.ColorNormalDef:
		color = normal
	default:
		color = normal
	}

	return color
}

// CreateScanKey processes the input data byte slice and returns a KeyType and the corresponding rune based on the input sequence.
func (l *VT100) CreateScanKey(data []byte) (interfaces.KeyType, rune) {
	escape := false
	var escapeSequence []byte
	var escapeParameter byte
	//var escapeIntermediate byte

	if len(data) <= 0 {
		return interfaces.KeyTypeNone, 0
	}

	//UTF8
	sequence := []rune(string(data))

	for _, key := range sequence {
		if escape {
			switch len(escapeSequence) {
			case 0:
				escape = false
			case 1:
				if key == 91 {
					escapeSequence = append(escapeSequence, byte(key))
				} else {
					escape = false
				}
			case 2, 3, 4:
				if key >= 0x30 && key <= 0x3F {
					escapeParameter = byte(key)
					escapeSequence = append(escapeSequence, byte(key))
				} else if key >= 0x20 && key <= 0x2F {
					//escapeIntermediate = byte(key)
					escapeSequence = append(escapeSequence, byte(key))
				} else if key >= 0x40 && key <= 0x7E {
					escapeSequence = append(escapeSequence, byte(key))
					switch byte(key) {
					case 65:
						return interfaces.KeyTypeCursor, rune(interfaces.CursorUpDef)
					case 66:
						return interfaces.KeyTypeCursor, rune(interfaces.CursorDownDef)
					case 67:
						return interfaces.KeyTypeCursor, rune(interfaces.CursorRightDef)
					case 68:
						return interfaces.KeyTypeCursor, rune(interfaces.CursorLeftDef)
					case 126:
						if escapeParameter == 51 {
							return interfaces.KeyTypeCancel, 127
						}
					}
					escape = false
				} else {
					escape = false
				}
			default:
				escape = false
			}
		} else {
			if key == 9 {
				return interfaces.KeyTypeTab, '\t'
			} else {
				switch key {
				case l.enterKey:
					return interfaces.KeyTypeEnter, '\n'
				case 3:
					return interfaces.KeyTypeCtrl, key
				case 4:
					return interfaces.KeyTypeCtrl, key
				case 8:
					return interfaces.KeyTypeBackspace, 8
				case 27:
					escape = true
					escapeSequence = nil
					//escapeIntermediate = 0
					escapeParameter = 0
					escapeSequence = append(escapeSequence, byte(key))
				case 127:
					return interfaces.KeyTypeBackspace, 8
				default:
					return interfaces.KeyTypeKey, key
				}
			}
		}
	}
	return interfaces.KeyTypeNone, 0
}

//func (l * VT100) RetrieveBuffer() string {
//	out := string(l.current)
//	l.doReset()
//	return out
//}

/*
func (l * VT100) RetrieveTab() (string, bool) {
	var tab []rune
	found := false
	if l.pos == len(l.current) {
		tab = l.current[:l.pos]
		found = true
	}
	return string(tab), found
}
*/

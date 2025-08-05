package interfaces

type CursorCodeDef rune
type ColorDef int
type ColorMode int

const (
	ModeNormal ColorMode = iota
	Mode8bit   ColorMode = iota
)

const (
	ColorNoneDef          ColorDef = iota
	ColorRedDef           ColorDef = iota
	ColorGreenDef         ColorDef = iota
	ColorYellowDef        ColorDef = iota
	ColorBlueDef          ColorDef = iota
	ColorMagentaDef       ColorDef = iota
	ColorCyanDef          ColorDef = iota
	ColorWhiteDef         ColorDef = iota
	ColorBlackDef         ColorDef = iota
	ColorBrightRedDef     ColorDef = iota
	ColorBrightGreenDef   ColorDef = iota
	ColorBrightYellowDef  ColorDef = iota
	ColorBrightBlueDef    ColorDef = iota
	ColorBrightMagentaDef ColorDef = iota
	ColorBrightCyanDef    ColorDef = iota
	ColorBrightWhiteDef   ColorDef = iota
	ColorBrightBlackDef   ColorDef = iota
	ColorGrayDef          ColorDef = iota
	ColorNormalDef        ColorDef = iota
)

type KeyFunc func(event *KeyData)

type ITerminal interface {
	CreateEmpty() string

	CreateColorize(text string, fg int, bg int, mode ColorMode) string

	CreateSaveCursor() []byte

	CreateRestoreCursor() []byte

	CreateMoveCursorLeft() []byte

	CreateMoveCursorRight() []byte

	CreateMoveCursorTopLeft() []byte

	CreateClearLine(line string) []byte

	CreateClearScreen() []byte

	CreateScanKey(data []byte) (KeyType, rune)
}

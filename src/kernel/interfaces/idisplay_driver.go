package interfaces

// IDisplayDriver defines methods for handling terminal video display functionalities including text color and cursor control.
type IDisplayDriver interface {
	Write(p []byte) (n int, err error)

	Colorize(text string, fg int, bg int, mode ColorMode) string

	CreateSaveCursor() []byte

	CreateRestoreCursor() []byte

	CreateMoveCursorLeft() []byte

	CreateMoveCursorRight() []byte

	CreateMoveCursorTopLeft() []byte

	CreateClearLine(line string) []byte

	CreateClearScreen() []byte
}

package interfaces

type IRender interface {
	IServer

	GetScreenSize() (int, int)

	SetScreenSize(width int, height int)

	CallPaintRequest(full bool)

	CallPaintExec()

	WindowsSelectionBegin()

	WindowsSelectionPrevious()

	WindowsSelectionNext()

	WindowsSelectionEnd()

	WindowsSelectionOptions(option rune, value float64)

	Colorize(text string, fg int, bg int, mode ColorMode) string

	Write(data string)

	WriteLn(data string)

	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	WriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	ClearScreen()

	ClearLine(line string)

	MoveCursorLeft()

	MoveCursorRight()

	MoveCursorTopLeft()

	SaveCursor()

	RestoreCursor()

	EOL() string

	WritePromptLine(prompt string, line string)

	WritePromptEOL(prompt string, eol bool)

	WriteNormal(line string)

	WriteHighlight(line string)

	WriteCritical(line string)
}

package interfaces

type IRender interface {
	GetScreenSize() (int, int)

	SetScreenSize(width int, height int)

	IsDirty() bool

	ExecPaint(fgTask IProcess, tasks []IProcess) bool

	PaintRequest(full bool) bool

	Read(data []byte) (int, error)

	Write(data string)

	WriteLn(data string)

	WriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	WriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	ClearScreen()

	Scan(data []byte)

	ClearLine(line string)

	MoveCursorLeft()

	MoveCursorRight()

	SaveCursor()

	RestoreCursor()

	EOL() string

	WriteLine(prompt string, line string)

	WriteEOL(prompt string, eol bool)

	WriteNormal(line string)

	WriteHighlight(line string)

	WriteCritical(line string)
}

package interfaces

type IRender interface {
	GetScreenSize() (int, int)

	SetScreenSize(width int, height int)

	IsDirty() bool

	ExecPaint(fgTask ITask, tasks []ITask) bool

	PaintRequest(full bool) bool

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
}

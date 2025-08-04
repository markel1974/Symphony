package interfaces

type IRender interface {
	IServer

	CallGetScreenSize() (int, int)

	CallSetScreenSize(width int, height int)

	CallPaintRequest(full bool)

	CallPaintExec()

	CallWindowsSelectionBegin()

	CallWindowsSelectionPrevious()

	CallWindowsSelectionNext()

	CallWindowsSelectionEnd()

	CallWritePromptLine(prompt string, line string)

	CallWritePromptEOL(prompt string, eol bool)

	CallWindowsSelectionOptions(option rune, value float64)

	CallWriteNormal(line string)

	CallWriteHighlight(line string)

	CallWriteCritical(line string)

	CallWrite(data string)

	CallWriteLn(data string)

	CallWriteColor(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteColorLn(data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallClearScreen()

	CallSaveCursor()

	CallRestoreCursor()

	CallMoveCursorLeft()

	CallMoveCursorRight()
}

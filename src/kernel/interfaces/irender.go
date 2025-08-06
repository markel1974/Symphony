package interfaces

type IRender interface {
	IServer

	CallGetScreenSize(process IRouter) (int, int)

	CallSetScreenSize(router IRouter, width int, height int)

	CallWindowsSelectionBegin(router IRouter)

	CallWindowsSelectionOptions(router IRouter, option rune, value float64)

	CallWindowsSelectionPrevious(router IRouter)

	CallWindowsSelectionNext(router IRouter)

	CallWindowsSelectionEnd(router IRouter)

	CallWritePromptLine(router IRouter, prompt string, line string)

	CallWritePromptEOL(router IRouter, prompt string, eol bool)

	CallWriteNormal(router IRouter, line string)

	CallWriteHighlight(router IRouter, line string)

	CallWriteCritical(router IRouter, line string)

	CallWrite(router IRouter, data string)

	CallWriteLn(router IRouter, data string)

	CallWriteColor(router IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallWriteColorLn(router IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode)

	CallClearScreen(router IRouter)

	CallSaveCursor(router IRouter)

	CallRestoreCursor(router IRouter)

	CallMoveCursorLeft(router IRouter)

	CallMoveCursorRight(router IRouter)
}

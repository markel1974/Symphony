package interfaces

type IRender interface {
	IServer

	CallGetScreenSize(process IRouter) (int, int)

	CallSetScreenSize(router IRouter, width int, height int)

	CallWritePromptLine(router IRouter, prompt string, line string)

	CallWritePromptEOL(router IRouter, prompt string, eol bool)

	CallWrite(router IRouter, data string, eol bool)

	CallWriteColor(router IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode, eol bool)

	CallClearScreen(router IRouter)

	CallSaveCursor(router IRouter)

	CallRestoreCursor(router IRouter)

	CallMoveCursorLeft(router IRouter)

	CallMoveCursorRight(router IRouter)
}

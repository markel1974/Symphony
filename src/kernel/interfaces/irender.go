package interfaces

type IRender interface {
	IServer

	CallGetScreenSize(process IRouter) (int, int)

	CallSetScreenSize(router IRouter, width int, height int)

	CallWrite(router IRouter, data string, eol bool)

	CallWriteColor(router IRouter, data string, fg ColorDef, bg ColorDef, mode ColorMode, eol bool)

	CallClearLine(router IRouter, line string)

	CallClearScreen(router IRouter)

	CallSaveCursor(router IRouter)

	CallRestoreCursor(router IRouter)

	CallMoveCursorLeft(router IRouter)

	CallMoveCursorRight(router IRouter)
}

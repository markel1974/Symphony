package xvi

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// XVI is a text editor type that integrates a text user interface (Tui), a text buffer, and a user process manager.
type XVI struct {
	tui     *Tui
	buffer  *Buffer
	mode    string
	process interfaces.IUserProcess
}

// NewXVI creates and initializes a new instance of the XVI structure with default settings.
func NewXVI() *XVI {
	p := &XVI{
		mode: "normal",
	}
	return p
}

// Setup initializes the process and assigns handlers for paint, read, activate, and error events.
func (p *XVI) Setup(process interfaces.IUserProcess, args []string) {
	p.process = process
	process.SetOnPaint(p.onPaint)
	process.SetOnKey(p.onKey)
	process.SetOnActivate(p.onActivate)
	process.SetOnError(p.onError)
}

// Start initializes the XVI instance by creating and configuring the buffer and TUI components.
func (p *XVI) Start() {
	p.buffer = NewBuffer("/home/user/welcome.txt", "Benvenuto in xvi!\nPremi 'i' per entrare in modalità inserimento.\nPremi 'Esc' per tornare in modalità normale.")
	p.tui = NewTui(p.buffer)
	p.process.PaintRequest()
}

// onPaint handles the painting of the user interface by delegating drawing operations to the Tui instance.
// It renders the buffer content and the status bar on the provided ISurface.
func (p *XVI) onPaint(surface interfaces.ISurface) {
	p.tui.Draw(p.process, surface)
}

// onError handles errors encountered during the execution of the process and triggers appropriate error-handling logic.
func (p *XVI) onError(err error) {
	//TODO
}

// activateHandler handles the activation event for the process, setting up initial state or responding to reactivation.
func (p *XVI) onActivate() {
	//TODO
}

// keyHandler processes keyboard input based on the current mode (normal or insert) and updates the buffer or mode accordingly.
func (p *XVI) onKey(code int, key rune) {
	kind := interfaces.KeyType(code)

	switch p.mode {
	case "normal":
		if kind == interfaces.KeyTypeCursor {
			switch interfaces.CursorCodeDef(key) {
			case interfaces.CursorUpDef:
				p.buffer.MoveCursor(0, -1)
			case interfaces.CursorDownDef:
				p.buffer.MoveCursor(0, 1)
			case interfaces.CursorLeftDef:
				p.buffer.MoveCursor(-1, 0)
			case interfaces.CursorRightDef:
				p.buffer.MoveCursor(1, 0)
			}
		}
		switch key {
		case 'i':
			p.mode = "insert"
		case 'h':
			p.buffer.MoveCursor(-1, 0)
		case 'j':
			p.buffer.MoveCursor(0, 1)
		case 'k':
			p.buffer.MoveCursor(0, -1)
		case 'l':
			p.buffer.MoveCursor(1, 0)
		}

	case "insert":
		switch kind {
		case interfaces.KeyTypeTab:
			p.mode = "normal"
		case interfaces.KeyTypeCancel:
			p.buffer.DeleteChar()
		default:
			if key != 0 {
				p.buffer.InsertChar(key)
			}
		}
	}
	p.process.PaintRequest()
}

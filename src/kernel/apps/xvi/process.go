package xvi

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
)

// XVI is a text editor type that integrates a text user interface (Tui), a text buffer, and a user process manager.
type XVI struct {
	tui     *Tui
	buffer  *Buffer
	process interfaces.IUserProcess
}

// NewXVI creates and initializes a new instance of the XVI structure with default settings.
func NewXVI() *XVI {
	p := &XVI{}
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
	switch p.tui.GetMode() {
	case "normal":
		p.doNormalMode(code, key)
	case "insert":
		p.doInsertMode(code, key)
	}
	p.process.PaintRequest()
}

// doInsertMode processes key input in insert mode, modifying the buffer or changing the mode based on the key type.
func (p *XVI) doInsertMode(code int, key rune) {
	if p.doCursor(code, key) {
		return
	}
	kind := interfaces.KeyType(code)
	switch kind {
	case interfaces.KeyTypeTab:
		p.tui.SetMode("normal")
	case interfaces.KeyTypeCancel:
		p.buffer.DeleteChar()
	default:
		if key < 32 {
			if key == '\n' {
				p.buffer.InsertRow()
			}
			//nothing for now
		} else {
			p.buffer.InsertChar(key)
		}
	}
}

// doNormalMode processes input in normal mode, handling cursor movements and switching to insert mode when 'i' is pressed.
func (p *XVI) doNormalMode(code int, key rune) {
	if p.doCursor(code, key) {
		return
	}
	switch key {
	case 'i':
		p.tui.SetMode("insert")
	case 'h':
		//p.buffer.MoveCursor(-1, 0)
	case 'j':
		//p.buffer.MoveCursor(0, 1)
	case 'k':
		//p.buffer.MoveCursor(0, -1)
	case 'l':
		//p.buffer.MoveCursor(1, 0)
	}
}

func (p *XVI) doCursor(code int, key rune) bool {
	kind := interfaces.KeyType(code)
	if kind == interfaces.KeyTypeCursor {
		switch interfaces.CursorCodeDef(key) {
		case interfaces.CursorUpDef:
			p.buffer.MoveCursor(0, -1)
			return true
		case interfaces.CursorDownDef:
			p.buffer.MoveCursor(0, 1)
			return true
		case interfaces.CursorLeftDef:
			p.buffer.MoveCursor(-1, 0)
			return true
		case interfaces.CursorRightDef:
			p.buffer.MoveCursor(1, 0)
			return true
		}
	}
	return false
}

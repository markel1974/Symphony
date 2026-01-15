package xvi

import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
)

// XVI is a text editor type that integrates a text user interface (Tui), a text buffer, and a user process manager.
type XVI struct {
	tui                    *Tui
	buffer                 *Buffer
	process                interfaces.IUserProcess
	yankBuffer             string
	lastChar               rune
	isWelcomeMessageActive bool
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
	p.isWelcomeMessageActive = true
}

// onPaint handles the painting of the user interface by delegating drawing operations to the Tui instance.
func (p *XVI) onPaint(surface interfaces.ISurface) {
	p.tui.Draw(p.process, surface)
}

// onError handles errors encountered during the execution of the process.
func (p *XVI) onError(err error) {
	p.tui.SetError(err)
}

// onActivate handles the activation event for the process.
func (p *XVI) onActivate() {
	//TODO
}

// onKey processes keyboard input based on the current mode.
func (p *XVI) onKey(code int, key rune) {
	if p.isWelcomeMessageActive {
		p.buffer.Clear()
		p.isWelcomeMessageActive = false
	}

	switch p.tui.GetMode() {
	case "normal":
		p.doNormalMode(code, key)
	case "insert":
		p.doInsertMode(code, key)
	case "command":
		p.doCommandMode(code, key)
	}
	p.process.PaintRequest()
}

// doInsertMode processes key input in insert mode.
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
		if key == '\n' {
			p.buffer.InsertRow()
		} else if key >= 32 {
			p.buffer.InsertChar(key)
		}
	}
}

// doNormalMode processes input in normal mode.
func (p *XVI) doNormalMode(code int, key rune) {
	if p.doCursor(code, key) {
		p.lastChar = 0 // Reset sequence on cursor movement
		return
	}

	// Handle command sequences like 'dd' and 'yy'
	if p.lastChar != 0 {
		if p.lastChar == 'd' && key == 'd' {
			_, y := p.buffer.Cursor()
			p.yankBuffer = p.buffer.DeleteLine(y)
			p.lastChar = 0
			return
		}
		if p.lastChar == 'y' && key == 'y' {
			_, y := p.buffer.Cursor()
			p.yankBuffer = p.buffer.GetLine(y)
			p.tui.SetError(nil)
			p.lastChar = 0
			return
		}
		p.lastChar = 0
	}

	switch key {
	case 'i':
		p.tui.SetMode("insert")
	case ':':
		p.tui.SetMode("command")
	case 'h':
		p.buffer.MoveCursor(-1, 0)
	case 'j':
		p.buffer.MoveCursor(0, 1)
	case 'k':
		p.buffer.MoveCursor(0, -1)
	case 'l':
		p.buffer.MoveCursor(1, 0)
	case 'd':
		p.lastChar = 'd'
	case 'y':
		p.lastChar = 'y'
	case 'p':
		if p.yankBuffer != "" {
			_, y := p.buffer.Cursor()
			p.buffer.InsertLineBelow(y, p.yankBuffer)
		}
	default:
		p.lastChar = 0
	}
}

// doCommandMode processes input in command mode.
func (p *XVI) doCommandMode(code int, key rune) {
	kind := interfaces.KeyType(code)
	if kind == interfaces.KeyTypeTab {
		p.tui.SetMode("normal")
		return
	}
}

// doCursor handles universal cursor movement keys.
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

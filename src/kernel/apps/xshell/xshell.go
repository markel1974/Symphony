package xshell

import (
	"fmt"
	"log"
	"unicode"

	"github.com/markel1974/symphony/src/kernel/interfaces"
)

// XShell defines a command-line interface entity with support for input management, authentication, rendering, and history.
type XShell struct {
	process         interfaces.IUserProcess
	current         []rune
	pos             int
	echo            bool
	history         *HistoryHandler
	tabData         string
	tabFound        bool
	tabCount        int
	defaultPrompt   string
	prompt          string
	currentUsername string
	passwordRetry   int
	selectionMode   bool
}

// NewXShell initializes and returns a new instance of *Shell configured with dependencies and initial settings.
func NewXShell(prompt string, autosave bool) *XShell {
	history := NewHistoryHandler(128, autosave)
	history.SetEnabled(true)
	c := &XShell{
		echo:          true,
		defaultPrompt: prompt,
		prompt:        prompt,
		history:       history,
		selectionMode: false,
	}
	return c
}

func (c *XShell) Setup(process interfaces.IUserProcess) {
	c.process = process
	process.SetOnKey(c.keyHandler)
	process.SetOnKeyBroadcast(c.broadcastKeyHandler)
	process.SetOnError(c.ErrorHandler)
	process.SetOnActivate(c.activateHandler)
}

func (c *XShell) Start(process interfaces.IUserProcess) {
	process.WriteForeground("Admin Console Ready", interfaces.ColorBlueDef, true)
}

// ErrorHandler handles errors by writing the error message to the process output and proceeding to the next line.
func (c *XShell) ErrorHandler(err error) {
	c.process.Write(err.Error(), true)
	c.nextLine(false)
}

// ActivateHandler triggers the processing loop for the specified process, ensuring the next line is handled without a delay.
func (c *XShell) activateHandler() {
	c.nextLine(false)
}

func (c *XShell) broadcastKeyHandler(code int, key rune) {
	kind := interfaces.KeyType(code)
	if kind == interfaces.KeyTypeCtrl {
		if key == 3 {
			//ctrl-c
			c.selectionMode = false
			c.process.KillForeground()
		} else if key == 4 {
			//ctrl-d
			c.selectionMode = !c.selectionMode
			if c.selectionMode {
				c.process.WindowsSelectionBegin()
			} else {
				c.process.WindowsSelectionEnd()
				c.process.ProcessSetSelfForeground()
			}
		}
		return
	}

	if c.selectionMode {
		c.keyHandlerSelection(code, key)
	}
}

// KeyHandler handles keyboard input events for the XShell context, adjusting behavior based on selection mode status.
func (c *XShell) keyHandler(code int, key rune) {
	if !c.selectionMode {
		c.keyHandlerNormal(code, key)
	}
}

// KeyHandler handles keyboard inputs based on the provided key type and key value, executing corresponding actions.
func (c *XShell) keyHandlerNormal(code int, key rune) {
	kind := interfaces.KeyType(code)
	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			c.nextLine(true)
		case 4:
			//c.handleExecActivate()
		}
		return
	}

	switch kind {
	case interfaces.KeyTypeEnter:
		c.enterPressed()
		c.tabCount = 0
	case interfaces.KeyTypeTab:
		c.tabPressed()
	case interfaces.KeyTypeCancel:
		c.tabCount = 0
		c.textCancel()
	case interfaces.KeyTypeBackspace:
		c.tabCount = 0
		c.textBackspace()
	case interfaces.KeyTypeKey:
		c.tabCount = 0
		c.keyPressed(key)
	case interfaces.KeyTypeCursor:
		c.tabCount = 0
		c.cursorPressed(interfaces.CursorCodeDef(key))
	default:
		log.Println("handlerKeyEvent: Unknown key type")
	}
}

func (c *XShell) keyHandlerSelection(code int, key rune) {
	kind := interfaces.KeyType(code)
	if kind == interfaces.KeyTypeCursor {
		switch interfaces.CursorCodeDef(key) {
		case interfaces.CursorUpDef:
			c.process.WindowsSelectionOptions('y', -1)
		case interfaces.CursorDownDef:
			c.process.WindowsSelectionOptions('y', 1)
		case interfaces.CursorLeftDef:
			c.process.WindowsSelectionOptions('x', -1)
		case interfaces.CursorRightDef:
			c.process.WindowsSelectionOptions('x', 1)
		}
		return
	}

	switch key {
	case 'w':
		c.process.WindowsSelectionOptions('y', -1)
	case 's':
		c.process.WindowsSelectionOptions('y', 1)
	case 'a':
		c.process.WindowsSelectionOptions('x', -1)
	case 'd':
		c.process.WindowsSelectionOptions('x', 1)
	case '+':
		c.process.WindowsSelectionOptions('z', 0.1)
	case '-':
		c.process.WindowsSelectionOptions('z', -0.1)
	case '\t':
		c.process.WindowsSelectionNext()
	case 'q':
		c.process.WindowsSelectionPrevious()
	}
}

// nextLine resets the input buffer and renders the prompt and EOL markers with specified colors and styles.
func (c *XShell) nextLine(eol bool) {
	c.current = nil
	c.pos = 0
	c.prompt = c.process.CWDName() + c.defaultPrompt
	c.process.WritePromptEOL(c.prompt, eol)
}

// Redraw refreshes the current shell display with the given line, updates internal state, and re-renders the prompt and line.
func (c *XShell) redraw(process interfaces.IUserProcess, line string) {
	c.current = []rune(line)
	c.pos = len(c.current)
	c.prompt = process.CWDName() + c.defaultPrompt
	process.WritePromptLine(c.prompt, line)
}

// HistoryApply performs actions on the command history based on the specified verb (list, clear, or execute at the given index).
func (c *XShell) historyApply(process interfaces.IUserProcess, verb interfaces.HistoryAction, idx int) string {
	switch verb {
	case interfaces.HistoryActionClear:
		c.history.Clear()
	case interfaces.HistoryActionExec:
		if arg, found := c.history.GetHistoryAtPos(idx); found {
			return arg
		}
	case interfaces.HistoryActionList:
		out := ""
		for n, x := range c.history.GetHistory() {
			out += "\r\n"
			out += fmt.Sprintf("%d: %s", n, x)
		}
		process.Write(out, false)
	}
	return ""
}

// HistorySuggest suggests autocompletion options based on input and handles cycling through suggestions on repeated calls.
func (c *XShell) historySuggest(process interfaces.IUserProcess, data string, suggestions []string, found bool) {
	if found && len(suggestions) > 0 {
		sLen := len(suggestions)
		if idx := c.tabCount % sLen; idx < sLen {
			if complete := suggestions[idx]; len(complete) > len(data) {
				tabLine := complete
				c.redraw(process, tabLine)
				c.history.SetDefault(tabLine)
				if sLen == 1 {
					c.tabCount = 0
					c.history.SetDefault(string(c.current))
				}
			}
		}
	}
}

// textBackspace removes the character at the current cursor position and updates the Shell state accordingly.
func (c *XShell) textBackspace() {
	if c.pos > 0 {
		c.pos--
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			c.process.MoveCursorLeft()
			c.process.SaveCursor()
			c.process.WriteForeground(string(c.current[c.pos:]), interfaces.ColorNoneDef, false)
			c.process.WriteForeground(string(' '), interfaces.ColorNoneDef, false)
			c.process.RestoreCursor()
		}
	}
}

// textCancel removes the character at the current cursor position and updates the display if echo mode is enabled.
func (c *XShell) textCancel() {
	if c.pos >= 0 {
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			c.process.SaveCursor()
			c.process.WriteForeground(string(c.current[c.pos:]), interfaces.ColorNoneDef, false)
			c.process.WriteForeground(string(' '), interfaces.ColorNoneDef, false)
			c.process.RestoreCursor()
		}
	}
}

// cursorPressed handles cursor navigation events based on the given CursorCodeDef.
func (c *XShell) cursorPressed(code interfaces.CursorCodeDef) {
	switch code {
	case interfaces.CursorUpDef:
		if data, valid := c.history.GetHistoryPrev(); valid {
			c.redraw(c.process, data)
		}
	case interfaces.CursorDownDef:
		if data, valid := c.history.GetHistoryNext(); valid {
			c.redraw(c.process, data)
		}
	case interfaces.CursorLeftDef:
		if c.pos > 0 {
			c.pos--
			c.process.MoveCursorLeft()
		}
	case interfaces.CursorRightDef:
		if c.pos >= 0 && c.pos < len(c.current) {
			c.pos++
			c.process.MoveCursorRight()
		}
	}
}

// EnterPressed handles the Enter key press event, processes the input based on the current shell state, and updates the state accordingly.
func (c *XShell) enterPressed() {
	buffer := string(c.current)
	if len(buffer) == 0 {
		c.nextLine(true)
		return
	}
	c.history.AddToHistory(buffer)
	c.history.SetDefault("")
	c.process.Write("", true)
	c.process.ProcessExec(buffer)
}

// TabPressed handles tab key events, providing intelligent autocompletion based on current input and command context.
func (c *XShell) tabPressed() {
	if c.tabCount == 0 {
		c.tabFound = false
		c.tabData = ""
		if c.pos >= 0 && c.pos <= len(c.current) {
			c.tabData = string(c.current)
			c.tabFound = true
		}
	}
	c.tabCount++

	if c.tabFound {
		data, suggestions, found := c.process.Suggestion(c.tabData, c.pos)
		c.historySuggest(c.process, data, suggestions, found)
	}
}

// keyPressed processes a printable key input and updates the current input buffer, cursor position, and visual rendering.
func (c *XShell) keyPressed(key rune) {
	if unicode.IsPrint(key) {
		if c.pos < 0 {
			log.Println("doTextInsert: negative pos", c.pos)
		} else if c.pos == len(c.current) {
			if c.echo {
				c.process.WriteForeground(string(key), interfaces.ColorNoneDef, false)
			}
			c.current = append(c.current, key)
			c.pos++
		} else if c.pos < len(c.current) {
			if c.echo {
				c.process.WriteForeground(string(key), interfaces.ColorNoneDef, false)
				c.process.SaveCursor()
				c.process.WriteForeground(string(c.current[c.pos:]), interfaces.ColorNoneDef, false)
			}
			c.process.RestoreCursor()
			c.current = insertAtPos(c.current, key, c.pos)
			c.pos++
		} else {
			log.Println("terminalKeyPressed: invalid pos", c.pos)
		}
	}
	c.history.SetDefault(string(c.current))
}

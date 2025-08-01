package xshell

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"log"
	"unicode"
)

// XShell defines a command-line interface entity with support for input management, authentication, rendering, and history.
type XShell struct {
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
	}
	return c
}

// KeyHandler handles keyboard inputs based on the provided key type and key value, executing corresponding actions.
func (c *XShell) KeyHandler(process interfaces.IProcess, code int, key rune) {
	kind := interfaces.KeyType(code)
	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			//c.selector.Clear()
			//c.CallTaskKillForeground()
			c.NextLine(process, true)
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
		c.textCancel(process)
	case interfaces.KeyTypeBackspace:
		c.tabCount = 0
		c.textBackspace(process)
	case interfaces.KeyTypeKey:
		c.tabCount = 0
		c.keyPressed(process, key)
	case interfaces.KeyTypeCursor:
		c.tabCount = 0
		c.cursorPressed(process, interfaces.CursorCodeDef(key))
	default:
		log.Println("handlerKeyEvent: Unknown key type")
	}

	if kind == interfaces.KeyTypeEnter {
		if buffer := c.InputBuffer(); len(buffer) > 0 {
			process.WriteLn("")
			_, _ = process.TaskExec(buffer)
			c.NextLine(process, false)
		} else {
			c.NextLine(process, true)
		}
	} else if kind == interfaces.KeyTypeTab {
		if c.tabFound {
			data, suggestions, found := process.Suggestion(c.tabData, c.pos)
			c.historySuggest(process, data, suggestions, found)
		}
	}
}

// NextLine resets the input buffer and renders the prompt and EOL markers with specified colors and styles.
func (c *XShell) NextLine(process interfaces.IProcess, eol bool) {
	c.resetBuffer()
	process.WritePromptEOL(c.prompt, eol)
}

// Redraw refreshes the current shell display with the given line, updates internal state, and re-renders the prompt and line.
func (c *XShell) redraw(process interfaces.IProcess, line string) {
	c.current = []rune(line)
	c.pos = len(c.current)
	process.WritePromptLine(c.prompt, line)
}

// SetPromptPrefix sets a custom prefix for the prompt by prepending the given prefix to the default prompt value.
func (c *XShell) SetPromptPrefix(prefix string) {
	c.prompt = prefix + c.defaultPrompt
}

// HistoryApply performs actions on the command history based on the specified verb (list, clear, or execute at the given index).
func (c *XShell) HistoryApply(task interfaces.IProcess, verb interfaces.HistoryAction, idx int) string {
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
		task.Write(out)
	}
	return ""
}

// HistorySuggest suggests autocompletion options based on input and handles cycling through suggestions on repeated calls.
func (c *XShell) historySuggest(process interfaces.IProcess, data string, suggestions []string, found bool) {
	if found && len(suggestions) > 0 {
		sLen := len(suggestions)
		if idx := c.tabCount % sLen; idx < sLen {
			if complete := suggestions[idx]; len(complete) > len(data) {
				tabLine := complete
				c.redraw(process, tabLine)
				c.history.SetDefault(tabLine)
				if sLen == 1 {
					c.TabReset()
				}
			}
		}
	}
}

// InputBuffer retrieves the current input buffer as a string.
func (c *XShell) InputBuffer() string {
	return string(c.current)
}

// TabData retrieves the tab-related data including a string representation, position, and tab count from the Shell instance.
//func (c *XShell) TabData() (string, int, bool) {
//	return c.tabData, c.pos, c.tabFound
//}

func (c *XShell) TabReset() {
	c.tabCount = 0
	c.history.SetDefault(string(c.current))
}

// textBackspace removes the character at the current cursor position and updates the Shell state accordingly.
func (c *XShell) textBackspace(task interfaces.IProcess) {
	if c.pos > 0 {
		c.pos--
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			task.MoveCursorLeft()
			task.SaveCursor()
			task.WriteNormal(string(c.current[c.pos:]))
			task.WriteNormal(string(' '))
			task.RestoreCursor()
		}
	}
}

// textCancel removes the character at the current cursor position and updates the display if echo mode is enabled.
func (c *XShell) textCancel(process interfaces.IProcess) {
	if c.pos >= 0 {
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			process.SaveCursor()
			process.WriteNormal(string(c.current[c.pos:]))
			process.WriteNormal(string(' '))
			process.RestoreCursor()
		}
	}
}

// resetBuffer clears the current input buffer by resetting it to nil and setting the buffer position to zero.
func (c *XShell) resetBuffer() {
	c.current = nil
	c.pos = 0
}

// cursorPressed handles cursor navigation events based on the given CursorCodeDef.
func (c *XShell) cursorPressed(process interfaces.IProcess, code interfaces.CursorCodeDef) {
	switch code {
	case interfaces.CursorUpDef:
		if data, valid := c.history.GetHistoryPrev(); valid {
			c.redraw(process, data)
		}
	case interfaces.CursorDownDef:
		if data, valid := c.history.GetHistoryNext(); valid {
			c.redraw(process, data)
		}
	case interfaces.CursorLeftDef:
		if c.pos > 0 {
			c.pos--
			process.MoveCursorLeft()
		}
	case interfaces.CursorRightDef:
		if c.pos >= 0 && c.pos < len(c.current) {
			c.pos++
			process.MoveCursorRight()
		}
	}
}

// EnterPressed handles the Enter key press event, processes the input based on the current shell state, and updates the state accordingly.
func (c *XShell) enterPressed() bool {
	buffer := string(c.current)
	if len(buffer) > 0 {
		c.history.AddToHistory(buffer)
		c.history.SetDefault("")
	}
	return false
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
}

// keyPressed processes a printable key input and updates the current input buffer, cursor position, and visual rendering.
func (c *XShell) keyPressed(process interfaces.IProcess, key rune) {
	if unicode.IsPrint(key) {
		if c.pos < 0 {
			log.Println("doTextInsert: negative pos", c.pos)
		} else if c.pos == len(c.current) {
			if c.echo {
				process.WriteNormal(string(key))
			}
			c.current = append(c.current, key)
			c.pos++
		} else if c.pos < len(c.current) {
			if c.echo {
				process.WriteNormal(string(key))
				process.SaveCursor()
				process.WriteNormal(string(c.current[c.pos:]))
			}
			process.RestoreCursor()
			c.current = insertAtPos(c.current, key, c.pos)
			c.pos++
		} else {
			log.Println("terminalKeyPressed: invalid pos", c.pos)
		}
	}

	c.history.SetDefault(string(c.current))
	//c.tabCount = 0
}

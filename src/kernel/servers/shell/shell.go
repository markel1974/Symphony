/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package shell

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"log"
	"unicode"
)

// stateUndefined represents an undefined state in the authentication process.
// stateUsernameRequired indicates that a username is required for authentication.
// statePasswordRequired indicates that a password is required for authentication.
// stateAuthenticated represents a successful authentication state.
const (
	stateUndefined        = iota
	stateUsernameRequired = iota
	statePasswordRequired = iota
	stateAuthenticated    = iota
)

// usernamePrompt is the prompt text displayed when asking for the username.
// passwordPrompt is the prompt text displayed when asking for the password.
// maxPasswordRetry defines the maximum number of password retry attempts allowed.
const (
	usernamePrompt   = "Username: "
	passwordPrompt   = "Password: "
	maxPasswordRetry = 3
)

// Shell defines a command-line interface entity with support for input management, authentication, rendering, and history.
type Shell struct {
	render          interfaces.IRender
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
	state           int
	auth            interfaces.IAuthenticator
}

// NewShell initializes and returns a new instance of *Shell configured with dependencies and initial settings.
func NewShell(auth interfaces.IAuthenticator, render interfaces.IRender, prompt string, autosave bool) *Shell {
	c := &Shell{
		render:        render,
		history:       NewHistoryHandler(128, autosave),
		echo:          true,
		auth:          auth,
		defaultPrompt: prompt,
		passwordRetry: 0,
		state:         stateUndefined,
	}
	if auth.IsAuthenticated() {
		c.setAuthenticatedState()
	}
	return c
}

// KeyHandler handles keyboard inputs based on the provided key type and key value, executing corresponding actions.
func (c *Shell) KeyHandler(kind interfaces.KeyType, key rune) {
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

// NextLine resets the input buffer and renders the prompt and EOL markers with specified colors and styles.
func (c *Shell) NextLine(eol bool) {
	c.resetBuffer()
	c.render.WritePromptEOL(c.prompt, eol)
}

// Redraw refreshes the current shell display with the given line, updates internal state, and re-renders the prompt and line.
func (c *Shell) Redraw(line string) {
	c.current = []rune(line)
	c.pos = len(c.current)
	c.render.WritePromptLine(c.prompt, line)
}

// SetPromptPrefix sets a custom prefix for the prompt by prepending the given prefix to the default prompt value.
func (c *Shell) SetPromptPrefix(prefix string) {
	c.prompt = prefix + c.defaultPrompt
}

// HistoryApply performs actions on the command history based on the specified verb (list, clear, or execute at the given index).
func (c *Shell) HistoryApply(verb interfaces.HistoryAction, idx int) string {
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
		c.render.Write(out)
	}
	return ""
}

// HistorySuggest suggests autocompletion options based on input and handles cycling through suggestions on repeated calls.
func (c *Shell) HistorySuggest(data string, suggestions []string, found bool) {
	if found && len(suggestions) > 0 {
		sLen := len(suggestions)
		if idx := c.tabCount % sLen; idx < sLen {
			if complete := suggestions[idx]; len(complete) > len(data) {
				tabLine := complete
				c.Redraw(tabLine)
				c.history.SetDefault(tabLine)
				if sLen == 1 {
					c.TabReset()
				}
			}
		}
	}
}

// Authenticated checks if the Shell instance is in an authenticated state and returns true if authenticated, otherwise false.
func (c *Shell) Authenticated() bool {
	return c.state == stateAuthenticated
}

// InputBuffer retrieves the current input buffer as a string.
func (c *Shell) InputBuffer() string {
	return string(c.current)
}

// TabData retrieves the tab-related data including a string representation, position, and tab count from the Shell instance.
func (c *Shell) TabData() (string, int, bool) {
	return c.tabData, c.pos, c.tabFound
}

func (c *Shell) TabReset() {
	c.tabCount = 0
	c.history.SetDefault(string(c.current))
}

// setUsernameRequiredState transitions the Shell into a state where the username is required and disables history.
func (c *Shell) setUsernameRequiredState() {
	c.echo = true
	c.prompt = usernamePrompt
	c.history.SetEnabled(false)
	c.state = stateUsernameRequired
}

// setPasswordRequiredState sets the shell state to require a password, disables history, hides input echo, and sets the prompt to "Password:".
func (c *Shell) setPasswordRequiredState() {
	c.echo = false
	c.prompt = passwordPrompt
	c.history.SetEnabled(false)
	c.state = statePasswordRequired
}

// setAuthenticatedState updates the shell to an authenticated state, enabling command history and setting the default prompt.
func (c *Shell) setAuthenticatedState() {
	c.echo = true
	c.prompt = c.defaultPrompt
	c.history.SetEnabled(true)
	c.state = stateAuthenticated
}

// textBackspace removes the character at the current cursor position and updates the Shell state accordingly.
func (c *Shell) textBackspace() {
	if c.pos > 0 {
		c.pos--
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			c.render.MoveCursorLeft()
			c.render.SaveCursor()
			c.render.WriteNormal(string(c.current[c.pos:]))
			c.render.WriteNormal(string(' '))
			c.render.RestoreCursor()
		}
	}
}

// textCancel removes the character at the current cursor position and updates the display if echo mode is enabled.
func (c *Shell) textCancel() {
	if c.pos >= 0 {
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			c.render.SaveCursor()
			c.render.WriteNormal(string(c.current[c.pos:]))
			c.render.WriteNormal(string(' '))
			c.render.RestoreCursor()
		}
	}
}

// resetBuffer clears the current input buffer by resetting it to nil and setting the buffer position to zero.
func (c *Shell) resetBuffer() {
	c.current = nil
	c.pos = 0
}

// cursorPressed handles cursor navigation events based on the given CursorCodeDef.
func (c *Shell) cursorPressed(code interfaces.CursorCodeDef) {
	switch code {
	case interfaces.CursorUpDef:
		if data, valid := c.history.GetHistoryPrev(); valid {
			c.Redraw(data)
		}
	case interfaces.CursorDownDef:
		if data, valid := c.history.GetHistoryNext(); valid {
			c.Redraw(data)
		}
	case interfaces.CursorLeftDef:
		if c.pos > 0 {
			c.pos--
			c.render.MoveCursorLeft()
		}
	case interfaces.CursorRightDef:
		if c.pos >= 0 && c.pos < len(c.current) {
			c.pos++
			c.render.MoveCursorRight()
		}
	}
}

// EnterPressed handles the Enter key press event, processes the input based on the current shell state, and updates the state accordingly.
func (c *Shell) enterPressed() bool {
	buffer := string(c.current)
	quit := false

	switch c.state {
	case stateUsernameRequired:
		c.passwordRetry = 0
		c.currentUsername = buffer
		c.setPasswordRequiredState()

	case statePasswordRequired:
		if c.auth.Authenticate(c.currentUsername, buffer) {
			c.setAuthenticatedState()
		} else {
			c.passwordRetry++
			eol := c.render.EOL()
			if c.passwordRetry >= maxPasswordRetry {
				c.render.WriteCritical(eol + "Unauthorized" + eol)
				quit = true
			} else {
				c.render.WriteCritical(eol + "Login incorrect" + eol)
			}
		}

	case stateAuthenticated:
		if len(buffer) > 0 {
			c.history.AddToHistory(buffer)
			c.history.SetDefault("")
		}
	default:
		quit = true
	}
	return quit
}

// TabPressed handles tab key events, providing intelligent autocompletion based on current input and command context.
func (c *Shell) tabPressed() {
	if c.state != stateAuthenticated {
		return
	}
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
func (c *Shell) keyPressed(key rune) {
	if unicode.IsPrint(key) {
		if c.pos < 0 {
			log.Println("doTextInsert: negative pos", c.pos)
		} else if c.pos == len(c.current) {
			if c.echo {
				c.render.WriteNormal(string(key))
			}
			c.current = append(c.current, key)
			c.pos++
		} else if c.pos < len(c.current) {
			if c.echo {
				c.render.WriteNormal(string(key))
				c.render.SaveCursor()
				c.render.WriteNormal(string(c.current[c.pos:]))
			}
			c.render.RestoreCursor()

			c.current = insertAtPos(c.current, key, c.pos)
			c.pos++
		} else {
			log.Println("terminalKeyPressed: invalid pos", c.pos)
		}
	}

	if c.state == stateAuthenticated {
		c.history.SetDefault(string(c.current))
		//c.tabCount = 0
	}
}

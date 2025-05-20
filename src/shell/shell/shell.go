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
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"log"
	"unicode"
)

const (
	stateUndefined        = iota
	stateUsernameRequired = iota
	statePasswordRequired = iota
	stateAuthenticated    = iota
)

const (
	usernamePrompt   = "Username: "
	passwordPrompt   = "Password: "
	maxPasswordRetry = 3
)

type IExecutor interface {
	ExecCommand(line string) (bool, error)
	ExecSuggestion(in string, cursor int, count int) (int, bool)
}

type Shell struct {
	current         []rune
	pos             int
	echo            bool
	history         *HistoryHandler
	tabData         string
	tabFound        bool
	tabCount        int
	render          interfaces.IRender
	executor        IExecutor
	defaultPrompt   string
	prompt          string
	currentUsername string
	passwordRetry   int
	state           int
	auth            interfaces.IAuthenticator
}

func NewShell(auth interfaces.IAuthenticator, render interfaces.IRender, executor IExecutor, prompt string, autosave bool) *Shell {
	c := &Shell{
		history:       NewHistoryHandler(128, autosave),
		echo:          true,
		render:        render,
		executor:      executor,
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

func (c *Shell) KeyEvent(kind interfaces.KeyType, key rune) bool {
	ret := false
	switch kind {
	case interfaces.KeyTypeEnter:
		c.tabCount = 0
		ret = c.enterPressed()
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
		log.Println("KeyEvent: Unknown key type")
	}
	return ret
}

func (c *Shell) ClearHistory() {
	c.history.Clear()
}

func (c *Shell) GetHistoryAtPos(idx int) (string, bool) {
	return c.history.GetHistoryAtPos(idx)
}

func (c *Shell) GetHistory() string {
	out := ""
	for n, x := range c.history.GetHistory() {
		out += "\r\n"
		out += fmt.Sprintf("%d: %s", n, x)
	}
	return out
}

func (c *Shell) SetHistoryDefault(data string) {
	c.history.SetDefault(data)
}

func (c *Shell) NextLine() {
	c.resetBuffer()
	c.render.WriteColor("\r\n", interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	c.render.WriteColor(c.prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

func (c *Shell) Redraw(line string) {
	c.current = []rune(line)
	c.pos = len(c.current)
	c.render.ClearLine(line)
	c.render.WriteColor(c.prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	c.render.WriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

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

func (c *Shell) enterPressed() bool {
	buffer := string(c.current)
	quit := false

	if len(buffer) > 0 {
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
				if c.passwordRetry >= maxPasswordRetry {
					c.render.WriteColor("\r\nUnauthorized\r\n", interfaces.ColorRedDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
					quit = true
				} else {
					c.render.WriteColor("\r\nLogin incorrect\r\n", interfaces.ColorRedDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
				}
			}

		case stateAuthenticated:
			c.history.AddToHistory(buffer)
			c.history.SetDefault("")
			_, _ = c.executor.ExecCommand(buffer)

			//c.quit = c.execCommand(buffer)
		default:
			quit = true
		}
	}

	c.NextLine()

	return quit
}

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
	if c.tabFound {
		if l, ok := c.executor.ExecSuggestion(c.tabData, c.pos, c.tabCount); l == 1 && ok {
			c.tabCount = 0
			c.history.SetDefault(string(c.current))
		}
	}
}

func (c *Shell) keyPressed(key rune) {
	if unicode.IsPrint(key) {
		if c.pos < 0 {
			log.Println("doTextInsert: negative pos", c.pos)
		} else if c.pos == len(c.current) {
			if c.echo {
				c.render.WriteColor(string(key), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			}
			c.current = append(c.current, key)
			c.pos++
		} else if c.pos < len(c.current) {
			if c.echo {
				c.render.WriteColor(string(key), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
				c.render.SaveCursor()
				c.render.WriteColor(string(c.current[c.pos:]), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
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

func (c *Shell) setUsernameRequiredState() {
	c.echo = true
	c.prompt = usernamePrompt
	c.history.SetEnabled(false)
	c.state = stateUsernameRequired
}

func (c *Shell) setPasswordRequiredState() {
	c.echo = false
	c.prompt = passwordPrompt
	c.history.SetEnabled(false)
	c.state = statePasswordRequired
}

func (c *Shell) setAuthenticatedState() {
	c.echo = true
	c.prompt = c.defaultPrompt
	c.history.SetEnabled(true)
	c.state = stateAuthenticated
}

func (c *Shell) textBackspace() {
	if c.pos > 0 {
		c.pos--
		c.current = removeAtPos(c.current, c.pos)
		if c.echo {
			c.render.MoveCursorLeft()
			c.render.SaveCursor()
			c.render.WriteColor(string(c.current[c.pos:]), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			c.render.WriteColor(string(' '), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			c.render.RestoreCursor()
		}
	}
}

func (c *Shell) textCancel() {
	if c.pos >= 0 {
		c.current = removeAtPos(c.current, c.pos)

		if c.echo {
			c.render.SaveCursor()
			c.render.WriteColor(string(c.current[c.pos:]), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			c.render.WriteColor(string(' '), interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
			c.render.RestoreCursor()
		}
	}
}

func (c *Shell) resetBuffer() {
	c.current = nil
	c.pos = 0
}

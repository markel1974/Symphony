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

package context

import (
	"github.com/markel1974/c64emu/src/shell/adaptiveticker"
	"github.com/markel1974/c64emu/src/shell/apps"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"github.com/markel1974/c64emu/src/shell/shell"
	"github.com/markel1974/c64emu/src/shell/terminal"
	"io"
)

const (
	contextMaQueueLen = 1024
)

const eol = "\r\n"

type Context struct {
	Exit        bool
	ticker      *adaptiveticker.AdaptiveTicker
	reader      io.Reader
	writer      io.Writer
	factory     *terminal.EquipmentFactory
	commands    *cli.Command
	terminal    interfaces.ITerminal
	auth        interfaces.IAuthenticator
	defaultApp  *shell.Shell
	enterKey    rune
	tasks       *TaskManager
	messageChan chan iMessage
	timersChan  chan *adaptiveticker.TimerHandler
	prompt      string
	autosave    bool
}

func NewContext(ticker *adaptiveticker.AdaptiveTicker, reader io.Reader, writer io.Writer, auth interfaces.IAuthenticator, factory *terminal.EquipmentFactory, commands *cli.Command, prompt string, autosave bool) *Context {
	ctx := &Context{
		ticker:      ticker,
		reader:      reader,
		writer:      writer,
		auth:        auth,
		factory:     factory,
		commands:    commands,
		Exit:        false,
		prompt:      prompt,
		enterKey:    -1,
		messageChan: make(chan iMessage, contextMaQueueLen),
		timersChan:  make(chan *adaptiveticker.TimerHandler, contextMaQueueLen),
		tasks:       nil,
		autosave:    autosave,
	}
	return ctx
}

func (c *Context) Setup() {
	c.terminal = c.factory.Create("VT100", c.writer)
	c.terminal.SetKeyFunc(c.keyHandler)
	if c.enterKey > -1 {
		c.terminal.SetEnterKey(c.enterKey)
	}
	system := apps.NewRoot()
	systemCommands, commands := system.Build(c.commands)
	c.tasks = NewTaskManager(c, c.ticker, c.timersChan, []interfaces.ICommand{systemCommands}, commands)

	c.defaultApp = shell.NewShell(c.auth, c.terminal, c.prompt, c.autosave)
	c.defaultApp.ExecCommand = c.execCommand
	c.defaultApp.ExecSuggestion = c.execSuggestion
}

func (c *Context) GetWriter() io.Writer {
	return c.writer
}

func (c *Context) SetScreenSize(width int, height int) {
	c.terminal.SetSize(width, height)
	c.tasks.SetScreenSize(width, height)
}

func (c *Context) keyHandler(event *interfaces.KeyData) {
	if event.GetType() == interfaces.KeyTypeCtrl {
		c.ctrlPressed(event.Key)
		return
	}

	if fgPid := c.tasks.GetForegroundPid(); fgPid != adaptiveticker.UnknownId {
		c.tasks.ExecRead(fgPid, int(event.GetType()), event.Key)
		return
	}

	quit := c.defaultApp.KeyEvent(event)
	if quit {
		c.Exit = true
	}
}

func (c *Context) SetEnterKey(key rune) {
	c.enterKey = key
}

func (c *Context) Close() {
}

func (c *Context) Exec() {
	d := make(chan bool)
	go func() {
		d <- true
		readBuffer := make([]byte, 1024)
		for {
			n, err := c.reader.Read(readBuffer)
			if err == nil {
				if n > 0 {
					re := newMessageRead(readBuffer, n)
					re.postEvent(c.messageChan)
				}
			} else {
				qe := newMessageQuit()
				qe.postEvent(c.messageChan)
				return
			}
		}
	}()
	_ = <-d
	c.eventLoop()
}

func (c *Context) execCommand(line string) (bool, error) {
	return c.tasks.Execute(line, nil)
}

func (c *Context) ctrlPressed(key rune) {
	switch key {
	case 3:
		c.tasks.SetSelectionDisabled()
		c.tasks.KillForeground()
		c.defaultApp.DoNext()
	case 4:
		c.tasks.ExecActivate()
	}
}

func (c *Context) execSuggestion(in string, count int) (int, bool) {
	ret := false
	data, suggestions, found := c.tasks.GetSuggestion(in)
	sLen := 0
	if found && len(suggestions) > 0 {
		sLen = len(suggestions)
		if idx := count % sLen; idx < sLen {
			if complete := suggestions[idx]; len(complete) > len(data) {
				tabLine := complete
				c.defaultApp.DoRedraw(tabLine)
				c.defaultApp.SetHistoryDefault(tabLine)
				ret = true
			}
		}
	}
	return sLen, ret
}

func (c *Context) eventLoop() {
	_, _ = c.terminal.WriteColor("Admin Console Ready", interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)
	c.defaultApp.DoNext()
	for {
		select {
		case m := <-c.messageChan:
			c.messageEventHandler(m)
		case t := <-c.timersChan:
			c.messageEventHandler(t.Event.(iMessage))
		}
		if c.Exit {
			c.shutdown()
			return
		}
	}
}

func (c *Context) messageEventHandler(m iMessage) {
	if m != nil {
		switch m.getType() {
		case MessageTypeRead:
			if mm, ok := m.(*MessageRead); ok {
				c.terminal.Scan(mm.data)
			}

		case MessageTypeTimer:
			if mt, ok := m.(*MessageTimer); ok {
				c.tasks.ExecTimer(mt.pid, mt.tid, mt.interval)
			}

		case MessageTypePaint:
			if _, ok := m.(*MessagePaint); ok {
				c.tasks.ExecPaint(c.terminal)
			}

		case MessageTypeQuit:
			if _, ok := m.(*MessageQuit); ok {
				c.Exit = true
			}
		}
	}
}

func (c *Context) shutdown() {
	c.tasks.KillAll("")
}

//CLI INTERFACE

func (c *Context) Write(data string) {
	_, _ = c.terminal.Write(data)
}

func (c *Context) WriteLn(data string) {
	c.Write(data)
	c.Write(eol)
}

func (c *Context) WriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	_, _ = c.terminal.WriteColor(data, fg, bg, mode)
}

func (c *Context) WriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	c.WriteColor(data, fg, bg, mode)
	c.Write(eol)
}

func (c *Context) ClearScreen() {
	_, _ = c.terminal.ClearScreen()
}

func (c *Context) SetExit() {
	c.Exit = true
}

func (c *Context) History(verb interfaces.HistoryAction, idx int) {
	switch verb {
	case interfaces.HistoryActionClear:
		c.defaultApp.ClearHistory()
	case interfaces.HistoryActionExec:
		if arg, found := c.defaultApp.GetHistoryAtPos(idx); found {
			c.execCommand(arg)
		}
	case interfaces.HistoryActionList:
		_, _ = c.terminal.Write(c.defaultApp.GetHistory())
	}
}

package core

import (
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/apps"
	"github.com/markel1974/c64emu/src/kernel/file_system"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/render"
	"github.com/markel1974/c64emu/src/kernel/shell"
	"io"
)

// contextMaQueueLen defines the default length for message and timer channels used within the context's queue system.
const (
	contextMaQueueLen = 1024
)

// Context is a container for managing and coordinating various dependencies and resources needed for application execution.
// It encapsulates input/output handlers, commands, shell interactions, rendering, authentication, and runtime configurations.
type Context struct {
	ticker   *adaptiveticker.AdaptiveTicker
	reader   io.Reader
	writer   io.Writer
	commands *shell.Command
	render   interfaces.IRender
	auth     interfaces.IAuthenticator
	sh       *shell.Shell
	enterKey rune
	kernel   *Kernel
	prompt   string
	autosave bool
}

// NewContext creates and initializes a new Context with the provided parameters, including ticker, reader, writer, and others.
func NewContext(ticker *adaptiveticker.AdaptiveTicker, reader io.Reader, writer io.Writer, auth interfaces.IAuthenticator, commands *shell.Command, prompt string, autosave bool) *Context {
	ctx := &Context{
		ticker:   ticker,
		reader:   reader,
		writer:   writer,
		auth:     auth,
		commands: commands,
		prompt:   prompt,
		kernel:   nil,
		autosave: autosave,
	}
	return ctx
}

// Setup initializes the context with the terminal, rendering, system commands, file system, kernel, and shell instances.
func (c *Context) Setup(terminal interfaces.ITerminal) {
	c.render = render.NewRender(terminal)
	system := apps.NewRoot()
	systemCommands, commands := system.Build(c.commands)
	fs := file_system.NewCommandInteractor(commands, []interfaces.ICommand{systemCommands})
	ioAdapter := interfaces.IInputOutput(c)
	c.kernel = NewKernel(c.ticker, c.render, ioAdapter, fs)
	c.sh = shell.NewShell(c.auth, c.render, c, c.prompt, c.autosave)
}

// Exec initializes the admin console display, advances the shell line, and starts the kernel.
func (c *Context) Exec() {
	c.render.WriteColor("Admin Console Ready", interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)
	c.sh.NextLine(true)
	c.kernel.Start()
}

// Type processes a key event based on its type and value, influencing execution, interaction, or system state.
func (c *Context) Type(kind interfaces.KeyType, key rune) {
	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			c.kernel.SetSelectionDisabled()
			c.kernel.KillForeground()
			c.sh.NextLine(true)
		case 4:
			c.kernel.ExecActivate()
		}
		return
	}
	if fgPid := c.kernel.GetForegroundPid(); fgPid != adaptiveticker.UnknownId {
		c.kernel.ExecRead(fgPid, int(kind), key)
		return
	}
	if quit := c.sh.KeyEvent(kind, key); quit {
		c.kernel.ExitRequested()
	}
}

// Read reads data into the provided byte slice from the underlying reader and returns the number of bytes read and any error.
func (c *Context) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

// Write writes the provided byte slice to the underlying writer in the context.
// Returns the number of bytes written and any write error encountered.
func (c *Context) Write(data []byte) (int, error) {
	return c.writer.Write(data)
}

// Close releases any resources held by the Context and ensures a clean termination of its operations.
func (c *Context) Close() {
}

// ExecCommand executes a provided command line within the current kernel context and returns the success status and an error.
func (c *Context) ExecCommand(line string) (bool, error) {
	return c.kernel.ExecCommand(line, nil)
}

// ExecSuggestion executes a suggestion mechanism based on input, cursor position, and count, returning total suggestions and success status.
func (c *Context) ExecSuggestion(in string, cursor int, count int) (int, bool) {
	ret := false
	data, suggestions, found := c.kernel.GetSuggestion(in, cursor)
	sLen := 0
	if found && len(suggestions) > 0 {
		sLen = len(suggestions)
		if idx := count % sLen; idx < sLen {
			if complete := suggestions[idx]; len(complete) > len(data) {
				tabLine := complete
				c.sh.Redraw(tabLine)
				c.sh.SetHistoryDefault(tabLine)
				ret = true
			}
		}
	}
	return sLen, ret
}

// SetScreenSize adjusts the terminal's display dimensions to the specified width (w) and height (h).
func (c *Context) SetScreenSize(w int, h int) {
	c.kernel.SetScreenSize(w, h)
}

// History performs actions on the command history based on the specified verb (list, clear, or execute at the given index).
func (c *Context) History(verb interfaces.HistoryAction, idx int) {
	switch verb {
	case interfaces.HistoryActionClear:
		c.sh.ClearHistory()
	case interfaces.HistoryActionExec:
		if arg, found := c.sh.GetHistoryAtPos(idx); found {
			_, _ = c.ExecCommand(arg)
		}
	case interfaces.HistoryActionList:
		c.render.Write(c.sh.GetHistory())
	}
}

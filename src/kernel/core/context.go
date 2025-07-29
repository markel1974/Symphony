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
	auth     interfaces.IAuthenticator
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
	terminalRender := render.NewRender(terminal)
	system := apps.NewRoot()
	systemCommands, commands := system.Build(c.commands)
	fs := file_system.NewCommandInteractor(commands, []interfaces.ICommand{systemCommands})
	ioAdapter := interfaces.IInputOutput(c)
	sh := shell.NewShell(c.auth, terminalRender, c, c.prompt, c.autosave)
	c.kernel = NewKernel(c.ticker, ioAdapter, fs, sh)
}

// Exec initializes the admin console display, advances the shell line, and starts the kernel.
func (c *Context) Exec() {
	c.kernel.Start()
}

// Type processes a key event based on its type and value, influencing execution, interaction, or system state.
func (c *Context) Type(kind interfaces.KeyType, key rune) {
	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			c.kernel.SetSelectionDisabled()
			c.kernel.KillForeground()
			c.kernel.NextLine(true)
		case 4:
			c.kernel.ExecActivate()
		}
		return
	}
	if fgPid := c.kernel.GetForegroundPid(); fgPid != adaptiveticker.UnknownId {
		c.kernel.ExecRead(fgPid, int(kind), key)
		return
	}
	if quit := c.kernel.KeyEvent(kind, key); quit {
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

// SetScreenSize adjusts the terminal's display dimensions to the specified width (w) and height (h).
func (c *Context) SetScreenSize(w int, h int) {
	c.kernel.SetScreenSize(w, h)
}

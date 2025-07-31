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

// NewContext creates and initializes a new Context with the provided parameters, including ticker, inputDriver, outputDriver, and others.
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
	system := apps.NewRoot()
	systemCommands, commands := system.Build(c.commands)
	terminalRender := render.NewRender(terminal)
	fs := file_system.NewCommandInteractor(commands, []interfaces.ICommand{systemCommands})
	sh := shell.NewShell(c.auth, terminalRender, c.prompt, c.autosave)
	c.kernel = NewKernel(c.ticker, c.reader, c.writer, terminalRender, fs, sh)

	terminal.SetIO(c.kernel)
}

// Exec initializes the admin console display, advances the shell line, and starts the kernel.
func (c *Context) Exec() {
	c.kernel.Start()
}

// Close releases any resources held by the Context and ensures a clean termination of its operations.
func (c *Context) Close() {
}

// SetScreenSize adjusts the terminal's display dimensions to the specified width (w) and height (h).
func (c *Context) SetScreenSize(w int, h int) {
	c.kernel.SetScreenSize(w, h)
}

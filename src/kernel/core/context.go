package core

import (
	"io"

	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/apps"
	"github.com/markel1974/c64emu/src/kernel/drivers"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"github.com/markel1974/c64emu/src/kernel/servers/file_system"
	"github.com/markel1974/c64emu/src/kernel/servers/render"
)

// contextMaQueueLen defines the default length for message and timer channels used within the context's queue system.
const (
	contextMaQueueLen = 1024
)

// Context is a container for managing and coordinating various dependencies and resources needed for application execution.
// It encapsulates input/output routingTable, commands, shell interactions, rendering, authentication, and runtime configurations.
type Context struct {
	ticker   *adaptiveticker.AdaptiveTicker
	reader   io.Reader
	writer   io.Writer
	commands *process.Command
	auth     interfaces.IAuthenticator
	enterKey rune
	kernel   *Kernel
	prompt   string
	autosave bool
}

// NewContext creates and initializes a new Context with the provided parameters, including ticker, inputDriver, outputDriver, and others.
func NewContext(ticker *adaptiveticker.AdaptiveTicker, reader io.Reader, writer io.Writer, auth interfaces.IAuthenticator, commands *process.Command, prompt string, autosave bool) *Context {
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
	xsh, systemCommands, commands := system.Build(c.commands)

	keyboardDriver := drivers.NewKeyboardTerminal(c.reader, terminal)
	videoDriver := drivers.NewDisplayTerminal(c.writer, terminal)

	admin := "root"

	terminalRender := render.NewRender(admin, videoDriver)
	fs := file_system.NewFileSystem(admin, commands, []interfaces.ICommand{systemCommands})
	c.kernel = NewKernel(admin, c.ticker, keyboardDriver, xsh)
	c.kernel.AddServer(terminalRender)
	c.kernel.AddServer(fs)
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

package render

import (
	"bytes"
	"log"
	"sort"

	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

const (
	renderQueueLen = 1024
	renderQueueMax = renderQueueLen - 1
)

// eol represents the end-of-line marker used for denoting line breaks in the output, set to "\r\n".
const eolDef = "\r\n"

// Render represents a rendering engine responsible for managing terminal dimensions, repainting logic, and paint tasks.
type Render struct {
	process        interfaces.IProcess
	driver         interfaces.IDisplayDriver
	surface        *Surface
	width          int
	height         int
	fullPaint      bool
	windowSelector *WindowSelector
	running        map[int]*Component
	foreground     *Component
	messageChan    chan interfaces.IMessage
	router         interfaces.IRouter
	handlers       map[interfaces.MessageType]func(interfaces.IMessage)
}

// NewRender creates and initializes a new Render instance with the provided terminal implementation.
// Returns a pointer to the newly created Render object.
func NewRender(driver interfaces.IDisplayDriver) *Render {
	const width = 80
	const height = 24
	r := &Render{
		driver:         driver,
		windowSelector: NewWindowSelector(),
		width:          width,
		height:         height,
		fullPaint:      true,
		running:        make(map[int]*Component),
		messageChan:    make(chan interfaces.IMessage, renderQueueLen),
		surface:        NewSurface(driver, height, width, ""),
		handlers:       make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}
	r.handlers[interfaces.MessageTypeSetScreenSize] = r.handleSetScreenSize
	r.handlers[interfaces.MessageTypePaintRequest] = r.handlePaintRequest
	r.handlers[interfaces.MessageTypePaintApply] = r.handlePaintApply
	r.handlers[interfaces.MessageTypeWindowsSelectionBegin] = r.handleWindowsSelectionBegin
	r.handlers[interfaces.MessageTypeWindowsSelectionOptions] = r.handleWindowsSelectionOptions
	r.handlers[interfaces.MessageTypeWindowsSelectionPrevious] = r.handleWindowsSelectionPrevious
	r.handlers[interfaces.MessageTypeWindowsSelectionNext] = r.handleWindowsSelectionNext
	r.handlers[interfaces.MessageTypeWindowsSelectionEnd] = r.handleWindowsSelectionEnd
	r.handlers[interfaces.MessageTypeClearLine] = r.handleClearLine
	r.handlers[interfaces.MessageTypeClearScreen] = r.handleClearScreen
	r.handlers[interfaces.MessageTypeSaveCursor] = r.handleSaveCursor
	r.handlers[interfaces.MessageTypeRestoreCursor] = r.handleRestoreCursor
	r.handlers[interfaces.MessageTypeMoveCursorLeft] = r.handleMoveCursorLeft
	r.handlers[interfaces.MessageTypeMoveCursorRight] = r.handleMoveCursorRight
	r.handlers[interfaces.MessageTypeWrite] = r.handleWrite
	r.handlers[interfaces.MessageTypeWriteColor] = r.handleWriteColor
	r.handlers[interfaces.MessageTypeGetScreenSize] = r.handleGetScreenSize
	r.handlers[interfaces.MessageTypeNotifyProcessCreate] = r.handleProcessCreate
	r.handlers[interfaces.MessageTypeNotifyProcessForeground] = r.handleProcessForeground
	r.handlers[interfaces.MessageTypeNotifyProcessTerminate] = r.handleProcessTerminate
	return r
}

// Name returns the name of the Render object as a string.
func (c *Render) Name() string {
	return "render"
}

// Process returns the process implementation adhering to the interfaces.IProcess interface.
func (c *Render) Process() interfaces.IProcess {
	return c.process
}

// PID returns an identifier for the file system process. It always returns a fixed value of -2.
func (c *Render) PID() int {
	return c.process.PID()
}

// User returns the default username as a string, typically "root".
func (c *Render) User() string {
	return c.process.User()
}

// Register binds the provided router to the Render instance and returns a list of supported message types.
func (c *Render) Register(router interfaces.IRouter) []interfaces.MessageType {
	c.router = router
	var out []interfaces.MessageType
	for id := range c.handlers {
		out = append(out, id)
	}
	return out
}

func (c *Render) Setup(process interfaces.IProcess) error {
	c.process = process
	b := make(chan bool)
	c.eventLoop(b)
	_ = <-b
	return nil
}

// PostMessage sends a message of type IMessage to the file system's message channel for further processing.
func (c *Render) PostMessage(m interfaces.IMessage) {
	if len(c.messageChan) >= renderQueueMax {
		log.Printf("Render: message queue full, dropping message: %d", m.GetType())
		return
	}
	c.messageChan <- m
}

// NotifyProcessCreation notifies the Render instance about the creation of a new process and updates internal state if necessary.
func (c *Render) handleProcessCreate(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageNotifyProcessCreate)
	if !ok {
		return
	}
	c.running[mt.CreatedPID()] = NewComponent(mt.CreatedPID(), mt.Name(), c.driver, c.height, c.width)
}

// NotifyProcessTermination handles the necessary cleanup and state updates when a process associated with the Render terminates.
func (c *Render) handleProcessTerminate(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageNotifyProcessTerminate)
	if !ok {
		return
	}
	c.windowSelector.Clear()
	delete(c.running, mt.TerminatedPID())
}

// NotifyProcessForeground updates the Render object with the process description currently in the foreground.
func (c *Render) handleProcessForeground(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageNotifyProcessForeground)
	if !ok {
		return
	}
	p, _ := c.running[mt.ForegroundPID()]
	if p == nil {
		c.foreground = nil
		return
	}
	c.foreground = p
}

// eventLoop continuously listens on the message channel and processes incoming messages until a quit message is received.
func (c *Render) eventLoop(r chan bool) {
	go func() {
		r <- true
		for {
			select {
			case m, ok := <-c.messageChan:
				if !ok {
					return
				}
				id := m.GetType()
				if id == interfaces.MessageTypeQuit {
					close(c.messageChan)
					return
				}
				if handler, _ := c.handlers[id]; handler != nil {
					handler(m)
				} else {
					log.Printf("Render: unknown message type: %d", id)
				}
			}
		}
	}()
}

// handleSetScreenSize handles setting the screen size by updating width, height, and triggering a full repaint.
func (c *Render) handleSetScreenSize(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageSetScreenSize)
	if !ok {
		return
	}
	c.width = mt.Width()
	c.height = mt.Height()
	c.fullPaint = true
}

// handlePaintRequest handles paint requests by triggering a repaint.
func (c *Render) handlePaintRequest(msg interfaces.IMessage) {
	mp := messages.NewMessagePaintPrepare(msg.PID(), NewInterpretedSurface(c.height, c.width))
	c.router.PostMessage(mp)
}

// handlePaintApply handles paint apply messages by applying the surface to the terminal and triggering a repaint.
func (c *Render) handlePaintApply(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessagePaintApply)
	if !ok {
		return
	}
	component, _ := c.running[msg.PID()]
	if component == nil {
		return
	}
	component.SetInterpretedSurface(mt.Surface())

	fullPaint := c.fullPaint
	c.fullPaint = false

	var surfaces []*Surface
	for _, process := range c.running {
		surface := process.Surface()
		if process.PID() == c.windowSelector.PID() {
			surface.SetZIndex(255)
			surface.SetSelectionMode(true)
		} else {
			surface.SetZIndex(0)
			surface.SetSelectionMode(false)
		}
		surfaces = append(surfaces, surface)
	}

	sort.SliceStable(surfaces, func(i, j int) bool {
		return surfaces[i].ZIndex() < surfaces[j].ZIndex()
	})

	var lines bytes.Buffer
	c.surface.Prepare(c.height, c.width)

	for _, surface := range surfaces {
		if surface.interpreted != nil {
			c.surface.Assign(surface)
			c.surface.Begin()
			surface.GetInterpretedSurface().Appy(c.surface)
			c.surface.End()
		}
	}
	c.surface.GetBuffer(&lines, fullPaint)

	c.doSaveCursor()
	c.doMoveCursorTopLeft()
	c.doWrite(string(lines.Bytes()), false)
	c.doRestoreCursor()
}

// handleWindowsSelectionBegin handles the selection of a process to be displayed in the terminal.
func (c *Render) handleWindowsSelectionBegin(msg interfaces.IMessage) {
	c.windowSelector.Clear()
	for idx, process := range c.running {
		c.windowSelector.AddAvailable(process.PID())
		if c.windowSelector.PID() == adaptiveticker.UnknownId {
			if c.foreground != nil {
				if c.foreground.PID() == process.PID() {
					c.windowSelector.Set(c.foreground.PID(), idx)
				}
			} else {
				c.windowSelector.Set(process.PID(), idx)
			}
		}
	}
	mp := messages.NewMessagePaintPrepare(msg.PID(), NewInterpretedSurface(c.height, c.width))
	c.router.PostMessage(mp)
}

// handleWindowsSelectionOptions processes a windows selection options message and updates the corresponding surface options.
func (c *Render) handleWindowsSelectionOptions(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageWindowsSelectionOptions)
	if !ok {
		return
	}
	pid := c.windowSelector.PID()
	process, _ := c.running[pid]
	if process == nil {
		return
	}
	process.Surface().SetOption(mt.Option(), mt.Value())
	c.fullPaint = true
	mp := messages.NewMessagePaintPrepare(msg.PID(), NewInterpretedSurface(c.height, c.width))
	c.router.PostMessage(mp)
}

// handleWindowsSelectionPrevious navigates to the previous window in the selection if possible and triggers a paint request.
func (c *Render) handleWindowsSelectionPrevious(msg interfaces.IMessage) {
	if c.windowSelector.Prev() {
		mp := messages.NewMessagePaintPrepare(msg.PID(), NewInterpretedSurface(c.height, c.width))
		c.router.PostMessage(mp)
	}
}

// handleWindowsSelectionNext moves the window selector to the next window and triggers a repaint if the selection changes.
func (c *Render) handleWindowsSelectionNext(msg interfaces.IMessage) {
	if c.windowSelector.Next() {
		mp := messages.NewMessagePaintPrepare(msg.PID(), NewInterpretedSurface(c.height, c.width))
		c.router.PostMessage(mp)
	}
}

// handleWindowsSelectionEnd handles the completion of the window selection process by clearing the window selector state.
func (c *Render) handleWindowsSelectionEnd(_ interfaces.IMessage) {
	c.windowSelector.Clear()
}

// handleClearLine processes a clear line message by invoking the doClearLine method if the message is of correct type.
func (c *Render) handleClearLine(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageClearLine)
	if !ok {
		return
	}
	c.doClearLine(mt.Line())
}

// handleClearScreen processes a clear screen message and invokes the method to clear the screen if the message is valid.
func (c *Render) handleClearScreen(msg interfaces.IMessage) {
	_, ok := msg.(*messages.MessageClearScreen)
	if !ok {
		return
	}
	c.doClearScreen()
}

// handleSaveCursor processes a save cursor message and triggers the logic to save the current cursor state if valid.
func (c *Render) handleSaveCursor(msg interfaces.IMessage) {
	_, ok := msg.(*messages.MessageSaveCursor)
	if !ok {
		return
	}
	c.doSaveCursor()
}

// handleRestoreCursor processes a restore cursor message and updates the rendering state if the message is valid.
func (c *Render) handleRestoreCursor(msg interfaces.IMessage) {
	_, ok := msg.(*messages.MessageRestoreCursor)
	if !ok {
		return
	}
	c.doRestoreCursor()
}

// handleMoveCursorLeft moves the cursor one position to the left using the associated driver commands.
func (c *Render) handleMoveCursorLeft(msg interfaces.IMessage) {
	_, ok := msg.(*messages.MessageMoveCursorLeft)
	if !ok {
		return
	}
	c.doMoveCursorLeft()
}

// handleMoveCursorRight moves the cursor one step to the right on the rendering interface using the associated driver.
func (c *Render) handleMoveCursorRight(msg interfaces.IMessage) {
	_, ok := msg.(*messages.MessageMoveCursorRight)
	if !ok {
		return
	}
	c.doMoveCursorRight()
}

// handleWrite writes the provided string to the terminal followed by an end-of-line character.
func (c *Render) handleWrite(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageWrite)
	if !ok {
		return
	}
	c.doWrite(mt.Data(), mt.Eol())
}

// CallGetScreenSize returns the current screen width and height of the Render instance.
func (c *Render) handleGetScreenSize(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageGetScreenSize)
	if !ok {
		return
	}
	mt.SetResult(c.width, c.height)
	c.router.PostMessage(mt)
}

// handleWriteColor writes the provided string to the terminal using the provided foreground and background colors.
func (c *Render) handleWriteColor(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessageWriteColor)
	if !ok {
		return
	}
	p := c.driver.CreateColorize(mt.Data(), int(mt.Fg()), int(mt.Bg()), mt.Mode())
	_, _ = c.driver.Write([]byte(p))
	if mt.Eol() {
		_, _ = c.driver.Write([]byte(eolDef))
	}
}

// doClearLine clears the specified line by generating a clear line sequence and writing it through the driver.
func (c *Render) doClearLine(line string) {
	p := c.driver.CreateClearLine(line)
	_, _ = c.driver.Write(p)
}

// doMoveCursorTopLeft moves the cursor to the top-left position of the rendering surface using the underlying driver.
func (c *Render) doMoveCursorTopLeft() {
	p := c.driver.CreateMoveCursorTopLeft()
	_, _ = c.driver.Write(p)
}

// doClearScreen clears the terminal screen by generating and writing a clear screen sequence via the driver.
func (c *Render) doClearScreen() {
	p := c.driver.CreateClearScreen()
	_, _ = c.driver.Write(p)
}

// doSaveCursor saves the current cursor position using the driver's CreateSaveCursor method and writes it to the output.
func (c *Render) doSaveCursor() {
	p := c.driver.CreateSaveCursor()
	_, _ = c.driver.Write(p)
}

// doRestoreCursor restores the cursor to its last saved position using the underlying driver's capabilities.
func (c *Render) doRestoreCursor() {
	p := c.driver.CreateRestoreCursor()
	_, _ = c.driver.Write(p)
}

// doMoveCursorLeft moves the cursor one position to the left using the associated driver commands.
func (c *Render) doMoveCursorLeft() {
	p := c.driver.CreateMoveCursorLeft()
	_, _ = c.driver.Write(p)
}

// doMoveCursorRight moves the cursor one step to the right on the rendering interface using the associated driver.
func (c *Render) doMoveCursorRight() {
	p := c.driver.CreateMoveCursorRight()
	_, _ = c.driver.Write(p)
}

// doWrite writes the provided data string to the driver and optionally appends an end-of-line delimiter if eol is true.
func (c *Render) doWrite(data string, eol bool) {
	_, _ = c.driver.Write([]byte(data))
	if eol {
		_, _ = c.driver.Write([]byte(eolDef))
	}
}

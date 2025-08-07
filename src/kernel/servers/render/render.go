package render

import (
	"bytes"
	"log"
	"sort"

	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
)

// eol represents the end-of-line marker used for denoting line breaks in the output, set to "\r\n".
const eolDef = "\r\n"

// Render represents a rendering engine responsible for managing terminal dimensions, repainting logic, and paint tasks.
type Render struct {
	driver         interfaces.IDisplayDriver
	user           string
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
func NewRender(user string, driver interfaces.IDisplayDriver) *Render {
	const width = 80
	const height = 24
	r := &Render{
		driver:         driver,
		user:           user,
		windowSelector: NewWindowSelector(),
		width:          width,
		height:         height,
		fullPaint:      true,
		running:        make(map[int]*Component),
		messageChan:    make(chan interfaces.IMessage, 128),
		surface:        NewSurface(driver, height, width, ""),
		handlers:       make(map[interfaces.MessageType]func(interfaces.IMessage)),
	}
	r.handlers[interfaces.MessageTypePaintRequest] = r.handlePaintRequest
	r.handlers[interfaces.MessageTypePaintApply] = r.handlePaintApply
	r.handlers[interfaces.MessageTypeWindowsSelectionBegin] = r.handleWindowsSelectionBegin
	r.handlers[interfaces.MessageTypeWindowsSelectionOptions] = r.handleWindowsSelectionOptions
	r.handlers[interfaces.MessageTypeWindowsSelectionPrevious] = r.handleWindowsSelectionPrevious
	r.handlers[interfaces.MessageTypeWindowsSelectionNext] = r.handleWindowsSelectionNext
	r.handlers[interfaces.MessageTypeWindowsSelectionEnd] = r.handleWindowsSelectionEnd
	return r
}

// Process returns a nil pointer to the Render object.
func (c *Render) Process() interfaces.IProcess {
	return nil
}

// PID returns a hardcoded integer value representing a Process ID.
func (c *Render) PID() int {
	return -1
}

// User returns the default username as a string, typically "root".
func (c *Render) User() string {
	return c.user
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

// Start begins the process by setting its state to running and initiating its event loop asynchronously.
func (c *Render) Start() {
	b := make(chan bool)
	c.eventLoop(b)
	_ = <-b
}

// PostMessage sends a message of type IMessage to the file system's message channel for further processing.
func (c *Render) PostMessage(m interfaces.IMessage) {
	c.messageChan <- m
}

// CallGetScreenSize returns the current screen width and height of the Render instance.
func (c *Render) CallGetScreenSize(router interfaces.IRouter) (int, int) {
	return c.width, c.height
}

// CallClearLine clears the specified line from the terminal screen using the underlying ITerminal implementation.
func (c *Render) CallClearLine(router interfaces.IRouter, line string) {
	c.doClearLine(line)
}

// CallSetScreenSize updates the screen's width and height, marks the screen for a full repaint, and sets the terminal size.
func (c *Render) CallSetScreenSize(router interfaces.IRouter, width int, height int) {
	c.width = width
	c.height = height
	c.fullPaint = true
}

// CallClearScreen clears the terminal screen using the underlying ITerminal implementation.
func (c *Render) CallClearScreen(router interfaces.IRouter) {
	p := c.driver.CreateClearScreen()
	_, _ = c.driver.Write(p)
}

// CallSaveCursor saves the current cursor position in the terminal for future restoration.
func (c *Render) CallSaveCursor(router interfaces.IRouter) {
	p := c.driver.CreateSaveCursor()
	_, _ = c.driver.Write(p)
}

// CallRestoreCursor restores the saved cursor position in the terminal using the associated ITerminal implementation.
func (c *Render) CallRestoreCursor(router interfaces.IRouter) {
	p := c.driver.CreateRestoreCursor()
	_, _ = c.driver.Write(p)
}

// CallMoveCursorLeft moves the terminal cursor one position to the left using the underlying terminal implementation.
func (c *Render) CallMoveCursorLeft(router interfaces.IRouter) {
	p := c.driver.CreateMoveCursorLeft()
	_, _ = c.driver.Write(p)
}

// CallMoveCursorRight moves the cursor one position to the right in the terminal.
func (c *Render) CallMoveCursorRight(router interfaces.IRouter) {
	p := c.driver.CreateMoveCursorRight()
	_, _ = c.driver.Write(p)
}

// CallWrite writes the provided string to the terminal followed by an end-of-line character.
func (c *Render) CallWrite(router interfaces.IRouter, data string, eol bool) {
	_, _ = c.driver.Write([]byte(data))
	if eol {
		_, _ = c.driver.Write([]byte(eolDef))
	}
}

// CallWriteColor writes the given text with specified foreground and background colors and mode, followed by a line break.
func (c *Render) CallWriteColor(router interfaces.IRouter, data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode, eol bool) {
	p := c.driver.CreateColorize(data, int(fg), int(bg), mode)
	_, _ = c.driver.Write([]byte(p))
	if eol {
		_, _ = c.driver.Write([]byte(eolDef))
	}
}

// NotifyProcessCreation notifies the Render instance about the creation of a new process and updates internal state if necessary.
func (c *Render) NotifyProcessCreation(desc *interfaces.ProcessDescription) {
	c.running[desc.PID()] = NewComponent(desc.PID(), desc.Name(), c.driver, c.height, c.width)
}

// NotifyProcessTermination handles the necessary cleanup and state updates when a process associated with the Render terminates.
func (c *Render) NotifyProcessTermination(desc *interfaces.ProcessDescription) {
	c.windowSelector.Clear()
	delete(c.running, desc.PID())
}

// NotifyProcessForeground updates the Render object with the process description currently in the foreground.
func (c *Render) NotifyProcessForeground(desc *interfaces.ProcessDescription) {
	p, _ := c.running[desc.PID()]
	if p == nil {
		c.foreground = nil
		return
	}
	c.foreground = p
}

// evenLoop continuously listens on the message channel and processes incoming messages until a quit message is received.
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

// handlePaintRequest handles paint requests by triggering a repaint.
func (c *Render) handlePaintRequest(msg interfaces.IMessage) {
	mp := messages.NewMessagePaintPrepare(msg.Router(), NewDescriptiveSurface(c.height, c.width))
	c.router.PostMessage(mp)
}

func (c *Render) handlePaintApply(msg interfaces.IMessage) {
	mt, ok := msg.(*messages.MessagePaintApply)
	if !ok {
		return
	}
	pid := msg.Router().PID()
	process, _ := c.running[pid]
	if process == nil {
		return
	}
	process.SetDescriptiveSurface(mt.Surface())
	c.doPaint(mt.Router(), false)
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
	mp := messages.NewMessagePaintPrepare(msg.Router(), NewDescriptiveSurface(c.height, c.width))
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
	mp := messages.NewMessagePaintPrepare(msg.Router(), NewDescriptiveSurface(c.height, c.width))
	c.router.PostMessage(mp)
}

// handleWindowsSelectionPrevious navigates to the previous window in the selection if possible and triggers a paint request.
func (c *Render) handleWindowsSelectionPrevious(msg interfaces.IMessage) {
	if c.windowSelector.Prev() {
		mp := messages.NewMessagePaintPrepare(msg.Router(), NewDescriptiveSurface(c.height, c.width))
		c.router.PostMessage(mp)
	}
}

// handleWindowsSelectionNext moves the window selector to the next window and triggers a repaint if the selection changes.
func (c *Render) handleWindowsSelectionNext(msg interfaces.IMessage) {
	if c.windowSelector.Next() {
		mp := messages.NewMessagePaintPrepare(msg.Router(), NewDescriptiveSurface(c.height, c.width))
		c.router.PostMessage(mp)
	}
}

// handleWindowsSelectionEnd handles the completion of the window selection process by clearing the window selector state.
func (c *Render) handleWindowsSelectionEnd(msg interfaces.IMessage) {
	c.windowSelector.Clear()
}

// clearLine clears the specified line from the terminal screen using the terminal implementation of the associated Render object.
func (c *Render) doClearLine(line string) {
	p := c.driver.CreateClearLine(line)
	_, _ = c.driver.Write(p)
}

// moveCursorTopLeft moves the terminal cursor to the top-left position using the underlying terminal implementation.
func (c *Render) doMoveCursorTopLeft() {
	p := c.driver.CreateMoveCursorTopLeft()
	_, _ = c.driver.Write(p)
}

// doPaint renders and updates the UI by processing components and their surfaces based on the given router and paint mode.
// Components are sorted by their zIndex for proper rendering order, and final output is sent to the router interface.
func (c *Render) doPaint(router interfaces.IRouter, full bool) {
	if full {
		c.fullPaint = true
	}
	fullPaint := c.fullPaint
	c.fullPaint = false
	rowMax := 0
	var tasks []*Component
	for _, process := range c.running {
		process.Compile(c.height, c.width, c.windowSelector.PID())
		if process.RowMax() > rowMax {
			rowMax = process.RowMax()
		}
		tasks = append(tasks, process)
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].Surface().ZIndex() < tasks[j].Surface().ZIndex()
	})

	var lines bytes.Buffer
	w, h := c.CallGetScreenSize(router)
	c.surface.Prepare(h, w)
	c.surface.SetRowMax(rowMax)
	for _, s := range tasks {
		c.surface.Merge(s.Surface())
	}
	c.surface.GetBuffer(&lines, fullPaint)

	c.CallSaveCursor(router)
	c.doMoveCursorTopLeft()
	c.CallWrite(router, string(lines.Bytes()), false)
	c.CallRestoreCursor(router)
}

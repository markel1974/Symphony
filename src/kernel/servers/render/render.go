package render

import (
	"bytes"
	"github.com/markel1974/c64emu/src/kernel/adaptiveticker"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/messages"
	"log"
)

// eol represents the end-of-line marker used for denoting line breaks in the output, set to "\r\n".
const eolDef = "\r\n"

// Render represents a rendering engine responsible for managing terminal dimensions, repainting logic, and paint tasks.
type Render struct {
	driver         interfaces.IDisplayDriver
	user           string
	surface        *Surface
	dirty          bool
	width          int
	height         int
	fullPaint      bool
	windowSelector *WindowSelector
	running        map[int]*Component
	foreground     *Component
	messageChan    chan interfaces.IMessage
	router         interfaces.IRouter
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
		dirty:          false,
		width:          width,
		height:         height,
		fullPaint:      true,
		running:        make(map[int]*Component),
		messageChan:    make(chan interfaces.IMessage, 128),
		surface:        NewSurface(driver, height, width),
	}

	return r
}

// PID returns a hardcoded integer value representing a Process ID.
func (c *Render) PID() int {
	return -1
}

// User returns the default username as a string, typically "root".
func (c *Render) User() string {
	return c.user
}

// SetRouter sets the instance of IRouter to be used by the Render for routing purposes.
func (c *Render) SetRouter(router interfaces.IRouter) {
	c.router = router
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

// Register returns a slice of message types that the Render object is set to handle, including MessageTypePaint.
func (c *Render) Register() []interfaces.MessageType {
	return []interfaces.MessageType{interfaces.MessageTypePaint, interfaces.MessageTypePaintRequest}
}

// CallGetScreenSize returns the current screen width and height of the Render instance.
func (c *Render) CallGetScreenSize() (int, int) {
	return c.width, c.height
}

// CallSetScreenSize updates the screen's width and height, marks the screen for a full repaint, and sets the terminal size.
func (c *Render) CallSetScreenSize(width int, height int) {
	c.width = width
	c.height = height
	c.fullPaint = true
}

// CallWindowsSelectionBegin updates the selection mode for a specific process and triggers a repaint without requesting a redraw.
func (c *Render) CallWindowsSelectionBegin() {
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
	c.handlePaintRequest(false)
}

// CallWindowsSelectionOptions modifies window selection options for a process and triggers a paint request if necessary.
func (c *Render) CallWindowsSelectionOptions(option rune, value float64) {
	process, _ := c.running[c.windowSelector.PID()]
	if process == nil {
		return
	}

	process.SetWindowOption(option, value)
	process.surface.SetWindowOptions(process.WindowOptions)
	c.handlePaintRequest(true)
}

// CallWindowsSelectionPrevious moves the task selection to the previous task and triggers a render update if successful.
func (c *Render) CallWindowsSelectionPrevious() {
	if c.windowSelector.Prev() {
		c.handlePaintRequest(false)
	}
}

// CallWindowsSelectionNext moves the task selection to the next task and triggers a render update if successful.
func (c *Render) CallWindowsSelectionNext() {
	if c.windowSelector.Next() {
		c.handlePaintRequest(false)
	}
}

// CallWindowsSelectionEnd clears the current selection in the windowSelector instance of the Render object.
func (c *Render) CallWindowsSelectionEnd() {
	c.windowSelector.Clear()
}

// CallWrite sends the given string data to the terminal's output stream.
func (c *Render) CallWrite(data string) {
	_, _ = c.driver.Write([]byte(data))
}

// CallWriteLn writes the provided string to the terminal followed by an end-of-line character.
func (c *Render) CallWriteLn(data string) {
	_, _ = c.driver.Write([]byte(data))
	_, _ = c.driver.Write([]byte(eolDef))
}

// CallWriteColor writes the given data string to the terminal with specified foreground and background colors, and color mode.
func (c *Render) CallWriteColor(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	p := c.driver.CreateColorize(data, int(fg), int(bg), mode)
	_, _ = c.driver.Write([]byte(p))
}

// CallWriteColorLn writes the given text with specified foreground and background colors and mode, followed by a line break.
func (c *Render) CallWriteColorLn(data string, fg interfaces.ColorDef, bg interfaces.ColorDef, mode interfaces.ColorMode) {
	p := c.driver.CreateColorize(data, int(fg), int(bg), mode)
	_, _ = c.driver.Write([]byte(p))
	_, _ = c.driver.Write([]byte(eolDef))
}

// CallClearScreen clears the terminal screen using the underlying ITerminal implementation.
func (c *Render) CallClearScreen() {
	p := c.driver.CreateClearScreen()
	_, _ = c.driver.Write(p)
}

// CallSaveCursor saves the current cursor position in the terminal for future restoration.
func (c *Render) CallSaveCursor() {
	p := c.driver.CreateSaveCursor()
	_, _ = c.driver.Write(p)
}

// CallRestoreCursor restores the saved cursor position in the terminal using the associated ITerminal implementation.
func (c *Render) CallRestoreCursor() {
	p := c.driver.CreateRestoreCursor()
	_, _ = c.driver.Write(p)
}

// CallMoveCursorLeft moves the terminal cursor one position to the left using the underlying terminal implementation.
func (c *Render) CallMoveCursorLeft() {
	p := c.driver.CreateMoveCursorLeft()
	_, _ = c.driver.Write(p)
}

// CallMoveCursorRight moves the cursor one position to the right in the terminal.
func (c *Render) CallMoveCursorRight() {
	p := c.driver.CreateMoveCursorRight()
	_, _ = c.driver.Write(p)
}

// CallWritePromptLine clears the given line and writes the prompt and line with specified color and mode configurations.
func (c *Render) CallWritePromptLine(prompt string, line string) {
	c.clearLine(line)
	c.CallWriteColor(prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	c.CallWriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// CallWritePromptEOL writes the provided prompt with green color and optionally appends an end-of-line marker if enabled.
func (c *Render) CallWritePromptEOL(prompt string, eol bool) {
	if eol {
		c.CallWriteColor(eolDef, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
	}
	c.CallWriteColor(prompt, interfaces.ColorGreenDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// CallWriteCritical writes a critical log line with predefined red color and normal mode formatting.
func (c *Render) CallWriteCritical(line string) {
	c.CallWriteColor(line, interfaces.ColorRedDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// CallWriteNormal writes a line with default colors and normal mode configuration using the CallWriteColor method.
func (c *Render) CallWriteNormal(line string) {
	c.CallWriteColor(line, interfaces.ColorNoneDef, interfaces.ColorNoneDef, interfaces.ModeNormal)
}

// CallWriteHighlight writes the given line with default blue foreground, red background, and normal display mode.
func (c *Render) CallWriteHighlight(line string) {
	c.CallWriteColor(line, interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)
}

// clearLine clears the specified line from the terminal screen using the terminal implementation of the associated Render object.
func (c *Render) clearLine(line string) {
	p := c.driver.CreateClearLine(line)
	_, _ = c.driver.Write(p)
}

// moveCursorTopLeft moves the terminal cursor to the top-left position using the underlying terminal implementation.
func (c *Render) moveCursorTopLeft() {
	p := c.driver.CreateMoveCursorTopLeft()
	_, _ = c.driver.Write(p)
}

// NotifyProcessCreation notifies the Render instance about the creation of a new process and updates internal state if necessary.
func (c *Render) NotifyProcessCreation(desc *interfaces.ProcessDescription) {
	if !desc.HasPaint() {
		return
	}
	c.running[desc.PID()] = NewComponent(desc, c.driver)
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
				if m.GetType() == interfaces.MessageTypeQuit {
					close(c.messageChan)
					return
				}
				switch m.GetType() {
				case interfaces.MessageTypePaint:
					c.handlePaintExec()
				case interfaces.MessageTypePaintRequest:
					c.handlePaintRequest(false)
					//
				default:
					log.Printf("Unknown message type: %v\n", m.GetType())
				}
			}
		}
	}()
}

// CallPaintExec executes a rendering operation if the surface is marked as dirty, processing selected and other tasks.
// Returns true if the rendering process is executed, false otherwise.
func (c *Render) handlePaintExec() {
	if !c.dirty {
		return
	}
	var selectedProcess *Component = nil
	var tasks []*Component
	for _, process := range c.running {
		if process.PID() == c.windowSelector.PID() {
			selectedProcess = process
			selectedProcess.surface.SetSelectionMode(true)
		} else {
			process.surface.SetSelectionMode(false)
			tasks = append(tasks, process)
		}
	}
	if selectedProcess != nil {
		//zOrder
		tasks = append(tasks, selectedProcess)
	}

	w, h := c.CallGetScreenSize()
	fullPaint := c.fullPaint
	c.fullPaint = false
	rMax := 0
	//zOrder
	for _, task := range tasks {
		task.surface.Prepare(h, w, fullPaint)
		task.surface.Begin()
		task.Paint(task.surface)
		task.surface.End()
		if task.surface.rMax > rMax {
			rMax = task.surface.rMax
		}
	}

	var lines bytes.Buffer
	c.surface.Prepare(h, w, fullPaint)
	c.surface.rMax = rMax
	for _, s := range tasks {
		c.surface.Merge(s.surface)
	}
	c.surface.GetBuffer(&lines)

	c.CallSaveCursor()
	c.moveCursorTopLeft()
	c.CallWrite(string(lines.Bytes()))
	c.CallRestoreCursor()

	c.dirty = false
}

// handlePaintRequest triggers a paint request by marking the object as dirty and setting up the necessary ticker for repainting.
func (c *Render) handlePaintRequest(full bool) {
	if full {
		c.fullPaint = true
	}
	if !c.dirty {
		c.dirty = true
		msg := messages.NewMessageTimedMessage(messages.NewMessagePaint(), -1, -1, -1)
		c.router.PostMessage(msg)
		//c.router.PostTimedMessage(messages.NewMessagePaint(), -1, -1, 1)
	}
}

/*
func (c *Render) MergeDisplays(tasks []*Component) [][]string {
	// 1. Calcola dimensioni massime
	maxRows := 0
	maxCols := 0
	for _, disp := range tasks {
		if len(disp.surface.surface) > maxRows {
			maxRows = len(disp.surface.surface)
		}
		if len(disp.surface.surface) > 0 && len(disp.surface.surface[0]) > maxCols {
			maxCols = len(disp.surface.surface[0])
		}
	}
	merged := make([][]string, maxRows)
	for i := range merged {
		merged[i] = make([]string, maxCols)
		for j := range merged[i] {
			merged[i][j] = " "
		}
	}
	for _, disp := range tasks {
		for i, row := range disp.surface.surface {
			for j, val := range row {
				if val != "$" {
					merged[i][j] = val
				}
			}
		}
	}
	return merged
}



*/

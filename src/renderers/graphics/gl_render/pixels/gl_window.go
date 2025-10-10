package pixels

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"runtime"
	"time"

	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels/executor"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// WindowConfig represents the configuration options for initializing a window.
// Title specifies the title of the window.
// Icon is a slice of IPicture used for the window's icon images.
// Bounds defines the initial size of the window as a Rect.
// Position specifies the initial position of the window using a Vector.
// Monitor specifies the monitor for full-screen mode, if applicable.
// Smooth determines whether rendering is smoothed.
// Resizable indicates if the window can be resized by the user.
// Undecorated defines whether the window has system decorations (e.g., borders, title bar).
// NoIconify prevents the window from being iconified (minimized).
// AlwaysOnTop makes the window stay above other windows.
// TransparentFramebuffer enables transparency in the framebuffer.
// VSync sets vertical synchronization for reducing screen tearing.
// Maximized determines whether the window starts maximized.
// Invisible makes the window initially invisible.
// SamplesMSAA specifies the number of samples for multisample anti-aliasing (MSAA).
type WindowConfig struct {
	Title string

	Icon []IPicture

	Bounds Rect

	Position Vector

	Monitor *GLMonitor

	Smooth bool

	Resizable bool

	Undecorated bool

	NoIconify bool

	AlwaysOnTop bool

	TransparentFramebuffer bool

	VSync bool

	Maximized bool

	Invisible bool

	SamplesMSAA int
}

// GLWindow represents a graphical window using GLFW, managing its bounds, input states, fullscreen restoration, and events.
type GLWindow struct {
	window             *glfw.Window
	bounds             Rect
	canvas             *GLCanvas
	cursorVisible      bool
	vSync              int
	cursorInsideWindow bool

	// need to save these to correctly restore a fullscreen window
	restore struct {
		xPos, yPos, width, height int
	}

	prevInp, currInp, tempInp struct {
		mouse   Vector
		buttons [KeyLast + 1]bool
		repeat  [KeyLast + 1]bool
		scroll  Vector
		typed   string
	}

	keysPressed                      map[Button]bool
	pressEvents, tempPressEvents     [KeyLast + 1]bool
	releaseEvents, tempReleaseEvents [KeyLast + 1]bool

	prevJoy, currJoy, tempJoy GLJoystick
}

// currWin holds a reference to the current GLWindow with an active OpenGL context.
var currWin *GLWindow

// NewGLWindow creates a new OpenGL window using the provided configuration, initializes it, and returns the GLWindow instance.
// Returns an error if the window creation fails or an invalid value is passed for certain configuration fields like msaaSamples.
func NewGLWindow(cfg WindowConfig) (*GLWindow, error) {
	bool2int := map[bool]int{
		true:  glfw.True,
		false: glfw.False,
	}
	w := &GLWindow{
		bounds:        cfg.Bounds,
		cursorVisible: true,
		keysPressed:   make(map[Button]bool),
	}

	flag := false
	for _, v := range []int{0, 2, 4, 8, 16} {
		if cfg.SamplesMSAA == v {
			flag = true
			break
		}
	}
	if !flag {
		return nil, fmt.Errorf("invalid value '%v' for msaaSamples", cfg.SamplesMSAA)
	}

	err := executor.GraphicThread.CallErr(func() error {
		var err error
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.Resizable, bool2int[cfg.Resizable])
		glfw.WindowHint(glfw.Decorated, bool2int[!cfg.Undecorated])
		glfw.WindowHint(glfw.Floating, bool2int[cfg.AlwaysOnTop])
		glfw.WindowHint(glfw.AutoIconify, bool2int[!cfg.NoIconify])
		glfw.WindowHint(glfw.TransparentFramebuffer, bool2int[cfg.TransparentFramebuffer])
		glfw.WindowHint(glfw.Maximized, bool2int[cfg.Maximized])
		glfw.WindowHint(glfw.Visible, bool2int[!cfg.Invisible])
		glfw.WindowHint(glfw.Samples, cfg.SamplesMSAA)
		if cfg.Position.X != 0 || cfg.Position.Y != 0 {
			glfw.WindowHint(glfw.Visible, glfw.False)
		}
		var share *glfw.Window
		if currWin != nil {
			share = currWin.window
		}
		_, _, width, height := intBounds(cfg.Bounds)
		w.window, err = glfw.CreateWindow(width, height, cfg.Title, nil, share)
		if err != nil {
			return err
		}
		if cfg.Position.X != 0 || cfg.Position.Y != 0 {
			w.window.SetPos(int(cfg.Position.X), int(cfg.Position.Y))
			w.window.Show()
		}
		w.begin()
		executor.Init()
		gl.Enable(gl.MULTISAMPLE)
		w.end()
		return nil
	})
	if err != nil {
		return nil, errors.New("creating window failed")
	}

	if len(cfg.Icon) > 0 {
		m := make([]image.Image, len(cfg.Icon))
		for i, icon := range cfg.Icon {
			pic := NewPictureFromPicture(icon)
			fmt.Println(pic, i)
			m[i] = pic.Image()
		}
		executor.GraphicThread.Call(func() {
			w.window.SetIcon(m)
		})
	}
	w.SetVSync(cfg.VSync)
	w.initInput()
	w.SetMonitor(cfg.Monitor)
	w.canvas = NewGLCanvas(cfg.Bounds, cfg.Smooth)
	w.Update()
	runtime.SetFinalizer(w, (*GLWindow).Destroy)
	return w, nil
}

// Destroy safely destroys the GLWindow instance by invoking its destruction on the graphics thread.
func (w *GLWindow) Destroy() {
	executor.GraphicThread.Call(func() {
		w.window.Destroy()
	})
}

// Update refreshes the GLWindow by handling resizing, framebuffer updates, and rendering tasks executed on the graphic thread.
func (w *GLWindow) Update() {
	bounds := w.bounds
	_, _, oldW, oldH := intBounds(bounds)
	newBounds := false

	executor.GraphicThread.Call(func() {
		newW, newH := w.window.GetSize()
		width := newW - oldW
		height := newH - oldH
		if width > 0 || height > 0 {
			bounds = bounds.ResizedMin(bounds.Size().Add(NewVec(float64(width), float64(height))))
			newBounds = true
		}

		w.begin()
		fbWidth, fbHeight := w.window.GetFramebufferSize()
		executor.Bounds(0, 0, fbWidth, fbHeight)
		executor.Clear(0, 0, 0, 0)
		w.canvas.gf.Frame().Begin()
		w.canvas.gf.Frame().Blit(nil, 0, 0, w.canvas.Texture().Width(), w.canvas.Texture().Height(), 0, 0, fbWidth, fbHeight)
		w.canvas.gf.Frame().End()
		glfw.SwapInterval(w.vSync)
		w.window.SwapBuffers()
		w.end()

		glfw.PollEvents()
	})
	if newBounds {
		w.canvas.SetBounds(bounds)
		w.bounds = bounds
	}
	w.doUpdateInput()
}

// ClipboardText retrieves the current text content from the system clipboard associated with the GLWindow instance.
func (w *GLWindow) ClipboardText() string {
	return w.window.GetClipboardString()
}

// SetClipboardText sets the given text to the system clipboard for the current GLWindow instance.
func (w *GLWindow) SetClipboardText(text string) {
	w.window.SetClipboardString(text)
}

// SetClosed sets the closed state of the GLWindow by signaling the graphics thread to update the window's close status.
func (w *GLWindow) SetClosed(closed bool) {
	executor.GraphicThread.Call(func() {
		w.window.SetShouldClose(closed)
	})
}

// Closed returns true if the window is marked for closing, otherwise false.
func (w *GLWindow) Closed() bool {
	var closed bool
	executor.GraphicThread.Call(func() {
		closed = w.window.ShouldClose()
	})
	return closed
}

// SetTitle sets the title of the GLWindow. The operation is performed on the graphical thread to ensure thread safety.
func (w *GLWindow) SetTitle(title string) {
	executor.GraphicThread.Call(func() {
		w.window.SetTitle(title)
	})
}

// SetBounds updates the dimensions of the GLWindow by setting its bounding rectangle.
// It adjusts the window's size on the graphic thread using the provided Rect dimensions.
func (w *GLWindow) SetBounds(bounds Rect) {
	w.bounds = bounds
	executor.GraphicThread.Call(func() {
		_, _, width, height := intBounds(bounds)
		w.window.SetSize(width, height)
	})
}

// SetPos sets the position of the window to the specified coordinates using the Vector type on the graphic thread.
func (w *GLWindow) SetPos(pos Vector) {
	executor.GraphicThread.Call(func() {
		left, top := int(pos.X), int(pos.Y)
		w.window.SetPos(left, top)
	})
}

// GetPos retrieves the position of the GLWindow as a Vector, making a thread-safe call to the graphic thread.
func (w *GLWindow) GetPos() Vector {
	var v Vector
	executor.GraphicThread.Call(func() {
		x, y := w.window.GetPos()
		v = NewVec(float64(x), float64(y))
	})
	return v
}

// Bounds returns the current bounds of the GLWindow as a Rect.
func (w *GLWindow) Bounds() Rect {
	return w.bounds
}

// setFullscreen switches the window to fullscreen mode on the specified monitor and stores the current window state.
func (w *GLWindow) setFullscreen(monitor *GLMonitor) {
	executor.GraphicThread.Call(func() {
		w.restore.xPos, w.restore.yPos = w.window.GetPos()
		w.restore.width, w.restore.height = w.window.GetSize()

		mode := monitor.monitor.GetVideoMode()

		w.window.SetMonitor(
			monitor.monitor,
			0,
			0,
			mode.Width,
			mode.Height,
			mode.RefreshRate,
		)
	})
}

// setWindowed switches the window back to windowed mode with its previous position and size, detaching it from any monitor.
func (w *GLWindow) setWindowed() {
	executor.GraphicThread.Call(func() {
		w.window.SetMonitor(
			nil,
			w.restore.xPos,
			w.restore.yPos,
			w.restore.width,
			w.restore.height,
			0,
		)
	})
}

// SetMonitor adjusts the window by associating it with a specific GLMonitor to enter fullscreen or detached to exit it.
func (w *GLWindow) SetMonitor(monitor *GLMonitor) {
	if w.Monitor() != monitor {
		if monitor != nil {
			w.setFullscreen(monitor)
		} else {
			w.setWindowed()
		}
	}
}

// Monitor retrieves the GLMonitor associated with the GLWindow, or nil if no monitor is currently associated.
func (w *GLWindow) Monitor() *GLMonitor {
	var monitor *glfw.Monitor
	executor.GraphicThread.Call(func() {
		monitor = w.window.GetMonitor()
	})
	if monitor == nil {
		return nil
	}
	return &GLMonitor{
		monitor: monitor,
	}
}

// Focused checks if the window is currently focused and returns true if focused, otherwise it returns false.
func (w *GLWindow) Focused() bool {
	var focused bool
	executor.GraphicThread.Call(func() {
		focused = w.window.GetAttrib(glfw.Focused) == glfw.True
	})
	return focused
}

// SetVSync configures whether vertical synchronization (VSync) is enabled or disabled for the GLWindow.
func (w *GLWindow) SetVSync(vSync bool) {
	if vSync {
		w.vSync = 1
	} else {
		w.vSync = 0
	}
}

// VSync checks if vertical synchronization (VSync) is enabled for the window, returning true if enabled, otherwise false.
func (w *GLWindow) VSync() bool {
	if w.vSync == 0 {
		return false
	}
	return true
}

// SetCursorVisible toggles the visibility of the window's cursor based on the provided boolean parameter.
func (w *GLWindow) SetCursorVisible(visible bool) {
	w.cursorVisible = visible
	executor.GraphicThread.Call(func() {
		if visible {
			w.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		} else {
			w.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		}
	})
}

// SetCursorDisabled sets the cursor mode to disabled, making it invisible and confined to the window.
func (w *GLWindow) SetCursorDisabled() {
	w.cursorVisible = false
	executor.GraphicThread.Call(func() {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	})
}

// CursorVisible returns a boolean indicating whether the cursor is currently visible in the window.
func (w *GLWindow) CursorVisible() bool {
	return w.cursorVisible
}

// begin sets the current OpenGL context to the window if it is not already the current context.
func (w *GLWindow) begin() {
	if currWin != w {
		w.window.MakeContextCurrent()
		currWin = w
	}
}

// end finalizes the rendering or processing for the current frame. It is a placeholder with no specific implementation.
func (w *GLWindow) end() {
	// nothing, really
}

// MakeTriangles creates a specialized copy of the supplied ITriangles that can draw onto the GLWindow's canvas.
func (w *GLWindow) MakeTriangles(t ITriangles) ITargetTriangles {
	return w.canvas.MakeTriangles(t)
}

// MakePicture creates a specialized copy of the provided IPicture that can be drawn onto the GLWindow's canvas.
func (w *GLWindow) MakePicture(p IPicture) ITargetPicture {
	return w.canvas.MakePicture(p)
}

// SetMatrix sets the transformation matrix for the window's canvas, applying spatial transforms like scaling or rotation.
func (w *GLWindow) SetMatrix(m Matrix) {
	w.canvas.SetMatrix(m)
}

// SetColorMask sets the color mask of the GLWindow by multiplying it with the provided color for rendering operations.
func (w *GLWindow) SetColorMask(c color.Color) {
	w.canvas.SetColorMask(c)
}

// SetComposeMethod sets the Porter-Duff composition method for rendering on the window's canvas.
func (w *GLWindow) SetComposeMethod(cmp ComposeMethod) {
	w.canvas.SetComposeMethod(cmp)
}

// SetSmooth sets whether the window's canvas should apply smoothing when resizing or scaling drawn content.
func (w *GLWindow) SetSmooth(smooth bool) {
	w.canvas.SetSmooth(smooth)
}

// Smooth returns whether the GLWindow's canvas is set to draw stretched pictures smoothly or not.
func (w *GLWindow) Smooth() bool {
	return w.canvas.Smooth()
}

// Clear fills the entire canvas of the GLWindow with the specified color.
func (w *GLWindow) Clear(c color.Color) {
	w.canvas.Clear(c)
}

// Color returns the RGBA color at the specified Vector position on the GLWindow's canvas.
func (w *GLWindow) Color(at Vector) RGBA {
	return w.canvas.Color(at)
}

// Canvas returns the associated instance of GLCanvas, which is an off-screen drawable rectangular area.
func (w *GLWindow) Canvas() *GLCanvas {
	return w.canvas
}

// Show makes the window visible by invoking the Show method on the internal window object within the graphic thread.
func (w *GLWindow) Show() {
	executor.GraphicThread.Call(func() {
		w.window.Show()
	})
}

// Clipboard retrieves the current text from the system clipboard associated with the GLWindow instance.
func (w *GLWindow) Clipboard() string {
	var clipboard string
	executor.GraphicThread.Call(func() {
		clipboard = w.window.GetClipboardString()
	})
	return clipboard
}

// SetClipboard sets the clipboard content to the provided string on the graphic thread.
func (w *GLWindow) SetClipboard(str string) {
	executor.GraphicThread.Call(func() {
		w.window.SetClipboardString(str)
	})
}

// KeysPressed returns a map where the keys currently pressed are marked as true, and others as false.
func (w *GLWindow) KeysPressed() map[Button]bool {
	return w.keysPressed
}

// Pressed returns true if the specified button is currently being pressed, otherwise false.
func (w *GLWindow) Pressed(button Button) bool {
	return w.currInp.buttons[button]
}

// PressedList returns a boolean array indicating the pressed state of keys up to the maximum defined key constant KeyLast.
func (w *GLWindow) PressedList() [KeyLast + 1]bool {
	return w.pressEvents
}

// ReleasedList returns an array of booleans indicating the release state of keys and buttons up to KeyLast.
func (w *GLWindow) ReleasedList() [KeyLast + 1]bool {
	return w.releaseEvents
}

// JustPressed checks if the specified button was just pressed during the current frame. Returns true if pressed, false otherwise.
func (w *GLWindow) JustPressed(button Button) bool {
	return w.pressEvents[button]
}

// JustReleased checks if the specified button was just released during the most recent input event.
func (w *GLWindow) JustReleased(button Button) bool {
	return w.releaseEvents[button]
}

// Repeated checks if the specified button is being held down and is generating repeated input events.
func (w *GLWindow) Repeated(button Button) bool {
	return w.currInp.repeat[button]
}

// MousePosition returns the current position of the mouse relative to the window as a Vector.
func (w *GLWindow) MousePosition() Vector {
	return w.currInp.mouse
}

// MousePositionX returns the current X-coordinate of the mouse position within the window as a float64.
func (w *GLWindow) MousePositionX() float64 {
	return w.currInp.mouse.X
}

// MousePositionY returns the current Y-coordinate of the mouse pointer relative to the window.
func (w *GLWindow) MousePositionY() float64 {
	return w.currInp.mouse.Y
}

// MousePositionXY retrieves the current mouse cursor position in the window as X and Y coordinates in float64 format.
func (w *GLWindow) MousePositionXY() (float64, float64) {
	return w.currInp.mouse.X, w.currInp.mouse.Y
}

// MousePreviousPosition returns the previous position of the mouse as a Vector.
func (w *GLWindow) MousePreviousPosition() Vector {
	return w.prevInp.mouse
}

// SetMousePosition sets the mouse position to the specified vector coordinates within the window bounds on the graphics thread.
func (w *GLWindow) SetMousePosition(v Vector) {
	executor.GraphicThread.Call(func() {
		if (v.X >= 0 && v.X <= w.bounds.W()) &&
			(v.Y >= 0 && v.Y <= w.bounds.H()) {
			w.window.SetCursorPos(v.X+w.bounds.Min.X, (w.bounds.H()-v.Y)+w.bounds.Min.Y)
			w.prevInp.mouse = v
			w.currInp.mouse = v
			w.tempInp.mouse = v
		}
	})
}

// MouseInsideWindow returns true if the mouse cursor is currently within the boundaries of the window.
func (w *GLWindow) MouseInsideWindow() bool {
	return w.cursorInsideWindow
}

// MouseScroll returns the current scroll offset as a 2D vector (Vector) representing the x and y scroll components.
func (w *GLWindow) MouseScroll() Vector {
	return w.currInp.scroll
}

// Typed returns the current typed input as a string from the GLWindow's internal input state.
func (w *GLWindow) Typed() string {
	return w.currInp.typed
}

// initInput initializes input callbacks for the window, handling mouse, keyboard, cursor, scroll, and character input events.
func (w *GLWindow) initInput() {
	executor.GraphicThread.Call(func() {
		w.window.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
			switch action {
			case glfw.Press:
				w.keysPressed[Button(button)] = true
				w.tempPressEvents[button] = true
				w.tempInp.buttons[button] = true
			case glfw.Release:
				delete(w.keysPressed, Button(button))
				w.tempReleaseEvents[button] = true
				w.tempInp.buttons[button] = false
			}
		})

		w.window.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeyUnknown {
				return
			}
			switch action {
			case glfw.Press:
				w.keysPressed[Button(key)] = true
				w.tempPressEvents[key] = true
				w.tempInp.buttons[key] = true
			case glfw.Release:
				delete(w.keysPressed, Button(key))
				w.tempReleaseEvents[key] = true
				w.tempInp.buttons[key] = false
			case glfw.Repeat:
				w.keysPressed[Button(key)] = true
				w.tempInp.repeat[key] = true
			}
		})

		w.window.SetCursorEnterCallback(func(_ *glfw.Window, entered bool) {
			w.cursorInsideWindow = entered
		})

		w.window.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
			w.tempInp.mouse = NewVec(x+w.bounds.Min.X, (w.bounds.H()-y)+w.bounds.Min.Y)
		})

		w.window.SetScrollCallback(func(_ *glfw.Window, xOff, yOff float64) {
			w.tempInp.scroll.X += xOff
			w.tempInp.scroll.Y += yOff
		})

		w.window.SetCharCallback(func(_ *glfw.Window, r rune) {
			w.tempInp.typed += string(r)
		})
	})
}

// UpdateInput poll window events.
// Call this function to poll window events without swapping buffers.
// Note that the Update method invoke UpdateInput.
//Func (w *GLWindow) UpdateInput() {
//	executor.GraphicThread.Call(func() { glfw.PollEvents() })
//	w.doUpdateInput()
//}

// UpdateInputWait processes input updates, pausing the thread to wait for events or a timeout on the graphics thread.
func (w *GLWindow) UpdateInputWait(timeout time.Duration) {
	executor.GraphicThread.Call(func() {
		if timeout <= 0 {
			glfw.WaitEvents()
		} else {
			glfw.WaitEventsTimeout(timeout.Seconds())
		}
	})
	w.doUpdateInput()
}

// doUpdateInput performs input state updates by transferring current input states to previous, clearing temporary states,
// and processing joystick inputs. It resets per-frame accumulators like key presses, releases, scroll data, and typing.
func (w *GLWindow) doUpdateInput() {
	w.prevInp = w.currInp
	w.currInp = w.tempInp

	//w.keysPressed = w.tempKeysPressed
	w.pressEvents = w.tempPressEvents
	w.releaseEvents = w.tempReleaseEvents

	// Clear last frame's temporary status
	//w.tempKeysPressed = []Button{}
	w.tempPressEvents = [KeyLast + 1]bool{}
	w.tempReleaseEvents = [KeyLast + 1]bool{}
	w.tempInp.repeat = [KeyLast + 1]bool{}
	w.tempInp.scroll = ZeroVector
	w.tempInp.typed = ""

	w.updateJoystickInput()
}

// JoystickPresent checks whether the specified joystick is currently connected.
func (w *GLWindow) JoystickPresent(js Joystick) bool {
	return w.currJoy.connected[js]
}

// JoystickName returns the name of the specified joystick.
func (w *GLWindow) JoystickName(js Joystick) string {
	return w.currJoy.name[js]
}

// JoystickButtonCount returns the number of buttons on the specified joystick.
func (w *GLWindow) JoystickButtonCount(js Joystick) int {
	return len(w.currJoy.buttons[js])
}

// JoystickAxisCount returns the number of axes for the specified joystick.
func (w *GLWindow) JoystickAxisCount(js Joystick) int {
	return len(w.currJoy.axis[js])
}

// JoystickPressed checks if the specified button on the given joystick is currently pressed and returns true or false.
func (w *GLWindow) JoystickPressed(js Joystick, button GamepadButton) bool {
	return w.currJoy.getButton(js, int(button))
}

// JoystickJustPressed returns true if the specified button on the given joystick was just pressed in the current frame.
func (w *GLWindow) JoystickJustPressed(js Joystick, button GamepadButton) bool {
	return w.currJoy.getButton(js, int(button)) && !w.prevJoy.getButton(js, int(button))
}

// JoystickJustReleased returns true if the specified button on the given joystick was just released in the current frame.
func (w *GLWindow) JoystickJustReleased(js Joystick, button GamepadButton) bool {
	return !w.currJoy.getButton(js, int(button)) && w.prevJoy.getButton(js, int(button))
}

// JoystickAxis retrieves the current value of the specified axis for the given joystick as a float64.
func (w *GLWindow) JoystickAxis(js Joystick, axis GamepadAxis) float64 {
	return w.currJoy.getAxis(js, int(axis))
}

// updateJoystickInput updates the state of all joysticks, including connection status, buttons, axes, and names.
func (w *GLWindow) updateJoystickInput() {
	for js := Joystick1; js <= JoystickLast; js++ {
		// Determine and store if the joystick was connected
		joystickPresent := glfw.Joystick(js).Present()
		w.tempJoy.connected[js] = joystickPresent
		if joystickPresent {
			if glfw.Joystick(js).IsGamepad() {
				if gamepadInputs := glfw.Joystick(js).GetGamepadState(); gamepadInputs != nil {
					w.tempJoy.buttons[js] = gamepadInputs.Buttons[:]
					w.tempJoy.axis[js] = gamepadInputs.Axes[:]
				}
			} else {
				w.tempJoy.buttons[js] = glfw.Joystick(js).GetButtons()
				w.tempJoy.axis[js] = glfw.Joystick(js).GetAxes()
			}

			if !w.currJoy.connected[js] {
				// The joystick was recently connected, we get the name
				w.tempJoy.name[js] = glfw.Joystick(js).GetName()
			} else {
				// Use the name from the previous one
				w.tempJoy.name[js] = w.currJoy.name[js]
			}
		} else {
			w.tempJoy.buttons[js] = []glfw.Action{}
			w.tempJoy.axis[js] = []float32{}
			w.tempJoy.name[js] = ""
		}
	}

	w.prevJoy = w.currJoy
	w.currJoy = w.tempJoy
}

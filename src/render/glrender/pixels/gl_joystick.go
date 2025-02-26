package pixels

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Joystick represents a joystick input device and maps to specific identifiers for each joystick supported.
type Joystick int

// Joystick1 to Joystick16 represent GLFW joystick mappings for respective joystick IDs.
// JoystickLast represents the last available joystick ID for GLFW.
const (
	Joystick1  = Joystick(glfw.Joystick1)
	Joystick2  = Joystick(glfw.Joystick2)
	Joystick3  = Joystick(glfw.Joystick3)
	Joystick4  = Joystick(glfw.Joystick4)
	Joystick5  = Joystick(glfw.Joystick5)
	Joystick6  = Joystick(glfw.Joystick6)
	Joystick7  = Joystick(glfw.Joystick7)
	Joystick8  = Joystick(glfw.Joystick8)
	Joystick9  = Joystick(glfw.Joystick9)
	Joystick10 = Joystick(glfw.Joystick10)
	Joystick11 = Joystick(glfw.Joystick11)
	Joystick12 = Joystick(glfw.Joystick12)
	Joystick13 = Joystick(glfw.Joystick13)
	Joystick14 = Joystick(glfw.Joystick14)
	Joystick15 = Joystick(glfw.Joystick15)
	Joystick16 = Joystick(glfw.Joystick16)

	JoystickLast = Joystick(glfw.JoystickLast)
)

// GamepadAxis represents an axis on a gamepad, typically associated with analog sticks or trigger controls.
type GamepadAxis int

// AxisLeftX represents the horizontal axis of the left thumbstick on a gamepad.
// AxisLeftY represents the vertical axis of the left thumbstick on a gamepad.
// AxisRightX represents the horizontal axis of the right thumbstick on a gamepad.
// AxisRightY represents the vertical axis of the right thumbstick on a gamepad.
// AxisLeftTrigger represents the axis of the left trigger on a gamepad.
// AxisRightTrigger represents the axis of the right trigger on a gamepad.
// AxisLast represents the last valid axis on a gamepad.
const (
	AxisLeftX        = GamepadAxis(glfw.AxisLeftX)
	AxisLeftY        = GamepadAxis(glfw.AxisLeftY)
	AxisRightX       = GamepadAxis(glfw.AxisRightX)
	AxisRightY       = GamepadAxis(glfw.AxisRightY)
	AxisLeftTrigger  = GamepadAxis(glfw.AxisLeftTrigger)
	AxisRightTrigger = GamepadAxis(glfw.AxisRightTrigger)
	AxisLast         = GamepadAxis(glfw.AxisLast)
)

// GamepadButton represents a specific button on a gamepad controller.
type GamepadButton int

// ButtonA represents the A button on a gamepad.
// ButtonB represents the B button on a gamepad.
// ButtonX represents the X button on a gamepad.
// ButtonY represents the Y button on a gamepad.
// ButtonLeftBumper represents the left bumper button on a gamepad.
// ButtonRightBumper represents the right bumper button on a gamepad.
// ButtonBack represents the back button on a gamepad.
// ButtonStart represents the start button on a gamepad.
// ButtonGuide represents the guide button on a gamepad.
// ButtonLeftThumb represents the left thumbstick button on a gamepad.
// ButtonRightThumb represents the right thumbstick button on a gamepad.
// ButtonDPadUp represents the up direction on the D-Pad of a gamepad.
// ButtonDPadRight represents the right direction on the D-Pad of a gamepad.
// ButtonDPadDown represents the down direction on the D-Pad of a gamepad.
// ButtonDPadLeft represents the left direction on the D-Pad of a gamepad.
// ButtonLast represents the last button enumerated on a gamepad.
// ButtonCross represents the cross button on a gamepad (common on PlayStation-style controllers).
// ButtonCircle represents the circle button on a gamepad (common on PlayStation-style controllers).
// ButtonSquare represents the square button on a gamepad (common on PlayStation-style controllers).
// ButtonTriangle represents the triangle button on a gamepad (common on PlayStation-style controllers).
const (
	ButtonA           = GamepadButton(glfw.ButtonA)
	ButtonB           = GamepadButton(glfw.ButtonB)
	ButtonX           = GamepadButton(glfw.ButtonX)
	ButtonY           = GamepadButton(glfw.ButtonY)
	ButtonLeftBumper  = GamepadButton(glfw.ButtonLeftBumper)
	ButtonRightBumper = GamepadButton(glfw.ButtonRightBumper)
	ButtonBack        = GamepadButton(glfw.ButtonBack)
	ButtonStart       = GamepadButton(glfw.ButtonStart)
	ButtonGuide       = GamepadButton(glfw.ButtonGuide)
	ButtonLeftThumb   = GamepadButton(glfw.ButtonLeftThumb)
	ButtonRightThumb  = GamepadButton(glfw.ButtonRightThumb)
	ButtonDPadUp      = GamepadButton(glfw.ButtonDpadUp)
	ButtonDPadRight   = GamepadButton(glfw.ButtonDpadRight)
	ButtonDPadDown    = GamepadButton(glfw.ButtonDpadDown)
	ButtonDPadLeft    = GamepadButton(glfw.ButtonDpadLeft)
	ButtonLast        = GamepadButton(glfw.ButtonLast)
	ButtonCross       = GamepadButton(glfw.ButtonCross)
	ButtonCircle      = GamepadButton(glfw.ButtonCircle)
	ButtonSquare      = GamepadButton(glfw.ButtonSquare)
	ButtonTriangle    = GamepadButton(glfw.ButtonTriangle)
)

// GLJoystick tracks the state of connected joysticks, including their connection status, names, buttons, and axis data.
type GLJoystick struct {
	connected [JoystickLast + 1]bool
	name      [JoystickLast + 1]string
	buttons   [JoystickLast + 1][]glfw.Action
	axis      [JoystickLast + 1][]float32
}

// getButton checks if the specified button on the given joystick is currently pressed and returns true or false accordingly.
func (js *GLJoystick) getButton(joystick Joystick, button int) bool {
	// Check that the joystick and button are valid, return false by default
	if js.buttons[joystick] == nil || button >= len(js.buttons[joystick]) || button < 0 {
		return false
	}
	return js.buttons[joystick][byte(button)] == glfw.Press
}

// getAxis retrieves the current value of the specified axis for the given joystick, returning 0 if invalid.
func (js *GLJoystick) getAxis(joystick Joystick, axis int) float64 {
	// Check that the joystick and axis are valid, return 0 by default.
	if js.axis[joystick] == nil || axis >= len(js.axis[joystick]) || axis < 0 {
		return 0
	}
	return float64(js.axis[joystick][axis])
}

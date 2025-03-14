package references

// IJoystick represents an interface for handling joystick operations and its state.
type IJoystick interface {
	Update()

	Reset()

	Emulate()

	Move(x uint, y uint, buttons uint)

	SetKey(pressed bool, jId int)

	Poll()
}

package references

// IJoystick defines an interface for handling joystick operations and interactions.
type IJoystick interface {
	Update()

	Reset()

	Emulate()

	Move(x uint, y uint, buttons uint)

	SetKey(pressed bool, jId int)

	Poll()
}

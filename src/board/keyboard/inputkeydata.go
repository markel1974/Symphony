package keyboard

type InputKeyData struct {
	c64Byte int
	c64Bit  int
	shifted bool
	joyKey  uint8
	pressed bool
	counter uint8
}

func NewInputKeyData(pressed bool, c64Byte int, c64bit int, shifted bool) *InputKeyData {
	return &InputKeyData{
		c64Byte: c64Byte,
		c64Bit:  c64bit,
		shifted: shifted,
		pressed: pressed,
		counter: CMD_COUNT,
	}
}

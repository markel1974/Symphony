package keyboard

type InputKeyData struct {
	keyM    int
	revM    int
	shifted bool
	joyKey  uint8
	pressed bool
}

func NewInputKeyData(pressed bool, keyM int, revM int, shifted bool) *InputKeyData {
	return &InputKeyData{
		keyM:    keyM,
		revM:    revM,
		shifted: shifted,
		pressed: pressed,
	}
}

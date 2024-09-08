package inputs

type InputKeyData struct {
	keyM        int
	revM        int
	shifted     bool
	joyKey      uint8
	pressed     bool
	persistence uint8
}

func NewInputKeyData(pressed bool, keyM int, revM int, shifted bool, persistence uint8) *InputKeyData {
	return &InputKeyData{
		keyM:        keyM,
		revM:        revM,
		shifted:     shifted,
		pressed:     pressed,
		persistence: persistence,
	}
}

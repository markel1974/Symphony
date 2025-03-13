package iecprotocol

type IGlobal interface {
	GetState() uint8
	SetState() uint8
}

type Global struct {
	state uint8
}

func NewGlobalState() *Global {
	return &Global{
		state: 0,
	}
}

func (v *Global) SetState(state uint8) {
	v.state = state
}

func (v *Global) GetState() uint8 {
	return v.state
}

var _gs = NewGlobalState()

package signals

type FnByte func(a uint8)

type SignalByte struct {
	receiver  FnByte
	container []FnByte
}

func NewSignalByte() *SignalByte {
	s := &SignalByte{}
	s.receiver = s.void
	return s
}

func (s *SignalByte) Bind(b FnByte) {
	s.container = append(s.container, b)
	if len(s.container) == 1 {
		s.receiver = b
	} else {
		s.receiver = s.multi
	}
}

func (s *SignalByte) Emit(v1 uint8) {
	s.receiver(v1)
}

func (s *SignalByte) void(_ uint8) {
}

func (s *SignalByte) multi(v uint8) {
	for _, r := range s.container {
		r(v)
	}
}

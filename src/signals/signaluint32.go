package signals

type FnUint32 func(a uint32)

type SignalUint32 struct {
	receiver  FnUint32
	container []FnUint32
}

func NewSignalUint32() *SignalUint32 {
	s := &SignalUint32{}
	s.receiver = s.void
	return s
}

func (s *SignalUint32) Bind(b FnUint32) {
	s.container = append(s.container, b)
	if len(s.container) == 1 {
		s.receiver = b
	} else {
		s.receiver = s.multi
	}
}

func (s *SignalUint32) Emit(v1 uint32) {
	s.receiver(v1)
}

func (s *SignalUint32) void(_ uint32) {
}

func (s *SignalUint32) multi(v uint32) {
	for _, r := range s.container {
		r(v)
	}
}

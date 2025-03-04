package signals

type FnBool func(a bool)

type SignalBool struct {
	receiver  FnBool
	container []FnBool
}

func NewSignalBool() *SignalBool {
	s := &SignalBool{}
	s.receiver = s.void
	return s
}

func (s *SignalBool) Bind(b FnBool) {
	s.container = append(s.container, b)
	if len(s.container) == 1 {
		s.receiver = b
	} else {
		s.receiver = s.multi
	}
}

func (s *SignalBool) Emit(v1 bool) {
	s.receiver(v1)
}

func (s *SignalBool) void(_ bool) {
}

func (s *SignalBool) multi(v bool) {
	for _, r := range s.container {
		r(v)
	}
}

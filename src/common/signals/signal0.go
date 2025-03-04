package signals

type Fn0 func()

type Signal struct {
	receiver  Fn0
	container []Fn0
}

func NewSignal() *Signal {
	s := &Signal{}
	s.receiver = s.void
	return s
}

func (s *Signal) Bind(b Fn0) {
	s.container = append(s.container, b)
	if len(s.container) == 1 {
		s.receiver = b
	} else {
		s.receiver = s.multi
	}
}

func (s *Signal) Emit() {
	s.receiver()
}

func (s *Signal) void() {
}

func (s *Signal) multi() {
	for _, r := range s.container {
		r()
	}
}

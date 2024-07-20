package signals

type Fn0 func()

type Signal struct {
	receivers []Fn0
}

func NewSignal() *Signal {
	return &Signal{}
}

func (s *Signal) Bind(b Fn0) {
	s.receivers = append(s.receivers, b)
}

func (s *Signal) Emit() {
	if s.receivers == nil {
		return
	}
	if len(s.receivers) == 1 {
		s.receivers[0]()
		return
	}
	for _, r := range s.receivers {
		r()
	}
}

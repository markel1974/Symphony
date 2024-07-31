package signals

type FnBool func(a bool)

type SignalBool struct {
	receivers []FnBool
}

func NewSignalBool() *SignalBool {
	return &SignalBool{}
}

func (s *SignalBool) Bind(b FnBool) {
	s.receivers = append(s.receivers, b)
}

func (s *SignalBool) Emit(v1 bool) {
	if s.receivers == nil {
		return
	}
	if len(s.receivers) == 1 {
		s.receivers[0](v1)
		return
	}
	for _, r := range s.receivers {
		r(v1)
	}
}

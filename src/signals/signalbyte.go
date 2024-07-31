package signals

type FnByte func(a uint8)

type SignalByte struct {
	receivers []FnByte
}

func NewSignalByte() *SignalByte {
	return &SignalByte{}
}

func (s *SignalByte) Bind(b FnByte) {
	s.receivers = append(s.receivers, b)
}

func (s *SignalByte) Emit(v1 uint8) {
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

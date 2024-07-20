package signals

type Fn1[T1 any] func(a T1)

type Signal1[T1 any] struct {
	receivers []Fn1[T1]
}

func NewSignal1[T1 any]() *Signal1[T1] {
	return &Signal1[T1]{}
}

func (s *Signal1[T1]) Bind(b Fn1[T1]) {
	s.receivers = append(s.receivers, b)
}

func (s *Signal1[T1]) Emit(v1 any) {
	if s.receivers == nil {
		return
	}
	t1, ok := v1.(T1)
	if !ok {
		return
	}
	if len(s.receivers) == 1 {
		s.receivers[0](t1)
		return
	}
	for _, r := range s.receivers {
		r(t1)
	}
}

package signals

type Fn2[T1 any, T2 any] func(a T1, b T2)

type Signal2[T1 any, T2 any] struct {
	receivers []Fn2[T1, T2]
}

func NewSignal2[T1 any, T2 any]() *Signal2[T1, T2] {
	return &Signal2[T1, T2]{}
}

func (s *Signal2[T1, T2]) Bind(b Fn2[T1, T2]) {
	s.receivers = append(s.receivers, b)
}

func (s *Signal2[T1, T2]) Emit(v1 any, v2 any) {
	if s.receivers == nil {
		return
	}
	t1, ok := v1.(T1)
	if !ok {
		return
	}
	t2, ok := v2.(T2)
	if !ok {
		return
	}
	if len(s.receivers) == 1 {
		s.receivers[0](t1, t2)
		return
	}
	for _, r := range s.receivers {
		r(t1, t2)
	}
}

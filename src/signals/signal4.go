package signals

type Fn4[T1 any, T2 any, T3 any, T4 any] func(a T1, b T2, c T3, d T4)

type Signal4[T1 any, T2 any, T3 any, T4 any] struct {
	receivers []Fn4[T1, T2, T3, T4]
}

func NewSignal4[T1 any, T2 any, T3 any, T4 any]() *Signal4[T1, T2, T3, T4] {
	return &Signal4[T1, T2, T3, T4]{}
}

func (s *Signal4[T1, T2, T3, T4]) Bind(b Fn4[T1, T2, T3, T4]) {
	s.receivers = append(s.receivers, b)
}

func (s *Signal4[T1, T2, T3, T4]) Emit(v1 any, v2 any, v3 any, v4 any) {
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
	t3, ok := v3.(T3)
	if !ok {
		return
	}
	t4, ok := v4.(T4)
	if !ok {
		return
	}
	if len(s.receivers) == 1 {
		s.receivers[0](t1, t2, t3, t4)
		return
	}
	for _, r := range s.receivers {
		r(t1, t2, t3, t4)
	}
}

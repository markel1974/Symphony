package signals

type Fn5[T1 any, T2 any, T3 any, T4 any, T5 any] func(a T1, b T2, c T3, d T4, e T5)

type Signal5[T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	receivers []Fn5[T1, T2, T3, T4, T5]
}

func NewSignal5[T1 any, T2 any, T3 any, T4 any, T5 any]() *Signal5[T1, T2, T3, T4, T5] {
	return &Signal5[T1, T2, T3, T4, T5]{}
}

func (s *Signal5[T1, T2, T3, T4, T5]) Bind(b Fn5[T1, T2, T3, T4, T5]) {
	s.receivers = append(s.receivers, b)
}

func (s *Signal5[T1, T2, T3, T4, T5]) Emit(v1 any, v2 any, v3 any, v4 any, v5 any) {
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
	t5, ok := v5.(T5)
	if !ok {
		return
	}
	if len(s.receivers) == 1 {
		s.receivers[0](t1, t2, t3, t4, t5)
		return
	}
	for _, r := range s.receivers {
		r(t1, t2, t3, t4, t5)
	}
}

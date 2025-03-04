package signals

type Fn3[T1 any, T2 any, T3 any] func(a T1, b T2, c T3)

type Signal3[T1 any, T2 any, T3 any] struct {
	receivers []Fn3[T1, T2, T3]
}

func NewSignal3[T1 any, T2 any, T3 any]() *Signal3[T1, T2, T3] {
	return &Signal3[T1, T2, T3]{}
}

func (s *Signal3[T1, T2, T3]) Bind(b Fn3[T1, T2, T3]) {
	s.receivers = append(s.receivers, b)
}

func (s *Signal3[T1, T2, T3]) Emit(v1 any, v2 any, v3 any) {
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
	if len(s.receivers) == 1 {
		s.receivers[0](t1, t2, t3)
		return
	}
	for _, r := range s.receivers {
		r(t1, t2, t3)
	}
}

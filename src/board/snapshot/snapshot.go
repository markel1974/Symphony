package snapshot

type Snapshot struct {
	container map[string]*Module
}

func NewSnapshot() *Snapshot {
	return &Snapshot{
		container: make(map[string]*Module),
	}
}

func (s *Snapshot) NewModule(name string, major int, minor int) *Module {
	m := NewModule(name, major, minor)
	s.container[name] = m
	return m
}

func (s *Snapshot) GetModule(name string) *Module {
	return s.container[name]
}

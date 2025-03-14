package snapshot

type Module struct {
	Name      string
	Major     int
	Minor     int
	Container map[string]interface{}
}

func NewModule(name string, major int, minor int) *Module {
	return &Module{
		Name:      name,
		Major:     major,
		Minor:     minor,
		Container: make(map[string]interface{}),
	}
}

func (m *Module) Add(id string, value interface{}) {
	m.Container[id] = value
}

func (m *Module) Get(id string) interface{} {
	if v, ok := m.Container[id]; ok {
		return v
	}
	return nil
}

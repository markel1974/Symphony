package mechanics

type Core struct {
	_current_halftrack int
	_gcr               uint32
	_write_protected   bool
	_disk_changed      bool
}

func NewCore() *Core {
	c := &Core{}
	c.Reset()
	return c
}

func (c *Core) Reset() {
	c._current_halftrack = 0
	c._gcr = 0
	c._write_protected = false
	c._disk_changed = false
}

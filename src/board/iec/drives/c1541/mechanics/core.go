package mechanics

type Core struct {
	currentHalfTrack int
	gcr              uint32
	writeProtected   bool
	diskChanged      bool
}

func NewCore() *Core {
	c := &Core{}
	c.Reset()
	return c
}

func (c *Core) Reset() {
	c.currentHalfTrack = 2
	c.gcr = 0
	c.writeProtected = false
	c.diskChanged = false
}

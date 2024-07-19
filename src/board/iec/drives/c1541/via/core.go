package via

const defaultViaTimeout = 0xffff

type Core struct {
	pra  uint8
	ddra uint8
	prb  uint8
	ddrb uint8
	t1c  uint16
	t1l  uint16
	t2c  uint16
	t2l  uint16
	sr   uint8
	acr  uint8
	pcr  uint8
	ifr  uint8
	ier  uint8
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) Reset() {
	c.pra = 0
	c.ddra = 0
	c.prb = 0
	c.ddrb = 0
	c.t1c = 0
	c.t1l = 0
	c.t2c = 0
	c.t2l = 0
	c.sr = 0
	c.acr = 0
	c.pcr = 0
	c.ifr = 0
	c.ier = 0
}

package via

const DEFAULT_VIA_TIMEOUT = 0xffff

type Core struct {
	_pra  uint8
	_ddra uint8
	_prb  uint8
	_ddrb uint8
	_t1c  uint16
	_t1l  uint16
	_t2c  uint16
	_t2l  uint16
	_sr   uint8
	_acr  uint8
	_pcr  uint8
	_ifr  uint8
	_ier  uint8
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) Reset() {
	c._pra = 0
	c._ddra = 0
	c._prb = 0
	c._ddrb = 0
	c._t1c = 0
	c._t1l = 0
	c._t2c = 0
	c._t2l = 0
	c._sr = 0
	c._acr = 0
	c._pcr = 0
	c._ifr = 0
	c._ier = 0
}

package ram

type Ram struct {
	r []uint8
}

func New(size int) *Ram {
	return &Ram{r: make([]uint8, size)}
}

func (r *Ram) Setup() {
}

func (r *Ram) Read(addr uint16) uint8 {
	return r.r[addr]
}

func (r *Ram) Write(addr uint16, data uint8) {
	r.r[addr] = data
}

func (r *Ram) Interval(start uint16, count uint16) []byte {
	return r.r[start : start+count]
}

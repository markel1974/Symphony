package Interrupt

type Interrupt uint32

func (i *Interrupt) Clear() {
	*i = 0
}

func (i *Interrupt) BitSet(n uint32) {
	*i = *i | (1 << n)
}

func (i *Interrupt) BitClear(n uint32) {
	*i = *i & ^(1 << n)
}

func (i *Interrupt) BitCheck(n uint32) bool {
	v := (*i >> n) & 1
	return v != 0
}

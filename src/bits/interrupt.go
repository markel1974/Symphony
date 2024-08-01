package bits

type Bits uint32

func (i *Bits) Clear() {
	*i = 0
}

func (i *Bits) BitSet(n uint32) {
	*i = *i | (1 << n)
}

func (i *Bits) BitClear(n uint32) {
	*i = *i & ^(1 << n)
}

func (i *Bits) BitCheck(n uint32) bool {
	v := (*i >> n) & 1
	return v != 0
}

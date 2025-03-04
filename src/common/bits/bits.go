package bits

var Uint8s = []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}

type Bits uint32

// Clear resets all bits in the Bits instance to zero.
func (i *Bits) Clear() {
	*i = 0
}

// BitSet sets the bit at position n to 1.
func (i *Bits) BitSet(n uint32) {
	*i = *i | (1 << (n & 0x1f))
}

// BitClear clears the bit at the specified position n
func (i *Bits) BitClear(n uint32) {
	*i = *i & ^(1 << (n & 0x1f))
}

// BitCheck checks if the bit at position n is set and returns true; otherwise, it returns false.
func (i *Bits) BitCheck(n uint32) bool {
	v := (*i >> (n & 0x1f)) & 1
	return v != 0
}

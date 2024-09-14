package inputs

type KeyData struct {
	r1      int
	shifted int
}

func NewKeyData(a int, b int, shifted int) *KeyData {
	return &KeyData{
		r1:      matrix(a, b),
		shifted: shifted,
	}
}

type Ascii struct {
	container []*KeyData
}

func NewAscii() *Ascii {
	a := &Ascii{}
	a.container = make([]*KeyData, 256)
	space := NewKeyData(7, 4, 0) //NewInvalidKeyData()
	for x := range a.container {
		a.container[x] = space
	}
	a.container['!'] = NewKeyData(7, 0, 1)
	a.container['"'] = NewKeyData(7, 3, 1)
	a.container['#'] = NewKeyData(1, 0, 1)
	a.container['$'] = NewKeyData(1, 3, 1)
	a.container['%'] = NewKeyData(2, 0, 1)
	a.container['&'] = NewKeyData(2, 3, 1)
	a.container['\\'] = NewKeyData(3, 0, 1)
	a.container['('] = NewKeyData(3, 3, 1)
	a.container[')'] = NewKeyData(4, 0, 1)
	a.container['>'] = NewKeyData(5, 4, 1)
	a.container['<'] = NewKeyData(5, 7, 1)
	a.container['?'] = NewKeyData(6, 2, 1)
	a.container['['] = NewKeyData(5, 5, 1)
	a.container[']'] = NewKeyData(6, 2, 1)

	a.container['\n'] = NewKeyData(0, 1, 0)
	a.container['\r'] = NewKeyData(0, 1, 0)
	a.container[' '] = NewKeyData(7, 4, 0)
	a.container['/'] = NewKeyData(6, 7, 0)
	a.container['^'] = NewKeyData(6, 6, 0)
	a.container['='] = NewKeyData(6, 5, 0)
	a.container[';'] = NewKeyData(6, 2, 0)
	a.container['*'] = NewKeyData(6, 1, 0)
	//container['£'] = NewKeyData(6, 0, 0)
	a.container[','] = NewKeyData(5, 7, 0)
	a.container['@'] = NewKeyData(5, 6, 0)
	a.container[':'] = NewKeyData(5, 5, 0)
	a.container['.'] = NewKeyData(5, 4, 0)
	a.container['-'] = NewKeyData(5, 3, 0)
	a.container['+'] = NewKeyData(5, 0, 0)
	a.container['0'] = NewKeyData(4, 3, 0)
	a.container['1'] = NewKeyData(7, 0, 0)
	a.container['2'] = NewKeyData(7, 3, 0)
	a.container['3'] = NewKeyData(1, 0, 0)
	a.container['4'] = NewKeyData(1, 3, 0)
	a.container['5'] = NewKeyData(2, 0, 0)
	a.container['6'] = NewKeyData(2, 3, 0)
	a.container['7'] = NewKeyData(3, 0, 0)
	a.container['8'] = NewKeyData(3, 3, 0)
	a.container['9'] = NewKeyData(4, 0, 0)
	a.container['A'] = NewKeyData(1, 2, 0)
	a.container['B'] = NewKeyData(3, 4, 0)
	a.container['C'] = NewKeyData(2, 4, 0)
	a.container['D'] = NewKeyData(2, 2, 0)
	a.container['E'] = NewKeyData(1, 6, 0)
	a.container['F'] = NewKeyData(2, 5, 0)
	a.container['G'] = NewKeyData(3, 2, 0)
	a.container['H'] = NewKeyData(3, 5, 0)
	a.container['I'] = NewKeyData(4, 1, 0)
	a.container['J'] = NewKeyData(4, 2, 0)
	a.container['K'] = NewKeyData(4, 5, 0)
	a.container['L'] = NewKeyData(5, 2, 0)
	a.container['M'] = NewKeyData(4, 4, 0)
	a.container['N'] = NewKeyData(4, 7, 0)
	a.container['O'] = NewKeyData(4, 6, 0)
	a.container['P'] = NewKeyData(5, 1, 0)
	a.container['Q'] = NewKeyData(7, 6, 0)
	a.container['R'] = NewKeyData(2, 1, 0)
	a.container['S'] = NewKeyData(1, 5, 0)
	a.container['T'] = NewKeyData(2, 6, 0)
	a.container['U'] = NewKeyData(3, 6, 0)
	a.container['V'] = NewKeyData(3, 7, 0)
	a.container['W'] = NewKeyData(1, 1, 0)
	a.container['X'] = NewKeyData(2, 7, 0)
	a.container['Y'] = NewKeyData(3, 1, 0)
	a.container['Z'] = NewKeyData(1, 4, 0)
	a.container['a'] = NewKeyData(1, 2, 0)
	a.container['b'] = NewKeyData(3, 4, 0)
	a.container['c'] = NewKeyData(2, 4, 0)
	a.container['d'] = NewKeyData(2, 2, 0)
	a.container['e'] = NewKeyData(1, 6, 0)
	a.container['f'] = NewKeyData(2, 5, 0)
	a.container['g'] = NewKeyData(3, 2, 0)
	a.container['h'] = NewKeyData(3, 5, 0)
	a.container['i'] = NewKeyData(4, 1, 0)
	a.container['j'] = NewKeyData(4, 2, 0)
	a.container['k'] = NewKeyData(4, 5, 0)
	a.container['l'] = NewKeyData(5, 2, 0)
	a.container['m'] = NewKeyData(4, 4, 0)
	a.container['n'] = NewKeyData(4, 7, 0)
	a.container['o'] = NewKeyData(4, 6, 0)
	a.container['p'] = NewKeyData(5, 1, 0)
	a.container['q'] = NewKeyData(7, 6, 0)
	a.container['r'] = NewKeyData(2, 1, 0)
	a.container['s'] = NewKeyData(1, 5, 0)
	a.container['t'] = NewKeyData(2, 6, 0)
	a.container['u'] = NewKeyData(3, 6, 0)
	a.container['v'] = NewKeyData(3, 7, 0)
	a.container['w'] = NewKeyData(1, 1, 0)
	a.container['x'] = NewKeyData(2, 7, 0)
	a.container['y'] = NewKeyData(3, 1, 0)
	a.container['z'] = NewKeyData(1, 4, 0)
	return a
}

func (a *Ascii) Reset() {

}

func (a *Ascii) FromAscii(v uint8) *KeyData {
	out := a.container[v]
	return out
}

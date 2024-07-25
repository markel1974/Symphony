package keyboard

type KeyData struct {
	r1 int
	r2 int
}

func NewSingleKeyData(a int, b int) *KeyData {
	return &KeyData{
		r1: matrix(a, b),
		r2: 0,
	}
}

func NewDoubleKeyData(a int, b int, c int, d int) *KeyData {
	return &KeyData{
		r1: matrix(a, b),
		r2: matrix(c, d),
	}
}

type Ascii struct {
	container []*KeyData
}

func NewAscii() *Ascii {
	a := &Ascii{}
	a.container = make([]*KeyData, 256)
	space := NewSingleKeyData(7, 4) //NewInvalidKeyData()
	for x := range a.container {
		a.container[x] = space
	}
	a.container['!'] = NewDoubleKeyData(6, 4, 7, 0)
	a.container['"'] = NewDoubleKeyData(6, 4, 7, 3)
	a.container['#'] = NewDoubleKeyData(6, 4, 1, 0)
	a.container['$'] = NewDoubleKeyData(6, 4, 1, 3)
	a.container['%'] = NewDoubleKeyData(6, 4, 2, 0)
	a.container['&'] = NewDoubleKeyData(6, 4, 2, 3)
	a.container['\\'] = NewDoubleKeyData(6, 4, 3, 0)
	a.container['('] = NewDoubleKeyData(6, 4, 3, 3)
	a.container[')'] = NewDoubleKeyData(6, 4, 4, 0)
	a.container['>'] = NewDoubleKeyData(6, 4, 5, 4)
	a.container['<'] = NewDoubleKeyData(6, 4, 5, 7)
	a.container['?'] = NewDoubleKeyData(6, 4, 6, 2)
	a.container['['] = NewDoubleKeyData(6, 4, 5, 5)
	a.container[']'] = NewDoubleKeyData(6, 4, 6, 2)
	a.container['\n'] = NewSingleKeyData(0, 1)
	a.container['\r'] = NewSingleKeyData(0, 1)
	a.container[' '] = NewSingleKeyData(7, 4)
	a.container['/'] = NewSingleKeyData(6, 7)
	a.container['^'] = NewSingleKeyData(6, 6)
	a.container['='] = NewSingleKeyData(6, 5)
	a.container[';'] = NewSingleKeyData(6, 2)
	a.container['*'] = NewSingleKeyData(6, 1)
	//container['£'] = NewSingleKeyData(6, 0)
	a.container[','] = NewSingleKeyData(5, 7)
	a.container['@'] = NewSingleKeyData(5, 6)
	a.container[':'] = NewSingleKeyData(5, 5)
	a.container['.'] = NewSingleKeyData(5, 4)
	a.container['-'] = NewSingleKeyData(5, 3)
	a.container['+'] = NewSingleKeyData(5, 0)
	a.container['0'] = NewSingleKeyData(4, 3)
	a.container['1'] = NewSingleKeyData(7, 0)
	a.container['2'] = NewSingleKeyData(7, 3)
	a.container['3'] = NewSingleKeyData(1, 0)
	a.container['4'] = NewSingleKeyData(1, 3)
	a.container['5'] = NewSingleKeyData(2, 0)
	a.container['6'] = NewSingleKeyData(2, 3)
	a.container['7'] = NewSingleKeyData(3, 0)
	a.container['8'] = NewSingleKeyData(3, 3)
	a.container['9'] = NewSingleKeyData(4, 0)
	a.container['A'] = NewSingleKeyData(1, 2)
	a.container['B'] = NewSingleKeyData(3, 4)
	a.container['C'] = NewSingleKeyData(2, 4)
	a.container['D'] = NewSingleKeyData(2, 2)
	a.container['E'] = NewSingleKeyData(1, 6)
	a.container['F'] = NewSingleKeyData(2, 5)
	a.container['G'] = NewSingleKeyData(3, 2)
	a.container['H'] = NewSingleKeyData(3, 5)
	a.container['I'] = NewSingleKeyData(4, 1)
	a.container['J'] = NewSingleKeyData(4, 2)
	a.container['K'] = NewSingleKeyData(4, 5)
	a.container['L'] = NewSingleKeyData(5, 2)
	a.container['M'] = NewSingleKeyData(4, 4)
	a.container['N'] = NewSingleKeyData(4, 7)
	a.container['O'] = NewSingleKeyData(4, 6)
	a.container['P'] = NewSingleKeyData(5, 1)
	a.container['Q'] = NewSingleKeyData(7, 6)
	a.container['R'] = NewSingleKeyData(2, 1)
	a.container['S'] = NewSingleKeyData(1, 5)
	a.container['T'] = NewSingleKeyData(2, 6)
	a.container['U'] = NewSingleKeyData(3, 6)
	a.container['V'] = NewSingleKeyData(3, 7)
	a.container['W'] = NewSingleKeyData(1, 1)
	a.container['X'] = NewSingleKeyData(2, 7)
	a.container['Y'] = NewSingleKeyData(3, 1)
	a.container['Z'] = NewSingleKeyData(1, 4)
	a.container['a'] = NewSingleKeyData(1, 2)
	a.container['b'] = NewSingleKeyData(3, 4)
	a.container['c'] = NewSingleKeyData(2, 4)
	a.container['d'] = NewSingleKeyData(2, 2)
	a.container['e'] = NewSingleKeyData(1, 6)
	a.container['f'] = NewSingleKeyData(2, 5)
	a.container['g'] = NewSingleKeyData(3, 2)
	a.container['h'] = NewSingleKeyData(3, 5)
	a.container['i'] = NewSingleKeyData(4, 1)
	a.container['j'] = NewSingleKeyData(4, 2)
	a.container['k'] = NewSingleKeyData(4, 5)
	a.container['l'] = NewSingleKeyData(5, 2)
	a.container['m'] = NewSingleKeyData(4, 4)
	a.container['n'] = NewSingleKeyData(4, 7)
	a.container['o'] = NewSingleKeyData(4, 6)
	a.container['p'] = NewSingleKeyData(5, 1)
	a.container['q'] = NewSingleKeyData(7, 6)
	a.container['r'] = NewSingleKeyData(2, 1)
	a.container['s'] = NewSingleKeyData(1, 5)
	a.container['t'] = NewSingleKeyData(2, 6)
	a.container['u'] = NewSingleKeyData(3, 6)
	a.container['v'] = NewSingleKeyData(3, 7)
	a.container['w'] = NewSingleKeyData(1, 1)
	a.container['x'] = NewSingleKeyData(2, 7)
	a.container['y'] = NewSingleKeyData(3, 1)
	a.container['z'] = NewSingleKeyData(1, 4)
	return a
}

func (a *Ascii) FromAscii(v uint8) *KeyData {
	out := a.container[v]
	return out
}

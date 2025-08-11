package formats

import (
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

type Buffer []byte

func (b *Buffer) Write(p []byte) {
	if len(*b)+len(p) > objects.MaxStringLen {
		panic(errors.ErrStringLimit)
	}

	*b = append(*b, p...)
}

func (b *Buffer) WriteString(s string) {
	if len(*b)+len(s) > objects.MaxStringLen {
		panic(errors.ErrStringLimit)
	}

	*b = append(*b, s...)
}

func (b *Buffer) WriteSingleByte(c byte) {
	if len(*b) >= objects.MaxStringLen {
		panic(errors.ErrStringLimit)
	}

	*b = append(*b, c)
}

func (b *Buffer) WriteRune(r rune) {
	if len(*b)+utf8.RuneLen(r) > objects.MaxStringLen {
		panic(errors.ErrStringLimit)
	}

	if r < utf8.RuneSelf {
		*b = append(*b, byte(r))
		return
	}

	b2 := *b
	n := len(b2)
	for n+utf8.UTFMax > cap(b2) {
		b2 = append(b2, 0)
	}
	w := utf8.EncodeRune(b2[n:n+utf8.UTFMax], r)
	*b = b2[:n+w]
}

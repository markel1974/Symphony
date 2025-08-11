package formats

import (
	"sync"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var _ppFree = sync.Pool{
	New: func() interface{} { return new(PrinterState) },
}

// PrinterState is used to store a printer's state and is reused with sync.Pool to avoid
// allocations.
type PrinterState struct {
	buf Buffer

	// arg holds the current item.
	arg objects.Object

	// fmt is used to format basic items such as integers or strings.
	fmt Formatter

	// reordered records whether the format string used argument reordering.
	reordered bool

	// goodArgNum records whether the most recent reordering directive was
	// valid.
	goodArgNum bool

	// erroring is set when printing an error string to guard against calling
	// handleMethods.
	erroring bool
}

// newPrinter allocates a new PrinterState struct or grabs a cached one.
func newPrinter() *PrinterState {
	p := _ppFree.Get().(*PrinterState)
	p.erroring = false
	p.fmt.init(&p.buf)
	return p
}

// free saves used PrinterState structs in ppFree; avoids an allocation per invocation.
func (p *PrinterState) free() {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized Buffer, we add a hard limit on the maximum
	// Buffer to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if cap(p.buf) > 64<<10 {
		return
	}

	p.buf = p.buf[:0]
	p.arg = nil
	_ppFree.Put(p)
}

func (p *PrinterState) Width() (wid int, ok bool) {
	return p.fmt.wid, p.fmt.widPresent
}

func (p *PrinterState) Precision() (prec int, ok bool) {
	return p.fmt.precision, p.fmt.precisionPresent
}

func (p *PrinterState) Flag(b int) bool {
	switch b {
	case '-':
		return p.fmt.minus
	case '+':
		return p.fmt.plus || p.fmt.plusV
	case '#':
		return p.fmt.sharp || p.fmt.sharpV
	case ' ':
		return p.fmt.space
	case '0':
		return p.fmt.zero
	}
	return false
}

func (p *PrinterState) Write(b []byte) (ret int, err error) {
	p.buf.Write(b)
	return len(b), nil
}

func (p *PrinterState) WriteString(s string) (ret int, err error) {
	p.buf.WriteString(s)
	return len(s), nil
}

func (p *PrinterState) WriteRune(r rune) (ret int, err error) {
	p.buf.WriteRune(r)
	return utf8.RuneLen(r), nil
}

func (p *PrinterState) WriteSingleByte(c byte) (ret int, err error) {
	p.buf.WriteSingleByte(c)
	return 1, nil
}

func (p *PrinterState) badVerb(verb rune) {
	p.erroring = true
	_, _ = p.WriteString(percentBangString)
	_, _ = p.WriteRune(verb)
	_, _ = p.WriteSingleByte('(')
	switch {
	case p.arg != nil:
		_, _ = p.WriteString(p.arg.String())
		_, _ = p.WriteSingleByte('=')
		p.printArg(p.arg, 'v')
	default:
		_, _ = p.WriteString(objects.UndefinedValue.String())
	}
	_, _ = p.WriteSingleByte(')')
	p.erroring = false
}

func (p *PrinterState) fmtBool(v bool, verb rune) {
	switch verb {
	case 't', 'v':
		p.fmt.fmtBoolean(v)
	default:
		p.badVerb(verb)
	}
}

// fmt0x64 formats a uint64 in hexadecimal and prefixes it with 0x or
// not, as requested, by temporarily setting the sharp flag.
func (p *PrinterState) fmt0x64(v uint64, leading0x bool) {
	sharp := p.fmt.sharp
	p.fmt.sharp = leading0x
	p.fmt.fmtInteger(v, 16, unsigned, 'v', lDigits)
	p.fmt.sharp = sharp
}

// fmtInteger formats a signed or unsigned integer.
func (p *PrinterState) fmtInteger(v uint64, isSigned bool, verb rune) {
	switch verb {
	case 'v':
		if p.fmt.sharpV && !isSigned {
			p.fmt0x64(v, true)
		} else {
			p.fmt.fmtInteger(v, 10, isSigned, verb, lDigits)
		}
	case 'd':
		p.fmt.fmtInteger(v, 10, isSigned, verb, lDigits)
	case 'b':
		p.fmt.fmtInteger(v, 2, isSigned, verb, lDigits)
	case 'o', 'O':
		p.fmt.fmtInteger(v, 8, isSigned, verb, lDigits)
	case 'x':
		p.fmt.fmtInteger(v, 16, isSigned, verb, lDigits)
	case 'X':
		p.fmt.fmtInteger(v, 16, isSigned, verb, uDigits)
	case 'c':
		p.fmt.fmtC(v)
	case 'q':
		if v <= utf8.MaxRune {
			p.fmt.fmtQc(v)
		} else {
			p.badVerb(verb)
		}
	case 'U':
		p.fmt.fmtUnicode(v)
	default:
		p.badVerb(verb)
	}
}

// fmtFloat formats a float. The default precision for each verb
// is specified as last argument in the call to fmt_float.
func (p *PrinterState) fmtFloat(v float64, size int, verb rune) {
	switch verb {
	case 'v':
		p.fmt.fmtFloat(v, size, 'g', -1)
	case 'b', 'g', 'G', 'x', 'X':
		p.fmt.fmtFloat(v, size, verb, -1)
	case 'f', 'e', 'E':
		p.fmt.fmtFloat(v, size, verb, 6)
	case 'F':
		p.fmt.fmtFloat(v, size, 'f', 6)
	default:
		p.badVerb(verb)
	}
}

func (p *PrinterState) fmtString(v string, verb rune) {
	switch verb {
	case 'v':
		if p.fmt.sharpV {
			p.fmt.fmtQ(v)
		} else {
			p.fmt.fmtS(v)
		}
	case 's':
		p.fmt.fmtS(v)
	case 'x':
		p.fmt.fmtSx(v, lDigits)
	case 'X':
		p.fmt.fmtSx(v, uDigits)
	case 'q':
		p.fmt.fmtQ(v)
	default:
		p.badVerb(verb)
	}
}

func (p *PrinterState) fmtBytes(v []byte, verb rune, typeString string) {
	switch verb {
	case 'v', 'd':
		if p.fmt.sharpV {
			_, _ = p.WriteString(typeString)
			if v == nil {
				_, _ = p.WriteString(nilParenString)
				return
			}
			_, _ = p.WriteSingleByte('{')
			for i, c := range v {
				if i > 0 {
					_, _ = p.WriteString(commaSpaceString)
				}
				p.fmt0x64(uint64(c), true)
			}
			_, _ = p.WriteSingleByte('}')
		} else {
			_, _ = p.WriteSingleByte('[')
			for i, c := range v {
				if i > 0 {
					_, _ = p.WriteSingleByte(' ')
				}
				p.fmt.fmtInteger(uint64(c), 10, unsigned, verb, lDigits)
			}
			_, _ = p.WriteSingleByte(']')
		}
	case 's':
		p.fmt.fmtBs(v)
	case 'x':
		p.fmt.fmtBx(v, lDigits)
	case 'X':
		p.fmt.fmtBx(v, uDigits)
	case 'q':
		p.fmt.fmtQ(string(v))
	}
}

func (p *PrinterState) printArg(arg objects.Object, verb rune) {
	p.arg = arg

	if arg == nil {
		arg = objects.UndefinedValue
	}

	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do
	// them first.
	switch verb {
	case 'T':
		p.fmt.fmtS(arg.TypeName())
		return
	case 'v':
		p.fmt.fmtS(arg.String())
		return
	}

	// Some types can be done without reflection.
	switch f := arg.(type) {
	case *objects.Bool:
		p.fmtBool(!f.IsFalsy(), verb)
	case *objects.Float:
		p.fmtFloat(f.Value, 64, verb)
	case *objects.Int:
		p.fmtInteger(uint64(f.Value), signed, verb)
	case *objects.String:
		p.fmtString(f.Value, verb)
	case *objects.Bytes:
		p.fmtBytes(f.Value, verb, "[]byte")
	default:
		p.fmtString(f.String(), verb)
	}
}

// argNumber returns the next argument to evaluate, which is either the value
// of the passed-in argNum or the value of the bracketed integer that begins
// format[i:]. It also returns the new value of i, that is, the index of the
// next byte of the format to process.
func (p *PrinterState) argNumber(
	argNum int,
	format string,
	i int,
	numArgs int,
) (newArgNum, newi int, found bool) {
	if len(format) <= i || format[i] != '[' {
		return argNum, i, false
	}
	p.reordered = true
	index, wid, ok := parseArgNumber(format[i:])
	if ok && 0 <= index && index < numArgs {
		return index, i + wid, true
	}
	p.goodArgNum = false
	return argNum, i + wid, ok
}

func (p *PrinterState) badArgNum(verb rune) {
	_, _ = p.WriteString(percentBangString)
	_, _ = p.WriteRune(verb)
	_, _ = p.WriteString(badIndexString)
}

func (p *PrinterState) missingArg(verb rune) {
	_, _ = p.WriteString(percentBangString)
	_, _ = p.WriteRune(verb)
	_, _ = p.WriteString(missingString)
}

func (p *PrinterState) doFormat(format string, a []objects.Object) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok && e == errors.ErrStringLimit {
				err = e
				return
			}
			panic(r)
		}
	}()

	end := len(format)
	argNum := 0         // we process one argument per non-trivial format
	afterIndex := false // previous item in format was an index like [3].
	p.reordered = false
formatLoop:
	for i := 0; i < end; {
		p.goodArgNum = true
		lastI := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lastI {
			_, _ = p.WriteString(format[lastI:i])
		}
		if i >= end {
			// done processing format string
			break
		}

		// Process one verb
		i++

		// Do we have flags?
		p.fmt.clearFlags()
	simpleFormat:
		for ; i < end; i++ {
			c := format[i]
			switch c {
			case '#':
				p.fmt.sharp = true
			case '0':
				// Only allow zero padding to the left.
				p.fmt.zero = !p.fmt.minus
			case '+':
				p.fmt.plus = true
			case '-':
				p.fmt.minus = true
				p.fmt.zero = false // Do not pad with zeros to the right.
			case ' ':
				p.fmt.space = true
			default:
				// Fast path for common case of ascii lower case simple verbs without precision or width or argument indices.
				if 'a' <= c && c <= 'z' && argNum < len(a) {
					if c == 'v' {
						// Go syntax
						p.fmt.sharpV = p.fmt.sharp
						p.fmt.sharp = false
						// Struct-field syntax
						p.fmt.plusV = p.fmt.plus
						p.fmt.plus = false
					}
					p.printArg(a[argNum], rune(c))
					argNum++
					i++
					continue formatLoop
				}
				// Format is more complex than simple flags and a verb or is
				// malformed.
				break simpleFormat
			}
		}

		// Do we have an explicit argument index?
		argNum, i, afterIndex = p.argNumber(argNum, format, i, len(a))

		// Do we have width?
		if i < end && format[i] == '*' {
			i++
			p.fmt.wid, p.fmt.widPresent, argNum = intFromArg(a, argNum)

			if !p.fmt.widPresent {
				_, _ = p.WriteString(badWidthString)
			}

			// We have a negative width, so take its value and ensure
			// that the minus flag is set
			if p.fmt.wid < 0 {
				p.fmt.wid = -p.fmt.wid
				p.fmt.minus = true
				p.fmt.zero = false // Do not pad with zeros to the right.
			}
			afterIndex = false
		} else {
			p.fmt.wid, p.fmt.widPresent, i = parseNum(format, i, end)
			if afterIndex && p.fmt.widPresent { // "%[3]2d"
				p.goodArgNum = false
			}
		}

		// Do we have precision?
		if i+1 < end && format[i] == '.' {
			i++
			if afterIndex { // "%[3].2d"
				p.goodArgNum = false
			}
			argNum, i, afterIndex = p.argNumber(argNum, format, i, len(a))
			if i < end && format[i] == '*' {
				i++
				p.fmt.precision, p.fmt.precisionPresent, argNum = intFromArg(a, argNum)
				// Negative precision arguments don't make sense
				if p.fmt.precision < 0 {
					p.fmt.precision = 0
					p.fmt.precisionPresent = false
				}
				if !p.fmt.precisionPresent {
					_, _ = p.WriteString(badPrecisionString)
				}
				afterIndex = false
			} else {
				p.fmt.precision, p.fmt.precisionPresent, i = parseNum(format, i, end)
				if !p.fmt.precisionPresent {
					p.fmt.precision = 0
					p.fmt.precisionPresent = true
				}
			}
		}

		if !afterIndex {
			argNum, i, afterIndex = p.argNumber(argNum, format, i, len(a))
		}

		if i >= end {
			_, _ = p.WriteString(noVerbString)
			break
		}

		verb, size := rune(format[i]), 1
		if verb >= utf8.RuneSelf {
			verb, size = utf8.DecodeRuneInString(format[i:])
		}
		i += size

		switch {
		case verb == '%':
			// Percent does not absorb operands and ignores f.wid and f.precision.
			_, _ = p.WriteSingleByte('%')
		case !p.goodArgNum:
			p.badArgNum(verb)
		case argNum >= len(a):
			// No argument left over to print for the current verb.
			p.missingArg(verb)
		case verb == 'v':
			// Go syntax
			p.fmt.sharpV = p.fmt.sharp
			p.fmt.sharp = false
			// Struct-field syntax
			p.fmt.plusV = p.fmt.plus
			p.fmt.plus = false
			fallthrough
		default:
			p.printArg(a[argNum], verb)
			argNum++
		}
	}

	// Check for extra arguments unless the call accessed the arguments
	// out of order, in which case it's too expensive to detect if they've all
	// been used and arguably OK if they're not.
	if !p.reordered && argNum < len(a) {
		p.fmt.clearFlags()
		_, _ = p.WriteString(extraString)
		for i, arg := range a[argNum:] {
			if i > 0 {
				_, _ = p.WriteString(commaSpaceString)
			}
			if arg == nil {
				_, _ = p.WriteString(objects.UndefinedValue.String())
			} else {
				_, _ = p.WriteString(arg.TypeName())
				_, _ = p.WriteSingleByte('=')
				p.printArg(arg, 'v')
			}
		}
		_, _ = p.WriteSingleByte(')')
	}

	return nil
}

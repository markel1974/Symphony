package formats

import (
	"strconv"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// maxIntegerMagnitude defines the maximum allowable magnitude for integer values used in certain operations.
const (
	maxIntegerMagnitude int = 1e6
)

// commaSpaceString represents a comma followed by a space, commonly used for separating elements in strings.
// nilParenString represents the string "(nil)", typically used to denote a nil value in formatted output.
// percentBangString represents "%!", used as a prefix for formatting errors in strings.
// missingString represents the string "(MISSING)", used when a placeholder value is missing in format operations.
// badIndexString represents the string "(BADINDEX)", used when an invalid index is encountered in format operations.
// extraString represents "%!(EXTRA ", used to denote extra arguments in formatted string operations.
// badWidthString represents "%!(BADWIDTH)", used for indicating an invalid width in formatting operations.
// badPrecisionString represents "%!(BADPREC)", used to denote an invalid precision in formatting operations.
// noVerbString represents "%!(NOVERB)", used when no verb is specified in formatted string operations.
const (
	commaSpaceString   = ", "
	nilParenString     = "(nil)"
	percentBangString  = "%!"
	missingString      = "(MISSING)"
	badIndexString     = "(BADINDEX)"
	extraString        = "%!(EXTRA "
	badWidthString     = "%!(BADWIDTH)"
	badPrecisionString = "%!(BADPREC)"
	noVerbString       = "%!(NOVERB)"
)

// lDigits represents lowercase hexadecimal digits and the prefix 'x' used in string formatting or parsing.
// uDigits represents uppercase hexadecimal digits and the prefix 'X' used in string formatting or parsing.
const (
	lDigits = "0123456789abcdefx"
	uDigits = "0123456789ABCDEFX"
)

// signed indicates whether a value is signed, set to true.
// unsigned indicates whether a value is unsigned, set to false.
const (
	signed   = true
	unsigned = false
)

// Formatter is a type for managing formatting operations such as padding, width, precision, and flag handling.
type Formatter struct {
	fmtFlags
	buf       *Buffer
	wid       int
	precision int
	intBuf    [68]byte
}

// clearFlags resets the formatter's flags to their default state by clearing all formatting options.
func (f *Formatter) clearFlags() {
	f.fmtFlags = fmtFlags{}
}

// init initializes the Formatter with the provided Buffer and resets its state by clearing existing formatting flags.
func (f *Formatter) init(buf *Buffer) {
	f.buf = buf
	f.clearFlags()
}

// writePadding writes n padding bytes into the buffer, using spaces by default or zeros if the zero flag is set.
// If n is zero or negative, the method does nothing. Panics if the resulting buffer exceeds the allowed maximum length.
func (f *Formatter) writePadding(n int) {
	if n <= 0 { // No padding bytes needed.
		return
	}
	buf := *f.buf
	oldLen := len(buf)
	newLen := oldLen + n
	if newLen > objects.MaxStringLen {
		panic(errors.ErrStringLimit)
	}
	if newLen > cap(buf) {
		buf = make(Buffer, cap(buf)*2+n)
		copy(buf, *f.buf)
	}
	padByte := byte(' ')
	if f.zero {
		padByte = byte('0')
	}
	padding := buf[oldLen:newLen]
	for i := range padding {
		padding[i] = padByte
	}
	*f.buf = buf[:newLen]
}

// pad applies width-padding to the given byte slice and writes it to the buffer based on the formatting flags.
func (f *Formatter) pad(b []byte) {
	if !f.widPresent || f.wid == 0 {
		f.buf.Write(b)
		return
	}
	width := f.wid - utf8.RuneCount(b)
	if !f.minus {
		// left padding
		f.writePadding(width)
		f.buf.Write(b)
	} else {
		// right padding
		f.buf.Write(b)
		f.writePadding(width)
	}
}

// padString adjusts the given string `s` according to the formatter's width and alignment settings, adding padding if needed.
func (f *Formatter) padString(s string) {
	if !f.widPresent || f.wid == 0 {
		f.buf.WriteString(s)
		return
	}
	width := f.wid - utf8.RuneCountInString(s)
	if !f.minus {
		f.writePadding(width)
		f.buf.WriteString(s)
	} else {
		f.buf.WriteString(s)
		f.writePadding(width)
	}
}

// fmtBoolean formats a boolean value and writes "true" or "false" to the buffer with appropriate padding.
func (f *Formatter) fmtBoolean(v bool) {
	if v {
		f.padString("true")
	} else {
		f.padString("false")
	}
}

// fmtUnicode formats a Unicode code point as a hexadecimal string with optional padding and precision adjustments.
func (f *Formatter) fmtUnicode(u uint64) {
	buf := f.intBuf[0:]
	precision := 4
	if f.precisionPresent && f.precision > 4 {
		precision = f.precision
		width := 2 + precision + 2 + utf8.UTFMax + 1
		if width > len(buf) {
			buf = make([]byte, width)
		}
	}
	i := len(buf)
	if f.sharp && u <= utf8.MaxRune && strconv.IsPrint(rune(u)) {
		i--
		buf[i] = '\''
		i -= utf8.RuneLen(rune(u))
		utf8.EncodeRune(buf[i:], rune(u))
		i--
		buf[i] = '\''
		i--
		buf[i] = ' '
	}
	for u >= 16 {
		i--
		buf[i] = uDigits[u&0xF]
		precision--
		u >>= 4
	}
	i--
	buf[i] = uDigits[u]
	precision--
	for precision > 0 {
		i--
		buf[i] = '0'
		precision--
	}
	i--
	buf[i] = '+'
	i--
	buf[i] = 'U'

	oldZero := f.zero
	f.zero = false
	f.pad(buf[i:])
	f.zero = oldZero
}

// fmtInteger formats an unsigned integer based on the specified base, width, precision, and formatting flags.
func (f *Formatter) fmtInteger(u uint64, base int, isSigned bool, verb rune, digits string) {
	negative := isSigned && int64(u) < 0
	if negative {
		u = -u
	}

	buf := f.intBuf[0:]
	if f.widPresent || f.precisionPresent {
		width := 3 + f.wid + f.precision
		if width > len(buf) {
			buf = make([]byte, width)
		}
	}

	precision := 0
	if f.precisionPresent {
		precision = f.precision
		if precision == 0 && u == 0 {
			oldZero := f.zero
			f.zero = false
			f.writePadding(f.wid)
			f.zero = oldZero
			return
		}
	} else if f.zero && f.widPresent {
		precision = f.wid
		if negative || f.plus || f.space {
			precision--
		}
	}
	i := len(buf)
	switch base {
	case 10:
		for u >= 10 {
			i--
			next := u / 10
			buf[i] = byte('0' + u - next*10)
			u = next
		}
	case 16:
		for u >= 16 {
			i--
			buf[i] = digits[u&0xF]
			u >>= 4
		}
	case 8:
		for u >= 8 {
			i--
			buf[i] = byte('0' + u&7)
			u >>= 3
		}
	case 2:
		for u >= 2 {
			i--
			buf[i] = byte('0' + u&1)
			u >>= 1
		}
	default:
		panic("fmt: unknown base; can't happen")
	}
	i--
	buf[i] = digits[u]
	for i > 0 && precision > len(buf)-i {
		i--
		buf[i] = '0'
	}
	if f.sharp {
		switch base {
		case 2:
			i--
			buf[i] = 'b'
			i--
			buf[i] = '0'
		case 8:
			if buf[i] != '0' {
				i--
				buf[i] = '0'
			}
		case 16:
			i--
			buf[i] = digits[16]
			i--
			buf[i] = '0'
		}
	}
	if verb == 'O' {
		i--
		buf[i] = 'o'
		i--
		buf[i] = '0'
	}

	if negative {
		i--
		buf[i] = '-'
	} else if f.plus {
		i--
		buf[i] = '+'
	} else if f.space {
		i--
		buf[i] = ' '
	}

	oldZero := f.zero
	f.zero = false
	f.pad(buf[i:])
	f.zero = oldZero
}

// truncateString truncates the provided string to a length determined by the Formatter's precision if precision is present.
func (f *Formatter) truncateString(s string) string {
	if f.precisionPresent {
		n := f.precision
		for i := range s {
			n--
			if n < 0 {
				return s[:i]
			}
		}
	}
	return s
}

// truncate trims the input byte slice `b` to a specified precision if `precisionPresent` is true.
func (f *Formatter) truncate(b []byte) []byte {
	if f.precisionPresent {
		n := f.precision
		for i := 0; i < len(b); {
			n--
			if n < 0 {
				return b[:i]
			}
			wid := 1
			if b[i] >= utf8.RuneSelf {
				_, wid = utf8.DecodeRune(b[i:])
			}
			i += wid
		}
	}
	return b
}

// fmtS formats a string by truncating it based on precision and applying width-based padding.
func (f *Formatter) fmtS(s string) {
	s = f.truncateString(s)
	f.padString(s)
}

// fmtBs formats a byte slice by truncating its content based on precision and applying width-based padding.
func (f *Formatter) fmtBs(b []byte) {
	b = f.truncate(b)
	f.pad(b)
}

// fmtSbx encodes a string or byte slice into a formatted hexadecimal representation according to the given flags and digits.
func (f *Formatter) fmtSbx(s string, b []byte, digits string) {
	length := len(b)
	if b == nil {
		length = len(s)
	}
	if f.precisionPresent && f.precision < length {
		length = f.precision
	}
	width := 2 * length
	if width > 0 {
		if f.space {
			if f.sharp {
				width *= 2
			}
			width += length - 1
		} else if f.sharp {
			width += 2
		}
	} else {
		if f.widPresent {
			f.writePadding(f.wid)
		}
		return
	}
	if f.widPresent && f.wid > width && !f.minus {
		f.writePadding(f.wid - width)
	}
	buf := *f.buf
	if f.sharp {
		buf = append(buf, '0', digits[16])
	}
	var c byte
	for i := 0; i < length; i++ {
		if f.space && i > 0 {
			buf = append(buf, ' ')
			if f.sharp {
				buf = append(buf, '0', digits[16])
			}
		}
		if b != nil {
			c = b[i]
		} else {
			c = s[i]
		}
		buf = append(buf, digits[c>>4], digits[c&0xF])
	}
	*f.buf = buf
	if f.widPresent && f.wid > width && f.minus {
		f.writePadding(f.wid - width)
	}
}

// fmtSx formats a string `s` according to the specified hexadecimal digits `digits`. It delegates to fmtSbx.
func (f *Formatter) fmtSx(s, digits string) {
	f.fmtSbx(s, nil, digits)
}

// fmtBx formats a byte slice using a specified set of digits and delegates formatting to fmtSbx with an empty prefix.
func (f *Formatter) fmtBx(b []byte, digits string) {
	f.fmtSbx("", b, digits)
}

// fmtQ formats a string as a quoted string, handling truncation, ASCII encoding, and proper padding based on flags.
func (f *Formatter) fmtQ(s string) {
	s = f.truncateString(s)
	if f.sharp && strconv.CanBackquote(s) {
		f.padString("`" + s + "`")
		return
	}
	buf := f.intBuf[:0]
	if f.plus {
		f.pad(strconv.AppendQuoteToASCII(buf, s))
	} else {
		f.pad(strconv.AppendQuote(buf, s))
	}
}

// fmtC formats the given Unicode code point (as uint64) into its UTF-8 encoded form and applies padding if necessary.
func (f *Formatter) fmtC(c uint64) {
	r := rune(c)
	if c > utf8.MaxRune {
		r = utf8.RuneError
	}
	buf := f.intBuf[:0]
	w := utf8.EncodeRune(buf[:utf8.UTFMax], r)
	f.pad(buf[:w])
}

// fmtQc formats a quoted character using its rune representation, handling ASCII quoting with optional '+' formatting.
func (f *Formatter) fmtQc(c uint64) {
	r := rune(c)
	if c > utf8.MaxRune {
		r = utf8.RuneError
	}
	buf := f.intBuf[:0]
	if f.plus {
		f.pad(strconv.AppendQuoteRuneToASCII(buf, r))
	} else {
		f.pad(strconv.AppendQuoteRune(buf, r))
	}
}

// fmtFloat formats a floating-point number v according to the specified size, verb, and precision parameters.
func (f *Formatter) fmtFloat(v float64, size int, verb rune, precision int) {
	if f.precisionPresent {
		precision = f.precision
	}
	num := strconv.AppendFloat(f.intBuf[:1], v, byte(verb), precision, size)
	if num[1] == '-' || num[1] == '+' {
		num = num[1:]
	} else {
		num[0] = '+'
	}
	if f.space && num[0] == '+' && !f.plus {
		num[0] = ' '
	}
	if num[1] == 'I' || num[1] == 'N' {
		oldZero := f.zero
		f.zero = false
		if num[1] == 'N' && !f.space && !f.plus {
			num = num[1:]
		}
		f.pad(num)
		f.zero = oldZero
		return
	}
	if f.sharp && verb != 'b' {
		digits := 0
		switch verb {
		case 'v', 'g', 'G', 'x':
			digits = precision
			if digits == -1 {
				digits = 6
			}
		}
		var tailBuf [6]byte
		tail := tailBuf[:0]

		hasDecimalPoint := false
		for i := 1; i < len(num); i++ {
			switch num[i] {
			case '.':
				hasDecimalPoint = true
			case 'p', 'P':
				tail = append(tail, num[i:]...)
				num = num[:i]
			case 'e', 'E':
				if verb != 'x' && verb != 'X' {
					tail = append(tail, num[i:]...)
					num = num[:i]
					break
				}
				fallthrough
			default:
				digits--
			}
		}
		if !hasDecimalPoint {
			num = append(num, '.')
		}
		for digits > 0 {
			num = append(num, '0')
			digits--
		}
		num = append(num, tail...)
	}
	if f.plus || num[0] != '+' {
		if f.zero && f.widPresent && f.wid > len(num) {
			f.buf.WriteSingleByte(num[0])
			f.writePadding(f.wid - len(num))
			f.buf.Write(num[1:])
			return
		}
		f.pad(num)
		return
	}
	f.pad(num[1:])
}

// Format formats a string using the specified format string and arguments implementing the IObject interface.
// Returns the formatted string and an error, if any occurred during formatting.
func Format(format string, a ...objects.IObject) (string, error) {
	p := newPrinter()
	err := p.doFormat(format, a)
	s := string(p.buf)
	p.free()
	return s, err
}

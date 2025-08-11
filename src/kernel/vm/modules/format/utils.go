package format

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// tooLarge reports whether the magnitude of the integer is
// too large to be used as a formatting width or precision.
func tooLarge(x int) bool {
	return x > maxIntegerMagnitude || x < -maxIntegerMagnitude
}

// parseNum converts ASCII to integer.  num is 0 (and isNum is false) if no
// number present.
func parseNum(s string, start, end int) (int, bool, int) {
	if start >= end {
		return 0, false, end
	}
	num := 0
	isNum := false
	newI := start
	for ; newI < end && '0' <= s[newI] && s[newI] <= '9'; newI++ {
		if tooLarge(num) {
			return 0, false, end // Overflow; crazy long number most likely.
		}
		num = num*10 + int(s[newI]-'0')
		isNum = true
	}
	return num, isNum, newI
}

func intFromArg(a []objects.IObject, argNum int) (int, bool, int) {
	num := 0
	isInt := false
	newArgNum := argNum
	if argNum < len(a) {
		var num64 int64
		num64, isInt = objects.ToInt64(a[argNum])
		num = int(num64)
		newArgNum = argNum + 1
		if tooLarge(num) {
			num = 0
			isInt = false
		}
	}
	return num, isInt, newArgNum
}

// parseArgNumber returns the value of the bracketed number, minus 1
// (explicit argument numbers are one-indexed but we want zero-indexed).
// The opening bracket is known to be present at format[0].
// The returned values are the index, the number of bytes to consume
// up to the closing paren, if present, and whether the number parsed
// ok. The bytes to consume will be 1 if no closing paren is present.
func parseArgNumber(format string) (int, int, bool) {
	// There must be at least 3 bytes: [n].
	if len(format) < 3 {
		return 0, 1, false
	}
	// Find closing bracket.
	for i := 1; i < len(format); i++ {
		if format[i] == ']' {
			width, ok, newI := parseNum(format, 1, i)
			if !ok || newI != i {
				return 0, i + 1, false
			}
			// arg numbers are one-indexed and skip paren.
			return width - 1, i + 1, true
		}
	}
	return 0, 1, false
}

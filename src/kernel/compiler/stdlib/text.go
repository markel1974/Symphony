package stdlib

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// textModule is a map of various text processing functions, including string manipulation, regex operations, and type conversions.
var _textModule = map[string]objects.IObject{
	"re_match":       objects.NewFunctionModule("re_match", textREMatch),                                   // re_match(pattern, text) => bool/error
	"re_find":        objects.NewFunctionModule("re_find", textREFind),                                     // re_find(pattern, text, count) => [[{text:,begin:,end:}]]/undefined
	"re_replace":     objects.NewFunctionModule("re_replace", textREReplace),                               // re_replace(pattern, text, repl) => string/error
	"re_split":       objects.NewFunctionModule("re_split", textRESplit),                                   // re_split(pattern, text, count) => [string]/error
	"re_compile":     objects.NewFunctionModule("re_compile", textRECompile),                               // re_compile(pattern) => Regexp/error
	"compare":        objects.NewFunctionModule("compare", objects.FuncIssOi(strings.Compare)),             // compare(a, b) => int
	"contains":       objects.NewFunctionModule("contains", objects.FuncIssOb(strings.Contains)),           // contains(s, substr) => bool
	"contains_any":   objects.NewFunctionModule("contains_any", objects.FuncIssOb(strings.ContainsAny)),    // contains_any(s, chars) => bool
	"count":          objects.NewFunctionModule("count", objects.FuncIssOi(strings.Count)),                 // count(s, substr) => int
	"equal_fold":     objects.NewFunctionModule("equal_fold", objects.FuncIssOb(strings.EqualFold)),        // "equal_fold(s, t) => bool
	"fields":         objects.NewFunctionModule("fields", objects.FuncIsOsS(strings.Fields)),               // fields(s) => [string]
	"has_prefix":     objects.NewFunctionModule("has_prefix", objects.FuncIssOb(strings.HasPrefix)),        // has_prefix(s, prefix) => bool
	"has_suffix":     objects.NewFunctionModule("has_suffix", objects.FuncIssOb(strings.HasSuffix)),        // has_suffix(s, suffix) => bool
	"index":          objects.NewFunctionModule("index", objects.FuncIssOi(strings.Index)),                 // index(s, substr) => int
	"index_any":      objects.NewFunctionModule("index_any", objects.FuncIssOi(strings.IndexAny)),          // index_any(s, chars) => int
	"join":           objects.NewFunctionModule("join", textJoin),                                          // join(arr, sep) => string
	"last_index":     objects.NewFunctionModule("last_index", objects.FuncIssOi(strings.LastIndex)),        // last_index(s, substr) => int
	"last_index_any": objects.NewFunctionModule("last_index_any", objects.FuncIssOi(strings.LastIndexAny)), // last_index_any(s, chars) => int
	"repeat":         objects.NewFunctionModule("repeat", textRepeat),                                      // repeat(s, count) => string
	"replace":        objects.NewFunctionModule("replace", textReplace),                                    // replace(s, old, new, n) => string
	"substr":         objects.NewFunctionModule("substr", textSubstring),                                   // substr(s, lower, upper) => string
	"split":          objects.NewFunctionModule("split", objects.FuncIssOsS(strings.Split)),                // split(s, sep) => [string]
	"split_after":    objects.NewFunctionModule("split_after", objects.FuncIssOsS(strings.SplitAfter)),     // split_after(s, sep) => [string]
	"split_after_n":  objects.NewFunctionModule("split_after_n", objects.FuncIssiOsS(strings.SplitAfterN)), // split_after_n(s, sep, n) => [string]
	"split_n":        objects.NewFunctionModule("split_n", objects.FuncIssiOsS(strings.SplitN)),            // split_n(s, sep, n) => [string]
	"title":          objects.NewFunctionModule("title", objects.FuncIsOs(strings.Title)),                  // title(s) => string
	"to_lower":       objects.NewFunctionModule("to_lower", objects.FuncIsOs(strings.ToLower)),             // to_lower(s) => string
	"to_title":       objects.NewFunctionModule("to_title", objects.FuncIsOs(strings.ToTitle)),             // to_title(s) => string
	"to_upper":       objects.NewFunctionModule("to_upper", objects.FuncIsOs(strings.ToUpper)),             // to_upper(s) => string
	"pad_left":       objects.NewFunctionModule("pad_left", textPadLeft),                                   // pad_left(s, pad_len, pad_with) => string
	"pad_right":      objects.NewFunctionModule("pad_right", textPadRight),                                 // pad_right(s, pad_len, pad_with) => string
	"trim":           objects.NewFunctionModule("trim", objects.FuncIssOs(strings.Trim)),                   // trim(s, cutset) => string
	"trim_left":      objects.NewFunctionModule("trim_left", objects.FuncIssOs(strings.TrimLeft)),          // trim_left(s, cutset) => string
	"trim_prefix":    objects.NewFunctionModule("trim_prefix", objects.FuncIssOs(strings.TrimPrefix)),      // trim_prefix(s, prefix) => string
	"trim_right":     objects.NewFunctionModule("trim_right", objects.FuncIssOs(strings.TrimRight)),        // trim_right(s, cutset) => string
	"trim_space":     objects.NewFunctionModule("trim_space", objects.FuncIsOs(strings.TrimSpace)),         // trim_space(s) => string
	"trim_suffix":    objects.NewFunctionModule("trim_suffix", objects.FuncIssOs(strings.TrimSuffix)),      // trim_suffix(s, suffix) => string
	"atoi":           objects.NewFunctionModule("atoi", objects.FuncIsOie(strconv.Atoi)),                   // atoi(str) => int/error
	"format_bool":    objects.NewFunctionModule("format_bool", textFormatBool),                             // format_bool(b) => string
	"format_float":   objects.NewFunctionModule("format_float", textFormatFloat),                           // format_float(f, fmt, prec, bits) => string
	"format_int":     objects.NewFunctionModule("format_int", textFormatInt),                               // format_int(i, base) => string
	"itoa":           objects.NewFunctionModule("itoa", objects.FuncIiOs(strconv.Itoa)),                    // itoa(i) => string
	"parse_bool":     objects.NewFunctionModule("parse_bool", textParseBool),                               // parse_bool(str) => bool/error
	"parse_float":    objects.NewFunctionModule("parse_float", textParseFloat),                             // parse_float(str, bits) => float/error
	"parse_number":   objects.NewFunctionModule("parse_number", textParseNumber),                           // parse_float(str, bits) => float/error
	"parse_int":      objects.NewFunctionModule("parse_int", textParseInt),                                 // parse_int(str, base, bits) => int/error
	"quote":          objects.NewFunctionModule("quote", objects.FuncIsOs(strconv.Quote)),                  // quote(str) => string
	"unquote":        objects.NewFunctionModule("unquote", objects.FuncIsOse(strconv.Unquote)),             // unquote(str) => string/error
}

// textREMatch checks if the first argument (regex pattern) matches the second argument (string) and returns a boolean object.
// Returns an error if arguments are not string or if there is an invalid number of arguments.
func textREMatch(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	matched, err := regexp.MatchString(s1, s2)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	if matched {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// textREFind performs a regular expression match on a string and extracts matched substrings with their positions.
func textREFind(args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 2 && numArgs != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	if numArgs < 3 {
		m := re.FindStringSubmatchIndex(s2)
		if m == nil {
			return objects.UndefinedValue, nil
		}
		arr := objects.NewArray(nil)
		for i := 0; i < len(m); i += 2 {
			arr.Append(objects.NewMapImmutable(map[string]objects.IObject{
				"text":  objects.NewStringNoSize(s2[m[i]:m[i+1]]),
				"begin": objects.NewInt(int64(m[i])),
				"end":   objects.NewInt(int64(m[i+1])),
			}))
		}
		return objects.NewArray([]objects.IObject{arr}), nil
	}
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	mFA := re.FindAllStringSubmatchIndex(s2, int(i3))
	if mFA == nil {
		return objects.UndefinedValue, nil
	}
	arr := objects.NewArray(nil)
	for _, m := range mFA {
		subMatch := objects.NewArray(nil)
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(objects.NewMapImmutable(map[string]objects.IObject{
				"text":  objects.NewStringNoSize(s2[m[i]:m[i+1]]),
				"begin": objects.NewInt(int64(m[i])),
				"end":   objects.NewInt(int64(m[i+1])),
			}))
		}
		arr.Append(subMatch)
	}
	return arr, nil
}

// textREReplace performs a regular expression replacement in a string.
// Accepts three arguments: regex pattern as the first argument, source string as the second, and replacement string as the third.
// Returns the modified string or an error if the input arguments are invalid or the operation fails.
func textREReplace(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s3, err := objects.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	s := doTextRegexpReplace(re, s2, s3)
	v, err := objects.NewString(s)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// textRESplit splits a string by a regular expression pattern and returns an array of resulting substrings.
// The first argument is the regex pattern, the second is the string to split, and the optional third is the max splits.
// Returns an error if arguments' types or counts are incorrect, or if the regex compilation fails.
func textRESplit(args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 2 && numArgs != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	var i3 = int64(-1)
	if numArgs > 2 {
		i3, err = objects.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	arr := &objects.Array{}
	for _, s := range re.Split(s2, int(i3)) {
		v, err := objects.NewString(s)
		if err != nil {
			return nil, err
		}
		arr.Append(v)
	}
	return arr, nil
}

// textRECompile compiles a string argument into a regular expression and returns a map of regex operations or an error.
func textRECompile(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, objects.NewInvalidArgumentError(0, "string(compatible)", args[0].TypeName())
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return makeTextRegexp(re), nil
}

// textReplace performs a string replacement operation using 4 arguments: original string, target, replacement, and limit.
// It returns the modified string as an object or an error if the input arguments are invalid or exceed string limits.
func textReplace(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s3, err := objects.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	s, ok := doTextReplace(s1, s2, s3, int(i4))
	if !ok {
		return nil, objects.ErrStringLimit
	}
	return objects.NewString(s)
}

// textSubstring extracts a substring from a string, using start and optional end indices provided as arguments.
// Returns an error if arguments are invalid or indices are out of bounds.
func textSubstring(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	strlen := int64(len(s1))
	i3 := strlen
	if argsLen == 3 {
		i3, err = objects.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	if i2 > i3 {
		return nil, objects.ErrInvalidIndexType
	}
	if i2 < 0 {
		i2 = 0
	} else if i2 > strlen {
		i2 = strlen
	}
	if i3 < 0 {
		i3 = 0
	} else if i3 > strlen {
		i3 = strlen
	}
	return objects.NewString(s1[i2:i3])
}

// textPadLeft pads a string from the left with a specified character or default space until the string reaches a given length.
// Returns an error if arguments are invalid or exceed defined limits.
// Accepts 2 or 3 arguments: the string to pad, the target length, and an optional padding string.
func textPadLeft(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	if i2 > objects.MaxStringLen {
		return nil, objects.ErrStringLimit
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return objects.NewString(s1)
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = objects.ToStringArg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return objects.NewString(s1)
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := strings.Repeat(s3, int(padCount)) + s1
	return objects.NewString(retStr[int64(len(retStr))-i2:])
}

// textPadRight appends padding characters to the right of a string until it reaches the desired length.
// Accepts 2-3 arguments: string, target length (int), and optional padding string (default is a space).
// Returns an error if arguments are of invalid type, the number of arguments is wrong, or string size exceeds the limit.
func textPadRight(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return objects.NewString(s1)
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = objects.ToStringArg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return objects.NewString(s1)
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := s1 + strings.Repeat(s3, int(padCount))
	return objects.NewString(retStr[:i2])
}

// textRepeat repeats a string a specified number of times and returns the result as a new string object.
// It requires exactly 2 arguments: a string and an integer.
// Returns an error if arguments are of invalid types or if the resulting string exceeds the maximum allowed length.
func textRepeat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strings.Repeat(s1, int(i2)))
}

// textJoin concatenates elements of the first argument (array or immutable-array of strings) using the second argument as a separator.
// Returns an error if the number of arguments is incorrect, types are invalid, or the resulting string exceeds the size limit.
func textJoin(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	var sLen int
	var ss1 []string
	switch arg0 := args[0].(type) {
	case *objects.Array:
		for idx, a := range arg0.Values() {
			as, err := objects.ToStringArg(idx, a)
			if err != nil {
				return nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	case *objects.ArrayImmutable:
		for idx, a := range arg0.Values() {
			as, err := objects.ToStringArg(idx, a)
			if err != nil {
				return nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	default:
		return nil, objects.NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strings.Join(ss1, s2))
}

// textFormatBool converts a single boolean argument to its string representation ("true" or "false").
// Returns an error if the argument count is not 1 or the argument type is not boolean.
func textFormatBool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	b1, err := objects.ToBoolArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if b1 {
		return objects.NewString("true")
	}
	return objects.NewString("false")
}

// textFormatFloat formats a float value according to a specific format, precision, and bit size.
// It expects four arguments: a float, a string format character, precision (int), and bit size (int).
// Returns a formatted string object or an error if arguments are invalid.
func textFormatFloat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	f1, err := objects.ToFloat64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strconv.FormatFloat(f1, s2[0], int(i3), int(i4)))
}

// textFormatInt formats an integer to a string using the specified base.
// The first argument must be an integer type, and the second argument must be an int-compatible base.
// It returns the formatted string or an error if the arguments are invalid.
func textFormatInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strconv.FormatInt(i1, int(i2)))
}

// textParseBool parses a string argument into a boolean value and returns it as an IObject.
// Returns an error if the argument is not a string or if it fails to parse as a boolean.
// Requires exactly one argument; otherwise, it returns an ErrWrongNumArguments error.
func textParseBool(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		err = objects.NewInvalidArgumentError(0, "string", args[0].TypeName())
		return
	}
	parsed, err := strconv.ParseBool(s1.Value())
	if err != nil {
		ret = objects.NewObjectError(err)
		return
	}
	if parsed {
		ret = objects.TrueValue
	} else {
		ret = objects.FalseValue
	}
	return
}

// textParseFloat parses a string into a float64 with a specified precision, returning a Float object or an error.
func textParseFloat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseFloat(s1, int(i2))
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewFloat(parsed), nil
}

// textParseNumber parses a numeric value from a string argument and returns it as a Float object.
// Returns an error if the argument count is not 1 or if the input type is invalid.
// Non-numeric characters in the input string are ignored during parsing.
func textParseNumber(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var target []rune
	for _, p := range s1 {
		if (p >= '0' && p <= '9') || p == '.' || p == '+' || p == '-' {
			target = append(target, p)
		}
	}
	parsed, err := strconv.ParseFloat(string(target), 64)
	if err != nil {
		return nil, err
	}
	return objects.NewFloat(parsed), nil
}

// textParseInt parses a string into an integer using the specified base and bit size, returning the resulting integer.
// The function expects exactly three arguments: a string, an integer for base, and an integer for bit size.
// Returns an error for invalid argument types, an incorrect number of arguments, or a parsing failure.
func textParseInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseInt(s1, int(i2), int(i3))
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(parsed), nil
}

// doTextReplace replaces occurrences of `old` with `new` in the string `s` up to `n` times and returns the result.
// If `n` is negative, replaces all occurrences. Returns a boolean indicating whether the operation succeeded.
// If the resulting string exceeds `objects.MaxStringLen`, returns an empty string and `false`.
func doTextReplace(s string, old string, new string, n int) (string, bool) {
	if old == new || n == 0 {
		return s, true // avoid allocation
	}
	if m := strings.Count(s, old); m == 0 {
		return s, true
	} else if n < 0 || m < n {
		n = m
	}
	t := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}
		ssj := s[start:j]
		if w+len(ssj)+len(new) > objects.MaxStringLen {
			return "", false
		}

		w += copy(t[w:], ssj)
		w += copy(t[w:], new)
		start = j + len(old)
	}
	ss := s[start:]
	if w+len(ss) > objects.MaxStringLen {
		return "", false
	}
	w += copy(t[w:], ss)
	return string(t[0:w]), true
}

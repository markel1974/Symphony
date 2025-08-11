package stdlib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// textModule is a map of various text processing functions, including string manipulation, regex operations, and type conversions.
var textModule = map[string]objects.Object{
	"re_match":       objects.NewUserFunction("re_match", textREMatch),                           // re_match(pattern, text) => bool/error
	"re_find":        objects.NewUserFunction("re_find", textREFind),                             // re_find(pattern, text, count) => [[{text:,begin:,end:}]]/undefined
	"re_replace":     objects.NewUserFunction("re_replace", textREReplace),                       // re_replace(pattern, text, repl) => string/error
	"re_split":       objects.NewUserFunction("re_split", textRESplit),                           // re_split(pattern, text, count) => [string]/error
	"re_compile":     objects.NewUserFunction("re_compile", textRECompile),                       // re_compile(pattern) => Regexp/error
	"compare":        objects.NewUserFunction("compare", FuncASSRI(strings.Compare)),             // compare(a, b) => int
	"contains":       objects.NewUserFunction("contains", FuncASSRB(strings.Contains)),           // contains(s, substr) => bool
	"contains_any":   objects.NewUserFunction("contains_any", FuncASSRB(strings.ContainsAny)),    // contains_any(s, chars) => bool
	"count":          objects.NewUserFunction("count", FuncASSRI(strings.Count)),                 // count(s, substr) => int
	"equal_fold":     objects.NewUserFunction("equal_fold", FuncASSRB(strings.EqualFold)),        // "equal_fold(s, t) => bool
	"fields":         objects.NewUserFunction("fields", FuncASRSs(strings.Fields)),               // fields(s) => [string]
	"has_prefix":     objects.NewUserFunction("has_prefix", FuncASSRB(strings.HasPrefix)),        // has_prefix(s, prefix) => bool
	"has_suffix":     objects.NewUserFunction("has_suffix", FuncASSRB(strings.HasSuffix)),        // has_suffix(s, suffix) => bool
	"index":          objects.NewUserFunction("index", FuncASSRI(strings.Index)),                 // index(s, substr) => int
	"index_any":      objects.NewUserFunction("index_any", FuncASSRI(strings.IndexAny)),          // index_any(s, chars) => int
	"join":           objects.NewUserFunction("join", textJoin),                                  // join(arr, sep) => string
	"last_index":     objects.NewUserFunction("last_index", FuncASSRI(strings.LastIndex)),        // last_index(s, substr) => int
	"last_index_any": objects.NewUserFunction("last_index_any", FuncASSRI(strings.LastIndexAny)), // last_index_any(s, chars) => int
	"repeat":         objects.NewUserFunction("repeat", textRepeat),                              // repeat(s, count) => string
	"replace":        objects.NewUserFunction("replace", textReplace),                            // replace(s, old, new, n) => string
	"substr":         objects.NewUserFunction("substr", textSubstring),                           // substr(s, lower, upper) => string
	"split":          objects.NewUserFunction("split", FuncASSRSs(strings.Split)),                // split(s, sep) => [string]
	"split_after":    objects.NewUserFunction("split_after", FuncASSRSs(strings.SplitAfter)),     // split_after(s, sep) => [string]
	"split_after_n":  objects.NewUserFunction("split_after_n", FuncASSIRSs(strings.SplitAfterN)), // split_after_n(s, sep, n) => [string]
	"split_n":        objects.NewUserFunction("split_n", FuncASSIRSs(strings.SplitN)),            // split_n(s, sep, n) => [string]
	"title":          objects.NewUserFunction("title", FuncASRS(strings.Title)),                  // title(s) => string
	"to_lower":       objects.NewUserFunction("to_lower", FuncASRS(strings.ToLower)),             // to_lower(s) => string
	"to_title":       objects.NewUserFunction("to_title", FuncASRS(strings.ToTitle)),             // to_title(s) => string
	"to_upper":       objects.NewUserFunction("to_upper", FuncASRS(strings.ToUpper)),             // to_upper(s) => string
	"pad_left":       objects.NewUserFunction("pad_left", textPadLeft),                           // pad_left(s, pad_len, pad_with) => string
	"pad_right":      objects.NewUserFunction("pad_right", textPadRight),                         // pad_right(s, pad_len, pad_with) => string
	"trim":           objects.NewUserFunction("trim", FuncASSRS(strings.Trim)),                   // trim(s, cutset) => string
	"trim_left":      objects.NewUserFunction("trim_left", FuncASSRS(strings.TrimLeft)),          // trim_left(s, cutset) => string
	"trim_prefix":    objects.NewUserFunction("trim_prefix", FuncASSRS(strings.TrimPrefix)),      // trim_prefix(s, prefix) => string
	"trim_right":     objects.NewUserFunction("trim_right", FuncASSRS(strings.TrimRight)),        // trim_right(s, cutset) => string
	"trim_space":     objects.NewUserFunction("trim_space", FuncASRS(strings.TrimSpace)),         // trim_space(s) => string
	"trim_suffix":    objects.NewUserFunction("trim_suffix", FuncASSRS(strings.TrimSuffix)),      // trim_suffix(s, suffix) => string
	"atoi":           objects.NewUserFunction("atoi", FuncASRIE(strconv.Atoi)),                   // atoi(str) => int/error
	"format_bool":    objects.NewUserFunction("format_bool", textFormatBool),                     // format_bool(b) => string
	"format_float":   objects.NewUserFunction("format_float", textFormatFloat),                   // format_float(f, fmt, prec, bits) => string
	"format_int":     objects.NewUserFunction("format_int", textFormatInt),                       // format_int(i, base) => string
	"itoa":           objects.NewUserFunction("itoa", FuncAIRS(strconv.Itoa)),                    // itoa(i) => string
	"parse_bool":     objects.NewUserFunction("parse_bool", textParseBool),                       // parse_bool(str) => bool/error
	"parse_float":    objects.NewUserFunction("parse_float", textParseFloat),                     // parse_float(str, bits) => float/error
	"parse_number":   objects.NewUserFunction("parse_number", textParseNumber),                   // parse_float(str, bits) => float/error
	"parse_int":      objects.NewUserFunction("parse_int", textParseInt),                         // parse_int(str, base, bits) => int/error
	"quote":          objects.NewUserFunction("quote", FuncASRS(strconv.Quote)),                  // quote(str) => string
	"unquote":        objects.NewUserFunction("unquote", FuncASRSE(strconv.Unquote)),             // unquote(str) => string/error
}

// textREMatch checks if the first argument (regex pattern) matches the second argument (string) and returns a boolean object.
// Returns an error if arguments are not string or if there is an invalid number of arguments.
func textREMatch(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		return
	}
	matched, err := regexp.MatchString(s1, s2)
	if err != nil {
		ret = wrapError(err)
		return
	}
	if matched {
		ret = objects.TrueValue
	} else {
		ret = objects.FalseValue
	}
	return
}

// textREFind performs a regular expression match on a string and extracts matched substrings with their positions.
func textREFind(args ...objects.Object) (ret objects.Object, err error) {
	numArgs := len(args)
	if numArgs != 2 && numArgs != 3 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		ret = wrapError(err)
		return
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		return
	}
	if numArgs < 3 {
		m := re.FindStringSubmatchIndex(s2)
		if m == nil {
			ret = objects.UndefinedValue
			return
		}
		arr := &objects.Array{}
		for i := 0; i < len(m); i += 2 {
			arr.Value = append(arr.Value,
				&objects.ImmutableMap{Value: map[string]objects.Object{
					"text":  &objects.String{Value: s2[m[i]:m[i+1]]},
					"begin": &objects.Int{Value: int64(m[i])},
					"end":   &objects.Int{Value: int64(m[i+1])},
				}})
		}
		ret = &objects.Array{Value: []objects.Object{arr}}
		return
	}
	i3, ok := objects.ToInt(args[2])
	if !ok {
		err = errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
		return
	}
	m := re.FindAllStringSubmatchIndex(s2, i3)
	if m == nil {
		ret = objects.UndefinedValue
		return
	}
	arr := &objects.Array{}
	for _, m := range m {
		subMatch := &objects.Array{}
		for i := 0; i < len(m); i += 2 {
			subMatch.Value = append(subMatch.Value,
				&objects.ImmutableMap{Value: map[string]objects.Object{
					"text":  &objects.String{Value: s2[m[i]:m[i+1]]},
					"begin": &objects.Int{Value: int64(m[i])},
					"end":   &objects.Int{Value: int64(m[i+1])},
				}})
		}

		arr.Value = append(arr.Value, subMatch)
	}
	ret = arr
	return
}

// textREReplace performs a regular expression replacement in a string.
// Accepts three arguments: regex pattern as the first argument, source string as the second, and replacement string as the third.
// Returns the modified string or an error if the input arguments are invalid or the operation fails.
func textREReplace(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 3 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		return
	}
	s3, ok := objects.ToString(args[2])
	if !ok {
		err = errors.NewInvalidArgumentType("third", "string(compatible)", args[2].TypeName())
		return
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		ret = wrapError(err)
	} else {
		s, ok := doTextRegexpReplace(re, s2, s3)
		if !ok {
			return nil, errors.ErrStringLimit
		}

		ret = &objects.String{Value: s}
	}
	return
}

// textRESplit splits a string by a regular expression pattern and returns an array of resulting substrings.
// The first argument is the regex pattern, the second is the string to split, and the optional third is the max splits.
// Returns an error if arguments' types or counts are incorrect, or if the regex compilation fails.
func textRESplit(args ...objects.Object) (ret objects.Object, err error) {
	numArgs := len(args)
	if numArgs != 2 && numArgs != 3 {
		err = errors.ErrWrongNumArguments
		return
	}

	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}

	s2, ok := objects.ToString(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		return
	}

	var i3 = -1
	if numArgs > 2 {
		i3, ok = objects.ToInt(args[2])
		if !ok {
			err = errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
			return
		}
	}

	re, err := regexp.Compile(s1)
	if err != nil {
		ret = wrapError(err)
		return
	}

	arr := &objects.Array{}
	for _, s := range re.Split(s2, i3) {
		arr.Value = append(arr.Value, &objects.String{Value: s})
	}

	ret = arr

	return
}

// textRECompile compiles a string argument into a regular expression and returns a map of regex operations or an error.
func textRECompile(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		ret = wrapError(err)
	} else {
		ret = makeTextRegexp(re)
	}
	return
}

// textReplace performs a string replacement operation using 4 arguments: original string, target, replacement, and limit.
// It returns the modified string as an object or an error if the input arguments are invalid or exceed string limits.
func textReplace(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 4 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		return
	}
	s3, ok := objects.ToString(args[2])
	if !ok {
		err = errors.NewInvalidArgumentType("third", "string(compatible)", args[2].TypeName())
		return
	}
	i4, ok := objects.ToInt(args[3])
	if !ok {
		err = errors.NewInvalidArgumentType("fourth", "int(compatible)", args[3].TypeName())
		return
	}
	s, ok := doTextReplace(s1, s2, s3, i4)
	if !ok {
		err = errors.ErrStringLimit
		return
	}
	ret = &objects.String{Value: s}
	return
}

// textSubstring extracts a substring from a string, using start and optional end indices provided as arguments.
// Returns an error if arguments are invalid or indices are out of bounds.
func textSubstring(args ...objects.Object) (ret objects.Object, err error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		return
	}
	strlen := len(s1)
	i3 := strlen
	if argsLen == 3 {
		i3, ok = objects.ToInt(args[2])
		if !ok {
			err = errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
			return
		}
	}
	if i2 > i3 {
		err = errors.ErrInvalidIndexType
		return
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
	ret = &objects.String{Value: s1[i2:i3]}
	return
}

// textPadLeft pads a string from the left with a specified character or default space until the string reaches a given length.
// Returns an error if arguments are invalid or exceed defined limits.
// Accepts 2 or 3 arguments: the string to pad, the target length, and an optional padding string.
func textPadLeft(args ...objects.Object) (ret objects.Object, err error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		return
	}
	if i2 > objects.MaxStringLen {
		return nil, errors.ErrStringLimit
	}
	sLen := len(s1)
	if sLen >= i2 {
		ret = &objects.String{Value: s1}
		return
	}
	s3 := " "
	if argsLen == 3 {
		if s3, ok = objects.ToString(args[2]); !ok {
			err = errors.NewInvalidArgumentType("third", "string(compatible)", args[2].TypeName())
			return
		}
	}
	padStrLen := len(s3)
	if padStrLen == 0 {
		ret = &objects.String{Value: s1}
		return
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := strings.Repeat(s3, padCount) + s1
	ret = &objects.String{Value: retStr[len(retStr)-i2:]}
	return
}

// textPadRight appends padding characters to the right of a string until it reaches the desired length.
// Accepts 2-3 arguments: string, target length (int), and optional padding string (default is a space).
// Returns an error if arguments are of invalid type, the number of arguments is wrong, or string size exceeds the limit.
func textPadRight(args ...objects.Object) (ret objects.Object, err error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		return
	}
	if i2 > objects.MaxStringLen {
		return nil, errors.ErrStringLimit
	}
	sLen := len(s1)
	if sLen >= i2 {
		ret = &objects.String{Value: s1}
		return
	}
	s3 := " "
	if argsLen == 3 {
		s3, ok = objects.ToString(args[2])
		if !ok {
			err = errors.NewInvalidArgumentType("third", "string(compatible)", args[2].TypeName())
			return
		}
	}
	padStrLen := len(s3)
	if padStrLen == 0 {
		ret = &objects.String{Value: s1}
		return
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := s1 + strings.Repeat(s3, padCount)
	ret = &objects.String{Value: retStr[:i2]}
	return
}

// textRepeat repeats a string a specified number of times and returns the result as a new string object.
// It requires exactly 2 arguments: a string and an integer.
// Returns an error if arguments are of invalid types or if the resulting string exceeds the maximum allowed length.
func textRepeat(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
	}
	if len(s1)*i2 > objects.MaxStringLen {
		return nil, errors.ErrStringLimit
	}
	return &objects.String{Value: strings.Repeat(s1, i2)}, nil
}

// textJoin concatenates elements of the first argument (array or immutable-array of strings) using the second argument as a separator.
// Returns an error if the number of arguments is incorrect, types are invalid, or the resulting string exceeds the size limit.
func textJoin(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	var sLen int
	var ss1 []string
	switch arg0 := args[0].(type) {
	case *objects.Array:
		for idx, a := range arg0.Value {
			as, ok := objects.ToString(a)
			if !ok {
				return nil, errors.NewInvalidArgumentType(fmt.Sprintf("first[%d]", idx), "string(compatible)", a.TypeName())
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	case *objects.ImmutableArray:
		for idx, a := range arg0.Value {
			as, ok := objects.ToString(a)
			if !ok {
				return nil, errors.NewInvalidArgumentType(fmt.Sprintf("first[%d]", idx), "string(compatible)", a.TypeName())
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	default:
		return nil, errors.NewInvalidArgumentType("first", "array", args[0].TypeName())
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
	}
	if sLen+len(s2)*(len(ss1)-1) > objects.MaxStringLen {
		return nil, errors.ErrStringLimit
	}
	return &objects.String{Value: strings.Join(ss1, s2)}, nil
}

// textFormatBool converts a single boolean argument to its string representation ("true" or "false").
// Returns an error if the argument count is not 1 or the argument type is not boolean.
func textFormatBool(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = errors.ErrWrongNumArguments
		return
	}
	b1, ok := args[0].(*objects.Bool)
	if !ok {
		err = errors.NewInvalidArgumentType("first", "bool", args[0].TypeName())
		return
	}
	if b1 == objects.TrueValue {
		ret = &objects.String{Value: "true"}
	} else {
		ret = &objects.String{Value: "false"}
	}
	return
}

// textFormatFloat formats a float value according to a specific format, precision, and bit size.
// It expects four arguments: a float, a string format character, precision (int), and bit size (int).
// Returns a formatted string object or an error if arguments are invalid.
func textFormatFloat(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 4 {
		err = errors.ErrWrongNumArguments
		return
	}
	f1, ok := args[0].(*objects.Float)
	if !ok {
		err = errors.NewInvalidArgumentType("first", "float", args[0].TypeName())
		return
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
		return
	}
	i3, ok := objects.ToInt(args[2])
	if !ok {
		err = errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
		return
	}
	i4, ok := objects.ToInt(args[3])
	if !ok {
		err = errors.NewInvalidArgumentType("fourth", "int(compatible)", args[3].TypeName())
		return
	}
	ret = &objects.String{Value: strconv.FormatFloat(f1.Value, s2[0], i3, i4)}
	return
}

// textFormatInt formats an integer to a string using the specified base.
// The first argument must be an integer type, and the second argument must be an int-compatible base.
// It returns the formatted string or an error if the arguments are invalid.
func textFormatInt(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = errors.ErrWrongNumArguments
		return
	}
	i1, ok := args[0].(*objects.Int)
	if !ok {
		err = errors.NewInvalidArgumentType("first", "int", args[0].TypeName())
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		return
	}
	ret = &objects.String{Value: strconv.FormatInt(i1.Value, i2)}
	return
}

// textParseBool parses a string argument into a boolean value and returns it as an Object.
// Returns an error if the argument is not a string or if it fails to parse as a boolean.
// Requires exactly one argument; otherwise, it returns an ErrWrongNumArguments error.
func textParseBool(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string", args[0].TypeName())
		return
	}
	parsed, err := strconv.ParseBool(s1.Value)
	if err != nil {
		ret = wrapError(err)
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
func textParseFloat(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 2 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string", args[0].TypeName())
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		return
	}
	parsed, err := strconv.ParseFloat(s1.Value, i2)
	if err != nil {
		ret = wrapError(err)
		return
	}
	ret = &objects.Float{Value: parsed}
	return
}

// textParseNumber parses a numeric value from a string argument and returns it as a Float object.
// Returns an error if the argument count is not 1 or if the input type is invalid.
// Non-numeric characters in the input string are ignored during parsing.
func textParseNumber(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 1 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		ret = &objects.Float{Value: 0}
		return
	}

	var target []rune
	for _, p := range s1.Value {
		if (p >= '0' && p <= '9') || p == '.' || p == '+' || p == '-' {
			target = append(target, p)
		}
	}
	parsed, err := strconv.ParseFloat(string(target), 64)
	if err != nil {
		ret = &objects.Float{Value: 0}
		return
	}
	ret = &objects.Float{Value: parsed}
	return
}

// textParseInt parses a string into an integer using the specified base and bit size, returning the resulting integer.
// The function expects exactly three arguments: a string, an integer for base, and an integer for bit size.
// Returns an error for invalid argument types, an incorrect number of arguments, or a parsing failure.
func textParseInt(args ...objects.Object) (ret objects.Object, err error) {
	if len(args) != 3 {
		err = errors.ErrWrongNumArguments
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		err = errors.NewInvalidArgumentType("first", "string", args[0].TypeName())
		return
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		err = errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		return
	}
	i3, ok := objects.ToInt(args[2])
	if !ok {
		err = errors.NewInvalidArgumentType("third", "int(compatible)", args[2].TypeName())
		return
	}
	parsed, err := strconv.ParseInt(s1.Value, i2, i3)
	if err != nil {
		ret = wrapError(err)
		return
	}
	ret = &objects.Int{Value: parsed}
	return
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

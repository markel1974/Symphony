package stdlib

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// _stringsModule is a map containing string manipulation and formatting functions, regex utilities, and conversions.
var _stringsModule = map[string]objects.IObject{
	"re_match":     objects.NewFunctionModule(objects.FunctionModuleDef, "re_match", stringsREMatch),                              // re_match(pattern, text) => bool/error
	"re_find":      objects.NewFunctionModule(objects.FunctionModuleDef, "re_find", stringsREFind),                                // re_find(pattern, text, count) => [[{text:,begin:,end:}]]/undefined
	"re_replace":   objects.NewFunctionModule(objects.FunctionModuleDef, "re_replace", stringsREReplace),                          // re_replace(pattern, text, repl) => string/error
	"re_split":     objects.NewFunctionModule(objects.FunctionModuleDef, "re_split", stringsRESplit),                              // re_split(pattern, text, count) => [string]/error
	"re_compile":   objects.NewFunctionModule(objects.FunctionModuleDef, "re_compile", stringsRECompile),                          // re_compile(pattern) => Regexp/error
	"Compare":      objects.NewFunctionModule(objects.FunctionModuleDef, "Compare", objects.FuncIssOi(strings.Compare)),           // compare(a, b) => int
	"Contains":     objects.NewFunctionModule(objects.FunctionModuleDef, "Contains", objects.FuncIssOb(strings.Contains)),         // contains(s, substr) => bool
	"ContainsAny":  objects.NewFunctionModule(objects.FunctionModuleDef, "ContainsAny", objects.FuncIssOb(strings.ContainsAny)),   // contains_any(s, chars) => bool
	"Count":        objects.NewFunctionModule(objects.FunctionModuleDef, "Count", objects.FuncIssOi(strings.Count)),               // count(s, substr) => int
	"EqualFold":    objects.NewFunctionModule(objects.FunctionModuleDef, "EqualFold", objects.FuncIssOb(strings.EqualFold)),       // "equal_fold(s, t) => bool
	"Fields":       objects.NewFunctionModule(objects.FunctionModuleDef, "Fields", objects.FuncIsOsS(strings.Fields)),             // fields(s) => [string]
	"HasPrefix":    objects.NewFunctionModule(objects.FunctionModuleDef, "HasPrefix", objects.FuncIssOb(strings.HasPrefix)),       // has_prefix(s, prefix) => bool
	"HasSuffix":    objects.NewFunctionModule(objects.FunctionModuleDef, "HasSuffix", objects.FuncIssOb(strings.HasSuffix)),       // has_suffix(s, suffix) => bool
	"Index":        objects.NewFunctionModule(objects.FunctionModuleDef, "Index", objects.FuncIssOi(strings.Index)),               // index(s, substr) => int
	"IndexAny":     objects.NewFunctionModule(objects.FunctionModuleDef, "IndexAny", objects.FuncIssOi(strings.IndexAny)),         // index_any(s, chars) => int
	"Join":         objects.NewFunctionModule(objects.FunctionModuleDef, "Join", stringsJoin),                                     // join(arr, sep) => string
	"LastIndex":    objects.NewFunctionModule(objects.FunctionModuleDef, "LastIndex", objects.FuncIssOi(strings.LastIndex)),       // last_index(s, substr) => int
	"LastIndexAny": objects.NewFunctionModule(objects.FunctionModuleDef, "LastIndexAny", objects.FuncIssOi(strings.LastIndexAny)), // last_index_any(s, chars) => int
	"Repeat":       objects.NewFunctionModule(objects.FunctionModuleDef, "Repeat", stringsRepeat),                                 // repeat(s, count) => string
	"Replace":      objects.NewFunctionModule(objects.FunctionModuleDef, "Replace", stringsReplace),                               // replace(s, old, new, n) => string
	"Substring":    objects.NewFunctionModule(objects.FunctionModuleDef, "Substring", stringsSubstring),                           // substr(s, lower, upper) => string
	"Split":        objects.NewFunctionModule(objects.FunctionModuleDef, "Split", objects.FuncIssOsS(strings.Split)),              // split(s, sep) => [string]
	"SplitAfter":   objects.NewFunctionModule(objects.FunctionModuleDef, "SplitAfter", objects.FuncIssOsS(strings.SplitAfter)),    // split_after(s, sep) => [string]
	"SplitAfterN":  objects.NewFunctionModule(objects.FunctionModuleDef, "SplitAfterN", objects.FuncIssiOsS(strings.SplitAfterN)), // split_after_n(s, sep, n) => [string]
	"SplitN":       objects.NewFunctionModule(objects.FunctionModuleDef, "SplitN", objects.FuncIssiOsS(strings.SplitN)),           // split_n(s, sep, n) => [string]
	"Title":        objects.NewFunctionModule(objects.FunctionModuleDef, "Title", objects.FuncIsOs(strings.Title)),                // title(s) => string
	"ToLower":      objects.NewFunctionModule(objects.FunctionModuleDef, "ToLower", objects.FuncIsOs(strings.ToLower)),            // to_lower(s) => string
	"ToTitle":      objects.NewFunctionModule(objects.FunctionModuleDef, "ToTitle", objects.FuncIsOs(strings.ToTitle)),            // to_title(s) => string
	"ToUpper":      objects.NewFunctionModule(objects.FunctionModuleDef, "ToUpper", objects.FuncIsOs(strings.ToUpper)),            // to_upper(s) => string
	"PadLeft":      objects.NewFunctionModule(objects.FunctionModuleDef, "PadLeft", stringsPadLeft),                               // pad_left(s, pad_len, pad_with) => string
	"PadRight":     objects.NewFunctionModule(objects.FunctionModuleDef, "PadRight", stringsPadRight),                             // pad_right(s, pad_len, pad_with) => string
	"Trim":         objects.NewFunctionModule(objects.FunctionModuleDef, "Trim", objects.FuncIssOs(strings.Trim)),                 // trim(s, cutset) => string
	"TrimLeft":     objects.NewFunctionModule(objects.FunctionModuleDef, "TrimLeft", objects.FuncIssOs(strings.TrimLeft)),         // trim_left(s, cutset) => string
	"TrimPrefix":   objects.NewFunctionModule(objects.FunctionModuleDef, "TrimPrefix", objects.FuncIssOs(strings.TrimPrefix)),     // trim_prefix(s, prefix) => string
	"TrimRight":    objects.NewFunctionModule(objects.FunctionModuleDef, "TrimRight", objects.FuncIssOs(strings.TrimRight)),       // trim_right(s, cutset) => string
	"TrimSpace":    objects.NewFunctionModule(objects.FunctionModuleDef, "TrimSpace", objects.FuncIsOs(strings.TrimSpace)),        // trim_space(s) => string
	"TrimSuffix":   objects.NewFunctionModule(objects.FunctionModuleDef, "TrimSuffix", objects.FuncIssOs(strings.TrimSuffix)),     // trim_suffix(s, suffix) => string
	"Atoi":         objects.NewFunctionModule(objects.FunctionModuleDef, "Atoi", objects.FuncIsOie(strconv.Atoi)),                 // atoi(str) => int/error
	"FormatBool":   objects.NewFunctionModule(objects.FunctionModuleDef, "FormatBool", stringsFormatBool),                         // format_bool(b) => string
	"FormatFloat":  objects.NewFunctionModule(objects.FunctionModuleDef, "FormatFloat", stringsFormatFloat),                       // format_float(f, fmt, prec, bits) => string
	"FormatInt":    objects.NewFunctionModule(objects.FunctionModuleDef, "FormatInt", stringsFormatInt),                           // format_int(i, base) => string
	"Itoa":         objects.NewFunctionModule(objects.FunctionModuleDef, "Itoa", objects.FuncIiOs(strconv.Itoa)),                  // itoa(i) => string
	"ParseBool":    objects.NewFunctionModule(objects.FunctionModuleDef, "ParseBool", stringsParseBool),                           // parse_bool(str) => bool/error
	"ParseFloat":   objects.NewFunctionModule(objects.FunctionModuleDef, "ParseFloat", stringsParseFloat),                         // parse_float(str, bits) => float/error
	"ParseNumber":  objects.NewFunctionModule(objects.FunctionModuleDef, "ParseNumber", stringsParseNumber),                       // parse_float(str, bits) => float/error
	"ParseInt":     objects.NewFunctionModule(objects.FunctionModuleDef, "ParseInt", stringsParseInt),                             // parse_int(str, base, bits) => int/error
	"Quote":        objects.NewFunctionModule(objects.FunctionModuleDef, "Quote", objects.FuncIsOs(strconv.Quote)),                // quote(str) => string
	"Unquote":      objects.NewFunctionModule(objects.FunctionModuleDef, "Unquote", objects.FuncIsOse(strconv.Unquote)),           // unquote(str) => string/error
}

// stringsREMatch matches a string against a given regular expression and returns a boolean IObject result or an error.
func stringsREMatch(args ...objects.IObject) (objects.IObject, error) {
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

// stringsREFind performs a regex search on a string. Returns matches with details or undefined if no match is found.
func stringsREFind(args ...objects.IObject) (objects.IObject, error) {
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

// stringsREReplace performs a regex-based replacement on a string using provided pattern, target, and replacement values.
// It takes exactly three arguments: a regex pattern string, a target string, and a replacement string.
// Returns a new string after performing the replacements or an error if inputs are invalid.
func stringsREReplace(args ...objects.IObject) (objects.IObject, error) {
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

// stringsRESplit takes a regex pattern, a string, and an optional limit, returning an array of substrings split by the pattern.
func stringsRESplit(args ...objects.IObject) (objects.IObject, error) {
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

// stringsRECompile compiles a string into a regular expression and returns it as an immutable map of regex functions.
func stringsRECompile(args ...objects.IObject) (objects.IObject, error) {
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

// stringsReplace replaces occurrences of a substring within a string with another substring up to a specified limit.
// Takes four arguments: the original string, the substring to replace, the replacement string, and the maximum replacements.
// Returns the resulting string as an IObject and an error if the inputs are invalid or the output exceeds size constraints.
func stringsReplace(args ...objects.IObject) (objects.IObject, error) {
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

// stringsSubstring extracts a substring from a given string based on start and optional end indices.
// Returns an error if arguments are invalid, indices are out of range, or the indices are incompatible.
// The function expects 2 or 3 arguments: the string, start index, and optional end index.
// Adjusts indices to be within the valid range of the string.
func stringsSubstring(args ...objects.IObject) (objects.IObject, error) {
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

// stringsPadLeft pads the input string on the left to a specified length with a given substring or space if not provided.
// Returns an error if the number of arguments is incorrect, input types are invalid, or the padded string exceeds limits.
func stringsPadLeft(args ...objects.IObject) (objects.IObject, error) {
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

// stringsPadRight pads the input string on the right to a specified length using a given padding string or space by default.
// Returns the padded string or an error if invalid arguments are provided.
func stringsPadRight(args ...objects.IObject) (objects.IObject, error) {
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

// stringsRepeat repeats a string a specified number of times and returns the result as an IObject.
// Accepts exactly two arguments: a string and an integer.
// Returns an error if the argument count is invalid or if conversions fail.
func stringsRepeat(args ...objects.IObject) (objects.IObject, error) {
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

// stringsJoin concatenates array elements into a single string using a specified delimiter.
// The first argument must be an array or immutable array, and the second argument must be a string.
// Returns a new string object or an error if the inputs are invalid.
func stringsJoin(args ...objects.IObject) (ret objects.IObject, err error) {
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

// stringsFormatBool formats a boolean IObject as a string, returning "true" or "false". Accepts exactly one argument.
func stringsFormatBool(args ...objects.IObject) (objects.IObject, error) {
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

// stringsFormatFloat formats a float64 according to the provided format, precision, and bit size arguments.
// Returns a string object with the formatted result or an error if the arguments are invalid.
func stringsFormatFloat(args ...objects.IObject) (objects.IObject, error) {
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

// stringsFormatInt formats an int64 value to a string in a specified numeric base and returns it as an IObject.
// Accepts two arguments: the integer to format and the base, both convertible to int64.
// Returns an error if the conversion or base is invalid, or if the argument count is incorrect.
func stringsFormatInt(args ...objects.IObject) (objects.IObject, error) {
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

// stringsParseBool parses a string argument as a boolean and returns TrueValue or FalseValue based on the parsed result.
// It requires exactly one string argument, returns an error for invalid arguments or parsing errors.
func stringsParseBool(args ...objects.IObject) (ret objects.IObject, err error) {
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

// stringsParseFloat parses a string as a floating-point number with precision specified by the second argument.
// Returns the parsed float as an IObject or an error if parsing or argument validation fails.
func stringsParseFloat(args ...objects.IObject) (objects.IObject, error) {
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

// stringsParseNumber extracts and parses numeric values from a string, converting them into a Float object if valid.
// Returns an error if the input is not a valid string argument or the numeric parsing fails.
func stringsParseNumber(args ...objects.IObject) (objects.IObject, error) {
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

// stringsParseInt parses a string as an integer of a specified base and bit size, returning the parsed value or an error.
// Takes three arguments: the string to parse, the base (e.g., 10 for decimal), and the size in bits (e.g., 64 for int64).
// Returns an IObject representing the parsed integer value or an error if parsing fails or arguments are invalid.
// Produces ErrWrongNumArguments if the number of arguments is not exactly three.
func stringsParseInt(args ...objects.IObject) (objects.IObject, error) {
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

// doTextReplace replaces occurrences of 'old' with 'new' in the string 's', up to 'n' times. Returns updated string and a success flag.
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

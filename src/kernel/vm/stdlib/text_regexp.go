package stdlib

import (
	"regexp"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// makeTextRegexp returns an ImmutableMap of functions for regex operations: match, find, replace, and split.
func makeTextRegexp(re *regexp.Regexp) *objects.ImmutableMap {
	return &objects.ImmutableMap{
		Value: map[string]objects.Object{
			// match(text) => bool
			"match": objects.NewUserFunction("match", func(args ...objects.Object) (objects.Object, error) { return doTextRegexpMatch(re, args...) }),
			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"find": objects.NewUserFunction("find", func(args ...objects.Object) (objects.Object, error) { return doTextRegexpFind(re, args...) }),
			// replace(src, repl) => string
			"replace": objects.NewUserFunction("replace", func(args ...objects.Object) (objects.Object, error) { return doTextRegexpREReplace(re, args...) }),
			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"split": objects.NewUserFunction("split", func(args ...objects.Object) (objects.Object, error) { return doTextRegexpRESplit(re, args...) }),
		},
	}
}

// doTextRegexpReplace performs a regex-based replacement in the input string, returning the modified string and success status.
// re is the compiled regular expression used for matching patterns in src.
// src is the source string to search within and replace matched substrings.
// repl defines the replacement pattern, which may include backreferences from the regex match.
// Returns the modified string and a boolean indicating if the operation was successful (e.g., did not exceed length limits).
func doTextRegexpReplace(re *regexp.Regexp, src, repl string) (string, bool) {
	idx := 0
	out := ""
	for _, m := range re.FindAllStringSubmatchIndex(src, -1) {
		var exp []byte
		exp = re.ExpandString(exp, repl, src, m)
		if len(out)+m[0]-idx+len(exp) > objects.MaxStringLen {
			return "", false
		}
		out += src[idx:m[0]] + string(exp)
		idx = m[1]
	}
	if idx < len(src) {
		if len(out)+len(src)-idx > objects.MaxStringLen {
			return "", false
		}
		out += src[idx:]
	}
	return out, true
}

// doTextRegexpMatch checks if the provided string argument matches the given regular expression and returns a boolean result.
// Returns an error if the number of arguments is not 1 or if the argument type is not a string (compatible).
func doTextRegexpMatch(re *regexp.Regexp, args ...objects.Object) (objects.Object, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
	}
	if re.MatchString(s1) {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// doTextRegexpFind attempts to find matches in a given string based on a specified regular expression.
// The first argument must be a string, and the second argument (optional) specifies a match limit (integer).
// Returns an array of match details (e.g., matched text, start index, end index) or undefined if no match is found.
// Returns an error if invalid argument types or wrong number of arguments are provided.
func doTextRegexpFind(re *regexp.Regexp, args ...objects.Object) (objects.Object, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
	}
	if numArgs == 1 {
		m := re.FindStringSubmatchIndex(s1)
		if m == nil {
			return objects.UndefinedValue, nil
		}
		arr := objects.NewArray(nil)
		for i := 0; i < len(m); i += 2 {
			arr.Value = append(arr.Value,
				&objects.ImmutableMap{
					Value: map[string]objects.Object{
						"text":  &objects.String{Value: s1[m[i]:m[i+1]]},
						"begin": &objects.Int{Value: int64(m[i])},
						"end":   &objects.Int{Value: int64(m[i+1])},
					}})
		}
		return objects.NewArray([]objects.Object{arr}), nil
	}
	i2, ok := objects.ToInt(args[1])
	if !ok {
		return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
	}
	m := re.FindAllStringSubmatchIndex(s1, i2)
	if m == nil {
		return objects.UndefinedValue, nil
	}
	arr := &objects.Array{}
	for _, m := range m {
		subMatch := &objects.Array{}
		for i := 0; i < len(m); i += 2 {
			subMatch.Value = append(subMatch.Value,
				&objects.ImmutableMap{
					Value: map[string]objects.Object{
						"text":  &objects.String{Value: s1[m[i]:m[i+1]]},
						"begin": &objects.Int{Value: int64(m[i])},
						"end":   &objects.Int{Value: int64(m[i+1])},
					}})
		}
		arr.Value = append(arr.Value, subMatch)
	}
	return arr, nil
}

// doTextRegexpREReplace performs a regex-based replacement with a source and replacement string, returning a new string or an error.
func doTextRegexpREReplace(re *regexp.Regexp, args ...objects.Object) (objects.Object, error) {
	if len(args) != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
	}
	s2, ok := objects.ToString(args[1])
	if !ok {
		return nil, errors.NewInvalidArgumentType("second", "string(compatible)", args[1].TypeName())
	}
	s, ok := doTextRegexpReplace(re, s1, s2)
	if !ok {
		return nil, errors.ErrStringLimit
	}
	return objects.NewString(s), nil
}

// doTextRegexpRESplit splits a string using a regular expression and returns an array of substrings. Accepts one or two arguments.
// The first argument must be a string, the second (optional) must be an integer specifying the maximum number of splits.
func doTextRegexpRESplit(re *regexp.Regexp, args ...objects.Object) (objects.Object, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, ok := objects.ToString(args[0])
	if !ok {
		return nil, errors.NewInvalidArgumentType("first", "string(compatible)", args[0].TypeName())
	}
	var i2 = -1
	if numArgs > 1 {
		i2, ok = objects.ToInt(args[1])
		if !ok {
			return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
		}
	}
	arr := objects.NewArray(nil)
	for _, s := range re.Split(s1, i2) {
		arr.Value = append(arr.Value, objects.NewString(s))
	}
	return arr, nil
}

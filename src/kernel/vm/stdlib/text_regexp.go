package stdlib

import (
	"regexp"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// makeTextRegexp returns an ImmutableMap of functions for regex operations: match, find, replace, and split.
func makeTextRegexp(re *regexp.Regexp) *objects.ImmutableMap {
	return objects.NewImmutableMap(
		map[string]objects.IObject{
			// match(text) => bool
			"match": objects.NewUserFunction("match", func(args ...objects.IObject) (objects.IObject, error) { return doTextRegexpMatch(re, args...) }),
			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"find": objects.NewUserFunction("find", func(args ...objects.IObject) (objects.IObject, error) { return doTextRegexpFind(re, args...) }),
			// replace(src, repl) => string
			"replace": objects.NewUserFunction("replace", func(args ...objects.IObject) (objects.IObject, error) { return doTextRegexpREReplace(re, args...) }),
			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"split": objects.NewUserFunction("split", func(args ...objects.IObject) (objects.IObject, error) { return doTextRegexpRESplit(re, args...) }),
		},
	)
}

// doTextRegexpReplace performs a regex-based replacement in the input string, returning the modified string and success status.
// re is the compiled regular expression used for matching patterns in src.
// src is the source string to search within and replace matched substrings.
// repl defines the replacement pattern, which may include backreferences from the regex match.
// Returns the modified string and a boolean indicating if the operation was successful (e.g., did not exceed length limits).
func doTextRegexpReplace(re *regexp.Regexp, src, repl string) string {
	idx := 0
	out := ""
	for _, m := range re.FindAllStringSubmatchIndex(src, -1) {
		var exp []byte
		exp = re.ExpandString(exp, repl, src, m)
		out += src[idx:m[0]] + string(exp)
		idx = m[1]
	}
	if idx < len(src) {
		out += src[idx:]
	}
	return out
}

// doTextRegexpMatch checks if the provided string argument matches the given regular expression and returns a boolean result.
// Returns an error if the number of arguments is not 1 or if the argument type is not a string (compatible).
func doTextRegexpMatch(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg("first", args[0])
	if err != nil {
		return nil, err
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
func doTextRegexpFind(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg("first", args[0])
	if err != nil {
		return nil, err
	}
	if numArgs == 1 {
		mRe := re.FindStringSubmatchIndex(s1)
		if mRe == nil {
			return objects.UndefinedValue, nil
		}
		arr := objects.NewArray(nil)
		for i := 0; i < len(mRe); i += 2 {
			arr.Append(objects.NewImmutableMap(
				map[string]objects.IObject{
					"text":  objects.NewStringNoSize(s1[mRe[i]:mRe[i+1]]),
					"begin": objects.NewInt(int64(mRe[i])),
					"end":   objects.NewInt(int64(mRe[i+1])),
				}))
		}
		return objects.NewArray([]objects.IObject{arr}), nil
	}
	i2, err := objects.ToIntArg("second", args[1])
	if err != nil {
		return nil, err
	}
	mRe := re.FindAllStringSubmatchIndex(s1, i2)
	if mRe == nil {
		return objects.UndefinedValue, nil
	}
	arr := objects.NewArray(nil)
	for _, m := range mRe {
		subMatch := objects.NewArray(nil)
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(objects.NewImmutableMap(
				map[string]objects.IObject{
					"text":  objects.NewStringNoSize(s1[m[i]:m[i+1]]),
					"begin": objects.NewInt(int64(m[i])),
					"end":   objects.NewInt(int64(m[i+1])),
				}))
		}
		arr.Append(subMatch)
	}
	return arr, nil
}

// doTextRegexpREReplace performs a regex-based replacement with a source and replacement string, returning a new string or an error.
func doTextRegexpREReplace(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg("first", args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg("second", args[1])
	if err != nil {
		return nil, err
	}
	s := doTextRegexpReplace(re, s1, s2)
	return objects.NewString(s)
}

// doTextRegexpRESplit splits a string using a regular expression and returns an array of substrings. Accepts one or two arguments.
// The first argument must be a string, the second (optional) must be an integer specifying the maximum number of splits.
func doTextRegexpRESplit(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, errors.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg("first", args[0])
	if err != nil {
		return nil, err
	}
	var i2 = -1
	if numArgs > 1 {
		i2, err = objects.ToIntArg("second", args[1])
		if err != nil {
			return nil, err
		}
	}
	arr := objects.NewArray(nil)
	for _, s := range re.Split(s1, i2) {
		v, err := objects.NewString(s)
		if err != nil {
			return nil, err
		}
		arr.Append(v)
	}
	return arr, nil
}

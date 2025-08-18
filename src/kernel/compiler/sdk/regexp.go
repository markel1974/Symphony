package sdk

import (
	"regexp"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

const (
	RegexpTextDef  = "Text"
	RegexpBeginDef = "Begin"
	RegexpEndDef   = "End"
)

var _regexpModule = map[string]objects.IObject{
	"Match":   objects.NewFunctionModule(objects.FunctionModuleDef, "Match", regexpMatch),     // re_match(pattern, text) => bool/error
	"Find":    objects.NewFunctionModule(objects.FunctionModuleDef, "Find", regexpFind),       // re_find(pattern, text, count) => [[{text:,begin:,end:}]]/undefined
	"Replace": objects.NewFunctionModule(objects.FunctionModuleDef, "Replace", regexpReplace), // re_replace(pattern, text, repl) => string/error
	"Split":   objects.NewFunctionModule(objects.FunctionModuleDef, "Split", regexpSplit),     // re_split(pattern, text, count) => [string]/error
	"Compile": objects.NewFunctionModule(objects.FunctionModuleDef, "Compile", regexpCompile),
}

// regexpMatch matches a string against a given regular expression and returns a boolean IObject result or an error.
func regexpMatch(args ...objects.IObject) (objects.IObject, error) {
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

// regexpFind performs a regex search on a string. Returns matches with details or undefined if no match is found.
func regexpFind(args ...objects.IObject) (objects.IObject, error) {
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
				RegexpTextDef:  objects.NewStringNoSize(s2[m[i]:m[i+1]]),
				RegexpBeginDef: objects.NewInt(int64(m[i])),
				RegexpEndDef:   objects.NewInt(int64(m[i+1])),
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
				RegexpTextDef:  objects.NewStringNoSize(s2[m[i]:m[i+1]]),
				RegexpBeginDef: objects.NewInt(int64(m[i])),
				RegexpEndDef:   objects.NewInt(int64(m[i+1])),
			}))
		}
		arr.Append(subMatch)
	}
	return arr, nil
}

// regexpSplit takes a regex pattern, a string, and an optional limit, returning an array of substrings split by the pattern.
func regexpSplit(args ...objects.IObject) (objects.IObject, error) {
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

// regexpCompile compiles a string into a regular expression and returns it as an immutable map of regex functions.
func regexpCompile(args ...objects.IObject) (objects.IObject, error) {
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
	return _regexpCompileResult(re), nil
}

// regexpReplace performs a regex-based replacement on a string using provided pattern, target, and replacement values.
// It takes exactly three arguments: a regex pattern string, a target string, and a replacement string.
// Returns a new string after performing the replacements or an error if inputs are invalid.
func regexpReplace(args ...objects.IObject) (objects.IObject, error) {
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
	s := _regexpReplaceInternal(re, s2, s3)
	v, err := objects.NewString(s)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// _regexpCompileResult returns an MapImmutable of functions for regex operations: match, find, replace, and split.
func _regexpCompileResult(re *regexp.Regexp) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			// match(text) => bool
			"Match": objects.NewFunctionModule(objects.FunctionModuleDef, "Match", func(args ...objects.IObject) (objects.IObject, error) { return _regexpMatch(re, args...) }),
			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"Find": objects.NewFunctionModule(objects.FunctionModuleDef, "Find", func(args ...objects.IObject) (objects.IObject, error) { return _regexpFind(re, args...) }),
			// replace(src, repl) => string
			"Replace": objects.NewFunctionModule(objects.FunctionModuleDef, "Replace", func(args ...objects.IObject) (objects.IObject, error) { return _regexpReplace(re, args...) }),
			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"Split": objects.NewFunctionModule(objects.FunctionModuleDef, "Split", func(args ...objects.IObject) (objects.IObject, error) { return _regexpSplit(re, args...) }),
		},
	)
}

// _regexpReplaceInternal performs a regex-based replacement in the input string, returning the modified string and success status.
// re is the compiled regular expression used for matching patterns in src.
// src is the source string to search within and replace matched substrings.
// repl defines the replacement pattern, which may include backreferences from the regex match.
// Returns the modified string and a boolean indicating if the operation was successful (e.g., did not exceed length limits).
func _regexpReplaceInternal(re *regexp.Regexp, src, repl string) string {
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

// _regexpMatch checks if the provided string argument matches the given regular expression and returns a boolean result.
// Returns an error if the number of arguments is not 1 or if the argument type is not a string (compatible).
func _regexpMatch(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if re.MatchString(s1) {
		return objects.TrueValue, nil
	}
	return objects.FalseValue, nil
}

// _regexpFind attempts to find matches in a given string based on a specified regular expression.
// The first argument must be a string, and the second argument (optional) specifies a match limit (integer).
// Returns an array of match details (e.g., matched text, start index, end index) or undefined if no match is found.
// Returns an error if invalid argument types or wrong number of arguments are provided.
func _regexpFind(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
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
			arr.Append(objects.NewMapImmutable(
				map[string]objects.IObject{
					RegexpTextDef:  objects.NewStringNoSize(s1[mRe[i]:mRe[i+1]]),
					RegexpBeginDef: objects.NewInt(int64(mRe[i])),
					RegexpEndDef:   objects.NewInt(int64(mRe[i+1])),
				}))
		}
		return objects.NewArray([]objects.IObject{arr}), nil
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	mRe := re.FindAllStringSubmatchIndex(s1, int(i2))
	if mRe == nil {
		return objects.UndefinedValue, nil
	}
	arr := objects.NewArray(nil)
	for _, m := range mRe {
		subMatch := objects.NewArray(nil)
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(objects.NewMapImmutable(
				map[string]objects.IObject{
					RegexpTextDef:  objects.NewStringNoSize(s1[m[i]:m[i+1]]),
					RegexpBeginDef: objects.NewInt(int64(m[i])),
					RegexpEndDef:   objects.NewInt(int64(m[i+1])),
				}))
		}
		arr.Append(subMatch)
	}
	return arr, nil
}

// _regexpReplace performs a regex-based replacement with a source and replacement string, returning a new string or an error.
func _regexpReplace(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
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
	s := _regexpReplaceInternal(re, s1, s2)
	return objects.NewString(s)
}

// _regexpSplit splits a string using a regular expression and returns an array of substrings. Accepts one or two arguments.
// The first argument must be a string, the second (optional) must be an integer specifying the maximum number of splits.
func _regexpSplit(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var i2 = int64(-1)
	if numArgs > 1 {
		i2, err = objects.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
	}
	arr := objects.NewArray(nil)
	for _, s := range re.Split(s1, int(i2)) {
		v, err := objects.NewString(s)
		if err != nil {
			return nil, err
		}
		arr.Append(v)
	}
	return arr, nil
}

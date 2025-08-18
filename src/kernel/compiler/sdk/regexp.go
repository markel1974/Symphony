package sdk

import (
	"regexp"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// RegexpTextDef defines the default key for regular expression text.
// RegexpBeginDef defines the default key for the beginning of a regular expression.
// RegexpEndDef defines the default key for the end of a regular expression.
const (
	RegexpTextDef  = "Text"
	RegexpBeginDef = "Begin"
	RegexpEndDef   = "End"
)

// Regexp represents a structure providing regular expression functionality through associated operations and methods.
type Regexp struct {
	module map[string]objects.IObject
}

// NewRegexp creates and returns a new instance of the Regexp struct with initialized module functions.
func NewRegexp() *Regexp {
	r := &Regexp{}
	r.module = map[string]objects.IObject{
		"Match":   objects.NewFunctionModule(objects.FunctionModuleDef, "Match", r.Match),
		"Find":    objects.NewFunctionModule(objects.FunctionModuleDef, "Find", r.Find),
		"Replace": objects.NewFunctionModule(objects.FunctionModuleDef, "Replace", r.Replace),
		"Split":   objects.NewFunctionModule(objects.FunctionModuleDef, "Split", r.Split),
		"Compile": objects.NewFunctionModule(objects.FunctionModuleDef, "Compile", r.Compile),
	}
	return r
}

// Name returns the name of Regexp module.
func (r *Regexp) Name() string {
	return "regexp"
}

// Module returns the module map associated with the Regexp, containing functionalities as key-value pairs.
func (r *Regexp) Module() map[string]objects.IObject {
	return r.module
}

// Match checks whether the second string argument matches the pattern defined by the first string argument.
// Returns a boolean object indicating match success or failure, and an error if any issues occur.
func (r *Regexp) Match(args ...objects.IObject) (objects.IObject, error) {
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

// Find searches a string for matches against a regular expression and returns match details or an array of matches.
// The first argument is the pattern, the second argument is the string to search, and optionally, the third is the limit.
// Returns an error if the number of arguments is incorrect or conversion fails for any argument.
func (r *Regexp) Find(args ...objects.IObject) (objects.IObject, error) {
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

// Split divides a string into substrings based on the regular expression pattern and returns them as an array.
func (r *Regexp) Split(args ...objects.IObject) (objects.IObject, error) {
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

// Compile compiles a single string argument into a regular expression and returns an object encapsulating its options.
func (r *Regexp) Compile(args ...objects.IObject) (objects.IObject, error) {
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
	return r.CompileOptions(re), nil
}

// Replace performs a regex-based replacement on the input string with the specified replacement string and returns the result.
func (r *Regexp) Replace(args ...objects.IObject) (objects.IObject, error) {
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
	s := r.replaceInternal(re, s2, s3)
	v, err := objects.NewString(s)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// CompileOptions returns a MapImmutable containing methods for operations like match, find, replace, and split using a compiled regexp.
func (r *Regexp) CompileOptions(re *regexp.Regexp) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			// match(text) => bool
			"Match": objects.NewFunctionModule(objects.FunctionModuleDef, "Match", func(args ...objects.IObject) (objects.IObject, error) { return r.CompileOptionMatch(re, args...) }),
			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"Find": objects.NewFunctionModule(objects.FunctionModuleDef, "Find", func(args ...objects.IObject) (objects.IObject, error) { return r.CompileOptionFind(re, args...) }),
			// replace(src, repl) => string
			"Replace": objects.NewFunctionModule(objects.FunctionModuleDef, "Replace", func(args ...objects.IObject) (objects.IObject, error) { return r.CompileOptionReplace(re, args...) }),
			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"Split": objects.NewFunctionModule(objects.FunctionModuleDef, "Split", func(args ...objects.IObject) (objects.IObject, error) { return r.CompileOptionSplit(re, args...) }),
		},
	)
}

// CompileOptionMatch checks if the given regular expression matches the provided string argument and returns a boolean result.
func (r *Regexp) CompileOptionMatch(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
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

// CompileOptionFind performs regex matching with options for finding single or multiple matches in a given string.
// Accepts a compiled regex and one or two additional arguments: a string to match and optionally the match limit.
// Returns a collection of matched substrings with their positions or an undefined value if no match is found.
// Returns an error if arguments are invalid or if type conversion fails.
func (r *Regexp) CompileOptionFind(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
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

// CompileOptionReplace replaces occurrences in the input string matching the regular expression with the replacement string.
func (r *Regexp) CompileOptionReplace(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
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
	s := r.replaceInternal(re, s1, s2)
	return objects.NewString(s)
}

// CompileOptionSplit splits the input string using the given compiled regular expression and optional max split count.
func (r *Regexp) CompileOptionSplit(re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
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

// replaceInternal performs a regex-based replacement on the input string using the specified pattern and replacement string.
func (r *Regexp) replaceInternal(re *regexp.Regexp, src, repl string) string {
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

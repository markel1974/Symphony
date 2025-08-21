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
	factory *objects.GateKeeper
	*Package
}

// NewRegexp creates and returns a new instance of the Regexp struct with initialized module functions.
func NewRegexp(factory *objects.GateKeeper) *Regexp {
	r := &Regexp{factory: factory}
	container := []*objects.FuncPackage{
		factory.NewFuncPackage(objects.FuncPackageDef, "Match", r.match),
		factory.NewFuncPackage(objects.FuncPackageDef, "Find", r.find),
		factory.NewFuncPackage(objects.FuncPackageDef, "Replace", r.replace),
		factory.NewFuncPackage(objects.FuncPackageDef, "Split", r.split),
		factory.NewFuncPackage(objects.FuncPackageDef, "Compile", r.compile),
	}
	r.Package = NewPackage("regexp", container, nil)
	return r
}

// Match checks whether the second string argument matches the pattern defined by the first string argument.
// Returns a boolean object indicating match success or failure, and an error if any issues occur.
func (r *Regexp) match(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := r.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	matched, err := regexp.MatchString(s1, s2)
	if err != nil {
		return r.factory.NewError(frame, err.Error()), nil
	}
	if matched {
		return r.factory.TrueValue(), nil
	}
	return r.factory.FalseValue(), nil
}

// Find searches a string for matches against a regular expression and returns match details or an array of matches.
// The first argument is the pattern, the second argument is the string to search, and optionally, the third is the limit.
// Returns an error if the number of arguments is incorrect or conversion fails for any argument.
func (r *Regexp) find(frame int, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 2 && numArgs != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := r.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return r.factory.NewError(frame, err.Error()), nil
	}
	if numArgs < 3 {
		m := re.FindStringSubmatchIndex(s2)
		if m == nil {
			return r.factory.UndefinedValue(), nil
		}
		arr := r.factory.NewArray(frame, nil)
		for i := 0; i < len(m); i += 2 {
			arr.Append(r.factory.NewMapImmutable(frame, map[string]objects.IObject{
				RegexpTextDef:  r.factory.NewString(frame, s2[m[i]:m[i+1]]),
				RegexpBeginDef: r.factory.NewInt(frame, int64(m[i])),
				RegexpEndDef:   r.factory.NewInt(frame, int64(m[i+1])),
			}))
		}
		return r.factory.NewArray(frame, []objects.IObject{arr}), nil
	}
	i3, err := r.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	mFA := re.FindAllStringSubmatchIndex(s2, int(i3))
	if mFA == nil {
		return r.factory.UndefinedValue(), nil
	}
	arr := r.factory.NewArray(frame, nil)
	for _, m := range mFA {
		subMatch := r.factory.NewArray(frame, nil)
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(r.factory.NewMapImmutable(frame, map[string]objects.IObject{
				RegexpTextDef:  r.factory.NewString(frame, s2[m[i]:m[i+1]]),
				RegexpBeginDef: r.factory.NewInt(frame, int64(m[i])),
				RegexpEndDef:   r.factory.NewInt(frame, int64(m[i+1])),
			}))
		}
		arr.Append(subMatch)
	}
	return arr, nil
}

// Split divides a string into substrings based on the regular expression pattern and returns them as an array.
func (r *Regexp) split(frame int, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 2 && numArgs != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := r.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	var i3 = int64(-1)
	if numArgs > 2 {
		i3, err = r.factory.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return r.factory.NewError(frame, err.Error()), nil
	}
	arr := r.factory.NewArray(frame, nil)
	for _, s := range re.Split(s2, int(i3)) {
		v := r.factory.NewString(frame, s)
		arr.Append(v)
	}
	return arr, nil
}

// Compile compiles a single string argument into a regular expression and returns an object encapsulating its options.
func (r *Regexp) compile(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, ok := r.factory.ToString(args[0])
	if !ok {
		return nil, objects.NewInvalidArgumentError(0, "string(compatible)", args[0].TypeName())
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return r.factory.NewError(frame, err.Error()), nil
	}
	return r.compileOptions(frame, re), nil
}

// Replace performs a regex-based replacement on the input string with the specified replacement string and returns the result.
func (r *Regexp) replace(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := r.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s3, err := r.factory.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return r.factory.NewError(frame, err.Error()), nil
	}
	s := r.replaceInternal(re, s2, s3)
	v := r.factory.NewString(frame, s)
	return v, nil
}

// CompileOptions returns a MapImmutable containing methods for operations like match, find, replace, and split using a compiled regexp.
func (r *Regexp) compileOptions(frame int, re *regexp.Regexp) *objects.MapImmutable {
	return r.factory.NewMapImmutable(frame,
		map[string]objects.IObject{
			// match(text) => bool
			"Match": r.factory.NewFuncFramePackage(frame, objects.FuncPackageDef, "Match", func(frame int, args ...objects.IObject) (objects.IObject, error) {
				return r.compileOptionMatch(frame, re, args...)
			}),
			"Find": r.factory.NewFuncFramePackage(frame, objects.FuncPackageDef, "Find", func(frame int, args ...objects.IObject) (objects.IObject, error) {
				return r.compileOptionFind(frame, re, args...)
			}),
			"Replace": r.factory.NewFuncFramePackage(frame, objects.FuncPackageDef, "Replace", func(frame int, args ...objects.IObject) (objects.IObject, error) {
				return r.compileOptionReplace(frame, re, args...)
			}),
			"Split": r.factory.NewFuncFramePackage(frame, objects.FuncPackageDef, "Split", func(frame int, args ...objects.IObject) (objects.IObject, error) {
				return r.compileOptionSplit(frame, re, args...)
			}),
		},
	)
}

// CompileOptionMatch checks if the given regular expression matches the provided string argument and returns a boolean result.
func (r *Regexp) compileOptionMatch(_ int, re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if re.MatchString(s1) {
		return r.factory.TrueValue(), nil
	}
	return r.factory.FalseValue(), nil
}

// CompileOptionFind performs regex matching with options for finding single or multiple matches in a given string.
// Accepts a compiled regex and one or two additional arguments: a string to match and optionally the match limit.
// Returns a collection of matched substrings with their positions or an undefined value if no match is found.
// Returns an error if arguments are invalid or if type conversion fails.
func (r *Regexp) compileOptionFind(frame int, re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if numArgs == 1 {
		mRe := re.FindStringSubmatchIndex(s1)
		if mRe == nil {
			return r.factory.UndefinedValue(), nil
		}
		arr := r.factory.NewArray(frame, nil)
		for i := 0; i < len(mRe); i += 2 {
			arr.Append(r.factory.NewMapImmutable(frame,
				map[string]objects.IObject{
					RegexpTextDef:  r.factory.NewString(frame, s1[mRe[i]:mRe[i+1]]),
					RegexpBeginDef: r.factory.NewInt(frame, int64(mRe[i])),
					RegexpEndDef:   r.factory.NewInt(frame, int64(mRe[i+1])),
				}))
		}
		return r.factory.NewArray(frame, []objects.IObject{arr}), nil
	}
	i2, err := r.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	mRe := re.FindAllStringSubmatchIndex(s1, int(i2))
	if mRe == nil {
		return r.factory.UndefinedValue(), nil
	}
	arr := r.factory.NewArray(frame, nil)
	for _, m := range mRe {
		subMatch := r.factory.NewArray(frame, nil)
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(r.factory.NewMapImmutable(frame,
				map[string]objects.IObject{
					RegexpTextDef:  r.factory.NewString(frame, s1[m[i]:m[i+1]]),
					RegexpBeginDef: r.factory.NewInt(frame, int64(m[i])),
					RegexpEndDef:   r.factory.NewInt(frame, int64(m[i+1])),
				}))
		}
		arr.Append(subMatch)
	}
	return arr, nil
}

// CompileOptionReplace replaces occurrences in the input string matching the regular expression with the replacement string.
func (r *Regexp) compileOptionReplace(frame int, re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := r.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s := r.replaceInternal(re, s1, s2)
	return r.factory.NewString(frame, s), nil
}

// CompileOptionSplit splits the input string using the given compiled regular expression and optional max split count.
func (r *Regexp) compileOptionSplit(frame int, re *regexp.Regexp, args ...objects.IObject) (objects.IObject, error) {
	numArgs := len(args)
	if numArgs != 1 && numArgs != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := r.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var i2 = int64(-1)
	if numArgs > 1 {
		i2, err = r.factory.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
	}
	arr := r.factory.NewArray(frame, nil)
	for _, s := range re.Split(s1, int(i2)) {
		v := r.factory.NewString(frame, s)
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

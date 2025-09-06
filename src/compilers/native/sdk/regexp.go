package sdk

import (
	"fmt"
	"regexp"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// RegexpTextDef defines the default key for regular expression text.
// RegexpBeginDef defines the default key for the beginning of a regular expression.
// RegexpEndDef defines the default key for the end of a regular expression.
const (
	RegexpTextDef  = "Text"
	RegexpBeginDef = "Begin"
	RegexpEndDef   = "End"
)

func init() {
	RegisterPackage(NewRegexp)
}

// Regexp represents a structure providing regular expression functionality through associated operations and methods.
type Regexp struct {
	*Package
}

// NewRegexp creates and returns a new instance of the Regexp struct with initialized module functions.
func NewRegexp(factory objects.IGateKeeper) IPackage {
	r := &Regexp{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Match", 2, r.match),
		factory.NewFuncImport(objects.FrameStatic, "Find", -1, r.find),
		factory.NewFuncImport(objects.FrameStatic, "Replace", 3, r.replace),
		factory.NewFuncImport(objects.FrameStatic, "Split", -1, r.split),
		factory.NewFuncImport(objects.FrameStatic, "Compile", 1, r.compile),
	}
	r.Package = NewExternalPackage("regexp", container, nil)
	return r
}

// Match checks whether the second string argument matches the pattern defined by the first string argument.
// Returns a boolean object indicating match success or failure, and an error if any issues occur.
func (r *Regexp) match(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	matched, err := regexp.MatchString(s1, s2)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	if matched {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// Find searches a string for matches against a regular expression and returns match details or an array of matches.
// The first argument is the pattern, the second argument is the string to search, and optionally, the third is the limit.
// Returns an error if the number of arguments is incorrect or conversion fails for any argument.
func (r *Regexp) find(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	numArgs := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	if numArgs < 3 {
		m := re.FindStringSubmatchIndex(s2)
		if m == nil {
			return 1, gk.UndefinedValue(), nil
		}
		obj := gk.NewArray(frame, nil)
		arr, ok := obj.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for i := 0; i < len(m); i += 2 {
			arr.Append(gk.NewMap(frame, map[string]objects.IObject{
				RegexpTextDef:  gk.NewString(frame, s2[m[i]:m[i+1]]),
				RegexpBeginDef: gk.NewInt(frame, int64(m[i])),
				RegexpEndDef:   gk.NewInt(frame, int64(m[i+1])),
			}))
		}
		return 1, gk.NewArray(frame, []objects.IObject{arr}), nil
	}
	i3, err := gk.ToInt64Arg(2, args)
	if err != nil {
		return 0, nil, err
	}
	mFA := re.FindAllStringSubmatchIndex(s2, int(i3))
	if mFA == nil {
		return 1, gk.UndefinedValue(), nil
	}
	base := gk.NewArray(frame, nil)
	arr, ok := base.(*objects.Array)
	if !ok {
		return 0, nil, fmt.Errorf("expected Array, got %T", base)
	}
	for _, m := range mFA {
		obj := gk.NewArray(frame, nil)
		subMatch, ok := obj.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(gk.NewMap(frame, map[string]objects.IObject{
				RegexpTextDef:  gk.NewString(frame, s2[m[i]:m[i+1]]),
				RegexpBeginDef: gk.NewInt(frame, int64(m[i])),
				RegexpEndDef:   gk.NewInt(frame, int64(m[i+1])),
			}))
		}
		arr.Append(subMatch)
	}
	return 1, arr, nil
}

// Split divides a string into substrings based on the regular expression pattern and returns them as an array.
func (r *Regexp) split(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	numArgs := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	var i3 = int64(-1)
	if numArgs > 2 {
		i3, err = gk.ToInt64Arg(2, args)
		if err != nil {
			return 0, nil, err
		}
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	base := gk.NewArray(frame, nil)
	arr, ok := base.(*objects.Array)
	if !ok {
		return 0, nil, fmt.Errorf("expected Array, got %T", base)
	}
	for _, s := range re.Split(s2, int(i3)) {
		v := gk.NewString(frame, s)
		arr.Append(v)
	}
	return 1, arr, nil
}

// Compile compiles a single string argument into a regular expression and returns an object encapsulating its options.
func (r *Regexp) compile(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}

	obj := gk.NewMap(frame,
		map[string]objects.IObject{
			// match(text) => bool
			"Match": gk.NewFuncImport(frame, "Match", 1, func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
				return r.compileOptionMatch(gk, frame, re, args...)
			}),
			"Find": gk.NewFuncImport(frame, "Find", -1, func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
				return r.compileOptionFind(gk, frame, re, args...)
			}),
			"Replace": gk.NewFuncImport(frame, "Replace", 2, func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
				return r.compileOptionReplace(gk, frame, re, args...)
			}),
			"Split": gk.NewFuncImport(frame, "Split", -1, func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
				return r.compileOptionSplit(gk, frame, re, args...)
			}),
		},
	)
	mImm, ok := obj.(*objects.Map)
	if !ok {
		return 0, nil, fmt.Errorf("unexpected object type: %T", obj)
	}
	return 1, mImm, nil
}

// Replace performs a regex-based replacement on the input string with the specified replacement string and returns the result.
func (r *Regexp) replace(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	s3, err := gk.ToStringArg(2, args)
	if err != nil {
		return 0, nil, err
	}
	re, err := regexp.Compile(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	s := r.replaceInternal(re, s2, s3)
	v := gk.NewString(frame, s)
	return 1, v, nil
}

// CompileOptionMatch checks if the given regular expression matches the provided string argument and returns a boolean result.
func (r *Regexp) compileOptionMatch(gk objects.IGateKeeper, _ int, re *regexp.Regexp, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	if re.MatchString(s1) {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// CompileOptionFind performs regex matching with options for finding single or multiple matches in a given string.
// Accepts a compiled regex and one or two additional arguments: a string to match and optionally the match limit.
// Returns a collection of matched substrings with their positions or an undefined value if no match is found.
// Returns an error if arguments are invalid or if type conversion fails.
func (r *Regexp) compileOptionFind(gk objects.IGateKeeper, frame int, re *regexp.Regexp, args ...objects.IObject) (uint, objects.IObject, error) {
	numArgs := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	if numArgs == 1 {
		mRe := re.FindStringSubmatchIndex(s1)
		if mRe == nil {
			return 1, gk.UndefinedValue(), nil
		}
		base := gk.NewArray(frame, nil)
		arr, ok := base.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", base)
		}
		for i := 0; i < len(mRe); i += 2 {
			arr.Append(gk.NewMap(frame,
				map[string]objects.IObject{
					RegexpTextDef:  gk.NewString(frame, s1[mRe[i]:mRe[i+1]]),
					RegexpBeginDef: gk.NewInt(frame, int64(mRe[i])),
					RegexpEndDef:   gk.NewInt(frame, int64(mRe[i+1])),
				}))
		}
		return 1, gk.NewArray(frame, []objects.IObject{arr}), nil
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	mRe := re.FindAllStringSubmatchIndex(s1, int(i2))
	if mRe == nil {
		return 1, gk.UndefinedValue(), nil
	}
	base := gk.NewArray(frame, nil)
	arr, ok := base.(*objects.Array)
	if !ok {
		return 0, nil, fmt.Errorf("expected Array, got %T", base)
	}
	for _, m := range mRe {
		obj := gk.NewArray(frame, nil)
		subMatch, ok := base.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for i := 0; i < len(m); i += 2 {
			subMatch.Append(gk.NewMap(frame,
				map[string]objects.IObject{
					RegexpTextDef:  gk.NewString(frame, s1[m[i]:m[i+1]]),
					RegexpBeginDef: gk.NewInt(frame, int64(m[i])),
					RegexpEndDef:   gk.NewInt(frame, int64(m[i+1])),
				}))
		}
		arr.Append(subMatch)
	}
	return 1, arr, nil
}

// CompileOptionReplace replaces occurrences in the input string matching the regular expression with the replacement string.
func (r *Regexp) compileOptionReplace(gk objects.IGateKeeper, frame int, re *regexp.Regexp, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	s := r.replaceInternal(re, s1, s2)
	return 1, gk.NewString(frame, s), nil
}

// CompileOptionSplit splits the input string using the given compiled regular expression and optional max split count.
func (r *Regexp) compileOptionSplit(gk objects.IGateKeeper, frame int, re *regexp.Regexp, args ...objects.IObject) (uint, objects.IObject, error) {
	numArgs := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	var i2 = int64(-1)
	if numArgs > 1 {
		i2, err = gk.ToInt64Arg(1, args)
		if err != nil {
			return 0, nil, err
		}
	}
	base := gk.NewArray(frame, nil)
	arr, ok := base.(*objects.Array)
	if !ok {
		return 0, nil, fmt.Errorf("expected Array, got %T", base)
	}
	for _, s := range re.Split(s1, int(i2)) {
		v := gk.NewString(frame, s)
		arr.Append(v)
	}
	return 1, arr, nil
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

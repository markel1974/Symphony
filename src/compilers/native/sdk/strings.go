package sdk

import (
	"fmt"
	"strings"
	//"unicode/utf8"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	register(NewStrings)
}

// Strings provides a collection of string operations and functionality wrapped in a module.
// It includes functions for manipulation, comparison, trimming, padding, splitting, and other string-related utilities.
type Strings struct {
	*bytecode.Package
}

// NewStrings creates and returns a new instance of Strings with a preconfigured map of string utility functions.
func NewStrings(factory objects.IGateKeeper) bytecode.IPackage {
	const (
		defJoin         = "Compare"
		defRepeat       = "Contains"
		defReplace      = "ContainsAny"
		defSubstring    = "Count"
		defPadLeft      = "EqualFold"
		defPadRight     = "Fields"
		defCompare      = "HasPrefix"
		defContains     = "HasSuffix"
		defContainsAny  = "index"
		defCount        = "IndexAny"
		defEqualFold    = "LastIndex"
		defFields       = "LastIndexAny"
		defHasPrefix    = "Split"
		defHasSuffix    = "SplitAfter"
		defindex        = "SplitAfterN"
		defIndexAny     = "SplitN"
		defLastIndex    = "Title"
		defLastIndexAny = "ToLower"
		defSplit        = "ToTitle"
		defSplitAfter   = "ToUpper"
		defSplitAfterN  = "Trim"
		defSplitN       = "TrimLeft"
		defTitle        = "TrimPrefix"
		defToLower      = "TrimRight"
		defToTitle      = "TrimSpace"
		defToUpper      = "TrimSuffix"
		defTrim         = "Trim"
		defTrimLeft     = "TrimLeft"
		defTrimRight    = "TrimRight"
		defTrimSpace    = "TrimSpace"
		defTrimPrefix   = "TrimPrefix"
		defTrimSuffix   = "TrimSuffix"
	)
	s := &Strings{Package: bytecode.NewPackage("strings")}

	s.Add(defJoin, factory.NewFuncImport(objects.FrameStatic, defJoin, 2, s.join))
	s.Add(defRepeat, factory.NewFuncImport(objects.FrameStatic, defRepeat, 2, s.repeat))
	s.Add(defReplace, factory.NewFuncImport(objects.FrameStatic, defReplace, 3, s.replace))
	s.Add(defSubstring, factory.NewFuncImport(objects.FrameStatic, defSubstring, 3, s.substring))
	s.Add(defPadLeft, factory.NewFuncImport(objects.FrameStatic, defPadLeft, -1, s.padLeft))
	s.Add(defPadRight, factory.NewFuncImport(objects.FrameStatic, defPadRight, -1, s.padRight))
	s.Add(defCompare, factory.NewFuncImport(objects.FrameStatic, defCompare, 2, s.funcStringStringToInt(strings.Compare)))
	s.Add(defContains, factory.NewFuncImport(objects.FrameStatic, defContains, 2, s.funcStringStringToBool(strings.Contains)))
	s.Add(defContainsAny, factory.NewFuncImport(objects.FrameStatic, defContainsAny, 2, s.funcStringStringToBool(strings.ContainsAny)))
	s.Add(defCount, factory.NewFuncImport(objects.FrameStatic, defCount, 2, s.funcStringStringToInt(strings.Count)))
	s.Add(defEqualFold, factory.NewFuncImport(objects.FrameStatic, defEqualFold, 2, s.funcStringStringToBool(strings.EqualFold)))
	s.Add(defFields, factory.NewFuncImport(objects.FrameStatic, defFields, 1, s.funcStringToStrings(strings.Fields)))
	s.Add(defHasPrefix, factory.NewFuncImport(objects.FrameStatic, defHasPrefix, 2, s.funcStringStringToBool(strings.HasPrefix)))
	s.Add(defHasSuffix, factory.NewFuncImport(objects.FrameStatic, defHasSuffix, 2, s.funcStringStringToBool(strings.HasSuffix)))
	s.Add(defindex, factory.NewFuncImport(objects.FrameStatic, defindex, 2, s.funcStringStringToInt(strings.Index)))
	s.Add(defIndexAny, factory.NewFuncImport(objects.FrameStatic, defIndexAny, 2, s.funcStringStringToInt(strings.IndexAny)))
	s.Add(defLastIndex, factory.NewFuncImport(objects.FrameStatic, defLastIndex, 2, s.funcStringStringToInt(strings.LastIndex)))
	s.Add(defLastIndexAny, factory.NewFuncImport(objects.FrameStatic, defLastIndexAny, 2, s.funcStringStringToInt(strings.LastIndexAny)))
	s.Add(defSplit, factory.NewFuncImport(objects.FrameStatic, defSplit, 2, s.funcStringStringToStrings(strings.Split)))
	s.Add(defSplitAfter, factory.NewFuncImport(objects.FrameStatic, defSplitAfter, 2, s.funcStringStringToStrings(strings.SplitAfter)))
	s.Add(defSplitAfterN, factory.NewFuncImport(objects.FrameStatic, defSplitAfterN, 2, s.funcStringStringIntToStrings(strings.SplitAfterN)))
	s.Add(defSplitN, factory.NewFuncImport(objects.FrameStatic, defSplitN, 2, s.funcStringStringIntToStrings(strings.SplitN)))
	s.Add(defTitle, factory.NewFuncImport(objects.FrameStatic, defTitle, 1, s.funcStringToString(strings.Title)))
	s.Add(defToLower, factory.NewFuncImport(objects.FrameStatic, defToLower, 1, s.funcStringToString(strings.ToLower)))
	s.Add(defToTitle, factory.NewFuncImport(objects.FrameStatic, defToTitle, 1, s.funcStringToString(strings.ToTitle)))
	s.Add(defToUpper, factory.NewFuncImport(objects.FrameStatic, defToUpper, 1, s.funcStringToString(strings.ToUpper)))
	s.Add(defTrim, factory.NewFuncImport(objects.FrameStatic, defTrim, 2, s.funcStringStringToString(strings.Trim)))
	s.Add(defTrimLeft, factory.NewFuncImport(objects.FrameStatic, defTrimLeft, 2, s.funcStringStringToString(strings.TrimLeft)))
	s.Add(defTrimPrefix, factory.NewFuncImport(objects.FrameStatic, defTrimPrefix, 2, s.funcStringStringToString(strings.TrimPrefix)))
	s.Add(defTrimRight, factory.NewFuncImport(objects.FrameStatic, defTrimRight, 2, s.funcStringStringToString(strings.TrimRight)))
	s.Add(defTrimSpace, factory.NewFuncImport(objects.FrameStatic, defTrimSpace, 1, s.funcStringToString(strings.TrimSpace)))
	s.Add(defTrimSuffix, factory.NewFuncImport(objects.FrameStatic, defTrimSuffix, 2, s.funcStringStringToString(strings.TrimSuffix)))

	return s
}

// Replace replaces occurrences of a substring within a string with the specified replacement string up to a given limit.
func (s *Strings) replace(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
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
	i4, err := gk.ToInt64Arg(3, args)
	if err != nil {
		return 0, nil, err
	}
	ret := strings.Replace(s1, s2, s3, int(i4))
	return 1, gk.NewString(frame, ret), nil
}

// Substring extracts a portion of a string based on the starting and ending indices provided as arguments.
func (s *Strings) substring(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	src, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	start, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	end := int64(len(src))
	if argsLen == 3 {
		end, err = gk.ToInt64Arg(2, args)
		if err != nil {
			return 0, nil, err
		}
	}
	if len(src) == 0 {
		return 1, gk.NewString(frame, ""), nil
	}
	if end < 0 {
		end = 0
	}
	if end > int64(len(src)) {
		end = int64(len(src))
	}
	if start < 0 {
		start = 0
	}
	if start > end {
		return 1, gk.NewString(frame, ""), nil
	}
	return 1, gk.NewString(frame, src[start:end]), nil
}

// PadLeft adds padding to the left of a string to ensure its total length is at least the specified value.
// The padding string can be optionally provided; otherwise, a space is used.
func (s *Strings) padLeft(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return 1, gk.NewString(frame, s1), nil
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = gk.ToStringArg(2, args)
		if err != nil {
			return 0, nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return 1, gk.NewString(frame, s1), nil
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := strings.Repeat(s3, int(padCount)) + s1
	return 1, gk.NewString(frame, retStr[int64(len(retStr))-i2:]), nil
}

// PadRight pads the input string on the right with a specified string or space until it reaches the desired length.
func (s *Strings) padRight(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return 1, gk.NewString(frame, s1), nil
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = gk.ToStringArg(2, args)
		if err != nil {
			return 0, nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return 1, gk.NewString(frame, s1), nil
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := s1 + strings.Repeat(s3, int(padCount))
	return 1, gk.NewString(frame, retStr[:i2]), nil
}

// Repeat repeats the input string a specified number of times and returns the concatenated result.
func (s *Strings) repeat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, strings.Repeat(s1, int(i2))), nil
}

// Join concatenates elements of an array into a single string, using a specified separator string.
func (s *Strings) join(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	sep, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	var elems []string
	switch arg0 := args[0].(type) {
	case *objects.Array:
		arg0.ForEach(func(i int, v objects.IObject) {
			elems = append(elems, v.AsString())
		})
	default:
		return 0, nil, objects.NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	return 1, gk.NewString(frame, strings.Join(elems, sep)), nil
}

// FuncIsOs creates a callable function using the provided transformation function that operates on string arguments.
// Returns a function that takes a frame and varargs, validates the input, applies the transformation, and returns the result.
func (s *Strings) funcStringToString(fn func(string) string) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		v := gk.NewString(frame, fn(s1))
		return 1, v, nil
	}
}

// funcStringToStrings transforms a string-to-string-slice function into a Invocable that operates on an IObject and returns an array.
func (s *Strings) funcStringToStrings(fn func(string) []string) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		res := fn(s1)
		obj := gk.NewArray(frame, nil)
		arr, ok := obj.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, elem := range res {
			v := gk.NewString(frame, elem)
			arr.Append(v)
		}
		return 1, arr, nil
	}
}

// funcStringToStringError converts a function with a string input and output into a Invocable for use in the current framework.
// It validates arguments, converts inputs/outputs to/from IObject, and handles errors gracefully.
func (s *Strings) funcStringToStringError(fn func(string) (string, error)) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return 0, gk.NewError(frame, err.Error()), nil
		}
		v := gk.NewString(frame, res)
		return 1, v, nil
	}
}

// funcStringStringToStrings creates a callable function wrapping a string-transforming function and returns results as an array object.
func (s *Strings) funcStringStringToStrings(fn func(string, string) []string) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		s2, err := gk.ToStringArg(1, args)
		if err != nil {
			return 0, nil, err
		}
		obj := gk.NewArray(frame, nil)
		arr, ok := obj.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, res := range fn(s1, s2) {
			v := gk.NewString(frame, res)
			arr.Append(v)
		}
		return 1, arr, nil
	}
}

// funcStringStringIntToStrings wraps a function of type func(string, string, int) []string into a Invocable compatible function.
func (s *Strings) funcStringStringIntToStrings(fn func(string, string, int) []string) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		s2, err := gk.ToStringArg(1, args)
		if err != nil {
			return 0, nil, err
		}
		i3, err := gk.ToInt64Arg(2, args)
		if err != nil {
			return 0, nil, err
		}
		obj := gk.NewArray(frame, nil)
		arr, ok := obj.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, res := range fn(s1, s2, int(i3)) {
			v := gk.NewString(frame, res)
			arr.Append(v)
		}
		return 1, arr, nil
	}
}

// funcStringStringToInt wraps a function accepting two string arguments and returning an int, into a Invocable type function.
func (s *Strings) funcStringStringToInt(fn func(string, string) int) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		s2, err := gk.ToStringArg(1, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewInt(frame, int64(fn(s1, s2))), nil
	}
}

// funcStringStringToString wraps a function that takes two strings and returns a string, converting it into a Invocable handler.
func (s *Strings) funcStringStringToString(fn func(string, string) string) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		s2, err := gk.ToStringArg(1, args)
		if err != nil {
			return 0, nil, err
		}
		v := gk.NewString(frame, fn(s1, s2))
		return 1, v, nil
	}
}

// funcStringStringToBool creates a Invocable that evaluates a function on two string arguments extracted from IObject inputs.
// It returns an IObject representing true or false based on the function's result or an error if arguments are invalid.
func (s *Strings) funcStringStringToBool(fn func(string, string) bool) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		s2, err := gk.ToStringArg(1, args)
		if err != nil {
			return 0, nil, err
		}
		if fn(s1, s2) {
			return 1, gk.TrueValue(), nil
		}
		return 1, gk.FalseValue(), nil
	}
}

// funcIntToIntError wraps a function matching func(string) (int, error) into a Invocable, enabling its usage with IObject arguments.
// It validates the argument count, converts the first argument to a string, applies the given function, and returns the result.
// If an error occurs during argument conversion or function execution, it returns appropriate error objects or messages.
func (s *Strings) funcIntToIntError(fn func(string) (int, error)) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (retCount uint, ret objects.IObject, err error) {
		s1, err := gk.ToStringArg(0, args)
		if err != nil {
			return 0, nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return 0, gk.NewError(frame, err.Error()), nil
		}
		return 1, gk.NewInt(frame, int64(res)), nil
	}
}

// funcIntToString wraps a function from int to string into a Invocable, ensuring compatibility with the IObject interface.
func (s *Strings) funcIntToString(fn func(int) string) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (retCount uint, ret objects.IObject, err error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		s := fn(int(i1))
		v := gk.NewString(frame, s)
		return 1, v, nil
	}
}

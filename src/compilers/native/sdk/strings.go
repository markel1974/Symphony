package sdk

import (
	"fmt"
	"strings"
	//"unicode/utf8"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewStrings)
}

// Strings provides a collection of string operations and functionality wrapped in a module.
// It includes functions for manipulation, comparison, trimming, padding, splitting, and other string-related utilities.
type Strings struct {
	*bytecode.Package
}

// NewStrings creates and returns a new instance of Strings with a preconfigured map of string utility functions.
func NewStrings(factory objects.IGateKeeper) bytecode.IPackage {
	s := &Strings{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Join", 2, s.join),
		factory.NewFuncImport(objects.FrameStatic, "Repeat", 2, s.repeat),
		factory.NewFuncImport(objects.FrameStatic, "Replace", 3, s.replace),
		factory.NewFuncImport(objects.FrameStatic, "Substring", 3, s.substring),
		factory.NewFuncImport(objects.FrameStatic, "PadLeft", -1, s.padLeft),
		factory.NewFuncImport(objects.FrameStatic, "PadRight", -1, s.padRight),
		factory.NewFuncImport(objects.FrameStatic, "Compare", 2, s.funcStringStringToInt(strings.Compare)),
		factory.NewFuncImport(objects.FrameStatic, "Contains", 2, s.funcStringStringToBool(strings.Contains)),
		factory.NewFuncImport(objects.FrameStatic, "ContainsAny", 2, s.funcStringStringToBool(strings.ContainsAny)),
		factory.NewFuncImport(objects.FrameStatic, "Count", 2, s.funcStringStringToInt(strings.Count)),
		factory.NewFuncImport(objects.FrameStatic, "EqualFold", 2, s.funcStringStringToBool(strings.EqualFold)),
		factory.NewFuncImport(objects.FrameStatic, "Fields", 1, s.funcStringToStrings(strings.Fields)),
		factory.NewFuncImport(objects.FrameStatic, "HasPrefix", 2, s.funcStringStringToBool(strings.HasPrefix)),
		factory.NewFuncImport(objects.FrameStatic, "HasSuffix", 2, s.funcStringStringToBool(strings.HasSuffix)),
		factory.NewFuncImport(objects.FrameStatic, "index", 2, s.funcStringStringToInt(strings.Index)),
		factory.NewFuncImport(objects.FrameStatic, "IndexAny", 2, s.funcStringStringToInt(strings.IndexAny)),
		factory.NewFuncImport(objects.FrameStatic, "LastIndex", 2, s.funcStringStringToInt(strings.LastIndex)),
		factory.NewFuncImport(objects.FrameStatic, "LastIndexAny", 2, s.funcStringStringToInt(strings.LastIndexAny)),
		factory.NewFuncImport(objects.FrameStatic, "Split", 2, s.funcStringStringToStrings(strings.Split)),
		factory.NewFuncImport(objects.FrameStatic, "SplitAfter", 2, s.funcStringStringToStrings(strings.SplitAfter)),
		factory.NewFuncImport(objects.FrameStatic, "SplitAfterN", 2, s.funcStringStringIntToStrings(strings.SplitAfterN)),
		factory.NewFuncImport(objects.FrameStatic, "SplitN", 2, s.funcStringStringIntToStrings(strings.SplitN)),
		factory.NewFuncImport(objects.FrameStatic, "Title", 1, s.funcStringToString(strings.Title)),
		factory.NewFuncImport(objects.FrameStatic, "ToLower", 1, s.funcStringToString(strings.ToLower)),
		factory.NewFuncImport(objects.FrameStatic, "ToTitle", 1, s.funcStringToString(strings.ToTitle)),
		factory.NewFuncImport(objects.FrameStatic, "ToUpper", 1, s.funcStringToString(strings.ToUpper)),
		factory.NewFuncImport(objects.FrameStatic, "Trim", 2, s.funcStringStringToString(strings.Trim)),
		factory.NewFuncImport(objects.FrameStatic, "TrimLeft", 2, s.funcStringStringToString(strings.TrimLeft)),
		factory.NewFuncImport(objects.FrameStatic, "TrimPrefix", 2, s.funcStringStringToString(strings.TrimPrefix)),
		factory.NewFuncImport(objects.FrameStatic, "TrimRight", 2, s.funcStringStringToString(strings.TrimRight)),
		factory.NewFuncImport(objects.FrameStatic, "TrimSpace", 1, s.funcStringToString(strings.TrimSpace)),
		factory.NewFuncImport(objects.FrameStatic, "TrimSuffix", 2, s.funcStringStringToString(strings.TrimSuffix)),
	}
	s.Package = bytecode.NewPackage("strings", container, nil)
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
		for _, a := range arg0.Values() {
			elems = append(elems, a.AsString())
		}
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

package sdk

import (
	"strings"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewStrings)
}

// Strings provides a collection of string operations and functionality wrapped in a module.
// It includes functions for manipulation, comparison, trimming, padding, splitting, and other string-related utilities.
type Strings struct {
	container map[string]objects.IObject
}

// NewStrings creates and returns a new instance of Strings with a preconfigured map of string utility functions.
func NewStrings(factory objects.IGateKeeper) IPackage {
	s := &Strings{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Compare", factory.FuncIssOi(strings.Compare)),
		factory.NewFuncImport(objects.FrameStatic, "Contains", factory.FuncIssOb(strings.Contains)),
		factory.NewFuncImport(objects.FrameStatic, "ContainsAny", factory.FuncIssOb(strings.ContainsAny)),
		factory.NewFuncImport(objects.FrameStatic, "Count", factory.FuncIssOi(strings.Count)),
		factory.NewFuncImport(objects.FrameStatic, "EqualFold", factory.FuncIssOb(strings.EqualFold)),
		factory.NewFuncImport(objects.FrameStatic, "Fields", factory.FuncIsOsS(strings.Fields)),
		factory.NewFuncImport(objects.FrameStatic, "HasPrefix", factory.FuncIssOb(strings.HasPrefix)),
		factory.NewFuncImport(objects.FrameStatic, "HasSuffix", factory.FuncIssOb(strings.HasSuffix)),
		factory.NewFuncImport(objects.FrameStatic, "Index", factory.FuncIssOi(strings.Index)),
		factory.NewFuncImport(objects.FrameStatic, "IndexAny", factory.FuncIssOi(strings.IndexAny)),
		factory.NewFuncImport(objects.FrameStatic, "Join", s.join),
		factory.NewFuncImport(objects.FrameStatic, "LastIndex", factory.FuncIssOi(strings.LastIndex)),
		factory.NewFuncImport(objects.FrameStatic, "LastIndexAny", factory.FuncIssOi(strings.LastIndexAny)),
		factory.NewFuncImport(objects.FrameStatic, "Repeat", s.repeat),
		factory.NewFuncImport(objects.FrameStatic, "Replace", s.replace),
		factory.NewFuncImport(objects.FrameStatic, "Substring", s.substring),
		factory.NewFuncImport(objects.FrameStatic, "Split", factory.FuncIssOsS(strings.Split)),
		factory.NewFuncImport(objects.FrameStatic, "SplitAfter", factory.FuncIssOsS(strings.SplitAfter)),
		factory.NewFuncImport(objects.FrameStatic, "SplitAfterN", factory.FuncIssiOsS(strings.SplitAfterN)),
		factory.NewFuncImport(objects.FrameStatic, "SplitN", factory.FuncIssiOsS(strings.SplitN)),
		factory.NewFuncImport(objects.FrameStatic, "Title", factory.FuncIsOs(strings.Title)),
		factory.NewFuncImport(objects.FrameStatic, "ToLower", factory.FuncIsOs(strings.ToLower)),
		factory.NewFuncImport(objects.FrameStatic, "ToTitle", factory.FuncIsOs(strings.ToTitle)),
		factory.NewFuncImport(objects.FrameStatic, "ToUpper", factory.FuncIsOs(strings.ToUpper)),
		factory.NewFuncImport(objects.FrameStatic, "PadLeft", s.padLeft),
		factory.NewFuncImport(objects.FrameStatic, "PadRight", s.padRight),
		factory.NewFuncImport(objects.FrameStatic, "Trim", factory.FuncIssOs(strings.Trim)),
		factory.NewFuncImport(objects.FrameStatic, "TrimLeft", factory.FuncIssOs(strings.TrimLeft)),
		factory.NewFuncImport(objects.FrameStatic, "TrimPrefix", factory.FuncIssOs(strings.TrimPrefix)),
		factory.NewFuncImport(objects.FrameStatic, "TrimRight", factory.FuncIssOs(strings.TrimRight)),
		factory.NewFuncImport(objects.FrameStatic, "TrimSpace", factory.FuncIsOs(strings.TrimSpace)),
		factory.NewFuncImport(objects.FrameStatic, "TrimSuffix", factory.FuncIssOs(strings.TrimSuffix)),
	}
	s.container = BuildContainer(container, nil)
	return s
}

// Name returns the name of the Strings structure as a string.
func (s *Strings) Name() string {
	return "strings"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (s *Strings) Get(name string) (objects.IObject, bool) {
	v, ok := s.container[name]
	return v, ok
}

// Replace replaces occurrences of a substring within a string with the specified replacement string up to a given limit.
func (s *Strings) replace(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) != 4 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	s3, err := gk.ToStringArg(2, args[2])
	if err != nil {
		return 0, nil, err
	}
	i4, err := gk.ToInt64Arg(3, args[3])
	if err != nil {
		return 0, nil, err
	}
	ret := s.stringsReplace(gk, frame, s1, s2, s3, int(i4))
	return 1, gk.NewString(frame, ret), nil
}

// Substring extracts a portion of a string based on the starting and ending indices provided as arguments.
func (s *Strings) substring(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	strlen := int64(len(s1))
	i3 := strlen
	if argsLen == 3 {
		i3, err = gk.ToInt64Arg(2, args[2])
		if err != nil {
			return 0, nil, err
		}
	}
	if i2 > i3 {
		return 0, nil, objects.ErrInvalidIndexType
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
	return 1, gk.NewString(frame, s1[i2:i3]), nil
}

// PadLeft adds padding to the left of a string to ensure its total length is at least the specified value.
// The padding string can be optionally provided; otherwise, a space is used.
func (s *Strings) padLeft(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return 1, gk.NewString(frame, s1), nil
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = gk.ToStringArg(2, args[2])
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
	if argsLen != 2 && argsLen != 3 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return 1, gk.NewString(frame, s1), nil
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = gk.ToStringArg(2, args[2])
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
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, strings.Repeat(s1, int(i2))), nil
}

// Join concatenates elements of an array into a single string, using a specified separator string.
func (s *Strings) join(gk objects.IGateKeeper, frame int, args ...objects.IObject) (retCount uint, ret objects.IObject, err error) {
	if len(args) != 2 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	var sLen int
	var ss1 []string
	switch arg0 := args[0].(type) {
	case *objects.Array:
		for idx, a := range arg0.Values() {
			as, err := gk.ToStringArg(idx, a)
			if err != nil {
				return 0, nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	default:
		return 0, nil, objects.NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	s2, err := gk.ToStringArg(1, args[1])
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, strings.Join(ss1, s2)), nil
}

// stringsReplace replaces up to n occurrences of the substring old with the substring new in the input string str.
// Returns the modified string and a boolean indicating success. Returns the original string if no replacement is needed.
func (s *Strings) stringsReplace(_ objects.IGateKeeper, _ int, str string, old string, new string, n int) string {
	if old == new || n == 0 {
		return str
	}
	if m := strings.Count(str, old); m == 0 {
		return str
	} else if n < 0 || m < n {
		n = m
	}
	t := make([]byte, len(str)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(str[start:])
				j += wid
			}
		} else {
			j += strings.Index(str[start:], old)
		}
		ssj := str[start:j]
		w += copy(t[w:], ssj)
		w += copy(t[w:], new)
		start = j + len(old)
	}
	ss := str[start:]
	w += copy(t[w:], ss)
	return string(t[0:w])
}

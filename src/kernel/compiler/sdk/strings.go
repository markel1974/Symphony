package sdk

import (
	"strings"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Strings provides a collection of string operations and functionality wrapped in a module.
// It includes functions for manipulation, comparison, trimming, padding, splitting, and other string-related utilities.
type Strings struct {
	factory *objects.Factory
	*Package
}

// NewStrings creates and returns a new instance of Strings with a preconfigured map of string utility functions.
func NewStrings(factory *objects.Factory) *Strings {
	s := &Strings{factory: factory}
	container := []*objects.FuncPackage{
		factory.NewFuncPackage(objects.FuncPackageDef, "Compare", factory.FuncIssOi(strings.Compare)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Contains", factory.FuncIssOb(strings.Contains)),
		factory.NewFuncPackage(objects.FuncPackageDef, "ContainsAny", factory.FuncIssOb(strings.ContainsAny)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Count", factory.FuncIssOi(strings.Count)),
		factory.NewFuncPackage(objects.FuncPackageDef, "EqualFold", factory.FuncIssOb(strings.EqualFold)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Fields", factory.FuncIsOsS(strings.Fields)),
		factory.NewFuncPackage(objects.FuncPackageDef, "HasPrefix", factory.FuncIssOb(strings.HasPrefix)),
		factory.NewFuncPackage(objects.FuncPackageDef, "HasSuffix", factory.FuncIssOb(strings.HasSuffix)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Index", factory.FuncIssOi(strings.Index)),
		factory.NewFuncPackage(objects.FuncPackageDef, "IndexAny", factory.FuncIssOi(strings.IndexAny)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Join", s.join),
		factory.NewFuncPackage(objects.FuncPackageDef, "LastIndex", factory.FuncIssOi(strings.LastIndex)),
		factory.NewFuncPackage(objects.FuncPackageDef, "LastIndexAny", factory.FuncIssOi(strings.LastIndexAny)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Repeat", s.repeat),
		factory.NewFuncPackage(objects.FuncPackageDef, "Replace", s.replace),
		factory.NewFuncPackage(objects.FuncPackageDef, "Substring", s.Substring),
		factory.NewFuncPackage(objects.FuncPackageDef, "Split", factory.FuncIssOsS(strings.Split)),
		factory.NewFuncPackage(objects.FuncPackageDef, "SplitAfter", factory.FuncIssOsS(strings.SplitAfter)),
		factory.NewFuncPackage(objects.FuncPackageDef, "SplitAfterN", factory.FuncIssiOsS(strings.SplitAfterN)),
		factory.NewFuncPackage(objects.FuncPackageDef, "SplitN", factory.FuncIssiOsS(strings.SplitN)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Title", factory.FuncIsOs(strings.Title)),
		factory.NewFuncPackage(objects.FuncPackageDef, "ToLower", factory.FuncIsOs(strings.ToLower)),
		factory.NewFuncPackage(objects.FuncPackageDef, "ToTitle", factory.FuncIsOs(strings.ToTitle)),
		factory.NewFuncPackage(objects.FuncPackageDef, "ToUpper", factory.FuncIsOs(strings.ToUpper)),
		factory.NewFuncPackage(objects.FuncPackageDef, "PadLeft", s.padLeft),
		factory.NewFuncPackage(objects.FuncPackageDef, "PadRight", s.padRight),
		factory.NewFuncPackage(objects.FuncPackageDef, "Trim", factory.FuncIssOs(strings.Trim)),
		factory.NewFuncPackage(objects.FuncPackageDef, "TrimLeft", factory.FuncIssOs(strings.TrimLeft)),
		factory.NewFuncPackage(objects.FuncPackageDef, "TrimPrefix", factory.FuncIssOs(strings.TrimPrefix)),
		factory.NewFuncPackage(objects.FuncPackageDef, "TrimRight", factory.FuncIssOs(strings.TrimRight)),
		factory.NewFuncPackage(objects.FuncPackageDef, "TrimSpace", factory.FuncIsOs(strings.TrimSpace)),
		factory.NewFuncPackage(objects.FuncPackageDef, "TrimSuffix", factory.FuncIssOs(strings.TrimSuffix)),
	}
	s.Package = NewPackage("strings", container, nil)
	return s
}

// Replace replaces occurrences of a substring within a string with the specified replacement string up to a given limit.
func (s *Strings) replace(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := s.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := s.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	s3, err := s.factory.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := s.factory.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	ret, ok := s.stringsReplace(frame, s1, s2, s3, int(i4))
	if !ok {
		return nil, objects.ErrStringLimit
	}
	return s.factory.NewString(frame, ret)
}

// Substring extracts a portion of a string based on the starting and ending indices provided as arguments.
func (s *Strings) Substring(frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := s.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := s.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	strlen := int64(len(s1))
	i3 := strlen
	if argsLen == 3 {
		i3, err = s.factory.ToInt64Arg(2, args[2])
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
	return s.factory.NewString(frame, s1[i2:i3])
}

// PadLeft adds padding to the left of a string to ensure its total length is at least the specified value.
// The padding string can be optionally provided; otherwise, a space is used.
func (s *Strings) padLeft(frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := s.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := s.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	if i2 > objects.MaxStringLen {
		return nil, objects.ErrStringLimit
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return s.factory.NewString(frame, s1)
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = s.factory.ToStringArg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return s.factory.NewString(frame, s1)
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := strings.Repeat(s3, int(padCount)) + s1
	return s.factory.NewString(frame, retStr[int64(len(retStr))-i2:])
}

// PadRight pads the input string on the right with a specified string or space until it reaches the desired length.
func (s *Strings) padRight(frame int, args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := s.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := s.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return s.factory.NewString(frame, s1)
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = s.factory.ToStringArg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return s.factory.NewString(frame, s1)
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := s1 + strings.Repeat(s3, int(padCount))
	return s.factory.NewString(frame, retStr[:i2])
}

// Repeat repeats the input string a specified number of times and returns the concatenated result.
func (s *Strings) repeat(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := s.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := s.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return s.factory.NewString(frame, strings.Repeat(s1, int(i2)))
}

// Join concatenates elements of an array into a single string, using a specified separator string.
func (s *Strings) join(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	var sLen int
	var ss1 []string
	switch arg0 := args[0].(type) {
	case *objects.Array:
		for idx, a := range arg0.Values() {
			as, err := s.factory.ToStringArg(idx, a)
			if err != nil {
				return nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	case *objects.ArrayImmutable:
		for idx, a := range arg0.Values() {
			as, err := s.factory.ToStringArg(idx, a)
			if err != nil {
				return nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	default:
		return nil, objects.NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	s2, err := s.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	return s.factory.NewString(frame, strings.Join(ss1, s2))
}

// stringsReplace replaces up to n occurrences of the substring old with the substring new in the input string str.
// Returns the modified string and a boolean indicating success. Returns the original string if no replacement is needed.
func (s *Strings) stringsReplace(_ int, str string, old string, new string, n int) (string, bool) {
	if old == new || n == 0 {
		return str, true // avoid allocation
	}
	if m := strings.Count(str, old); m == 0 {
		return str, true
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
		if w+len(ssj)+len(new) > objects.MaxStringLen {
			return "", false
		}

		w += copy(t[w:], ssj)
		w += copy(t[w:], new)
		start = j + len(old)
	}
	ss := str[start:]
	if w+len(ss) > objects.MaxStringLen {
		return "", false
	}
	w += copy(t[w:], ss)
	return string(t[0:w]), true
}

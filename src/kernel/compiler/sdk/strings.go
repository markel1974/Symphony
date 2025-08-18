package sdk

import (
	"strings"
	"unicode/utf8"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Strings provides a collection of string operations and functionality wrapped in a module.
// It includes functions for manipulation, comparison, trimming, padding, splitting, and other string-related utilities.
type Strings struct {
	*Package
}

// NewStrings creates and returns a new instance of Strings with a preconfigured map of string utility functions.
func NewStrings() *Strings {
	s := &Strings{}
	container := []*objects.FuncPackage{
		objects.NewFuncPackage(objects.FuncPackageDef, "Compare", objects.FuncIssOi(strings.Compare)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Contains", objects.FuncIssOb(strings.Contains)),
		objects.NewFuncPackage(objects.FuncPackageDef, "ContainsAny", objects.FuncIssOb(strings.ContainsAny)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Count", objects.FuncIssOi(strings.Count)),
		objects.NewFuncPackage(objects.FuncPackageDef, "EqualFold", objects.FuncIssOb(strings.EqualFold)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Fields", objects.FuncIsOsS(strings.Fields)),
		objects.NewFuncPackage(objects.FuncPackageDef, "HasPrefix", objects.FuncIssOb(strings.HasPrefix)),
		objects.NewFuncPackage(objects.FuncPackageDef, "HasSuffix", objects.FuncIssOb(strings.HasSuffix)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Index", objects.FuncIssOi(strings.Index)),
		objects.NewFuncPackage(objects.FuncPackageDef, "IndexAny", objects.FuncIssOi(strings.IndexAny)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Join", s.Join),
		objects.NewFuncPackage(objects.FuncPackageDef, "LastIndex", objects.FuncIssOi(strings.LastIndex)),
		objects.NewFuncPackage(objects.FuncPackageDef, "LastIndexAny", objects.FuncIssOi(strings.LastIndexAny)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Repeat", s.Repeat),
		objects.NewFuncPackage(objects.FuncPackageDef, "Replace", s.Replace),
		objects.NewFuncPackage(objects.FuncPackageDef, "Substring", s.Substring),
		objects.NewFuncPackage(objects.FuncPackageDef, "Split", objects.FuncIssOsS(strings.Split)),
		objects.NewFuncPackage(objects.FuncPackageDef, "SplitAfter", objects.FuncIssOsS(strings.SplitAfter)),
		objects.NewFuncPackage(objects.FuncPackageDef, "SplitAfterN", objects.FuncIssiOsS(strings.SplitAfterN)),
		objects.NewFuncPackage(objects.FuncPackageDef, "SplitN", objects.FuncIssiOsS(strings.SplitN)),
		objects.NewFuncPackage(objects.FuncPackageDef, "Title", objects.FuncIsOs(strings.Title)),
		objects.NewFuncPackage(objects.FuncPackageDef, "ToLower", objects.FuncIsOs(strings.ToLower)),
		objects.NewFuncPackage(objects.FuncPackageDef, "ToTitle", objects.FuncIsOs(strings.ToTitle)),
		objects.NewFuncPackage(objects.FuncPackageDef, "ToUpper", objects.FuncIsOs(strings.ToUpper)),
		objects.NewFuncPackage(objects.FuncPackageDef, "PadLeft", s.PadLeft),
		objects.NewFuncPackage(objects.FuncPackageDef, "PadRight", s.PadRight),
		objects.NewFuncPackage(objects.FuncPackageDef, "Trim", objects.FuncIssOs(strings.Trim)),
		objects.NewFuncPackage(objects.FuncPackageDef, "TrimLeft", objects.FuncIssOs(strings.TrimLeft)),
		objects.NewFuncPackage(objects.FuncPackageDef, "TrimPrefix", objects.FuncIssOs(strings.TrimPrefix)),
		objects.NewFuncPackage(objects.FuncPackageDef, "TrimRight", objects.FuncIssOs(strings.TrimRight)),
		objects.NewFuncPackage(objects.FuncPackageDef, "TrimSpace", objects.FuncIsOs(strings.TrimSpace)),
		objects.NewFuncPackage(objects.FuncPackageDef, "TrimSuffix", objects.FuncIssOs(strings.TrimSuffix)),
	}
	s.Package = NewPackage("strings", container, nil)
	return s
}

// Replace replaces occurrences of a substring within a string with the specified replacement string up to a given limit.
func (s *Strings) Replace(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
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
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	ret, ok := s.stringsReplace(s1, s2, s3, int(i4))
	if !ok {
		return nil, objects.ErrStringLimit
	}
	return objects.NewString(ret)
}

// Substring extracts a portion of a string based on the starting and ending indices provided as arguments.
func (s *Strings) Substring(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	strlen := int64(len(s1))
	i3 := strlen
	if argsLen == 3 {
		i3, err = objects.ToInt64Arg(2, args[2])
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
	return objects.NewString(s1[i2:i3])
}

// PadLeft adds padding to the left of a string to ensure its total length is at least the specified value.
// The padding string can be optionally provided; otherwise, a space is used.
func (s *Strings) PadLeft(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	if i2 > objects.MaxStringLen {
		return nil, objects.ErrStringLimit
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return objects.NewString(s1)
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = objects.ToStringArg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return objects.NewString(s1)
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := strings.Repeat(s3, int(padCount)) + s1
	return objects.NewString(retStr[int64(len(retStr))-i2:])
}

// PadRight pads the input string on the right with a specified string or space until it reaches the desired length.
func (s *Strings) PadRight(args ...objects.IObject) (objects.IObject, error) {
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	sLen := int64(len(s1))
	if sLen >= i2 {
		return objects.NewString(s1)
	}
	s3 := " "
	if argsLen == 3 {
		s3, err = objects.ToStringArg(2, args[2])
		if err != nil {
			return nil, err
		}
	}
	padStrLen := int64(len(s3))
	if padStrLen == 0 {
		return objects.NewString(s1)
	}
	padCount := ((i2 - padStrLen) / padStrLen) + 1
	retStr := s1 + strings.Repeat(s3, int(padCount))
	return objects.NewString(retStr[:i2])
}

// Repeat repeats the input string a specified number of times and returns the concatenated result.
func (s *Strings) Repeat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strings.Repeat(s1, int(i2)))
}

// Join concatenates elements of an array into a single string, using a specified separator string.
func (s *Strings) Join(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	var sLen int
	var ss1 []string
	switch arg0 := args[0].(type) {
	case *objects.Array:
		for idx, a := range arg0.Values() {
			as, err := objects.ToStringArg(idx, a)
			if err != nil {
				return nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	case *objects.ArrayImmutable:
		for idx, a := range arg0.Values() {
			as, err := objects.ToStringArg(idx, a)
			if err != nil {
				return nil, err
			}
			sLen += len(as)
			ss1 = append(ss1, as)
		}
	default:
		return nil, objects.NewInvalidArgumentError(0, "array", args[0].TypeName())
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strings.Join(ss1, s2))
}

// stringsReplace replaces up to n occurrences of the substring old with the substring new in the input string str.
// Returns the modified string and a boolean indicating success. Returns the original string if no replacement is needed.
func (s *Strings) stringsReplace(str string, old string, new string, n int) (string, bool) {
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

package sdk

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Strconv is a type that provides a module containing string conversion functions implemented using strconv.
type Strconv struct {
	factory *objects.Factory
	*Package
}

// NewStrconv initializes and returns a pointer to a new Strconv instance containing predefined module functions.
func NewStrconv(factory *objects.Factory) *Strconv {
	s := &Strconv{
		factory: factory,
	}
	container := []*objects.FuncPackage{
		factory.NewFuncPackage(objects.FuncPackageDef, "Atoi", factory.FuncIsOie(strconv.Atoi)),
		factory.NewFuncPackage(objects.FuncPackageDef, "FormatBool", s.formatBool),
		factory.NewFuncPackage(objects.FuncPackageDef, "FormatFloat", s.formatFloat),
		factory.NewFuncPackage(objects.FuncPackageDef, "FormatInt", s.formatInt),
		factory.NewFuncPackage(objects.FuncPackageDef, "Itoa", factory.FuncIiOs(strconv.Itoa)),
		factory.NewFuncPackage(objects.FuncPackageDef, "ParseBool", s.parseBool),
		factory.NewFuncPackage(objects.FuncPackageDef, "ParseFloat", s.parseFloat),
		factory.NewFuncPackage(objects.FuncPackageDef, "ParseNumber", s.parseNumber),
		factory.NewFuncPackage(objects.FuncPackageDef, "ParseInt", s.parseInt),
		factory.NewFuncPackage(objects.FuncPackageDef, "Quote", factory.FuncIsOs(strconv.Quote)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Unquote", factory.FuncIsOse(strconv.Unquote)),
	}
	s.Package = NewPackage("strconv", container, nil)
	return s
}

// FormatBool converts a boolean argument to its string representation ("true" or "false"). Returns an error if the argument is invalid.
func (s *Strconv) formatBool(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	b1, err := s.factory.ToBoolArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if b1 {
		return s.factory.NewString(frame, "true")
	}
	return s.factory.NewString(frame, "false")
}

// FormatFloat converts a float64 into a string representation according to the specified format, precision, and bit size.
func (s *Strconv) formatFloat(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	f1, err := s.factory.ToFloat64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := s.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := s.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := s.factory.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	return s.factory.NewString(frame, strconv.FormatFloat(f1, s2[0], int(i3), int(i4)))
}

// FormatInt formats an int64 number as a string in the specified base, provided by the second argument.
func (s *Strconv) formatInt(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := s.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := s.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return s.factory.NewString(frame, strconv.FormatInt(i1, int(i2)))
}

// ParseBool parses a string representation of a boolean value and returns the corresponding boolean object or an error.
func (s *Strconv) parseBool(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		err = objects.ErrWrongNumArguments
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		err = objects.NewInvalidArgumentError(0, "string", args[0].TypeName())
		return
	}
	parsed, err := strconv.ParseBool(s1.Value())
	if err != nil {
		ret = s.factory.NewError(frame, err.Error())
		return
	}
	if parsed {
		ret = s.factory.TrueValue()
	} else {
		ret = s.factory.FalseValue()
	}
	return
}

// ParseFloat parses a string argument as a floating-point number using the specified precision and returns the result.
func (s *Strconv) parseFloat(frame int, args ...objects.IObject) (objects.IObject, error) {
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
	parsed, err := strconv.ParseFloat(s1, int(i2))
	if err != nil {
		return s.factory.NewError(frame, err.Error()), nil
	}
	return s.factory.NewFloat(frame, parsed), nil
}

// ParseNumber extracts and parses numeric values from a string, returning a float object or an error if parsing fails.
func (s *Strconv) parseNumber(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := s.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var target []rune
	for _, p := range s1 {
		if (p >= '0' && p <= '9') || p == '.' || p == '+' || p == '-' {
			target = append(target, p)
		}
	}
	parsed, err := strconv.ParseFloat(string(target), 64)
	if err != nil {
		return nil, err
	}
	return s.factory.NewFloat(frame, parsed), nil
}

// ParseInt converts a string argument to an integer with the specified base and bit size after validating arguments.
func (s *Strconv) parseInt(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
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
	i3, err := s.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseInt(s1, int(i2), int(i3))
	if err != nil {
		return s.factory.NewError(frame, err.Error()), nil
	}
	return s.factory.NewInt(frame, parsed), nil
}

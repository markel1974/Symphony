package sdk

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewStrconv)
}

// Strconv is a type that provides a module containing string conversion functions implemented using strconv.
type Strconv struct {
	container map[string]objects.IObject
}

// NewStrconv initializes and returns a pointer to a new Strconv instance containing predefined module functions.
func NewStrconv(factory objects.IGateKeeper) IPackage {
	s := &Strconv{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Atoi", factory.FuncIsOie(strconv.Atoi)),
		factory.NewFuncImport(objects.FrameStatic, "FormatBool", s.formatBool),
		factory.NewFuncImport(objects.FrameStatic, "FormatFloat", s.formatFloat),
		factory.NewFuncImport(objects.FrameStatic, "FormatInt", s.formatInt),
		factory.NewFuncImport(objects.FrameStatic, "Itoa", factory.FuncIiOs(strconv.Itoa)),
		factory.NewFuncImport(objects.FrameStatic, "ParseBool", s.parseBool),
		factory.NewFuncImport(objects.FrameStatic, "ParseFloat", s.parseFloat),
		factory.NewFuncImport(objects.FrameStatic, "ParseNumber", s.parseNumber),
		factory.NewFuncImport(objects.FrameStatic, "ParseInt", s.parseInt),
		factory.NewFuncImport(objects.FrameStatic, "Quote", factory.FuncIsOs(strconv.Quote)),
		factory.NewFuncImport(objects.FrameStatic, "Unquote", factory.FuncIsOse(strconv.Unquote)),
	}
	s.container = BuildContainer(container, nil)
	return s
}

// Name returns the name of the Strconv module as a string.
func (s *Strconv) Name() string {
	return "strconv"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (s *Strconv) Get(name string) (objects.IObject, bool) {
	v, ok := s.container[name]
	return v, ok
}

// FormatBool converts a boolean argument to its string representation ("true" or "false"). Returns an error if the argument is invalid.
func (s *Strconv) formatBool(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	b1, err := gk.ToBoolArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if b1 {
		return gk.NewString(frame, "true"), nil
	}
	return gk.NewString(frame, "false"), nil
}

// FormatFloat converts a float64 into a string representation according to the specified format, precision, and bit size.
func (s *Strconv) formatFloat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	f1, err := gk.ToFloat64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := gk.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := gk.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := gk.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	return gk.NewString(frame, strconv.FormatFloat(f1, s2[0], int(i3), int(i4))), nil
}

// FormatInt formats an int64 number as a string in the specified base, provided by the second argument.
func (s *Strconv) formatInt(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return gk.NewString(frame, strconv.FormatInt(i1, int(i2))), nil
}

// ParseBool parses a string representation of a boolean value and returns the corresponding boolean object or an error.
func (s *Strconv) parseBool(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		err = objects.ErrInvalidArgumentsNumber
		return
	}
	s1, ok := args[0].(*objects.String)
	if !ok {
		err = objects.NewInvalidArgumentError(0, "string", args[0].TypeName())
		return
	}
	parsed, err := strconv.ParseBool(s1.Value())
	if err != nil {
		ret = gk.NewError(frame, err.Error())
		return
	}
	if parsed {
		ret = gk.TrueValue()
	} else {
		ret = gk.FalseValue()
	}
	return
}

// ParseFloat parses a string argument as a floating-point number using the specified precision and returns the result.
func (s *Strconv) parseFloat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseFloat(s1, int(i2))
	if err != nil {
		return gk.NewError(frame, err.Error()), nil
	}
	return gk.NewFloat(frame, parsed), nil
}

// ParseNumber extracts and parses numeric values from a string, returning a float object or an error if parsing fails.
func (s *Strconv) parseNumber(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
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
	return gk.NewFloat(frame, parsed), nil
}

// ParseInt converts a string argument to an integer with the specified base and bit size after validating arguments.
func (s *Strconv) parseInt(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := gk.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := gk.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseInt(s1, int(i2), int(i3))
	if err != nil {
		return gk.NewError(frame, err.Error()), nil
	}
	return gk.NewInt(frame, parsed), nil
}

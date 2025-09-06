package sdk

import (
	"strconv"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the strconv package by appending its registration function to the internal package list.
func init() {
	RegisterPackage(NewStrconv)
}

// Strconv represents a container for storing string conversion-related functions in a map-like structure.
type Strconv struct {
	*Package
}

// NewStrconv creates a new Strconv package using the provided IGateKeeper factory and initializes its container with functions.
func NewStrconv(factory objects.IGateKeeper) IPackage {
	s := &Strconv{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Atoi", 1, s.atoi),
		factory.NewFuncImport(objects.FrameStatic, "FormatBool", 1, s.formatBool),
		factory.NewFuncImport(objects.FrameStatic, "FormatFloat", 3, s.formatFloat),
		factory.NewFuncImport(objects.FrameStatic, "FormatInt", 2, s.formatInt),
		factory.NewFuncImport(objects.FrameStatic, "Itoa", 1, s.itoa),
		factory.NewFuncImport(objects.FrameStatic, "ParseBool", 1, s.parseBool),
		factory.NewFuncImport(objects.FrameStatic, "ParseFloat", 2, s.parseFloat),
		factory.NewFuncImport(objects.FrameStatic, "ParseNumber", 1, s.parseNumber),
		factory.NewFuncImport(objects.FrameStatic, "ParseInt", 3, s.parseInt),
		factory.NewFuncImport(objects.FrameStatic, "Quote", 1, s.quote),
		factory.NewFuncImport(objects.FrameStatic, "Unquote", 1, s.unquote),
	}
	s.Package = NewExternalPackage("strconv", container, nil)
	return s
}

// formatBool converts a boolean argument to its string representation ("true" or "false") and returns it as an object.
func (s *Strconv) formatBool(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	b1, err := gk.ToBoolArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	if b1 {
		return 1, gk.NewString(frame, "true"), nil
	}
	return 1, gk.NewString(frame, "false"), nil
}

// formatFloat formats a floating-point number according to the specified format, precision, and bit size.
// It takes 4 arguments: float64, format string, integer precision, and integer bit size, returning a string result.
func (s *Strconv) formatFloat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	f1, err := gk.ToFloat64Arg(0, args)
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
	i4, err := gk.ToInt64Arg(3, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, strconv.FormatFloat(f1, s2[0], int(i3), int(i4))), nil
}

// formatInt formats an integer as a string in a specified base using strconv.FormatInt. Accepts base and integer arguments.
func (s *Strconv) formatInt(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, strconv.FormatInt(i1, int(i2))), nil
}

// parseBool parses a single string argument into a boolean value and returns associated objects and error if any.
// It expects exactly one string argument and returns an error for invalid argument types or number of arguments.
func (s *Strconv) parseBool(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	parsed, err := strconv.ParseBool(s1)
	if err != nil {
		return 0, nil, err
	}
	if parsed {
		return 1, gk.TrueValue(), nil
	}
	return 1, gk.FalseValue(), nil
}

// parseFloat parses a given string into a floating-point number based on the specified precision.
func (s *Strconv) parseFloat(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	parsed, err := strconv.ParseFloat(s1, int(i2))
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewFloat(frame, parsed), nil
}

// parseNumber parses a numeric value from the input argument and returns it as a float object or an error if invalid.
func (s *Strconv) parseNumber(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	var target []rune
	for _, p := range s1 {
		if (p >= '0' && p <= '9') || p == '.' || p == '+' || p == '-' {
			target = append(target, p)
		}
	}
	parsed, err := strconv.ParseFloat(string(target), 64)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewFloat(frame, parsed), nil
}

// parseInt parses a string into an integer based on the supplied base and bit size; returns parsed value or error.
func (s *Strconv) parseInt(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args)
	if err != nil {
		return 0, nil, err
	}
	parsed, err := strconv.ParseInt(s1, int(i2), int(i3))
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewInt(frame, parsed), nil
}

// quote converts the given string argument into a double-quoted string literal and returns it as a new string object.
func (s *Strconv) quote(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := gk.NewString(frame, strconv.Quote(s1))
	return 1, v, nil
}

// atoi converts a string argument to an integer if the argument is a valid numeric string, returning the result or an error.
func (s *Strconv) atoi(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := strconv.Atoi(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewInt(frame, int64(res)), nil
}

// unquote removes surrounding quotes from a string if present and interprets escape sequences. Returns the unquoted string.
func (s *Strconv) unquote(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := strconv.Unquote(s1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewString(frame, res), nil
}

// itoa converts an integer to its decimal string representation and returns it as an IObject.
func (s *Strconv) itoa(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, strconv.Itoa(int(i1))), nil
}

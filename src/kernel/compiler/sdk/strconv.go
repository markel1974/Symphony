package sdk

import (
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Strconv is a type that provides a module containing string conversion functions implemented using strconv.
type Strconv struct {
	module map[string]objects.IObject
}

// NewStrconv initializes and returns a pointer to a new Strconv instance containing predefined module functions.
func NewStrconv() *Strconv {
	s := &Strconv{}
	s.module = map[string]objects.IObject{
		"Atoi":        objects.NewFunctionModule(objects.FunctionModuleDef, "Atoi", objects.FuncIsOie(strconv.Atoi)),
		"FormatBool":  objects.NewFunctionModule(objects.FunctionModuleDef, "FormatBool", s.FormatBool),
		"FormatFloat": objects.NewFunctionModule(objects.FunctionModuleDef, "FormatFloat", s.FormatFloat),
		"FormatInt":   objects.NewFunctionModule(objects.FunctionModuleDef, "FormatInt", s.FormatInt),
		"Itoa":        objects.NewFunctionModule(objects.FunctionModuleDef, "Itoa", objects.FuncIiOs(strconv.Itoa)),
		"ParseBool":   objects.NewFunctionModule(objects.FunctionModuleDef, "ParseBool", s.ParseBool),
		"ParseFloat":  objects.NewFunctionModule(objects.FunctionModuleDef, "ParseFloat", s.ParseFloat),
		"ParseNumber": objects.NewFunctionModule(objects.FunctionModuleDef, "ParseNumber", s.ParseNumber),
		"ParseInt":    objects.NewFunctionModule(objects.FunctionModuleDef, "ParseInt", s.ParseInt),
		"Quote":       objects.NewFunctionModule(objects.FunctionModuleDef, "Quote", objects.FuncIsOs(strconv.Quote)),
		"Unquote":     objects.NewFunctionModule(objects.FunctionModuleDef, "Unquote", objects.FuncIsOse(strconv.Unquote)),
	}
	return s
}

// Module retrieves the `module` field from the `Strconv` struct and returns it as a map of string keys and IObject values.
func (s *Strconv) Module() map[string]objects.IObject {
	return s.module
}

// FormatBool converts a boolean argument to its string representation ("true" or "false"). Returns an error if the argument is invalid.
func (s *Strconv) FormatBool(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	b1, err := objects.ToBoolArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if b1 {
		return objects.NewString("true")
	}
	return objects.NewString("false")
}

// FormatFloat converts a float64 into a string representation according to the specified format, precision, and bit size.
func (s *Strconv) FormatFloat(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	f1, err := objects.ToFloat64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strconv.FormatFloat(f1, s2[0], int(i3), int(i4)))
}

// FormatInt formats an int64 number as a string in the specified base, provided by the second argument.
func (s *Strconv) FormatInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 2 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	return objects.NewString(strconv.FormatInt(i1, int(i2)))
}

// ParseBool parses a string representation of a boolean value and returns the corresponding boolean object or an error.
func (s *Strconv) ParseBool(args ...objects.IObject) (ret objects.IObject, err error) {
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
		ret = objects.NewObjectError(err)
		return
	}
	if parsed {
		ret = objects.TrueValue
	} else {
		ret = objects.FalseValue
	}
	return
}

// ParseFloat parses a string argument as a floating-point number using the specified precision and returns the result.
func (s *Strconv) ParseFloat(args ...objects.IObject) (objects.IObject, error) {
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
	parsed, err := strconv.ParseFloat(s1, int(i2))
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewFloat(parsed), nil
}

// ParseNumber extracts and parses numeric values from a string, returning a float object or an error if parsing fails.
func (s *Strconv) ParseNumber(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
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
	return objects.NewFloat(parsed), nil
}

// ParseInt converts a string argument to an integer with the specified base and bit size after validating arguments.
func (s *Strconv) ParseInt(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
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
	i3, err := objects.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	parsed, err := strconv.ParseInt(s1, int(i2), int(i3))
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(parsed), nil
}

package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Fmt represents a struct that provides formatted output methods using a map of predefined functions.
type Fmt struct {
	module map[string]objects.IObject
}

// NewFmt initializes and returns a new Fmt instance with predefined formatting functions as module properties.
func NewFmt() *Fmt {
	f := &Fmt{}
	f.module = map[string]objects.IObject{
		"Print":   objects.NewFunctionModule(objects.FunctionModuleDef, "Print", f.Print),
		"Printf":  objects.NewFunctionModule(objects.FunctionModuleDef, "Printf", f.Printf),
		"Println": objects.NewFunctionModule(objects.FunctionModuleDef, "Println", f.Println),
		"Sprint":  objects.NewFunctionModule(objects.FunctionModuleDef, "Sprintf", f.Sprint),
		"Sprintf": objects.NewFunctionModule(objects.FunctionModuleDef, "Sprintf", f.Sprintf),
		"Errorf":  objects.NewFunctionModule(objects.FunctionModuleDef, "Errorf", f.Errorf),
	}
	return f
}

// Name returns the name of Fmt module.
func (f *Fmt) Name() string {
	return "fmt"
}

// Module returns the module map containing string keys and corresponding IObject values from the Fmt struct.
func (f *Fmt) Module() map[string]objects.IObject {
	return f.module
}

// Print writes the string representations of the provided arguments to the standard output without appending a newline.
func (f *Fmt) Print(args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, objects.ToInterface(arg))
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

// Printf formats and outputs a string using the provided format and arguments, implementing similar behavior to fmt.Printf.
// The first argument must be a format string, with additional arguments used to populate the format specifiers.
// Returns an error if the number of arguments is insufficient or if the format argument is incompatible.
func (f *Fmt) Printf(args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if argsLen == 1 {
		fmt.Print(s1)
		return nil, nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, objects.ToInterface(v))
	}
	fmt.Printf(s1, ar...)
	return nil, nil
}

// Println writes the given arguments to the standard output with a newline and returns nil and no error.
func (f *Fmt) Println(args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, objects.ToInterface(arg))
	}
	_, _ = fmt.Println(printArgs...)
	return nil, nil
}

// Sprint formats and concatenates the provided arguments into a single string and returns it as a new String object.
func (f *Fmt) Sprint(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	var ar []interface{}
	for _, v := range args {
		ar = append(ar, objects.ToInterface(v))
	}
	return objects.NewString(fmt.Sprint(ar))
}

// Sprintf formats a string using a format specifier and optional arguments, returning it as a new string object.
func (f *Fmt) Sprintf(args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return objects.NewString(s1)
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, objects.ToInterface(v))
	}
	return objects.NewString(fmt.Sprintf(s1, ar...))
}

// Errorf formats an error message using a format string and arguments, returning an IObject error representation.
func (f *Fmt) Errorf(args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		v, err := objects.NewString(fmt.Errorf(s1).Error())
		if err != nil {
			return nil, err
		}
		return objects.NewError(v), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, objects.ToInterface(v))
	}
	v, err := objects.NewString(fmt.Errorf(s1, ar...).Error())
	if err != nil {
		return nil, err
	}
	return objects.NewError(v), nil
}

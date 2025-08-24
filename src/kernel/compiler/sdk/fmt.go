package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Fmt represents a struct that provides formatted output methods using a map of predefined functions.
type Fmt struct {
	factory objects.IGateKeeper
	*Package
}

// NewFmt initializes and returns a new Fmt instance with predefined formatting functions as module properties.
func NewFmt(factory objects.IGateKeeper) *Fmt {
	f := &Fmt{
		factory: factory,
	}
	container := []objects.IObject{
		factory.NewFuncPackage(objects.FuncPackageDef, "Print", f.print),
		factory.NewFuncPackage(objects.FuncPackageDef, "Printf", f.printf),
		factory.NewFuncPackage(objects.FuncPackageDef, "Println", f.println),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sprintf", f.sprint),
		factory.NewFuncPackage(objects.FuncPackageDef, "Sprintf", f.sprintf),
		factory.NewFuncPackage(objects.FuncPackageDef, "Errorf", f.errorf),
	}
	f.Package = NewPackage("fmt", container, nil)
	return f
}

// Print writes the string representations of the provided arguments to the standard output without appending a newline.
func (f *Fmt) print(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, f.factory.ToInterface(arg))
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

// Printf formats and outputs a string using the provided format and arguments, implementing similar behavior to fmt.Printf.
// The first argument must be a format string, with additional arguments used to populate the format specifiers.
// Returns an error if the number of arguments is insufficient or if the format argument is incompatible.
func (f *Fmt) printf(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := f.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if argsLen == 1 {
		fmt.Print(s1)
		return nil, nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, f.factory.ToInterface(v))
	}
	fmt.Printf(s1, ar...)
	return nil, nil
}

// Println writes the given arguments to the standard output with a newline and returns nil and no error.
func (f *Fmt) println(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, f.factory.ToInterface(arg))
	}
	_, _ = fmt.Println(printArgs...)
	return nil, nil
}

// Sprint formats and concatenates the provided arguments into a single string and returns it as a new String object.
func (f *Fmt) sprint(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	var ar []interface{}
	for _, v := range args {
		ar = append(ar, f.factory.ToInterface(v))
	}
	return f.factory.NewString(frame, fmt.Sprint(ar)), nil
}

// Sprintf formats a string using a format specifier and optional arguments, returning it as a new string object.
func (f *Fmt) sprintf(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := f.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return f.factory.NewString(frame, s1), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, f.factory.ToInterface(v))
	}
	return f.factory.NewString(frame, fmt.Sprintf(s1, ar...)), nil
}

// Errorf formats an error message using a format string and arguments, returning an IObject error representation.
func (f *Fmt) errorf(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := f.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return f.factory.NewError(frame, s1), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, f.factory.ToInterface(v))
	}
	return f.factory.NewError(frame, s1), nil
}

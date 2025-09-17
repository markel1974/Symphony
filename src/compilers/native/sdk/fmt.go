package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewFmt)
}

// Fmt represents a struct that provides formatted output methods using a map of predefined functions.
type Fmt struct {
	*Package
}

// NewFmt initializes and returns a new Fmt instance with predefined formatting functions as module properties.
func NewFmt(factory objects.IGateKeeper) IPackage {
	f := &Fmt{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Print", -1, f.print),
		factory.NewFuncImport(objects.FrameStatic, "Printf", -1, f.printf),
		factory.NewFuncImport(objects.FrameStatic, "Println", -1, f.println),
		factory.NewFuncImport(objects.FrameStatic, "Sprintf", -1, f.sprint),
		factory.NewFuncImport(objects.FrameStatic, "Sprintf", -1, f.sprintf),
		factory.NewFuncImport(objects.FrameStatic, "Errorf", -1, f.errorf),
	}
	f.Package = NewExternalPackage("fmt", container, nil)
	return f
}

// Print writes the string representations of the provided arguments to the standard output without appending a newline.
func (f *Fmt) print(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, arg.AsInterface())
	}
	_, _ = fmt.Print(printArgs...)
	return 0, nil, nil
}

// Printf formats and outputs a string using the provided format and arguments, implementing similar behavior to fmt.Printf.
// The first argument must be a format string, with additional arguments used to populate the format specifiers.
// Returns an error if the number of arguments is insufficient or if the format argument is incompatible.
func (f *Fmt) printf(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	if argsLen == 1 {
		fmt.Print(s1)
		return 0, nil, nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, v.AsInterface())
	}
	fmt.Printf(s1, ar...)
	return 0, nil, nil
}

// Println writes the given arguments to the standard output with a newline and returns nil and no error.
func (f *Fmt) println(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, arg.AsInterface())
	}
	_, _ = fmt.Println(printArgs...)
	return 0, nil, nil
}

// Sprint formats and concatenates the provided arguments into a single string and returns it as a new AsString object.
func (f *Fmt) sprint(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	if len(args) == 0 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	var ar []interface{}
	for _, v := range args {
		ar = append(ar, v.AsInterface())
	}
	return 1, gk.NewString(frame, fmt.Sprint(ar)), nil
}

// Sprintf formats a string using a format specifier and optional arguments, returning it as a new string object.
func (f *Fmt) sprintf(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	if len(args) == 1 {
		return 1, gk.NewString(frame, s1), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, v.AsInterface())
	}
	return 1, gk.NewString(frame, fmt.Sprintf(s1, ar...)), nil
}

// Errorf formats an error message using a format string and arguments, returning an IObject error representation.
func (f *Fmt) errorf(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	if len(args) == 1 {
		return 0, gk.NewError(frame, s1), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, v.AsInterface())
	}
	return 1, gk.NewError(frame, s1), nil
}

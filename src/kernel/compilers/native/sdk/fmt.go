package sdk

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewFmt)
}

// Fmt represents a struct that provides formatted output methods using a map of predefined functions.
type Fmt struct {
	container map[string]objects.IObject
}

// NewFmt initializes and returns a new Fmt instance with predefined formatting functions as module properties.
func NewFmt(factory objects.IGateKeeper) IPackage {
	f := &Fmt{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Print", f.print),
		factory.NewFuncImport(objects.FrameStatic, "Printf", f.printf),
		factory.NewFuncImport(objects.FrameStatic, "Println", f.println),
		factory.NewFuncImport(objects.FrameStatic, "Sprintf", f.sprint),
		factory.NewFuncImport(objects.FrameStatic, "Sprintf", f.sprintf),
		factory.NewFuncImport(objects.FrameStatic, "Errorf", f.errorf),
	}
	f.container = BuildContainer(container, nil)
	return f
}

// Name returns the name of the Fmt struct as a string.
func (f *Fmt) Name() string {
	return "fmt"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (f *Fmt) Get(name string) (objects.IObject, bool) {
	v, ok := f.container[name]
	return v, ok
}

// Print writes the string representations of the provided arguments to the standard output without appending a newline.
func (f *Fmt) print(gk objects.IGateKeeper, _ int, args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, gk.ToInterface(arg))
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

// Printf formats and outputs a string using the provided format and arguments, implementing similar behavior to fmt.Printf.
// The first argument must be a format string, with additional arguments used to populate the format specifiers.
// Returns an error if the number of arguments is insufficient or if the format argument is incompatible.
func (f *Fmt) printf(gk objects.IGateKeeper, _ int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if argsLen == 1 {
		fmt.Print(s1)
		return nil, nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, gk.ToInterface(v))
	}
	fmt.Printf(s1, ar...)
	return nil, nil
}

// Println writes the given arguments to the standard output with a newline and returns nil and no error.
func (f *Fmt) println(gk objects.IGateKeeper, _ int, args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, gk.ToInterface(arg))
	}
	_, _ = fmt.Println(printArgs...)
	return nil, nil
}

// Sprint formats and concatenates the provided arguments into a single string and returns it as a new AsString object.
func (f *Fmt) sprint(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) == 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	var ar []interface{}
	for _, v := range args {
		ar = append(ar, gk.ToInterface(v))
	}
	return gk.NewString(frame, fmt.Sprint(ar)), nil
}

// Sprintf formats a string using a format specifier and optional arguments, returning it as a new string object.
func (f *Fmt) sprintf(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return gk.NewString(frame, s1), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, gk.ToInterface(v))
	}
	return gk.NewString(frame, fmt.Sprintf(s1, ar...)), nil
}

// Errorf formats an error message using a format string and arguments, returning an IObject error representation.
func (f *Fmt) errorf(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return gk.NewError(frame, s1), nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, gk.ToInterface(v))
	}
	return gk.NewError(frame, s1), nil
}

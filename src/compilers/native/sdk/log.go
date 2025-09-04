package sdk

import (
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	RegisterPackage(NewLog)
}

// Log represents a struct that provides formatted output methods using a map of predefined functions.
type Log struct {
	container map[string]objects.IObject
}

// NewLog initializes and returns a new Log instance with predefined formatting functions as module properties.
func NewLog(factory objects.IGateKeeper) IPackage {
	f := &Log{}
	container := []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Print", f.print),
		factory.NewFuncImport(objects.FrameStatic, "Printf", f.printf),
		factory.NewFuncImport(objects.FrameStatic, "Println", f.println),
		factory.NewFuncImport(objects.FrameStatic, "Fatalf", f.fatalf),
	}
	f.container = BuildContainer(container, nil)
	return f
}

// Name returns the name of the Fmt struct as a string.
func (f *Log) Name() string {
	return "log"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (f *Log) Get(name string) (objects.IObject, bool) {
	v, ok := f.container[name]
	return v, ok
}

// Print writes the string representations of the provided arguments to the standard output without appending a newline.
func (f *Log) print(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, gk.ToInterface(arg))
	}
	log.Print(printArgs...)
	return 0, nil, nil
}

// Printf formats and outputs a string using the provided format and arguments, implementing similar behavior to fmt.Printf.
// The first argument must be a format string, with additional arguments used to populate the format specifiers.
// Returns an error if the number of arguments is insufficient or if the format argument is incompatible.
func (f *Log) printf(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	if argsLen == 1 {
		log.Printf(s1)
		return 0, nil, nil
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, gk.ToInterface(v))
	}
	log.Printf(s1, ar...)
	return 0, nil, nil
}

// Println writes the given arguments to the standard output with a newline and returns nil and no error.
func (f *Log) println(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, gk.ToInterface(arg))
	}
	log.Println(printArgs...)
	return 0, nil, nil
}

// fatalf formats and logs a fatal error message using the provided arguments and terminates execution with an error.
func (f *Log) fatalf(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	argsLen := len(args)
	if argsLen == 0 {
		return 0, nil, objects.ErrInvalidArgumentsNumber
	}
	s1, err := gk.ToStringArg(0, args[0])
	if err != nil {
		return 0, nil, err
	}
	if argsLen == 1 {
		log.Printf(s1)
		return 0, nil, fmt.Errorf("fatal")
	}
	var ar []interface{}
	for _, v := range args[1:] {
		ar = append(ar, gk.ToInterface(v))
	}
	log.Printf(s1, ar...)
	return 0, nil, fmt.Errorf("fatal")
}

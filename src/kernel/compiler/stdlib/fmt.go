package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// fmtSafeModule is a map containing predefined functions from the fmt package wrapped as IObjects, like "sprintf".
//var _fmtSafeModule = map[string]objects.IObject{
//	"Sprintf": objects.NewFunctionModule(objects.FunctionModuleDef, "Sprintf", fmtSprintf),
//}

// fmtModule is a map that associates string keys with user-defined function objects for various formatted output operations.
var _fmtModule = map[string]objects.IObject{
	"Print":   objects.NewFunctionModule(objects.FunctionModuleDef, "Print", fmtPrint),
	"Printf":  objects.NewFunctionModule(objects.FunctionModuleDef, "Printf", fmtPrintf),
	"Println": objects.NewFunctionModule(objects.FunctionModuleDef, "Println", fmtPrintln),
	"Sprint":  objects.NewFunctionModule(objects.FunctionModuleDef, "Sprintf", fmtSprint),
	"Sprintf": objects.NewFunctionModule(objects.FunctionModuleDef, "Sprintf", fmtSprintf),
	"Errorf":  objects.NewFunctionModule(objects.FunctionModuleDef, "Errorf", fmtErrorf),
}

// fmtPrint prints the string representations of the provided IObject arguments without a newline.
func fmtPrint(args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, objects.ToInterface(arg))
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

// fmtPrintf formats and prints a string using given arguments; returns nil or an error on invalid input or formatting.
func fmtPrintf(args ...objects.IObject) (ret objects.IObject, err error) {
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

// fmtPrintln outputs arguments to the standard output with a newline, converting them to strings if necessary.
// Returns nil and an error if any argument conversion fails.
func fmtPrintln(args ...objects.IObject) (ret objects.IObject, err error) {
	var printArgs []interface{}
	for _, arg := range args {
		printArgs = append(printArgs, objects.ToInterface(arg))
	}
	_, _ = fmt.Println(printArgs...)
	return nil, nil
}

// fmtSprintf formats a string based on a format string and arguments, returning the result or an error if formatting fails.
func fmtSprint(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	var ar []interface{}
	for _, v := range args {
		ar = append(ar, objects.ToInterface(v))
	}
	return objects.NewString(fmt.Sprint(ar))
}

// fmtSprintf formats a string based on a format string and arguments, returning the result or an error if formatting fails.
func fmtSprintf(args ...objects.IObject) (ret objects.IObject, err error) {
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

// fmtSprintf formats a string based on a format string and arguments, returning the result or an error if formatting fails.
func fmtErrorf(args ...objects.IObject) (ret objects.IObject, err error) {
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

package stdlib

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/compiler/modules/format"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// fmtSafeModule is a map containing predefined functions from the fmt package wrapped as IObjects, like "sprintf".
var fmtSafeModule = map[string]objects.IObject{
	"sprintf": objects.NewFunctionUser("sprintf", fmtSprintf),
}

// fmtModule is a map that associates string keys with user-defined function objects for various formatted output operations.
var fmtModule = map[string]objects.IObject{
	"print":   objects.NewFunctionUser("print", fmtPrint),
	"printf":  objects.NewFunctionUser("printf", fmtPrintf),
	"println": objects.NewFunctionUser("println", fmtPrintln),
	"sprintf": objects.NewFunctionUser("sprintf", fmtSprintf),
}

// fmtPrint prints the string representations of the provided IObject arguments without a newline.
func fmtPrint(args ...objects.IObject) (ret objects.IObject, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

// fmtPrintf formats and prints a string using given arguments; returns nil or an error on invalid input or formatting.
func fmtPrintf(args ...objects.IObject) (ret objects.IObject, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	data, ok := args[0].(*objects.String)
	if !ok {
		return nil, objects.NewInvalidArgumentError("format", "string", args[0].TypeName())
	}
	if numArgs == 1 {
		fmt.Print(data)
		return nil, nil
	}

	s, err := format.Format(data.Value(), args[1:]...)
	if err != nil {
		return nil, err
	}
	fmt.Print(s)
	return nil, nil
}

// fmtPrintln outputs arguments to the standard output with a newline, converting them to strings if necessary.
// Returns nil and an error if any argument conversion fails.
func fmtPrintln(args ...objects.IObject) (ret objects.IObject, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	printArgs = append(printArgs, "\n")
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

// fmtSprintf formats a string based on a format string and arguments, returning the result or an error if formatting fails.
func fmtSprintf(args ...objects.IObject) (ret objects.IObject, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	data, ok := args[0].(*objects.String)
	if !ok {
		return nil, objects.NewInvalidArgumentError("format", "string", args[0].TypeName())
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return data, nil
	}
	s, err := format.Format(data.Value(), args[1:]...)
	if err != nil {
		return nil, err
	}
	return objects.NewString(s)
}

// getPrintArgs converts IObject arguments to their string representations while ensuring the total length does not exceed the limit.
// Returns a slice of string representations as interface{} or an error if the length exceeds MaxStringLen.
func getPrintArgs(args ...objects.IObject) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := objects.ToString(arg)
		sLen := len(s)
		// make sure length does not exceed the limit
		if l+sLen > objects.MaxStringLen {
			return nil, objects.ErrStringLimit
		}
		l += sLen
		printArgs = append(printArgs, s)
	}
	return printArgs, nil
}

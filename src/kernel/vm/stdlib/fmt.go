package stdlib

import (
	"fmt"
	//"github.com/markel1974/injector/src/vm/compiler"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/modules/format"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

var fmtSafeModule = map[string]objects.IObject{
	"sprintf": objects.NewUserFunction("sprintf", fmtSprintf),
}

var fmtModule = map[string]objects.IObject{
	"print":   objects.NewUserFunction("print", fmtPrint),
	"printf":  objects.NewUserFunction("printf", fmtPrintf),
	"println": objects.NewUserFunction("println", fmtPrintln),
	"sprintf": objects.NewUserFunction("sprintf", fmtSprintf),
}

func fmtPrint(args ...objects.IObject) (ret objects.IObject, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

func fmtPrintf(args ...objects.IObject) (ret objects.IObject, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, errors.ErrWrongNumArguments
	}
	data, ok := args[0].(*objects.String)
	if !ok {
		return nil, errors.NewInvalidArgumentType("format", "string", args[0].TypeName())
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

func fmtPrintln(args ...objects.IObject) (ret objects.IObject, err error) {
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	printArgs = append(printArgs, "\n")
	_, _ = fmt.Print(printArgs...)
	return nil, nil
}

func fmtSprintf(args ...objects.IObject) (ret objects.IObject, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, errors.ErrWrongNumArguments
	}
	data, ok := args[0].(*objects.String)
	if !ok {
		return nil, errors.NewInvalidArgumentType("format", "string", args[0].TypeName())
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

func getPrintArgs(args ...objects.IObject) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := objects.ToString(arg)
		sLen := len(s)
		// make sure length does not exceed the limit
		if l+sLen > objects.MaxStringLen {
			return nil, errors.ErrStringLimit
		}
		l += sLen
		printArgs = append(printArgs, s)
	}
	return printArgs, nil
}

package process

import (
	"errors"
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Library represents a collection of modules that interact with system processes and provide various functionalities.
// It contains a reference to a Process and a map of module names to their respective objects.
type Library struct {
	factory *objects.GateKeeper
	process *Process
	pkg     map[string]objects.IObject
}

// NewLibrary creates and initializes a new Library instance with the provided Process.
func NewLibrary(factory *objects.GateKeeper, process *Process) *Library {
	l := &Library{
		factory: factory,
		process: process,
		pkg:     make(map[string]objects.IObject),
	}
	container := []objects.IObject{
		factory.NewFuncPackage(objects.FuncPackageDef, "Printf", l.doPrintf),
		factory.NewFuncPackage(objects.FuncPackageDef, "CreateTimer", l.doCreateTimer),
		factory.NewFuncPackage(objects.FuncPackageDef, "IsActive", l.doIsActive),
		factory.NewFuncPackage(objects.FuncPackageDef, "Kill", l.doKill),
		factory.NewFuncPackage(objects.FuncPackageDef, "KillForeground", l.doKillForeground),
		factory.NewFuncPackage(objects.FuncPackageDef, "KillAll", l.doKillAll),
		factory.NewFuncPackage(objects.FuncPackageDef, "CWDSet", l.doCWDSet),
		factory.NewFuncPackage(objects.FuncPackageDef, "CWDName", l.doCWDName),
		factory.NewFuncPackage(objects.FuncPackageDef, "CWDPath", l.doCWDPath),
		factory.NewFuncPackage(objects.FuncPackageDef, "CWDDirectoryListing", l.doCWDDirectoryListing),
		factory.NewFuncPackage(objects.FuncPackageDef, "GetScreenSize", l.doGetScreenSize),
		factory.NewFuncPackage(objects.FuncPackageDef, "PaintRequest", l.doPaintRequest),
		factory.NewFuncPackage(objects.FuncPackageDef, "ProcessExec", l.doProcessExec),
		factory.NewFuncPackage(objects.FuncPackageDef, "WindowsSelectionBegin", l.doWindowsSelectionBegin),
		factory.NewFuncPackage(objects.FuncPackageDef, "CWDSet", l.doWindowsSelectionEnd),
		factory.NewFuncPackage(objects.FuncPackageDef, "WindowsSelectionOptions", l.doWindowsSelectionOptions),
		factory.NewFuncPackage(objects.FuncPackageDef, "WindowsSelectionNext", l.doWindowsSelectionNext),
		factory.NewFuncPackage(objects.FuncPackageDef, "WindowsSelectionPrevious", l.doWindowsSelectionPrevious),
		factory.NewFuncPackage(objects.FuncPackageDef, "ProcessList", l.doProcessList),
		factory.NewFuncPackage(objects.FuncPackageDef, "ProcessSetForeground", l.doProcessSetForeground),
		factory.NewFuncPackage(objects.FuncPackageDef, "ProcessSetSelfForeground", l.doProcessSetSelfForeground),
		factory.NewFuncPackage(objects.FuncPackageDef, "Write", l.doWrite),
		factory.NewFuncPackage(objects.FuncPackageDef, "WritePromptEOL", l.doWritePromptEOL),
		factory.NewFuncPackage(objects.FuncPackageDef, "WritePromptLine", l.doWritePromptLine),
		factory.NewFuncPackage(objects.FuncPackageDef, "WriteColor", l.doWriteColor),
		factory.NewFuncPackage(objects.FuncPackageDef, "WriteForeground", l.doWriteForeground),
		factory.NewFuncPackage(objects.FuncPackageDef, "MoveCursorLeft", l.doMoveCursorLeft),
		factory.NewFuncPackage(objects.FuncPackageDef, "MoveCursorRight", l.doMoveCursorRight),
		factory.NewFuncPackage(objects.FuncPackageDef, "SaveCursor", l.doSaveCursor),
		factory.NewFuncPackage(objects.FuncPackageDef, "RestoreCursor", l.doRestoreCursor),
		factory.NewFuncPackage(objects.FuncPackageDef, "ClearScreen", l.doClearScreen),
		factory.NewFuncPackage(objects.FuncPackageDef, "SetExit", l.doSetExit),
		factory.NewFuncPackage(objects.FuncPackageDef, "Suggestion", l.doSuggestion),
		factory.NewFuncPackage(objects.FuncPackageDef, "Help", l.doHelp),
	}
	for _, obj := range container {
		if fn, ok := obj.(*objects.FuncPackage); ok {
			l.pkg[fn.Name()] = fn
		}
	}
	return l
}

// Package returns a map where keys are strings and values implement the IObject interface, representing the library's package.
func (l *Library) Package() map[string]objects.IObject {
	return l.pkg
}

// doCreateTimer validates and extracts three integer arguments, then creates a timer using these arguments. Returns nil or an error.
func (l *Library) doCreateTimer(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := l.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := l.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := l.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	l.process.CreateTimer(int(i1), int(i2), int(i3))
	return nil, nil
}

// doPrintf formats and writes a string output using provided arguments; returns an error for invalid input or formatting.
func (l *Library) doPrintf(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	argsLen := len(args)
	if argsLen == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var val string
	if len(args) == 1 {
		val = s1
	} else {
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, l.factory.ToInterface(v))
		}
		val = fmt.Sprintf(s1, ar...)
	}
	l.process.Write(val, true)
	return nil, nil
}

func (l *Library) doIsActive(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := l.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	v := l.process.IsActive(int(i1))
	return l.factory.FromBool(v), nil
}

// doKill terminates a process identified by the provided argument, which must be convertible to an int64. Returns no value.
func (l *Library) doKill(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := l.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.process.Kill(int(i1))
	return nil, nil
}

// doKillForeground terminates the foreground process if no arguments are provided, returning nil or an error.
func (l *Library) doKillForeground(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.KillForeground()
	return nil, nil
}

// doKillAll terminates all processes matching the given context name, passed as a single string argument.
// Returns an error if the argument is invalid or conversion to string fails.
func (l *Library) doKillAll(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.process.KillAll(s1)
	return nil, nil
}

// doWindowsSelectionEnd finalizes a windows selection operation, ensuring no arguments are passed and invoking WindowsSelectionEnd.
func (l *Library) doWindowsSelectionEnd(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.WindowsSelectionEnd()
	return nil, nil
}

// doWindowsSelectionBegin initializes the window selection process in the current library context.
// Returns an error if invalid arguments are passed.
func (l *Library) doWindowsSelectionBegin(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.WindowsSelectionBegin()
	return nil, nil
}

// doProcessExec validates arguments and executes a process based on a string argument. Returns an error if invalid.
func (l *Library) doProcessExec(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.process.ProcessExec(s1)
	return nil, nil
}

// doPaintRequest performs a paint request by invoking the PaintRequest method of the Library's internal painter instance.
// It validates that no arguments are passed; otherwise, it returns an error.
// Returns nil if the operation succeeds or an error if the argument validation fails.
func (l *Library) doPaintRequest(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.PaintRequest()
	return nil, nil
}

// doGetScreenSize retrieves the screen dimensions as a map containing "width" and "height" with their respective values.
func (l *Library) doGetScreenSize(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	w, h := l.process.GetScreenSize()
	return l.factory.NewMap(frame, map[string]objects.IObject{
		"width":  l.factory.NewInt(frame, int64(w)),
		"height": l.factory.NewInt(frame, int64(h)),
	}), nil
}

// doCWDDirectoryListing returns the directory listing of the current working directory as an IObject array.
// It expects no arguments and raises an error if arguments are provided.
func (l *Library) doCWDDirectoryListing(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	v := l.process.CWDDirectoryListing()
	return l.factory.FromStringArray(frame, v)
}

// doCWDSet sets the current working directory (CWD) using the provided string argument. Returns a boolean as IObject.
func (l *Library) doCWDSet(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	v := l.process.CWDSet(s1)
	return l.factory.FromBool(v), nil
}

// doCWDPath returns the current working directory's path as a string object or an error if arguments are provided.
func (l *Library) doCWDPath(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	v := l.process.CWDPath()
	return l.factory.NewString(frame, v), nil
}

// doCWDName retrieves the current working directory name, returns it as a string object, and validates argument count.
func (l *Library) doCWDName(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	v := l.process.CWDName()
	return l.factory.NewString(frame, v), nil
}

// doHelp retrieves and returns a help string for the provided argument, which must be a single string-compatible object.
func (l *Library) doHelp(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	v, err := l.process.Help(s1)
	if err != nil {
		return nil, err
	}
	return l.factory.NewString(frame, v), nil
}

// doSuggestion processes two arguments: a string and an integer, to trigger the Suggestion functionality within the library.
// Returns an error for invalid argument types or count.
func (l *Library) doSuggestion(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := l.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.process.Suggestion(s1, int(i2))
	return nil, nil
}

// doSetExit sets the exit condition for the associated process and returns nil if successful or an error otherwise.
// Accepts no arguments; an error is returned if arguments are provided.
func (l *Library) doSetExit(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.process.SetExit()
	return nil, nil
}

// doClearScreen clears the terminal screen and returns nil if successful or an error if invalid arguments are passed.
func (l *Library) doClearScreen(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.ClearScreen()
	return nil, nil
}

// doRestoreCursor restores the cursor to its previously saved position, returning an error for invalid arguments.
func (l *Library) doRestoreCursor(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.RestoreCursor()
	return nil, nil
}

// doSaveCursor saves the current cursor position. Returns an error for invalid argument count.
func (l *Library) doSaveCursor(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.process.SaveCursor()
	return nil, nil
}

// doMoveCursorRight moves the cursor to the right and accepts no arguments, returning an error if any are provided.
func (l *Library) doMoveCursorRight(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.MoveCursorRight()
	return nil, nil
}

// doMoveCursorLeft moves the cursor one position to the left if the number of arguments is zero, otherwise returns an error.
func (l *Library) doMoveCursorLeft(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.MoveCursorLeft()
	return nil, nil
}

// doWriteForeground writes text to the foreground using specified string, color, and a boolean flag for emphasis.
// Accepts three arguments: a string, an int64 representing color, and a boolean indicating emphasis.
// Returns an error if the argument count is invalid or conversion of arguments fails.
func (l *Library) doWriteForeground(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := l.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	b3, err := l.factory.ToBoolArg(2, args[2])
	if err != nil {
		return nil, err
	}
	l.process.WriteForeground(s1, interfaces.ColorDef(i2), b3)
	return nil, nil
}

// doWriteColor writes colored text to the output using specified color definitions, mode, and a boolean flag.
func (l *Library) doWriteColor(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 5 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := l.factory.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	i3, err := l.factory.ToInt64Arg(2, args[2])
	if err != nil {
		return nil, err
	}
	i4, err := l.factory.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	b5, err := l.factory.ToBoolArg(4, args[4])
	if err != nil {
		return nil, err
	}
	l.process.WriteColor(s1, interfaces.ColorDef(i2), interfaces.ColorDef(i3), interfaces.ColorMode(i4), b5)
	return nil, nil
}

// doWritePromptLine writes a prompt line using two string arguments from objects.IObject and returns nil or an error.
func (l *Library) doWritePromptLine(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := l.factory.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.process.WritePromptLine(s1, s2)
	return nil, nil
}

// doWritePromptEOL writes a prompt with an end-of-line flag based on the provided string and boolean arguments.
func (l *Library) doWritePromptEOL(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	b2, err := l.factory.ToBoolArg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.process.WritePromptEOL(s1, b2)
	return nil, nil
}

// doWrite writes a string with optional newline control using provided arguments; validates input and returns errors if invalid.
func (l *Library) doWrite(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := l.factory.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	b2, err := l.factory.ToBoolArg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.process.Write(s1, b2)
	return nil, nil
}

// doProcessSetSelfForeground sets the current process as the foreground and ensures no arguments are passed during invocation.
func (l *Library) doProcessSetSelfForeground(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.ProcessSetSelfForeground()
	return nil, nil
}

// doProcessSetForeground sets the specified process as the foreground process using the provided integer argument.
func (l *Library) doProcessSetForeground(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := l.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.process.ProcessSetForeground(int(i1))
	return nil, nil
}

// doProcessList retrieves a map of running processes, including details like name, PID, line, and start time.
// Returns an error if any arguments are provided.
func (l *Library) doProcessList(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	res := make(map[string]objects.IObject)
	for _, p := range l.process.ProcessList() {
		c := map[string]objects.IObject{
			"line": l.factory.NewString(frame, p.Line()),
			"name": l.factory.NewString(frame, p.Name()),
			"time": l.factory.NewTime(frame, p.Time()),
			"pid":  l.factory.NewInt(frame, int64(p.PID())),
		}
		res[p.Name()] = l.factory.NewMap(frame, c)
	}
	return l.factory.NewMap(frame, res), nil
}

// doWindowsSelectionPrevious navigates to the previous selection in the Windows selection context.
// Returns an error if any arguments are provided.
// Calls the underlying WindowsSelectionPrevious method to perform the operation.
func (l *Library) doWindowsSelectionPrevious(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.WindowsSelectionPrevious()
	return nil, nil
}

// doWindowsSelectionNext advances to the next item in the windows selection process. Returns an error for invalid arguments.
func (l *Library) doWindowsSelectionNext(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.process.WindowsSelectionNext()
	return nil, nil
}

// doWindowsSelectionOptions configures windows selection behavior using provided arguments as input parameters.
// Expects exactly two arguments: a rune (int64) and a float64.
// Returns an error if argument count is invalid or conversion of arguments fails.
// Invokes WindowsSelectionOptions method with the parsed inputs.
func (l *Library) doWindowsSelectionOptions(_ int, args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := l.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	f2, err := l.factory.ToFloat64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.process.WindowsSelectionOptions(rune(i1), f2)
	return nil, nil
}

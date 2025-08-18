package process

import (
	"errors"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Library represents a collection of modules that interact with system processes and provide various functionalities.
// It contains a reference to a Process and a map of module names to their respective objects.
type Library struct {
	p      *Process
	module map[string]objects.IObject
}

// NewLibrary creates and initializes a new Library instance with the provided Process.
func NewLibrary(p *Process) *Library {
	l := &Library{
		p: p,
	}
	l.setup()
	return l
}

// setup initializes the Library by mapping function names to their respective FunctionModule implementations.
func (l *Library) setup() {
	l.module = map[string]objects.IObject{
		"SetContext":               objects.NewFunctionModule(objects.FunctionModuleDef, "SetContext", l.doSetContext),
		"GetContext":               objects.NewFunctionModule(objects.FunctionModuleDef, "GetContext", l.doGetContext),
		"CreateTimer":              objects.NewFunctionModule(objects.FunctionModuleDef, "CreateTimer", l.doCreateTimer),
		"IsActive":                 objects.NewFunctionModule(objects.FunctionModuleDef, "IsActive", l.doIsActive),
		"Kill":                     objects.NewFunctionModule(objects.FunctionModuleDef, "Kill", l.doKill),
		"KillForeground":           objects.NewFunctionModule(objects.FunctionModuleDef, "KillForeground", l.doKillForeground),
		"KillAll":                  objects.NewFunctionModule(objects.FunctionModuleDef, "KillAll", l.doKillAll),
		"CWDSet":                   objects.NewFunctionModule(objects.FunctionModuleDef, "CWDSet", l.doCWDSet),
		"CWDName":                  objects.NewFunctionModule(objects.FunctionModuleDef, "CWDName", l.doCWDName),
		"CWDPath":                  objects.NewFunctionModule(objects.FunctionModuleDef, "CWDPath", l.doCWDPath),
		"CWDDirectoryListing":      objects.NewFunctionModule(objects.FunctionModuleDef, "CWDDirectoryListing", l.doCWDDirectoryListing),
		"GetScreenSize":            objects.NewFunctionModule(objects.FunctionModuleDef, "GetScreenSize", l.doGetScreenSize),
		"PaintRequest":             objects.NewFunctionModule(objects.FunctionModuleDef, "PaintRequest", l.doPaintRequest),
		"ProcessExec":              objects.NewFunctionModule(objects.FunctionModuleDef, "ProcessExec", l.doProcessExec),
		"WindowsSelectionBegin":    objects.NewFunctionModule(objects.FunctionModuleDef, "WindowsSelectionBegin", l.doWindowsSelectionBegin),
		"WindowsSelectionEnd":      objects.NewFunctionModule(objects.FunctionModuleDef, "CWDSet", l.doWindowsSelectionEnd),
		"WindowsSelectionOptions":  objects.NewFunctionModule(objects.FunctionModuleDef, "WindowsSelectionOptions", l.doWindowsSelectionOptions),
		"WindowsSelectionNext":     objects.NewFunctionModule(objects.FunctionModuleDef, "WindowsSelectionNext", l.doWindowsSelectionNext),
		"WindowsSelectionPrevious": objects.NewFunctionModule(objects.FunctionModuleDef, "WindowsSelectionPrevious", l.doWindowsSelectionPrevious),
		"ProcessList":              objects.NewFunctionModule(objects.FunctionModuleDef, "ProcessList", l.doProcessList),
		"ProcessSetForeground":     objects.NewFunctionModule(objects.FunctionModuleDef, "ProcessSetForeground", l.doProcessSetForeground),
		"ProcessSetSelfForeground": objects.NewFunctionModule(objects.FunctionModuleDef, "ProcessSetSelfForeground", l.doProcessSetSelfForeground),
		"Write":                    objects.NewFunctionModule(objects.FunctionModuleDef, "Write", l.doWrite),
		"WritePromptEOL":           objects.NewFunctionModule(objects.FunctionModuleDef, "WritePromptEOL", l.doWritePromptEOL),
		"WritePromptLine":          objects.NewFunctionModule(objects.FunctionModuleDef, "WritePromptLine", l.doWritePromptLine),
		"WriteColor":               objects.NewFunctionModule(objects.FunctionModuleDef, "WriteColor", l.doWriteColor),
		"WriteForeground":          objects.NewFunctionModule(objects.FunctionModuleDef, "WriteForeground", l.doWriteForeground),
		"MoveCursorLeft":           objects.NewFunctionModule(objects.FunctionModuleDef, "MoveCursorLeft", l.doMoveCursorLeft),
		"MoveCursorRight":          objects.NewFunctionModule(objects.FunctionModuleDef, "MoveCursorRight", l.doMoveCursorRight),
		"SaveCursor":               objects.NewFunctionModule(objects.FunctionModuleDef, "SaveCursor", l.doSaveCursor),
		"RestoreCursor":            objects.NewFunctionModule(objects.FunctionModuleDef, "RestoreCursor", l.doRestoreCursor),
		"ClearScreen":              objects.NewFunctionModule(objects.FunctionModuleDef, "ClearScreen", l.doClearScreen),
		"SetExit":                  objects.NewFunctionModule(objects.FunctionModuleDef, "SetExit", l.doSetExit),
		"Suggestion":               objects.NewFunctionModule(objects.FunctionModuleDef, "Suggestion", l.doSuggestion),
		"Help":                     objects.NewFunctionModule(objects.FunctionModuleDef, "Help", l.doHelp),
	}
}

// Module returns a map of strings to objects implementing the IObject interface, representing the module's content.
func (l *Library) Module() map[string]objects.IObject {
	return l.module
}

// doGetContext retrieves the current execution context within the library. Returns an error if arguments are provided.
func (l *Library) doGetContext(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	_ = l.p.GetContext()
	//TODO implement
	return nil, nil
}

// doCreateTimer validates and extracts three integer arguments, then creates a timer using these arguments. Returns nil or an error.
func (l *Library) doCreateTimer(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := objects.ToInt64Arg(0, args[0])
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
	l.p.CreateTimer(int(i1), int(i2), int(i3))
	return nil, nil
}

func (l *Library) doIsActive(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	v := l.p.IsActive(int(i1))
	return objects.FromBool(v), nil
}

// doKill terminates a process identified by the provided argument, which must be convertible to an int64. Returns no value.
func (l *Library) doKill(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.p.Kill(int(i1))
	return nil, nil
}

// doKillForeground terminates the foreground process if no arguments are provided, returning nil or an error.
func (l *Library) doKillForeground(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.KillForeground()
	return nil, nil
}

// doKillAll terminates all processes matching the given context name, passed as a single string argument.
// Returns an error if the argument is invalid or conversion to string fails.
func (l *Library) doKillAll(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.p.KillAll(s1)
	return nil, nil
}

// doWindowsSelectionEnd finalizes a windows selection operation, ensuring no arguments are passed and invoking WindowsSelectionEnd.
func (l *Library) doWindowsSelectionEnd(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.WindowsSelectionEnd()
	return nil, nil
}

// doWindowsSelectionBegin initializes the window selection process in the current library context.
// Returns an error if invalid arguments are passed.
func (l *Library) doWindowsSelectionBegin(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.WindowsSelectionBegin()
	return nil, nil
}

// doProcessExec validates arguments and executes a process based on a string argument. Returns an error if invalid.
func (l *Library) doProcessExec(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.p.ProcessExec(s1)
	return nil, nil
}

// doPaintRequest performs a paint request by invoking the PaintRequest method of the Library's internal painter instance.
// It validates that no arguments are passed; otherwise, it returns an error.
// Returns nil if the operation succeeds or an error if the argument validation fails.
func (l *Library) doPaintRequest(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.PaintRequest()
	return nil, nil
}

// doGetScreenSize retrieves the screen dimensions as a map containing "width" and "height" with their respective values.
func (l *Library) doGetScreenSize(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	w, h := l.p.GetScreenSize()
	return objects.NewMap(map[string]objects.IObject{
		"width":  objects.NewInt(int64(w)),
		"height": objects.NewInt(int64(h)),
	}), nil
}

// doCWDDirectoryListing returns the directory listing of the current working directory as an IObject array.
// It expects no arguments and raises an error if arguments are provided.
func (l *Library) doCWDDirectoryListing(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	v := l.p.CWDDirectoryListing()
	return objects.FromStringArray(v)
}

// doCWDSet sets the current working directory (CWD) using the provided string argument. Returns a boolean as IObject.
func (l *Library) doCWDSet(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	v := l.p.CWDSet(s1)
	return objects.FromBool(v), nil
}

// doCWDPath returns the current working directory's path as a string object or an error if arguments are provided.
func (l *Library) doCWDPath(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	v := l.p.CWDPath()
	return objects.NewString(v)
}

// doCWDName retrieves the current working directory name, returns it as a string object, and validates argument count.
func (l *Library) doCWDName(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	v := l.p.CWDName()
	return objects.NewString(v)
}

// doSetContext updates the library's context with the provided IObject. Validates a single argument is provided.
func (l *Library) doSetContext(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.SetContext(args[0])
	return nil, nil
}

// doHelp retrieves and returns a help string for the provided argument, which must be a single string-compatible object.
func (l *Library) doHelp(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	v, err := l.p.Help(s1)
	if err != nil {
		return nil, err
	}
	return objects.NewString(v)
}

// doSuggestion processes two arguments: a string and an integer, to trigger the Suggestion functionality within the library.
// Returns an error for invalid argument types or count.
func (l *Library) doSuggestion(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.p.Suggestion(s1, int(i2))
	return nil, nil
}

// doSetExit sets the exit condition for the associated process and returns nil if successful or an error otherwise.
// Accepts no arguments; an error is returned if arguments are provided.
func (l *Library) doSetExit(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.p.SetExit()
	return nil, nil
}

// doClearScreen clears the terminal screen and returns nil if successful or an error if invalid arguments are passed.
func (l *Library) doClearScreen(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.p.ClearScreen()
	return nil, nil
}

// doRestoreCursor restores the cursor to its previously saved position, returning an error for invalid arguments.
func (l *Library) doRestoreCursor(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.p.RestoreCursor()
	return nil, nil
}

// doSaveCursor saves the current cursor position. Returns an error for invalid argument count.
func (l *Library) doSaveCursor(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.p.SaveCursor()
	return nil, nil
}

// doMoveCursorRight moves the cursor to the right and accepts no arguments, returning an error if any are provided.
func (l *Library) doMoveCursorRight(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.p.MoveCursorRight()
	return nil, nil
}

// doMoveCursorLeft moves the cursor one position to the left if the number of arguments is zero, otherwise returns an error.
func (l *Library) doMoveCursorLeft(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}

	l.p.MoveCursorLeft()
	return nil, nil
}

// doWriteForeground writes text to the foreground using specified string, color, and a boolean flag for emphasis.
// Accepts three arguments: a string, an int64 representing color, and a boolean indicating emphasis.
// Returns an error if the argument count is invalid or conversion of arguments fails.
func (l *Library) doWriteForeground(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 3 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	i2, err := objects.ToInt64Arg(1, args[1])
	if err != nil {
		return nil, err
	}
	b3, err := objects.ToBoolArg(2, args[2])
	if err != nil {
		return nil, err
	}
	l.p.WriteForeground(s1, interfaces.ColorDef(i2), b3)
	return nil, nil
}

// doWriteColor writes colored text to the output using specified color definitions, mode, and a boolean flag.
func (l *Library) doWriteColor(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 5 {
		return nil, errors.New("invalid number of arguments")
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
	i4, err := objects.ToInt64Arg(3, args[3])
	if err != nil {
		return nil, err
	}
	b5, err := objects.ToBoolArg(4, args[4])
	if err != nil {
		return nil, err
	}
	l.p.WriteColor(s1, interfaces.ColorDef(i2), interfaces.ColorDef(i3), interfaces.ColorMode(i4), b5)
	return nil, nil
}

// doWritePromptLine writes a prompt line using two string arguments from objects.IObject and returns nil or an error.
func (l *Library) doWritePromptLine(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	s2, err := objects.ToStringArg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.p.WritePromptLine(s1, s2)
	return nil, nil
}

// doWritePromptEOL writes a prompt with an end-of-line flag based on the provided string and boolean arguments.
func (l *Library) doWritePromptEOL(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	b2, err := objects.ToBoolArg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.p.WritePromptEOL(s1, b2)
	return nil, nil
}

// doWrite writes a string with optional newline control using provided arguments; validates input and returns errors if invalid.
func (l *Library) doWrite(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	b2, err := objects.ToBoolArg(1, args[1])
	if err != nil {
		return nil, err
	}
	l.p.Write(s1, b2)
	return nil, nil
}

// doProcessSetSelfForeground sets the current process as the foreground and ensures no arguments are passed during invocation.
func (l *Library) doProcessSetSelfForeground(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.ProcessSetSelfForeground()
	return nil, nil
}

// doProcessSetForeground sets the specified process as the foreground process using the provided integer argument.
func (l *Library) doProcessSetForeground(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.p.ProcessSetForeground(int(i1))
	return nil, nil
}

// doProcessList retrieves a map of running processes, including details like name, PID, line, and start time.
// Returns an error if any arguments are provided.
func (l *Library) doProcessList(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	res := make(map[string]objects.IObject)
	for _, p := range l.p.ProcessList() {
		c := map[string]objects.IObject{
			"line": objects.NewStringNoSize(p.Line()),
			"name": objects.NewStringNoSize(p.Name()),
			"time": objects.NewTime(p.Time()),
			"pid":  objects.NewInt(int64(p.PID())),
		}
		res[p.Name()] = objects.NewMap(c)
	}
	return objects.NewMap(res), nil
}

// doWindowsSelectionPrevious navigates to the previous selection in the Windows selection context.
// Returns an error if any arguments are provided.
// Calls the underlying WindowsSelectionPrevious method to perform the operation.
func (l *Library) doWindowsSelectionPrevious(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.WindowsSelectionPrevious()
	return nil, nil
}

// doWindowsSelectionNext advances to the next item in the windows selection process. Returns an error for invalid arguments.
func (l *Library) doWindowsSelectionNext(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	l.p.WindowsSelectionNext()
	return nil, nil
}

// doWindowsSelectionOptions configures windows selection behavior using provided arguments as input parameters.
// Expects exactly two arguments: a rune (int64) and a float64.
// Returns an error if argument count is invalid or conversion of arguments fails.
// Invokes WindowsSelectionOptions method with the parsed inputs.
func (l *Library) doWindowsSelectionOptions(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments")
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	f2, err := objects.ToFloat64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	l.p.WindowsSelectionOptions(rune(i1), f2)
	return nil, nil
}

package process

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Library represents a collection of modules that interact with system processes and provide various functionalities.
// It contains a reference to a Process and a map of module names to their respective objects.
type Library struct {
	process   *Process
	functions []objects.IObject
}

// NewLibrary creates and initializes a new Library instance with the provided Process.
func NewLibrary(factory objects.IGateKeeper, process *Process) *Library {
	l := &Library{
		process:   process,
		functions: []objects.IObject{},
	}
	l.functions = []objects.IObject{
		factory.NewFuncImport(objects.FrameStatic, "Printf", -1, l.doPrintf),
		factory.NewFuncImport(objects.FrameStatic, "CreateTimer", 3, l.doCreateTimer),
		factory.NewFuncImport(objects.FrameStatic, "IsActive", 1, l.doIsActive),
		factory.NewFuncImport(objects.FrameStatic, "Kill", 1, l.doKill),
		factory.NewFuncImport(objects.FrameStatic, "KillForeground", 0, l.doKillForeground),
		factory.NewFuncImport(objects.FrameStatic, "KillAll", 1, l.doKillAll),
		factory.NewFuncImport(objects.FrameStatic, "CWDSet", 1, l.doCWDSet),
		factory.NewFuncImport(objects.FrameStatic, "CWDName", 0, l.doCWDName),
		factory.NewFuncImport(objects.FrameStatic, "CWDPath", 0, l.doCWDPath),
		factory.NewFuncImport(objects.FrameStatic, "CWDDirectoryListing", 0, l.doCWDDirectoryListing),
		factory.NewFuncImport(objects.FrameStatic, "GetScreenSize", 0, l.doGetScreenSize),
		factory.NewFuncImport(objects.FrameStatic, "PaintRequest", 0, l.doPaintRequest),
		factory.NewFuncImport(objects.FrameStatic, "ProcessExec", 1, l.doProcessExec),
		factory.NewFuncImport(objects.FrameStatic, "WindowsSelectionBegin", 0, l.doWindowsSelectionBegin),
		factory.NewFuncImport(objects.FrameStatic, "CWDSet", 0, l.doWindowsSelectionEnd),
		factory.NewFuncImport(objects.FrameStatic, "WindowsSelectionOptions", 2, l.doWindowsSelectionOptions),
		factory.NewFuncImport(objects.FrameStatic, "WindowsSelectionNext", 0, l.doWindowsSelectionNext),
		factory.NewFuncImport(objects.FrameStatic, "WindowsSelectionPrevious", 0, l.doWindowsSelectionPrevious),
		factory.NewFuncImport(objects.FrameStatic, "ProcessList", 0, l.doProcessList),
		factory.NewFuncImport(objects.FrameStatic, "ProcessSetForeground", 1, l.doProcessSetForeground),
		factory.NewFuncImport(objects.FrameStatic, "ProcessSetSelfForeground", 0, l.doProcessSetSelfForeground),
		factory.NewFuncImport(objects.FrameStatic, "Write", 2, l.doWrite),
		factory.NewFuncImport(objects.FrameStatic, "WritePromptEOL", 2, l.doWritePromptEOL),
		factory.NewFuncImport(objects.FrameStatic, "WritePromptLine", 2, l.doWritePromptLine),
		factory.NewFuncImport(objects.FrameStatic, "WriteColor", 5, l.doWriteColor),
		factory.NewFuncImport(objects.FrameStatic, "WriteForeground", 3, l.doWriteForeground),
		factory.NewFuncImport(objects.FrameStatic, "MoveCursorLeft", 0, l.doMoveCursorLeft),
		factory.NewFuncImport(objects.FrameStatic, "MoveCursorRight", 0, l.doMoveCursorRight),
		factory.NewFuncImport(objects.FrameStatic, "SaveCursor", 0, l.doSaveCursor),
		factory.NewFuncImport(objects.FrameStatic, "RestoreCursor", 0, l.doRestoreCursor),
		factory.NewFuncImport(objects.FrameStatic, "ClearScreen", 0, l.doClearScreen),
		factory.NewFuncImport(objects.FrameStatic, "SetExit", 0, l.doSetExit),
		factory.NewFuncImport(objects.FrameStatic, "Suggestion", 2, l.doSuggestion),
		factory.NewFuncImport(objects.FrameStatic, "Help", 1, l.doHelp),
	}

	return l
}

// Functions returns a slice of IObject representing the functions available in the Library.
func (l *Library) Functions() []objects.IObject {
	return l.functions
}

// doCreateTimer validates and extracts three integer arguments, then creates a timer using these arguments. Returns nil or an error.
func (l *Library) doCreateTimer(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.CreateTimer(int(i1), int(i2), int(i3))
	return 0, nil, nil
}

// doPrintf formats and writes a string output using provided arguments; returns an error for invalid input or formatting.
func (l *Library) doPrintf(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	var val string
	if len(args) == 1 {
		val = s1
	} else {
		var ar []interface{}
		for _, v := range args[1:] {
			ar = append(ar, gk.ToInterface(v))
		}
		val = fmt.Sprintf(s1, ar...)
	}
	l.process.Write(val, true)
	return 0, nil, nil
}

func (l *Library) doIsActive(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := l.process.IsActive(int(i1))
	return 1, gk.FromBool(v), nil
}

// doKill terminates a process identified by the provided argument, which must be convertible to an int64. Returns no value.
func (l *Library) doKill(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.Kill(int(i1))
	return 0, nil, nil
}

// doKillForeground terminates the foreground process if no arguments are provided, returning nil or an error.
func (l *Library) doKillForeground(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.KillForeground()
	return 0, nil, nil
}

// doKillAll terminates all processes matching the given context name, passed as a single string argument.
// Returns an error if the argument is invalid or conversion to string fails.
func (l *Library) doKillAll(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.KillAll(s1)
	return 0, nil, nil
}

// doWindowsSelectionEnd finalizes a windows selection operation, ensuring no arguments are passed and invoking WindowsSelectionEnd.
func (l *Library) doWindowsSelectionEnd(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.WindowsSelectionEnd()
	return 0, nil, nil
}

// doWindowsSelectionBegin initializes the window selection process in the current library context.
// Returns an error if invalid arguments are passed.
func (l *Library) doWindowsSelectionBegin(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.WindowsSelectionBegin()
	return 0, nil, nil
}

// doProcessExec validates arguments and executes a process based on a string argument. Returns an error if invalid.
func (l *Library) doProcessExec(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.ProcessExec(s1)
	return 0, nil, nil
}

// doPaintRequest performs a paint request by invoking the PaintRequest method of the Library's internal painter instance.
// It validates that no arguments are passed; otherwise, it returns an error.
// Returns nil if the operation succeeds or an error if the argument validation fails.
func (l *Library) doPaintRequest(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.PaintRequest()
	return 0, nil, nil
}

// doGetScreenSize retrieves the screen dimensions as a map containing "width" and "height" with their respective values.
func (l *Library) doGetScreenSize(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	w, h := l.process.GetScreenSize()
	return 1, gk.NewMap(frame, map[string]objects.IObject{
		"width":  gk.NewInt(frame, int64(w)),
		"height": gk.NewInt(frame, int64(h)),
	}), nil
}

// doCWDDirectoryListing returns the directory listing of the current working directory as an IObject array.
// It expects no arguments and raises an error if arguments are provided.
func (l *Library) doCWDDirectoryListing(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	v := l.process.CWDDirectoryListing()
	obj, err := gk.FromStringArray(frame, v)
	if err != nil {
		return 0, nil, err
	}
	return 1, obj, nil
}

// doCWDSet sets the current working directory (CWD) using the provided string argument. Returns a boolean as IObject.
func (l *Library) doCWDSet(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v := l.process.CWDSet(s1)
	return 1, gk.FromBool(v), nil
}

// doCWDPath returns the current working directory's path as a string object or an error if arguments are provided.
func (l *Library) doCWDPath(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	v := l.process.CWDPath()
	return 1, gk.NewString(frame, v), nil
}

// doCWDName retrieves the current working directory name, returns it as a string object, and validates argument count.
func (l *Library) doCWDName(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	v := l.process.CWDName()
	return 1, gk.NewString(frame, v), nil
}

// doHelp retrieves and returns a help string for the provided argument, which must be a single string-compatible object.
func (l *Library) doHelp(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	v, err := l.process.Help(s1)
	if err != nil {
		return 0, nil, err
	}
	return 1, gk.NewString(frame, v), nil
}

// doSuggestion processes two arguments: a string and an integer, to trigger the Suggestion functionality within the library.
// Returns an error for invalid argument types or count.
func (l *Library) doSuggestion(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.Suggestion(s1, int(i2))
	return 0, nil, nil
}

// doSetExit sets the exit condition for the associated process and returns nil if successful or an error otherwise.
// Accepts no arguments; an error is returned if arguments are provided.
func (l *Library) doSetExit(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.SetExit()
	return 0, nil, nil
}

// doClearScreen clears the terminal screen and returns nil if successful or an error if invalid arguments are passed.
func (l *Library) doClearScreen(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.ClearScreen()
	return 0, nil, nil
}

// doRestoreCursor restores the cursor to its previously saved position, returning an error for invalid arguments.
func (l *Library) doRestoreCursor(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.RestoreCursor()
	return 0, nil, nil
}

// doSaveCursor saves the current cursor position. Returns an error for invalid argument count.
func (l *Library) doSaveCursor(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.SaveCursor()
	return 0, nil, nil
}

// doMoveCursorRight moves the cursor to the right and accepts no arguments, returning an error if any are provided.
func (l *Library) doMoveCursorRight(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.MoveCursorRight()
	return 0, nil, nil
}

// doMoveCursorLeft moves the cursor one position to the left if the number of arguments is zero, otherwise returns an error.
func (l *Library) doMoveCursorLeft(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.MoveCursorLeft()
	return 0, nil, nil
}

// doWriteForeground writes text to the foreground using specified string, color, and a boolean flag for emphasis.
// Accepts three arguments: a string, an int64 representing color, and a boolean indicating emphasis.
// Returns an error if the argument count is invalid or conversion of arguments fails.
func (l *Library) doWriteForeground(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	b3, err := gk.ToBoolArg(2, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.WriteForeground(s1, interfaces.ColorDef(i2), b3)
	return 0, nil, nil
}

// doWriteColor writes colored text to the output using specified color definitions, mode, and a boolean flag.
func (l *Library) doWriteColor(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	i2, err := gk.ToInt64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	i3, err := gk.ToInt64Arg(2, args)
	if err != nil {
		return 0, nil, err
	}
	i4, err := gk.ToInt64Arg(3, args)
	if err != nil {
		return 0, nil, err
	}
	b5, err := gk.ToBoolArg(4, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.WriteColor(s1, interfaces.ColorDef(i2), interfaces.ColorDef(i3), interfaces.ColorMode(i4), b5)
	return 0, nil, nil
}

// doWritePromptLine writes a prompt line using two string arguments from objects.IObject and returns nil or an error.
func (l *Library) doWritePromptLine(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	s2, err := gk.ToStringArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.WritePromptLine(s1, s2)
	return 0, nil, nil
}

// doWritePromptEOL writes a prompt with an end-of-line flag based on the provided string and boolean arguments.
func (l *Library) doWritePromptEOL(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	b2, err := gk.ToBoolArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.WritePromptEOL(s1, b2)
	return 0, nil, nil
}

// doWrite writes a string with optional newline control using provided arguments; validates input and returns errors if invalid.
func (l *Library) doWrite(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	s1, err := gk.ToStringArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	b2, err := gk.ToBoolArg(1, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.Write(s1, b2)
	return 0, nil, nil
}

// doProcessSetSelfForeground sets the current process as the foreground and ensures no arguments are passed during invocation.
func (l *Library) doProcessSetSelfForeground(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.ProcessSetSelfForeground()
	return 0, nil, nil
}

// doProcessSetForeground sets the specified process as the foreground process using the provided integer argument.
func (l *Library) doProcessSetForeground(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.ProcessSetForeground(int(i1))
	return 0, nil, nil
}

// doProcessList retrieves a map of running processes, including details like name, PID, line, and start time.
// Returns an error if any arguments are provided.
func (l *Library) doProcessList(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	res := make(map[string]objects.IObject)
	for _, p := range l.process.ProcessList() {
		c := map[string]objects.IObject{
			"line": gk.NewString(frame, p.Line()),
			"name": gk.NewString(frame, p.Name()),
			"time": gk.NewTime(frame, p.Time()),
			"pid":  gk.NewInt(frame, int64(p.PID())),
		}
		res[p.Name()] = gk.NewMap(frame, c)
	}
	return 1, gk.NewMap(frame, res), nil
}

// doWindowsSelectionPrevious navigates to the previous selection in the Windows selection context.
// Returns an error if any arguments are provided.
// Calls the underlying WindowsSelectionPrevious method to perform the operation.
func (l *Library) doWindowsSelectionPrevious(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.WindowsSelectionPrevious()
	return 0, nil, nil
}

// doWindowsSelectionNext advances to the next item in the windows selection process. Returns an error for invalid arguments.
func (l *Library) doWindowsSelectionNext(_ objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	l.process.WindowsSelectionNext()
	return 0, nil, nil
}

// doWindowsSelectionOptions configures windows selection behavior using provided arguments as input parameters.
// Expects exactly two arguments: a rune (int64) and a float64.
// Returns an error if argument count is invalid or conversion of arguments fails.
// Invokes WindowsSelectionOptions method with the parsed inputs.
func (l *Library) doWindowsSelectionOptions(gk objects.IGateKeeper, _ int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	f2, err := gk.ToFloat64Arg(1, args)
	if err != nil {
		return 0, nil, err
	}
	l.process.WindowsSelectionOptions(rune(i1), f2)
	return 0, nil, nil
}

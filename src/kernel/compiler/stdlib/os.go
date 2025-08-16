package stdlib

import (
	"io"
	"os"
	"os/exec"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// osModule provides a mapping of OS-related constants and functions to their respective object representations or operations.
var osModule = map[string]objects.IObject{
	"o_rdonly":            objects.NewInt(int64(os.O_RDONLY)),
	"o_wronly":            objects.NewInt(int64(os.O_WRONLY)),
	"o_rdwr":              objects.NewInt(int64(os.O_RDWR)),
	"o_append":            objects.NewInt(int64(os.O_APPEND)),
	"o_create":            objects.NewInt(int64(os.O_CREATE)),
	"o_excl":              objects.NewInt(int64(os.O_EXCL)),
	"o_sync":              objects.NewInt(int64(os.O_SYNC)),
	"o_trunc":             objects.NewInt(int64(os.O_TRUNC)),
	"mode_dir":            objects.NewInt(int64(os.ModeDir)),
	"mode_append":         objects.NewInt(int64(os.ModeAppend)),
	"mode_exclusive":      objects.NewInt(int64(os.ModeExclusive)),
	"mode_temporary":      objects.NewInt(int64(os.ModeTemporary)),
	"mode_symlink":        objects.NewInt(int64(os.ModeSymlink)),
	"mode_device":         objects.NewInt(int64(os.ModeDevice)),
	"mode_named_pipe":     objects.NewInt(int64(os.ModeNamedPipe)),
	"mode_socket":         objects.NewInt(int64(os.ModeSocket)),
	"mode_setuid":         objects.NewInt(int64(os.ModeSetuid)),
	"mode_setgui":         objects.NewInt(int64(os.ModeSetgid)),
	"mode_char_device":    objects.NewInt(int64(os.ModeCharDevice)),
	"mode_sticky":         objects.NewInt(int64(os.ModeSticky)),
	"mode_type":           objects.NewInt(int64(os.ModeType)),
	"mode_perm":           objects.NewInt(int64(os.ModePerm)),
	"path_separator":      objects.NewChar(os.PathSeparator),
	"path_list_separator": objects.NewChar(os.PathListSeparator),
	"dev_null":            objects.NewStringNoSize(os.DevNull),
	"seek_set":            objects.NewInt(int64(io.SeekStart)),
	"seek_cur":            objects.NewInt(int64(io.SeekCurrent)),
	"seek_end":            objects.NewInt(int64(io.SeekEnd)),
	"args":                objects.NewFunctionModule(objects.FunctionModuleDef, "args", osArgs),                                     // args() => array(string)
	"chdir":               objects.NewFunctionModule(objects.FunctionModuleDef, "chdir", objects.FuncIsOe(os.Chdir)),                // chdir(dir string) => error
	"chmod":               osFuncASFmRE("chmod", os.Chmod),                                                                          // chmod(name string, mode int) => error
	"chown":               objects.NewFunctionModule(objects.FunctionModuleDef, "chown", objects.FuncIsiiOe(os.Chown)),              // chown(name string, uid int, gid int) => error
	"clearenv":            objects.NewFunctionModule(objects.FunctionModuleDef, "clearenv", objects.FuncInOn(os.Clearenv)),          // clearenv()
	"environ":             objects.NewFunctionModule(objects.FunctionModuleDef, "environ", objects.FuncInOsS(os.Environ)),           // environ() => array(string)
	"exit":                objects.NewFunctionModule(objects.FunctionModuleDef, "exit", objects.FuncIiOn(os.Exit)),                  // exit(code int)
	"expand_env":          objects.NewFunctionModule(objects.FunctionModuleDef, "expand_env", osExpandEnv),                          // expand_env(s string) => string
	"getegid":             objects.NewFunctionModule(objects.FunctionModuleDef, "getegid", objects.FuncInOi(os.Getegid)),            // getegid() => int
	"getenv":              objects.NewFunctionModule(objects.FunctionModuleDef, "getenv", objects.FuncIsOs(os.Getenv)),              // getenv(s string) => string
	"geteuid":             objects.NewFunctionModule(objects.FunctionModuleDef, "geteuid", objects.FuncInOi(os.Geteuid)),            // geteuid() => int
	"getgid":              objects.NewFunctionModule(objects.FunctionModuleDef, "getgid", objects.FuncInOi(os.Getgid)),              // getgid() => int
	"getgroups":           objects.NewFunctionModule(objects.FunctionModuleDef, "getgroups", objects.FuncInOiSe(os.Getgroups)),      // getgroups() => array(string)/error
	"getpagesize":         objects.NewFunctionModule(objects.FunctionModuleDef, "getpagesize", objects.FuncInOi(os.Getpagesize)),    // getpagesize() => int
	"getpid":              objects.NewFunctionModule(objects.FunctionModuleDef, "getpid", objects.FuncInOi(os.Getpid)),              // getpid() => int
	"getppid":             objects.NewFunctionModule(objects.FunctionModuleDef, "getppid", objects.FuncInOi(os.Getppid)),            // getppid() => int
	"getuid":              objects.NewFunctionModule(objects.FunctionModuleDef, "getuid", objects.FuncInOi(os.Getuid)),              // getuid() => int
	"getwd":               objects.NewFunctionModule(objects.FunctionModuleDef, "getwd", objects.FuncInOse(os.Getwd)),               // getwd() => string/error
	"hostname":            objects.NewFunctionModule(objects.FunctionModuleDef, "hostname", objects.FuncInOse(os.Hostname)),         // hostname() => string/error
	"lchown":              objects.NewFunctionModule(objects.FunctionModuleDef, "lchown", objects.FuncIsiiOe(os.Lchown)),            // lchown(name string, uid int, gid int) => error
	"link":                objects.NewFunctionModule(objects.FunctionModuleDef, "link", objects.FuncIssOe(os.Link)),                 // link(oldname string, newname string) => error
	"lookup_env":          objects.NewFunctionModule(objects.FunctionModuleDef, "lookup_env", osLookupEnv),                          // lookup_env(key string) => string/false
	"mkdir":               osFuncASFmRE("mkdir", os.Mkdir),                                                                          // mkdir(name string, perm int) => error
	"mkdir_all":           osFuncASFmRE("mkdir_all", os.MkdirAll),                                                                   // mkdir_all(name string, perm int) => error
	"readlink":            objects.NewFunctionModule(objects.FunctionModuleDef, "readlink", objects.FuncIsOse(os.Readlink)),         // readlink(name string) => string/error
	"remove":              objects.NewFunctionModule(objects.FunctionModuleDef, "remove", objects.FuncIsOe(os.Remove)),              // remove(name string) => error
	"remove_all":          objects.NewFunctionModule(objects.FunctionModuleDef, "remove_all", objects.FuncIsOe(os.RemoveAll)),       // remove_all(name string) => error
	"rename":              objects.NewFunctionModule(objects.FunctionModuleDef, "rename", objects.FuncIssOe(os.Rename)),             // rename(oldpath string, newpath string) => error
	"setenv":              objects.NewFunctionModule(objects.FunctionModuleDef, "setenv", objects.FuncIssOe(os.Setenv)),             // setenv(key string, value string) => error
	"symlink":             objects.NewFunctionModule(objects.FunctionModuleDef, "symlink", objects.FuncIssOe(os.Symlink)),           // symlink(oldname string newname string) => error
	"temp_dir":            objects.NewFunctionModule(objects.FunctionModuleDef, "temp_dir", objects.FuncInOs(os.TempDir)),           // temp_dir() => string
	"truncate":            objects.NewFunctionModule(objects.FunctionModuleDef, "truncate", objects.FuncIsi64Oe(os.Truncate)),       // truncate(name string, size int) => error
	"unsetenv":            objects.NewFunctionModule(objects.FunctionModuleDef, "unsetenv", objects.FuncIsOe(os.Unsetenv)),          // unsetenv(key string) => error
	"create":              objects.NewFunctionModule(objects.FunctionModuleDef, "create", osCreate),                                 // create(name string) => imap(file)/error
	"open":                objects.NewFunctionModule(objects.FunctionModuleDef, "open", osOpen),                                     // open(name string) => imap(file)/error
	"open_file":           objects.NewFunctionModule(objects.FunctionModuleDef, "open_file", osOpenFile),                            // open_file(name string, flag int, perm int) => imap(file)/error
	"find_process":        objects.NewFunctionModule(objects.FunctionModuleDef, "find_process", osFindProcess),                      // find_process(pid int) => imap(process)/error
	"start_process":       objects.NewFunctionModule(objects.FunctionModuleDef, "start_process", osStartProcess),                    // start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"exec_look_path":      objects.NewFunctionModule(objects.FunctionModuleDef, "exec_look_path", objects.FuncIsOse(exec.LookPath)), // exec_look_path(file) => string/error
	"exec":                objects.NewFunctionModule(objects.FunctionModuleDef, "exec", osExec),                                     // exec(name, args...) => command
	"stat":                objects.NewFunctionModule(objects.FunctionModuleDef, "stat", osStat),                                     // stat(name) => imap(fileinfo)/error
	"read_file":           objects.NewFunctionModule(objects.FunctionModuleDef, "read_file", osReadFile),                            // readfile(name) => array(byte)/error
}

// osReadFile reads the content of a file specified by its path and returns the content as an IObject or an error.
// Returns objects.ErrWrongNumArguments if the argument count is incorrect.
// Converts the input argument to a string using objects.ToStringArg.
// Reads the file using os.ReadFile; wraps any errors in objects.NewObjectError.
// Returns objects.ErrBytesLimit if the file content exceeds objects.MaxBytesLen.
// Wraps the file content in a Bytes object using objects.NewBytes on success.
func osReadFile(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	if len(bytes) > objects.MaxBytesLen {
		return nil, objects.ErrBytesLimit
	}
	return objects.NewBytes(bytes), nil
}

// osStat retrieves file or directory metadata based on the given path, returning it as an IObject map or an error.
func osStat(args ...objects.IObject) (ret objects.IObject, err error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	dir := objects.FalseValue
	if stat.IsDir() {
		dir = objects.TrueValue
	}
	fstat := objects.NewMapImmutable(map[string]objects.IObject{
		"name":      objects.NewStringNoSize(stat.Name()),
		"mtime":     objects.NewTime(stat.ModTime()),
		"size":      objects.NewInt(stat.Size()),
		"mode":      objects.NewInt(int64(stat.Mode())),
		"directory": dir,
	})
	return fstat, nil
}

// osCreate opens or creates a file specified by the single string argument and returns an IObject representing the file.
// Returns an error if the argument count is not 1, conversion to string fails, or the file cannot be created.
func osCreate(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := os.Create(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return makeOSFile(res), nil
}

// osOpen opens a file from the given string path argument and returns an IObject representing the file or an error.
func osOpen(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := os.Open(s1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return makeOSFile(res), nil
}

// osOpenFile opens a file with specified name, flag, and permission mode. Returns a wrapped os.File object or an error.
func osOpenFile(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 3 {
		return nil, objects.ErrWrongNumArguments
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
	res, err := os.OpenFile(s1, int(i2), os.FileMode(i3))
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return makeOSFile(res), nil
}

// osArgs retrieves the command-line arguments passed to the program as an array of IObject strings.
// Returns an error if the provided arguments are not zero or if any string exceeds the maximum allowed length.
func osArgs(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 0 {
		return nil, objects.ErrWrongNumArguments
	}
	arr := objects.NewArray(nil)
	for _, osArg := range os.Args {
		if len(osArg) > objects.MaxStringLen {
			return nil, objects.ErrStringLimit
		}
		v, err := objects.NewString(osArg)
		if err != nil {
			return nil, err
		}
		arr.Append(v)
	}
	return arr, nil
}

// osFuncASFmRE creates a FunctionModule object that wraps the provided name and a function processing string and FileMode.
// The wrapped function converts its arguments to string and FileMode, calls the passed function, and returns an error object.
// Returns an error if the number of arguments is not exactly 2 or if argument conversions fail.
func osFuncASFmRE(name string, fn func(string, os.FileMode) error) *objects.FunctionModule {
	return objects.NewFunctionModule(objects.FunctionModuleDef, name, func(args ...objects.IObject) (objects.IObject, error) {
		if len(args) != 2 {
			return nil, objects.ErrWrongNumArguments
		}
		s1, err := objects.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := objects.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return objects.NewObjectError(fn(s1, os.FileMode(i2))), nil
	})
}

// osLookupEnv retrieves the value of the environment variable named by the first argument and returns it as a string.
// If the variable does not exist, it returns a `FalseValue`.
// Returns an error if the number or type of arguments is invalid.
func osLookupEnv(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, ok := os.LookupEnv(s1)
	if !ok {
		return objects.FalseValue, nil
	}
	return objects.NewString(res)
}

// osExpandEnv replaces ${var} or $var in the input string based on the current values of the environment variables.
// Accepts exactly one string argument and returns a new string with variables expanded or an error if input is invalid.
func osExpandEnv(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var vLen int
	var failed bool
	s := os.Expand(s1, func(k string) string {
		if failed {
			return ""
		}
		v := os.Getenv(k)
		vLen += len(v)
		if vLen > objects.MaxStringLen {
			failed = true
			return ""
		}
		return v
	})
	return objects.NewString(s)
}

// osExec executes an external command with given arguments and returns a command object or an error.
// Returns ErrWrongNumArguments if no arguments are provided.
// Converts inputs to strings, constructing the command with initial arguments.
// Wraps the exec.Cmd instance into a map-like object with callable methods like Run, Start, and SetPath.
func osExec(args ...objects.IObject) (objects.IObject, error) {
	if len(args) == 0 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var execArgs []string
	for idx, arg := range args[1:] {
		execArg, err := objects.ToStringArg(idx, arg)
		if err != nil {
			return nil, err
		}
		execArgs = append(execArgs, execArg)
	}
	return makeOSExecCommand(exec.Command(s1, execArgs...)), nil
}

// osFindProcess retrieves an operating system process by its ID, represented as an IObject, and returns a wrapped process object.
// Returns an error if the provided arguments are invalid or if the process cannot be found.
func osFindProcess(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	proc, err := os.FindProcess(int(i1))
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return makeOSProcess(proc), nil
}

// osStartProcess starts a new process with the provided executable, arguments, working directory, and environment variables.
// Requires exactly 4 arguments: a string for the executable path, an array of string arguments, a string for the working directory,
// and an array of string environment variables. Returns an IObject representing the created process or an error.
func osStartProcess(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 4 {
		return nil, objects.ErrWrongNumArguments
	}
	s1, err := objects.ToStringArg(0, args[0])
	if err != nil {
		return nil, err
	}
	var argv []string
	switch arg1 := args[1].(type) {
	case *objects.Array:
		argv, err = objects.ToStringArrayArg(1, arg1.Values())
		if err != nil {
			return nil, err
		}
	case *objects.ArrayImmutable:
		argv, err = objects.ToStringArrayArg(1, arg1.Values())
		if err != nil {
			return nil, err
		}
	default:
		return nil, objects.NewInvalidArgumentError(1, "array", arg1.TypeName())
	}

	s2, err := objects.ToStringArg(2, args[2])
	if err != nil {
		return nil, err
	}
	var env []string
	switch arg3 := args[3].(type) {
	case *objects.Array:
		env, err = objects.ToStringArrayArg(3, arg3.Values())
		if err != nil {
			return nil, err
		}
	case *objects.ArrayImmutable:
		env, err = objects.ToStringArrayArg(3, arg3.Values())
		if err != nil {
			return nil, err
		}
	default:
		return nil, objects.NewInvalidArgumentError(3, "array", arg3.TypeName())
	}
	proc, err := os.StartProcess(s1, argv, &os.ProcAttr{Dir: s2, Env: env})
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return makeOSProcess(proc), nil
}

package stdlib

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// osModule provides a mapping of OS-related constants and functions to their respective object representations or operations.
var osModule = map[string]objects.IObject{
	"O_RDONLY":          objects.NewInt(int64(os.O_RDONLY)),
	"O_WRONLY":          objects.NewInt(int64(os.O_WRONLY)),
	"O_RDWR":            objects.NewInt(int64(os.O_RDWR)),
	"O_APPEND":          objects.NewInt(int64(os.O_APPEND)),
	"O_CREATE":          objects.NewInt(int64(os.O_CREATE)),
	"O_EXCL":            objects.NewInt(int64(os.O_EXCL)),
	"O_SYNC":            objects.NewInt(int64(os.O_SYNC)),
	"O_TRUNC":           objects.NewInt(int64(os.O_TRUNC)),
	"ModeDir":           objects.NewInt(int64(os.ModeDir)),
	"ModeAppend":        objects.NewInt(int64(os.ModeAppend)),
	"ModeExclusive":     objects.NewInt(int64(os.ModeExclusive)),
	"ModeTemporary":     objects.NewInt(int64(os.ModeTemporary)),
	"ModeSymlink":       objects.NewInt(int64(os.ModeSymlink)),
	"ModeDevice":        objects.NewInt(int64(os.ModeDevice)),
	"ModeNamedPipe":     objects.NewInt(int64(os.ModeNamedPipe)),
	"ModeSocket":        objects.NewInt(int64(os.ModeSocket)),
	"ModeSetuid":        objects.NewInt(int64(os.ModeSetuid)),
	"ModeSetgid":        objects.NewInt(int64(os.ModeSetgid)),
	"ModeCharDevice":    objects.NewInt(int64(os.ModeCharDevice)),
	"ModeSticky":        objects.NewInt(int64(os.ModeSticky)),
	"ModeType":          objects.NewInt(int64(os.ModeType)),
	"ModePerm":          objects.NewInt(int64(os.ModePerm)),
	"PathSeparator":     objects.NewChar(os.PathSeparator),
	"PathListSeparator": objects.NewChar(os.PathListSeparator),
	"DevNull":           objects.NewStringNoSize(os.DevNull),
	"SeekStart":         objects.NewInt(int64(io.SeekStart)),
	"SeekCurrent":       objects.NewInt(int64(io.SeekCurrent)),
	"SeekEnd":           objects.NewInt(int64(io.SeekEnd)),
	"Args":              objects.NewFunctionModule(objects.FunctionModuleDef, "Args", osArgs),                                   // args() => array(string)
	"Chdir":             objects.NewFunctionModule(objects.FunctionModuleDef, "Chdir", objects.FuncIsOe(os.Chdir)),              // chdir(dir string) => error
	"Chmod":             osFuncASFmRE("Chmod", os.Chmod),                                                                        // chmod(name string, mode int) => error
	"Chown":             objects.NewFunctionModule(objects.FunctionModuleDef, "Chown", objects.FuncIsiiOe(os.Chown)),            // chown(name string, uid int, gid int) => error
	"Clearenv":          objects.NewFunctionModule(objects.FunctionModuleDef, "Clearenv", objects.FuncInOn(os.Clearenv)),        // clearenv()
	"Environ":           objects.NewFunctionModule(objects.FunctionModuleDef, "Environ", objects.FuncInOsS(os.Environ)),         // environ() => array(string)
	"Exit":              objects.NewFunctionModule(objects.FunctionModuleDef, "Exit", objects.FuncIiOn(os.Exit)),                // exit(code int)
	"Expand":            objects.NewFunctionModule(objects.FunctionModuleDef, "Expand", osExpandEnv),                            // expand_env(s string) => string
	"Getegid":           objects.NewFunctionModule(objects.FunctionModuleDef, "Getegid", objects.FuncInOi(os.Getegid)),          // getegid() => int
	"Getenv":            objects.NewFunctionModule(objects.FunctionModuleDef, "Getenv", objects.FuncIsOs(os.Getenv)),            // getenv(s string) => string
	"geteuid":           objects.NewFunctionModule(objects.FunctionModuleDef, "Geteuid", objects.FuncInOi(os.Geteuid)),          // geteuid() => int
	"Getgid":            objects.NewFunctionModule(objects.FunctionModuleDef, "Getgid", objects.FuncInOi(os.Getgid)),            // getgid() => int
	"Getgroups":         objects.NewFunctionModule(objects.FunctionModuleDef, "Getgroups", objects.FuncInOiSe(os.Getgroups)),    // getgroups() => array(string)/error
	"Getpagesize":       objects.NewFunctionModule(objects.FunctionModuleDef, "Getpagesize", objects.FuncInOi(os.Getpagesize)),  // getpagesize() => int
	"Getpid":            objects.NewFunctionModule(objects.FunctionModuleDef, "Getpid", objects.FuncInOi(os.Getpid)),            // getpid() => int
	"Getppid":           objects.NewFunctionModule(objects.FunctionModuleDef, "Getppid", objects.FuncInOi(os.Getppid)),          // getppid() => int
	"Getuid":            objects.NewFunctionModule(objects.FunctionModuleDef, "Getuid", objects.FuncInOi(os.Getuid)),            // getuid() => int
	"Getwd":             objects.NewFunctionModule(objects.FunctionModuleDef, "Getwd", objects.FuncInOse(os.Getwd)),             // getwd() => string/error
	"Hostname":          objects.NewFunctionModule(objects.FunctionModuleDef, "Hostname", objects.FuncInOse(os.Hostname)),       // hostname() => string/error
	"Lchown":            objects.NewFunctionModule(objects.FunctionModuleDef, "Lchown", objects.FuncIsiiOe(os.Lchown)),          // lchown(name string, uid int, gid int) => error
	"Link":              objects.NewFunctionModule(objects.FunctionModuleDef, "Link", objects.FuncIssOe(os.Link)),               // link(oldname string, newname string) => error
	"LookupEnv":         objects.NewFunctionModule(objects.FunctionModuleDef, "LookupEnv", osLookupEnv),                         // lookup_env(key string) => string/false
	"Mkdir":             osFuncASFmRE("Mkdir", os.Mkdir),                                                                        // mkdir(name string, perm int) => error
	"MkdirAll":          osFuncASFmRE("MkdirAll", os.MkdirAll),                                                                  // mkdir_all(name string, perm int) => error
	"Readlink":          objects.NewFunctionModule(objects.FunctionModuleDef, "Readlink", objects.FuncIsOse(os.Readlink)),       // readlink(name string) => string/error
	"Remove":            objects.NewFunctionModule(objects.FunctionModuleDef, "Remove", objects.FuncIsOe(os.Remove)),            // remove(name string) => error
	"RemoveAll":         objects.NewFunctionModule(objects.FunctionModuleDef, "RemoveAll", objects.FuncIsOe(os.RemoveAll)),      // remove_all(name string) => error
	"Rename":            objects.NewFunctionModule(objects.FunctionModuleDef, "Rename", objects.FuncIssOe(os.Rename)),           // rename(oldpath string, newpath string) => error
	"Setenv":            objects.NewFunctionModule(objects.FunctionModuleDef, "Setenv", objects.FuncIssOe(os.Setenv)),           // setenv(key string, value string) => error
	"Symlink":           objects.NewFunctionModule(objects.FunctionModuleDef, "Symlink", objects.FuncIssOe(os.Symlink)),         // symlink(oldname string newname string) => error
	"TempDir":           objects.NewFunctionModule(objects.FunctionModuleDef, "TempDir", objects.FuncInOs(os.TempDir)),          // temp_dir() => string
	"Truncate":          objects.NewFunctionModule(objects.FunctionModuleDef, "Truncate", objects.FuncIsi64Oe(os.Truncate)),     // truncate(name string, size int) => error
	"Unsetenv":          objects.NewFunctionModule(objects.FunctionModuleDef, "Unsetenv", objects.FuncIsOe(os.Unsetenv)),        // unsetenv(key string) => error
	"Create":            objects.NewFunctionModule(objects.FunctionModuleDef, "Create", osCreate),                               // create(name string) => imap(file)/error
	"Open":              objects.NewFunctionModule(objects.FunctionModuleDef, "Open", osOpen),                                   // open(name string) => imap(file)/error
	"OpenFile":          objects.NewFunctionModule(objects.FunctionModuleDef, "OpenFile", osOpenFile),                           // open_file(name string, flag int, perm int) => imap(file)/error
	"FindProcess":       objects.NewFunctionModule(objects.FunctionModuleDef, "FindProcess", osFindProcess),                     // find_process(pid int) => imap(process)/error
	"StartProcess":      objects.NewFunctionModule(objects.FunctionModuleDef, "StartProcess", osStartProcess),                   // start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"ExecLookPath":      objects.NewFunctionModule(objects.FunctionModuleDef, "ExecLookPath", objects.FuncIsOse(exec.LookPath)), // exec_look_path(file) => string/error
	"Exec":              objects.NewFunctionModule(objects.FunctionModuleDef, "Exec", osExec),                                   // exec(name, args...) => command
	"Stat":              objects.NewFunctionModule(objects.FunctionModuleDef, "Stat", osStat),                                   // stat(name) => imap(fileinfo)/error
	"ReadFile":          objects.NewFunctionModule(objects.FunctionModuleDef, "ReadFile", osReadFile),                           // readfile(name) => array(byte)/error
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

func makeOSProcessState(state *os.ProcessState) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"Exited":  objects.NewFunctionModule(objects.FunctionModuleDef, "Exited", objects.FuncInOb(state.Exited)),
			"Pid":     objects.NewFunctionModule(objects.FunctionModuleDef, "Pid", objects.FuncInOi(state.Pid)),
			"String":  objects.NewFunctionModule(objects.FunctionModuleDef, "String", objects.FuncInOs(state.String)),
			"Success": objects.NewFunctionModule(objects.FunctionModuleDef, "Success", objects.FuncInOb(state.Success)),
		},
	)
}

func makeOSProcess(proc *os.Process) *objects.MapImmutable {
	return objects.NewMapImmutable(map[string]objects.IObject{
		"Kill":    objects.NewFunctionModule(objects.FunctionModuleDef, "Kill", objects.FuncInOe(proc.Kill)),
		"Release": objects.NewFunctionModule(objects.FunctionModuleDef, "Release", objects.FuncInOe(proc.Release)),
		"Signal": objects.NewFunctionModule(objects.FunctionModuleDef, "Signal", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			i1, err := objects.ToInt64Arg(0, args[0])
			if err != nil {
				return nil, err
			}
			return objects.NewObjectError(proc.Signal(syscall.Signal(i1))), nil
		}),
		"Wait": objects.NewFunctionModule(objects.FunctionModuleDef, "Wait", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}
			state, err := proc.Wait()
			if err != nil {
				return objects.NewObjectError(err), nil
			}
			return makeOSProcessState(state), nil
		}),
	})
}

// makeOSExecCommand returns an immutable map exposing methods to manipulate and control an exec.Cmd instance.
func makeOSExecCommand(cmd *exec.Cmd) *objects.MapImmutable {
	return objects.NewMapImmutable(map[string]objects.IObject{
		// combined_output() => bytes/error
		"CombinedOutput": objects.NewFunctionModule(objects.FunctionModuleDef, "CombinedOutput", objects.FuncInObSe(cmd.CombinedOutput)),
		// output() => bytes/error
		"Output": objects.NewFunctionModule(objects.FunctionModuleDef, "Output", objects.FuncInObSe(cmd.Output)), //
		// run() => error
		"Run": objects.NewFunctionModule(objects.FunctionModuleDef, "Run", objects.FuncInOe(cmd.Run)), //
		// start() => error
		"Start": objects.NewFunctionModule(objects.FunctionModuleDef, "Start", objects.FuncInOe(cmd.Start)), //
		// wait() => error
		"Wait": objects.NewFunctionModule(objects.FunctionModuleDef, "Wait", objects.FuncInOe(cmd.Wait)), //
		// set_path(path string)
		"SetPath": objects.NewFunctionModule(objects.FunctionModuleDef, "SetPath", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			s1, err := objects.ToStringArg(0, args[0])
			if err != nil {
				return nil, err
			}
			cmd.Path = s1
			return objects.UndefinedValue, nil
		}),
		// set_dir(dir string)
		"SetDir": objects.NewFunctionModule(objects.FunctionModuleDef, "SetDir", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			s1, err := objects.ToStringArg(0, args[0])
			if err != nil {
				return nil, err
			}
			cmd.Dir = s1
			return objects.UndefinedValue, nil
		}),
		// set_env(env array(string))
		"SetEnv": objects.NewFunctionModule(objects.FunctionModuleDef, "SetEnv", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			var env []string
			var err error
			switch arg0 := args[0].(type) {
			case *objects.Array:
				env, err = objects.ToStringArrayArg(0, arg0.Values())
				if err != nil {
					return nil, err
				}
			case *objects.ArrayImmutable:
				env, err = objects.ToStringArrayArg(0, arg0.Values())
				if err != nil {
					return nil, err
				}
			default:
				return nil, objects.NewInvalidArgumentError(0, "array", arg0.TypeName())
			}
			cmd.Env = env
			return objects.UndefinedValue, nil
		}),
		// process() => imap(process)
		"Process": objects.NewFunctionModule(objects.FunctionModuleDef, "Process", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}
			return makeOSProcess(cmd.Process), nil
		}),
	})
}

// makeOSFile creates an MapImmutable containing methods applicable to an os.File object as IObject values.
func makeOSFile(file *os.File) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			// chdir() => true/error
			"Chdir": objects.NewFunctionModule(objects.FunctionModuleDef, "Chdir", objects.FuncInOe(file.Chdir)), //
			// chown(uid int, gid int) => true/error
			"Chown": objects.NewFunctionModule(objects.FunctionModuleDef, "Chown", objects.FuncIiiOe(file.Chown)), //
			// close() => error
			"Close": objects.NewFunctionModule(objects.FunctionModuleDef, "Close", objects.FuncInOe(file.Close)), //
			// name() => string
			"Name": objects.NewFunctionModule(objects.FunctionModuleDef, "Name", objects.FuncInOs(file.Name)), //
			// readdirnames(n int) => array(string)/error
			"Readdirnames": objects.NewFunctionModule(objects.FunctionModuleDef, "Readdirnames", objects.FuncIiOsSe(file.Readdirnames)), //
			// sync() => error
			"Sync": objects.NewFunctionModule(objects.FunctionModuleDef, "Sync", objects.FuncInOe(file.Sync)), //
			// write(bytes) => int/error
			"Write": objects.NewFunctionModule(objects.FunctionModuleDef, "Write", objects.FuncIbSOie(file.Write)), //
			// write(string) => int/error
			"WriteString": objects.NewFunctionModule(objects.FunctionModuleDef, "WriteString", objects.FuncIsOie(file.WriteString)), //
			// read(bytes) => int/error
			"Read": objects.NewFunctionModule(objects.FunctionModuleDef, "Read", objects.FuncIbSOie(file.Read)), //
			// chmod(mode int) => error
			"Chmod": objects.NewFunctionModule(objects.FunctionModuleDef, "Chmod", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 1 {
					return nil, objects.ErrWrongNumArguments
				}
				i1, err := objects.ToInt64Arg(0, args[0])
				if err != nil {
					return nil, err
				}
				return objects.NewObjectError(file.Chmod(os.FileMode(i1))), nil
			}),
			// seek(offset int, whence int) => int/error
			"Seek": objects.NewFunctionModule(objects.FunctionModuleDef, "Seek", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 2 {
					return nil, objects.ErrWrongNumArguments
				}
				i1, err := objects.ToInt64Arg(0, args[0])
				if err != nil {
					return nil, err
				}
				i2, err := objects.ToInt64Arg(1, args[1])
				if err != nil {
					return nil, err
				}
				res, err := file.Seek(i1, int(i2))
				if err != nil {
					return objects.NewObjectError(err), nil
				}
				return objects.NewInt(res), nil
			}),
			// stat() => imap(fileinfo)/error
			"Stat": objects.NewFunctionModule(objects.FunctionModuleDef, "Stat", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 0 {
					return nil, objects.ErrWrongNumArguments
				}
				return osStat(objects.NewStringNoSize(file.Name()))
			}),
		})
}

package stdlib

import (
	"io"
	"os"
	"os/exec"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

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
	"args":                objects.NewFunctionUser("args", osArgs),                                     // args() => array(string)
	"chdir":               objects.NewFunctionUser("chdir", objects.FuncASRE(os.Chdir)),                // chdir(dir string) => error
	"chmod":               osFuncASFmRE("chmod", os.Chmod),                                             // chmod(name string, mode int) => error
	"chown":               objects.NewFunctionUser("chown", objects.FuncASIIRE(os.Chown)),              // chown(name string, uid int, gid int) => error
	"clearenv":            objects.NewFunctionUser("clearenv", objects.FuncAR(os.Clearenv)),            // clearenv()
	"environ":             objects.NewFunctionUser("environ", objects.FuncARSs(os.Environ)),            // environ() => array(string)
	"exit":                objects.NewFunctionUser("exit", objects.FuncAIR(os.Exit)),                   // exit(code int)
	"expand_env":          objects.NewFunctionUser("expand_env", osExpandEnv),                          // expand_env(s string) => string
	"getegid":             objects.NewFunctionUser("getegid", objects.FuncARI(os.Getegid)),             // getegid() => int
	"getenv":              objects.NewFunctionUser("getenv", objects.FuncASRS(os.Getenv)),              // getenv(s string) => string
	"geteuid":             objects.NewFunctionUser("geteuid", objects.FuncARI(os.Geteuid)),             // geteuid() => int
	"getgid":              objects.NewFunctionUser("getgid", objects.FuncARI(os.Getgid)),               // getgid() => int
	"getgroups":           objects.NewFunctionUser("getgroups", objects.FuncARIsE(os.Getgroups)),       // getgroups() => array(string)/error
	"getpagesize":         objects.NewFunctionUser("getpagesize", objects.FuncARI(os.Getpagesize)),     // getpagesize() => int
	"getpid":              objects.NewFunctionUser("getpid", objects.FuncARI(os.Getpid)),               // getpid() => int
	"getppid":             objects.NewFunctionUser("getppid", objects.FuncARI(os.Getppid)),             // getppid() => int
	"getuid":              objects.NewFunctionUser("getuid", objects.FuncARI(os.Getuid)),               // getuid() => int
	"getwd":               objects.NewFunctionUser("getwd", objects.FuncARSE(os.Getwd)),                // getwd() => string/error
	"hostname":            objects.NewFunctionUser("hostname", objects.FuncARSE(os.Hostname)),          // hostname() => string/error
	"lchown":              objects.NewFunctionUser("lchown", objects.FuncASIIRE(os.Lchown)),            // lchown(name string, uid int, gid int) => error
	"link":                objects.NewFunctionUser("link", objects.FuncASSRE(os.Link)),                 // link(oldname string, newname string) => error
	"lookup_env":          objects.NewFunctionUser("lookup_env", osLookupEnv),                          // lookup_env(key string) => string/false
	"mkdir":               osFuncASFmRE("mkdir", os.Mkdir),                                             // mkdir(name string, perm int) => error
	"mkdir_all":           osFuncASFmRE("mkdir_all", os.MkdirAll),                                      // mkdir_all(name string, perm int) => error
	"readlink":            objects.NewFunctionUser("readlink", objects.FuncASRSE(os.Readlink)),         // readlink(name string) => string/error
	"remove":              objects.NewFunctionUser("remove", objects.FuncASRE(os.Remove)),              // remove(name string) => error
	"remove_all":          objects.NewFunctionUser("remove_all", objects.FuncASRE(os.RemoveAll)),       // remove_all(name string) => error
	"rename":              objects.NewFunctionUser("rename", objects.FuncASSRE(os.Rename)),             // rename(oldpath string, newpath string) => error
	"setenv":              objects.NewFunctionUser("setenv", objects.FuncASSRE(os.Setenv)),             // setenv(key string, value string) => error
	"symlink":             objects.NewFunctionUser("symlink", objects.FuncASSRE(os.Symlink)),           // symlink(oldname string newname string) => error
	"temp_dir":            objects.NewFunctionUser("temp_dir", objects.FuncARS(os.TempDir)),            // temp_dir() => string
	"truncate":            objects.NewFunctionUser("truncate", objects.FuncASI64RE(os.Truncate)),       // truncate(name string, size int) => error
	"unsetenv":            objects.NewFunctionUser("unsetenv", objects.FuncASRE(os.Unsetenv)),          // unsetenv(key string) => error
	"create":              objects.NewFunctionUser("create", osCreate),                                 // create(name string) => imap(file)/error
	"open":                objects.NewFunctionUser("open", osOpen),                                     // open(name string) => imap(file)/error
	"open_file":           objects.NewFunctionUser("open_file", osOpenFile),                            // open_file(name string, flag int, perm int) => imap(file)/error
	"find_process":        objects.NewFunctionUser("find_process", osFindProcess),                      // find_process(pid int) => imap(process)/error
	"start_process":       objects.NewFunctionUser("start_process", osStartProcess),                    // start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"exec_look_path":      objects.NewFunctionUser("exec_look_path", objects.FuncASRSE(exec.LookPath)), // exec_look_path(file) => string/error
	"exec":                objects.NewFunctionUser("exec", osExec),                                     // exec(name, args...) => command
	"stat":                objects.NewFunctionUser("stat", osStat),                                     // stat(name) => imap(fileinfo)/error
	"read_file":           objects.NewFunctionUser("read_file", osReadFile),                            // readfile(name) => array(byte)/error
}

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

func osFuncASFmRE(name string, fn func(string, os.FileMode) error) *objects.FunctionUser {
	return objects.NewFunctionUser(name, func(args ...objects.IObject) (objects.IObject, error) {
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

package stdlib

import (
	"os/exec"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// makeOSExecCommand returns an immutable map exposing methods to manipulate and control an exec.Cmd instance.
func makeOSExecCommand(cmd *exec.Cmd) *objects.ImmutableMap {
	return objects.NewImmutableMap(map[string]objects.IObject{
		// combined_output() => bytes/error
		"combined_output": objects.NewUserFunction("combined_output", objects.FuncARYE(cmd.CombinedOutput)),
		// output() => bytes/error
		"output": objects.NewUserFunction("output", objects.FuncARYE(cmd.Output)), //
		// run() => error
		"run": objects.NewUserFunction("run", objects.FuncARE(cmd.Run)), //
		// start() => error
		"start": objects.NewUserFunction("start", objects.FuncARE(cmd.Start)), //
		// wait() => error
		"wait": objects.NewUserFunction("wait", objects.FuncARE(cmd.Wait)), //
		// set_path(path string)
		"set_path": objects.NewUserFunction("set_path", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
			}
			cmd.Path = s1
			return objects.UndefinedValue, nil
		}),
		// set_dir(dir string)
		"set_dir": objects.NewUserFunction("set_dir", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			s1, ok := objects.ToString(args[0])
			if !ok {
				return nil, objects.NewInvalidArgumentError("first", "string(compatible)", args[0].TypeName())
			}
			cmd.Dir = s1
			return objects.UndefinedValue, nil
		}),
		// set_env(env array(string))
		"set_env": objects.NewUserFunction("set_env", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			var env []string
			var err error
			switch arg0 := args[0].(type) {
			case *objects.Array:
				env, err = stringArray(arg0.Values(), "first")
				if err != nil {
					return nil, err
				}
			case *objects.ImmutableArray:
				env, err = stringArray(arg0.Values(), "first")
				if err != nil {
					return nil, err
				}
			default:
				return nil, objects.NewInvalidArgumentError("first", "array", arg0.TypeName())
			}
			cmd.Env = env
			return objects.UndefinedValue, nil
		}),
		// process() => imap(process)
		"process": objects.NewUserFunction("process", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 0 {
				return nil, objects.ErrWrongNumArguments
			}
			return makeOSProcess(cmd.Process), nil
		}),
	})
}

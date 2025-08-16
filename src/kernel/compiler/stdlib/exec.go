package stdlib

import (
	"os/exec"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

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

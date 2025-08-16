package stdlib

import (
	"os"
	"syscall"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

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

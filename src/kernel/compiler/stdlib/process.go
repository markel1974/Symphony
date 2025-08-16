package stdlib

import (
	"os"
	"syscall"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func makeOSProcessState(state *os.ProcessState) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"Exited":  objects.NewFunctionModule("Exited", objects.FuncInOb(state.Exited)),
			"Pid":     objects.NewFunctionModule("Pid", objects.FuncInOi(state.Pid)),
			"String":  objects.NewFunctionModule("String", objects.FuncInOs(state.String)),
			"Success": objects.NewFunctionModule("Success", objects.FuncInOb(state.Success)),
		},
	)
}

func makeOSProcess(proc *os.Process) *objects.MapImmutable {
	return objects.NewMapImmutable(map[string]objects.IObject{
		"Kill":    objects.NewFunctionModule("Kill", objects.FuncInOe(proc.Kill)),
		"Release": objects.NewFunctionModule("Release", objects.FuncInOe(proc.Release)),
		"Signal": objects.NewFunctionModule("Signal", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			i1, err := objects.ToInt64Arg(0, args[0])
			if err != nil {
				return nil, err
			}
			return objects.NewObjectError(proc.Signal(syscall.Signal(i1))), nil
		}),
		"Wait": objects.NewFunctionModule("Wait", func(args ...objects.IObject) (objects.IObject, error) {
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

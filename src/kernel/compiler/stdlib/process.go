package stdlib

import (
	"os"
	"syscall"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func makeOSProcessState(state *os.ProcessState) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"Exited":  objects.NewFunctionUser("Exited", objects.FuncARB(state.Exited)),
			"Pid":     objects.NewFunctionUser("Pid", objects.FuncARI(state.Pid)),
			"String":  objects.NewFunctionUser("String", objects.FuncARS(state.String)),
			"Success": objects.NewFunctionUser("Success", objects.FuncARB(state.Success)),
		},
	)
}

func makeOSProcess(proc *os.Process) *objects.MapImmutable {
	return objects.NewMapImmutable(map[string]objects.IObject{
		"Kill":    objects.NewFunctionUser("Kill", objects.FuncARE(proc.Kill)),
		"Release": objects.NewFunctionUser("Release", objects.FuncARE(proc.Release)),
		"Signal": objects.NewFunctionUser("Signal", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			i1, err := objects.ToInt64Arg("first", args[0])
			if err != nil {
				return nil, err
			}
			return objects.NewObjectError(proc.Signal(syscall.Signal(i1))), nil
		}),
		"Wait": objects.NewFunctionUser("Wait", func(args ...objects.IObject) (objects.IObject, error) {
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

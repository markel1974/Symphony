package stdlib

import (
	"os"
	"syscall"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func makeOSProcessState(state *os.ProcessState) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"exited":  objects.NewFunctionUser("exited", objects.FuncARB(state.Exited)),
			"pid":     objects.NewFunctionUser("pid", objects.FuncARI(state.Pid)),
			"string":  objects.NewFunctionUser("string", objects.FuncARS(state.String)),
			"success": objects.NewFunctionUser("success", objects.FuncARB(state.Success)),
		},
	)
}

func makeOSProcess(proc *os.Process) *objects.MapImmutable {
	return objects.NewMapImmutable(map[string]objects.IObject{
		"kill":    objects.NewFunctionUser("kill", objects.FuncARE(proc.Kill)),
		"release": objects.NewFunctionUser("release", objects.FuncARE(proc.Release)),
		"signal": objects.NewFunctionUser("signal", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			i1, err := objects.ToInt64Arg("first", args[0])
			if err != nil {
				return nil, err
			}
			return objects.NewObjectError(proc.Signal(syscall.Signal(i1))), nil
		}),
		"wait": objects.NewFunctionUser("wait", func(args ...objects.IObject) (objects.IObject, error) {
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

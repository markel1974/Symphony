package stdlib

import (
	"os"
	"syscall"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func makeOSProcessState(state *os.ProcessState) *objects.ImmutableMap {
	return objects.NewImmutableMap(
		map[string]objects.IObject{
			"exited":  objects.NewUserFunction("exited", objects.FuncARB(state.Exited)),
			"pid":     objects.NewUserFunction("pid", objects.FuncARI(state.Pid)),
			"string":  objects.NewUserFunction("string", objects.FuncARS(state.String)),
			"success": objects.NewUserFunction("success", objects.FuncARB(state.Success)),
		},
	)
}

func makeOSProcess(proc *os.Process) *objects.ImmutableMap {
	return objects.NewImmutableMap(map[string]objects.IObject{
		"kill":    objects.NewUserFunction("kill", objects.FuncARE(proc.Kill)),
		"release": objects.NewUserFunction("release", objects.FuncARE(proc.Release)),
		"signal": objects.NewUserFunction("signal", func(args ...objects.IObject) (objects.IObject, error) {
			if len(args) != 1 {
				return nil, objects.ErrWrongNumArguments
			}
			i1, err := objects.ToInt64Arg("first", args[0])
			if err != nil {
				return nil, err
			}
			return objects.NewObjectError(proc.Signal(syscall.Signal(i1))), nil
		}),
		"wait": objects.NewUserFunction("wait", func(args ...objects.IObject) (objects.IObject, error) {
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

package stdlib

import (
	"os"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

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

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
			"Chdir": objects.NewFunctionModule("Chdir", objects.FuncInOe(file.Chdir)), //
			// chown(uid int, gid int) => true/error
			"Chown": objects.NewFunctionModule("Chown", objects.FuncIiiOe(file.Chown)), //
			// close() => error
			"Close": objects.NewFunctionModule("Close", objects.FuncInOe(file.Close)), //
			// name() => string
			"Name": objects.NewFunctionModule("Name", objects.FuncInOs(file.Name)), //
			// readdirnames(n int) => array(string)/error
			"Readdirnames": objects.NewFunctionModule("Readdirnames", objects.FuncIiOsSe(file.Readdirnames)), //
			// sync() => error
			"Sync": objects.NewFunctionModule("Sync", objects.FuncInOe(file.Sync)), //
			// write(bytes) => int/error
			"Write": objects.NewFunctionModule("Write", objects.FuncIbSOie(file.Write)), //
			// write(string) => int/error
			"WriteString": objects.NewFunctionModule("WriteString", objects.FuncIsOie(file.WriteString)), //
			// read(bytes) => int/error
			"Read": objects.NewFunctionModule("Read", objects.FuncIbSOie(file.Read)), //
			// chmod(mode int) => error
			"Chmod": objects.NewFunctionModule("Chmod", func(args ...objects.IObject) (objects.IObject, error) {
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
			"Seek": objects.NewFunctionModule("Seek", func(args ...objects.IObject) (objects.IObject, error) {
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
			"Stat": objects.NewFunctionModule("Stat", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 0 {
					return nil, objects.ErrWrongNumArguments
				}
				return osStat(objects.NewStringNoSize(file.Name()))
			}),
		})
}

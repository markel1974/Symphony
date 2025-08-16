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
			"Chdir": objects.NewFunctionUser("Chdir", objects.FuncARE(file.Chdir)), //
			// chown(uid int, gid int) => true/error
			"Chown": objects.NewFunctionUser("Chown", objects.FuncAIIRE(file.Chown)), //
			// close() => error
			"Close": objects.NewFunctionUser("Close", objects.FuncARE(file.Close)), //
			// name() => string
			"Name": objects.NewFunctionUser("Name", objects.FuncARS(file.Name)), //
			// readdirnames(n int) => array(string)/error
			"Readdirnames": objects.NewFunctionUser("Readdirnames", objects.FuncAIRSsE(file.Readdirnames)), //
			// sync() => error
			"Sync": objects.NewFunctionUser("Sync", objects.FuncARE(file.Sync)), //
			// write(bytes) => int/error
			"Write": objects.NewFunctionUser("Write", objects.FuncAYRIE(file.Write)), //
			// write(string) => int/error
			"WriteString": objects.NewFunctionUser("WriteString", objects.FuncASRIE(file.WriteString)), //
			// read(bytes) => int/error
			"Read": objects.NewFunctionUser("Read", objects.FuncAYRIE(file.Read)), //
			// chmod(mode int) => error
			"Chmod": objects.NewFunctionUser("Chmod", func(args ...objects.IObject) (objects.IObject, error) {
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
			"Seek": objects.NewFunctionUser("Seek", func(args ...objects.IObject) (objects.IObject, error) {
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
			"Stat": objects.NewFunctionUser("Stat", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 0 {
					return nil, objects.ErrWrongNumArguments
				}
				return osStat(objects.NewStringNoSize(file.Name()))
			}),
		})
}

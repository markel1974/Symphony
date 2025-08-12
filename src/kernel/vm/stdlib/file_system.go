package stdlib

import (
	"os"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// makeOSFile creates an ImmutableMap containing methods applicable to an os.File object as IObject values.
func makeOSFile(file *os.File) *objects.ImmutableMap {
	return objects.NewImmutableMap(
		map[string]objects.IObject{
			// chdir() => true/error
			"chdir": objects.NewUserFunction("chdir", objects.FuncARE(file.Chdir)), //
			// chown(uid int, gid int) => true/error
			"chown": objects.NewUserFunction("chown", objects.FuncAIIRE(file.Chown)), //
			// close() => error
			"close": objects.NewUserFunction("close", objects.FuncARE(file.Close)), //
			// name() => string
			"name": objects.NewUserFunction("name", objects.FuncARS(file.Name)), //
			// readdirnames(n int) => array(string)/error
			"readdirnames": objects.NewUserFunction("readdirnames", objects.FuncAIRSsE(file.Readdirnames)), //
			// sync() => error
			"sync": objects.NewUserFunction("sync", objects.FuncARE(file.Sync)), //
			// write(bytes) => int/error
			"write": objects.NewUserFunction("write", objects.FuncAYRIE(file.Write)), //
			// write(string) => int/error
			"write_string": objects.NewUserFunction("write_string", objects.FuncASRIE(file.WriteString)), //
			// read(bytes) => int/error
			"read": objects.NewUserFunction("read", objects.FuncAYRIE(file.Read)), //
			// chmod(mode int) => error
			"chmod": objects.NewUserFunction("chmod", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 1 {
					return nil, objects.ErrWrongNumArguments
				}
				i1, ok := objects.ToInt64(args[0])
				if !ok {
					return nil, objects.NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
				}
				return objects.NewObjectError(file.Chmod(os.FileMode(i1))), nil
			}),
			// seek(offset int, whence int) => int/error
			"seek": objects.NewUserFunction("seek", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 2 {
					return nil, objects.ErrWrongNumArguments
				}
				i1, ok := objects.ToInt64(args[0])
				if !ok {
					return nil, objects.NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
				}
				i2, ok := objects.ToInt(args[1])
				if !ok {
					return nil, objects.NewInvalidArgumentError("second", "int(compatible)", args[1].TypeName())
				}
				res, err := file.Seek(i1, i2)
				if err != nil {
					return objects.NewObjectError(err), nil
				}
				return objects.NewInt(res), nil
			}),
			// stat() => imap(fileinfo)/error
			"stat": objects.NewUserFunction("stat", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 0 {
					return nil, objects.ErrWrongNumArguments
				}
				return osStat(objects.NewStringNoSize(file.Name()))
			}),
		})
}

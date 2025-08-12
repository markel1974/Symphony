package stdlib

import (
	"os"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// makeOSFile creates an ImmutableMap containing methods applicable to an os.File object as IObject values.
func makeOSFile(file *os.File) *objects.ImmutableMap {
	return objects.NewImmutableMap(
		map[string]objects.IObject{
			// chdir() => true/error
			"chdir": objects.NewUserFunction("chdir", FuncARE(file.Chdir)), //
			// chown(uid int, gid int) => true/error
			"chown": objects.NewUserFunction("chown", FuncAIIRE(file.Chown)), //
			// close() => error
			"close": objects.NewUserFunction("close", FuncARE(file.Close)), //
			// name() => string
			"name": objects.NewUserFunction("name", FuncARS(file.Name)), //
			// readdirnames(n int) => array(string)/error
			"readdirnames": objects.NewUserFunction("readdirnames", FuncAIRSsE(file.Readdirnames)), //
			// sync() => error
			"sync": objects.NewUserFunction("sync", FuncARE(file.Sync)), //
			// write(bytes) => int/error
			"write": objects.NewUserFunction("write", FuncAYRIE(file.Write)), //
			// write(string) => int/error
			"write_string": objects.NewUserFunction("write_string", FuncASRIE(file.WriteString)), //
			// read(bytes) => int/error
			"read": objects.NewUserFunction("read", FuncAYRIE(file.Read)), //
			// chmod(mode int) => error
			"chmod": objects.NewUserFunction("chmod", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 1 {
					return nil, errors.ErrWrongNumArguments
				}
				i1, ok := objects.ToInt64(args[0])
				if !ok {
					return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
				}
				return wrapError(file.Chmod(os.FileMode(i1))), nil
			}),
			// seek(offset int, whence int) => int/error
			"seek": objects.NewUserFunction("seek", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 2 {
					return nil, errors.ErrWrongNumArguments
				}
				i1, ok := objects.ToInt64(args[0])
				if !ok {
					return nil, errors.NewInvalidArgumentType("first", "int(compatible)", args[0].TypeName())
				}
				i2, ok := objects.ToInt(args[1])
				if !ok {
					return nil, errors.NewInvalidArgumentType("second", "int(compatible)", args[1].TypeName())
				}
				res, err := file.Seek(i1, i2)
				if err != nil {
					return wrapError(err), nil
				}
				return objects.NewInt(res), nil
			}),
			// stat() => imap(fileinfo)/error
			"stat": objects.NewUserFunction("stat", func(args ...objects.IObject) (objects.IObject, error) {
				if len(args) != 0 {
					return nil, errors.ErrWrongNumArguments
				}
				return osStat(objects.NewStringNoSize(file.Name()))
			}),
		})
}

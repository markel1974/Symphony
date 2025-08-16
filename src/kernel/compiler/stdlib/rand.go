package stdlib

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// randModule is a map containing random-related utility functions for generating numbers, seeding, permutations, and more.
var _randModule = map[string]objects.IObject{
	"Int63":       objects.NewFunctionModule(objects.FunctionModuleDef, "Int63", objects.FuncInOi64(rand.Int63)),
	"Float64":     objects.NewFunctionModule(objects.FunctionModuleDef, "Float64", objects.FuncInOf64(rand.Float64)),
	"Int63n":      objects.NewFunctionModule(objects.FunctionModuleDef, "Int63n", objects.FuncIi64Oi64(rand.Int63n)),
	"ExpFloat64":  objects.NewFunctionModule(objects.FunctionModuleDef, "ExpFloat64", objects.FuncInOf64(rand.ExpFloat64)),
	"NormFloat64": objects.NewFunctionModule(objects.FunctionModuleDef, "NormFloat64", objects.FuncInOf64(rand.NormFloat64)),
	"Perm":        objects.NewFunctionModule(objects.FunctionModuleDef, "Perm", objects.FuncIiOiS(rand.Perm)),
	"Seed":        objects.NewFunctionModule(objects.FunctionModuleDef, "Seed", objects.FuncIi64On(rand.Seed)),
	"Read":        objects.NewFunctionModule(objects.FunctionModuleDef, "Read", doRandRead),
	"Rand":        objects.NewFunctionModule(objects.FunctionModuleDef, "Rand", doRandRand),
}

// randRand generates an immutable map containing random-related functions derived from the provided rand.Rand instance.
func randRand(r *rand.Rand) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"Int63":       objects.NewFunctionModule(objects.FunctionModuleDef, "Int63", objects.FuncInOi64(r.Int63)),
			"Float64":     objects.NewFunctionModule(objects.FunctionModuleDef, "Float64", objects.FuncInOf64(r.Float64)),
			"Int63n":      objects.NewFunctionModule(objects.FunctionModuleDef, "Int63n", objects.FuncIi64Oi64(r.Int63n)),
			"ExpFloat64":  objects.NewFunctionModule(objects.FunctionModuleDef, "ExpFloat64", objects.FuncInOf64(r.ExpFloat64)),
			"NormFloat64": objects.NewFunctionModule(objects.FunctionModuleDef, "NormFloat64", objects.FuncInOf64(r.NormFloat64)),
			"Perm":        objects.NewFunctionModule(objects.FunctionModuleDef, "Perm", objects.FuncIiOiS(r.Perm)),
			"Seed":        objects.NewFunctionModule(objects.FunctionModuleDef, "Seed", objects.FuncIi64On(r.Seed)),
			"Read":        objects.NewFunctionModule(objects.FunctionModuleDef, "Read", func(args ...objects.IObject) (objects.IObject, error) { return doRRandRand(r, args...) }),
		})
}

// doRandRead reads random bytes into the given bytes object and returns the number of bytes read as an integer.
// It expects exactly one argument of type *objects.Bytes. Returns an error if the argument is missing or invalid.
func doRandRead(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := objects.ToByteSliceArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := rand.Read(bs1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(int64(res)), nil
}

// doRandRand initializes a new random generator with the given seed and returns a map of random-related functions.
// Expects one argument: an int(compatible) seed value. Returns an error for invalid arguments or wrong argument count.
func doRandRand(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := objects.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	src := rand.NewSource(i1)
	return randRand(rand.New(src)), nil
}

// doRRandRand reads random data into a byte array and returns the number of bytes written or an error if any occurs.
// The function expects exactly one argument of type *objects.Bytes; otherwise, it returns an argument or type error.
func doRRandRand(r *rand.Rand, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := objects.ToByteSliceArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := r.Read(bs1)
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(int64(res)), nil
}

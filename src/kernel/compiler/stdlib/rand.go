package stdlib

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// randModule is a map containing random-related utility functions for generating numbers, seeding, permutations, and more.
var _randModule = map[string]objects.IObject{
	"Int63":       objects.NewFunctionUser("Int63", objects.FuncARI64(rand.Int63)),
	"Float64":     objects.NewFunctionUser("Float64", objects.FuncARF(rand.Float64)),
	"Int63n":      objects.NewFunctionUser("Int63n", objects.FuncAI64RI64(rand.Int63n)),
	"ExpFloat64":  objects.NewFunctionUser("ExpFloat64", objects.FuncARF(rand.ExpFloat64)),
	"NormFloat64": objects.NewFunctionUser("NormFloat64", objects.FuncARF(rand.NormFloat64)),
	"Perm":        objects.NewFunctionUser("Perm", objects.FuncAIRIs(rand.Perm)),
	"Seed":        objects.NewFunctionUser("Seed", objects.FuncAI64R(rand.Seed)),
	"Read":        objects.NewFunctionUser("Read", doRandRead),
	"Rand":        objects.NewFunctionUser("Rand", doRandRand),
}

// randRand generates an immutable map containing random-related functions derived from the provided rand.Rand instance.
func randRand(r *rand.Rand) *objects.MapImmutable {
	return objects.NewMapImmutable(
		map[string]objects.IObject{
			"Int63":       objects.NewFunctionUser("Int63", objects.FuncARI64(r.Int63)),
			"Float64":     objects.NewFunctionUser("Float64", objects.FuncARF(r.Float64)),
			"Int63n":      objects.NewFunctionUser("Int63n", objects.FuncAI64RI64(r.Int63n)),
			"ExpFloat64":  objects.NewFunctionUser("ExpFloat64", objects.FuncARF(r.ExpFloat64)),
			"NormFloat64": objects.NewFunctionUser("NormFloat64", objects.FuncARF(r.NormFloat64)),
			"Perm":        objects.NewFunctionUser("Perm", objects.FuncAIRIs(r.Perm)),
			"Seed":        objects.NewFunctionUser("Seed", objects.FuncAI64R(r.Seed)),
			"Read":        objects.NewFunctionUser("Read", func(args ...objects.IObject) (objects.IObject, error) { return doRRandRand(r, args...) }),
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

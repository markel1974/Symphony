package stdlib

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// randModule is a map containing random-related utility functions for generating numbers, seeding, permutations, and more.
var randModule = map[string]objects.IObject{
	"int":        objects.NewUserFunction("int", objects.FuncARI64(rand.Int63)),
	"float":      objects.NewUserFunction("float", objects.FuncARF(rand.Float64)),
	"intn":       objects.NewUserFunction("intn", objects.FuncAI64RI64(rand.Int63n)),
	"exp_float":  objects.NewUserFunction("exp_float", objects.FuncARF(rand.ExpFloat64)),
	"norm_float": objects.NewUserFunction("norm_float", objects.FuncARF(rand.NormFloat64)),
	"perm":       objects.NewUserFunction("perm", objects.FuncAIRIs(rand.Perm)),
	"seed":       objects.NewUserFunction("seed", objects.FuncAI64R(rand.Seed)),
	"read":       objects.NewUserFunction("read", doRandRead),
	"rand":       objects.NewUserFunction("rand", doRandRand),
}

// randRand generates an immutable map containing random-related functions derived from the provided rand.Rand instance.
func randRand(r *rand.Rand) *objects.ImmutableMap {
	return objects.NewImmutableMap(
		map[string]objects.IObject{
			"int":        objects.NewUserFunction("int", objects.FuncARI64(r.Int63)),
			"float":      objects.NewUserFunction("float", objects.FuncARF(r.Float64)),
			"intn":       objects.NewUserFunction("intn", objects.FuncAI64RI64(r.Int63n)),
			"exp_float":  objects.NewUserFunction("exp_float", objects.FuncARF(r.ExpFloat64)),
			"norm_float": objects.NewUserFunction("norm_float", objects.FuncARF(r.NormFloat64)),
			"perm":       objects.NewUserFunction("perm", objects.FuncAIRIs(r.Perm)),
			"seed":       objects.NewUserFunction("seed", objects.FuncAI64R(r.Seed)),
			"read":       objects.NewUserFunction("read", func(args ...objects.IObject) (objects.IObject, error) { return doRRandRand(r, args...) }),
		})
}

// doRandRead reads random bytes into the given bytes object and returns the number of bytes read as an integer.
// It expects exactly one argument of type *objects.Bytes. Returns an error if the argument is missing or invalid.
func doRandRead(args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	y1, ok := args[0].(*objects.Bytes)
	if !ok {
		return nil, objects.NewInvalidArgumentError("first", "bytes", args[0].TypeName())
	}
	res, err := rand.Read(y1.Value())
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
	i1, ok := objects.ToInt64(args[0])
	if !ok {
		return nil, objects.NewInvalidArgumentError("first", "int(compatible)", args[0].TypeName())
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
	y1, ok := args[0].(*objects.Bytes)
	if !ok {
		return nil, objects.NewInvalidArgumentError("first", "bytes", args[0].TypeName())
	}
	res, err := r.Read(y1.Value())
	if err != nil {
		return objects.NewObjectError(err), nil
	}
	return objects.NewInt(int64(res)), nil
}

package sdk

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Rand is a struct that encapsulates a module mapping of string keys to objects implementing the IObject interface.
type Rand struct {
	*Package
	factory *objects.GateKeeper
}

// NewRand creates a new instance of Rand with a pre-defined set of random number generation functions.
func NewRand(factory *objects.GateKeeper) *Rand {
	z := &Rand{
		factory: factory,
	}
	container := []objects.IObject{
		factory.NewFuncPackage(objects.FuncPackageDef, "Int63", factory.FuncInOi64(rand.Int63)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Float64", factory.FuncInOf64(rand.Float64)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Int63n", factory.FuncIi64Oi64(rand.Int63n)),
		factory.NewFuncPackage(objects.FuncPackageDef, "ExpFloat64", factory.FuncInOf64(rand.ExpFloat64)),
		factory.NewFuncPackage(objects.FuncPackageDef, "NormFloat64", factory.FuncInOf64(rand.NormFloat64)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Perm", factory.FuncIiOiS(rand.Perm)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Seed", factory.FuncIi64On(rand.Seed)),
		factory.NewFuncPackage(objects.FuncPackageDef, "Read", z.read),
		factory.NewFuncPackage(objects.FuncPackageDef, "Rand", z.rand),
	}
	z.Package = NewPackage("rand", container, nil)
	return z
}

// Read reads random data into a byte slice and returns the number of bytes written as an integer or an error if it occurs.
func (z *Rand) read(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := z.factory.ToByteSliceArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := rand.Read(bs1)
	if err != nil {
		return z.factory.NewError(frame, err.Error()), nil
	}
	return z.factory.NewInt(frame, int64(res)), nil
}

// Rand generates a new random number generator using the provided seed argument and returns its options as a map.
func (z *Rand) rand(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := z.factory.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	src := rand.NewSource(i1)
	r := rand.New(src)
	return z.factory.NewMapImmutable(frame,
		map[string]objects.IObject{
			"Int63":       z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Int63", z.factory.FuncInOi64(r.Int63)),
			"Float64":     z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Float64", z.factory.FuncInOf64(r.Float64)),
			"Int63n":      z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Int63n", z.factory.FuncIi64Oi64(r.Int63n)),
			"ExpFloat64":  z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "ExpFloat64", z.factory.FuncInOf64(r.ExpFloat64)),
			"NormFloat64": z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "NormFloat64", z.factory.FuncInOf64(r.NormFloat64)),
			"Perm":        z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Perm", z.factory.FuncIiOiS(r.Perm)),
			"Seed":        z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Seed", z.factory.FuncIi64On(r.Seed)),
			"Read": z.factory.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Read", func(frame int, args ...objects.IObject) (objects.IObject, error) {
				return z.randOptionsRead(r, frame, args...)
			}),
		}), nil
}

// RandOptionsRead reads random data into a provided byte slice, returning the number of bytes read or an error.
func (z *Rand) randOptionsRead(r *rand.Rand, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := z.factory.ToByteSliceArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := r.Read(bs1)
	if err != nil {
		return z.factory.NewError(frame, err.Error()), nil
	}
	return z.factory.NewInt(frame, int64(res)), nil
}

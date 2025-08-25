package sdk

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Rand is a struct that encapsulates a module mapping of string keys to objects implementing the IObject interface.
type Rand struct {
	*Package
	gk objects.IGateKeeper
}

// NewRand creates a new instance of Rand with a pre-defined set of random number generation functions.
func NewRand(gk objects.IGateKeeper) *Rand {
	z := &Rand{
		gk: gk,
	}
	container := []objects.IObject{
		gk.NewFuncPackage(objects.FuncPackageDef, "Int63", z.int63(rand.Int63)),
		gk.NewFuncPackage(objects.FuncPackageDef, "Float64", z.funcInOf64(rand.Float64)),
		gk.NewFuncPackage(objects.FuncPackageDef, "Int63n", z.int63n(rand.Int63n)),
		gk.NewFuncPackage(objects.FuncPackageDef, "ExpFloat64", z.funcInOf64(rand.ExpFloat64)),
		gk.NewFuncPackage(objects.FuncPackageDef, "NormFloat64", z.funcInOf64(rand.NormFloat64)),
		gk.NewFuncPackage(objects.FuncPackageDef, "Perm", gk.FuncIiOiS(rand.Perm)),
		gk.NewFuncPackage(objects.FuncPackageDef, "Seed", gk.FuncIi64On(rand.Seed)),
		gk.NewFuncPackage(objects.FuncPackageDef, "Read", z.read),
		gk.NewFuncPackage(objects.FuncPackageDef, "Rand", z.rand),
	}
	z.Package = NewPackage("rand", container, nil)
	return z
}

// Read reads random data into a byte slice and returns the number of bytes written as an integer or an error if it occurs.
func (z *Rand) read(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := z.gk.ToBytesArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := rand.Read(bs1)
	if err != nil {
		return z.gk.NewError(frame, err.Error()), nil
	}
	return z.gk.NewInt(frame, int64(res)), nil
}

// Rand generates a new random number generator using the provided seed argument and returns its options as a map.
func (z *Rand) rand(frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	i1, err := z.gk.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	src := rand.NewSource(i1)
	r := rand.New(src)
	return z.gk.NewMapImmutable(frame,
		map[string]objects.IObject{
			"Int63":       z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Int63", z.int63(r.Int63)),
			"Float64":     z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Float64", z.funcInOf64(r.Float64)),
			"Int63n":      z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Int63n", z.int63n(r.Int63n)),
			"ExpFloat64":  z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "ExpFloat64", z.funcInOf64(r.ExpFloat64)),
			"NormFloat64": z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "NormFloat64", z.funcInOf64(r.NormFloat64)),
			"Perm":        z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Perm", z.gk.FuncIiOiS(r.Perm)),
			"Seed":        z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Seed", z.gk.FuncIi64On(r.Seed)),
			"Read": z.gk.NewFuncPackageFrame(frame, objects.FuncPackageDef, "Read", func(frame int, args ...objects.IObject) (objects.IObject, error) {
				return z.randOptionsRead(r, frame, args...)
			}),
		}), nil
}

// RandOptionsRead reads random data into a provided byte slice, returning the number of bytes read or an error.
func (z *Rand) randOptionsRead(r *rand.Rand, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrWrongNumArguments
	}
	bs1, err := z.gk.ToBytesArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := r.Read(bs1)
	if err != nil {
		return z.gk.NewError(frame, err.Error()), nil
	}
	return z.gk.NewInt(frame, int64(res)), nil
}

// int63 is a method that returns a callable function producing a random int64 via the provided generator function (fn).
// The callable function accepts zero arguments and raises an error if arguments are provided.
func (z *Rand) int63(fn func() int64) objects.FuncCallable {
	return func(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, objects.ErrWrongNumArguments
		}
		return z.gk.NewInt(frame, fn()), nil
	}
}

// int63n returns a function that applies a provided transformation to a 63-bit signed integer argument.
// The function takes a single int64 argument, applies the provided function to it, and returns the result as an IObject.
// An error is returned if the number of arguments is different from one.
func (z *Rand) int63n(fn func(int64) int64) objects.FuncCallable {
	return func(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, objects.ErrWrongNumArguments
		}
		i1, err := z.gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return z.gk.NewInt(frame, fn(i1)), nil
	}
}

// funcInOf64 wraps a no-argument function returning float64 into a FuncCallable to integrate with the GateAdapter system.
// It enforces no arguments and converts the float64 result to an IObject, returning ErrWrongNumArguments for invalid input.
func (z *Rand) funcInOf64(fn func() float64) objects.FuncCallable {
	return func(frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, objects.ErrWrongNumArguments
		}
		return z.gk.NewFloat(frame, fn()), nil
	}
}

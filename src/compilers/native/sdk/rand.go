package sdk

import (
	"fmt"
	"math/rand"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// init initializes the package by registering the NewRand package using the RegisterPackage function.
func init() {
	RegisterPackage(NewRand)
}

// Rand is a type that provides a container for random-related functionalities and operations.
type Rand struct {
	*bytecode.Package
}

// NewRand creates and initializes a new instance of the Rand package, registering its functions with the provided IGateKeeper.
func NewRand(gk objects.IGateKeeper) bytecode.IPackage {
	z := &Rand{}
	container := []objects.IObject{
		gk.NewFuncImport(objects.FrameStatic, "Int63", 0, z.int63(rand.Int63)),
		gk.NewFuncImport(objects.FrameStatic, "Float64", 0, z.float64(rand.Float64)),
		gk.NewFuncImport(objects.FrameStatic, "Int63n", 1, z.int63n(rand.Int63n)),
		gk.NewFuncImport(objects.FrameStatic, "ExpFloat64", 0, z.float64(rand.ExpFloat64)),
		gk.NewFuncImport(objects.FrameStatic, "NormFloat64", 0, z.float64(rand.NormFloat64)),
		gk.NewFuncImport(objects.FrameStatic, "Perm", 1, z.perm(rand.Perm)),
		gk.NewFuncImport(objects.FrameStatic, "Seed", 1, z.seed(rand.Seed)),
		gk.NewFuncImport(objects.FrameStatic, "Operand", 1, z.read),
		gk.NewFuncImport(objects.FrameStatic, "Rand", 1, z.rand),
	}
	z.Package = bytecode.NewPackage("rand", container, nil)
	return z
}

// read reads random data into the provided byte array argument. Returns the number of bytes read or an error.
func (z *Rand) read(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := rand.Read(bs1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewInt(frame, int64(res)), nil
}

// rand initializes a new random number generator using the provided seed and returns a map of random generator methods.
func (z *Rand) rand(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	i1, err := gk.ToInt64Arg(0, args)
	if err != nil {
		return 0, nil, err
	}
	src := rand.NewSource(i1)
	r := rand.New(src)
	return 1, gk.NewMap(frame,
		map[string]objects.IObject{
			"Int63":       gk.NewFuncImport(frame, "Int63", 0, z.int63(r.Int63)),
			"Float64":     gk.NewFuncImport(frame, "Float64", 0, z.float64(r.Float64)),
			"Int63n":      gk.NewFuncImport(frame, "Int63n", 1, z.int63n(r.Int63n)),
			"ExpFloat64":  gk.NewFuncImport(frame, "ExpFloat64", 0, z.float64(r.ExpFloat64)),
			"NormFloat64": gk.NewFuncImport(frame, "NormFloat64", 0, z.float64(r.NormFloat64)),
			"Perm":        gk.NewFuncImport(frame, "Perm", 1, z.perm(r.Perm)),
			"Seed":        gk.NewFuncImport(frame, "Seed", 1, z.seed(r.Seed)),
			"Operand": gk.NewFuncImport(frame, "Operand", 1, func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
				return z.randOptionsRead(gk, r, frame, args...)
			}),
		}), nil
}

// randOptionsRead reads random bytes into a byte slice using a specific random number generator and returns the result size.
func (z *Rand) randOptionsRead(gk objects.IGateKeeper, r *rand.Rand, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
	bs1, err := gk.ToBytesArg(0, args)
	if err != nil {
		return 0, nil, err
	}
	res, err := r.Read(bs1)
	if err != nil {
		return 0, gk.NewError(frame, err.Error()), nil
	}
	return 1, gk.NewInt(frame, int64(res)), nil
}

// int63 returns a Invocable that generates a 63-bit non-negative integer using the provided function.
// It accepts no arguments and returns the generated integer as an IObject.
// Returns an error if any arguments are provided.
func (z *Rand) int63(fn func() int64) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		return 1, gk.NewInt(frame, fn()), nil
	}
}

// int63n returns a Invocable that applies a provided int64 function to a single int64 argument and returns the result.
func (z *Rand) int63n(fn func(int64) int64) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		return 1, gk.NewInt(frame, fn(i1)), nil
	}
}

// seed returns a Invocable that sets a seed using the provided function, ensuring the argument is a single int64 value.
func (z *Rand) seed(fn func(int64)) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		fn(i1)
		return 0, gk.UndefinedValue(), nil
	}
}

// float64 generates a Invocable that produces a float64 value using the provided function and returns it as an IObject.
func (z *Rand) float64(fn func() float64) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		return 1, gk.NewFloat(frame, fn()), nil
	}
}

// perm generates a Invocable that calculates permutations of size n using the provided function.
// The input is a single integer argument, and the result is an Array containing the permutation sequence.
// Returns an error if the input is invalid or the conversion fails.
func (z *Rand) perm(fn func(int) []int) objects.Invocable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
		i1, err := gk.ToInt64Arg(0, args)
		if err != nil {
			return 0, nil, err
		}
		res := fn(int(i1))
		obj := gk.NewArray(frame, nil)
		arr, ok := obj.(*objects.Array)
		if !ok {
			return 0, nil, fmt.Errorf("expected Array, got %T", obj)
		}
		for _, v := range res {
			arr.Append(gk.NewInt(frame, int64(v)))
		}
		return 1, arr, nil
	}
}

package sdk

import (
	"fmt"
	"math/rand"

	"github.com/markel1974/symphony/src/vm/bytecode"
	"github.com/markel1974/symphony/src/vm/objects"
)

// init initializes the package by registering the NewRand package using the register function.
func init() {
	register(NewRand)
}

const (
	defRandInt63       = "Int63"
	defRandFloat64     = "Float64"
	defRandInt63n      = "Int63n"
	defRandExpFloat64  = "ExpFloat64"
	defRandNormFloat64 = "NormFloat64"
	defRandPerm        = "Perm"
	defRandSeed        = "Seed"
	defRandOperand     = "Operand"
	defRandRand        = "Rand"
)

// Rand is a type that provides a container for random-related functionalities and operations.
type Rand struct {
	*bytecode.Package
}

// NewRand creates and initializes a new instance of the Rand package, registering its functions with the provided IGateKeeper.
func NewRand(gk objects.IGateKeeper) bytecode.IPackage {
	z := &Rand{Package: bytecode.NewPackage("rand")}

	z.Add(defRandInt63, gk.NewFuncImport(objects.FrameStatic, defRandInt63, 0, z.int63(rand.Int63)))
	z.Add(defRandFloat64, gk.NewFuncImport(objects.FrameStatic, defRandFloat64, 0, z.float64(rand.Float64)))
	z.Add(defRandInt63n, gk.NewFuncImport(objects.FrameStatic, defRandInt63n, 1, z.int63n(rand.Int63n)))
	z.Add(defRandExpFloat64, gk.NewFuncImport(objects.FrameStatic, defRandExpFloat64, 0, z.float64(rand.ExpFloat64)))
	z.Add(defRandNormFloat64, gk.NewFuncImport(objects.FrameStatic, defRandNormFloat64, 0, z.float64(rand.NormFloat64)))
	z.Add(defRandPerm, gk.NewFuncImport(objects.FrameStatic, defRandPerm, 1, z.perm(rand.Perm)))
	z.Add(defRandSeed, gk.NewFuncImport(objects.FrameStatic, defRandSeed, 1, z.seed(rand.Seed)))
	z.Add(defRandOperand, gk.NewFuncImport(objects.FrameStatic, defRandOperand, 1, z.read))
	z.Add(defRandRand, gk.NewFuncImport(objects.FrameStatic, defRandRand, 1, z.rand))

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
			defRandInt63:       gk.NewFuncImport(frame, defRandInt63, 0, z.int63(r.Int63)),
			defRandFloat64:     gk.NewFuncImport(frame, defRandFloat64, 0, z.float64(r.Float64)),
			defRandInt63n:      gk.NewFuncImport(frame, defRandInt63n, 1, z.int63n(r.Int63n)),
			defRandExpFloat64:  gk.NewFuncImport(frame, defRandExpFloat64, 0, z.float64(r.ExpFloat64)),
			defRandNormFloat64: gk.NewFuncImport(frame, defRandNormFloat64, 0, z.float64(r.NormFloat64)),
			defRandPerm:        gk.NewFuncImport(frame, defRandPerm, 1, z.perm(r.Perm)),
			defRandSeed:        gk.NewFuncImport(frame, defRandSeed, 1, z.seed(r.Seed)),
			defRandOperand: gk.NewFuncImport(frame, defRandOperand, 1, func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (uint, objects.IObject, error) {
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

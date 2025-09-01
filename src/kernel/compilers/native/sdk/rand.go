package sdk

import (
	"math/rand"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	RegisterPackage(NewRand)
}

// Rand is a struct that encapsulates a module mapping of string keys to objects implementing the IObject interface.
type Rand struct {
	container map[string]objects.IObject
}

// NewRand creates a new instance of Rand with a pre-defined set of random number generation functions.
func NewRand(gk objects.IGateKeeper) IPackage {
	z := &Rand{}
	container := []objects.IObject{
		gk.NewFuncPackage("Int63", z.int63(rand.Int63)),
		gk.NewFuncPackage("Float64", z.funcInOf64(rand.Float64)),
		gk.NewFuncPackage("Int63n", z.int63n(rand.Int63n)),
		gk.NewFuncPackage("ExpFloat64", z.funcInOf64(rand.ExpFloat64)),
		gk.NewFuncPackage("NormFloat64", z.funcInOf64(rand.NormFloat64)),
		gk.NewFuncPackage("Perm", gk.FuncIiOiS(rand.Perm)),
		gk.NewFuncPackage("Seed", gk.FuncIi64On(rand.Seed)),
		gk.NewFuncPackage("Read", z.read),
		gk.NewFuncPackage("Rand", z.rand),
	}
	z.container = BuildContainer(container, nil)
	return z
}

// Name returns the name of the Rand module as a string.
func (z *Rand) Name() string {
	return "rand"
}

// Get retrieves an object associated with the given name from the container. It returns the object and a boolean indicating success.
func (z *Rand) Get(name string) (objects.IObject, bool) {
	v, ok := z.container[name]
	return v, ok
}

// Read reads random data into a byte slice and returns the number of bytes written as an integer or an error if it occurs.
func (z *Rand) read(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	bs1, err := gk.ToBytesArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := rand.Read(bs1)
	if err != nil {
		return gk.NewError(frame, err.Error()), nil
	}
	return gk.NewInt(frame, int64(res)), nil
}

// Rand generates a new random number generator using the provided seed argument and returns its options as a map.
func (z *Rand) rand(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	i1, err := gk.ToInt64Arg(0, args[0])
	if err != nil {
		return nil, err
	}
	src := rand.NewSource(i1)
	r := rand.New(src)
	return gk.NewMap(frame,
		map[string]objects.IObject{
			"Int63":       gk.NewFuncPackageFrame(frame, "Int63", z.int63(r.Int63)),
			"Float64":     gk.NewFuncPackageFrame(frame, "Float64", z.funcInOf64(r.Float64)),
			"Int63n":      gk.NewFuncPackageFrame(frame, "Int63n", z.int63n(r.Int63n)),
			"ExpFloat64":  gk.NewFuncPackageFrame(frame, "ExpFloat64", z.funcInOf64(r.ExpFloat64)),
			"NormFloat64": gk.NewFuncPackageFrame(frame, "NormFloat64", z.funcInOf64(r.NormFloat64)),
			"Perm":        gk.NewFuncPackageFrame(frame, "Perm", gk.FuncIiOiS(r.Perm)),
			"Seed":        gk.NewFuncPackageFrame(frame, "Seed", gk.FuncIi64On(r.Seed)),
			"Read": gk.NewFuncPackageFrame(frame, "Read", func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (objects.IObject, error) {
				return z.randOptionsRead(gk, r, frame, args...)
			}),
		}), nil
}

// RandOptionsRead reads random data into a provided byte slice, returning the number of bytes read or an error.
func (z *Rand) randOptionsRead(gk objects.IGateKeeper, r *rand.Rand, frame int, args ...objects.IObject) (objects.IObject, error) {
	if len(args) != 1 {
		return nil, objects.ErrInvalidArgumentsNumber
	}
	bs1, err := gk.ToBytesArg(0, args[0])
	if err != nil {
		return nil, err
	}
	res, err := r.Read(bs1)
	if err != nil {
		return gk.NewError(frame, err.Error()), nil
	}
	return gk.NewInt(frame, int64(res)), nil
}

// int63 is a method that returns a callable function producing a random int64 via the provided generator function (fn).
// The callable function accepts zero arguments and raises an error if arguments are provided.
func (z *Rand) int63(fn func() int64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, objects.ErrInvalidArgumentsNumber
		}
		return gk.NewInt(frame, fn()), nil
	}
}

// int63n returns a function that applies a provided transformation to a 63-bit signed integer argument.
// The function takes a single int64 argument, applies the provided function to it, and returns the result as an IObject.
// An error is returned if the number of arguments is different from one.
func (z *Rand) int63n(fn func(int64) int64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 1 {
			return nil, objects.ErrInvalidArgumentsNumber
		}
		i1, err := gk.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return gk.NewInt(frame, fn(i1)), nil
	}
}

// funcInOf64 wraps a no-argument function returning float64 into a FuncCallable to integrate with the GateAdapter system.
// It enforces no arguments and converts the float64 result to an IObject, returning ErrInvalidArgumentsNumber for invalid input.
func (z *Rand) funcInOf64(fn func() float64) objects.FuncCallable {
	return func(gk objects.IGateKeeper, frame int, args ...objects.IObject) (ret objects.IObject, err error) {
		if len(args) != 0 {
			return nil, objects.ErrInvalidArgumentsNumber
		}
		return gk.NewFloat(frame, fn()), nil
	}
}

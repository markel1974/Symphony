package objects

import (
	"time"
)

const (
	maxDepth     = 256
	maxBytesLen  = 2147483647
	maxMapLen    = 100000
	maxArrayLen  = 100000
	maxStructLen = 100000
)

// GateKeeper is a type responsible for creating and managing IObject instances, including primitive and complex types.
// It provides pre-instantiated objects for `true`, `false`, and `undefined` values for efficient reuse.
// The GateKeeper may also include object pooling for specific types to optimize memory usage and performance.
type GateKeeper struct {
	trueValue      IObject
	falseValue     IObject
	undefinedValue IObject

	// Aggiungiamo i pool per gli oggetti
	//intPool   sync.Pool
	//floatPool sync.Pool
	//charPool  sync.Pool
}

const (
	FrameStatic = -1
)

// NewFactory initializes a new GateKeeper instance and sets up default bool and undefined values.
func NewFactory() *GateKeeper {
	f := &GateKeeper{}
	f.trueValue = newBool(f, FrameStatic, true)
	f.falseValue = newBool(f, FrameStatic, false)
	f.undefinedValue = newUndefined(f, FrameStatic)

	//f.intPool.New = func() interface{} {
	//	return &_newInt(f, 0) // Crea un Int con valore di default
	//}
	return f
}

// FalseValue returns the false representation as an IObject from the GateKeeper instance.
func (f *GateKeeper) FalseValue() IObject {
	return f.falseValue
}

// TrueValue returns the IObject instance representing the true value from the GateKeeper.
func (f *GateKeeper) TrueValue() IObject {
	return f.trueValue
}

// UndefinedValue returns the undefined value of the GateKeeper as an IObject.
func (f *GateKeeper) UndefinedValue() IObject {
	return f.undefinedValue
}

// NewObject creates and returns a new instance of Object with its factory field set to the receiving GateKeeper instance.
func (f *GateKeeper) NewObject(frame int) *Object {
	return newObject(f, frame)
}

// NewArray creates and returns a new Array populated with the provided slice of IObject elements.
func (f *GateKeeper) NewArray(frame int, values []IObject) *Array {
	return newArray(f, frame, values)
}

// NewArrayImmutable constructs a new ArrayImmutable instance with the provided slice of IObject, ensuring immutability.
func (f *GateKeeper) NewArrayImmutable(frame int, values []IObject) *ArrayImmutable {
	return newArrayImmutable(f, frame, values)
}

// NewArrayIterator creates a new ArrayIterator for iterating over the provided slice of IObject values.
func (f *GateKeeper) NewArrayIterator(frame int, values []IObject) *ArrayIterator {
	return newArrayIterator(f, frame, values)
}

// NewBool creates and returns a new Bool object initialized with the specified boolean value.
func (f *GateKeeper) NewBool(frame int, value bool) *Bool {
	return newBool(f, frame, value)
}

// NewBuiltin creates a new Builtin object with the specified name and index using the GateKeeper.
func (f *GateKeeper) NewBuiltin(frame int, name string, index int) *Builtin {
	return newBuiltin(f, frame, name, index)
}

// NewBytes creates and returns a new instance of Bytes initialized with the provided byte slice and factory context.
func (f *GateKeeper) NewBytes(frame int, value []byte) *Bytes {
	return newBytes(f, frame, value)
}

// NewBytesIterator creates a new BytesIterator for iterating over the provided byte slice `v` using the specified GateKeeper.
func (f *GateKeeper) NewBytesIterator(frame int, v []byte) *BytesIterator {
	return newBytesIterator(f, frame, v)
}

// NewChar creates a new Char instance associated with the GateKeeper, initialized with the given rune value.
func (f *GateKeeper) NewChar(frame int, value rune) *Char {
	return newChar(f, frame, value)
}

// NewError creates and returns a new Error instance based on the provided IObject value and the associated GateKeeper.
func (f *GateKeeper) NewError(frame int, e string) *Error {
	return newError(f, frame, e)
}

// NewFuncCompiled creates and returns a new FuncCompiled instance using the provided function metadata and bytecode.
func (f *GateKeeper) NewFuncCompiled(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) *FuncCompiled {
	return newFuncCompiled(f, frame, name, instructions, numLocals, numParameters, varArgs, sourceMap, free)
}

// NewFuncPackage creates a new instance of FuncPackage with the specified kind, name, and callable function.
func (f *GateKeeper) NewFuncPackage(kind string, name string, fn FuncCallable) *FuncPackage {
	return newFuncPackage(f, FrameStatic, kind, name, fn)
}

// NewFuncFramePackage creates a new instance of FuncPackage with the specified kind, name, and callable function.
func (f *GateKeeper) NewFuncFramePackage(frame int, kind string, name string, fn FuncCallable) *FuncPackage {
	return newFuncPackage(f, frame, kind, name, fn)
}

// NewFloat creates a new Float instance with the given float64 value, using the GateKeeper for initialization.
func (f *GateKeeper) NewFloat(frame int, v float64) *Float {
	return newFloat(f, frame, v)
}

// NewInt creates and returns a new instance of Int initialized with the given int64 value.
func (f *GateKeeper) NewInt(frame int, v int64) *Int {
	//obj := f.intPool.Get().(*Int)
	//obj.value = v
	//return obj
	return newInt(f, frame, v)
}

//func (f *GateKeeper) ReleaseInt(obj *Int) {
//	// It's good practice to reset the object's state before putting it back in the pool
//	obj.value = 0
//	f.intPool.Put(obj)
//}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.
func (f *GateKeeper) NewObjectPointer(frame int, value *IObject) *ObjectPointer {
	return newObjectPointer(f, frame, value)
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys and IObject values.
func (f *GateKeeper) NewMap(frame int, v map[string]IObject) *Map {
	return newMap(f, frame, v)
}

// NewMapImmutable creates a new immutable map with string keys and IObject values from the provided map.
func (f *GateKeeper) NewMapImmutable(frame int, v map[string]IObject) *MapImmutable {
	return newMapImmutable(f, frame, v)
}

// NewMapIterator creates and returns a new MapIterator for the provided map of string keys and IObject values.
func (f *GateKeeper) NewMapIterator(frame int, v map[string]IObject) *MapIterator {
	return newMapIterator(f, frame, v)
}

// NewString creates a new instance of String with the given value, utilizing the GateKeeper for initialization.
func (f *GateKeeper) NewString(frame int, value string) *String {
	return newString(f, frame, value)
}

// NewStringIterator creates a new StringIterator instance for a given slice of runes, enabling character traversal.
func (f *GateKeeper) NewStringIterator(frame int, v []rune) *StringIterator {
	return newStringIterator(f, frame, v)
}

// NewStruct creates and returns a new instance of Struct using the provided map of string keys and IObject values.
func (f *GateKeeper) NewStruct(frame int, value map[string]IObject) *Struct {
	return newStruct(f, frame, value)
}

// NewStructIterator creates a new StructIterator instance for iterating over a map with string keys and IObject values.
func (f *GateKeeper) NewStructIterator(frame int, v map[string]IObject) *StructIterator {
	return newStructIterator(f, frame, v)
}

// NewTime creates a new instance of Time using the provided time.Time value and initializes it with the factory instance.
func (f *GateKeeper) NewTime(frame int, value time.Time) *Time {
	return newTime(f, frame, value)
}

func (f *GateKeeper) NewUndefined(frame int) *Undefined {
	return newUndefined(f, frame)
}

// ToInterface converts an IObject to its corresponding native Go representation, such as int, string, float64, bool, etc.
func (f *GateKeeper) ToInterface(in IObject) interface{} {
	return toInterface(f, in)
}

// FromInterface converts a native Go value of various types into a corresponding IObject implementation.
func (f *GateKeeper) FromInterface(frame int, in interface{}) IObject {
	return fromInterface(f, frame, in)
}

// ToMap converts an IObject to a map[string]interface{} if the object is a *Map, recursively applying ToInterface.
func (f *GateKeeper) ToMap(o IObject) map[string]interface{} {
	return toMap(f, o)
}

// FromMap converts a map with string keys and interface{} values into a map with string keys and IObject values.
func (f *GateKeeper) FromMap(frame int, v map[string]interface{}) map[string]IObject {
	return fromMap(f, frame, v)
}

// ToInt64 attempts to convert the given IObject to an int64 value.
// It returns the converted value and a boolean indicating success or failure.
func (f *GateKeeper) ToInt64(o IObject) (int64, bool) {
	return toInt64(f, o)
}

// ToInt64Arg converts an IObject to an int64, returning an error if the conversion is not possible or the type is invalid.
func (f *GateKeeper) ToInt64Arg(index int, o IObject) (int64, error) {
	return toInt64Arg(f, index, o)
}

// ToRune attempts to convert an IObject to a rune if it is of type *Int or *Char, returning the rune and a boolean success flag.
func (f *GateKeeper) ToRune(o IObject) (rune, bool) {
	return toRune(o)
}

// ToString converts an IObject to its string representation and determines whether the conversion is valid.
func (f *GateKeeper) ToString(o IObject) (string, bool) {
	return toString(f, o)
}

// ToStringArg attempts to convert an IObject to a string. Returns an error if conversion fails or type is incompatible.
func (f *GateKeeper) ToStringArg(index int, o IObject) (string, error) {
	return toStringArg(f, index, o)
}

// ToStringArrayArg attempts to convert an array of IObjects to a slice of strings.
func (f *GateKeeper) ToStringArrayArg(index int, arr []IObject) ([]string, error) {
	return toStringArrayArg(f, index, arr)
}

// ToByteSlice converts an IObject to a byte slice if the object is of type *Bytes or *String.
// It returns the converted byte slice and a boolean indicating success.
func (f *GateKeeper) ToByteSlice(o IObject) ([]byte, bool) {
	return toByteSlice(o)
}

// ToByteSliceArg attempts to convert an IObject to a byte slice. Returns an error if the conversion fails or the type is incompatible.
func (f *GateKeeper) ToByteSliceArg(index int, o IObject) ([]byte, error) {
	return toByteSliceArg(f, index, o)
}

// ToFloat64 attempts to convert an IObject to a float64 and returns the values along with a success flag.
func (f *GateKeeper) ToFloat64(o IObject) (float64, bool) {
	return toFloat64(f, o)
}

// ToFloat64Arg converts an IObject to a float64 and returns an error if the conversion fails or the type is incompatible.
func (f *GateKeeper) ToFloat64Arg(index int, o IObject) (float64, error) {
	return toFloat64Arg(f, index, o)
}

// ToTime converts an IObject into a time.Time if it is time-compatible (e.g., *Time or *Int). Returns the time and a boolean.
func (f *GateKeeper) ToTime(o IObject) (time.Time, bool) {
	return toTime(o)
}

// ToTimeArg attempts to convert an IObject to a time.Time. Returns an error if the conversion fails or the type is incompatible.
func (f *GateKeeper) ToTimeArg(index int, o IObject) (time.Time, error) {
	return toTimeArg(index, o)
}

// ToBool converts the given IObject to a bool based on its Boolean() method and returns the result along with a success flag.
func (f *GateKeeper) ToBool(o IObject) (v bool, ok bool) {
	return toBool(o)
}

// FromBool converts a boolean values into its corresponding IObject representation, returning TrueValue or FalseValue.
func (f *GateKeeper) FromBool(v bool) IObject {
	return fromBool(f, v)
}

// ToBoolArg converts the given IObject to a boolean if possible or returns an error indicating an invalid argument type.
func (f *GateKeeper) ToBoolArg(index int, o IObject) (bool, error) {
	return toBoolArg(index, o)
}

// FromStringArray converts a slice of strings into an array of IObjects.
func (f *GateKeeper) FromStringArray(frame int, in []string) (IObject, error) {
	return fromStringArray(f, frame, in)
}

// FuncInOn converts a no-argument, no-return Go function into a FuncCallable type that can be called with zero arguments.
// Returns ErrWrongNumArguments if any arguments are passed.
// Invokes the provided function and returns UndefinedValue upon successful execution.
func (f *GateKeeper) FuncInOn(fn func()) FuncCallable {
	return funcInOn(f, fn)
}

// FuncInOi wraps a no-argument integer-returning function into a callable functional interface of type FuncCallable.
// Returns an error if arguments are provided. Converts the integer result into an IObject using NewInt.
func (f *GateKeeper) FuncInOi(fn func() int) FuncCallable {
	return funcInOi(f, fn)
}

// FuncInOi64 wraps a function returning int64 into a FuncCallable with no arguments.
// Returns ErrWrongNumArguments if arguments are passed.
// Converts the result to an IObject using NewInt before returning.
func (f *GateKeeper) FuncInOi64(fn func() int64) FuncCallable {
	return funcInOi64(f, fn)
}

// FuncIi64Oi64 wraps a function that takes int64 and returns int64, into a FuncCallable compatible with IObject interface.
func (f *GateKeeper) FuncIi64Oi64(fn func(int64) int64) FuncCallable {
	return funcIi64Oi64(f, fn)
}

// FuncIi64On wraps a function that accepts a single int64 argument into a FuncCallable that works with IObject arguments.
func (f *GateKeeper) FuncIi64On(fn func(int64)) FuncCallable {
	return funcIi64On(f, fn)
}

// FuncInOb wraps a zero-argument boolean function into a FuncCallable that returns TrueValue or FalseValue.
func (f *GateKeeper) FuncInOb(fn func() bool) FuncCallable {
	return funcInOb(f, fn)
}

// FuncInOe creates a FuncCallable wrapper around a zero-argument function that returns an error.
// Returns ErrWrongNumArguments if arguments are provided.
// Wraps the error returned by the given function into an IObject-compatible error object.
func (f *GateKeeper) FuncInOe(fn func() error) FuncCallable {
	return funcInOe(f, fn)
}

// FuncInOs wraps a function that returns a string, creating a FuncCallable with IObject arguments and results.
// If called with arguments, it returns ErrWrongNumArguments. Otherwise, it returns a string-wrapped IObject result.
func (f *GateKeeper) FuncInOs(fn func() string) FuncCallable {
	return funcInOs(f, fn)
}

// FuncInOse wraps a function that returns a string and error into a FuncCallable that accepts no arguments.
// Returns an error if arguments are provided or if the wrapped function encounters an error.
func (f *GateKeeper) FuncInOse(fn func() (string, error)) FuncCallable {
	return funcInOse(f, fn)
}

// FuncInObSe converts a function returning ([]byte, error) into a FuncCallable that adheres to IObject function standards.
// It ensures the argument count is zero, wraps errors into IObject-compatible errors, and enforces byte slice size limits.
func (f *GateKeeper) FuncInObSe(fn func() ([]byte, error)) FuncCallable {
	return funcInObSe(f, fn)
}

// FuncInOf64 wraps a zero-argument function that returns a float64 into a FuncCallable returning an IObject and an error.
// Returns ErrWrongNumArguments if called with arguments.
// Converts the float64 output of the provided function into an IObject using NewFloat.
func (f *GateKeeper) FuncInOf64(fn func() float64) FuncCallable {
	return funcInOf64(f, fn)
}

// FuncInOsS takes a function that returns a slice of strings and wraps it into a FuncCallable returning an Array of strings.
// The FuncCallable expects zero arguments; passing others results in ErrWrongNumArguments.
// Converts each string from the slice into a String object and appends it to the Array.
func (f *GateKeeper) FuncInOsS(fn func() []string) FuncCallable {
	return funcInOsS(f, fn)
}

// FuncInOiSe wraps a function that returns a slice of integers and an error into a FuncCallable compatible function.
// It validates zero arguments, invokes the wrapped function, wraps any error, and converts the slice to an array of IObject.
// Returns an IObject array containing the integers or a wrapped error if the wrapped function fails.
func (f *GateKeeper) FuncInOiSe(fn func() ([]int, error)) FuncCallable {
	return funcInOiSe(f, fn)
}

// FuncIiOiS takes a function that converts an integer to a slice of integers and returns it as a callable function.
func (f *GateKeeper) FuncIiOiS(fn func(int) []int) FuncCallable {
	return funcIiOiS(f, fn)
}

// FuncIf64Of64 converts a single-argument float64 function into a FuncCallable compatible with IObject arguments.
// It validates the input argument as a float-compatible type.
// Returns a new IObject representing the result or an appropriate error if validation fails.
func (f *GateKeeper) FuncIf64Of64(fn func(float64) float64) FuncCallable {
	return funcIf64Of64(f, fn)
}

// FuncIiOn wraps a function with an int parameter to conform to the FuncCallable signature for custom runtime calls.
// It validates the argument count and type, invoking the provided function with the argument as an integer.
// Returns UndefinedValue on success or an error if the argument is invalid.
func (f *GateKeeper) FuncIiOn(fn func(int)) FuncCallable {
	return funcIiOn(f, fn)
}

// FuncIiOf64 wraps a function of type func(int) float64 as a FuncCallable, enabling its use within the IObject interface ecosystem.
// It validates that exactly one argument is provided and converts it to an int before calling the wrapped function.
// If the argument type is incompatible or the wrong number of arguments are passed, an appropriate error is returned.
func (f *GateKeeper) FuncIiOf64(fn func(int) float64) FuncCallable {
	return funcIiOf64(f, fn)
}

// FuncIf64Oi wraps a function transforming a float64 to an int, making it callable with IObject arguments.
// Returns an error if incorrect number or type of arguments are provided.
func (f *GateKeeper) FuncIf64Oi(fn func(float64) int) FuncCallable {
	return funcIf64Oi(f, fn)
}

// FuncIf64f64Of64 creates a FuncCallable that applies the given binary float64 function to two converted IObject arguments.
// Returns an error if arguments are not exactly two or cannot be converted to float64.
func (f *GateKeeper) FuncIf64f64Of64(fn func(float64, float64) float64) FuncCallable {
	return funcIf64f64Of64(f, fn)
}

// FuncIif64Of64 wraps a provided function accepting an int and float64, returning it as a FuncCallable compatible with IObject arguments.
// It enforces argument type validation and handles potential type mismatches with descriptive errors.
func (f *GateKeeper) FuncIif64Of64(fn func(int, float64) float64) FuncCallable {
	return funcIif64Of64(f, fn)
}

// FuncIf64iOf64 creates a FuncCallable wrapping a function that takes a float64 and int and returns a float64.
// It validates input argument types and converts them to the appropriate types expected by the wrapped function.
// Returns an IObject representing the result of the wrapped function or an error if argument validation fails.
func (f *GateKeeper) FuncIf64iOf64(fn func(float64, int) float64) FuncCallable {
	return funcIf64iOf64(f, fn)
}

// FuncIf64iOb wraps a function that processes a float64 and an int, exposing it as a FuncCallable compatible with the IObject interface.
// It converts the first argument to a float64 and the second to an int, then applies the provided function.
// Returns TrueValue if the function evaluates to true; otherwise, returns FalseValue.
// Returns ErrWrongNumArguments if the argument count is not 2 or NewInvalidArgumentError on type conversion failures.
func (f *GateKeeper) FuncIf64iOb(fn func(float64, int) bool) FuncCallable {
	return funcIf64iOb(f, fn)
}

// FuncIf64Ob wraps a function accepting a float64 and returning a boolean into a FuncCallable compatible with the IObject interface.
func (f *GateKeeper) FuncIf64Ob(fn func(float64) bool) FuncCallable {
	return funcIf64Ob(f, fn)
}

// FuncIsOs creates a FuncCallable that applies a provided string-to-string function to the first argument and returns the result.
func (f *GateKeeper) FuncIsOs(fn func(string) string) FuncCallable {
	return funcIsOs(f, fn)
}

// FuncIsOsS converts a string-to-string-array function into a FuncCallable that operates on IObject arguments.
// It takes one string-compatible argument, applies the provided function, and returns the result as an Array of strings.
// If argument count or type is invalid, it returns an error.
func (f *GateKeeper) FuncIsOsS(fn func(string) []string) FuncCallable {
	return funcIsOsS(f, fn)
}

// FuncIsOse wraps a string transformation function and adapts it to a FuncCallable with argument validation logic.
func (f *GateKeeper) FuncIsOse(fn func(string) (string, error)) FuncCallable {
	return funcIsOse(f, fn)
}

// FuncIsOe converts a string-to-error function into a FuncCallable that operates on IObject arguments.
// It expects exactly one argument convertible to a string and returns an IObject error or result.
// Returns ErrWrongNumArguments if called with an incorrect number of arguments.
// Returns an invalid argument error if the first argument is not string-compatible.
func (f *GateKeeper) FuncIsOe(fn func(string) error) FuncCallable {
	return funcIsOe(f, fn)
}

// FuncIssOe wraps a function accepting two strings and returning an error into a FuncCallable compatible with the IObject interface.
// It ensures the function is called with exactly two string arguments and returns an appropriate error for incorrect usage.
func (f *GateKeeper) FuncIssOe(fn func(string, string) error) FuncCallable {
	return funcIssOe(f, fn)
}

// FuncIssOsS converts a function that takes two strings and returns a slice of strings into a FuncCallable.
// The returned FuncCallable validates its arguments, invokes the provided function, and returns the results as an array.
func (f *GateKeeper) FuncIssOsS(fn func(string, string) []string) FuncCallable {
	return funcIssOsS(f, fn)
}

// FuncIssiOsS converts a function with parameters (string, string, int) -> []string into a FuncCallable.
// It validates arguments, applies the function, and wraps the output in an IObject-compatible Array.
// Returns an error if argument validation fails or function results cannot be converted to a String.
func (f *GateKeeper) FuncIssiOsS(fn func(string, string, int) []string) FuncCallable {
	return funcIssiOsS(f, fn)
}

// FuncIssOi converts a function with two string inputs and an int output into a FuncCallable type.
// The returned FuncCallable validates that exactly two arguments are passed and they are string-compatible.
// If arguments are valid, the wrapped function is invoked, and its integer result is wrapped in an IObject.
// Returns an error if the number of arguments is incorrect or conversion to strings fails.
func (f *GateKeeper) FuncIssOi(fn func(string, string) int) FuncCallable {
	return funcIssOi(f, fn)
}

// FuncIssOs wraps a function that takes two strings and returns a string into a FuncCallable accepting IObject arguments.
// It validates argument types and ensures exactly two arguments are passed or returns an appropriate error.
// The wrapped function's result is converted to an IObject before being returned.
func (f *GateKeeper) FuncIssOs(fn func(string, string) string) FuncCallable {
	return funcIssOs(f, fn)
}

// FuncIssOb wraps a binary comparison function for strings as a callable function in the IObject system.
// The returned FuncCallable validates arguments, applies the provided function, and returns TrueValue or FalseValue.
// It expects the function to take two string arguments and return a boolean indicating the comparison result.
// Returns an error if the number of arguments is incorrect or arguments are not string-compatible.
func (f *GateKeeper) FuncIssOb(fn func(string, string) bool) FuncCallable {
	return funcIssOb(f, fn)
}

// FuncIsSsOs creates a FuncCallable that processes a string slice and a string, applying the given transformation function.
func (f *GateKeeper) FuncIsSsOs(fn func([]string, string) string) FuncCallable {
	return funcIsSsOs(f, fn)
}

// FuncIsi64Oe transforms a function accepting a string and int64 into a FuncCallable that operates on IObject arguments.
// Takes exactly two arguments; the first must be string-compatible, the second int64-compatible, or errors are returned.
// Wraps the result of the provided function into an IObject or returns an appropriate error if validation fails.
func (f *GateKeeper) FuncIsi64Oe(fn func(string, int64) error) FuncCallable {
	return funcIsi64Oe(f, fn)
}

// FuncIiiOe wraps a function taking two integers and returning an error into a FuncCallable accepting two IObject arguments.
func (f *GateKeeper) FuncIiiOe(fn func(int, int) error) FuncCallable {
	return funcIiiOe(f, fn)
}

// FuncIsiOs wraps a function that takes a string and int as inputs and returns a string, converting it to a FuncCallable.
// It validates the arguments, calls the wrapped function, and converts the result to an IObject.
// Returns an error if arguments are of invalid types or wrong number of arguments is supplied.
func (f *GateKeeper) FuncIsiOs(fn func(string, int) string) FuncCallable {
	return funcIsiOs(f, fn)
}

// FuncIsiiOe converts a function with string, int, int inputs, and an error return into a FuncCallable with variadic IObject arguments.
func (f *GateKeeper) FuncIsiiOe(fn func(string, int, int) error) FuncCallable {
	return funcIsiiOe(f, fn)
}

// FuncIbSOie wraps a function that takes a byte slice and returns an int and error into a FuncCallable for IObject use.
// It ensures the input argument is a single byte-compatible IObject and converts its result to IObject format.
// Returns ErrWrongNumArguments if called with more or less than one argument.
// Returns NewInvalidArgumentError if the input argument isn't byte-compatible.
// Converts the function's error output into an appropriate IObject error.
func (f *GateKeeper) FuncIbSOie(fn func([]byte) (int, error)) FuncCallable {
	return funcIbSOie(f, fn)
}

// FuncIbSOs wraps a function that converts a byte slice to a string, returning it as a FuncCallable in the custom object system.
// It ensures the input is a single argument of type bytes-compatible, and returns an error for invalid or unsupported types.
// The resulting FuncCallable checks argument validity, applies the provided function, and returns a new String object.
func (f *GateKeeper) FuncIbSOs(fn func([]byte) string) FuncCallable {
	return funcIbSOs(f, fn)
}

// FuncIsOie wraps a string-to-int function into a FuncCallable compatible with IObject interface arguments and error handling.
func (f *GateKeeper) FuncIsOie(fn func(string) (int, error)) FuncCallable {
	return funcIsOie(f, fn)
}

// FuncIsObSe returns a FuncCallable that wraps a function converting a string to a byte slice and error output.
// It validates input, reports invalid arguments, enforces byte length limits, and converts output to IObject format.
// Uses ErrWrongNumArguments, NewInvalidArgumentError, and ErrBytesLimit for error handling.
func (f *GateKeeper) FuncIsObSe(fn func(string) ([]byte, error)) FuncCallable {
	return funcIsObSe(f, fn)
}

// FuncIiOsSe converts a function mapping an integer to a slice of strings and an error into a FuncCallable.
func (f *GateKeeper) FuncIiOsSe(fn func(int) ([]string, error)) FuncCallable {
	return funcIiOsSe(f, fn)
}

// FuncIiOs wraps a function of type `func(int) string` into a FuncCallable compatible with the IObject interface system.
// It validates argument count and type, executes the provided function, and converts the result into an IObject.
func (f *GateKeeper) FuncIiOs(fn func(int) string) FuncCallable {
	return funcIiOs(f, fn)
}

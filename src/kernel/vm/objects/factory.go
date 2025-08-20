package objects

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Code		Meaning		Go-Type			Description
//	I		Input			-			Prefix indicating input parameters.
//	O		Output			-			Prefix indicating return values.
//	n		None			()			Indicates absence of parameters or return values.
//	i		int				int			An integer.
//	b		bool			bool		A boolean value.
//	e		error			error		An error object.
//	i64		int64			int64		A 64-bit integer.
//	f64		float64			float64		A 64-bit floating point number.
//	s		string			string		A string.
//	sS		string-Slice 	[]string	A slice (array) of strings.
//	bS		bytes-Slice		[]byte		A slice (array) of bytes.
//  iS      int-Slice		[]int       A slice (array) of int.

// Factory is a type responsible for creating and managing IObject instances, including primitive and complex types.
// It provides pre-instantiated objects for `true`, `false`, and `undefined` values for efficient reuse.
// The Factory may also include object pooling for specific types to optimize memory usage and performance.
type Factory struct {
	trueValue      IObject
	falseValue     IObject
	undefinedValue IObject

	// Aggiungiamo i pool per gli oggetti
	//intPool   sync.Pool
	//floatPool sync.Pool
	//charPool  sync.Pool
}

const (
	FrameStatic      = -1
	FrameReturnValue = -2
)

// NewFactory initializes a new Factory instance and sets up default bool and undefined values.
func NewFactory() *Factory {
	f := &Factory{}
	f.trueValue = _newBool(f, FrameStatic, true)
	f.falseValue = _newBool(f, FrameStatic, false)
	f.undefinedValue = _newUndefined(f, FrameStatic)

	//f.intPool.New = func() interface{} {
	//	return &_newInt(f, 0) // Crea un Int con valore di default
	//}
	return f
}

// FalseValue returns the false representation as an IObject from the Factory instance.
func (f *Factory) FalseValue() IObject {
	return f.falseValue
}

// TrueValue returns the IObject instance representing the true value from the Factory.
func (f *Factory) TrueValue() IObject {
	return f.trueValue
}

// UndefinedValue returns the undefined value of the Factory as an IObject.
func (f *Factory) UndefinedValue() IObject {
	return f.undefinedValue
}

// NewObject creates and returns a new instance of Object with its factory field set to the receiving Factory instance.
func (f *Factory) NewObject(frame int) *Object {
	return _newObject(f, frame)
}

// NewArray creates and returns a new Array populated with the provided slice of IObject elements.
func (f *Factory) NewArray(frame int, values []IObject) *Array {
	return _newArray(f, frame, values)
}

// NewArrayImmutable constructs a new ArrayImmutable instance with the provided slice of IObject, ensuring immutability.
func (f *Factory) NewArrayImmutable(frame int, values []IObject) *ArrayImmutable {
	return _newArrayImmutable(f, frame, values)
}

// NewArrayIterator creates a new ArrayIterator for iterating over the provided slice of IObject values.
func (f *Factory) NewArrayIterator(frame int, values []IObject) *ArrayIterator {
	return _newArrayIterator(f, frame, values)
}

// NewBool creates and returns a new Bool object initialized with the specified boolean value.
func (f *Factory) NewBool(frame int, value bool) *Bool {
	return _newBool(f, frame, value)
}

// NewBuiltin creates a new Builtin object with the specified name and index using the Factory.
func (f *Factory) NewBuiltin(frame int, name string, index int) *Builtin {
	return _newBuiltin(f, frame, name, index)
}

// NewBytes creates and returns a new instance of Bytes initialized with the provided byte slice and factory context.
func (f *Factory) NewBytes(frame int, value []byte) *Bytes {
	return _newBytes(f, frame, value)
}

// NewBytesIterator creates a new BytesIterator for iterating over the provided byte slice `v` using the specified Factory.
func (f *Factory) NewBytesIterator(frame int, v []byte) *BytesIterator {
	return _newBytesIterator(f, frame, v)
}

// NewChar creates a new Char instance associated with the Factory, initialized with the given rune value.
func (f *Factory) NewChar(frame int, value rune) *Char {
	return _newChar(f, frame, value)
}

// NewError creates and returns a new Error instance based on the provided IObject value and the associated Factory.
func (f *Factory) NewError(frame int, e string) *Error {
	return _newError(f, frame, e)
}

// NewFuncCompiled creates and returns a new FuncCompiled instance using the provided function metadata and bytecode.
func (f *Factory) NewFuncCompiled(frame int, name string, instructions []byte, numLocals int, numParameters int, varArgs bool, sourceMap map[int]int, free []*ObjectPointer) *FuncCompiled {
	return _newFuncCompiled(f, frame, name, instructions, numLocals, numParameters, varArgs, sourceMap, free)
}

// NewFuncPackage creates a new instance of FuncPackage with the specified kind, name, and callable function.
func (f *Factory) NewFuncPackage(kind string, name string, fn FuncCallable) *FuncPackage {
	return _newFuncPackage(f, kind, name, fn)
}

// NewFloat creates a new Float instance with the given float64 value, using the Factory for initialization.
func (f *Factory) NewFloat(frame int, v float64) *Float {
	return _newFloat(f, frame, v)
}

// NewInt creates and returns a new instance of Int initialized with the given int64 value.
func (f *Factory) NewInt(frame int, v int64) *Int {
	//obj := f.intPool.Get().(*Int)
	//obj.value = v
	//return obj
	return _newInt(f, frame, v)
}

//func (f *Factory) ReleaseInt(obj *Int) {
//	// It's good practice to reset the object's state before putting it back in the pool
//	obj.value = 0
//	f.intPool.Put(obj)
//}

// NewObjectPointer creates a new ObjectPointer instance wrapping the provided IObject pointer.
func (f *Factory) NewObjectPointer(frame int, value *IObject) *ObjectPointer {
	return _newObjectPointer(f, frame, value)
}

// NewMap creates and returns a new instance of Map initialized with the provided map of string keys and IObject values.
func (f *Factory) NewMap(frame int, v map[string]IObject) *Map {
	return _newMap(f, frame, v)
}

// NewMapImmutable creates a new immutable map with string keys and IObject values from the provided map.
func (f *Factory) NewMapImmutable(frame int, v map[string]IObject) *MapImmutable {
	return _newMapImmutable(f, frame, v)
}

// NewMapIterator creates and returns a new MapIterator for the provided map of string keys and IObject values.
func (f *Factory) NewMapIterator(frame int, v map[string]IObject) *MapIterator {
	return _newMapIterator(f, frame, v)
}

// NewStringNoSize creates a new String instance with the provided value, omitting size initialization.
func (f *Factory) NewStringNoSize(frame int, value string) *String {
	return _newStringNoSize(f, frame, value)
}

// NewString creates a new instance of String with the given value, utilizing the Factory for initialization.
func (f *Factory) NewString(frame int, value string) (*String, error) {
	return _newString(f, frame, value)
}

// NewStringIterator creates a new StringIterator instance for a given slice of runes, enabling character traversal.
func (f *Factory) NewStringIterator(frame int, v []rune) *StringIterator {
	return _newStringIterator(f, frame, v)
}

// NewStruct creates and returns a new instance of Struct using the provided map of string keys and IObject values.
func (f *Factory) NewStruct(frame int, value map[string]IObject) *Struct {
	return _newStruct(f, frame, value)
}

// NewStructIterator creates a new StructIterator instance for iterating over a map with string keys and IObject values.
func (f *Factory) NewStructIterator(frame int, v map[string]IObject) *StructIterator {
	return _newStructIterator(f, frame, v)
}

// NewTime creates a new instance of Time using the provided time.Time value and initializes it with the factory instance.
func (f *Factory) NewTime(frame int, value time.Time) *Time {
	return _newTime(f, frame, value)
}

// ToInterface converts an IObject to its corresponding native Go representation, such as int, string, float64, bool, etc.
func (f *Factory) ToInterface(in IObject) (res interface{}) {
	switch o := in.(type) {
	case *Int:
		res = o.value
	case *String:
		res = o.value
	case *Float:
		res = o.value
	case *Bool:
		res = o == f.TrueValue()
	case *Char:
		res = o.value
	case *Bytes:
		res = o.values
	case *Array:
		res = make([]interface{}, len(o.Values()))
		for i, val := range o.Values() {
			res.([]interface{})[i] = f.ToInterface(val)
		}
	case *ArrayImmutable:
		res = make([]interface{}, o.Length())
		for i, val := range o.Values() {
			res.([]interface{})[i] = f.ToInterface(val)
		}
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res.(map[string]interface{})[key] = f.ToInterface(v)
		}
	case *MapImmutable:
		res = make(map[string]interface{})
		for key, v := range o.Values() {
			res.(map[string]interface{})[key] = f.ToInterface(v)
		}
	case *Time:
		res = o.value
	case *Error:
		res = errors.New(o.String())
	case *Undefined:
		res = nil
	case IObject:
		return o
	}
	return
}

// FromInterface converts a native Go value of various types into a corresponding IObject implementation.
func (f *Factory) FromInterface(frame int, in interface{}) IObject {
	switch v := in.(type) {
	case nil:
		return f.UndefinedValue()
	case string:
		if len(v) > MaxStringLen {
			return f.NewStringNoSize(frame, v[0:MaxStringLen])
		}
		return f.NewStringNoSize(frame, v)
	case int64:
		return f.NewInt(frame, v)
	case int:
		return f.NewInt(frame, int64(v))
	case bool:
		if v {
			return f.TrueValue()
		}
		return f.FalseValue()
	case rune:
		return f.NewChar(frame, v)
	case byte:
		return f.NewChar(frame, rune(v))
	case float64:
		return f.NewFloat(frame, v)
	case []byte:
		if len(v) > MaxBytesLen {
			return f.NewBytes(frame, v[0:MaxBytesLen])
		}
		return f.NewBytes(frame, v)
	case error:
		return f.NewError(frame, v.Error())
	case map[string]IObject:
		return f.NewMap(frame, v)
	case map[string]interface{}:
		kv := f.FromMap(frame, v)
		return f.NewMap(frame, kv)
	case []bool:
		arr := make([]IObject, len(v))
		for i, e := range v {
			if e {
				arr[i] = f.TrueValue()
			} else {
				arr[i] = f.FalseValue()
			}
		}
		return f.NewArray(frame, arr)
	case []int:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = f.NewInt(frame, int64(e))
		}
		return f.NewArray(frame, arr)
	case []map[string]interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			kv := f.FromMap(frame, e)
			vo := f.FromInterface(frame, kv)
			arr[i] = vo
		}
		return f.NewArray(frame, arr)
	case []IObject:
		return f.NewArray(frame, v)
	case []interface{}:
		arr := make([]IObject, len(v))
		for i, e := range v {
			arr[i] = f.FromInterface(frame, e)
		}
		return f.NewArray(frame, arr)
	case time.Time:
		return f.NewTime(frame, v)
	case IObject:
		return v
	case FuncCallable:
		return f.NewFuncPackage(FuncPackageDef, "FuncCallable", v)
	}
	return f.UndefinedValue()
}

// ToMap converts an IObject to a map[string]interface{} if the object is a *Map, recursively applying ToInterface.
func (f *Factory) ToMap(o IObject) (res map[string]interface{}) {
	switch o := o.(type) {
	case *Map:
		res = make(map[string]interface{})
		for key, v := range o.values {
			res[key] = f.ToInterface(v)
		}
	}
	return
}

// FromMap converts a map with string keys and interface{} values into a map with string keys and IObject values.
func (f *Factory) FromMap(frame int, v map[string]interface{}) map[string]IObject {
	kv := make(map[string]IObject)
	for key, val := range v {
		kv[key] = f.FromInterface(frame, val)
	}
	return kv
}

// ToInt64 attempts to convert the given IObject to an int64 value.
// It returns the converted value and a boolean indicating success or failure.
func (f *Factory) ToInt64(o IObject) (int64, bool) {
	switch o := o.(type) {
	case *Int:
		return o.value, true
	case *Float:
		return int64(o.value), true
	case *Char:
		return int64(o.value), true
	case *Bool:
		if o == f.TrueValue() {
			return 1, true
		}
		return 0, true
	case *String:
		c, err := strconv.ParseInt(o.value, 10, 64)
		if err == nil {
			return c, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ToInt64Arg converts an IObject to an int64, returning an error if the conversion is not possible or the type is invalid.
func (f *Factory) ToInt64Arg(index int, o IObject) (int64, error) {
	v, ok := f.ToInt64(o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "int(compatible)", o.TypeName())
	}
	return v, nil
}

// ToRune attempts to convert an IObject to a rune if it is of type *Int or *Char, returning the rune and a boolean success flag.
func (f *Factory) ToRune(o IObject) (v rune, ok bool) {
	switch o := o.(type) {
	case *Int:
		v = rune(o.value)
		ok = true
	case *Char:
		v = o.value
		ok = true
	}
	return
}

// ToString converts an IObject to its string representation and determines whether the conversion is valid.
func (f *Factory) ToString(o IObject) (string, bool) {
	if o == nil {
		return "", false
	}
	if o == f.UndefinedValue() {
		return "", false
	}
	if str, isStr := o.(*String); isStr {
		return str.value, true
	}
	return o.String(), true
}

// ToStringArg attempts to convert an IObject to a string. Returns an error if conversion fails or type is incompatible.
func (f *Factory) ToStringArg(index int, o IObject) (string, error) {
	v, ok := f.ToString(o)
	if !ok {
		return "", NewInvalidArgumentError(index, "string(compatible)", o.TypeName())
	}
	return v, nil
}

// ToStringArrayArg attempts to convert an array of IObjects to a slice of strings.
func (f *Factory) ToStringArrayArg(index int, arr []IObject) ([]string, error) {
	var sArr []string
	for idx, elem := range arr {
		str, ok := f.ToString(elem)
		if !ok {
			return nil, NewInvalidArgumentError(index, fmt.Sprintf("%d - string array(compatible)", idx), elem.TypeName())
		}
		sArr = append(sArr, str)
	}
	return sArr, nil
}

// ToByteSlice converts an IObject to a byte slice if the object is of type *Bytes or *String.
// It returns the converted byte slice and a boolean indicating success.
func (f *Factory) ToByteSlice(o IObject) ([]byte, bool) {
	switch o := o.(type) {
	case *Bytes:
		return o.values, true
	case *String:
		return []byte(o.value), true
	default:
		return nil, false
	}
}

// ToByteSliceArg attempts to convert an IObject to a byte slice. Returns an error if the conversion fails or the type is incompatible.
func (f *Factory) ToByteSliceArg(index int, o IObject) ([]byte, error) {
	b, ok := f.ToByteSlice(o)
	if !ok {
		return nil, NewInvalidArgumentError(index, "byte slice(compatible)", o.TypeName())
	}
	return b, nil
}

// ToFloat64 attempts to convert an IObject to a float64 and returns the values along with a success flag.
func (f *Factory) ToFloat64(o IObject) (float64, bool) {
	switch o := o.(type) {
	case *Int:
		return float64(o.value), true
	case *Float:
		return o.value, true
	case *Char:
		return float64(o.value), true
	case *Bool:
		if o == f.TrueValue() {
			return 1, true
		}
		return 0, true
	case *String:
		c, err := strconv.ParseFloat(o.value, 64)
		if err == nil {
			return c, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ToFloat64Arg converts an IObject to a float64 and returns an error if the conversion fails or the type is incompatible.
func (f *Factory) ToFloat64Arg(index int, o IObject) (float64, error) {
	v, ok := f.ToFloat64(o)
	if !ok {
		return 0, NewInvalidArgumentError(index, "float64(compatible)", o.TypeName())
	}
	return v, nil
}

// ToTime converts an IObject into a time.Time if it is time-compatible (e.g., *Time or *Int). Returns the time and a boolean.
func (f *Factory) ToTime(o IObject) (time.Time, bool) {
	switch o := o.(type) {
	case *Time:
		return o.value, true
	case *Int:
		return time.Unix(o.value, 0), true
	}
	return time.Time{}, false
}

// ToTimeArg attempts to convert an IObject to a time.Time. Returns an error if the conversion fails or the type is incompatible.
func (f *Factory) ToTimeArg(index int, o IObject) (time.Time, error) {
	v, ok := f.ToTime(o)
	if !ok {
		return time.Time{}, NewInvalidArgumentError(index, "time(compatible)", o.TypeName())
	}
	return v, nil
}

// ToBool converts the given IObject to a bool based on its Boolean() method and returns the result along with a success flag.
func (f *Factory) ToBool(o IObject) (v bool, ok bool) {
	ok = true
	v = !o.Boolean()
	return
}

// FromBool converts a boolean values into its corresponding IObject representation, returning TrueValue or FalseValue.
func (f *Factory) FromBool(v bool) IObject {
	if v {
		return f.TrueValue()
	}
	return f.FalseValue()
}

// ToBoolArg converts the given IObject to a boolean if possible or returns an error indicating an invalid argument type.
func (f *Factory) ToBoolArg(index int, o IObject) (bool, error) {
	b1, ok := o.(*Bool)
	if !ok {
		return false, NewInvalidArgumentError(index, "bool(compatible)", o.TypeName())
	}
	return b1.value, nil
}

// FromStringArray converts a slice of strings into an array of IObjects.
func (f *Factory) FromStringArray(frame int, in []string) (IObject, error) {
	var data []IObject
	if len(in) > 0 {
		data = make([]IObject, len(in))
		for idx, v := range in {
			r, err := f.NewString(frame, v)
			if err != nil {
				return nil, err
			}
			data[idx] = r
		}
	}
	return f.NewArray(frame, data), nil
}

// FuncInOn converts a no-argument, no-return Go function into a FuncCallable type that can be called with zero arguments.
// Returns ErrWrongNumArguments if any arguments are passed.
// Invokes the provided function and returns UndefinedValue upon successful execution.
func (f *Factory) FuncInOn(fn func()) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		fn()
		return f.UndefinedValue(), nil
	}
}

// FuncInOi wraps a no-argument integer-returning function into a callable functional interface of type FuncCallable.
// Returns an error if arguments are provided. Converts the integer result into an IObject using NewInt.
func (f *Factory) FuncInOi(fn func() int) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return f.NewInt(FrameReturnValue, int64(fn())), nil
	}
}

// FuncInOi64 wraps a function returning int64 into a FuncCallable with no arguments.
// Returns ErrWrongNumArguments if arguments are passed.
// Converts the result to an IObject using NewInt before returning.
func (f *Factory) FuncInOi64(fn func() int64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return f.NewInt(FrameReturnValue, fn()), nil
	}
}

// FuncIi64Oi64 wraps a function that takes int64 and returns int64, into a FuncCallable compatible with IObject interface.
func (f *Factory) FuncIi64Oi64(fn func(int64) int64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewInt(FrameReturnValue, fn(i1)), nil
	}
}

// FuncIi64On wraps a function that accepts a single int64 argument into a FuncCallable that works with IObject arguments.
func (f *Factory) FuncIi64On(fn func(int64)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(i1)
		return f.UndefinedValue(), nil
	}
}

// FuncInOb wraps a zero-argument boolean function into a FuncCallable that returns TrueValue or FalseValue.
func (f *Factory) FuncInOb(fn func() bool) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		if fn() {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// FuncInOe creates a FuncCallable wrapper around a zero-argument function that returns an error.
// Returns ErrWrongNumArguments if arguments are provided.
// Wraps the error returned by the given function into an IObject-compatible error object.
func (f *Factory) FuncInOe(fn func() error) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		err := fn()
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// FuncInOs wraps a function that returns a string, creating a FuncCallable with IObject arguments and results.
// If called with arguments, it returns ErrWrongNumArguments. Otherwise, it returns a string-wrapped IObject result.
func (f *Factory) FuncInOs(fn func() string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		v, err := f.NewString(FrameReturnValue, fn())
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncInOse wraps a function that returns a string and error into a FuncCallable that accepts no arguments.
// Returns an error if arguments are provided or if the wrapped function encounters an error.
func (f *Factory) FuncInOse(fn func() (string, error)) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		v, err := f.NewString(FrameReturnValue, res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncInObSe converts a function returning ([]byte, error) into a FuncCallable that adheres to IObject function standards.
// It ensures the argument count is zero, wraps errors into IObject-compatible errors, and enforces byte slice size limits.
func (f *Factory) FuncInObSe(fn func() ([]byte, error)) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		if len(res) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return f.NewBytes(FrameReturnValue, res), nil
	}
}

// FuncInOf64 wraps a zero-argument function that returns a float64 into a FuncCallable returning an IObject and an error.
// Returns ErrWrongNumArguments if called with arguments.
// Converts the float64 output of the provided function into an IObject using NewFloat.
func (f *Factory) FuncInOf64(fn func() float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		return f.NewFloat(FrameReturnValue, fn()), nil
	}
}

// FuncInOsS takes a function that returns a slice of strings and wraps it into a FuncCallable returning an Array of strings.
// The FuncCallable expects zero arguments; passing others results in ErrWrongNumArguments.
// Converts each string from the slice into a String object and appends it to the Array.
func (f *Factory) FuncInOsS(fn func() []string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		arr := f.NewArray(FrameReturnValue, nil)
		for _, elem := range fn() {
			v, err := f.NewString(FrameReturnValue, elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncInOiSe wraps a function that returns a slice of integers and an error into a FuncCallable compatible function.
// It validates zero arguments, invokes the wrapped function, wraps any error, and converts the slice to an array of IObject.
// Returns an IObject array containing the integers or a wrapped error if the wrapped function fails.
func (f *Factory) FuncInOiSe(fn func() ([]int, error)) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 0 {
			return nil, ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		arr := f.NewArray(FrameReturnValue, nil)
		for _, v := range res {
			arr.Append(f.NewInt(FrameReturnValue, int64(v)))
		}
		return arr, nil
	}
}

// FuncIiOiS takes a function that converts an integer to a slice of integers and returns it as a callable function.
func (f *Factory) FuncIiOiS(fn func(int) []int) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(int(i1))
		arr := f.NewArray(FrameReturnValue, nil)
		for _, v := range res {
			arr.Append(f.NewInt(FrameReturnValue, int64(v)))
		}
		return arr, nil
	}
}

// FuncIf64Of64 converts a single-argument float64 function into a FuncCallable compatible with IObject arguments.
// It validates the input argument as a float-compatible type.
// Returns a new IObject representing the result or an appropriate error if validation fails.
func (f *Factory) FuncIf64Of64(fn func(float64) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(FrameReturnValue, fn(f1)), nil
	}
}

// FuncIiOn wraps a function with an int parameter to conform to the FuncCallable signature for custom runtime calls.
// It validates the argument count and type, invoking the provided function with the argument as an integer.
// Returns UndefinedValue on success or an error if the argument is invalid.
func (f *Factory) FuncIiOn(fn func(int)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		fn(int(i1))
		return f.UndefinedValue(), nil
	}
}

// FuncIiOf64 wraps a function of type func(int) float64 as a FuncCallable, enabling its use within the IObject interface ecosystem.
// It validates that exactly one argument is provided and converts it to an int before calling the wrapped function.
// If the argument type is incompatible or the wrong number of arguments are passed, an appropriate error is returned.
func (f *Factory) FuncIiOf64(fn func(int) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(FrameReturnValue, fn(int(i1))), nil
	}
}

// FuncIf64Oi wraps a function transforming a float64 to an int, making it callable with IObject arguments.
// Returns an error if incorrect number or type of arguments are provided.
func (f *Factory) FuncIf64Oi(fn func(float64) int) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		return f.NewInt(FrameReturnValue, int64(fn(f1))), nil
	}
}

// FuncIf64f64Of64 creates a FuncCallable that applies the given binary float64 function to two converted IObject arguments.
// Returns an error if arguments are not exactly two or cannot be converted to float64.
func (f *Factory) FuncIf64f64Of64(fn func(float64, float64) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := f.ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(FrameReturnValue, fn(f1, f2)), nil
	}
}

// FuncIif64Of64 wraps a provided function accepting an int and float64, returning it as a FuncCallable compatible with IObject arguments.
// It enforces argument type validation and handles potential type mismatches with descriptive errors.
func (f *Factory) FuncIif64Of64(fn func(int, float64) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		f2, err := f.ToFloat64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(FrameReturnValue, fn(int(i1), f2)), nil
	}
}

// FuncIf64iOf64 creates a FuncCallable wrapping a function that takes a float64 and int and returns a float64.
// It validates input argument types and converts them to the appropriate types expected by the wrapped function.
// Returns an IObject representing the result of the wrapped function or an error if argument validation fails.
func (f *Factory) FuncIf64iOf64(fn func(float64, int) float64) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewFloat(FrameReturnValue, fn(f1, int(i2))), nil
	}
}

// FuncIf64iOb wraps a function that processes a float64 and an int, exposing it as a FuncCallable compatible with the IObject interface.
// It converts the first argument to a float64 and the second to an int, then applies the provided function.
// Returns TrueValue if the function evaluates to true; otherwise, returns FalseValue.
// Returns ErrWrongNumArguments if the argument count is not 2 or NewInvalidArgumentError on type conversion failures.
func (f *Factory) FuncIf64iOb(fn func(float64, int) bool) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(f1, int(i2)) {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// FuncIf64Ob wraps a function accepting a float64 and returning a boolean into a FuncCallable compatible with the IObject interface.
func (f *Factory) FuncIf64Ob(fn func(float64) bool) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		f1, err := f.ToFloat64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		if fn(f1) {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// FuncIsOs creates a FuncCallable that applies a provided string-to-string function to the first argument and returns the result.
func (f *Factory) FuncIsOs(fn func(string) string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(FrameReturnValue, fn(s1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsOsS converts a string-to-string-array function into a FuncCallable that operates on IObject arguments.
// It takes one string-compatible argument, applies the provided function, and returns the result as an Array of strings.
// If argument count or type is invalid, it returns an error.
func (f *Factory) FuncIsOsS(fn func(string) []string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res := fn(s1)
		arr := f.NewArray(FrameReturnValue, nil)
		for _, elem := range res {
			v, err := f.NewString(FrameReturnValue, elem)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIsOse wraps a string transformation function and adapts it to a FuncCallable with argument validation logic.
func (f *Factory) FuncIsOse(fn func(string) (string, error)) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		v, err := f.NewString(FrameReturnValue, res)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsOe converts a string-to-error function into a FuncCallable that operates on IObject arguments.
// It expects exactly one argument convertible to a string and returns an IObject error or result.
// Returns ErrWrongNumArguments if called with an incorrect number of arguments.
// Returns an invalid argument error if the first argument is not string-compatible.
func (f *Factory) FuncIsOe(fn func(string) error) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		err = fn(s1)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// FuncIssOe wraps a function accepting two strings and returning an error into a FuncCallable compatible with the IObject interface.
// It ensures the function is called with exactly two string arguments and returns an appropriate error for incorrect usage.
func (f *Factory) FuncIssOe(fn func(string, string) error) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		err = fn(s1, s2)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// FuncIssOsS converts a function that takes two strings and returns a slice of strings into a FuncCallable.
// The returned FuncCallable validates its arguments, invokes the provided function, and returns the results as an array.
func (f *Factory) FuncIssOsS(fn func(string, string) []string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		arr := f.NewArray(FrameReturnValue, nil)
		for _, res := range fn(s1, s2) {
			v, err := f.NewString(FrameReturnValue, res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIssiOsS converts a function with parameters (string, string, int) -> []string into a FuncCallable.
// It validates arguments, applies the function, and wraps the output in an IObject-compatible Array.
// Returns an error if argument validation fails or function results cannot be converted to a String.
func (f *Factory) FuncIssiOsS(fn func(string, string, int) []string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := f.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		arr := f.NewArray(FrameReturnValue, nil)
		for _, res := range fn(s1, s2, int(i3)) {
			v, err := f.NewString(FrameReturnValue, res)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIssOi converts a function with two string inputs and an int output into a FuncCallable type.
// The returned FuncCallable validates that exactly two arguments are passed and they are string-compatible.
// If arguments are valid, the wrapped function is invoked, and its integer result is wrapped in an IObject.
// Returns an error if the number of arguments is incorrect or conversion to strings fails.
func (f *Factory) FuncIssOi(fn func(string, string) int) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		return f.NewInt(FrameReturnValue, int64(fn(s1, s2))), nil
	}
}

// FuncIssOs wraps a function that takes two strings and returns a string into a FuncCallable accepting IObject arguments.
// It validates argument types and ensures exactly two arguments are passed or returns an appropriate error.
// The wrapped function's result is converted to an IObject before being returned.
func (f *Factory) FuncIssOs(fn func(string, string) string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(FrameReturnValue, fn(s1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil

	}
}

// FuncIssOb wraps a binary comparison function for strings as a callable function in the IObject system.
// The returned FuncCallable validates arguments, applies the provided function, and returns TrueValue or FalseValue.
// It expects the function to take two string arguments and return a boolean indicating the comparison result.
// Returns an error if the number of arguments is incorrect or arguments are not string-compatible.
func (f *Factory) FuncIssOb(fn func(string, string) bool) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		if fn(s1, s2) {
			return f.TrueValue(), nil
		}
		return f.FalseValue(), nil
	}
}

// FuncIsSsOs creates a FuncCallable that processes a string slice and a string, applying the given transformation function.
func (f *Factory) FuncIsSsOs(fn func([]string, string) string) FuncCallable {
	return func(args ...IObject) (IObject, error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		var ss1 []string
		switch arg0 := args[0].(type) {
		case *Array:
			for idx, a := range arg0.Values() {
				as, err := f.ToStringArg(idx, a)
				if err != nil {
					return nil, err
				}
				ss1 = append(ss1, as)
			}
		case *ArrayImmutable:
			for idx, a := range arg0.Values() {
				as, err := f.ToStringArg(idx, a)
				if err != nil {
					return nil, err
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, NewInvalidArgumentError(0, "array", args[0].TypeName())
		}
		s2, err := f.ToStringArg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(FrameReturnValue, fn(ss1, s2))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsi64Oe transforms a function accepting a string and int64 into a FuncCallable that operates on IObject arguments.
// Takes exactly two arguments; the first must be string-compatible, the second int64-compatible, or errors are returned.
// Wraps the result of the provided function into an IObject or returns an appropriate error if validation fails.
func (f *Factory) FuncIsi64Oe(fn func(string, int64) error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		err = fn(s1, i2)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// FuncIiiOe wraps a function taking two integers and returning an error into a FuncCallable accepting two IObject arguments.
func (f *Factory) FuncIiiOe(fn func(int, int) error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		err = fn(int(i1), int(i2))
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// FuncIsiOs wraps a function that takes a string and int as inputs and returns a string, converting it to a FuncCallable.
// It validates the arguments, calls the wrapped function, and converts the result to an IObject.
// Returns an error if arguments are of invalid types or wrong number of arguments is supplied.
func (f *Factory) FuncIsiOs(fn func(string, int) string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 2 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(FrameReturnValue, fn(s1, int(i2)))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsiiOe converts a function with string, int, int inputs, and an error return into a FuncCallable with variadic IObject arguments.
func (f *Factory) FuncIsiiOe(fn func(string, int, int) error) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 3 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		i2, err := f.ToInt64Arg(1, args[1])
		if err != nil {
			return nil, err
		}
		i3, err := f.ToInt64Arg(2, args[2])
		if err != nil {
			return nil, err
		}
		err = fn(s1, int(i2), int(i3))
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.TrueValue(), nil
	}
}

// FuncIbSOie wraps a function that takes a byte slice and returns an int and error into a FuncCallable for IObject use.
// It ensures the input argument is a single byte-compatible IObject and converts its result to IObject format.
// Returns ErrWrongNumArguments if called with more or less than one argument.
// Returns NewInvalidArgumentError if the input argument isn't byte-compatible.
// Converts the function's error output into an appropriate IObject error.
func (f *Factory) FuncIbSOie(fn func([]byte) (int, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		bs1, err := f.ToByteSliceArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(bs1)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.NewInt(FrameReturnValue, int64(res)), nil
	}
}

// FuncIbSOs wraps a function that converts a byte slice to a string, returning it as a FuncCallable in the custom object system.
// It ensures the input is a single argument of type bytes-compatible, and returns an error for invalid or unsupported types.
// The resulting FuncCallable checks argument validity, applies the provided function, and returns a new String object.
func (f *Factory) FuncIbSOs(fn func([]byte) string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		bs1, err := f.ToByteSliceArg(0, args[0])
		if err != nil {
			return nil, err
		}
		v, err := f.NewString(FrameReturnValue, fn(bs1))
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// FuncIsOie wraps a string-to-int function into a FuncCallable compatible with IObject interface arguments and error handling.
func (f *Factory) FuncIsOie(fn func(string) (int, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		return f.NewInt(FrameReturnValue, int64(res)), nil
	}
}

// FuncIsObSe returns a FuncCallable that wraps a function converting a string to a byte slice and error output.
// It validates input, reports invalid arguments, enforces byte length limits, and converts output to IObject format.
// Uses ErrWrongNumArguments, NewInvalidArgumentError, and ErrBytesLimit for error handling.
func (f *Factory) FuncIsObSe(fn func(string) ([]byte, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		s1, err := f.ToStringArg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(s1)
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		if len(res) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return f.NewBytes(FrameReturnValue, res), nil
	}
}

// FuncIiOsSe converts a function mapping an integer to a slice of strings and an error into a FuncCallable.
func (f *Factory) FuncIiOsSe(fn func(int) ([]string, error)) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		res, err := fn(int(i1))
		if err != nil {
			return f.NewError(FrameReturnValue, err.Error()), nil
		}
		arr := f.NewArray(FrameReturnValue, nil)
		for _, r := range res {
			if len(r) > MaxStringLen {
				return nil, ErrStringLimit
			}
			v, err := f.NewString(FrameReturnValue, r)
			if err != nil {
				return nil, err
			}
			arr.Append(v)
		}
		return arr, nil
	}
}

// FuncIiOs wraps a function of type `func(int) string` into a FuncCallable compatible with the IObject interface system.
// It validates argument count and type, executes the provided function, and converts the result into an IObject.
func (f *Factory) FuncIiOs(fn func(int) string) FuncCallable {
	return func(args ...IObject) (ret IObject, err error) {
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		i1, err := f.ToInt64Arg(0, args[0])
		if err != nil {
			return nil, err
		}
		s := fn(int(i1))
		v, err := f.NewString(FrameReturnValue, s)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

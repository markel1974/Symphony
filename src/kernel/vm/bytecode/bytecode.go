package bytecode

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/vm/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers various object types with the gob encoder for serialization and deserialization.
func init() {
	gob.Register(&Files{})
	//gob.Register(&compiler.SourceFile{})
	gob.Register(&objects.Array{})
	gob.Register(&objects.Bool{})
	gob.Register(&objects.Bytes{})
	gob.Register(&objects.Char{})
	gob.Register(&objects.FunctionCompiled{})
	gob.Register(&objects.Error{})
	gob.Register(&objects.Float{})
	gob.Register(&objects.ArrayImmutable{})
	gob.Register(&objects.MapImmutable{})
	gob.Register(&objects.Int{})
	gob.Register(&objects.Map{})
	gob.Register(&objects.String{})
	gob.Register(&objects.Time{})
	gob.Register(&objects.Undefined{})
	gob.Register(&objects.FunctionUser{})
}

// Bytecode represents a compiled set of bytecode instructions, constants, and metadata for program execution.
type Bytecode struct {
	files        *Files
	mainFunction *objects.FunctionCompiled
	constants    []objects.IObject
}

func NewBytecode() *Bytecode {
	return &Bytecode{
		files: NewFiles(),
	}
}

func (b *Bytecode) AddFile(f IFile) error {
	return b.files.AddFile(f)
}

// Position returns the position in the source file corresponding to the given input offset.
// It provides filename, line, and column information or an error if the position is invalid.
func (b *Bytecode) Position(p int) (*FilePos, error) {
	return b.files.Position(p)
}

// SourceFiles returns the current Files instance associated with the Bytecode.
func (b *Bytecode) SourceFiles() *Files {
	return b.files
}

// Constants returns a slice of IObject that represents the constants used in the bytecode.
func (b *Bytecode) Constants() []objects.IObject {
	return b.constants
}

// MainFunction returns the main compiled function associated with the bytecode.
func (b *Bytecode) MainFunction() (*objects.FunctionCompiled, error) {
	if b.mainFunction == nil {
		return nil, fmt.Errorf("main function not found")
	}
	return b.mainFunction, nil
}

// SetMainFunction sets the main compiled function associated with the bytecode.
func (b *Bytecode) SetMainFunction(f *objects.FunctionCompiled) {
	b.mainFunction = f
}

// SetConstants sets the constants used in the bytecode.
func (b *Bytecode) SetConstants(constants []objects.IObject) {
	b.constants = constants
}

// Encode serializes the Bytecode object to the provided io.Writer using gob encoding and returns any encountered error.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(b.files); err != nil {
		return err
	}
	if err := enc.Encode(b.mainFunction); err != nil {
		return err
	}
	return enc.Encode(b.constants)
}

// CountObjects returns the total number of objects within the Bytecode's constants, including nested objects.
func (b *Bytecode) CountObjects() int {
	n := 0
	for _, c := range b.constants {
		n += objects.CountObjects(c)
	}
	return n
}

// FormatInstructions formats and returns the string representation of the main function's bytecode instructions.
func (b *Bytecode) FormatInstructions() []string {
	return FormatInstructions(b.mainFunction.Data(), 0)
}

// FormatConstants formats and returns a slice of strings representing the constants in the Bytecode object.
func (b *Bytecode) FormatConstants() (output []string) {
	for cIdx, constant := range b.constants {
		switch cn := constant.(type) {
		case *objects.FunctionCompiled:
			output = append(output, fmt.Sprintf("[% 3d] (Compiled Function|%p)", cIdx, &cn))
			for _, l := range FormatInstructions(cn.Data(), 0) {
				output = append(output, fmt.Sprintf("     %s", l))
			}
		default:
			output = append(output, fmt.Sprintf("[% 3d] %s (%s|%p)", cIdx, cn, reflect.TypeOf(cn).Elem().Name(), &cn))
		}
	}
	return
}

// Decode reads and decodes Bytecode data from the provided io.Reader and resolves constants using the given Modules map.
func (b *Bytecode) Decode(r io.Reader, mods *modules.Modules) error {
	if mods == nil {
		mods = modules.NewModuleMap()
	}
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&b.files); err != nil {
		return err
	}
	if err := dec.Decode(&b.mainFunction); err != nil {
		return err
	}
	if err := dec.Decode(&b.constants); err != nil {
		return err
	}
	for i, v := range b.constants {
		fv, err := fixDecodedObject(v, mods)
		if err != nil {
			return err
		}
		b.constants[i] = fv
	}
	return nil
}

// RemoveDuplicates removes duplicate constants from the Bytecode while maintaining their first occurrences.
// It updates constant indices throughout the Bytecode to reflect the changes and ensures consistency for any nested data.
// Returns an error if any unsupported constant type is encountered or if index updates fail.
func (b *Bytecode) RemoveDuplicates() error {
	var deDuped []objects.IObject

	indexMap := make(map[int]int) // mapping from old constant index to new index
	fns := make(map[*objects.FunctionCompiled]int)
	ints := make(map[int64]int)
	strings := make(map[string]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	immutableMaps := make(map[string]int) // for modules

	for curIdx, in := range b.constants {
		switch c := in.(type) {
		case *objects.FunctionCompiled:
			if newIdx, ok := fns[c]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				fns[c] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.MapImmutable:
			modName, err := inferModuleName(c)
			if err != nil {
				return err
			}
			newIdx, ok := immutableMaps[modName]
			if modName != "" && ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				immutableMaps[modName] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.Int:
			if newIdx, ok := ints[c.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				ints[c.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.String:
			if newIdx, ok := strings[c.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				strings[c.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.Float:
			if newIdx, ok := floats[c.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				floats[c.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.Char:
			if newIdx, ok := chars[c.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				chars[c.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		default:
			return fmt.Errorf("unsupported top-level constant type: %s", reflect.TypeOf(c).Elem().Name())
		}
	}
	b.constants = deDuped
	if err := updateConstIndexes(b.mainFunction.Data(), indexMap); err != nil {
		return err
	}
	for _, c := range b.constants {
		switch c := c.(type) {
		case *objects.FunctionCompiled:
			if err := updateConstIndexes(c.Data(), indexMap); err != nil {
				return err
			}
		}
	}
	return nil
}

// fixDecodedObject processes and adjusts decoded objects, handling specific types or values to ensure valid output.
// Returns a potentially modified object or an error if adjustments fail for certain types or constraints.
func fixDecodedObject(o objects.IObject, mods *modules.Modules) (objects.IObject, error) {
	switch o := o.(type) {
	case *objects.Bool:
		if o.Boolean() {
			return objects.FalseValue, nil
		}
		return objects.TrueValue, nil
	case *objects.Undefined:
		return objects.UndefinedValue, nil
	case *objects.Array:
		for i, v := range o.Values() {
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.ArrayImmutable:
		for i, v := range o.Values() {
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.Map:
		for k, v := range o.Values() {
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.Set(k, fv)
		}
	case *objects.MapImmutable:
		modName, err := inferModuleName(o)
		if err != nil {
			return nil, err
		}
		if mod := mods.GetBuiltinModule(modName); mod != nil {
			return mod.AsImmutableMap(modName), nil
		}
		for k, v := range o.Values() {
			if _, isUserFunction := v.(*objects.FunctionUser); isUserFunction {
				return nil, fmt.Errorf("user function not decodable")
			}
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.SetValue(k, fv)
		}
	}
	return o, nil
}

// updateConstIndexes modifies constant and closure indexes in the instructions slice based on the provided index map.
// It processes opcodes, updates corresponding indexes using the map, and returns an error for unsupported opcodes or missing indexes.
func updateConstIndexes(instances []byte, indexMap map[int]int) error {
	i := 0
	for i < len(instances) {
		op := instances[i]
		numOperands := OpcodeOperands[op]
		_, read := ReadOperands(numOperands, instances[i+1:])

		switch op {
		case OpConstant:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			newIdx, ok := indexMap[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			copy(instances[i:], MakeInstruction(op, newIdx))
		case OpClosure:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			numFree := int(instances[i+3])
			newIdx, ok := indexMap[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			copy(instances[i:], MakeInstruction(op, newIdx, numFree))
		default:
			return fmt.Errorf("unsupported opcode: %s", OpcodeNames[op])
		}
		i += 1 + read
	}
	return nil
}

// inferModuleName extracts the module name from an MapImmutable by retrieving the __module_name__ key as a string.
// Returns an empty string if the key is missing or not a valid string type.
func inferModuleName(mod *objects.MapImmutable) (string, error) {
	m, ok := mod.GetValue("__module_name__")
	if !ok {
		return "", fmt.Errorf("missing __module_name__ key")
	}
	modName, ok := m.(*objects.String)
	if !ok {
		return "", fmt.Errorf("invalid __module_name__ value: %s", m.TypeName())
	}
	return modName.Value(), nil
}

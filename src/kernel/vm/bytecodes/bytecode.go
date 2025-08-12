package bytecodes

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/compiler"
	"github.com/markel1974/c64emu/src/kernel/vm/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/opcodes"
)

// init registers various struct types used within the package for gob encoding and decoding.
func init() {
	gob.Register(&compiler.SourceFileSet{})
	gob.Register(&compiler.SourceFile{})
	gob.Register(&objects.Array{})
	gob.Register(&objects.Bool{})
	gob.Register(&objects.Bytes{})
	gob.Register(&objects.Char{})
	gob.Register(&objects.CompiledFunction{})
	gob.Register(&objects.Error{})
	gob.Register(&objects.Float{})
	gob.Register(&objects.ImmutableArray{})
	gob.Register(&objects.ImmutableMap{})
	gob.Register(&objects.Int{})
	gob.Register(&objects.Map{})
	gob.Register(&objects.String{})
	gob.Register(&objects.Time{})
	gob.Register(&objects.Undefined{})
	gob.Register(&objects.UserFunction{})
}

// Bytecode represents the compiled output of a program, including source files, main function, and constant pool.
type Bytecode struct {
	FileSet      *compiler.SourceFileSet
	MainFunction *objects.CompiledFunction
	Constants    []objects.IObject
}

// Encode serializes the Bytecode's FileSet, MainFunction, and Constants into the provided io.Writer using a gob encoder.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(b.FileSet); err != nil {
		return err
	}
	if err := enc.Encode(b.MainFunction); err != nil {
		return err
	}
	return enc.Encode(b.Constants)
}

// CountObjects calculates the total number of objects within the Bytecode's constants, including nested objects.
func (b *Bytecode) CountObjects() int {
	n := 0
	for _, c := range b.Constants {
		n += objects.CountObjects(c)
	}
	return n
}

// FormatInstructions retrieves string representations of the main function's bytecode instructions in the Bytecode object.
func (b *Bytecode) FormatInstructions() []string {
	return FormatInstructions(b.MainFunction.Instructions(), 0)
}

// FormatConstants generates a slice of formatted strings representing the constants in the bytecode.
func (b *Bytecode) FormatConstants() (output []string) {
	for cIdx, cn := range b.Constants {
		switch cn := cn.(type) {
		case *objects.CompiledFunction:
			output = append(output, fmt.Sprintf(
				"[% 3d] (Compiled Function|%p)", cIdx, &cn))
			for _, l := range FormatInstructions(cn.Instructions(), 0) {
				output = append(output, fmt.Sprintf("     %s", l))
			}
		default:
			output = append(output, fmt.Sprintf("[% 3d] %s (%s|%p)",
				cIdx, cn, reflect.TypeOf(cn).Elem().Name(), &cn))
		}
	}
	return
}

// Decode deserializes bytecode data from the provided reader and resolves constants using the given module map.
func (b *Bytecode) Decode(r io.Reader, mods *modules.Modules) error {
	if mods == nil {
		mods = modules.NewModuleMap()
	}

	dec := gob.NewDecoder(r)
	if err := dec.Decode(&b.FileSet); err != nil {
		return err
	}
	// TODO: files in b.FileSet.File does not have their 'set' field properly
	//  set to b.FileSet as it's private field and not serialized by gob
	//  encoder/decoder.
	if err := dec.Decode(&b.MainFunction); err != nil {
		return err
	}
	if err := dec.Decode(&b.Constants); err != nil {
		return err
	}
	for i, v := range b.Constants {
		fv, err := fixDecodedObject(v, mods)
		if err != nil {
			return err
		}
		b.Constants[i] = fv
	}
	return nil
}

// RemoveDuplicates identifies and removes duplicate objects from the Bytecode's constants, updating references accordingly.
func (b *Bytecode) RemoveDuplicates() error {
	var deDuped []objects.IObject

	indexMap := make(map[int]int) // mapping from old constant index to new index
	fns := make(map[*objects.CompiledFunction]int)
	ints := make(map[int64]int)
	strings := make(map[string]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	immutableMaps := make(map[string]int) // for modules

	for curIdx, in := range b.Constants {
		switch c := in.(type) {
		case *objects.CompiledFunction:
			if newIdx, ok := fns[c]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				fns[c] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.ImmutableMap:
			modName := inferModuleName(c)
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
	b.Constants = deDuped
	if err := updateConstIndexes(b.MainFunction.Instructions(), indexMap); err != nil {
		return err
	}
	for _, c := range b.Constants {
		switch c := c.(type) {
		case *objects.CompiledFunction:
			if err := updateConstIndexes(c.Instructions(), indexMap); err != nil {
				return err
			}
		}
	}
	return nil
}

// fixDecodedObject processes a decoded object to ensure proper structure and values, applying modifications as needed.
// It recursively fixes elements in arrays and maps, while handling specific object types like Bool and Undefined.
// Returns the fixed object or an error if processing fails.
func fixDecodedObject(o objects.IObject, mods *modules.Modules) (objects.IObject, error) {
	switch o := o.(type) {
	case *objects.Bool:
		if o.Falsy() {
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
	case *objects.ImmutableArray:
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
	case *objects.ImmutableMap:
		modName := inferModuleName(o)
		if mod := mods.GetBuiltinModule(modName); mod != nil {
			return mod.AsImmutableMap(modName), nil
		}

		for k, v := range o.Values() {
			// encoding of user function not supported
			if _, isUserFunction := v.(*objects.UserFunction); isUserFunction {
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

// updateConstIndexes updates constant indexes in the given bytecode using the provided index mapping. Returns an error if a mapping is missing.
func updateConstIndexes(instances []byte, indexMap map[int]int) error {
	i := 0
	for i < len(instances) {
		op := instances[i]
		numOperands := opcodes.OpcodeOperands[op]
		_, read := opcodes.ReadOperands(numOperands, instances[i+1:])

		switch op {
		case opcodes.OpConstant:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			newIdx, ok := indexMap[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			copy(instances[i:], MakeInstruction(op, newIdx))
		case opcodes.OpClosure:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			numFree := int(instances[i+3])
			newIdx, ok := indexMap[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			copy(instances[i:], MakeInstruction(op, newIdx, numFree))
		default:
			return fmt.Errorf("unsupported opcode: %s", opcodes.OpcodeNames[op])
		}
		i += 1 + read
	}
	return nil
}

// inferModuleName extracts the module name from the given ImmutableMap by looking up the "__module_name__" key.
// Returns an empty string if the key is absent or if the value is not of type String.
func inferModuleName(mod *objects.ImmutableMap) string {
	m, ok := mod.GetValue("__module_name__")
	if !ok {
		return ""
	}
	modName, ok := m.(*objects.String)
	if !ok {
		return ""
	}
	return modName.Value()
}

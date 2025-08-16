package bytecode

import (
	"encoding/gob"
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers various types with the gob package to enable serialization and deserialization.
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
	gob.Register(&objects.FunctionModule{})
}

// Bytecode represents a construct that encapsulates compiled code, associated constants, and object references.
// It aggregates information like source files, constant pool, and referenced objects required for execution.
type Bytecode struct {
	files      *Files
	constants  []objects.IObject
	references []objects.IObject
}

// NewBytecode creates and returns a new instance of Bytecode with an initialized Files object.
func NewBytecode() *Bytecode {
	return &Bytecode{
		files: NewFiles(),
	}
}

// AddFile adds an IFile to the internal files collection of the Bytecode.
// Returns an error if the file cannot be added.
func (b *Bytecode) AddFile(f IFile) error {
	return b.files.AddFile(f)
}

// Position retrieves the FilePos structure for a given position p in the bytecode's source files.
// Returns an error if the position is invalid or does not map to a specific file.
func (b *Bytecode) Position(p int) (*FilePos, error) {
	return b.files.Position(p)
}

// SourceFiles returns the collection of source files associated with the Bytecode.
func (b *Bytecode) SourceFiles() *Files {
	return b.files
}

// Constants returns the list of constant objects stored within the Bytecode instance.
func (b *Bytecode) Constants() []objects.IObject {
	return b.constants
}

// References retrieves the list of IObject references stored in the Bytecode.
func (b *Bytecode) References() []objects.IObject {
	return b.references
}

// SetConstants updates the constants field of the Bytecode with a new slice of IObject.
func (b *Bytecode) SetConstants(constants []objects.IObject) {
	b.constants = constants
}

// SetReferences sets the references field of the Bytecode to the provided slice of objects.
func (b *Bytecode) SetReferences(references []objects.IObject) {
	b.references = references
}

// Encode serializes the Bytecode object and writes it to the provided io.Writer in gob format. Returns an error if encoding fails.
func (b *Bytecode) Encode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(b.files); err != nil {
		return err
	}
	if err := enc.Encode(b.constants); err != nil {
		return err
	}
	if err := enc.Encode(b.references); err != nil {
		return err
	}
	return nil
}

// Decode deserializes the bytecode from the given io.Reader and resolves its components using the provided modules.
func (b *Bytecode) Decode(r io.Reader, loader ILoader) error {
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&b.files); err != nil {
		return err
	}
	if err := dec.Decode(&b.constants); err != nil {
		return err
	}
	for i, v := range b.constants {
		fv, err := fixDecodedObject(v, loader)
		if err != nil {
			return err
		}
		b.constants[i] = fv
	}
	if err := dec.Decode(&b.references); err != nil {
		return err
	}
	for i, v := range b.references {
		fv, err := fixDecodedObject(v, loader)
		if err != nil {
			return err
		}
		b.references[i] = fv
	}
	return nil
}

// RemoveDuplicates removes duplicate entries from the constants slice of the Bytecode instance.
// It updates references within the constants to match the deduplicated list.
// Returns an error if the deduplication process encounters issues.
func (b *Bytecode) RemoveDuplicates() error {
	constantsDeduped, constantsIndexMap, err := b.removeDuplicates(b.constants)
	if err != nil {
		return err
	}
	b.constants = constantsDeduped
	for _, in := range b.constants {
		switch c := in.(type) {
		case *objects.FunctionCompiled:
			if err = updateConstIndexes(c.Data(), constantsIndexMap); err != nil {
				return err
			}
		}
	}
	return nil
}

// removeDuplicates processes a container of objects, removing duplicates and mapping old indices to new indices.
// It returns a deduplicated list of objects, a mapping of old to new indices, and an error if encountered.
func (b *Bytecode) removeDuplicates(container []objects.IObject) ([]objects.IObject, map[int]int, error) {
	var deDuped []objects.IObject
	indexMap := make(map[int]int) // mapping from old constant index to new index
	fns := make(map[*objects.FunctionCompiled]int)
	ints := make(map[int64]int)
	strings := make(map[string]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	immutableMaps := make(map[string]int) // for modules

	for curIdx, in := range container {
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
				return nil, nil, err
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
			return nil, nil, fmt.Errorf("unsupported top-level constant type: %s", reflect.TypeOf(c).Elem().Name())
		}
	}
	return deDuped, indexMap, nil
}

// CompileInstruction returns a bytecode for an opcode and the operands.
func CompileInstruction(opcode Opcode, operands ...int) []byte {
	numOperands := OpcodeToOperands(opcode)
	totalLen := 1
	for _, w := range numOperands {
		totalLen += w
	}
	instruction := make([]byte, totalLen)
	instruction[0] = opcode
	offset := 1
	for i, o := range operands {
		width := numOperands[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			n := uint16(o)
			instruction[offset] = byte(n >> 8)
			instruction[offset+1] = byte(n)
		}
		offset += width
	}
	return instruction
}

// fixDecodedObject ensures that a decoded object is properly reconstructed and compatible with the runtime environment.
// It recursively processes composite objects like arrays and maps, fixing or transforming their elements if necessary.
// Returns the modified object or an error if reconstruction fails.
func fixDecodedObject(o objects.IObject, loader ILoader) (objects.IObject, error) {
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
			fv, err := fixDecodedObject(v, loader)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.ArrayImmutable:
		for i, v := range o.Values() {
			fv, err := fixDecodedObject(v, loader)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.Map:
		for k, v := range o.Values() {
			fv, err := fixDecodedObject(v, loader)
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
		if mod, _ := loader.CompileModule(modName); mod != nil {
			return mod, nil
		}
		for k, v := range o.Values() {
			if _, isUserFunction := v.(*objects.FunctionModule); isUserFunction {
				return nil, fmt.Errorf("user function not decodable")
			}
			fv, err := fixDecodedObject(v, loader)
			if err != nil {
				return nil, err
			}
			o.SetValue(k, fv)
		}
	}
	return o, nil
}

// updateConstIndexes modifies bytecode instructions to remap constant indexes based on the provided index map.
// It updates OpConstant and OpClosure instructions with new constant indexes or returns an error if mapping fails.
func updateConstIndexes(instances []byte, indexMap map[int]int) error {
	i := 0
	for i < len(instances) {
		op := instances[i]
		offset := OpcodeToOperandsOffset(op)
		switch op {
		case OpConstant:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			newIdx, ok := indexMap[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			copy(instances[i:], CompileInstruction(op, newIdx))
		case OpClosure:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			numFree := int(instances[i+3])
			newIdx, ok := indexMap[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			copy(instances[i:], CompileInstruction(op, newIdx, numFree))
		default:
			return fmt.Errorf("unsupported opcode: %s", OpcodeNames(op))
		}
		i += 1 + offset
	}
	return nil
}

// inferModuleName extracts the value of the __module_name__ key from the given MapImmutable if it exists and is a String.
// Returns the extracted string on success or an error if the key is missing or its value is not of type String.
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

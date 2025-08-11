package bytecodes

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/vm/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/opcodes"
)

func init() {
	gob.Register(&objects.SourceFileSet{})
	gob.Register(&objects.SourceFile{})
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

// Bytecode is a compiled instructions and constants.
type Bytecode struct {
	FileSet      *objects.SourceFileSet
	MainFunction *objects.CompiledFunction
	Constants    []objects.Object
}

// Encode writes Bytecode data to the writer.
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

// CountObjects returns the number of objects found in Constants.
func (b *Bytecode) CountObjects() int {
	n := 0
	for _, c := range b.Constants {
		n += objects.CountObjects(c)
	}
	return n
}

// FormatInstructions returns human readable string representations of compiled instructions.
func (b *Bytecode) FormatInstructions() []string {
	return FormatInstructions(b.MainFunction.Instructions, 0)
}

// FormatConstants returns human readable string representations of compiled constants.
func (b *Bytecode) FormatConstants() (output []string) {
	for cIdx, cn := range b.Constants {
		switch cn := cn.(type) {
		case *objects.CompiledFunction:
			output = append(output, fmt.Sprintf(
				"[% 3d] (Compiled Function|%p)", cIdx, &cn))
			for _, l := range FormatInstructions(cn.Instructions, 0) {
				output = append(output, fmt.Sprintf("     %s", l))
			}
		default:
			output = append(output, fmt.Sprintf("[% 3d] %s (%s|%p)",
				cIdx, cn, reflect.TypeOf(cn).Elem().Name(), &cn))
		}
	}
	return
}

// Decode reads Bytecode data from the reader.
func (b *Bytecode) Decode(r io.Reader, mods *modules.ModuleMap) error {
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

// RemoveDuplicates finds and remove the duplicate values in Constants. Note this function mutates Bytecode.
func (b *Bytecode) RemoveDuplicates() {
	var deDuped []objects.Object

	indexMap := make(map[int]int) // mapping from old constant index to new index
	fns := make(map[*objects.CompiledFunction]int)
	ints := make(map[int64]int)
	strings := make(map[string]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	immutableMaps := make(map[string]int) // for modules

	for curIdx, c := range b.Constants {
		switch c := c.(type) {
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
			if newIdx, ok := ints[c.Value]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				ints[c.Value] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.String:
			if newIdx, ok := strings[c.Value]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				strings[c.Value] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.Float:
			if newIdx, ok := floats[c.Value]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				floats[c.Value] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		case *objects.Char:
			if newIdx, ok := chars[c.Value]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				chars[c.Value] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, c)
			}
		default:
			panic(fmt.Errorf("unsupported top-level constant type: %s",
				c.TypeName()))
		}
	}

	// replace with de-duplicated constants
	b.Constants = deDuped

	// update CONST instructions with new indexes
	// main function
	updateConstIndexes(b.MainFunction.Instructions, indexMap)
	// other compiled functions in constants
	for _, c := range b.Constants {
		switch c := c.(type) {
		case *objects.CompiledFunction:
			updateConstIndexes(c.Instructions, indexMap)
		}
	}
}

func fixDecodedObject(o objects.Object, mods *modules.ModuleMap) (objects.Object, error) {
	switch o := o.(type) {
	case *objects.Bool:
		if o.Falsy() {
			return objects.FalseValue, nil
		}
		return objects.TrueValue, nil
	case *objects.Undefined:
		return objects.UndefinedValue, nil
	case *objects.Array:
		for i, v := range o.Value {
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.Value[i] = fv
		}
	case *objects.ImmutableArray:
		for i, v := range o.Value {
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.Value[i] = fv
		}
	case *objects.Map:
		for k, v := range o.Value {
			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.Value[k] = fv
		}
	case *objects.ImmutableMap:
		modName := inferModuleName(o)
		if mod := mods.GetBuiltinModule(modName); mod != nil {
			return mod.AsImmutableMap(modName), nil
		}

		for k, v := range o.Value {
			// encoding of user function not supported
			if _, isUserFunction := v.(*objects.UserFunction); isUserFunction {
				return nil, fmt.Errorf("user function not decodable")
			}

			fv, err := fixDecodedObject(v, mods)
			if err != nil {
				return nil, err
			}
			o.Value[k] = fv
		}
	}
	return o, nil
}

func updateConstIndexes(instances []byte, indexMap map[int]int) {
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
				panic(fmt.Errorf("constant index not found: %d", curIdx))
			}
			copy(instances[i:], MakeInstruction(op, newIdx))
		case opcodes.OpClosure:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			numFree := int(instances[i+3])
			newIdx, ok := indexMap[curIdx]
			if !ok {
				panic(fmt.Errorf("constant index not found: %d", curIdx))
			}
			copy(instances[i:], MakeInstruction(op, newIdx, numFree))
		default:
			log.Printf("unhandled default case")
		}
		i += 1 + read
	}
}

func inferModuleName(mod *objects.ImmutableMap) string {
	if modName, ok := mod.Value["__module_name__"].(*objects.String); ok {
		return modName.Value
	}
	return ""
}

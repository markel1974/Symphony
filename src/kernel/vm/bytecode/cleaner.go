package bytecode

import (
	"fmt"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Cleaner is responsible for processing, fixing, and reconstructing objects, ensuring compatibility with the runtime environment.
type Cleaner struct {
	gk      objects.IGateKeeper
	opcodes *Opcodes
	loader  ILoader
}

// NewCleaner creates and returns a new instance of Cleaner initialized with the provided IGateKeeper, ILoader, and Opcodes.
func NewCleaner(gk objects.IGateKeeper, loader ILoader, opcodes *Opcodes) *Cleaner {
	return &Cleaner{
		gk:      gk,
		loader:  loader,
		opcodes: opcodes,
	}
}

// FixObjects processes a slice of IObject instances, ensuring each object is fixed and reconstructed correctly.
// Returns a slice of fixed IObject instances or an error if the fixing process fails.
func (c *Cleaner) FixObjects(o []objects.IObject) ([]objects.IObject, error) {
	out := make([]objects.IObject, len(o))
	for i, v := range o {
		fv, err := c.FixObject(v)
		if err != nil {
			return nil, err
		}
		out[i] = fv
	}
	return out, nil
}

// FixObject ensures that a decoded object is properly reconstructed and compatible with the runtime environment.
// It recursively processes composite objects like arrays and maps, fixing or transforming their elements if necessary.
// Returns the modified object or an error if reconstruction fails.
func (c *Cleaner) FixObject(o objects.IObject) (objects.IObject, error) {
	switch o := o.(type) {
	case *objects.Bool:
		if o.Falsy() {
			return c.gk.FalseValue(), nil
		}
		return c.gk.TrueValue(), nil
	case *objects.Undefined:
		return c.gk.UndefinedValue(), nil
	case *objects.Array:
		for i, v := range o.Values() {
			fv, err := c.FixObject(v)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.ArrayImmutable:
		for i, v := range o.Values() {
			fv, err := c.FixObject(v)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.Map:
		for k, v := range o.Values() {
			fv, err := c.FixObject(v)
			if err != nil {
				return nil, err
			}
			o.Set(k, fv)
		}
	case *objects.MapImmutable:
		//nothing
	}
	return o, nil
}

// RemoveDuplicates removes duplicate entries from the constants slice of the Bytecode instance.
// It updates references within the constants to match the deduplicated list.
// Returns an error if the deduplication process encounters issues.
func (c *Cleaner) RemoveDuplicates(in []objects.IObject) ([]objects.IObject, error) {
	outDeduped, outIndexContainer, err := c.processDuplicates(in)
	if err != nil {
		return nil, err
	}
	for _, in := range outDeduped {
		switch z := in.(type) {
		case *objects.FuncCompiled:
			if err = c.updateIndexes(z.Data(), outIndexContainer); err != nil {
				return nil, err
			}
		}
	}
	return outDeduped, nil
}

// processDuplicates processes a container of objects, removing duplicates and mapping old indices to new indices.
// It returns a deduplicated list of objects, a mapping of old to new indices, and an error if encountered.
func (c *Cleaner) processDuplicates(container []objects.IObject) ([]objects.IObject, map[int]int, error) {
	var deDuped []objects.IObject
	indexMap := make(map[int]int) // mapping from old constant index to new index
	fns := make(map[*objects.FuncCompiled]int)
	ints := make(map[int64]int)
	strings := make(map[string]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)

	for curIdx, in := range container {
		switch z := in.(type) {
		case *objects.FuncCompiled:
			if newIdx, ok := fns[z]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				fns[z] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, z)
			}
		case *objects.Int:
			if newIdx, ok := ints[z.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				ints[z.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, z)
			}
		case *objects.String:
			if newIdx, ok := strings[z.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				strings[z.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, z)
			}
		case *objects.Float:
			if newIdx, ok := floats[z.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				floats[z.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, z)
			}
		case *objects.Char:
			if newIdx, ok := chars[z.Value()]; ok {
				indexMap[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				chars[z.Value()] = newIdx
				indexMap[curIdx] = newIdx
				deDuped = append(deDuped, z)
			}
		default:
			return nil, nil, fmt.Errorf("unsupported top-level object type: %s", reflect.TypeOf(c).Elem().Name())
		}
	}
	return deDuped, indexMap, nil
}

// updateConstIndexes modifies bytecode instructions to remap constant indexes based on the provided index map.
// It updates OpConstant and OpClosure instructions with new constant indexes or returns an error if mapping fails.
func (c *Cleaner) updateIndexes(instances []byte, indexContainer map[int]int) error {
	i := 0
	for i < len(instances) {
		op := instances[i]
		details := c.opcodes.Opcode(op)
		offset := details.Offset()
		switch op {
		case OpConstant:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			newIdx, ok := indexContainer[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			code, err := c.opcodes.CompileInstruction(op, newIdx)
			if err != nil {
			}
			copy(instances[i:], code)
		case OpClosure:
			curIdx := int(instances[i+2]) | int(instances[i+1])<<8
			numFree := int(instances[i+3])
			newIdx, ok := indexContainer[curIdx]
			if !ok {
				return fmt.Errorf("constant index not found: %d", curIdx)
			}
			code, err := c.opcodes.CompileInstruction(op, newIdx, numFree)
			if err != nil {
				return err
			}
			copy(instances[i:], code)
		default:
			return fmt.Errorf("unsupported opcode: %s", details.Name())
		}
		i += 1 + offset
	}
	return nil
}

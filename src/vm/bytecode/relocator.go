package bytecode

import (
	"fmt"
	"reflect"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Relocator is responsible for processing, fixing, and reconstructing objects, ensuring compatibility with the runtime environment.
type Relocator struct {
	gk        objects.IGateKeeper
	opcodes   *Opcodes
	loader    ILoader
	preserved map[string]bool
}

// NewRelocator creates and returns a new instance of Relocator initialized with the provided IGateKeeper, ILoader, and Opcodes.
func NewRelocator(gk objects.IGateKeeper, loader ILoader, opcodes *Opcodes, preserve []string) *Relocator {
	p := make(map[string]bool)
	for _, v := range preserve {
		p[v] = true
	}
	return &Relocator{
		gk:        gk,
		loader:    loader,
		opcodes:   opcodes,
		preserved: p,
	}
}

// Relocate processes a slice of Bytecode instances, ensuring each bytecode is fixed and reconstructed correctly.
// Returns a new Bytecode instance or an error if the fixing process fails.
func (c *Relocator) Relocate(codes []*Bytecode) (*Bytecode, error) {
	var constants []objects.IObject
	var imports []objects.IObject
	var globals []objects.IObject
	var sourceFiles []IFile
	for _, bc := range codes {
		imports = append(imports, bc.Imports()...)
		constants = append(constants, bc.Constants()...)
		globals = append(globals, bc.Globals()...)
		sourceFiles = append(sourceFiles, bc.SourceFiles().files...)
	}
	var err error
	if constants, err = c.RelocateObjects(constants); err != nil {
		return nil, err
	}
	if imports, err = c.RelocateObjects(imports); err != nil {
		return nil, err
	}
	if globals, err = c.RelocateObjects(globals); err != nil {
		return nil, err
	}
	out := NewBytecode(constants, imports, globals, nil)
	for _, sf := range sourceFiles {
		out.AddFile(sf)
	}
	return out, nil
}

// Fix processes a slice of IObject instances, ensuring each object is fixed and reconstructed correctly.
// Returns a slice of fixed IObject instances or an error if the fixing process fails.
func (c *Relocator) Fix(o []objects.IObject) ([]objects.IObject, error) {
	out := make([]objects.IObject, len(o))
	for i, v := range o {
		fv, err := c.fixObject(v)
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
func (c *Relocator) fixObject(o objects.IObject) (objects.IObject, error) {
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
			fv, err := c.fixObject(v)
			if err != nil {
				return nil, err
			}
			o.SetValue(i, fv)
		}
	case *objects.Map:
		for k, v := range o.Values() {
			fv, err := c.fixObject(v)
			if err != nil {
				return nil, err
			}
			o.Set(k, fv)
		}
	}
	return o, nil
}

// RelocateObjects modifies a slice of IObject instances by deduplicating input and updating bytecode constant indexes accordingly.
func (c *Relocator) RelocateObjects(in []objects.IObject) ([]objects.IObject, error) {
	outDeduped, outIndexContainer, err := c.processDuplicates(in)
	if err != nil {
		return nil, err
	}
	for _, in := range outDeduped {
		switch obj := in.(type) {
		case *objects.FuncCompiled:
			if err = c.updateFuncIndexes(obj, outIndexContainer); err != nil {
				return nil, err
			}
		}
	}
	return outDeduped, nil
}

// processDuplicates processes a container of objects, removing duplicates and mapping old indices to new indices.
// It returns a deduplicated list of objects, a mapping of old to new indices, and an error if encountered.
func (c *Relocator) processDuplicates(container []objects.IObject) ([]objects.IObject, map[int]int, error) {
	var deDuped []objects.IObject
	indexContainer := make(map[int]int)
	ints := make(map[int64]int)
	strings := make(map[string]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	fns := make(map[*objects.FuncCompiled]int)

	for curIdx, in := range container {
		switch obj := in.(type) {
		case *objects.FuncCompiled:
			newIdx := -1
			if _, preserve := c.preserved[obj.Name()]; !preserve {
				if v, ok := fns[obj]; ok {
					newIdx = v
				}
			}
			if newIdx >= 0 {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				fns[obj] = newIdx
				indexContainer[curIdx] = newIdx
				deDuped = append(deDuped, obj)
			}
		case *objects.Int:
			if newIdx, ok := ints[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				ints[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				deDuped = append(deDuped, obj)
			}
		case *objects.String:
			if newIdx, ok := strings[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				strings[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				deDuped = append(deDuped, obj)
			}
		case *objects.Float:
			if newIdx, ok := floats[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				floats[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				deDuped = append(deDuped, obj)
			}
		case *objects.Char:
			if newIdx, ok := chars[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(deDuped)
				chars[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				deDuped = append(deDuped, obj)
			}
		default:
			return nil, nil, fmt.Errorf("unsupported top-level object type: %s", reflect.TypeOf(c).Elem().Name())
		}
	}
	return deDuped, indexContainer, nil
}

// updateConstIndexes modifies bytecode instructions to remap constant indexes based on the provided index map.
// It updates OpConstant and OpClosure instructions with new constant indexes or returns an error if mapping fails.
func (c *Relocator) updateFuncIndexes(fc *objects.FuncCompiled, indexContainer map[int]int) error {
	data := fc.Data()
	i := 0
	for i < len(data) {
		opcode := data[i]
		details := c.opcodes.Opcode(opcode)
		offset := details.Offset()
		switch details.Relocatable() {
		case OpRelocatable:
			curIdx, ok := get16(data, uint(i), 1)
			if !ok {
				return fmt.Errorf("index not found: %d", curIdx)
			}
			newIdx, ok := indexContainer[curIdx]
			if !ok {
				return fmt.Errorf("index not found: %d", curIdx)
			}
			code, err := c.opcodes.CompileInstruction(opcode, newIdx)
			if err != nil {
				return err
			}
			copy(data[i:], code)
		case OpRelocatableFree:
			curIdx, ok := get16(data, uint(i), 1)
			if !ok {
				return fmt.Errorf("index not found: %d", curIdx)
			}
			numFree, ok := get8(data, uint(i), 3)
			if !ok {
				return fmt.Errorf("index not found: %d", curIdx)
			}
			newIdx, ok := indexContainer[curIdx]
			if !ok {
				return fmt.Errorf("index not found: %d", curIdx)
			}
			code, err := c.opcodes.CompileInstruction(opcode, newIdx, numFree)
			if err != nil {
				return err
			}
			copy(data[i:], code)
		default:
			//nothing to do
		}
		i += 1 + offset
	}
	return nil
}

// get16 extracts a 16-bit integer from the provided byte slice using the given base and offset positions.
// Returns the integer and a boolean indicating success or failure due to an out-of-bounds read.
func get16(data []byte, base uint, offset uint) (int, bool) {
	p1 := base + offset
	p2 := p1 + 1
	if p2 >= uint(len(data)) {
		return 0, false
	}
	res := int(data[p2]) | int(data[p1])<<8
	return res, true
}

// get8 retrieves an 8-bit value from the provided data slice at the calculated position (base + offset).
// Returns the value as an int and a boolean indicating success or failure.
// If the position exceeds the data slice boundary, it returns 0 and false.
func get8(data []byte, base uint, offset uint) (int, bool) {
	p1 := base + offset
	if p1 >= uint(len(data)) {
		return 0, false
	}
	return int(data[p1]), true
}

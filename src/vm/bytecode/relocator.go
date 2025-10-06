package bytecode

import (
	"fmt"
	"reflect"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// Relocator is responsible for processing, fixing, and reconstructing objects, ensuring compatibility with the runtime environment.
type Relocator struct {
	gk           objects.IGateKeeper
	opcodes      opcodes.IOpcodes
	instructions *opcodes.Instructions
	loader       ILoader
	preserveFunc map[string]bool
	relocator    *Bytecode
}

// NewRelocator creates and returns a new instance of Relocator initialized with the provided IGateKeeper, ILoader, and Opcodes.
func NewRelocator(gk objects.IGateKeeper, loader ILoader, op opcodes.IOpcodes, preserve ...string) *Relocator {
	p := make(map[string]bool)
	for _, v := range preserve {
		p[v] = true
	}
	return &Relocator{
		gk:           gk,
		loader:       loader,
		opcodes:      op,
		instructions: opcodes.NewInstructions(nil),
		preserveFunc: p,
		relocator:    NewBytecodeEmpty(gk),
	}
}

func (c *Relocator) Add(bc *Bytecode) {
	c.relocator.Append(ConstantsType, bc.Constants())
	c.relocator.Append(ImportsType, bc.Imports())
	c.relocator.Append(GlobalsType, bc.Globals())
	c.relocator.AddFiles(bc.Files())
}

func (c *Relocator) Relocate() (*Bytecode, error) {
	for _, r := range c.relocator.Containers() {
		data, err := c.relocateObjects(r.Objects())
		if err != nil {
			return nil, err
		}
		c.relocator.Assign(r.kind, data)
	}
	return c.relocator, nil
}

/*
// Relocate processes a slice of Bytecode instances, ensuring each bytecode is fixed and reconstructed correctly.
// Returns a new Bytecode instance or an error if the fixing process fails.
func (c *Relocator) Relocate(codes []*Bytecode) (*Bytecode, error) {
	relocator := NewBytecodeEmpty(c.gk)
	for _, bc := range codes {
		relocator.Append(ConstantsType, bc.Constants())
		relocator.Append(ImportsType, bc.Imports())
		relocator.Append(GlobalsType, bc.Globals())
		relocator.AddFiles(bc.Files())
	}
	for _, r := range relocator.Containers() {
		data, err := c.relocateObjects(r.Objects())
		if err != nil {
			return nil, err
		}
		relocator.Assign(r.kind, data)
	}
	return relocator, nil
}

*/

// RelocateObjects modifies a slice of IObject instances by deduplicating input and updating bytecode constant indexes accordingly.
func (c *Relocator) relocateObjects(inObj []objects.IObject) ([]objects.IObject, error) {
	result, indices, err := c.processObjects(inObj)
	if err != nil {
		return nil, err
	}
	for _, in := range result {
		switch obj := in.(type) {
		case *objects.Func:
			if err = c.relocateFunction(obj, indices); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

// processDuplicates processes a container of objects, removing duplicates and mapping old indices to new indices.
// It returns a deduplicated list of objects, a mapping of old to new indices, and an error if encountered.
func (c *Relocator) processObjects(container []objects.IObject) ([]objects.IObject, map[int]int, error) {
	var result []objects.IObject
	indices := make(map[int]int)
	ints := make(map[int64]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	strings := make(map[string]int)
	fns := make(map[*objects.Func]int)

	for idx, in := range container {
		newIndex := len(result)
		foundIndex := -1
		found := false
		switch obj := in.(type) {
		case *objects.Func:
			if _, preserve := c.preserveFunc[obj.Name()]; !preserve {
				var v int
				if v, found = fns[obj]; found {
					foundIndex = v
				}
			}
			if !found {
				fns[obj] = newIndex
			}
		case *objects.Int:
			if foundIndex, found = ints[obj.Value()]; !found {
				ints[obj.Value()] = newIndex
			}
		case *objects.Float:
			if foundIndex, found = floats[obj.Value()]; !found {
				floats[obj.Value()] = newIndex
			}
		case *objects.Char:
			if foundIndex, found = chars[obj.Value()]; !found {
				chars[obj.Value()] = newIndex
			}
		case *objects.String:
			value := obj.GetValue()
			if foundIndex, found = strings[value]; !found {
				strings[value] = newIndex
			}
		default:
			return nil, nil, fmt.Errorf("unsupported top-level object type: %s", reflect.TypeOf(c).Elem().Name())
		}
		if found {
			indices[idx] = foundIndex
		} else {
			indices[idx] = newIndex
			result = append(result, in)
		}
	}
	return result, indices, nil
}

// updateConstIndexes modifies bytecode instructions to remap constant indexes based on the provided index map.
// It updates OpConstant and OpCreateClosure instructions with new constant indexes or returns an error if mapping fails.
func (c *Relocator) relocateFunction(fc *objects.Func, indices map[int]int) error {
	bc := fc.Code()
	var end int
	for i := 0; i < len(bc); {
		c.instructions.Assign(bc)
		targetOpcode, headerSize, ok := c.instructions.Header(uint(i))
		if !ok {
			return fmt.Errorf("invalid instruction at offset %d", i)
		}
		opcode := c.opcodes.Opcode(targetOpcode)
		totalBytes := int(headerSize) + opcode.OperandsBytes()
		if end = i + totalBytes; end > len(bc) {
			return fmt.Errorf("invalid range %d-%d", i, end)
		}
		if opcode.IsRelocatable() {
			instructions := bc[i:end]
			operands, err := opcode.DecompileOperands(headerSize, instructions)
			if err != nil {
				return err
			}
			if len(operands) != opcode.OperandsLen() {
				return fmt.Errorf("invalid operand count: %d", len(operands))
			}
			modified := false
			for _, rel := range opcode.Relocate() {
				if rel >= len(operands) {
					return fmt.Errorf("invalid relocation: %d", rel)
				}
				curr := operands[rel]
				newIdx, ok := indices[curr]
				if !ok {
					return fmt.Errorf("index not found: %d", curr)
				}
				if newIdx != curr {
					operands[rel] = newIdx
					modified = true
				}
			}
			if modified {
				compiled, err := opcode.Compile(operands)
				if err != nil {
					return err
				}
				copy(bc[i:end], compiled)
			}
		}
		i = end
	}
	return nil
}

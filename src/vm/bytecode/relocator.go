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
	if imports, err = c.relocateObjects(imports); err != nil {
		return nil, err
	}
	if constants, err = c.relocateObjects(constants); err != nil {
		return nil, err
	}
	if globals, err = c.relocateObjects(globals); err != nil {
		return nil, err
	}
	out := NewBytecode(constants, imports, globals, nil)
	for _, sf := range sourceFiles {
		out.AddFile(sf)
	}
	return out, nil
}

// RelocateObjects modifies a slice of IObject instances by deduplicating input and updating bytecode constant indexes accordingly.
func (c *Relocator) relocateObjects(inObj []objects.IObject) ([]objects.IObject, error) {
	outDeduped, outIndexContainer, err := c.processDuplicates(inObj)
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
	var outDuped []objects.IObject
	indexContainer := make(map[int]int)
	ints := make(map[int64]int)
	floats := make(map[float64]int)
	chars := make(map[rune]int)
	strings := make(map[string]int)
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
				newIdx = len(outDuped)
				fns[obj] = newIdx
				indexContainer[curIdx] = newIdx
				outDuped = append(outDuped, obj)
			}
		case *objects.Int:
			if newIdx, ok := ints[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(outDuped)
				ints[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				outDuped = append(outDuped, obj)
			}
		case *objects.Float:
			if newIdx, ok := floats[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(outDuped)
				floats[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				outDuped = append(outDuped, obj)
			}
		case *objects.Char:
			if newIdx, ok := chars[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(outDuped)
				chars[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				outDuped = append(outDuped, obj)
			}
		case *objects.String:
			if newIdx, ok := strings[obj.Value()]; ok {
				indexContainer[curIdx] = newIdx
			} else {
				newIdx = len(outDuped)
				strings[obj.Value()] = newIdx
				indexContainer[curIdx] = newIdx
				outDuped = append(outDuped, obj)
			}
		default:
			return nil, nil, fmt.Errorf("unsupported top-level object type: %s", reflect.TypeOf(c).Elem().Name())
		}
	}
	return outDuped, indexContainer, nil
}

// updateConstIndexes modifies bytecode instructions to remap constant indexes based on the provided index map.
// It updates OpConstant and OpClosure instructions with new constant indexes or returns an error if mapping fails.
func (c *Relocator) updateFuncIndexes(fc *objects.FuncCompiled, indexContainer map[int]int) error {
	bc := fc.Data()
	unknownOpcode := c.opcodes.Opcode(OpUnknown)
	var end int
	for i := 0; i < len(bc); {
		if end = i + OpcodeWidth; end > len(bc) {
			return fmt.Errorf("invalid range %d-%d", i, end)
		}
		var targetOpcode OpcodeId
		if unk, err := unknownOpcode.Decompile(bc[i:end]); err != nil {
			return err
		} else if len(unk) == 0 {
			return fmt.Errorf("invalid instruction length: %d", len(unk))
		} else {
			targetOpcode = OpcodeId(unk[0])
		}
		opcode := c.opcodes.Opcode(targetOpcode)
		if end = i + opcode.FullWidth(); end > len(bc) {
			return fmt.Errorf("invalid range %d-%d", i, end)
		}
		if opcode.IsRelocatable() {
			decompiled, err := opcode.Decompile(bc[i:end])
			if err != nil {
				return err
			}
			if len(decompiled) < OpcodeWidth {
				return fmt.Errorf("invalid instruction length: %d", len(decompiled))
			}
			operands := decompiled[OpcodeWidth:]
			modified := false
			for _, rel := range opcode.Relocate() {
				if rel >= len(operands) {
					return fmt.Errorf("invalid relocation: %d", rel)
				}
				curr := operands[rel]
				newIdx, ok := indexContainer[curr]
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

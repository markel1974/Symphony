package bytecode

import (
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

type DisassemblerData struct {
	name string
	data []objects.IObject
}

// Disassembler represents a utility for analyzing and processing bytecode by dissecting its constants and imports.
type Disassembler struct {
	dd           []DisassemblerData
	opcodes      opcodes.IOpcodes
	instructions *opcodes.Instructions
}

// NewDisassembler creates a new Disassembler instance linked to the provided Bytecode object.
func NewDisassembler(b *Bytecode, op opcodes.IOpcodes) *Disassembler {
	d := &Disassembler{
		opcodes:      op,
		instructions: opcodes.NewInstructions(nil),
	}
	d.dd = append(d.dd, DisassemblerData{"Constants", b.Constants()})
	d.dd = append(d.dd, DisassemblerData{"Imports", b.Imports()})
	d.dd = append(d.dd, DisassemblerData{"Globals", b.Globals()})
	return d
}

// Disassemble parses and logs opcode of objects, constants, and imports within the associated bytecode.
func (d *Disassembler) Disassemble(writer io.Writer) error {
	for _, container := range d.dd {
		_, _ = fmt.Fprintf(writer, "---- %s ----\n", container.name)
		count := 0
		for cIdx, obj := range container.data {
			data, err := d.disassembleObject(cIdx, obj)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(writer, "%04d => %s\n", cIdx, data)
			count += obj.Count()
		}
		_, _ = fmt.Fprintf(writer, "total objects: %d\n", count)
	}
	return nil
}

// disassembleObject generates a disassembled representation of a constant object, including detailed instructions for functions.
func (d *Disassembler) disassembleObject(idx int, obj objects.IObject) ([]string, error) {
	var output []string
	if obj == nil {
		return []string{fmt.Sprintf("[% 3d] nil", idx)}, nil
	}
	switch cn := obj.(type) {
	case *objects.FuncCompiled:
		data, err := d.disassembleInstructions(cn.Data())
		if err != nil {
			return nil, err
		}
		output = append(output, fmt.Sprintf("[% 3d] %s (Compiled Function|%p)", idx, cn.Name(), &cn))
		for _, l := range data {
			output = append(output, fmt.Sprintf("\t\t%s", l))
		}
	default:
		kind := reflect.TypeOf(cn)
		output = append(output, fmt.Sprintf("[% 3d] %s -> '%s' (%s|%p)", idx, cn.TypeName(), cn.AsString(), kind.Elem().Name(), &cn))
	}
	return output, nil
}

// disassembleInstructions parses a sequence of bytecode instructions and generates a human-readable representation of the instructions.
func (d *Disassembler) disassembleInstructions(bc []byte) ([]string, error) {
	var out []string
	var end int
	for start := 0; start < len(bc); {
		d.instructions.Assign(bc)
		targetOpcode, headerBytes, ok := d.instructions.Header(uint(start))
		if !ok {
			return nil, fmt.Errorf("invalid instruction at offset %d", start)
		}
		opcode := d.opcodes.Opcode(targetOpcode)
		totalBytes := int(headerBytes) + opcode.OperandsBytes()
		if end = start + totalBytes; end > len(bc) {
			return nil, fmt.Errorf("invalid range %d-%d", start, end)
		}
		instructions := bc[start:end]
		operands, err := opcode.DecompileOperands(headerBytes, instructions)
		if err != nil {
			return nil, err
		}
		if len(operands) != opcode.OperandsLen() {
			return nil, fmt.Errorf("invalid operand count: %d", len(operands))
		}
		k := fmt.Sprintf("%04d %-16s", start, opcode.Name())
		for _, v := range operands {
			k += fmt.Sprintf(" %-5d", v)
		}
		out = append(out, k)
		start = end
	}
	return out, nil
}

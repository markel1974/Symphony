package bytecode

import (
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

type DisassemblerOpcode struct {
	name     string
	start    int
	operands []int
}

// Disassembler represents a utility for analyzing and processing bytecode by dissecting its constants and imports.
type Disassembler struct {
	cd           *Bytecode
	opcodes      opcodes.IOpcodes
	instructions *opcodes.Instructions
}

// NewDisassembler creates a new Disassembler instance linked to the provided Bytecode object.
func NewDisassembler(b *Bytecode, op opcodes.IOpcodes) *Disassembler {
	d := &Disassembler{
		opcodes:      op,
		instructions: opcodes.NewInstructions(nil),
		cd:           b,
	}
	return d
}

// Disassemble parses and logs opcode of objects, constants, and imports within the associated bytecode.
func (d *Disassembler) Disassemble(writer io.Writer) error {
	for _, container := range d.cd.Containers() {
		_, _ = fmt.Fprintf(writer, "---- %s ----\n", container.Type())
		count := 0
		for cIdx, obj := range container.Objects() {
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
		return []string{fmt.Sprintf("<% 4d> nil", idx)}, nil
	}
	switch cn := obj.(type) {
	case *objects.FuncCompiled:
		result, err := d.disassembleInstructions(cn.Data())
		if err != nil {
			return nil, err
		}
		output = append(output, fmt.Sprintf("<% 4d> %s (Compiled Function|%p)\n", idx, cn.Name(), &cn))
		for _, entry := range result {
			header := fmt.Sprintf("%04d %-16s", entry.start, entry.name)
			for _, v := range entry.operands {
				header += fmt.Sprintf(" %-5d", v)
			}
			output = append(output, header+"\n")
		}
	default:
		kind := reflect.TypeOf(cn)
		output = append(output, fmt.Sprintf("<% 4d> %s -> '%s' (%s|%p)", idx, cn.TypeName(), cn.AsString(), kind.Elem().Name(), &cn))
	}
	return output, nil
}

// disassembleInstructions parses a sequence of bytecode instructions and generates a human-readable representation of the instructions.
func (d *Disassembler) disassembleInstructions(bc []byte) ([]DisassemblerOpcode, error) {
	var out []DisassemblerOpcode
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
		out = append(out, DisassemblerOpcode{name: opcode.Name(), start: start, operands: operands})
		start = end
	}
	return out, nil
}

package bytecode

import (
	"fmt"
	"io"
	"strings"

	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// DisassemblerOpcode represents a single opcode with its name, start index, and associated operands.
type DisassemblerOpcode struct {
	name     string
	start    int
	operands []int
}

// Disassembler is used to analyze and break down bytecode into readable instructions through opcode resolution.
// It links a Bytecode object and opcode manager to interpret operations and their respective instructions.
type Disassembler struct {
	cd           *Bytecode
	opcodes      opcodes.IOpcodes
	instructions *opcodes.Instructions
}

// NewDisassembler creates and returns a new instance of Disassembler initialized with the provided Bytecode and IOpcodes.
func NewDisassembler(b *Bytecode, op opcodes.IOpcodes) *Disassembler {
	d := &Disassembler{
		opcodes:      op,
		instructions: opcodes.NewInstructions(nil),
		cd:           b,
	}
	return d
}

// Disassemble writes a textual representation of the disassembled bytecode to the provided writer.
// It processes each container and its objects, formatting their data into readable output.
// Returns an error if disassembling any object fails.
func (d *Disassembler) Disassemble(writer io.Writer) error {
	for _, container := range d.cd.Containers() {
		_, _ = fmt.Fprintf(writer, "---- %s ----\n", container.Type())
		count := 0
		for cIdx, obj := range container.Objects() {
			data, err := d.disassembleObject(obj)
			if err != nil {
				return err
			}
			count += obj.Count()
			for idx, v := range data {
				escaped := EscapeNonPrintable(v)
				if idx == 0 {
					_, _ = fmt.Fprintf(writer, "%04d => %s\n", cIdx, escaped)
				} else {
					_, _ = fmt.Fprintf(writer, "        %s\n", escaped)
				}
			}
		}
		_, _ = fmt.Fprintf(writer, "total objects: %d\n", count)
	}
	return nil
}

// disassembleObject generates a disassembly of the provided object and returns it as a slice of strings or an error.
func (d *Disassembler) disassembleObject(obj objects.IObject) ([]string, error) {
	if obj == nil {
		return []string{fmt.Sprintf("nil")}, nil
	}
	var output []string
	switch cn := obj.(type) {
	case *objects.Func:
		header := fmt.Sprintf("%-16s '%s' (args: %d, locals: %d)", cn.TypeName(), cn.Name(), cn.NumParameters(), cn.NumLocals())
		output = append(output, header)
		instructions, err := d.disassembleInstructions(cn.Code())
		if err != nil {
			return nil, err
		}
		for _, entry := range instructions {
			var operandsStr []string
			for _, v := range entry.operands {
				operandsStr = append(operandsStr, fmt.Sprintf("%-5d", v))
			}
			line := fmt.Sprintf("%04d %-16s %s", entry.start, entry.name, strings.Join(operandsStr, " "))
			output = append(output, line)
		}
	default:
		line := fmt.Sprintf("%-16s %s", cn.TypeName(), cn.AsString())
		output = append(output, line)
	}
	return output, nil
}

// disassembleInstructions parses bytecode into a sequence of disassembled opcodes including their operands and positions.
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

package bytecode

import (
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// Disassembler represents a utility for analyzing and processing bytecode by dissecting its constants and imports.
type Disassembler struct {
	bc           *Bytecode
	opcodes      *opcodes.Opcodes
	instructions *opcodes.Instructions
}

// NewDisassembler creates a new Disassembler instance linked to the provided Bytecode object.
func NewDisassembler(b *Bytecode, op *opcodes.Opcodes) *Disassembler {
	return &Disassembler{
		bc:           b,
		opcodes:      op,
		instructions: opcodes.NewInstructions(nil),
	}
}

// Disassemble parses and logs opcode of objects, constants, and imports within the associated bytecode.
func (d *Disassembler) Disassemble(writer io.Writer) error {
	_, _ = fmt.Fprintf(writer, "--- Object Count ---\n")
	_, _ = fmt.Fprintf(writer, "%d\n", d.CountObjects())
	_, _ = fmt.Fprintf(writer, "--- Constants ---\n")
	constants, err := d.disassembleObjects(d.bc.Constants())
	if err != nil {
		return err
	}
	for idx, v := range constants {
		_, _ = fmt.Fprintf(writer, "%04d => %s\n", idx, v)
	}
	imports, err := d.disassembleObjects(d.bc.Imports())
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(writer, "--- Imports ---\n")
	for idx, v := range imports {
		_, _ = fmt.Fprintf(writer, "%04d => %s\n", idx, v)
	}
	_, _ = fmt.Fprintf(writer, "--- Globals ---\n")
	globals, err := d.disassembleObjects(d.bc.Globals())
	if err != nil {
		return err
	}
	for idx, v := range globals {
		_, _ = fmt.Fprintf(writer, "%04d => %s\n", idx, v)
	}
	return nil
}

// disassembleConstants iterates through bytecode constants, disassembles each, and returns the results as a slice of strings.
func (d *Disassembler) disassembleObjects(e []objects.IObject) ([]string, error) {
	var output []string
	for cIdx, constant := range e {
		data, err := d.disassembleObject(cIdx, constant)
		if err != nil {
			return nil, err
		}
		output = append(output, data...)
	}
	return output, nil
}

// disassembleObject generates a disassembled representation of a constant object, including detailed instructions for functions.
func (d *Disassembler) disassembleObject(cIdx int, constant objects.IObject) ([]string, error) {
	var output []string
	if constant == nil {
		return []string{fmt.Sprintf("[% 3d] nil", cIdx)}, nil
	}
	switch cn := constant.(type) {
	case *objects.FuncCompiled:
		output = append(output, fmt.Sprintf("[% 3d] %s (Compiled Function|%p)", cIdx, cn.Name(), &cn))
		data, err := d.disassembleInstructions(cn.Data(), 0)
		if err != nil {
			return nil, err
		}
		for _, l := range data {
			output = append(output, fmt.Sprintf("\t\t%s", l))
		}
	default:
		z := reflect.TypeOf(cn)
		output = append(output, fmt.Sprintf("[% 3d] %s -> '%s' (%s|%p)", cIdx, cn.TypeName(), cn.AsString(), z.Elem().Name(), &cn))
	}
	return output, nil
}

// disassembleInstructions parses a sequence of bytecode instructions and generates a human-readable representation of the instructions.
func (d *Disassembler) disassembleInstructions(bc []byte, posOffset int) ([]string, error) {
	var out []string
	var end int
	for i := 0; i < len(bc); {
		d.instructions.Assign(bc)
		targetOpcode, headerBytes, ok := d.instructions.Header(uint(i))
		if !ok {
			return nil, fmt.Errorf("invalid instruction at offset %d", i)
		}
		opcode := d.opcodes.Opcode(targetOpcode)
		totalBytes := int(headerBytes) + opcode.OperandsBytes()
		if end = i + totalBytes; end > len(bc) {
			return nil, fmt.Errorf("invalid range %d-%d", i, end)
		}
		instructions := bc[i:end]
		operands, err := opcode.DecompileOperands(headerBytes, instructions)
		if err != nil {
			return nil, err
		}
		if len(operands) != opcode.OperandsLen() {
			return nil, fmt.Errorf("invalid operand count: %d", len(operands))
		}
		k := fmt.Sprintf("%04d %-16s", posOffset+i, opcode.Name())
		for _, v := range operands {
			k += fmt.Sprintf(" %-5d", v)
		}
		out = append(out, k)
		i = end
	}
	return out, nil

}

// CountObjects computes the total number of objects in the Bytecode's constants and imports, including nested objects.
func (d *Disassembler) CountObjects() int {
	n := 0
	for _, c := range d.bc.Constants() {
		n += d.countObjects(c)
	}
	for _, c := range d.bc.Imports() {
		n += d.countObjects(c)
	}
	return n
}

// countObjects recursively counts the total number of objects contained in the given IObject, including nested structures.
func (d *Disassembler) countObjects(in objects.IObject) int {
	c := 1
	switch o := in.(type) {
	case *objects.Array:
		for _, v := range o.Values() {
			c += d.countObjects(v)
		}
	case *objects.Map:
		for _, v := range o.Values() {
			c += d.countObjects(v)
		}
	case *objects.Struct:
		for _, v := range o.Values() {
			c += d.countObjects(v)
		}
	case *objects.Error:
		c += d.countObjects(o.Value())
	}
	return c
}

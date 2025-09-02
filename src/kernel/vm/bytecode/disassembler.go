package bytecode

import (
	"fmt"
	"io"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Disassembler represents a utility for analyzing and processing bytecode by dissecting its constants and imports.
type Disassembler struct {
	bc *Bytecode
}

// NewDisassembler creates a new Disassembler instance linked to the provided Bytecode object.
func NewDisassembler(b *Bytecode) *Disassembler {
	return &Disassembler{
		bc: b,
	}
}

// Disassemble parses and logs details of objects, constants, and imports within the associated bytecode.
func (d *Disassembler) Disassemble(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "--- Object Count ---\n")
	_, _ = fmt.Fprintf(writer, "%d\n", d.CountObjects())
	_, _ = fmt.Fprintf(writer, "--- Constants ---\n")
	for idx, v := range d.disassembleConstants() {
		_, _ = fmt.Fprintf(writer, "%04d => %s\n", idx, v)
	}
	_, _ = fmt.Fprintf(writer, "--- Imports ---\n")
	for idx, v := range d.disassembleReferences() {
		_, _ = fmt.Fprintf(writer, "%04d => %s\n", idx, v)
	}
	_, _ = fmt.Fprintf(writer, "--- Globals ---\n")
	for idx, v := range d.disassembleGlobal() {
		_, _ = fmt.Fprintf(writer, "%04d => %s\n", idx, v)
	}
}

// disassembleConstants iterates through bytecode constants, disassembles each, and returns the results as a slice of strings.
func (d *Disassembler) disassembleConstants() []string {
	var output []string
	for cIdx, constant := range d.bc.Constants() {
		output = append(output, d.disassembleObject(cIdx, constant)...)
	}
	return output
}

// disassembleConstants iterates through bytecode constants, disassembles each, and returns the results as a slice of strings.
func (d *Disassembler) disassembleGlobal() []string {
	var output []string
	for cIdx, constant := range d.bc.Globals() {
		output = append(output, d.disassembleObject(cIdx, constant)...)
	}
	return output
}

// disassembleReferences disassembles and formats the imports section of the bytecode into a slice of string representations.
func (d *Disassembler) disassembleReferences() []string {
	var output []string
	for cIdx, constant := range d.bc.Imports() {
		output = append(output, d.disassembleObject(cIdx, constant)...)
	}
	return output
}

// disassembleObject generates a disassembled representation of a constant object, including detailed instructions for functions.
func (d *Disassembler) disassembleObject(cIdx int, constant objects.IObject) []string {
	var output []string
	if constant == nil {
		return []string{fmt.Sprintf("[% 3d] nil", cIdx)}
	}
	switch cn := constant.(type) {
	case *objects.FuncCompiled:
		output = append(output, fmt.Sprintf("[% 3d] %s (Compiled Function|%p)", cIdx, cn.Name(), &cn))
		for _, l := range d.disassembleInstructions(cn.Data(), 0) {
			output = append(output, fmt.Sprintf("\t\t%s", l))
		}
	default:
		z := reflect.TypeOf(cn)
		output = append(output, fmt.Sprintf("[% 3d] %s -> '%s' (%s|%p)", cIdx, cn.TypeName(), cn.AsString(), z.Elem().Name(), &cn))
	}
	return output
}

// disassembleInstructions parses a sequence of bytecode instructions and generates a human-readable representation of the instructions.
func (d *Disassembler) disassembleInstructions(bc []byte, posOffset int) []string {
	var out []string
	for i := 0; i < len(bc); {
		details := d.bc.opcodes.Opcode(bc[i])
		numOperands, operands, read := d.computeOperands(details, bc[i+1:])
		switch len(numOperands) {
		case 0:
			out = append(out, fmt.Sprintf("%04d %-7s", posOffset+i, details.Name()))
		case 1:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d", posOffset+i, details.Name(), operands[0]))
		case 2:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d %-5d", posOffset+i, details.Name(), operands[0], operands[1]))
		}
		i += 1 + read
	}
	return out
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
	case *objects.Error:
		c += d.countObjects(o.Value())
	}
	return c
}

// computeOperands extracts operand details from a given opcode and instruction sequence, returning operand widths, values, and bytes read.
func (d *Disassembler) computeOperands(details *Opcode, ins []byte) ([]int, []int, int) {
	if len(details.Operands()) == 0 {
		return nil, nil, 0
	}
	var retOperands []int
	var offset int
	for _, width := range details.Operands() {
		switch width {
		case ByteSize:
			if offset >= len(ins) {
				return nil, nil, 0
			}
			retOperands = append(retOperands, int(ins[offset]))
		case Uint16Size:
			if offset+1 >= len(ins) {
				return nil, nil, 0
			}
			retOperands = append(retOperands, int(ins[offset+1])|int(ins[offset])<<8)
		}
		offset += width
	}
	return details.Operands(), retOperands, offset
}

package bytecode

import (
	"fmt"
	"log"
	"reflect"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Disassembler represents a utility for analyzing and processing bytecode by dissecting its constants and references.
type Disassembler struct {
	bc *Bytecode
}

// NewDisassembler creates a new Disassembler instance linked to the provided Bytecode object.
func NewDisassembler(b *Bytecode) *Disassembler {
	return &Disassembler{
		bc: b,
	}
}

// Disassemble parses and logs details of objects, constants, and references within the associated bytecode.
func (d *Disassembler) Disassemble() {
	log.Println("--- Object Count ---")
	log.Println(d.countObjects())
	log.Println("--- Constants ---")
	for idx, v := range d.disassembleConstants() {
		log.Printf("%04d => %s\n", idx, v)
	}
	log.Println("--- References ---")
	for idx, v := range d.disassembleReferences() {
		log.Printf("%04d => %s\n", idx, v)
	}
}

// CountObjects computes the total number of objects in the Bytecode's constants and references, including nested objects.
func (d *Disassembler) countObjects() int {
	n := 0
	for _, c := range d.bc.Constants() {
		n += objects.CountObjects(c)
	}
	for _, c := range d.bc.References() {
		n += objects.CountObjects(c)
	}
	return n
}

// disassembleConstants iterates through bytecode constants, disassembles each, and returns the results as a slice of strings.
func (d *Disassembler) disassembleConstants() []string {
	var output []string
	for cIdx, constant := range d.bc.Constants() {
		output = append(output, d.disassembleObject(cIdx, constant)...)
	}
	return output
}

// disassembleReferences disassembles and formats the references section of the bytecode into a slice of string representations.
func (d *Disassembler) disassembleReferences() []string {
	var output []string
	for cIdx, constant := range d.bc.References() {
		output = append(output, d.disassembleObject(cIdx, constant)...)
	}
	return output
}

// disassembleObject generates a disassembled representation of a constant object, including detailed instructions for functions.
func (d *Disassembler) disassembleObject(cIdx int, constant objects.IObject) []string {
	var output []string
	switch cn := constant.(type) {
	case *objects.FunctionCompiled:
		output = append(output, fmt.Sprintf("[% 3d] %s (Compiled Function|%p)", cIdx, cn.Name(), &cn))
		for _, l := range d.disassembleInstructions(cn.Data(), 0) {
			output = append(output, fmt.Sprintf("\t\t%s", l))
		}
	default:
		output = append(output, fmt.Sprintf("[% 3d] %s (%s|%p)", cIdx, cn, reflect.TypeOf(cn).Elem().Name(), &cn))
	}
	return output
}

// disassembleInstructions parses a sequence of bytecode instructions and generates a human-readable representation of the instructions.
func (d *Disassembler) disassembleInstructions(bc []byte, posOffset int) []string {
	var out []string
	for i := 0; i < len(bc); {
		opcode := bc[i]
		numOperands, operands, read := OpcodeToOperandsDetails(opcode, bc[i+1:])
		switch len(numOperands) {
		case 0:
			out = append(out, fmt.Sprintf("%04d %-7s", posOffset+i, OpcodeNames(opcode)))
		case 1:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d", posOffset+i, OpcodeNames(opcode), operands[0]))
		case 2:
			out = append(out, fmt.Sprintf("%04d %-7s %-5d %-5d", posOffset+i, OpcodeNames(opcode), operands[0], operands[1]))
		}
		i += 1 + read
	}
	return out
}

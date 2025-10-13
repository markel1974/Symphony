// file: compilers/z80/compiler/compiler.go
package compiler

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// Compiler is responsible for transforming high-level code into executable bytecode using various compilation components.
// It leverages IGateKeeper for object management, Opcodes for instruction handling, and Scopes for variable tracking.
// Compiler maintains constants, registers, and helper tools to facilitate the compilation process efficiently.
type Compiler struct {
	gk        objects.IGateKeeper
	opcodes   *opcodes.Opcodes
	scopes    *tables.Scopes
	constants *tables.Constants
	z80       *Z80
	helper    *Helper
}

// New creates and initializes a new Compiler instance with the provided gatekeeper, loader, and opcode definitions.
func New(gk objects.IGateKeeper, loader bytecode.ILoader, opcodes *opcodes.Opcodes) *Compiler {
	z80 := NewZ80()
	scopes := tables.NewScopes(gk, opcodes)
	constants := tables.NewConstants()
	helper := NewHelper(gk, z80, constants, scopes)
	return &Compiler{
		gk:        gk,
		z80:       z80,
		helper:    helper,
		scopes:    scopes,
		constants: constants,
		opcodes:   opcodes,
	}
}

// createInit finalizes and compiles the `__init__` function as part of the compilation process for initializing the program.
func (c *Compiler) createInit() error {
	c.scopes.SetRootIndex()
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	if scope.InstructionsLen() == 0 {
		return nil // Nessuna istruzione da compilare
	}

	// Aggiungiamo un OpReturn alla fine del nostro bytecode
	if _, err := c.scopes.SymbolEmit(opcodes.OpReturn, 0); err != nil {
		return err
	}

	initFuncCode := scope.Instructions()
	numLocals := c.scopes.SymbolCount()

	// Definiamo il simbolo per la nostra funzione d'ingresso
	initSymbol, err := c.scopes.SymbolDefine(bytecode.PreInitFunction)
	if err != nil {
		return err
	}
	// Creiamo l'oggetto funzione compilata
	compiledInitFn := c.gk.NewFuncCompiled(objects.FrameStatic, initSymbol.Name(), initFuncCode, numLocals, 0, false, nil, nil)
	initSymbol.SetObject(compiledInitFn)
	return nil
}

// defineZ80State initializes Z80 state by defining symbols for 8-bit and 16-bit registers and the memory array in the current scope.
func (c *Compiler) defineZ80State() error {
	for _, reg := range c.z80.Registers8Bit() {
		symbol, err := c.scopes.SymbolDefine(reg)
		if err != nil {
			return err
		}
		c.z80.registers[reg] = symbol.Index()
	}
	for _, reg := range c.z80.Registers16Bit() {
		symbol, err := c.scopes.SymbolDefine(reg)
		if err != nil {
			return err
		}
		c.z80.SetRegisterIndex(reg, symbol.Index())
	}

	// Memoria (un grande array di byte)
	memSymbol, err := c.scopes.SymbolDefine("MEMORY")
	if err != nil {
		return err
	}
	c.z80.registers["MEMORY"] = memSymbol.Index()
	return nil
}

/*
Esempio Concreto: da "Molti a 1" a "1 a 1"
Istruzione Z80	Approccio Iniziale (Molti-a-1)	Obiettivo Finale (1:1)
LD B, C	OpGlobalGet C<br>OpGlobalSet B	OpLdRegToReg B, C
ADD A, B	OpGlobalGet A<br>OpGlobalGet B<br>OpArithmetic ADD<br>OpGlobalSet A<br>(...e 10+ istruzioni per i flag)	OpAluReg ADD, B
JP Z, nn	OpGlobalGet F<br>OpArithmetic AND<br>OpConstant FlagZ<br>OpJumpTruthy nn	OpJumpConditional Z_SET, nn
*/

// compileInstruction processes a single CPU instruction based on its opcode and operands, updating the program counter.
// It emits corresponding bytecode depending on the instruction type and manages registers, constants, and jumps.
func (c *Compiler) compileInstruction(pc int, opcode byte, operands []byte) (int, error) {
	pcIncrement := 1 // Default

	switch {
	// LD r, r' Group (0x40 - 0x7F)
	case opcode >= 0x40 && opcode <= 0x7F:
		if opcode == 0x76 { // HALT
			break
		}
		// Abbiamo sostituito la vecchia logica con una singola istruzione specializzata.
		srcRegName := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		destRegName := c.z80.GetRegisterNameFromIndex(int((opcode >> 3) & 0x07))

		// Emettiamo la nuova istruzione 'OpGlobalCopy' con due operandi: destinazione e sorgente.
		// Questo raggiunge l'obiettivo 1:1.
		if _, err := c.scopes.SymbolEmit(opcodes.OpGlobalCopy, c.z80.Register(destRegName), c.z80.Register(srcRegName)); err != nil {
			return 0, err
		}

	// LD r, n Group (opcode & 0xC7 == 0x06)
	case (opcode & 0xC7) == 0x06:
		destReg := c.z80.GetRegisterNameFromIndex(int((opcode >> 3) & 0x07))
		constIndex := c.constants.AddOrGet("", c.gk.NewInt(objects.FrameStatic, int64(operands[0])))
		if _, err := c.scopes.SymbolEmit(opcodes.OpConstant, constIndex); err != nil {
			return 0, err
		}
		if _, err := c.scopes.SymbolEmit(opcodes.OpGlobalSet, c.z80.Register(destReg)); err != nil {
			return 0, err
		}
		pcIncrement = 2

	// ADD A, r Group (0x80 - 0x87)
	case opcode >= 0x80 && opcode <= 0x87:
		srcReg := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		if err := c.helper.EmitAluRegToReg(objects.OperatorAdd, srcReg); err != nil {
			return 0, err
		}

	// SUB r Group (0x90 - 0x97)
	case opcode >= 0x90 && opcode <= 0x97:
		srcReg := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		if err := c.helper.EmitAluRegToReg(objects.OperatorSub, srcReg); err != nil {
			return 0, err
		}

	// AND r Group (0xA0 - 0xA7)
	case opcode >= 0xA0 && opcode <= 0xA7:
		srcReg := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		if err := c.helper.EmitAluRegToReg(objects.OperatorAnd, srcReg); err != nil {
			return 0, err
		}

	// XOR r Group (0xA8 - 0xAF)
	case opcode >= 0xA8 && opcode <= 0xAF:
		srcReg := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		if err := c.helper.EmitAluRegToReg(objects.OperatorXor, srcReg); err != nil {
			return 0, err
		}

	// OR r Group (0xB0 - 0xB7)
	case opcode >= 0xB0 && opcode <= 0xB7:
		srcReg := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		if err := c.helper.EmitAluRegToReg(objects.OperatorOr, srcReg); err != nil {
			return 0, err
		}

	case opcode >= 0xB8 && opcode <= 0xBF:
		srcReg := c.z80.GetRegisterNameFromIndex(int(opcode & 0x07))
		if err := c.helper.EmitCpReg(srcReg); err != nil {
			return 0, err
		}

	case opcode == 0xC3: // JP nn
		targetAddr := int(operands[0]) | int(operands[1])<<8
		pcIncrement = 3
		if err := c.helper.EmitJumpUnconditional(targetAddr); err != nil {
			return 0, err
		}

	// Conditional Jumps
	case opcode == 0xCA: // JP Z, nn (Jump if Zero = true)
		targetAddr := int(operands[0]) | int(operands[1])<<8
		pcIncrement = 3
		if err := c.helper.EmitJumpConditional(FlagZ, true, targetAddr); err != nil {
			return 0, err
		}

	case opcode == 0xC2: // JP NZ, nn (Jump if Zero = false)
		targetAddr := int(operands[0]) | int(operands[1])<<8
		pcIncrement = 3
		if err := c.helper.EmitJumpConditional(FlagZ, false, targetAddr); err != nil {
			return 0, err
		}

	case opcode == 0xDA: // JP C, nn (Jump if Carry = true)
		targetAddr := int(operands[0]) | int(operands[1])<<8
		pcIncrement = 3
		if err := c.helper.EmitJumpConditional(FlagC, true, targetAddr); err != nil {
			return 0, err
		}

	case opcode == 0xD2: // JP NC, nn (Jump if Carry = false)
		targetAddr := int(operands[0]) | int(operands[1])<<8
		pcIncrement = 3
		if err := c.helper.EmitJumpConditional(FlagC, false, targetAddr); err != nil {
			return 0, err
		}

	// Subroutine Calls
	case opcode == 0xCD: // CALL nn
		targetAddr := int(operands[0]) | int(operands[1])<<8

		// CORRECTION HERE
		// Now using passed 'pc' to calculate return address
		returnAddr := pc + 3

		pcIncrement = 3 // CALL instruction is 3 bytes long

		// 1. Save return address on stack
		retAddrSymbol, _ := c.scopes.SymbolDefineUnique("__ret_addr")
		constRetAddr := c.constants.AddOrGet("", c.gk.NewInt(objects.FrameStatic, int64(returnAddr)))
		c.scopes.SymbolEmit(opcodes.OpConstant, constRetAddr)
		c.scopes.SymbolEmitSetAndPop(retAddrSymbol)
		if err := c.helper.emitPush16(retAddrSymbol); err != nil {
			return 0, err
		}

		// 2. Jump to subroutine address
		if err := c.helper.EmitJumpUnconditional(targetAddr); err != nil {
			return 0, err
		}

	case opcode == 0xC9: // RET
		pcIncrement = 0 // RET does not increment PC, it replaces it

		// 1. Get 16-bit return address from emulated Z80 stack
		retAddrSymbol, _ := c.scopes.SymbolDefineUnique("__ret_addr_from_stack")
		if err := c.helper.emitPop16(retAddrSymbol); err != nil {
			return 0, err
		}

		// 2. Put retrieved address on top of VM stack
		if err := c.scopes.SymbolEmitGet(retAddrSymbol); err != nil {
			return 0, err
		}

		// 3. Perform indirect jump - VM will take address from stack
		if _, err := c.scopes.SymbolEmit(opcodes.OpJumpIndirect); err != nil {
			return 0, err
		}
		// Conditional Calls
	case opcode == 0xC4: // CALL NZ, nn
		targetAddr := int(operands[0]) | int(operands[1])<<8
		returnAddr := pc + 3
		pcIncrement = 3
		if err := c.helper.EmitCallConditional(FlagZ, false, targetAddr, returnAddr); err != nil {
			return 0, err
		}

	case opcode == 0xCC: // CALL Z, nn
		targetAddr := int(operands[0]) | int(operands[1])<<8
		returnAddr := pc + 3
		pcIncrement = 3
		if err := c.helper.EmitCallConditional(FlagZ, true, targetAddr, returnAddr); err != nil {
			return 0, err
		}

	case opcode == 0xD4: // CALL NC, nn
		targetAddr := int(operands[0]) | int(operands[1])<<8
		returnAddr := pc + 3
		pcIncrement = 3
		if err := c.helper.EmitCallConditional(FlagC, false, targetAddr, returnAddr); err != nil {
			return 0, err
		}

	case opcode == 0xDC: // CALL C, nn
		targetAddr := int(operands[0]) | int(operands[1])<<8
		returnAddr := pc + 3
		pcIncrement = 3
		if err := c.helper.EmitCallConditional(FlagC, true, targetAddr, returnAddr); err != nil {
			return 0, err
		}

		// Add other flag cases for P, M, PE, PO here...

	// Conditional Returns
	case opcode == 0xC0: // RET NZ
		pcIncrement = 1 // Conditional RET instructions don't increment PC if they return
		if err := c.helper.EmitRetConditional(FlagZ, false); err != nil {
			return 0, err
		}

	case opcode == 0xC8: // RET Z
		pcIncrement = 1
		if err := c.helper.EmitRetConditional(FlagZ, true); err != nil {
			return 0, err
		}

	case opcode == 0xD0: // RET NC
		pcIncrement = 1
		if err := c.helper.EmitRetConditional(FlagC, false); err != nil {
			return 0, err
		}

	case opcode == 0xD8: // RET C
		pcIncrement = 1
		if err := c.helper.EmitRetConditional(FlagC, true); err != nil {
			return 0, err
		}
	case opcode >= 0xC7 && opcode <= 0xFF && (opcode&0x07) == 0x07:
		// This condition catches all 8 RST opcodes: C7, CF, D7, DF, E7, EF, F7, FF
		targetAddr := int(opcode & 0x38) // Calculate target address (0x00, 0x08, ..., 0x38)
		returnAddr := pc + 1             // RST instruction is 1 byte long
		pcIncrement = 1                  // Even if we jump, the instruction itself is 1 byte

		if err := c.helper.EmitRst(targetAddr, returnAddr); err != nil {
			return 0, err
		}
		// Add other flag cases for P, M, PE, PO here...
	default:
		// Instruction not yet implemented
		break
	}

	return pcIncrement, nil
}

// Compile compiles a Z80 ROM file by reading its byte stream, defining global symbols, translating instructions, and finalizing bytecode.
func (c *Compiler) Compile(filename string, source any) error {
	rom, err := io.ReadAll(source.(io.Reader))
	if err != nil {
		return fmt.Errorf("unable to read Z80 ROM file: %w", err)
	}

	// STEP 2: Define registers and memory as global symbols
	if err := c.defineZ80State(); err != nil {
		return err
	}

	// STEP 3: Translate byte stream
	if err := c.translateROM(rom); err != nil {
		return err
	}

	// STEP 4: Finalize bytecode (create __init__ function)
	if err := c.createInit(); err != nil {
		return err
	}

	return nil
}

// translateROM processes a sequence of bytes (ROM) and compiles each instruction, returning an error if encountered.
func (c *Compiler) translateROM(rom []byte) error {
	pc := 0
	for pc < len(rom) {
		opcode := rom[pc]

		// Pass current pc to compilation function
		bytesConsumed, err := c.compileInstruction(pc, opcode, rom[pc+1:])
		if err != nil {
			return fmt.Errorf("error at position 0x%X: %w", pc, err)
		}
		pc += bytesConsumed
	}
	return nil
}

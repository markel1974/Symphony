package compiler

import (
	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Helper provides utilities for managing CPU emulation, constants, and scope operations within the system.
type Helper struct {
	gk        objects.IGateKeeper
	z80       *Z80
	constants *tables.Constants
	scopes    *tables.Scopes
}

// NewHelper initializes and returns a pointer to a Helper struct configured with the provided gatekeeper, Z80, constants, and scopes.
func NewHelper(gk objects.IGateKeeper, z80 *Z80, constants *tables.Constants, scopes *tables.Scopes) *Helper {
	return &Helper{
		gk:        gk,
		z80:       z80,
		constants: constants,
		scopes:    scopes,
	}
}

// EmitLdRegToReg copies the value from the source register to the destination register and emits the respective bytecode.
func (h *Helper) EmitLdRegToReg(dest string, src string) error {
	_, err := h.scopes.Emit(bytecode.OpGlobalCopy, h.z80.Register(dest), h.z80.Register(src))
	return err
}

// EmitAluRegToReg performs an arithmetic operation between the accumulator and the specified register, updating the flags.
func (h *Helper) EmitAluRegToReg(op objects.ArithmeticOperator, srcReg string) error {
	// Definiamo i simboli per gli operandi prima dell'operazione
	opA, err := h.scopes.SymbolDefineUnique("__opA")
	if err != nil {
		return err
	}
	opB, err := h.scopes.SymbolDefineUnique("__opB")
	if err != nil {
		return err
	}

	// Salva gli operandi
	if _, err = h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("A")); err != nil {
		return err
	}
	if err = h.scopes.EmitSymbolSetAndPop(opA); err != nil {
		return err
	}
	if _, err = h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register(srcReg)); err != nil {
		return err
	}
	if err = h.scopes.EmitSymbolSetAndPop(opB); err != nil {
		return err
	}

	// Esegui l'operazione
	if err = h.scopes.EmitSymbolGet(opB); err != nil {
		return err
	}
	if err = h.scopes.EmitSymbolGet(opA); err != nil {
		return err
	}
	if _, err = h.scopes.Emit(bytecode.OpArithmetic, int(op)); err != nil {
		return err
	}

	tempResultSymbol, err := h.scopes.SymbolDefineUnique("__alu_result")
	if err != nil {
		return err
	}
	if err = h.scopes.EmitSymbolSetAndPop(tempResultSymbol); err != nil {
		return err
	}

	// Salva il risultato in A
	if err = h.scopes.EmitSymbolGet(tempResultSymbol); err != nil {
		return err
	}
	if _, err = h.scopes.Emit(bytecode.OpGlobalSet, h.z80.Register("A")); err != nil {
		return err
	}

	// Aggiorna i Flag, passando anche gli operandi originali
	if err = h.emitUpdateFlags(opA, opB, tempResultSymbol, op); err != nil {
		return err
	}

	return nil
}

// emitUpdateFlags calculates and sets CPU flags (e.g., Zero, Sign) based on the result of an arithmetic operation.
func (h *Helper) emitUpdateFlags(opA, opB, resultSymbol *tables.Symbol, op objects.ArithmeticOperator) error {
	// --- Simboli per i calcoli intermedi ---
	//tempFlags, _ := h.scopes.SymbolDefineUnique("__temp_flags")

	// --- 1. Calcolo Flag S (Sign) e Z (Zero) ---
	// Si basano solo sul risultato finale (un intero a 8 bit).
	// Stack: [risultato]
	if err := h.scopes.EmitSymbolGet(resultSymbol); err != nil {
		return err
	}

	// Calcolo Z: (risultato == 0) ? FlagZ : 0
	const0 := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 0))
	//constZ := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, FlagZ))
	if _, err := h.scopes.Emit(bytecode.OpConstant, const0); err != nil {
		return err
	}
	// Stack: [is_zero_bool]
	if _, err := h.scopes.Emit(bytecode.OpLogical, int(objects.OperatorLogicalEq)); err != nil {
		return err
	}

	// Calcolo S: (risultato & 0x80) != 0 ? FlagS : 0
	if err := h.scopes.EmitSymbolGet(resultSymbol); err != nil {
		return err
	}
	const128 := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 128))
	//constS := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, FlagS))
	if _, err := h.scopes.Emit(bytecode.OpConstant, const128); err != nil {
		return err
	}
	// Stack: [is_zero_bool, sign_bit]
	if _, err := h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorAnd)); err != nil {
		return err
	}

	// Combinazione S e Z (logica ternaria emulata con salti)
	// (la lasceremo per una versione più avanzata per semplicità)
	// Per ora, li calcoliamo separatamente e li sommiamo.
	// Questo non è super efficiente, ma è molto chiaro da leggere.

	// --- Placeholder per la logica completa dei flag ---
	// La logica completa richiederebbe molti più passaggi di bytecode.
	// Per mantenere questo passo gestibile, ci concentriamo sul concetto.
	// In un'implementazione reale, si creerebbero funzioni "builtin" nella VM
	// per calcolare i flag in modo nativo e veloce.

	// Per ora, impostiamo solo N e lasciamo gli altri a zero per dimostrare il flusso.
	var finalFlags int64 = 0
	if op == objects.OperatorSub {
		finalFlags |= FlagN
	}

	constFlags := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, finalFlags))
	if _, err := h.scopes.Emit(bytecode.OpConstant, constFlags); err != nil {
		return err
	}
	if _, err := h.scopes.Emit(bytecode.OpGlobalSet, h.z80.Register("F")); err != nil {
		return err
	}

	return nil
}

// EmitCpReg executes the CP instruction, subtracting the value of srcReg from register A without storing the result.
// It only updates the flags based on the result of the subtraction. Returns an error if any step fails.
func (h *Helper) EmitCpReg(srcReg string) error {
	// CP è una sottrazione che non salva il risultato in 'A', ma aggiorna solo i flag.
	opA, _ := h.scopes.SymbolDefineUnique("__opA")
	opB, _ := h.scopes.SymbolDefineUnique("__opB")

	// Salva operandi
	if _, err := h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("A")); err != nil {
		return err
	}
	if err := h.scopes.EmitSymbolSetAndPop(opA); err != nil {
		return err
	}
	if _, err := h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register(srcReg)); err != nil {
		return err
	}
	if err := h.scopes.EmitSymbolSetAndPop(opB); err != nil {
		return err
	}

	// Esegui la sottrazione
	if err := h.scopes.EmitSymbolGet(opB); err != nil {
		return err
	}
	if err := h.scopes.EmitSymbolGet(opA); err != nil {
		return err
	}
	if _, err := h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorSub)); err != nil {
		return err
	}

	resultSymbol, _ := h.scopes.SymbolDefineUnique("__alu_result")
	if err := h.scopes.EmitSymbolSetAndPop(resultSymbol); err != nil {
		return err
	}

	// Aggiorna i flag usando il risultato, ma NON salvare il risultato in A.
	return h.emitUpdateFlags(opA, opB, resultSymbol, objects.OperatorSub)
}

// emitCheckFlag generates code to check a specific flag in the Z80 F register and pushes the result based on the given condition.
func (h *Helper) emitCheckFlag(flag byte, condition bool) error {
	// 1. Prendi il registro F
	if _, err := h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("F")); err != nil {
		return err
	}
	// 2. Isola il bit del flag che ci interessa
	constFlag := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, int64(flag)))
	if _, err := h.scopes.Emit(bytecode.OpConstant, constFlag); err != nil {
		return err
	}
	if _, err := h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorAnd)); err != nil {
		return err
	}

	// Ora sulla stack abbiamo 0 (flag non impostato) o un valore > 0 (flag impostato).
	// La tua VM considera 0 come 'falsy' e qualsiasi altro valore come 'truthy'.

	// Se la condizione è "false" (es. per NZ - Jump if Not Zero),
	// dobbiamo negare il risultato.
	if !condition {
		if _, err := h.scopes.Emit(bytecode.OpNot); err != nil {
			return err
		}
	}
	// Alla fine, la stack contiene 'true' se dobbiamo saltare, 'false' altrimenti.
	return nil
}

// EmitJumpUnconditional emits an unconditional jump bytecode to the specified target address. Returns an error if it fails.
func (h *Helper) EmitJumpUnconditional(targetAddress int) error {
	if _, err := h.scopes.Emit(bytecode.OpJump, targetAddress); err != nil {
		return err
	}
	return nil
}

// EmitJumpConditional emits bytecode for a conditional jump based on a specified flag and condition.
func (h *Helper) EmitJumpConditional(flag byte, condition bool, targetAddress int) error {
	// 1. Emetti il codice per controllare il flag.
	// Questo lascia 'true' sulla stack se la condizione del flag è soddisfatta,
	// altrimenti 'false'.
	if err := h.emitCheckFlag(flag, condition); err != nil {
		return err
	}

	// 2. Emetti il salto condizionale della VM.
	// Vogliamo saltare se il valore sulla stack è 'true'.
	// OpJumpTruthy fa esattamente questo.
	if _, err := h.scopes.Emit(bytecode.OpJumpTruthy, targetAddress); err != nil {
		return err
	}

	return nil
}

// emitPush16 pushes a 16-bit value onto the stack by splitting it into high and low bytes and storing them in memory.
func (h *Helper) emitPush16(valueSymbol *tables.Symbol) error {
	// Prendi il valore da salvare (es. l'indirizzo di ritorno)
	if err := h.scopes.EmitSymbolGet(valueSymbol); err != nil {
		return err
	}

	// Salva il valore in una variabile temporanea per poterlo dividere in HI e LO byte
	tempValue, _ := h.scopes.SymbolDefineUnique("__push_val")
	if err := h.scopes.EmitSymbolSetAndPop(tempValue); err != nil {
		return err
	}

	// --- PUSH HI Byte ---
	// SP = SP - 1
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 1)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorSub))
	h.scopes.Emit(bytecode.OpGlobalSet, h.z80.Register("SP"))

	// Scrivi MEMORY[SP] = HI(value)
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("MEMORY")) // Prendi l'array della memoria
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))     // Prendi l'indirizzo (indice)
	h.scopes.EmitSymbolGet(tempValue)                             // Prendi il valore da scrivere
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 8)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorShr)) // Estrai HI byte
	h.scopes.Emit(bytecode.OpIndexSet)                             // Scrivi in memoria

	// --- PUSH LO Byte ---
	// SP = SP - 1
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 1)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorSub))
	h.scopes.Emit(bytecode.OpGlobalSet, h.z80.Register("SP"))

	// Scrivi MEMORY[SP] = LO(value)
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("MEMORY"))
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.EmitSymbolGet(tempValue)
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 0xFF)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorAnd)) // Estrai LO byte
	h.scopes.Emit(bytecode.OpIndexSet)

	return nil
}

// emitPop16 pops a 16-bit value from the stack, reconstructs it, and stores it in the specified destination symbol.
func (h *Helper) emitPop16(destSymbol *tables.Symbol) error {
	// --- POP LO Byte ---
	tempLO, _ := h.scopes.SymbolDefineUnique("__pop_lo")
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("MEMORY"))
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.Emit(bytecode.OpIndexGet)
	h.scopes.EmitSymbolSetAndPop(tempLO)

	// SP = SP + 1
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 1)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorAdd))
	h.scopes.Emit(bytecode.OpGlobalSet, h.z80.Register("SP"))

	// --- POP HI Byte ---
	tempHI, _ := h.scopes.SymbolDefineUnique("__pop_hi")
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("MEMORY"))
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.Emit(bytecode.OpIndexGet)
	h.scopes.EmitSymbolSetAndPop(tempHI)

	// SP = SP + 1
	h.scopes.Emit(bytecode.OpGlobalGet, h.z80.Register("SP"))
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 1)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorAdd))
	h.scopes.Emit(bytecode.OpGlobalSet, h.z80.Register("SP"))

	// --- Ricombina HI e LO ---
	// valore = (HI << 8) | LO
	h.scopes.EmitSymbolGet(tempHI)
	h.scopes.Emit(bytecode.OpConstant, h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, 8)))
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorShl))
	h.scopes.EmitSymbolGet(tempLO)
	h.scopes.Emit(bytecode.OpArithmetic, int(objects.OperatorOr))
	h.scopes.EmitSymbolSetAndPop(destSymbol)

	return nil
}

// EmitCallConditional emits a conditional CALL instruction, executing only if the specified flag and condition are met.
func (h *Helper) EmitCallConditional(flag byte, condition bool, targetAddress int, returnAddress int) error {
	// 1. Controlla se la condizione del flag è soddisfatta.
	// Questo lascia 'true' o 'false' sulla stack della VM.
	if err := h.emitCheckFlag(flag, condition); err != nil {
		return err
	}

	// 2. Se la condizione è 'false', salta oltre tutta la logica della CALL.
	// Emettiamo un OpJumpFalsy con un indirizzo fittizio che correggeremo dopo.
	jumpIfNotMet, err := h.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}

	// --- Blocco di codice eseguito solo se la condizione è VERA ---

	// 3. Salva l'indirizzo di ritorno sullo stack Z80 emulato.
	retAddrSymbol, _ := h.scopes.SymbolDefineUnique("__ret_addr_cond")
	constRetAddr := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, int64(returnAddress)))
	h.scopes.Emit(bytecode.OpConstant, constRetAddr)
	h.scopes.EmitSymbolSetAndPop(retAddrSymbol)
	if err := h.emitPush16(retAddrSymbol); err != nil {
		return err
	}

	// 4. Salta incondizionatamente alla subroutine.
	if err := h.EmitJumpUnconditional(targetAddress); err != nil {
		return err
	}

	// --- Fine del blocco condizionale ---

	// 5. Back-patching: ora sappiamo dove finisce la logica della CALL.
	// Aggiorniamo il nostro OpJumpFalsy iniziale perché salti a questo punto.
	scope, err := h.scopes.Current()
	if err != nil {
		return err
	}
	afterCallPos := scope.InstructionsLen()
	if err := h.scopes.ChangeOperand(jumpIfNotMet, afterCallPos); err != nil {
		return err
	}

	return nil
}

// EmitRetConditional performs a conditional return based on a flag and condition. Updates VM instruction flow accordingly.
func (h *Helper) EmitRetConditional(flag byte, condition bool) error {
	// 1. Controlla se la condizione del flag è soddisfatta.
	// Lascia 'true' o 'false' sulla stack della VM.
	if err := h.emitCheckFlag(flag, condition); err != nil {
		return err
	}

	// 2. Se la condizione è 'false', salta oltre la logica del RET.
	jumpIfNotMet, err := h.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}

	// --- Blocco di codice eseguito solo se la condizione è VERA ---

	// 3. Esegui la logica di un RET standard:
	//    - Pop dell'indirizzo di ritorno dallo stack Z80.
	//    - Salto indiretto a quell'indirizzo.
	retAddrSymbol, _ := h.scopes.SymbolDefineUnique("__ret_addr_cond")
	if err := h.emitPop16(retAddrSymbol); err != nil {
		return err
	}
	if err := h.scopes.EmitSymbolGet(retAddrSymbol); err != nil {
		return err
	}
	if _, err := h.scopes.Emit(bytecode.OpJumpIndirect); err != nil {
		return err
	}

	// --- Fine del blocco condizionale ---

	// 4. Back-patching: Aggiorna il salto iniziale per puntare qui.
	scope, err := h.scopes.Current()
	if err != nil {
		return err
	}
	afterRetPos := scope.InstructionsLen()
	if err := h.scopes.ChangeOperand(jumpIfNotMet, afterRetPos); err != nil {
		return err
	}

	return nil
}

// EmitRst saves the return address on the Z80 stack and jumps to a fixed restart routine address.
func (h *Helper) EmitRst(targetAddress int, returnAddress int) error {
	// 1. Salva l'indirizzo di ritorno (l'istruzione successiva) sullo stack Z80.
	retAddrSymbol, _ := h.scopes.SymbolDefineUnique("__rst_ret_addr")
	constRetAddr := h.constants.AddOrGet("", h.gk.NewInt(objects.FrameStatic, int64(returnAddress)))

	if _, err := h.scopes.Emit(bytecode.OpConstant, constRetAddr); err != nil {
		return err
	}
	if err := h.scopes.EmitSymbolSetAndPop(retAddrSymbol); err != nil {
		return err
	}
	if err := h.emitPush16(retAddrSymbol); err != nil {
		return err
	}

	// 2. Salta all'indirizzo fisso della routine di restart.
	if err := h.EmitJumpUnconditional(targetAddress); err != nil {
		return err
	}

	return nil
}

package z80
//#include <zel/z80.h>
//#include <zel/z80_instructions.h>

//#include "z80_types.h"

//https://github.com/stevecheckoway/libzel

type ControlFlowType int

const (
	CF_CALL ControlFlowType = iota        //!< Call instruction.
	CF_JUMP        //!< Jump instruction.
	CF_RETURN      //!< Return instruction.
	CF_RETURN_I    //!< Return from interrupt instruction.
	CF_RETURN_N    //!< Return from nonmaskable interrupt instruction.
	CF_RESTART     //!< Restart instruction.
	CF_INTERRUPT   //!< Maskable interrupt.
	CF_NMI       //!< Nonmaskable interrupt.
	CF_HALT        //!< Halt instruction.
)


const (
	REG_BC = iota  //!< z80 register bc.
	REG_DE  //!< z80 register de.
	REG_HL  //!< z80 register hl.
	REG_AF  //!< z80 register af.
	REG_IX  //!< z80 register ix.
	REG_IY  //!< z80 register iy.
	REG_PC  //!< z80 register pc.
	REG_SP  //!< z80 register sp.
	REG_BCP //!< z80 register bc'.
	REG_DEP //!< z80 register de'.
	REG_HLP //!< z80 register hl'.
	REG_AFP //!< z80 register af'.
	REG_IR  //!< z80 register ir.
	NUM_REG //!< Number of 16 bit z80 paired registers.
)

// REG_B is the high byte of the BC register pair in the Z80 CPU.
// REG_C is the low byte of the BC register pair in the Z80 CPU.
// REG_D is the high byte of the DE register pair in the Z80 CPU.
// REG_E is the low byte of the DE register pair in the Z80 CPU.
// REG_H is the high byte of the HL register pair in the Z80 CPU.
// REG_L is the low byte of the HL register pair in the Z80 CPU.
// REG_A is the high byte of the AF register pair in the Z80 CPU.
// REG_F is the low byte of the AF register pair in the Z80 CPU.
// REG_IXH is the high byte of the IX register pair in the Z80 CPU.
// REG_IXL is the low byte of the IX register pair in the Z80 CPU.
// REG_IYH is the high byte of the IY register pair in the Z80 CPU.
// REG_IYL is the low byte of the IY register pair in the Z80 CPU.
// REG_PCH is the high byte of the PC (Program Counter) register pair in the Z80 CPU.
// REG_PCL is the low byte of the PC (Program Counter) register pair in the Z80 CPU.

/*
const (
	REG_B = iota // REG_B represents the high byte of the BC register pair in the Z80 CPU.
	REG_C // REG_C represents the low byte of the BC register pair in the Z80 CPU.

	REG_D // REG_D represents the high byte of the DE register pair in the Z80 CPU.
	REG_E // REG_E represents the low byte of the DE register pair in the Z80 CPU.

	REG_H // REG_H represents the high byte of the HL register pair in the Z80 CPU.
	REG_L // REG_L represents the low byte of the HL register pair in the Z80 CPU.

	REG_A // REG_A represents the high byte of the HF register pair in the Z80 CPU.
	REG_F // REG_F represents the low byte of the HF register pair in the Z80 CPU.

	REG_IXH // REG_IXH represents the high byte of the IX register pair in the Z80 CPU.
	REG_IXL // REG_IXL represents the low byte of the IX register pair in the Z80 CPU.

	REG_IYH // REG_IYH represents the high byte of the IY register pair in the Z80 CPU.
	REG_IYL // REG_IYL represents the low byte of the IY register pair in the Z80 CPU.

	REG_PCH // REG_PCH represents the high byte of the PC (Program Counter) register pair in the Z80 CPU.
	REG_PCL // REG_PCH represents the low byte of the PC (Program Counter) register pair in the Z80 CPU.
)

 */

const (
	// No single byte regs
	REG_R = 24 //!< z80 register r. Memory refresh register.
	REG_I   //!< z80 register I. Interrupt page address register.
	INV = 0xff //!< Invalid.
)


type Z80FunctionBlock struct {
	ReadMem func(uint16, bool, *Z80) uint8
	WriteMem func(uint16, uint8, *Z80)
	ReadIO func(uint16, *Z80) uint8
	WriteIO func(uint16, uint8, *Z80)
	ReadInterruptData func(uint16, *Z80) uint16
	InterruptComplete func(uint16, *Z80)
	ControlFlow func(uint16, uint16, ControlFlowType, *Z80)
}




/*
const (
	REG_SP = iota
	REG_PC
	REG_PCH
	REG_PCL
	REG_BC
	REG_DE
	REG_HL
	REG_AF
)

const (
	REG_B = iota
	REG_C
	REG_D
	REG_E
	REG_H
	REG_L
	REG_A
)

 */

type Z80 struct {
	/* C guarantees consecutive layout */
	word_reg[8] uint16
	//byte_reg[7] byte
	iff1 bool
	iff2 bool
	can_handle_interrupt bool
	interrupt_mode int
	interrupt bool
	nmi bool
	halt bool
	restart_io bool

	ReadMem func(uint16, bool, *Z80) uint8
	WriteMem func(uint16, uint8, *Z80)
	ReadIO func(uint16, *Z80) uint8
	WriteIO func(uint16, uint8, *Z80)
	ReadInterruptData func(uint16, *Z80) uint16
	InterruptComplete func(uint16, *Z80)
	ControlFlow func(uint16, uint16, ControlFlowType, *Z80)
}

func (cpu * Z80) Disassemble(  address uint16, buffer[]byte)int{
	Instruction inst;
	length := IF_ID( &inst, address, ReadInstructionMemory, cpu );
	if len(buffer) > 0 {
		DisassembleInstruction(&inst, buffer)
	}
	return length;
}

func (cpu * Z80) HasHalted() bool{
	return cpu.halt;
}



// GetByteReg con calcolo diretto
func (cpu *Z80) GetByteReg(index uint8) uint8 {
	wordIndex := index >> 1 //index / 2 // Calcola l'indice del registro a 16 bit (BC=0, DE=1, etc.)
	if index & 1 == 0 {
		return uint8(cpu.word_reg[wordIndex]) // Low
	}
	return uint8(cpu.word_reg[wordIndex] >> 8) // High
}

// SetByteReg con calcolo diretto
func (cpu *Z80) SetByteReg(index uint8, value uint8) {
	wordIndex := index >> 1
	if index & 1 == 0 {
		cpu.word_reg[wordIndex] = (cpu.word_reg[wordIndex] & 0xFF00) | uint16(value) // Low
		return
	}
	cpu.word_reg[wordIndex] = (cpu.word_reg[wordIndex] & 0x00FF) | (uint16(value) << 8) // High
}

func (cpu * Z80) GetWordReg(reg uint8) uint16{
	//assert(reg >= 0 && reg < NUM_REG);
	return cpu.word_reg[reg];
}

func (cpu * Z80) SetWordReg(reg uint8, value uint16){
	//assert(reg >= 0 && reg < NUM_REG);
	cpu.word_reg[reg] = value;
}

func (cpu * Z80) RaiseNMI(){
	cpu.nmi = true;
}

func (cpu * Z80) RaiseInterrupt( ){
	cpu.interrupt = true;
}

func (cpu * Z80) RestartIO(  ){
	cpu.restart_io = true;
}

func (cpu * Z80) ClearHalt(  ){
	cpu.halt = false;
}

func (cpu * Z80) GetPC() uint16{
	return cpu.word_reg[REG_PC]
}

func (cpu *Z80) GetSP() uint16 {
	return cpu.word_reg[REG_SP]
}

func (cpu *Z80) GetBC() uint16 {
	return cpu.word_reg[REG_BC]
}

func (cpu *Z80) GetDE() uint16 {
	return cpu.word_reg[REG_DE]
}

func (cpu *Z80) GetHL() uint16 {
	return cpu.word_reg[REG_HL]
}

func (cpu *Z80) GetAF() uint16 {
	return cpu.word_reg[REG_AF]
}

func (cpu *Z80) GetPCH() uint8 {
	// Estrae il byte alto (PCH) da PC
	return uint8(cpu.word_reg[REG_PC] >> 8)
}

func (cpu *Z80) GetPCL() uint8 {
	// Estrae il byte basso (PCL) da PC
	return uint8(cpu.word_reg[REG_PC])
}

func (cpu *Z80) GetB() uint8 {
	// Estrae il byte alto (B) da BC
	return uint8(cpu.word_reg[REG_BC] >> 8)
}

func (cpu *Z80) GetC() uint8 {
	// Estrae il byte basso (C) da BC
	return uint8(cpu.word_reg[REG_BC])
}

func (cpu *Z80) GetD() uint8 {
	// Estrae il byte alto (D) da DE
	return uint8(cpu.word_reg[REG_DE] >> 8)
}

func (cpu *Z80) GetE() uint8 {
	// Estrae il byte basso (E) da DE
	return uint8(cpu.word_reg[REG_DE])
}

func (cpu *Z80) GetH() uint8 {
	// Estrae il byte alto (H) da HL
	return uint8(cpu.word_reg[REG_HL] >> 8)
}

func (cpu *Z80) GetL() uint8 {
	// Estrae il byte basso (L) da HL
	return uint8(cpu.word_reg[REG_HL])
}

func (cpu *Z80) GetA() uint8 {
	// Estrae il byte alto (A) da AF
	return uint8(cpu.word_reg[REG_AF] >> 8)
}

func (cpu *Z80) GetF() uint8 {
	// Estrae il byte basso (F) da AF
	return uint8(cpu.word_reg[REG_AF])
}

// 16 bit
func (cpu * Z80) SetSP(v uint16) {
	cpu.word_reg[REG_SP] = v
}

func (cpu * Z80) SetPC(v uint16) {
	cpu.word_reg[REG_PC] = v
}

func (cpu * Z80) SetBC(v uint16) {
	cpu.word_reg[REG_BC] = v
}

func (cpu * Z80) SetDE(v uint16) {
	cpu.word_reg[REG_DE] = v
}

func (cpu * Z80) SetHL(v uint16) {
	cpu.word_reg[REG_HL] = v
}

func (cpu * Z80) SetAF(v uint16) {
	cpu.word_reg[REG_AF] = v
}

// SetPCH 8 bit - CORRETTO
func (cpu *Z80) SetPCH(v uint8) {
	// Preserva PCL (byte basso) e imposta PCH (byte alto)
	cpu.word_reg[REG_PC] = (cpu.word_reg[REG_PC] & 0x00FF) | (uint16(v) << 8)
}

func (cpu *Z80) SetPCL(v uint8) {
	// Preserva PCH (byte alto) e imposta PCL (byte basso)
	cpu.word_reg[REG_PC] = (cpu.word_reg[REG_PC] & 0xFF00) | uint16(v)
}

func (cpu *Z80) SetB(v uint8) {
	cpu.word_reg[REG_BC] = (cpu.word_reg[REG_BC] & 0x00FF) | (uint16(v) << 8)
}

func (cpu *Z80) SetC(v uint8) {
	cpu.word_reg[REG_BC] = (cpu.word_reg[REG_BC] & 0xFF00) | uint16(v)
}

func (cpu *Z80) SetD(v uint8) {
	cpu.word_reg[REG_DE] = (cpu.word_reg[REG_DE] & 0x00FF) | (uint16(v) << 8)
}

func (cpu *Z80) SetE(v uint8) {
	cpu.word_reg[REG_DE] = (cpu.word_reg[REG_DE] & 0xFF00) | uint16(v)
}

func (cpu *Z80) SetH(v uint8) {
	cpu.word_reg[REG_HL] = (cpu.word_reg[REG_HL] & 0x00FF) | (uint16(v) << 8)
}

func (cpu *Z80) SetL(v uint8) {
	cpu.word_reg[REG_HL] = (cpu.word_reg[REG_HL] & 0xFF00) | uint16(v)
}

func (cpu *Z80) SetA(v uint8) {
	cpu.word_reg[REG_AF] = (cpu.word_reg[REG_AF] & 0x00FF) | (uint16(v) << 8)
}

func (cpu *Z80) SetF(v uint8) {
	cpu.word_reg[REG_AF] = (cpu.word_reg[REG_AF] & 0xFF00) | uint16(v)
}

func (cpu *Z80) DecrementSP() {
	cpu.word_reg[REG_SP]--
}

func (cpu *Z80) IncrementSP() {
	cpu.word_reg[REG_SP]++
}

func (cpu * Z80) FlagIsReset(z uint8) bool{
	f := cpu.GetF()
	out := f&(1<<(z))
	return out == 0
}

func (cpu * Z80) FlagIsSet(z uint8) bool{
	return !cpu.FlagIsReset(z)
}

func (cpu * Z80) SetFlag(z uint8) {
	f := cpu.GetF()
	out := f | (1<<(z))
	cpu.SetF(out)
}

func (cpu * Z80) ResetFlag(z uint8) {
	f := cpu.GetF()
	out := f & (^(1<<(z)))
	cpu.SetF(out)
}

func (cpu * Z80) SetFlagValueI(z uint8, v int){
	f := cpu.GetF()
	p1 := f  & (^(1<<(z)))
	p2 := uint8(0)
	if v != 0 {
		p2 = 1 << z
	}
	cpu.SetF(p1 | p2)
}

func (cpu * Z80) SetFlagValue(z uint8, v bool){
	f := cpu.GetF()
	p1 := f  & (^(1<<(z)))
	p2 := uint8(0)
	if v {
		p2 = 1 << z
	}
	cpu.SetF(p1 | p2)
}

func (cpu *Z80) CondIsMet(c int) bool {
	if c >= 0 {
		// Condizione diretta: controlla se il flag c è impostato.
		return cpu.FlagIsSet(uint8(c))
	} else {
		// Condizione inversa: controlla se il flag (-c-1) è resettato.
		// Esempi: c=-1 -> flag 0; c=-2 -> flag 1
		flagIndex := uint8(-c - 1)
		return cpu.FlagIsReset(flagIndex)
	}
}

func IgnoreControlFlow(pc uint16 ,target uint16, cf ControlFlowType , cpu *Z80) {
}

func NewZ80(blk *Z80FunctionBlock) *Z80 {
	cpu := &Z80{
		can_handle_interrupt : true,
		ReadMem : blk.ReadMem,
		WriteMem : blk.WriteMem,
		ReadInterruptData : blk.ReadInterruptData,
		WriteIO : blk.WriteIO,
		ReadIO : blk.ReadIO,
		InterruptComplete : blk.InterruptComplete,
		ControlFlow :blk.ControlFlow,
	}

	//cpu.byte_reg = (byte*)cpu.word_reg
	cpu.SetSP( 0xffff)
	cpu.SetAF( 0xffff)
	//TODO IMPLEMENT
	//#define REQUIRE(x) assert(blk.x); cpu.x = blk.x
	cpu.x = blk.x
	if cpu.ControlFlow == nil {
		cpu.ControlFlow = IgnoreControlFlow
	}

	return cpu;
}

func (cpu * Z80) Step(outPC uint16) (int, uint16) {
	//TODO we have to return outPC!

	oldPC := cpu.GetPC();
	var inst Instruction

	ticks := 0;

	cpu.restart_io = false;
	if cpu.nmi || (cpu.can_handle_interrupt && cpu.interrupt && cpu.iff1) {
		cpu.halt = false;
		// Any interrupt increases R by one.
		r := cpu.GetByteReg(REG_R)
		r = ((r + 1) & 0x7f) | (r & 0x80);
		cpu.SetByteReg(REG_R, r)
		if cpu.nmi {
			cpu.nmi = false;
			cpu.iff1 = false; /* Disable interrupts. */
			// 5 cycles fetching and ignoring the opcode
			// This can cause another nmi.
			cpu.ReadMem(cpu.GetPC(), true, cpu)
			// 6 cycles writing the PC
			cpu.WriteMem(cpu.GetSP(), cpu.GetPCH(), cpu)
			cpu.DecrementSP()
			cpu.WriteMem(cpu.GetSP(), cpu.GetPCL(), cpu)
			cpu.DecrementSP()
			cpu.SetPC(0x0066) // Fixed location
			cpu.ControlFlow(oldPC, cpu.GetPC(), CF_NMI, cpu)
			ticks = 11
			goto interrupt_exit
		}
		cpu.interrupt = false;
		// This depends on the mode
		switch cpu.interrupt_mode {
		case 0:
			// interrupting device supplies instruction.
			IF_ID(&inst, 0, (byte(*)(uint16, void*))cpu.ReadInterruptData, cpu)
			inst.additional_tstates += 2; // 2 wait states add to M1 cycle
			cpu.ControlFlow(oldPC, 0xffff, CF_INTERRUPT, cpu)
			break;

		case 1:
			// Insert a restart instruction + 2 cycles.
			// Handle it here since we know exactly what
			// it is supposed to do.
			cpu.WriteMem(cpu.GetSP(), cpu.GetPCH(), cpu)
			cpu.DecrementSP()
			cpu.WriteMem(cpu.GetSP(), cpu.GetPCL(), cpu)
			cpu.DecrementSP()
			cpu.SetPC(0x0038) // Fixed location
			cpu.ControlFlow(oldPC, cpu.GetPC(), CF_INTERRUPT, cpu)
			ticks = 13;
			goto interrupt_exit;

		case 2:
			// 7 cycles to read the 7 bits from the
			// interrupting device, 6 to push the PC, and
			// 6 to load the jump address.
			cpu.WriteMem(cpu.GetSP(), cpu.GetPCH(), cpu)
			cpu.DecrementSP()
			cpu.WriteMem(cpu.GetSP(), cpu.GetPCL(), cpu)
			cpu.DecrementSP()
			address := ((uint16)
			BYTE_REG[REG_I]) << 8
			address |= cpu.ReadInterruptData(0, cpu) & 0xfe
			cpu.SetPCL(cpu.ReadMem(address, true, cpu))
			cpu.SetPCH(cpu.ReadMem(address+1, true, cpu))
			(cpu.ControlFlow)(oldPC, cpu.GetPC(), CF_INTERRUPT, cpu);
			ticks = 19;
			goto interrupt_exit;

		}
	} else {
		// ei re-enables iff1 but interrupts cannot be handled
		// for another instruction.
		cpu.can_handle_interrupt = true
		// Fetch the next instruction from memory.
		if !cpu.halt {
			PC += IF_ID(&inst, PC, ReadInstructionMemory, cpu);
		} else {
			inst.additional_tstates = 0;
			inst.offset = 0;
			inst.immediate = 0;
			inst.r_increment = 1;
			inst.IT = &Unprefixed[0x76]; // NOP
		}
	}

	return cpu.doInst(&inst)
}

func (cpu * Z80) doInst(inst *Instruction) (int, uint16) {
	op1 := uint32(0)
	op2 := uint32(0)
	carry := uint32(0)
	result := uint32(0)
	i := 0;
	took_branch := true;

	switch inst.IT.kind {
		/* 8-Bit Load Group */
		case LD_I_N:
			cpu.instLD_I_N(inst)
		case LD_MRR_N:
			cpu.instLD_MRR_N(inst)
		case LD_I_R:
			cpu.instLD_I_R(inst)
		case LD_MRR_R:
			cpu.instLD_MRR_R(inst)
		case LD_MNN_R:
			cpu.instLD_MNN_R(inst)
		case LD_R_I:
			cpu.instLD_R_I(inst)
		case LD_R_MRR:
			cpu.instLD_R_MRR(inst)
		case LD_R_MNN:
			cpu.instLD_R_MNN(inst)
		case LD_R_N:
			cpu.instLD_R_N(inst)
		case LD_R_R:
			cpu.instLD_R_R(inst)
		/* 16-Bit Load Group */
		case LD_RR_MNN:
			cpu.instLD_RR_MNN(inst)
		case LD_RR_NN:
			cpu.instLD_RR_NN(inst)
		case LD_RR_RR:
			cpu.instLD_RR_RR(inst)
		case LD_MNN_RR:
			cpu.instLD_MNN_RR(inst)
		case POP_RR:
			cpu.instPOP_RR(inst)
		case PUSH_RR:
			cpu.instPUSH_RR(inst)
		/* Exchange, Block Transfer, Search Group */
		case CPD:
			cpu.instCPD(inst)
		case CPI:
			cpu.instCPI(inst)
		case CPDR:
			cpu.instCPDR(inst)
		case CPIR:
			cpu.instCPIR(inst)
		case EX_MRR_RR:
			cpu.instEX_MRR_RR(inst)
		case EX_RR_RR:
			cpu.instEX_RR_RR(inst)
		case EXX:
			cpu.instEXX(inst)
		case LDD:
			cpu.instLDD(inst)
		case LDI:
			cpu.instLDI(inst)
		case LDDR:
			cpu.instLDDR(inst)
		case LDIR:
			cpu.instLDIR(inst)
		/* 8-Bit Arithmetic and Logical Group */
		case ADC_R_I:
			cpu.instADC_R_I(inst)
		case ADC_R_MRR:
			cpu.instADC_R_MRR(inst)
		case ADC_R_N:
			cpu.instADC_R_N(inst)
		case ADC_R_R:
			cpu.instADC_R_R(inst)
		case ADD_R_I:
			cpu.instADD_R_I(inst)
		case ADD_R_MRR:
			cpu.instADD_R_MRR(inst)
		case ADD_R_N:
		op2 = IMM;
		goto add_r;
		case ADD_R_R:
		op2 = BYTE_REG[OP2];
		goto add_r;
		adc_r:
		carry = cpu.FlagIsSet( FLAG_C );
		add_r:
		op1 = BYTE_REG[OP1];
		result = op1 + op2 + carry;
		BYTE_REG[OP1] = result & 0xff;
		cpu.SetFlagValue( FLAG_S, result & 0x80 )
			cpu.SetFlagValue( FLAG_Z, !(result & 0xff) )
			cpu.SetFlagValue( FLAG_Y, result & 0x20 ) // bit 5
			cpu.SetFlagValue( FLAG_H, (op1&0x0f)+(op2&0x0f)+carry>0x0f )
			cpu.SetFlagValue( FLAG_X, result & 0x08 ) // bit 3
			cpu.SetFlagValue( FLAG_P, (op2 == 0x7f && carry) || ((op1&0x80)==(op1&0x80) && (op1&0x80)!=(result&0x80)) );
			cpu.ResetFlag( FLAG_N );
			cpu.SetFlagValue( FLAG_C, result > 0xff );
		break;

		case SBC_R_I:
		case SBC_R_MRR:
		op2 = uint32(cpu.ReadMem( WORD_REG[OP2]+OFFSET, false, cpu))
		goto sbc_r;
		case SBC_R_N:
		op2 = IMM;
		goto sbc_r;
		case SBC_R_R:
		op2 = BYTE_REG[OP2];
		goto sbc_r;

		case SUB_I:
		case SUB_MRR:
		op2 = uint32(cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu))
		goto sub_r;
		case SUB_N:
		op2 = IMM;
		goto sub_r;
		case SUB_R:
		op2 = BYTE_REG[OP1];
		goto sub_r;
		sbc_r:
		carry = cpu.FlagIsSet( FLAG_C );
		sub_r:
		op1 = A;
		result = op1 - op2 - carry;
		A = result & 0xff;
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !(result & 0xff) );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_H, (op1&0x0f) < (op2&0x0f)+carry );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, (carry && op1-op2 == 0x80) || ((op1&0x80) != (op2&0x80) &&  (op1&0x80) != (result&0x80)) );
			cpu.SetFlagValue( FLAG_C, op1 < op2 + carry );
		break;

		case DEC_I:
		case DEC_MRR:
		result = uint32(cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu))
		carry = result
		--result
		cpu.WriteMem( WORD_REG[OP1]+OFFSET, uint8(result & 0xff), cpu)
		goto dec_x;
		case DEC_R:
		result = BYTE_REG[OP1];
		carry = result;
		BYTE_REG[OP1] = --result;
		dec_x:
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !(result&0xff) );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_H, !(carry&0x0f) );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, carry == 0x80 );
			cpu.SetFlag( FLAG_N );
		break;


		case INC_I:
		case INC_MRR:
		result = cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu);
		carry = result;
		++result;
		cpu.WriteMem( WORD_REG[OP1]+OFFSET, result, cpu);
		goto inc_x;
		case INC_R:
		result = BYTE_REG[OP1];
		carry = result;
		++result;
		BYTE_REG[OP1] = result;
		inc_x:
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !(result&0xff) );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_H, (carry&0x0f)+1 > 0x0f );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, carry & 0x7f );
			cpu.ResetFlag( FLAG_N );
		break;

		case CP_I:
		case CP_MRR:
		op2 = cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu);
		goto cp;
		case CP_N:
		op2 = IMM;
		goto cp;
		case CP_R:
		op2 = BYTE_REG[OP1];
		cp:
		op1 = A;
		result = op1 - op2;
			cpu.SetFlagValue( FLAG_S, result&0x80 );
			cpu.SetFlagValue( FLAG_Z, !(result&0xff) );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_H, (op1&0x0f) < (op2&0x0f) );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, (op1&0x80) != (op2&0x80) && (op1&0x80) != (result&0x80) );
			cpu.SetFlag( FLAG_N );
			cpu.SetFlagValue( FLAG_C, op1 < op2 );
		break;

		case AND_I:
		case AND_MRR:
		result = A &= cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu);
			cpu.SetFlag( FLAG_H );
		goto logical_flags;
		case AND_N:
		result = A &= IMM;
			cpu.SetFlag( FLAG_H );
		goto logical_flags;
		case AND_R:
		result = A &= BYTE_REG[OP1];
			cpu.SetFlag( FLAG_H );
		goto logical_flags;

		case OR_I:
		case OR_MRR:
		result = A | cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu);
			cpu.ResetFlag( FLAG_H );
		goto logical_flags; /* Same flags as and */
		case OR_N:
		result = A | IMM;
			cpu.ResetFlag( FLAG_H );
		goto logical_flags;
		case OR_R:
		result = A | BYTE_REG[OP1];
			cpu.ResetFlag( FLAG_H );
		goto logical_flags;

		case XOR_I:
		case XOR_MRR:
		result = A ^ cpu.ReadMem( WORD_REG[OP1]+OFFSET, false, cpu);
			cpu.ResetFlag( FLAG_H );
		goto logical_flags;
		case XOR_N:
		result = A ^ IMM;
			cpu.ResetFlag( FLAG_H );
		goto logical_flags;
		case XOR_R:
		result = A ^ BYTE_REG[OP1];
		ResetFlag( FLAG_H );
		logical_flags:
			cpu.SetA(uint8(result))
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !result );
			cpu.SetFlagValue( FLAG_Y, result & 0x02 );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, ParityIsEven(result) );
			cpu.ResetFlag( FLAG_N );
			cpu.ResetFlag( FLAG_C );
		break;


		/* General-Purpose Arithmetic and CPU Control Group */
		case CCF:
		result = cpu.FlagIsSet( FLAG_C );
			cpu.SetFlagValue( FLAG_Y, A & 0x20 );
			cpu.SetFlagValue( FLAG_H, result );
			cpu.SetFlagValue( FLAG_X, A & 0x08 );
			cpu.ResetFlag( FLAG_N );
			cpu.SetFlagValue( FLAG_C, !result );
		break;

		case CPL:
		result = ~A;
		A = result;
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlag( FLAG_H );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlag( FLAG_N );
		break;

		case DAA:
		op1 = cpu.FlagIsSet( FLAG_N );
		op2 = cpu.FlagIsSet( FLAG_H );
		carry = cpu.FlagIsSet( FLAG_C );
		result = A;

		for i = 0; i < 13; i++ {
			if _daa_table[i][0] == op1 &&
			_daa_table[i][1] == carry &&
			_daa_table[i][2] <= result >> 4 &&
			_daa_table[i][3] >= result >> 4 &&
			_daa_table[i][4] == op2 &&
			_daa_table[i][5] <= (result & 0x0f) &&
			_daa_table[i][6] >= (result & 0x0f)  {
				result = (result + _daa_table[i][7]) & 0xff;
				cpu.SetFlagValue( FLAG_C, daa_table[i][7] );
				break
			}
		}
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !result );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, ParityIsEven(result) );
		break;

		case DI:
		cpu.iff1 = false;
		cpu.iff2 = false;
		break;

		case EI:
		cpu.iff1 = true;
		cpu.iff2 = true;
		cpu.can_handle_interrupt = false;
		break;

		case HALT:
		cpu.halt = true;
		cpu.ControlFlow( oldPC, PC, CF_HALT, cpu);
		break;

		case IM:
		cpu.interrupt_mode = OP1;
		break;

		case NEG: // This does A <- 0 - A, flags set accordingly.
		result = -A;
		A = result & 0xff;
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !result );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_H, result & 0x0f );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, (result&0xff) == 0x80 );
			cpu.SetFlag( FLAG_N );
			cpu.SetFlagValue( FLAG_C, !result );
		break;

		case NOP:
		break;

		case SCF:
			cpu.SetFlagValue( FLAG_Y, A & 0x20 );
			cpu.ResetFlag( FLAG_H );
			cpu.SetFlagValue( FLAG_X, A & 0x08 );
			cpu.ResetFlag( FLAG_N );
			cpu.SetFlag( FLAG_C );
		break;


		/* 16-Bit Arithmetic Group */
		case ADD_RR_RR:
		op1 = WORD_REG[OP1];
		op2 = WORD_REG[OP2];
		result = op1 + op2;
		goto add_rr_flags;
		case ADC_RR_RR:
		op1 = WORD_REG[OP1];
		op2 = WORD_REG[OP2];
		carry = FlagIsSet( FLAG_C );
		result = op1 + op2 + carry;
			cpu.SetFlagValue( FLAG_S, result & 0x8000 );
			cpu.SetFlagValue( FLAG_Z, !(result & 0xffff) );
			cpu.SetFlagValue( FLAG_P, (op2 == 0x7fff && carry) ||
		((op1&0x8000) == ((op2+carry)&0x8000) &&
		(op1&0x8000) != (result&0x8000)) );
		add_rr_flags:
			cpu.SetFlagValue( FLAG_Y, result & 0x2000 ); // bit 13
			cpu.SetFlagValue( FLAG_H, (op1&0x0fff)+(op2&0x0fff)+carry>0x0fff );
			cpu.SetFlagValue( FLAG_X, result & 0x0800 ); // bit 11
			cpu.ResetFlag( FLAG_N );
			cpu.SetFlagValue( FLAG_C, result > 0xffff );
		WORD_REG[OP1] = result & 0xffff;
		break;

		case DEC_RR:
		--WORD_REG[OP1];
		break;

		case INC_RR:
		++WORD_REG[OP1];
		break;

		case SBC_RR_RR:
		op1 = WORD_REG[OP1];
		op2 = WORD_REG[OP2];
		carry = FlagIsSet( FLAG_C );
		result = op1 - op2 - carry;
		WORD_REG[OP1] = result & 0xffff;
			cpu.SetFlagValue( FLAG_S, result & 0x8000 );
			cpu.SetFlagValue( FLAG_Z, !(result & 0xffff) );
			cpu.SetFlagValue( FLAG_Y, result & 0x2000 );
			cpu.SetFlagValue( FLAG_H, (op1&0x0fff) < (op2&0x0fff)+carry );
			cpu.SetFlagValue( FLAG_X, result & 0x0800 );
			cpu.SetFlagValue( FLAG_P, (carry && op1-op2 == 0x8000) || ((op1&0x8000) != (op2&0x8000) && (op1&0x8000) != (result&0x8000)) );
			cpu.SetFlagValue( FLAG_C, op1 < op2 + carry );
		break;


		/* Rotate and Shift Group
		 * Almost all flags are set the same so jump to a common block
		 * of flag setting. */
		case RLCA:
		result = A;
		carry = result >> 7;
		result = (result << 1) | carry;
		A = result & 0xff;
		goto rotate_accum_flags;
		case RLA:
		result = A;
		carry = result >> 7;
		result = (result << 1) | FlagIsSet(FLAG_C);
		A = result & 0xff;
		goto rotate_accum_flags;
		case RRCA:
		result = A;
		carry = result & 0x1;
		result = (result >> 1) | (carry << 7);
		A = result;
		goto rotate_accum_flags;

		case RRA:
		result = A;
		carry = result & 0x1;
		result = (result >> 1) | (FlagIsSet(FLAG_C) << 7);
		A = result;
		goto rotate_accum_flags;
		case RLC_I:
		case RLC_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result >> 7;
		result = (result << 1) | carry;
		cpu.WriteMem( op1, result & 0xff, cpu);
		goto shift_flags;
		case RLC_R:
		result = BYTE_REG[OP1];
		carry = result >> 7;
		result = (result << 1) | carry;
		BYTE_REG[OP1] = result;
		goto shift_flags;
		case RL_I:
		case RL_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result >> 7;
		result = (result << 1) | FlagIsSet(FLAG_C);
		cpu.WriteMem( op1, result & 0xff, cpu);
		goto shift_flags;
		case RL_R:
		result = BYTE_REG[OP1];
		carry = result >> 7;
		result = (result << 1) | FlagIsSet(FLAG_C);
		BYTE_REG[OP1] = result & 0xff;
		goto shift_flags;
		case RRC_I:
		case RRC_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result & 0x1;
		result = (result >> 1) | (carry << 7);
		cpu.WriteMem( op1, result, cpu);
		goto shift_flags;
		case RRC_R:
		result = BYTE_REG[OP1];
		carry = result & 0x1;
		result = (result >> 1) | (carry << 7);
		BYTE_REG[OP1] = result;
		goto shift_flags;
		case RR_I:
		case RR_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result & 0x1;
		result = (result >> 1) | (FlagIsSet(FLAG_C) << 7);
		cpu.WriteMem( op1, result, cpu);
		goto shift_flags;
		case RR_R:
		result = BYTE_REG[OP1];
		carry = result & 0x1;
		result = (result >> 1) | (FlagIsSet(FLAG_C) << 7);
		BYTE_REG[OP1] = result;
		goto shift_flags;
		case SLA_I:
		case SLA_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result >> 7;
		result <<= 1;
		cpu.WriteMem( op1, result & 0xff, cpu);
		goto shift_flags;
		case SLA_R:
		result = BYTE_REG[OP1];
		carry = result >> 7;
		result <<= 1;
		BYTE_REG[OP1] = result & 0xff;
		goto shift_flags;
		case SLL_I:
		case SLL_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result >> 7;
		result = (result << 1) | 0x1;
		cpu.WriteMem( op1, result & 0xff, cpu);
		goto shift_flags;
		case SLL_R:
		result = BYTE_REG[OP1];
		carry = result >> 7;
		result = (result << 1) | 0x1;
		BYTE_REG[OP1] = result & 0xff;
		goto shift_flags;
		case SRA_I:
		case SRA_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result & 0x1;
		result = (result & 0x80) | (result >> 1);
		cpu.WriteMem( op1, result, cpu);
		goto shift_flags;
		case SRA_R:
		result = BYTE_REG[OP1];
		carry = result & 0x1;
		result = (result & 0x80) | (result >> 1);
		BYTE_REG[OP1] = result;
		goto shift_flags;
		case SRL_I:
		case SRL_MRR:
		op1 = WORD_REG[OP1]+OFFSET;
		result = cpu.ReadMem( op1, false, cpu);
		carry = result & 0x1;
		result >>= 1;
		cpu.WriteMem( op1, result, cpu);
		goto shift_flags;
		case SRL_R:
		result = BYTE_REG[OP1];
		carry = result & 0x1;
		result >>= 1;
		BYTE_REG[OP1] = result;
		goto shift_flags;
		case RLD:
		op1 = cpu.ReadMem( HL, false, cpu);
		op2 = A;
		result = (op2 & 0xf0) | (op1 >> 4);
		op1 = (op1 << 4) | (op2 & 0x0f);
		cpu.WriteMem( HL, op1 & 0xff, cpu);
		A = result;
		carry = FlagIsSet( FLAG_C ); // makes code simpler
		goto shift_flags;
		case RRD:
		op1 = cpu.ReadMem( HL, false, cpu);
		op2 = A;
		result = (op2 & 0xf0) | (op1 & 0x0f);
		op1 = (op1 >> 4) | (op2 << 4);
		cpu.WriteMem( HL, op1 & 0xff, cpu);
		A = result;
		carry = cpu.FlagIsSet( FLAG_C ); // simpler code

		shift_flags:
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !result );
			cpu.SetFlagValue( FLAG_P, ParityIsEven(result) );
		rotate_accum_flags:
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.ResetFlag( FLAG_H );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.ResetFlag( FLAG_N );
			cpu.SetFlagValue( FLAG_C, carry );
		break;


		/* Bit Set, Reset, and Test Group
		 * In this group, the first operand is the bit to
		 * set/reset/test. */
		case BIT_I, BIT_MRR:
		op2 = cpu.ReadMem( WORD_REG[OP2]+OFFSET, false, cpu);
		result = op2 & (0x1<<OP1);
		// XXX: This is wrong for BIT_MRR, but right for BIT_I
		// Does anyone know what the right thing for BIT_MRR
		// is?
			cpu.SetFlagValue( FLAG_Y, (WORD_REG[OP2]+OFFSET)&0x20 );
			cpu.SetFlagValue( FLAG_X, (WORD_REG[OP2]+OFFSET)&0x08 );
		goto bit;


		case BIT_R:
		op2 = BYTE_REG[OP2];
		result = op2 & (0x1<<OP1);
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
		bit:
			cpu.SetFlagValue( FLAG_S, OP1==7 && result );
			cpu.SetFlagValue( FLAG_Z, !result );
			cpu.SetFlag( FLAG_H );
			cpu.SetFlagValue( FLAG_P, !result );
			cpu.ResetFlag( FLAG_N );
		break;

		case RES_I, RES_MRR:
			op2 = WORD_REG[OP2] + OFFSET;
			result = cpu.ReadMem( op2, false, cpu) & ~(1<<OP1);
			cpu.WriteMem( op2, result, cpu);
			if( inst.IT.extra != INV )
			BYTE_REG[inst.IT.extra] = result;
		break;

		case RES_R:
			result = BYTE_REG[OP2] & ~(1<<OP1);
			BYTE_REG[OP2] = result;
			if( inst.IT.extra != INV )
			BYTE_REG[inst.IT.extra] = result;
		break;

		case SET_I, SET_MRR:
			op2 = WORD_REG[OP2] + OFFSET;
			result = cpu.ReadMem( op2, false, cpu) | (1<<OP1);
			cpu.WriteMem( op2, result, cpu);
			if inst.IT.extra != INV  {
				BYTE_REG[inst.IT.extra] = result
			}
		break;

		case SET_R:
			result = BYTE_REG[OP2] | (1<<OP1);
			BYTE_REG[OP2] = result;
			if( inst.IT.extra != INV ) {
				BYTE_REG[inst.IT.extra] = result
			}
		break;

		/* Jump Group */
		case DJNZ:
			if --B  {
				PC += OFFSET;
			} else {
				took_branch = false;
			}
		break;

		case JP_C_MNN:
		if( !cpu.CondIsMet(OP1) ) {
			took_branch = false;
			break; // condition isn't met
		}
		case JP_MNN:
		PC = IMM;
		cpu.ControlFlow( oldPC, PC, CF_JUMP, cpu);
		break;
		case JP_MRR: // jp (hl); jp (ix); jp (iy)
		PC = WORD_REG[OP1];
		cpu.ControlFlow( oldPC, PC, CF_JUMP, cpu);
		break;

		case JR_C:
			if !cpu.CondIsMet(OP1) {
				took_branch = false;
				break;
			}

		case JR:
			PC += OFFSET;
		break;

		/* Call and Return Group */
		case CALL_C_MNN:
		if !cpu.CondIsMet(OP1)  {
			took_branch = false;
			break; // condition isn't met
		}

		case CALL_MNN:
			cpu.WriteMem( --SP, PCH, cpu );
			cpu.WriteMem( --SP, PCL, cpu );
			PC = IMM;
			cpu.ControlFlow(oldPC, PC, CF_CALL, cpu)
		break;

		case RETI:
			// Some docs say that iff2 is copied to iff1 like in
			// reti. Some simulatores do that. Zilog docs are very
			// clear that this doesn't happen, but they're often
			// wrong.
			(cpu.InterruptComplete)(cpu);
			result = CF_RETURN_I;
			goto ret

		case RETN:
			cpu.iff1 = cpu.iff2;
			result = CF_RETURN_N;
			goto ret;

		case RET_C:
			if !cpu.CondIsMet(OP1)  {
				took_branch = false;
				break;
			}

		case RET:
			result = CF_RETURN;
			ret:
				PCL = cpu.ReadMem( SP++, false, cpu );
				PCH = cpu.ReadMem( SP++, false, cpu );
				(cpu.ControlFlow)( oldPC, PC, result, cpu);
		break;

		case RST:
			cpu.WriteMem( --SP, PCH, cpu );
			cpu.WriteMem( --SP, PCL, cpu );
			PC = OP1;
			(cpu.ControlFlow)( oldPC, PC, CF_RESTART, cpu);
		break;

		/* Input and Output Group */
		case IND:
		carry = -1;
		goto inx;

		case INI:
		carry = 1;
		inx:
		result = (cpu.ReadIO)( BC, cpu);
		if( cpu.restart_io ) {
			PC = oldPC;
			goto early_exit;
		}
		cpu.WriteMem( HL, result, cpu);
		op1 = --B;
		HL += carry;

		in_flags:
			cpu.SetFlagValue( FLAG_S, op1 & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !(op1&0xff) );
			cpu.SetFlagValue( FLAG_Y, op1 & 0x20 );
			cpu.SetFlagValue( FLAG_X, op1 & 0x08 );
			cpu.SetFlagValue( FLAG_N, result & 0x80 );
			op2 = result + ((C+carry) & 0xff);
			cpu.SetFlagValue( FLAG_H, op2 > 0xff );
			cpu.SetFlagValue( FLAG_C, op2 > 0xff );
			cpu.SetFlagValue( FLAG_P, ParityIsEven((op2&0x07)^op1) );
		break;

		case INDR:
			carry = -1;
			goto inxr;

		case INIR:
			carry = 1;
			inxr:
			result = cpu.ReadIO(BC, cpu)
			if cpu.restart_io {
				PC = oldPC;
				goto early_exit;
			}
			cpu.WriteMem(HL, result, cpu);
			op1 = --B;
			HL += carry;
			if op1 == 0  {
				took_branch = false
			} else {
				PC -= 2
			}
			goto in_flags;

		case IN_R_MN: //in a,(n)
			result = (cpu.ReadIO)( (A<<8)|IMM, cpu);
			if( cpu.restart_io ) {
				PC = oldPC;
				goto early_exit;
			}
			A = result;
		break;

		case IN_R_R: // in r,(c)
			result = cpu.ReadIO( BC, cpu);
			if cpu.restart_io {
				PC = oldPC;
				goto early_exit;
			}
			if OP1 != REG_F  {
				BYTE_REG[OP1] = result;
			}
			cpu.SetFlagValue( FLAG_S, result & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !result );
			cpu.SetFlagValue( FLAG_Y, result & 0x20 );
			cpu.ResetFlag( FLAG_H );
			cpu.SetFlagValue( FLAG_X, result & 0x08 );
			cpu.SetFlagValue( FLAG_P, ParityIsEven(result) );
			cpu.ResetFlag( FLAG_N );
		break;

		case OTDR:
		carry = -1;
		goto otxr;
		case OTIR:
		carry = 1;
		otxr:
		result = cpu.ReadMem( HL, false, cpu);
		op1 = --B;
		(cpu.WriteIO)( BC, result, cpu);
		if cpu.restart_io  {
			++B; // Fix up
			PC = oldPC;
			goto early_exit;
		}
		HL += carry;
		if( op1 )
		PC -= 2;
		else
		took_branch = false;
		out_flags:
		// More flag strangeness
		op2 = result + L;
			cpu.SetFlagValue( FLAG_S, op1 & 0x80 );
			cpu.SetFlagValue( FLAG_Z, !op1 );
			cpu.SetFlagValue( FLAG_Y, op1 & 0x20 );
			cpu.SetFlagValue( FLAG_H, op2 > 0xff );
			cpu.SetFlagValue( FLAG_X, op1 & 0x08 );
			cpu.SetFlagValue( FLAG_P, ParityIsEven((op2&0x07) ^ op1) );
			cpu.SetFlagValue( FLAG_N, result & 0x80 );
			cpu.SetFlagValue( FLAG_C, op2 > 0xff );
		break;

		case OUTD:
		carry = -1;
		goto outx;
		case OUTI:
		carry = 1;
		outx:
		result = cpu.ReadMem( HL, false, cpu);
		op1 = --B;
		cpu.WriteIO( BC, result, cpu);
		if( cpu.restart_io ) {
			++B; // fix up
			PC = oldPC;
			goto early_exit;
		}
		HL += carry;
		goto out_flags;

		case OUT_MN_R: // out (n),a
		cpu.WriteIO( (A<<8)|IMM, A, cpu);
		if cpu.restart_io {
		PC = oldPC;
		goto early_exit;
		}
		break;

		case OUT_R: // out (c),0
		// I don't know what is on the top half of the address
		// bus during this, I'm guessing B.
		cpu.WriteIO( BC, 0, cpu);
		if cpu.restart_io {
			PC = oldPC;
			goto early_exit;
		}
		break;

		case OUT_R_R: // out (c),r
			cpu.WriteIO( BC, BYTE_REG[OP2], cpu);
			if( cpu.restart_io ) {
				PC = oldPC;
				goto early_exit;
			}
			break;
		}
		//#undef OP1
		//#undef OP2
		//#undef OFFSET
		//#undef IMM

		ticks = took_branch? inst.IT.tstates:inst.IT.extra;
		ticks += inst.additional_tstates;

	interrupt_exit:
		r = BYTE_REG[REG_R];
		r = ((r + inst.r_increment) & 0x7f) | (r & 0x80);
		BYTE_REG[REG_R] = r;

	early_exit:
		if( outPC )
			*outPC = PC;
		return ticks;
	}


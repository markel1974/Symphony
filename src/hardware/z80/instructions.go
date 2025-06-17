package z80

const (
	FLAG_C = iota //!< Carry flag.
	FLAG_N        //!< Add/subtract flag.
	FLAG_P        //!< Parity/overflow flag (P/V).
	FLAG_X        //!< Flag 3.
	FLAG_H        //!< Half-carry flag.
	FLAG_Y        //!< Flag 5.
	FLAG_Z        //!< Zero flag.
	FLAG_S        //!< Sign flag.
)

const (
	ADC_RR_RR  = iota //!< adc
	ADC_R_I           //!< adc
	ADC_R_MRR         //!< adc
	ADC_R_N           //!< adc
	ADC_R_R           //!< adc
	ADD_RR_RR         //!< add
	ADD_R_I           //!< add
	ADD_R_MRR         //!< add
	ADD_R_N           //!< add
	ADD_R_R           //!< add
	AND_I             //!< and
	AND_MRR           //!< and
	AND_N             //!< and
	AND_R             //!< and
	BIT_I             //!< bit
	BIT_MRR           //!< bit
	BIT_R             //!< bit
	CALL_C_MNN        //!< call
	CALL_MNN          //!< call
	CCF               //!< ccf
	CPD               //!< cpd
	CPDR              //!< cpdr
	CPI               //!< cpi
	CPIR              //!< cpir
	CPL               //!< cpl
	CP_I              //!< cp
	CP_MRR            //!< cp
	CP_N              //!< cp
	CP_R              //!< cp
	DAA               //!< daa
	DEC_I             //!< dec
	DEC_MRR           //!< dec
	DEC_R             //!< dec
	DEC_RR            //!< dec
	DI                //!< di
	DJNZ              //!< djnz
	EI                //!< ei
	EXX               //!< exx
	EX_MRR_RR         //!< ex
	EX_RR_RR          //!< ex
	HALT              //!< halt
	IM                //!< im
	INC_I             //!< inc
	INC_MRR           //!< inc
	INC_R             //!< inc
	INC_RR            //!< inc
	IND               //!< ind
	INDR              //!< indr
	INI               //!< ini
	INIR              //!< inir
	IN_R_MN           //!< in
	IN_R_R            //!< in
	JP_C_MNN          //!< jp
	JP_MNN            //!< jp
	JP_MRR            //!< jp
	JR                //!< jr
	JR_C              //!< jr
	LDD               //!< ldd
	LDDR              //!< lddr
	LDI               //!< ldi
	LDIR              //!< ldir
	LD_I_N            //!< ld
	LD_I_R            //!< ld
	LD_MNN_R          //!< ld
	LD_MNN_RR         //!< ld
	LD_MRR_N          //!< ld
	LD_MRR_R          //!< ld
	LD_RR_MNN         //!< ld
	LD_RR_NN          //!< ld
	LD_RR_RR          //!< ld
	LD_R_I            //!< ld
	LD_R_MNN          //!< ld
	LD_R_MRR          //!< ld
	LD_R_N            //!< ld
	LD_R_R            //!< ld
	NEG               //!< neg
	NOP               //!< nop
	OR_I              //!< or
	OR_MRR            //!< or
	OR_N              //!< or
	OR_R              //!< or
	OTDR              //!< otdr
	OTIR              //!< otir
	OUTD              //!< outd
	OUTI              //!< outi
	OUT_MN_R          //!< out
	OUT_R             //!< out
	OUT_R_R           //!< out
	POP_RR            //!< pop
	PUSH_RR           //!< push
	RES_I             //!< res
	RES_MRR           //!< res
	RES_R             //!< res
	RET               //!< ret
	RETI              //!< reti
	RETN              //!< retn
	RET_C             //!< ret
	RLA               //!< rla
	RLCA              //!< rlca
	RLC_I             //!< rlc
	RLC_MRR           //!< rlc
	RLC_R             //!< rlc
	RLD               //!< rld
	RL_I              //!< rl
	RL_MRR            //!< rl
	RL_R              //!< rl
	RRA               //!< rra
	RRCA              //!< rrca
	RRC_I             //!< rrc
	RRC_MRR           //!< rrc
	RRC_R             //!< rrc
	RRD               //!< rrd
	RR_I              //!< rr
	RR_MRR            //!< rr
	RR_R              //!< rr
	RST               //!< rst
	SBC_RR_RR         //!< sbc
	SBC_R_I           //!< sbc
	SBC_R_MRR         //!< sbc
	SBC_R_N           //!< sbc
	SBC_R_R           //!< sbc
	SCF               //!< scf
	SET_I             //!< set
	SET_MRR           //!< set
	SET_R             //!< set
	SLA_I             //!< sla
	SLA_MRR           //!< sla
	SLA_R             //!< sla
	SLL_I             //!< sll
	SLL_MRR           //!< sll
	SLL_R             //!< sll
	SRA_I             //!< sra
	SRA_MRR           //!< sra
	SRA_R             //!< sra
	SRL_I             //!< srl
	SRL_MRR           //!< srl
	SRL_R             //!< srl
	SUB_I             //!< sub
	SUB_MRR           //!< sub
	SUB_N             //!< sub
	SUB_R             //!< sub
	XOR_I             //!< xor
	XOR_MRR           //!< xor
	XOR_N             //!< xor
	XOR_R             //!< xor
)

type InstructionType int

// Exactly two instructions contain both an offset and an immediate:
// dd36 d n: ld (ix+d),n
// fd36 d n: ld (iy+d),n */
//
//	Describes the layout of the operands for an instruction. */
const (
	TYPE_NONE         InstructionType = iota //!< No operands.
	TYPE_IMM_N                               //!< 8 bit immediate.
	TYPE_IMM_NN                              //!< 16 bit immediate.
	TYPE_OFFSET                              //!< 8 bit signed offset.
	TYPE_DISP                                //!< 8 bit signed offset - 2.
	TYPE_OFFSET_IMM_N                        //!< 8 bit signed offset and 8 bit immediate.
)

// 4 (3) words native sized words on 32 (64) bit machine.
//
//	Uniform template for the instruction tables.
//
// headerfile z80_instructions.h zel/z80_instructions.h
type InstructionTemplate struct {
	kind          InstructionType //!< Type of the instruction.
	operand1      int16           //!< Type of the first operand, if any.
	operand2      int16           //!< Type of the second operand, if any.
	extra         int16           //!< Type of third operand or tstates when a branch is taken, if applicable.
	tstates       uint8           //!< Base number of clock ticks the instruction takes.
	operand_types uint8           //!< Operand layout in memory for the instruction.
	format        []uint8         //!< Format specifier string for disassembly.
}

// ! Completely describes an instruction.
//
//	headerfile z80_instructions.h zel/z80_instructions.h
type Instruction struct {
	IT                 *InstructionTemplate //!< Template for the instruction.
	immediate          uint                 //!< Immediate value, if any.
	additional_tstates uint                 //!< Additional clock ticks.
	r_increment        uint                 //!< Amount by which the \c r register is incremented.
	offset             int                  //!< Offset, if any.
}

/*

// Instruction table for unprefixed instructions.
extern const InstructionTemplate Unprefixed[256];
//  Instruction table for cb prefixed instructions.
extern const InstructionTemplate CB_Prefixed[256];
//  Instruction table for dd prefixed instructions.
extern const InstructionTemplate DD_Prefixed[256];
// Instruction table for ddcb prefixed instructions.
extern const InstructionTemplate DDCB_Prefixed[256];
//  Instruction table for ed prefixed instructions.
extern const InstructionTemplate ED_Prefixed[256];
//  Instruction table for fd prefixed instructions.
extern const InstructionTemplate FD_Prefixed[256];
//  Instruction table for fdcb prefixed instructions.
extern const InstructionTemplate FDCB_Prefixed[256];
*/

// Memory reading callback.
//\param addr The address to read.
//\param data Callback data from IF_ID().
//\return The value at address \a addr.
//
//typedef uint8_t (*ReadMemFunction)(uint16_t addr, void *data);

/*! Instruction fetch and instruction decode. Fetchs and decodes the
 * instruction pointed to by \a address into \a *inst.
 * \param inst Pointer to an \c Instruction. \a *inst is set to the
 * decoded instruction.
 * \param address The address of the instruction.
 * \param ReadMem Called repeatedly to get bytes of the instruction.
 * \a data is passed as the \c data argument.
 * \param data Arbitrary data passed to the \a ReadMem callback.
 * \return The length of the instruction.
 */
//int IF_ID( Instruction *inst, uint16_t address, ReadMemFunction ReadMem, void *data );

/*! Disassemble the instruction pointed to by \a inst into \a buffer.
 * \param inst Pointer to an \c Instruction.
 * \param buffer Buffer into which the disassembly is written. It must
 * be large enough to hold 25 bytes.
 */
//void DisassembleInstruction( const Instruction *inst, char *buffer );

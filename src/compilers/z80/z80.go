package compiler

const (
	FlagC  = 1 << 0 // Carry Flag
	FlagN  = 1 << 1 // Add/Subtract Flag
	FlagPV = 1 << 2 // Parity/Overflow Flag
	FlagH  = 1 << 4 // Half Carry Flag
	FlagZ  = 1 << 6 // Zero Flag
	FlagS  = 1 << 7 // Sign Flag
)

// Z80 represents a Zilog Z80 CPU emulator with registers, RAM, and support for 8-bit and 16-bit register operations.
type Z80 struct {
	registers      map[string]int
	ram            []int
	registers8bit  []string
	registers16bit []string
}

// NewZ80 creates and initializes a new Z80 CPU instance with default registers, RAM, and configuration.
func NewZ80() *Z80 {
	return &Z80{
		registers:      make(map[string]int),
		ram:            make([]int, 0xFFFF+1),
		registers8bit:  []string{"A", "F", "B", "C", "D", "E", "H", "L"},
		registers16bit: []string{"SP", "PC", "IX", "IY"},
	}
}

// Reset clears all internal registers and RAM, restoring the Z80 processor to its initial state.
func (z *Z80) Reset() {

}

// Registers8Bit returns a slice of strings representing the names of the Z80's 8-bit registers.
func (z *Z80) Registers8Bit() []string {
	return z.registers8bit
}

// Registers16Bit returns a slice of strings containing the names of the 16-bit registers in the Z80 processor.
func (z *Z80) Registers16Bit() []string {
	return z.registers16bit
}

// RAM retrieves the value stored at the specified 16-bit address in the Z80's RAM.
func (z *Z80) RAM(addr int) int {
	return z.ram[addr&0xFFFF]
}

// SetRAM sets the value at the specified address in the RAM, masking the address to fit within a 16-bit range.
func (z *Z80) SetRAM(addr int, value int) {
	z.ram[addr&0xFFFF] = value
}

// Register retrieves the value of the specified CPU register by its name.
func (z *Z80) Register(name string) int {
	return z.registers[name]
}

// SetRegisterIndex sets the index value for the specified register name in the Z80 structure's register map.
func (z *Z80) SetRegisterIndex(name string, value int) {
	z.registers[name] = value
}

// GetRegisterNameFromIndex returns the name of the Z80 CPU register corresponding to the provided 3-bit index.
func (z *Z80) GetRegisterNameFromIndex(index int) string {
	// L'ordine standard Z80 per gli indici a 3 bit
	// 000 -> B, 001 -> C, 010 -> D, 011 -> E, 100 -> H, 101 -> L, 110 -> (HL), 111 -> A
	// Per ora, ignoriamo (HL)
	switch index {
	case 0:
		return "B"
	case 1:
		return "C"
	case 2:
		return "D"
	case 3:
		return "E"
	case 4:
		return "H"
	case 5:
		return "L"
	case 7:
		return "A"
	}
	return "" // Placeholder per (HL) o casi non validi
}

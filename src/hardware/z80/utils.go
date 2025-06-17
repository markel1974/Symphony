package z80

func ReadInstructionMemory(address uint16, cpu *Z80) byte {
	return cpu.ReadMem(address, true, cpu)
}

func ParityIsEven(a uint) bool {
	b := a & 0x55
	c := (a >> 1) & 0x55
	a = b + c
	b = a & 0x33
	c = (a >> 2) & 0x33
	a = b + c
	b = a & 0x0f
	c = (a >> 4) & 0x0f
	val := (b + c) & 1
	if val != 0 {
		return true
	}
	return false
	//return !((b + c)&1);
}

package flash

// flashType defines the structure for flash memory properties, including identification, size, and operational parameters.
// It contains fields for manufacturer and device IDs, size details, sector configuration, magic addresses, and status bits.
// Additionally, it stores timing parameters for erase and operation cycles.
type flashType struct {
	manufacturerId           uint8
	deviceId                 uint8
	deviceIdAddr             uint8
	size                     uint
	sectorMask               uint
	sectorSize               uint
	sectorShift              uint
	magic1Addr               uint
	magic2Addr               uint
	magic1Mask               uint
	magic2Mask               uint
	statusToggleBits         uint
	eraseSectorTimeoutCycles uint
	eraseSectorCycles        uint
	eraseChipCycles          uint
}

// flashTypes defines an array of flashType configurations for different flash memory models, indexed by KindNum.
var flashTypes = [KindNum]flashType{
	/* AM29F040 */
	{0x01, 0xa4, 1,
		0x80000,
		0x70000, 0x10000, 16,
		0x5555, 0x2aaa, 0x7fff, 0x7fff,
		0x40,
		80, 2000000, 14000000}, /* may take up to 30s and 120s */
	/* AM29F040B */
	{0x01, 0xa4, 1,
		0x80000,
		0x70000, 0x10000, 16,
		0x555, 0x2aa, 0x7ff, 0x7ff,
		0x40,
		50, 1000000, 8000000}, /* may take up to 8s and 64s */
	/* 29F010 */
	{0x01, 0x20, 1,
		0x20000,
		0x1c000, 0x04000, 14,
		0x5555, 0x2aaa, 0x7fff, 0x7fff,
		0x40,
		80, 1000000, 1000000}, /* may take up to 15s */
	/* 29F032B with A0/1 swap */
	{0x01, 0x41, 1,
		0x400000,
		0x3f0000, 0x10000, 16,
		0x556, 0x2a9, 0x7ff, 0x7ff,
		0x44,
		50, 1000000, 64000000}, /* may take up to 8s */
	/* Expansion S29GL064N */
	{0x01, 0x7e, 2,
		0x800000,
		// FIXME: some models support non-uniform sector layout
		0x7f0000, 0x10000, 16,
		0xaaa, 0x555, 0xfff, 0xfff,
		0x40,
		50, 500000, 64000000}, /* may take up to 3.5s and 128s */
}

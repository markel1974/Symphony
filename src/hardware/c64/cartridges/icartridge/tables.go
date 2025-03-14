package icartridge

type RomInterval int

const (
	ROM_LO   = RomInterval(1)
	ROM_HI_1 = RomInterval(2)
	ROM_HI_2 = RomInterval(4)
)

type CartridgeMode int

const (
	CartridgeMode16K = CartridgeMode(iota)
	CartridgeMode8K
	CartridgeModeUltimax
	CartridgeModeOff
)

type CartridgeSpec struct {
	Game         uint8
	ExRom        uint8
	IntervalLow  RomInterval
	IntervalHigh RomInterval
}

var _cartridgesSpec []*CartridgeSpec

func init() {
	_cartridgesSpec = make([]*CartridgeSpec, CartridgeModeOff+1)
	_cartridgesSpec[CartridgeMode16K] = &CartridgeSpec{Game: 0, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_1}
	_cartridgesSpec[CartridgeMode8K] = &CartridgeSpec{Game: 0, ExRom: 1, IntervalLow: ROM_LO, IntervalHigh: 0}
	_cartridgesSpec[CartridgeModeUltimax] = &CartridgeSpec{Game: 1, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_2}
	_cartridgesSpec[CartridgeModeOff] = &CartridgeSpec{Game: 1, ExRom: 1, IntervalLow: 0, IntervalHigh: 0}
}

func GetCartridgeSpec(ct CartridgeMode) *CartridgeSpec {
	return _cartridgesSpec[ct]
}

//OFF
//GAME = 1, EXROM = 1

//8K Cartridge, $8000-$9FFF (ROML).
//GAME = 1, EXROM = 0
//ROML is read only. Basic ROM and Kernal ROM are available.

//16K Cartridge, $8000-$9FFF / $A000-$BFFF (ROML / ROMH).
//GAME = 0, EXROM = 0
//ROML/ROMH are read only, Basic ROM is overwritten by ROMH.

//16K Cartridge, $8000-$9FFF / $E000-$FFFF (ROML / ROMH). Ultimax mode.
//GAME = 0, EXROM = 1
//Ultimax mode is an emulation of the Japanese CBM machine called “MAX”. It is a predecessor of the C64 with less RAM. In Ultimax mode ROMH replaces the kernal at $E000. You do not need ROML for a cartridge to function and can be left out.

package via

const intrVIA1Id = 3
const intrVIA2Id = 4

type IMechanics interface {
	WriteProtectionState() uint8
	SyncFound() bool
	RotateDisk()
	MoveHeadOut()
	MoveHeadIn()
	ReadByte() uint8
	WriteByte(uint8)
	SetMotor(bool)
}

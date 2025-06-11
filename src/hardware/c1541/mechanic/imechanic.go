package mechanic

import "github.com/markel1974/c64emu/src/hardware/c1541/disk"

//see
//http://www.baltissen.org/newhtm/1541c.htm
//https://sta.c64.org/cbm1541mem.html
//https://c64os.com/post/howdoes1541work

//Dati tecnici originali:
//Tempo di seek traccia-traccia: 3 ms (documentazione tecnica del 1541)
//Modalità movimento: Full-step (1 traccia per step, non half-step)
//Velocità massima: ~333 tracce/secondo (1 / 0.003)

// syncByte represents a constant value used for synchronization detection in disk emulation processes.
// syncTolerance defines the tolerance level for synchronization adjustments, set as 3 or 5.
// headStep defines the step value used for calculating head movements in a specific context.
// headHalfStep represents the number of half-tracks the head can move, calculated as twice the value of headStep.
// motorSpinUpDelay specifies the delay in microseconds required for the motor to spin up, set as 300ms.
// stepDelay specifies the delay in microseconds per track movement, simulating a stepper motor's operation.
const (
	syncByte         = 0xff
	syncTolerance    = 5 // 3 or 5
	minHeadHalfStep  = 2
	maxHeadHalfStep  = 72
	motorSpinUpDelay = 300_000 // 300ms per avviamento
	headActivation   = 1500
	headInwardDelay  = headActivation + 1300 //2800
	headOutwardDelay = headActivation + 1700 //3200

	//headHalfStepDelay = 3000    // 3ms per traccia (simula il motore passo-passo)
)

type IMechanic interface {
	Reset()

	Setup() error

	InsertDisk(d disk.IDisk) error

	RemoveDisk() error

	SetWrite(w bool)

	EmulationRequired() bool

	Emulate()

	ReadByte() uint8

	WriteByte(data uint8)

	SyncFound() bool

	ByteReady() bool

	SetMotor(m bool)

	HasDisk() bool

	WriteProtectionState() uint8

	MoveHeadOut()

	MoveHeadIn()
}

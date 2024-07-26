package via

type IInterrupts interface {
	ClearVIA1IRQ()
	TriggerVIA1IRQ()
	ClearVIA2IRQ()
	TriggerVIA2IRQ()
	ClearNMI()
	TriggerNMI()
}

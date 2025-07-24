package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"reflect"
)

// Costanti per i nomi delle proprietà della CPU (per l'introspezione).
const (
	aId  = "a"  // Accumulatore
	xId  = "x"  // Registro X
	yId  = "y"  // Registro Y
	pcId = "pc" // Program Counter
	spId = "sp" // Stack Pointer
	srId = "sr" // Status Register (tutti i flag)
	// ... potresti aggiungere costanti per i singoli flag (nFlagId, vFlagId, ecc.) ...
)

// Reflect è una struct che incapsula una CPU e le sue proprietà,
// per l'accesso tramite introspezione.
type Reflect struct {
	props *component.Properties
	cpu   *CPU
}

// NewReflect crea un nuovo oggetto Reflect per una data CPU.
func NewReflect(c *CPU) *Reflect {
	r := &Reflect{
		props: nil,
		cpu:   c,
	}
	r.props = component.NewProperties()

	// Registra le proprietà della CPU.
	r.props.Add(component.NewPropertyInfo(aId, "Accumulator", false, r.getA, r.setA))
	r.props.Add(component.NewPropertyInfo(xId, "X Register", false, r.getX, r.setX))
	r.props.Add(component.NewPropertyInfo(yId, "Y Register", false, r.getY, r.setY))
	r.props.Add(component.NewPropertyInfo(pcId, "Program Counter", false, r.getPC, r.setPC))
	r.props.Add(component.NewPropertyInfo(spId, "Stack Pointer", false, r.getSP, r.setSP))
	//r.props.add(board.NewPropertyInfo(srId, "Status Register", false, r.getSR, r.setSR))
	// ... registra le altre proprietà (flag, ecc.) ...
	return r
}

// GetProperties restituisce la mappa delle proprietà della CPU.
func (r *Reflect) GetProperties() *component.Properties {
	return r.props
}

func (r *Reflect) getA() uint8 { // Restituisce *direttamente* uint8
	return r.cpu.a
}

func (r *Reflect) setA(v uint8) { // Accetta *direttamente* uint8
	r.cpu.a = v
}

func (r *Reflect) getX() uint8 {
	return r.cpu.x
}

func (r *Reflect) setX(v uint8) {
	r.cpu.x = v
}
func (r *Reflect) getY() uint8 {
	return r.cpu.y
}

func (r *Reflect) setY(v uint8) {
	r.cpu.y = v
}

func (r *Reflect) getPC() uint16 {
	return r.cpu.pc
}

func (r *Reflect) setPC(v uint16) {
	r.cpu.pc = v
}

func (r *Reflect) getSP() uint8 {
	return r.cpu.sp
}

func (r *Reflect) setSP(v uint8) {
	r.cpu.sp = v
}

func (r *Reflect) getNext() string {
	if r.cpu.next == nil {
		return ""
	}
	name, ok := r.cpu.executor.opContainer[reflect.ValueOf(r.cpu.next)]
	if !ok {
		return ""
	}
	return name
}

func (r *Reflect) setNext(name string) {
	f, ok := r.cpu.executor.opReverseContainer[name]
	if !ok {
		return
	}
	next, ok := f.Interface().(func(*CPU))
	if !ok {
		return
	}
	r.cpu.next = next
}

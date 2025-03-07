package mos6526

import (
	"fmt"
	"github.com/markel1974/c64emu/src/components/board"
)

// Costanti per i nomi delle proprietà del TOD (per l'introspezione).
const (
	tod10thsId       = "tod10ths"
	todSecId         = "todSec"
	todMinId         = "todMin"
	todHrId          = "todHr"
	todHaltId        = "todHalt"
	todDividerId     = "todDivider"
	todShadow10thsId = "todShadow10ths"
	todShadowSecId   = "todShadowSec"
	todShadowMinId   = "todShadowMin"
	alm10thsId       = "alm10ths"
	almSecId         = "almSec"
	almMinId         = "almMin"
	almHrId          = "almHr"
)

// Reflect è una struct che incapsula un TOD e le sue proprietà,
// per l'accesso tramite introspezione.
type Reflect struct {
	props *board.Properties
	tod   *TOD
}

// NewReflect crea un nuovo oggetto Reflect per un dato TOD.
func NewReflect(t *TOD) *Reflect {
	r := &Reflect{
		props: nil,
		tod:   t,
	}
	r.props = board.NewProperties(r.RunCommand)

	//TODO
	/*
		r.props.Add(board.CreatePropertyInfo(tod10thsId,  "TOD 10ths", false, r.getTod10ths, r.setTod10ths))
		r.props.Add(board.CreatePropertyInfo(todSecId,  "TOD Seconds", false, r.getTodSec, r.setTodSec))
		r.props.Add(board.CreatePropertyInfo(todMinId,  "TOD Minutes", false, r.getTodMin, r.setTodMin))
		r.props.Add(board.CreatePropertyInfo(todHrId,  "TOD Hours", false, r.getTodHr, r.setTodHr))
		r.props.Add(board.CreatePropertyInfo(todHaltId,  "TOD Halted", false, r.getTodHalt, r.setTodHalt))
		r.props.Add(board.CreatePropertyInfo(todDividerId,  "TOD Divider", false, r.getTodDivider, r.setTodDivider))
		r.props.Add(board.CreatePropertyInfo(todShadow10thsId,  "TOD Shadow 10ths", true, r.getTodShadow10ths, nil)) // Sola lettura
		r.props.Add(board.CreatePropertyInfo(todShadowSecId,  "TOD Shadow Seconds", true, r.getTodShadowSec, nil))   // Sola lettura
		r.props.Add(board.CreatePropertyInfo(todShadowMinId,  "TOD Shadow Minutes", true, r.getTodShadowMin, nil))   // Sola lettura
		r.props.Add(board.CreatePropertyInfo(alm10thsId,  "Alarm 10ths", false, r.getAlm10ths, r.setAlm10ths))
		r.props.Add(board.CreatePropertyInfo(almSecId, "Alarm Seconds", false, r.getAlmSec, r.setAlmSec))
		r.props.Add(board.CreatePropertyInfo(almMinId,  "Alarm Minutes", false, r.getAlmMin, r.setAlmMin))
		r.props.Add(board.CreatePropertyInfo(almHrId,  "Alarm Hours", false, r.getAlmHr, r.setAlmHr))

	*/
	return r
}

func (r *Reflect) GetProperties() *board.Properties {
	return r.props
}

func (r *Reflect) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (r *Reflect) getTod10ths() (interface{}, error) {
	return r.tod.tod10ths, nil
}
func (r *Reflect) setTod10ths(v uint8) error {
	//TODO
	//r.tod.tod10ths = v
	//r.tod.Set10ths(false, value) // Imposta l'ora corrente, non l'allarme.
	return nil
}

// ... Implementa getter/setter simili per tutti gli altri campi di TOD ...
//(todSec, todMin, todHr, todHalt, todDivider, todShadow10ths, todShadowSec, todShadowMin, alm10ths, almSec, almMin, almHr)

// Esempio per Get/Set di un campo bool:
func (r *Reflect) getTodHalt() bool {
	return r.tod.todHalt
}
func (r *Reflect) setTodHalt(v bool) {
	r.tod.todHalt = v
}

func (r *Reflect) getTodShadow10ths() int {
	return r.tod.todShadow10ths
}

func (r *Reflect) getTodDivider() int {
	return r.tod.todDivider
}

func (r *Reflect) setTodDivider(v int) {
	r.tod.todDivider = v
}

// Funzioni mancanti

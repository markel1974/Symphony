package mos6526

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
	tod *TOD
}

// NewReflect crea un nuovo oggetto Reflect per un dato TOD.
func NewReflect(t *TOD) *Reflect {
	r := &Reflect{
		tod: t,
	}

	//TODO
	/*
		r.AddProperty(tod10thsId,  "TOD 10ths", false, r.getTod10ths, r.setTod10ths)
		r.AddProperty(todSecId,  "TOD Seconds", false, r.getTodSec, r.setTodSec)
		r.AddProperty(todMinId,  "TOD Minutes", false, r.getTodMin, r.setTodMin)
		r.pAddProperty(todHrId,  "TOD Hours", false, r.getTodHr, r.setTodHr)
		r.AddProperty(todHaltId,  "TOD Halted", false, r.getTodHalt, r.setTodHalt)
		r.AddProperty(todDividerId,  "TOD Divider", false, r.getTodDivider, r.setTodDivider)
		r.AddProperty(todShadow10thsId,  "TOD Shadow 10ths", true, r.getTodShadow10ths, nil)
		r.AddProperty(todShadowSecId,  "TOD Shadow Seconds", true, r.getTodShadowSec, nil)
		r.AddProperty(todShadowMinId,  "TOD Shadow Minutes", true, r.getTodShadowMin, nil)
		r.AddProperty(alm10thsId,  "Alarm 10ths", false, r.getAlm10ths, r.setAlm10ths)
		r.AddProperty(almSecId, "Alarm Seconds", false, r.getAlmSec, r.setAlmSec)
		r.AddProperty(almMinId,  "Alarm Minutes", false, r.getAlmMin, r.setAlmMin)
		r.AddProperty(almHrId,  "Alarm Hours", false, r.getAlmHr, r.setAlmHr)
	*/
	return r
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

//func (r *Reflect) getTodDivider() int {
//	return r.tod.todDivider
//}

//func (r *Reflect) setTodDivider(v int) {
//	r.tod.todDivider = v
//}

// Funzioni mancanti

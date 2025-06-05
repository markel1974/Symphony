package mos6526

const (
	crId            = "cr"
	crNewId         = "crNew"
	crNewPendingId  = "crNewPending"
	timerId         = "timer"
	timerLatchId    = "timerLatch"
	toggleModeId    = "toggleMode"
	timerLatchLowId = "timerLatchLow" //Potrebbe non servire
	cntId           = "cnt"
	timerStateId    = "timerState"
	countModeId     = "countMode"
)

type TimerReflect struct {
	timer *Timer
}

func NewTimerReflect(t *Timer) *TimerReflect {
	r := &TimerReflect{timer: t}
	t.AddProperty(crId, "Control Register (CR)", false, r.getCr, r.setCr)
	t.AddProperty(crNewId, "New Control Register (CRNew)", false, r.getCrNew, r.setCrNew)
	t.AddProperty(crNewPendingId, "CRNew Pending Flag", false, r.getCrNewPending, r.setCrNewPending)
	t.AddProperty(timerId, "Timer Value", false, r.getTimer, r.setTimer)
	t.AddProperty(timerLatchId, "Timer Latch", false, r.getTimerLatch, r.setTimerLatch)
	t.AddProperty(toggleModeId, "Toggle Mode", false, r.getToggleMode, r.setToggleMode)
	t.AddProperty(timerLatchLowId, "Timer Latch (Low Byte)", false, r.getTimerLatchLow, r.setTimerLatchLow)
	t.AddProperty(cntId, "CNT Flag", false, r.getCNT, r.setCNT)
	t.AddProperty(timerStateId, "Timer State", false, r.getTimerState, r.setTimerState)
	t.AddProperty(countModeId, "Count Mode", false, r.getCountMode, r.setCountMode)
	return r
}

func (r *TimerReflect) getCr() uint8 {
	return r.timer.cr
}
func (r *TimerReflect) setCr(v uint8) error {
	r.timer.cr = v
	return nil
}

func (r *TimerReflect) getCrNew() uint8 {
	return r.timer.crNew
}
func (r *TimerReflect) setCrNew(v uint8) error {
	r.timer.crNew = v
	return nil
}

func (r *TimerReflect) getCrNewPending() bool {
	return r.timer.crNewPending
}
func (r *TimerReflect) setCrNewPending(v bool) error {
	r.timer.crNewPending = v
	return nil
}

func (r *TimerReflect) getTimer() uint16 {
	return r.timer.timer
}
func (r *TimerReflect) setTimer(v uint16) error {
	r.timer.timer = v
	return nil
}

func (r *TimerReflect) getTimerLatch() uint16 {
	return r.timer.timerLatch
}
func (r *TimerReflect) setTimerLatch(v uint16) error {
	r.timer.timerLatch = v
	return nil
}

func (r *TimerReflect) getToggleMode() bool {
	return r.timer.toggleMode
}
func (r *TimerReflect) setToggleMode(v bool) error {
	r.timer.toggleMode = v
	return nil
}

func (r *TimerReflect) getTimerLatchLow() uint16 {
	return r.timer.timerLatchLow
}

func (r *TimerReflect) setTimerLatchLow(v uint16) error {
	r.timer.timerLatchLow = v
	return nil
}

func (r *TimerReflect) getCNT() bool {
	return r.timer.cntLevel
}
func (r *TimerReflect) setCNT(v bool) error {
	r.timer.cntLevel = v
	return nil
}

func (r *TimerReflect) getTimerState() uint8 {
	return uint8(r.timer.timerState)
}
func (r *TimerReflect) setTimerState(v uint8) error {
	r.timer.timerState = TimerState(v)
	return nil
}

func (r *TimerReflect) getCountMode() uint8 {
	return r.timer.countMode
}
func (r *TimerReflect) setCountMode(v uint8) error {
	r.timer.updateCountMode(v)
	return nil
}

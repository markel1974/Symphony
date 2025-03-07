package mos6526

import (
	"fmt"
	"github.com/markel1974/c64emu/src/components/board"
)

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
	props *board.Properties
	timer *Timer
}

func NewTimerReflect(t *Timer) *TimerReflect {
	r := &TimerReflect{
		props: nil,
		timer: t,
	}
	r.props = board.NewProperties(r.RunCommand)

	r.props.Add(board.CreatePropertyInfo(crId, "Control Register (CR)", false, r.getCr, r.setCr))
	r.props.Add(board.CreatePropertyInfo(crNewId, "New Control Register (CRNew)", false, r.getCrNew, r.setCrNew))
	r.props.Add(board.CreatePropertyInfo(crNewPendingId, "CRNew Pending Flag", false, r.getCrNewPending, r.setCrNewPending))
	r.props.Add(board.CreatePropertyInfo(timerId, "Timer Value", false, r.getTimer, r.setTimer))
	r.props.Add(board.CreatePropertyInfo(timerLatchId, "Timer Latch", false, r.getTimerLatch, r.setTimerLatch))
	r.props.Add(board.CreatePropertyInfo(toggleModeId, "Toggle Mode", false, r.getToggleMode, r.setToggleMode))
	r.props.Add(board.CreatePropertyInfo(timerLatchLowId, "Timer Latch (Low Byte)", false, r.getTimerLatchLow, r.setTimerLatchLow)) //Potrebbe non essere necessario
	r.props.Add(board.CreatePropertyInfo(cntId, "CNT Flag", false, r.getCNT, r.setCNT))
	r.props.Add(board.CreatePropertyInfo(timerStateId, "Timer State", false, r.getTimerState, r.setTimerState))
	r.props.Add(board.CreatePropertyInfo(countModeId, "Count Mode", false, r.getCountMode, r.setCountMode))

	return r
}

func (r *TimerReflect) GetProperties() *board.Properties {
	return r.props
}

func (r *TimerReflect) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (r *TimerReflect) getCr() uint8 {
	return r.timer.cr
}
func (r *TimerReflect) setCr(v uint8) {
	r.timer.cr = v
}

func (r *TimerReflect) getCrNew() uint8 {
	return r.timer.crNew
}
func (r *TimerReflect) setCrNew(v uint8) {
	r.timer.crNew = v
}

func (r *TimerReflect) getCrNewPending() bool {
	return r.timer.crNewPending
}
func (r *TimerReflect) setCrNewPending(v bool) {
	r.timer.crNewPending = v
}

func (r *TimerReflect) getTimer() uint16 {
	return r.timer.timer
}
func (r *TimerReflect) setTimer(v uint16) {
	r.timer.timer = v
}

func (r *TimerReflect) getTimerLatch() uint16 {
	return r.timer.timerLatch
}
func (r *TimerReflect) setTimerLatch(v uint16) {
	r.timer.timerLatch = v
}

func (r *TimerReflect) getToggleMode() bool {
	return r.timer.toggleMode
}
func (r *TimerReflect) setToggleMode(v bool) {
	r.timer.toggleMode = v
}

func (r *TimerReflect) getTimerLatchLow() uint16 {
	return r.timer.timerLatchLow
}

func (r *TimerReflect) setTimerLatchLow(v uint16) {
	r.timer.timerLatchLow = v
}

func (r *TimerReflect) getCNT() bool {
	return r.timer.cnt
}
func (r *TimerReflect) setCNT(v bool) {
	r.timer.cnt = v
}

func (r *TimerReflect) getTimerState() uint8 {
	return uint8(r.timer.timerState)
}
func (r *TimerReflect) setTimerState(v uint8) {
	r.timer.timerState = TimerState(v)
}

func (r *TimerReflect) getCountMode() uint8 {
	return r.timer.countMode
}
func (r *TimerReflect) setCountMode(v uint8) {
	r.timer.updateCountMode(v)
}

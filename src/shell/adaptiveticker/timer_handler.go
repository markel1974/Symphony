/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package adaptiveticker

// TimerMode represents the operational mode of a timer, which can be continuous or counter-based.
type TimerMode int

// TimerModeContinuous represents a timer mode that operates continuously without stopping.
// TimerModeCounter represents a timer mode that operates based on a set counter or limit.
const (
	TimerModeContinuous TimerMode = iota
	TimerModeCounter    TimerMode = iota
)

// TimerHandler represents a timer event with scheduling information and execution control.
type TimerHandler struct {
	id        int
	target    chan *TimerHandler
	Event     interface{}
	first     int64
	interval  int64
	mode      TimerMode
	deadline  int64
	min       int64
	counter   int64
	loopCount int64
	removed   bool
}

// rounding rounds the given value up or down based on the provided rounding factor.
// If the value is less than the factor, it returns the factor.
// Otherwise, it calculates the largest multiple of the factor less than or equal to the value.
func rounding(val int64, r int64) int64 {
	if val < r {
		return r
	}
	return (val / r) * r
}

// NewTimerHandler initializes and returns a new TimerHandler with specified parameters for event scheduling and execution.
func NewTimerHandler(target chan *TimerHandler, event interface{}, first int64, interval int64, loopCount int64, min int64) *TimerHandler {
	var mode TimerMode

	if loopCount <= 0 {
		mode = TimerModeContinuous
	} else {
		mode = TimerModeCounter
	}

	t := &TimerHandler{
		target:    target,
		Event:     event,
		mode:      mode,
		min:       min,
		loopCount: loopCount,
		first:     rounding(first, min),
		interval:  rounding(interval, min),
		deadline:  0,
		counter:   0,
		removed:   false,
	}

	return t
}

// Prepare calculates and sets the next deadline for the TimerHandler based on the current time and its interval settings.
func (t *TimerHandler) Prepare(now int64) {
	interval := t.interval
	if t.counter == 0 {
		interval = t.first
	}
	t.deadline = rounding(now+interval, t.min)
	t.counter++
}

// SetId sets the identifier for the TimerHandler instance.
func (t *TimerHandler) SetId(id int) {
	t.id = id
}

// Unset marks the TimerHandler instance as removed by setting the `removed` flag to true.
func (t *TimerHandler) Unset() {
	t.removed = true
}

// IsUsable determines if the TimerHandler instance is usable based on its mode and current counter state.
func (t *TimerHandler) IsUsable() bool {
	var ret bool
	switch t.mode {
	case TimerModeCounter:
		if t.counter == t.loopCount {
			ret = false
		}
	case TimerModeContinuous:
		ret = true
	default:
		ret = false
	}
	return ret
}

// PostEvent sends the TimerHandler instance to its target channel asynchronously using a goroutine.
func (t *TimerHandler) PostEvent() {
	go func() {
		t.target <- t
	}()
}

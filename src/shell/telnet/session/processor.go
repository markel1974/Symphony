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

package session

import (
	"fmt"
)

type processorState int

const (
	stateBase   processorState = iota
	stateInIAC  processorState = iota
	stateInSB   processorState = iota
	stateCapSB  processorState = iota
	stateEscIAC processorState = iota
)

type processor struct {
	state     processorState
	currentSB IOCode

	capturedBytes []byte
	subData       map[IOCode][]byte
	cleanData     []byte
	listenFunc    func(IOCode, []byte)

	debug bool
}

func newProcessor() *processor {
	tp := &processor{
		state:     stateBase,
		debug:     false,
		currentSB: NUL,
	}
	return tp
}

func (tp *processor) Read(p []byte) (int, error) {
	maxLen := len(p)
	n := 0

	if maxLen >= len(tp.cleanData) {
		n = len(tp.cleanData)
	} else {
		n = maxLen
	}

	for i := 0; i < n; i++ {
		p[i] = tp.cleanData[i]
	}

	tp.cleanData = tp.cleanData[n:]

	return n, nil
}

func (tp *processor) capture(b byte) {
	if tp.debug {
		fmt.Println("Captured:", ByteToCodeString(b))
	}

	tp.capturedBytes = append(tp.capturedBytes, b)
}

func (tp *processor) dontCapture(b byte) {
	tp.cleanData = append(tp.cleanData, b)
}

func (tp *processor) resetSubDataField(code IOCode) {
	if tp.subData == nil {
		tp.subData = map[IOCode][]byte{}
	}

	tp.subData[code] = []byte{}
}

func (tp *processor) captureSubData(code IOCode, b byte) {
	if tp.debug {
		fmt.Println("Captured sub data:", CodeToString(code), b, string(b))
	}

	if tp.subData == nil {
		tp.subData = map[IOCode][]byte{}
	}

	tp.subData[code] = append(tp.subData[code], b)
}

func (tp *processor) addBytes(bytes []byte) {
	for _, b := range bytes {
		tp.addByte(b)
	}
}

func (tp *processor) addByte(b byte) {
	code := byteToCode[b]

	switch tp.state {
	case stateBase:
		if code == IAC {
			tp.state = stateInIAC
			tp.capture(b)
		} else {
			tp.dontCapture(b)
		}

	case stateInIAC:
		if code == WILL || code == WONT || code == DO || code == DONT {
			// Stay in this state
		} else if code == SB {
			tp.state = stateInSB
		} else {
			tp.state = stateBase
		}
		tp.capture(b)

	case stateInSB:
		tp.capture(b)
		tp.currentSB = code
		tp.state = stateCapSB
		tp.resetSubDataField(code)

	case stateCapSB:
		if code == IAC {
			tp.state = stateEscIAC
		} else {
			tp.captureSubData(tp.currentSB, b)
		}

	case stateEscIAC:
		if code == IAC {
			tp.state = stateCapSB
			tp.captureSubData(tp.currentSB, b)
		} else {
			tp.subDataFinished(tp.currentSB)
			tp.currentSB = NUL
			tp.state = stateBase
			tp.addByte(codeToByte[IAC])
			tp.addByte(b)
		}
	}
}

func (tp *processor) subDataFinished(code IOCode) {
	if tp.listenFunc != nil {
		tp.listenFunc(code, tp.subData[code])
	}
}

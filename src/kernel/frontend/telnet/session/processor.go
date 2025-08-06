package session

import (
	"fmt"
)

// processorState represents the state of a processor in a finite-state machine context.
type processorState int

// stateBase represents the default processor state.
// stateInIAC represents the processor state where an Interpret As Command (IAC) sequence is started.
// stateInSB represents the processor state where a Subnegotiation (SB) sequence is started.
// stateCapSB represents the processor state for capturing a Subnegotiation (SB) payload.
// stateEscIAC represents the processor state for handling an escaped Interpret As Command (IAC) sequence.
const (
	stateBase   processorState = iota
	stateInIAC  processorState = iota
	stateInSB   processorState = iota
	stateCapSB  processorState = iota
	stateEscIAC processorState = iota
)

// processor is a state-driven structure handling input/output operations and data processing at the byte level.
type processor struct {
	state     processorState
	currentSB IOCode

	capturedBytes []byte
	subData       map[IOCode][]byte
	cleanData     []byte
	listenFunc    func(IOCode, []byte)

	debug bool
}

// newProcessor initializes a new instance of the processor with default state and configurations.
func newProcessor() *processor {
	tp := &processor{
		state:     stateBase,
		debug:     false,
		currentSB: NUL,
	}
	return tp
}

// Read reads up to len(p) bytes from the processor's cleanData buffer into p and removes the read bytes from cleanData.
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

// capture adds the given byte to the capturedBytes slice and optionally logs it if debugging mode is enabled.
func (tp *processor) capture(b byte) {
	if tp.debug {
		fmt.Println("Captured:", ByteToCodeString(b))
	}

	tp.capturedBytes = append(tp.capturedBytes, b)
}

// dontCapture appends the given byte to the cleanData field without marking it as captured.
func (tp *processor) dontCapture(b byte) {
	tp.cleanData = append(tp.cleanData, b)
}

// resetSubDataField clears the subData map entry for the specified IOCode or initializes the map if it is nil.
func (tp *processor) resetSubDataField(code IOCode) {
	if tp.subData == nil {
		tp.subData = map[IOCode][]byte{}
	}

	tp.subData[code] = []byte{}
}

// captureSubData appends a byte to the subData map under the specified IOCode key, initializing the map if nil.
func (tp *processor) captureSubData(code IOCode, b byte) {
	if tp.debug {
		fmt.Println("Captured sub data:", CodeToString(code), b, string(b))
	}

	if tp.subData == nil {
		tp.subData = map[IOCode][]byte{}
	}

	tp.subData[code] = append(tp.subData[code], b)
}

// addBytes processes a slice of bytes by sequentially adding each byte to the processor using the addByte method.
func (tp *processor) addBytes(bytes []byte) {
	for _, b := range bytes {
		tp.addByte(b)
	}
}

// addByte processes a single byte of input based on the current state of the processor and updates its state accordingly.
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

// subDataFinished triggers the `listenFunc` callback with the specified IOCode and associated subData if `listenFunc` is defined.
func (tp *processor) subDataFinished(code IOCode) {
	if tp.listenFunc != nil {
		tp.listenFunc(code, tp.subData[code])
	}
}

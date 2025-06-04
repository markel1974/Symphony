package wasm_render

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// Inputs encapsulates input configurations, mappings, and state for joystick, keyboard, and clipboard interactions.
type Inputs struct {
	board        references.IBoard
	cfg          *config.Config
	keyMapper    map[string]func(bool)
	joyKeys      bool
	lastX        uint8
	lastY        uint8
	hasClipboard bool
}

// NewInputs initializes and returns a pointer to a new Inputs instance with default values.
func NewInputs() *Inputs {
	return &Inputs{
		board:        nil,
		cfg:          nil,
		keyMapper:    make(map[string]func(bool)),
		joyKeys:      true,
		lastX:        0,
		lastY:        0,
		hasClipboard: false,
	}
}

// Setup initializes the input mappings for a given board and configuration.
// It assigns key mappings to appropriate actions and handles setup logic for keyboard and joystick interactions.
func (g *Inputs) Setup(b references.IBoard, cfg *config.Config) error {
	g.board = b
	g.cfg = cfg
	g.keyMapper["Escape"] = func(p bool) { b.KeyboardSetKey(p, component.VKEscape) }
	g.keyMapper["Enter"] = func(p bool) { b.KeyboardSetKey(p, component.VKReturn) }
	g.keyMapper["Del"] = func(p bool) { b.KeyboardSetKey(p, component.VKDelete) }
	g.keyMapper["F1"] = func(p bool) { b.KeyboardSetKey(p, component.VKF1) }
	g.keyMapper["F2"] = func(p bool) { b.KeyboardSetKey(p, component.VKF2) }
	g.keyMapper["F3"] = func(p bool) { b.KeyboardSetKey(p, component.VKF3) }
	g.keyMapper["F4"] = func(p bool) { b.KeyboardSetKey(p, component.VKF4) }
	g.keyMapper["F5"] = func(p bool) { b.KeyboardSetKey(p, component.VKF5) }
	g.keyMapper["F6"] = func(p bool) { b.KeyboardSetKey(p, component.VKF6) }
	g.keyMapper["F7"] = func(p bool) { b.KeyboardSetKey(p, component.VKF7) }
	g.keyMapper["F8"] = func(p bool) { b.KeyboardSetKey(p, component.VKF8) }
	//g.keyMapper["backspace"] = func(p bool) { b.KeyboardSetKey(p, component.VKBack) }
	g.keyMapper[" "] = func(p bool) { b.KeyboardSetKey(p, component.VKSpace) }
	g.keyMapper["."] = func(p bool) { b.KeyboardSetKey(p, component.VKComma) }
	//g.keyMapper["period"] = func(p bool) { b.KeyboardSetKey(p, component.VKPeriod) }
	g.keyMapper[":"] = func(p bool) { b.KeyboardSetKey(p, component.VKSemicolon) }
	g.keyMapper[";"] = func(p bool) { b.KeyboardSetKey(p, component.VKQuote) }
	g.keyMapper["*"] = func(p bool) { b.KeyboardSetKey(p, component.VKAsterisk) }
	g.keyMapper["*"] = func(p bool) { b.KeyboardSetKey(p, component.VKAsterisk) }
	g.keyMapper["-"] = func(p bool) { b.KeyboardSetKey(p, component.VKMinus) }
	g.keyMapper["="] = func(p bool) { b.KeyboardSetKey(p, component.VKEqual) }
	g.keyMapper["+"] = func(p bool) { b.KeyboardSetKey(p, component.VKPlus) }

	g.keyMapper["0"] = func(p bool) { b.KeyboardSetKey(p, '0') }
	g.keyMapper["1"] = func(p bool) { b.KeyboardSetKey(p, '1') }
	g.keyMapper["2"] = func(p bool) { b.KeyboardSetKey(p, '2') }
	g.keyMapper["3"] = func(p bool) { b.KeyboardSetKey(p, '3') }
	g.keyMapper["4"] = func(p bool) { b.KeyboardSetKey(p, '4') }
	g.keyMapper["5"] = func(p bool) { b.KeyboardSetKey(p, '5') }
	g.keyMapper["6"] = func(p bool) { b.KeyboardSetKey(p, '6') }
	g.keyMapper["7"] = func(p bool) { b.KeyboardSetKey(p, '7') }
	g.keyMapper["8"] = func(p bool) { b.KeyboardSetKey(p, '8') }
	g.keyMapper["9"] = func(p bool) { b.KeyboardSetKey(p, '9') }

	g.keyMapper["A"] = func(p bool) { b.KeyboardSetKey(p, 'A') }
	g.keyMapper["B"] = func(p bool) { b.KeyboardSetKey(p, 'B') }
	g.keyMapper["C"] = func(p bool) { b.KeyboardSetKey(p, 'C') }
	g.keyMapper["D"] = func(p bool) { b.KeyboardSetKey(p, 'D') }
	g.keyMapper["E"] = func(p bool) { b.KeyboardSetKey(p, 'E') }
	g.keyMapper["F"] = func(p bool) { b.KeyboardSetKey(p, 'F') }
	g.keyMapper["G"] = func(p bool) { b.KeyboardSetKey(p, 'G') }
	g.keyMapper["H"] = func(p bool) { b.KeyboardSetKey(p, 'H') }
	g.keyMapper["I"] = func(p bool) { b.KeyboardSetKey(p, 'I') }
	g.keyMapper["J"] = func(p bool) { b.KeyboardSetKey(p, 'J') }
	g.keyMapper["K"] = func(p bool) { b.KeyboardSetKey(p, 'K') }
	g.keyMapper["L"] = func(p bool) { b.KeyboardSetKey(p, 'L') }
	g.keyMapper["M"] = func(p bool) { b.KeyboardSetKey(p, 'M') }
	g.keyMapper["N"] = func(p bool) { b.KeyboardSetKey(p, 'N') }
	g.keyMapper["O"] = func(p bool) { b.KeyboardSetKey(p, 'O') }
	g.keyMapper["P"] = func(p bool) { b.KeyboardSetKey(p, 'P') }
	g.keyMapper["Q"] = func(p bool) { b.KeyboardSetKey(p, 'Q') }
	g.keyMapper["R"] = func(p bool) { b.KeyboardSetKey(p, 'R') }
	g.keyMapper["S"] = func(p bool) { b.KeyboardSetKey(p, 'S') }
	g.keyMapper["T"] = func(p bool) { b.KeyboardSetKey(p, 'T') }
	g.keyMapper["U"] = func(p bool) { b.KeyboardSetKey(p, 'U') }
	g.keyMapper["V"] = func(p bool) { b.KeyboardSetKey(p, 'V') }
	g.keyMapper["W"] = func(p bool) { b.KeyboardSetKey(p, 'W') }
	g.keyMapper["X"] = func(p bool) { b.KeyboardSetKey(p, 'X') }
	g.keyMapper["Y"] = func(p bool) { b.KeyboardSetKey(p, 'Y') }
	g.keyMapper["Z"] = func(p bool) { b.KeyboardSetKey(p, 'Z') }

	g.keyMapper["a"] = func(p bool) { b.KeyboardSetKey(p, 'A') }
	g.keyMapper["b"] = func(p bool) { b.KeyboardSetKey(p, 'B') }
	g.keyMapper["c"] = func(p bool) { b.KeyboardSetKey(p, 'C') }
	g.keyMapper["d"] = func(p bool) { b.KeyboardSetKey(p, 'D') }
	g.keyMapper["e"] = func(p bool) { b.KeyboardSetKey(p, 'E') }
	g.keyMapper["f"] = func(p bool) { b.KeyboardSetKey(p, 'F') }
	g.keyMapper["g"] = func(p bool) { b.KeyboardSetKey(p, 'G') }
	g.keyMapper["h"] = func(p bool) { b.KeyboardSetKey(p, 'H') }
	g.keyMapper["i"] = func(p bool) { b.KeyboardSetKey(p, 'I') }
	g.keyMapper["j"] = func(p bool) { b.KeyboardSetKey(p, 'J') }
	g.keyMapper["k"] = func(p bool) { b.KeyboardSetKey(p, 'K') }
	g.keyMapper["l"] = func(p bool) { b.KeyboardSetKey(p, 'L') }
	g.keyMapper["m"] = func(p bool) { b.KeyboardSetKey(p, 'M') }
	g.keyMapper["n"] = func(p bool) { b.KeyboardSetKey(p, 'N') }
	g.keyMapper["o"] = func(p bool) { b.KeyboardSetKey(p, 'O') }
	g.keyMapper["p"] = func(p bool) { b.KeyboardSetKey(p, 'P') }
	g.keyMapper["q"] = func(p bool) { b.KeyboardSetKey(p, 'Q') }
	g.keyMapper["r"] = func(p bool) { b.KeyboardSetKey(p, 'R') }
	g.keyMapper["s"] = func(p bool) { b.KeyboardSetKey(p, 'S') }
	g.keyMapper["t"] = func(p bool) { b.KeyboardSetKey(p, 'T') }
	g.keyMapper["u"] = func(p bool) { b.KeyboardSetKey(p, 'U') }
	g.keyMapper["v"] = func(p bool) { b.KeyboardSetKey(p, 'V') }
	g.keyMapper["w"] = func(p bool) { b.KeyboardSetKey(p, 'W') }
	g.keyMapper["x"] = func(p bool) { b.KeyboardSetKey(p, 'X') }
	g.keyMapper["y"] = func(p bool) { b.KeyboardSetKey(p, 'Y') }
	g.keyMapper["z"] = func(p bool) { b.KeyboardSetKey(p, 'Z') }

	g.keyMapper["ArrowUp"] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJUp)
		} else {
			b.KeyboardSetKey(p, component.VKUp)
		}
	}
	g.keyMapper["ArrowDown"] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJDown)
		} else {
			b.KeyboardSetKey(p, component.VKDown)
		}
	}
	g.keyMapper["ArrowLeft"] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJLeft)
		} else {
			b.KeyboardSetKey(p, component.VKLeft)
		}
	}
	g.keyMapper["ArrowRight"] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJRight)
		} else {
			b.KeyboardSetKey(p, component.VKRight)
		}
	}
	g.keyMapper["Control"] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJFire)
		} else {
			b.KeyboardSetKey(p, component.VKControl)
		}
	}

	//err := clipboard.Init()
	//if err != nil {
	//	log.Printf("can't init clipboard: %s", err)
	//} else {
	//	g.hasClipboard = true
	//}
	return nil
}

// Key updates the state of a specified virtual key and triggers its associated function if mapped.
func (g *Inputs) Key(p string, pressed bool) {
	if f, ok := g.keyMapper[p]; ok {
		f(pressed)
	} else {
		fmt.Println("Can't find key", p, pressed)
	}
}

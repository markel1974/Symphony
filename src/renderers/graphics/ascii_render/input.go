package ascii_render

import (
	"log"
	"os"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
	"golang.design/x/clipboard"
)

type Inputs struct {
	board        references.IC64Board
	cfg          *config.Config
	keyMapper    map[byte]func(bool)
	joyKeys      bool
	lastX        uint8
	lastY        uint8
	hasClipboard bool
}

func NewInputs() *Inputs {
	return &Inputs{
		board:        nil,
		cfg:          nil,
		keyMapper:    make(map[byte]func(bool)),
		joyKeys:      true,
		lastX:        0,
		lastY:        0,
		hasClipboard: false,
	}
}

func (g *Inputs) Setup(b references.IC64Board, cfg *config.Config) error {
	g.board = b
	g.cfg = cfg

	/*


		g.keyMapper[pixels.KeyCapsLock] = func(_ bool) { b.KeyboardCapitalToggle() }
		g.keyMapper[pixels.KeyEscape] = func(p bool) { b.KeyboardSetKey(p, component.VKEscape) }

		g.keyMapper[pixels.KeyF9] = func(p bool) {
			if p {
				g.joyKeys = !g.joyKeys
				log.Printf("joyKeys changed: %v", g.joyKeys)
			}
		}
		g.keyMapper[pixels.KeyF10] = func(p bool) {
			if p {
				b.JoySwap()
			}
		}
		g.keyMapper[pixels.KeyF11] = func(p bool) {
			if !p {
				if fp, err := g.cfg.SwitchDisk(); err != nil {
					log.Printf("can't switch disk: %s", err)
				} else {
					log.Printf("swapping disk: %s", fp)
				}
			}
		}
		g.keyMapper[pixels.KeyF12] = func(p bool) {
			if p {
				if !g.hasClipboard {
					return
				}
				data := clipboard.Read(clipboard.FmtText)
				g.board.KeyboardSetCommand(string(data))
			}
		}
		g.keyMapper[pixels.KeyF1] = func(p bool) { b.KeyboardSetKey(p, component.VKF1) }
		g.keyMapper[pixels.KeyF2] = func(p bool) { b.KeyboardSetKey(p, component.VKF2) }
		g.keyMapper[pixels.KeyF3] = func(p bool) { b.KeyboardSetKey(p, component.VKF3) }
		g.keyMapper[pixels.KeyF4] = func(p bool) { b.KeyboardSetKey(p, component.VKF4) }
		g.keyMapper[pixels.KeyF5] = func(p bool) { b.KeyboardSetKey(p, component.VKF5) }
		g.keyMapper[pixels.KeyF6] = func(p bool) { b.KeyboardSetKey(p, component.VKF6) }
		g.keyMapper[pixels.KeyF7] = func(p bool) { b.KeyboardSetKey(p, component.VKF7) }
		g.keyMapper[pixels.KeyF8] = func(p bool) { b.KeyboardSetKey(p, component.VKF8) }

		g.keyMapper[pixels.KeyRightControl] = func(p bool) { b.KeyboardSetKey(p, component.VKControl) }
		g.keyMapper[pixels.KeyLeftControl] = func(p bool) { b.KeyboardSetKey(p, component.VKControl) }
		g.keyMapper[pixels.KeyLeftShift] = func(p bool) { b.KeyboardSetKey(p, component.VKShift) }
		g.keyMapper[pixels.KeyRightShift] = func(p bool) { b.KeyboardSetKey(p, component.VKShift) }

	*/

	g.keyMapper[13] = func(p bool) { b.KeyboardSetKey(p, component.VKReturn) }
	g.keyMapper[127] = func(p bool) { b.KeyboardSetKey(p, component.VKDelete) }
	g.keyMapper[' '] = func(p bool) { b.KeyboardSetKey(p, component.VKSpace) }
	g.keyMapper['.'] = func(p bool) { b.KeyboardSetKey(p, component.VKPeriod) }
	g.keyMapper[':'] = func(p bool) { b.KeyboardSetKey(p, component.VKColon) }
	g.keyMapper[';'] = func(p bool) { b.KeyboardSetKey(p, component.VKSemiColon) }
	g.keyMapper['*'] = func(p bool) { b.KeyboardSetKey(p, component.VKAsterisk) }
	g.keyMapper['-'] = func(p bool) { b.KeyboardSetKey(p, component.VKMinus) }
	g.keyMapper['='] = func(p bool) { b.KeyboardSetKey(p, component.VKEqual) }
	g.keyMapper['+'] = func(p bool) { b.KeyboardSetKey(p, component.VKPlus) }
	g.keyMapper[','] = func(p bool) { b.KeyboardSetKey(p, component.VKComma) }
	g.keyMapper['"'] = func(p bool) { b.KeyboardSetKey(p, component.VKShift); b.KeyboardSetKey(p, '2') }
	g.keyMapper['$'] = func(p bool) { b.KeyboardSetKey(p, component.VKShift); b.KeyboardSetKey(p, '4') }
	g.keyMapper['^'] = func(p bool) { os.Exit(0) }

	g.keyMapper['0'] = func(p bool) { b.KeyboardSetKey(p, '0') }
	g.keyMapper['1'] = func(p bool) { b.KeyboardSetKey(p, '1') }
	g.keyMapper['2'] = func(p bool) { b.KeyboardSetKey(p, '2') }
	g.keyMapper['3'] = func(p bool) { b.KeyboardSetKey(p, '3') }
	g.keyMapper['4'] = func(p bool) { b.KeyboardSetKey(p, '4') }
	g.keyMapper['5'] = func(p bool) { b.KeyboardSetKey(p, '5') }
	g.keyMapper['6'] = func(p bool) { b.KeyboardSetKey(p, '6') }
	g.keyMapper['7'] = func(p bool) { b.KeyboardSetKey(p, '7') }
	g.keyMapper['8'] = func(p bool) { b.KeyboardSetKey(p, '8') }
	g.keyMapper['9'] = func(p bool) { b.KeyboardSetKey(p, '9') }

	g.keyMapper['A'] = func(p bool) { b.KeyboardSetKey(p, 'A') }
	g.keyMapper['B'] = func(p bool) { b.KeyboardSetKey(p, 'B') }
	g.keyMapper['C'] = func(p bool) { b.KeyboardSetKey(p, 'C') }
	g.keyMapper['D'] = func(p bool) { b.KeyboardSetKey(p, 'D') }
	g.keyMapper['E'] = func(p bool) { b.KeyboardSetKey(p, 'E') }
	g.keyMapper['F'] = func(p bool) { b.KeyboardSetKey(p, 'F') }
	g.keyMapper['G'] = func(p bool) { b.KeyboardSetKey(p, 'G') }
	g.keyMapper['H'] = func(p bool) { b.KeyboardSetKey(p, 'H') }
	g.keyMapper['I'] = func(p bool) { b.KeyboardSetKey(p, 'I') }
	g.keyMapper['J'] = func(p bool) { b.KeyboardSetKey(p, 'J') }
	g.keyMapper['K'] = func(p bool) { b.KeyboardSetKey(p, 'K') }
	g.keyMapper['L'] = func(p bool) { b.KeyboardSetKey(p, 'L') }
	g.keyMapper['M'] = func(p bool) { b.KeyboardSetKey(p, 'M') }
	g.keyMapper['N'] = func(p bool) { b.KeyboardSetKey(p, 'N') }
	g.keyMapper['O'] = func(p bool) { b.KeyboardSetKey(p, 'O') }
	g.keyMapper['P'] = func(p bool) { b.KeyboardSetKey(p, 'P') }
	g.keyMapper['Q'] = func(p bool) { b.KeyboardSetKey(p, 'Q') }
	g.keyMapper['R'] = func(p bool) { b.KeyboardSetKey(p, 'R') }
	g.keyMapper['S'] = func(p bool) { b.KeyboardSetKey(p, 'S') }
	g.keyMapper['T'] = func(p bool) { b.KeyboardSetKey(p, 'T') }
	g.keyMapper['U'] = func(p bool) { b.KeyboardSetKey(p, 'U') }
	g.keyMapper['V'] = func(p bool) { b.KeyboardSetKey(p, 'V') }
	g.keyMapper['W'] = func(p bool) { b.KeyboardSetKey(p, 'W') }
	g.keyMapper['X'] = func(p bool) { b.KeyboardSetKey(p, 'X') }
	g.keyMapper['Y'] = func(p bool) { b.KeyboardSetKey(p, 'Y') }
	g.keyMapper['Z'] = func(p bool) { b.KeyboardSetKey(p, 'Z') }

	g.keyMapper['a'] = func(p bool) { b.KeyboardSetKey(p, 'A') }
	g.keyMapper['b'] = func(p bool) { b.KeyboardSetKey(p, 'B') }
	g.keyMapper['c'] = func(p bool) { b.KeyboardSetKey(p, 'C') }
	g.keyMapper['d'] = func(p bool) { b.KeyboardSetKey(p, 'D') }
	g.keyMapper['e'] = func(p bool) { b.KeyboardSetKey(p, 'E') }
	g.keyMapper['f'] = func(p bool) { b.KeyboardSetKey(p, 'F') }
	g.keyMapper['g'] = func(p bool) { b.KeyboardSetKey(p, 'G') }
	g.keyMapper['h'] = func(p bool) { b.KeyboardSetKey(p, 'H') }
	g.keyMapper['i'] = func(p bool) { b.KeyboardSetKey(p, 'I') }
	g.keyMapper['j'] = func(p bool) { b.KeyboardSetKey(p, 'J') }
	g.keyMapper['k'] = func(p bool) { b.KeyboardSetKey(p, 'K') }
	g.keyMapper['l'] = func(p bool) { b.KeyboardSetKey(p, 'L') }
	g.keyMapper['m'] = func(p bool) { b.KeyboardSetKey(p, 'M') }
	g.keyMapper['n'] = func(p bool) { b.KeyboardSetKey(p, 'N') }
	g.keyMapper['o'] = func(p bool) { b.KeyboardSetKey(p, 'O') }
	g.keyMapper['p'] = func(p bool) { b.KeyboardSetKey(p, 'P') }
	g.keyMapper['q'] = func(p bool) { b.KeyboardSetKey(p, 'Q') }
	g.keyMapper['r'] = func(p bool) { b.KeyboardSetKey(p, 'R') }
	g.keyMapper['s'] = func(p bool) { b.KeyboardSetKey(p, 'S') }
	g.keyMapper['t'] = func(p bool) { b.KeyboardSetKey(p, 'T') }
	g.keyMapper['u'] = func(p bool) { b.KeyboardSetKey(p, 'U') }
	g.keyMapper['v'] = func(p bool) { b.KeyboardSetKey(p, 'V') }
	g.keyMapper['w'] = func(p bool) { b.KeyboardSetKey(p, 'W') }
	g.keyMapper['x'] = func(p bool) { b.KeyboardSetKey(p, 'X') }
	g.keyMapper['y'] = func(p bool) { b.KeyboardSetKey(p, 'Y') }
	g.keyMapper['z'] = func(p bool) { b.KeyboardSetKey(p, 'Z') }

	/*
		g.keyMapper[Up] = func(p bool) {
			if g.joyKeys {
				b.Joy1SetKey(p, component.KeyJUp)
			} else {
				b.KeyboardSetKey(p, component.VKUp)
			}
		}
		g.keyMapper[pixels.KeyDown] = func(p bool) {
			if g.joyKeys {
				b.Joy1SetKey(p, component.KeyJDown)
			} else {
				b.KeyboardSetKey(p, component.VKDown)
			}
		}
		g.keyMapper[pixels.KeyLeft] = func(p bool) {
			if g.joyKeys {
				b.Joy1SetKey(p, component.KeyJLeft)
			} else {
				b.KeyboardSetKey(p, component.VKLeft)
			}
		}
		g.keyMapper[pixels.KeyRight] = func(p bool) {
			if g.joyKeys {
				b.Joy1SetKey(p, component.KeyJRight)
			} else {
				b.KeyboardSetKey(p, component.VKRight)
			}
		}
		g.keyMapper[pixels.KeyTab] = func(p bool) {
			if g.joyKeys {
				b.Joy1SetKey(p, component.KeyJFire)
			} else {
				b.KeyboardSetKey(p, component.VKTab)
			}
		}


		g.keyMapper[pixels.MouseButton1] = func(p bool) { b.Joy1SetKey(p, component.KeyJFire) }
		g.keyMapper[pixels.MouseButton2] = func(p bool) { b.Joy1SetKey(p, component.KeyJUp) }
	*/

	err := clipboard.Init()
	if err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		g.hasClipboard = true
	}
	return nil
}

func (g *Inputs) Key(p byte, pressed bool) {
	if f, ok := g.keyMapper[p]; ok {
		f(pressed)
	}
}

package glrender

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/render/glrender/pixels"
	"golang.design/x/clipboard"
	"log"
)

type Inputs struct {
	board        references.IBoard
	cfg          *config.Config
	keyMapper    []func(bool)
	activeKeys   map[pixels.Button]bool
	joyKeys      bool
	lastX        uint8
	lastY        uint8
	hasClipboard bool
}

func NewInputs() *Inputs {
	return &Inputs{
		board:        nil,
		cfg:          nil,
		keyMapper:    nil,
		joyKeys:      true,
		activeKeys:   make(map[pixels.Button]bool),
		lastX:        0,
		lastY:        0,
		hasClipboard: false,
	}
}

func (g *Inputs) Setup(b references.IBoard, cfg *config.Config) error {
	const max = int(pixels.KeyLast + 1)
	g.board = b
	g.cfg = cfg
	g.keyMapper = make([]func(bool), max)
	for x := 0; x < max; x++ {
		g.keyMapper[x] = func(b bool) {}
	}

	g.keyMapper[pixels.KeyCapsLock] = func(_ bool) { b.KeyboardCapitalToggle() }
	g.keyMapper[pixels.KeyEscape] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKEscape) }

	g.keyMapper[pixels.KeyF9] = func(p bool) {
		if p {
			g.joyKeys = !g.joyKeys
			fmt.Println("joyKeys", g.joyKeys)
		}
	}
	g.keyMapper[pixels.KeyF10] = func(p bool) {
		if p {
			b.JoySwap()
		}
	}
	g.keyMapper[pixels.KeyF11] = func(p bool) {
		if p {
			g.cfg.SwitchDisk()
			g.cfg.SetDriveOpt("", 8)
			fmt.Println("swapping disk")
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
	g.keyMapper[pixels.KeyF1] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF1) }
	g.keyMapper[pixels.KeyF2] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF2) }
	g.keyMapper[pixels.KeyF3] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF3) }
	g.keyMapper[pixels.KeyF4] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF4) }
	g.keyMapper[pixels.KeyF5] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF5) }
	g.keyMapper[pixels.KeyF6] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF6) }
	g.keyMapper[pixels.KeyF7] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF7) }
	g.keyMapper[pixels.KeyF8] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKF8) }

	g.keyMapper[pixels.KeyRightControl] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKControl) }
	g.keyMapper[pixels.KeyLeftControl] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKControl) }
	g.keyMapper[pixels.KeyLeftShift] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKShift) }
	g.keyMapper[pixels.KeyRightShift] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKShift) }
	g.keyMapper[pixels.KeyEnter] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKReturn) }
	g.keyMapper[pixels.KeyDelete] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKDelete) }
	g.keyMapper[pixels.KeyBackspace] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKBack) }
	g.keyMapper[pixels.KeySpace] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKSpace) }
	g.keyMapper[pixels.KeyComma] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKComma) }
	g.keyMapper[pixels.KeyPeriod] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKPeriod) }
	g.keyMapper[pixels.KeySemicolon] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKSemicolon) }
	g.keyMapper[pixels.KeyApostrophe] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKQuote) }
	g.keyMapper[pixels.KeyRightBracket] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKAsterisk) }
	g.keyMapper[pixels.KeyLeftBracket] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKAsterisk) }

	g.keyMapper[pixels.KeyMinus] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKMinus) }
	g.keyMapper[pixels.KeyEqual] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKEqual) }
	g.keyMapper[pixels.KeyBackslash] = func(p bool) { b.KeyboardSetVirtualKey(p, component.VKPlus) }

	g.keyMapper[pixels.Key0] = func(p bool) { b.KeyboardSetVirtualKey(p, '0') }
	g.keyMapper[pixels.Key1] = func(p bool) { b.KeyboardSetVirtualKey(p, '1') }
	g.keyMapper[pixels.Key2] = func(p bool) { b.KeyboardSetVirtualKey(p, '2') }
	g.keyMapper[pixels.Key3] = func(p bool) { b.KeyboardSetVirtualKey(p, '3') }
	g.keyMapper[pixels.Key4] = func(p bool) { b.KeyboardSetVirtualKey(p, '4') }
	g.keyMapper[pixels.Key5] = func(p bool) { b.KeyboardSetVirtualKey(p, '5') }
	g.keyMapper[pixels.Key6] = func(p bool) { b.KeyboardSetVirtualKey(p, '6') }
	g.keyMapper[pixels.Key7] = func(p bool) { b.KeyboardSetVirtualKey(p, '7') }
	g.keyMapper[pixels.Key8] = func(p bool) { b.KeyboardSetVirtualKey(p, '8') }
	g.keyMapper[pixels.Key9] = func(p bool) { b.KeyboardSetVirtualKey(p, '9') }

	g.keyMapper[pixels.KeyA] = func(p bool) { b.KeyboardSetVirtualKey(p, 'A') }
	g.keyMapper[pixels.KeyB] = func(p bool) { b.KeyboardSetVirtualKey(p, 'B') }
	g.keyMapper[pixels.KeyC] = func(p bool) { b.KeyboardSetVirtualKey(p, 'C') }
	g.keyMapper[pixels.KeyD] = func(p bool) { b.KeyboardSetVirtualKey(p, 'D') }
	g.keyMapper[pixels.KeyE] = func(p bool) { b.KeyboardSetVirtualKey(p, 'E') }
	g.keyMapper[pixels.KeyF] = func(p bool) { b.KeyboardSetVirtualKey(p, 'F') }
	g.keyMapper[pixels.KeyG] = func(p bool) { b.KeyboardSetVirtualKey(p, 'G') }
	g.keyMapper[pixels.KeyH] = func(p bool) { b.KeyboardSetVirtualKey(p, 'H') }
	g.keyMapper[pixels.KeyI] = func(p bool) { b.KeyboardSetVirtualKey(p, 'I') }
	g.keyMapper[pixels.KeyJ] = func(p bool) { b.KeyboardSetVirtualKey(p, 'J') }
	g.keyMapper[pixels.KeyK] = func(p bool) { b.KeyboardSetVirtualKey(p, 'K') }
	g.keyMapper[pixels.KeyL] = func(p bool) { b.KeyboardSetVirtualKey(p, 'L') }
	g.keyMapper[pixels.KeyM] = func(p bool) { b.KeyboardSetVirtualKey(p, 'M') }
	g.keyMapper[pixels.KeyN] = func(p bool) { b.KeyboardSetVirtualKey(p, 'N') }
	g.keyMapper[pixels.KeyO] = func(p bool) { b.KeyboardSetVirtualKey(p, 'O') }
	g.keyMapper[pixels.KeyP] = func(p bool) { b.KeyboardSetVirtualKey(p, 'P') }
	g.keyMapper[pixels.KeyQ] = func(p bool) { b.KeyboardSetVirtualKey(p, 'Q') }
	g.keyMapper[pixels.KeyR] = func(p bool) { b.KeyboardSetVirtualKey(p, 'R') }
	g.keyMapper[pixels.KeyS] = func(p bool) { b.KeyboardSetVirtualKey(p, 'S') }
	g.keyMapper[pixels.KeyT] = func(p bool) { b.KeyboardSetVirtualKey(p, 'T') }
	g.keyMapper[pixels.KeyU] = func(p bool) { b.KeyboardSetVirtualKey(p, 'U') }
	g.keyMapper[pixels.KeyV] = func(p bool) { b.KeyboardSetVirtualKey(p, 'V') }
	g.keyMapper[pixels.KeyW] = func(p bool) { b.KeyboardSetVirtualKey(p, 'W') }
	g.keyMapper[pixels.KeyX] = func(p bool) { b.KeyboardSetVirtualKey(p, 'X') }
	g.keyMapper[pixels.KeyY] = func(p bool) { b.KeyboardSetVirtualKey(p, 'Y') }
	g.keyMapper[pixels.KeyZ] = func(p bool) { b.KeyboardSetVirtualKey(p, 'Z') }

	g.keyMapper[pixels.KeyUp] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJUp)
		} else {
			b.KeyboardSetVirtualKey(p, component.VKUp)
		}
	}
	g.keyMapper[pixels.KeyDown] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJDown)
		} else {
			b.KeyboardSetVirtualKey(p, component.VKDown)
		}
	}
	g.keyMapper[pixels.KeyLeft] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJLeft)
		} else {
			b.KeyboardSetVirtualKey(p, component.VKLeft)
		}
	}
	g.keyMapper[pixels.KeyRight] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJRight)
		} else {
			b.KeyboardSetVirtualKey(p, component.VKRight)
		}
	}
	g.keyMapper[pixels.KeyTab] = func(p bool) {
		if g.joyKeys {
			b.Joy1SetKey(p, component.KeyJFire)
		} else {
			b.KeyboardSetVirtualKey(p, component.VKTab)
		}
	}
	g.keyMapper[pixels.MouseButton1] = func(p bool) { b.Joy1SetKey(p, component.KeyJFire) }
	g.keyMapper[pixels.MouseButton2] = func(p bool) { b.Joy1SetKey(p, component.KeyJUp) }

	err := clipboard.Init()
	if err != nil {
		log.Printf("can't init clipboard: %s", err)
	} else {
		g.hasClipboard = true
	}
	return nil
}

func (g *Inputs) Keys(pressed map[pixels.Button]bool) {
	if len(g.activeKeys) > 0 {
		for v := range g.activeKeys {
			if _, ok := pressed[v]; !ok {
				g.keyMapper[v](false)
				delete(g.activeKeys, v)
			}
		}
	}
	if len(pressed) > 0 {
		for v := range pressed {
			if _, ok := g.activeKeys[v]; ok {
				continue
			}
			g.activeKeys[v] = true
			g.keyMapper[v](true)
		}
	}
}

func (g *Inputs) MouseMove(x float64, y float64) {
	x1 := uint8(x)
	y1 := uint8(y)
	if g.lastX != x1 || g.lastY != y1 {
		g.lastX = x1
		g.lastY = y1
		g.board.SetMouse(uint8(x), uint8(y))
	}
}

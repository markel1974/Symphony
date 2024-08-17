package render

import (
	"github.com/markel1974/c64emu/src/board"
	"github.com/markel1974/c64emu/src/board/keyboard"
	"github.com/markel1974/c64emu/src/pixels"
)

type Keyboard struct {
	keyMapper  []func(bool)
	activeKeys map[pixels.Button]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		keyMapper:  nil,
		activeKeys: make(map[pixels.Button]bool),
	}
}

func (g *Keyboard) Setup(b *board.Board) {
	const max = int(pixels.KeyLast + 1)
	g.keyMapper = make([]func(bool), max)
	for x := 0; x < max; x++ {
		g.keyMapper[x] = func(b bool) {}
	}

	g.keyMapper[pixels.KeyCapsLock] = func(p bool) { b.KeyboardSetCapital(p) }
	g.keyMapper[pixels.KeyEscape] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_escape) }
	g.keyMapper[pixels.KeyF8] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_up) }
	g.keyMapper[pixels.KeyF9] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_down) }

	//g.keyMapper[pixels.KeyF9] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_bracketright) }
	g.keyMapper[pixels.KeyF10] = func(p bool) { b.KeyboardSwapJoystick(p) }
	g.keyMapper[pixels.KeyF11] = func(p bool) { b.KeyboardSetExt(p) }
	g.keyMapper[pixels.KeyF12] = func(p bool) { b.KeyboardPaste(p) }

	g.keyMapper[pixels.KeyF1] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f1) }
	g.keyMapper[pixels.KeyF2] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f2) }
	g.keyMapper[pixels.KeyF3] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f3) }
	g.keyMapper[pixels.KeyF4] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f4) }
	g.keyMapper[pixels.KeyF5] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f5) }
	g.keyMapper[pixels.KeyF6] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f6) }
	g.keyMapper[pixels.KeyF7] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f7) }
	g.keyMapper[pixels.KeyF8] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_f8) }

	g.keyMapper[pixels.KeyRightControl] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_control) }
	g.keyMapper[pixels.KeyLeftControl] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_control) }
	g.keyMapper[pixels.KeyLeftShift] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_shift) }
	g.keyMapper[pixels.KeyRightShift] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_shift) }
	g.keyMapper[pixels.KeyEnter] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_return) }
	g.keyMapper[pixels.KeyDelete] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_delete) }
	g.keyMapper[pixels.KeyBackspace] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_back) }
	g.keyMapper[pixels.KeySpace] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_space) }
	g.keyMapper[pixels.KeyComma] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_comma) }
	g.keyMapper[pixels.KeyPeriod] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_period) }
	g.keyMapper[pixels.KeySemicolon] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_semicolon) }
	g.keyMapper[pixels.KeyApostrophe] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_quote) }
	g.keyMapper[pixels.KeyRightBracket] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_asterisk) }

	g.keyMapper[pixels.KeyUp] = func(p bool) { b.KeyboardSetJoyKey(p, keyboard.KEY_JUP) }
	g.keyMapper[pixels.KeyDown] = func(p bool) { b.KeyboardSetJoyKey(p, keyboard.KEY_JDN) }
	g.keyMapper[pixels.KeyLeft] = func(p bool) { b.KeyboardSetJoyKey(p, keyboard.KEY_JLF) }
	g.keyMapper[pixels.KeyRight] = func(p bool) { b.KeyboardSetJoyKey(p, keyboard.KEY_JRT) }
	g.keyMapper[pixels.KeyTab] = func(p bool) { b.KeyboardSetJoyKey(p, keyboard.KEY_FIRE) }

	g.keyMapper[pixels.KeyMinus] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_minus) }
	g.keyMapper[pixels.KeyEqual] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_equal) }
	g.keyMapper[pixels.KeyBackslash] = func(p bool) { b.KeyboardSetVirtualKey(p, keyboard.VK_plus) }

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
}

func (g *Keyboard) Keys(pressed map[pixels.Button]bool) {
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
			//fmt.Println("CURRENT", v)
			g.activeKeys[v] = true
			g.keyMapper[v](true)
		}
	}
}

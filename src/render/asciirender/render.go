package asciirender

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"os"
)

type Render struct {
	board      references.IBoard
	player     references.IAudioRender
	cfg        *config.Config
	display    *DisplayBuffer
	textBuffer []byte
	ch         chan []byte
	run        bool
}

func New() *Render {
	g := &Render{
		board:      nil,
		cfg:        nil,
		display:    nil,
		textBuffer: make([]byte, 65000),
		ch:         make(chan []byte),
	}
	return g
}

func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.display = NewDisplayBuffer()
	return g.display, nil
}

func (g *Render) Start(board references.IBoard) error {
	g.board = board
	g.board.VBlankSignal().Bind(g.vBlank)
	if err := MakeStdInRaw(); err != nil {
		return err
	}
	go func() {
		textData := make([]byte, 16)
		for {
			if n, err := os.Stdin.Read(textData); err == nil && n > 0 {
				g.ch <- textData[:n]
			}
		}
	}()
	for g.run {
		g.board.Emulate()
	}
	return nil
}

func (g *Render) vBlank() {
	dt := g.board.Throttle()
	dt.Throttle()

	select {
	case text := <-g.ch:
		switch text[0] {
		case 'A':
			g.board.Joy1SetKey(true, component.KeyJLeft)
			g.board.Joy1SetKey(false, component.KeyJLeft)
		case 'D':
			g.board.Joy1SetKey(true, component.KeyJRight)
			g.board.Joy1SetKey(false, component.KeyJRight)
		case 'W':
			g.board.Joy1SetKey(true, component.KeyJUp)
			g.board.Joy1SetKey(false, component.KeyJUp)
		case 'S':
			g.board.Joy1SetKey(true, component.KeyJDown)
			g.board.Joy1SetKey(false, component.KeyJDown)
		case 'F':
			g.board.Joy1SetKey(true, component.KeyJFire)
			g.board.Joy1SetKey(false, component.KeyJFire)
		case 'Q':
			g.run = false
		default:
			g.board.KeyboardSetCommand(string(text))
		}
	default:
	}
	b := g.board.GetText()
	if v := bytes.Compare(b, g.textBuffer[:len(b)]); v == 0 {
		return
	}
	g.printBuffer(b)
	copy(g.textBuffer, b)
	//counter++
	//fmt.Printf("vblank %d\r\n", counter)
}

func (g *Render) printBuffer(textBuffer []byte) {
	const asciiEsc = 27
	const newLine = "\r\n"
	fmt.Printf("%c[2J", asciiEsc)
	for x, v := range textBuffer {
		if (x % 40) == 0 {
			fmt.Print(newLine)
		}
		fmt.Printf("%c", v)
	}
}

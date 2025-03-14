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
	cfg          *config.Config
	c64Board     references.IBoard
	scale        float64
	fullscreen   bool
	showMap      bool
	screenWidth  int
	screenHeight int
	maxW         float64
	maxH         float64
	display      *DisplayBuffer
	audio        *Audio
}

func New(c64board references.IBoard, cfg *config.Config) *Render {
	w, h := c64board.GetScreenSize()
	g := &Render{
		c64Board:     c64board,
		cfg:          cfg,
		fullscreen:   false,
		screenWidth:  w,
		screenHeight: h,
		scale:        2.0,
		display:      NewDisplayBuffer(),
		audio:        NewAudio(),
	}
	g.maxW = float64(g.screenWidth) * g.scale
	g.maxH = float64(g.screenHeight) * g.scale
	return g
}

func (g *Render) setup() {
	//g.c64Board = board.NewBoard()
	_ = g.c64Board.Setup(g.display, g.audio, g.cfg)
	//g.inputs.Setup(g.c64Board, g.maxW, g.maxH)
}

func (g *Render) PrintBuffer(textBuffer []byte) {
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

func (g *Render) Start() error {
	g.setup()
	dt := g.c64Board.Throttle()

	if err := MakeStdInRaw(); err != nil {
		return err
	}

	run := true

	ch := make(chan []byte)

	go func() {
		textData := make([]byte, 16)
		for {
			if n, err := os.Stdin.Read(textData); err == nil && n > 0 {
				ch <- textData[:n]
			}

		}
	}()

	counter := 0
	textBuffer := make([]byte, 65000)
	for run {
		dt.Throttle()

		select {
		case text := <-ch:
			switch text[0] {
			case 'A':
				g.c64Board.Joy1SetKey(true, component.KeyJLeft)
				g.c64Board.Joy1SetKey(false, component.KeyJLeft)
			case 'D':
				g.c64Board.Joy1SetKey(true, component.KeyJRight)
				g.c64Board.Joy1SetKey(false, component.KeyJRight)
			case 'W':
				g.c64Board.Joy1SetKey(true, component.KeyJUp)
				g.c64Board.Joy1SetKey(false, component.KeyJUp)
			case 'S':
				g.c64Board.Joy1SetKey(true, component.KeyJDown)
				g.c64Board.Joy1SetKey(false, component.KeyJDown)
			case 'F':
				g.c64Board.Joy1SetKey(true, component.KeyJFire)
				g.c64Board.Joy1SetKey(false, component.KeyJFire)
			case 'Q':
				run = false
			default:
				g.c64Board.KeyboardSetCommand(string(text))
			}
		default:
		}

		for {
			if vBlank := g.c64Board.Emulate(); vBlank {
				break
			}
		}
		b := g.c64Board.GetText()
		if v := bytes.Compare(b, textBuffer[:len(b)]); v == 0 {
			continue
		}
		g.PrintBuffer(b)
		copy(textBuffer, b)
		counter++
		fmt.Printf("vblank %d\r\n", counter)
	}
	return nil
}

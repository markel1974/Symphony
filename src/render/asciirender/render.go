package asciirender

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/c64/board"
	"github.com/markel1974/c64emu/src/components/vic"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/render/common"
	"os"
)

type Render struct {
	cfg          *config.Config
	c64Board     *board.Board
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

func New(cfg *config.Config) *Render {
	g := &Render{
		cfg:          cfg,
		fullscreen:   false,
		screenWidth:  mos6569.DisplayX,
		screenHeight: mos6569.DisplayY,
		scale:        2.0,
		display:      NewDisplayBuffer(),
		audio:        NewAudio(),
	}
	g.maxW = float64(g.screenWidth) * g.scale
	g.maxH = float64(g.screenHeight) * g.scale
	return g
}

func (g *Render) setup() {
	g.c64Board = board.NewBoard(g.display, g.audio)
	_ = g.c64Board.Setup(g.cfg)
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

func (g *Render) Start() {
	g.setup()
	dt := common.NewDynamicThrottling(mos6569.FrameInterval)

	MakeStdInRaw()

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
		dt.DynamicThrottling()

		select {
		case text := <-ch:
			g.c64Board.KeyboardSetCommand(string(text))
			//TODO REMOVE
			if string(text) == "Q" {
				run = false
			}
			//g.c64Board.Joy1SetKey(true, inputs.KeyJFire)
			//g.c64Board.Joy1SetKey(false, inputs.KeyJFire)
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
}

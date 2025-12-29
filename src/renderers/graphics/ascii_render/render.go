package ascii_render

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type Render struct {
	board      references.IC64Board
	player     references.IAudioRender
	cfg        *config.Config
	display    *DisplayBuffer
	textBuffer []byte
	ch         chan []byte
	run        bool
	input      *Inputs
}

func New() *Render {
	g := &Render{
		board:      nil,
		cfg:        nil,
		display:    nil,
		textBuffer: make([]byte, 65000),
		ch:         make(chan []byte),
		run:        true,
		input:      NewInputs(),
	}
	return g
}

func (g *Render) Setup(board references.IC64Board, cfg *config.Config) error {
	g.board = board
	if err := g.board.Wire(g); err != nil {
		return err
	}
	if err := g.input.Setup(g.board, cfg); err != nil {
		return err
	}
	if err := MakeStdInRaw(); err != nil {
		log.Printf("can't make stdin raw: %s", err)
	}
	return nil
}

func (g *Render) signalHandler() {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		_ = <-sig
		os.Exit(0)
	}()
}

func (g *Render) Start() error {
	//g.signalHandler()
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

func (g *Render) VBlank() {
	select {
	case text := <-g.ch:
		for _, v := range text {
			//fmt.Println(v)
			g.input.Key(v, true)
			g.input.Key(v, false)
		}

	//switch text[0] {
	/*
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

	*/
	//default:
	//g.board.KeyboardSetCommand(string(text))
	//}
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

func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.display = NewDisplayBuffer()
	return g.display, nil
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

func (g *Render) LedActivity(deviceNumber uint8, led bool) {
	//device := uint8(state & 0xf)
	//led := uint8((state >> 8) & 0xf)
	//fmt.Println("LED STATE", device, led)
	fmt.Println("LED STATE", deviceNumber, led)
}

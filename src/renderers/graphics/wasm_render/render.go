//go:build js && wasm

package wasm_render

import (
	"syscall/js"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

//TODO 5000 cicli max per emulate e creare callback per vblank
//Fette Ragionevoli: Dividere il lavoro di un VBlank (19656 cicli) in 2-4 fette (quindi ogni fetta esegue 5000-10000 cicli C64 e dura circa 5-10ms di tempo reale) è spesso un buon compromesso. Questo crea 100-200 "slot liberi" al secondo per il loop dell'emulatore, il che è ampiamente sufficiente per far girare uno scheduler audio con un intervallo di 10ms in modo molto più affidabile.

// Render manages the display rendering process and user interaction through inputs and a linked board implementation.
// It integrates with a DisplayBuffer for screen updates and relies on IC64Board for controlling board operations.
type Render struct {
	board         references.IC64Board
	displayBuffer *DisplayBuffer
	w             int
	h             int
	input         *Inputs
	onVBlank      js.Value
	emulateRate   int
}

// NewRender initializes and returns a pointer to a new instance of the Render struct with default values.
func NewRender() *Render {
	return &Render{
		displayBuffer: nil, //NewDisplayBuffer(320, 200),
		board:         nil,
		input:         NewInputs(),
		emulateRate:   5000,
	}
}

// CreateDisplayBuffer initializes a new display buffer with the specified width and height, stores it, and returns it.
func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.w = w
	g.h = h
	g.displayBuffer = NewDisplayBuffer(w, h)
	return g.displayBuffer, nil
}

// Setup initializes the Render object with the given IC64Board and Config instances, setting up the necessary components.
// Returns an error if the board or input setup fails.
func (g *Render) Setup(board references.IC64Board, cfg *config.Config) error {
	g.board = board
	if err := g.board.Mount(g); err != nil {
		return err
	}
	if err := g.input.Setup(g.board, cfg); err != nil {
		return err
	}
	if err := g.initWasm(); err != nil {
		return err
	}
	return nil
}

// Start initializes the rendering process and exposes Go functions for JavaScript interaction via the global scope.
// It blocks indefinitely using a channel to keep the Go routine alive.
func (g *Render) Start() error {
	c := make(chan struct{})
	<-c
	return nil
}

// VBlank sets the vBlank flag to true, signaling the start of the vertical blanking interval in the rendering process.
func (g *Render) VBlank() {
	g.onVBlank.Invoke()
}

// LedActivity toggles the LED state for the specified device number based on the given boolean value.
func (g *Render) LedActivity(deviceNumber uint8, led bool) {
	//g.led = led
	//fmt.Println("LED STATE", deviceNumber, led)
}

func (g *Render) initWasm() error {
	initRender := func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			println("Error: initRender emulateRate, vBlankCallback")
			return nil
		}
		g.emulateRate = args[0].Int()
		g.onVBlank = args[1]
		return nil
	}
	emulate := func(this js.Value, args []js.Value) interface{} {
		for x := 0; x < g.emulateRate; x++ {
			g.board.Emulate()
		}
		return nil
	}
	getSurfacePointer := func(this js.Value, args []js.Value) interface{} {
		surfacePtr := g.displayBuffer.GetSurfacePointer()
		return uintptr(surfacePtr)
	}
	getSurfaceLen := func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(g.displayBuffer.GetSurfaceLen())
	}
	getDisplayBuffer := func(this js.Value, args []js.Value) interface{} {
		surfaceBytes := g.displayBuffer.GetSurface()
		jsBuffer := js.Global().Get("Uint8Array").New(len(surfaceBytes))
		js.CopyBytesToJS(jsBuffer, surfaceBytes)
		return jsBuffer
	}
	getDisplayWidth := func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(g.w)
	}
	getDisplayHeight := func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(g.h)
	}
	keyPressed := func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			if args[0].Type() == js.TypeString {
				keyCode := args[0].String()
				g.input.Key(keyCode, true)
			}
		}
		return nil
	}
	keyReleased := func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			if args[0].Type() == js.TypeString {
				keyCode := args[0].String()
				g.input.Key(keyCode, false)
			}
		}
		return nil
	}
	js.Global().Set("initRender", js.FuncOf(initRender))
	js.Global().Set("getDisplayBuffer", js.FuncOf(getDisplayBuffer))
	js.Global().Set("getDisplayWidth", js.FuncOf(getDisplayWidth))
	js.Global().Set("getDisplayHeight", js.FuncOf(getDisplayHeight))
	js.Global().Set("getSurfacePointer", js.FuncOf(getSurfacePointer))
	js.Global().Set("getSurfaceLen", js.FuncOf(getSurfaceLen))
	js.Global().Set("emulate", js.FuncOf(emulate))
	js.Global().Set("keyPressed", js.FuncOf(keyPressed))
	js.Global().Set("keyReleased", js.FuncOf(keyReleased))
	return nil
}

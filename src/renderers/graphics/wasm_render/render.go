//go:build js && wasm

package wasm_render

import (
	"syscall/js"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

/*
//go:embed g64.wasm
var wasmFS embed.FS
func main() {
	http.Handle("/", http.FileServer(http.FS(wasmFS)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/

type Render struct {
	board         references.IBoard
	displayBuffer *DisplayBuffer
	vBlank        bool
	w             int
	h             int
	input         *Inputs
}

func NewRender() *Render {
	return &Render{
		displayBuffer: nil, //NewDisplayBuffer(320, 200),
		board:         nil,
		vBlank:        false,
		input:         NewInputs(),
	}
}

func (g *Render) Setup(board references.IBoard, cfg *config.Config) error {
	g.board = board
	if err := g.board.Mount(g); err != nil {
		return err
	}
	if err := g.input.Setup(g.board, cfg); err != nil {
		return err
	}
	return nil
}

func (g *Render) CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error) {
	g.w = w
	g.h = h
	g.displayBuffer = NewDisplayBuffer(w, h)
	return g.displayBuffer, nil
}

func (g *Render) Start() error {
	c := make(chan struct{}, 0)
	emulateFrame := func(this js.Value, args []js.Value) interface{} {
		for {
			g.board.Emulate()
			if g.vBlank {
				g.vBlank = false
				return nil
			}
		}
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

	// Esponi le funzioni Go a JavaScript.
	js.Global().Set("getDisplayBuffer", js.FuncOf(getDisplayBuffer))
	js.Global().Set("getDisplayWidth", js.FuncOf(getDisplayWidth))
	js.Global().Set("getDisplayHeight", js.FuncOf(getDisplayHeight))
	js.Global().Set("getSurfacePointer", js.FuncOf(getSurfacePointer))
	js.Global().Set("getSurfaceLen", js.FuncOf(getSurfaceLen))
	js.Global().Set("emulateFrame", js.FuncOf(emulateFrame))
	js.Global().Set("keyPressed", js.FuncOf(keyPressed))
	js.Global().Set("keyReleased", js.FuncOf(keyReleased))
	<-c

	return nil
}

func (g *Render) VBlank() {
	g.vBlank = true
	//if g.win.MouseInsideWindow() {
	//	g.inputs.MouseMove(g.win.MousePositionXY())
	//}
	//g.inputs.Keys(g.win.KeysPressed())
	//g.surface.Draw(g.win, g.surfaceM)
	//if g.led {
	//	g.ledSurface.Draw(g.win, g.ledSurfaceM)
	//}
	//g.win.Update()
	//g.run = !g.win.Closed()
}

func (g *Render) LedActivity(deviceNumber uint8, led bool) {
	//g.led = led
	//fmt.Println("LED STATE", deviceNumber, led)
}

//TODO WASM
// https://garciat.com/posts/go-wasm/
// https://github.com/seqsense/webgl-go/tree/master

//go:build js && wasm

package wasm_render

import (
	"fmt"
	"syscall/js"
	"unsafe"

	"github.com/markel1974/c64emu/src/config"
)

// Audio represents a structure handling audio virtualization and integration with JavaScript AudioContext.
type Audio struct {
	cfg            *config.Config
	audioCtx       js.Value  // Reference to the JavaScript AudioContext
	sampleRate     float64   // Sample rate of the AudioContext
	goBuffer       []float32 // Buffer in Go to accumulate samples
	onSamplesReady js.Value  // JS callback function to send samples
}

// NewAudio creates and returns a new instance of the Audio structure with default initialization for playback management.
func NewAudio() *Audio {
	return &Audio{}
}

// Setup initializes the Audio instance with the given configuration and prepares it for further setup steps.
func (a *Audio) Setup(cfg *config.Config) error {
	a.cfg = cfg
	if err := a.initWasm(); err != nil {
		return err
	}
	return nil
}

// Audio.Write modificato per gestire l'output attuale di calcBuffer
func (a *Audio) Write(bufferInput []float32, samples int) {
	if a.audioCtx.IsUndefined() {
		fmt.Println("Go Audio.Write: Errore - AudioContext non inizializzato")
		return
	}
	if a.onSamplesReady.IsUndefined() || !a.onSamplesReady.Truthy() {
		fmt.Println("Go Audio.Write: Errore - OnSamplesReady non inizializzato")
		return
	}

	if len(a.goBuffer) != len(bufferInput) {
		a.goBuffer = make([]float32, len(bufferInput))
	}
	copy(a.goBuffer, bufferInput[:samples])

	header := (*[3]uintptr)(unsafe.Pointer(&a.goBuffer))
	byteSlice := (*[1 << 30]byte)(unsafe.Pointer(header[0]))[: samples*4 : samples*4]

	jsArrayBuffer := js.Global().Get("ArrayBuffer").New(len(byteSlice))
	jsUint8Array := js.Global().Get("Uint8Array").New(jsArrayBuffer)

	js.CopyBytesToJS(jsUint8Array, byteSlice)

	jsFloat32Array := js.Global().Get("Float32Array").New(jsArrayBuffer)

	a.onSamplesReady.Invoke(jsFloat32Array)
}

// Play resumes audio playback if the AudioContext is suspended. Resets or starts AudioBufferSourceNode if needed.
func (a *Audio) Play() {
	if a.audioCtx.IsUndefined() {
		fmt.Println("Play Error: AudioContext not initialized")
	}
	if a.audioCtx.Truthy() && a.audioCtx.Get("state").String() == "suspended" {
		a.audioCtx.Call("resume")
		println("Go Audio: Resumed")
	}
}

// Pause suspends audio playback if the underlying AudioContext is currently in the "running" state.
func (a *Audio) Pause() {
	if a.audioCtx.IsUndefined() {
		fmt.Println("Pause Error: AudioContext not initialized")
	}
	if a.audioCtx.Truthy() && a.audioCtx.Get("state").String() == "running" {
		a.audioCtx.Call("suspend")
		println("Go Audio: Suspended")
	}
}

// Resume restarts playback by invoking the Play method, commonly acting as a synonym for it in AudioContext scenarios.
func (a *Audio) Resume() {
	a.Play()
}

// GlobalSetupWasmAudio exports Go methods to JavaScript for initializing, playing, and pausing the AudioContext.
func (a *Audio) initWasm() error {
	// init initializes the AudioContext with provided arguments and sets up a sample-ready callback.
	init := func(this js.Value, args []js.Value) interface{} {
		if len(args) < 3 {
			println("Error: initAudioContextAndGetCallback expects audioCtx, sampleRate, callback")
			return nil
		}
		a.audioCtx = args[0]
		a.sampleRate = args[1].Float()
		a.onSamplesReady = args[2] // Callback JS: function(float32Array)
		println("Go Audio: AudioContext initialized, sampleRate:", a.sampleRate)
		return nil
	}
	play := func(this js.Value, args []js.Value) interface{} {
		a.Play()
		return nil
	}
	pause := func(this js.Value, args []js.Value) interface{} {
		a.Pause()
		return nil
	}
	flush := func(this js.Value, args []js.Value) interface{} {
		return nil
	}

	js.Global().Set("wasmAudioInit", js.FuncOf(init))
	js.Global().Set("wasmAudioPlay", js.FuncOf(play))
	js.Global().Set("wasmAudioPause", js.FuncOf(pause))
	js.Global().Set("wasmAudioFlush", js.FuncOf(flush))

	return nil
}

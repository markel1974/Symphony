//go:build js && wasm

package wasm_render

import (
	"syscall/js"
	"unsafe"

	"github.com/markel1974/c64emu/src/config"
)

// Audio represents a structure handling audio virtualization and integration with JavaScript AudioContext.
type Audio struct {
	cfg        *config.Config
	pos        int
	audioCtx   js.Value  // Reference to the JavaScript AudioContext
	sampleRate float64   // Sample rate of the AudioContext
	goBuffer   []float32 // Buffer in Go to accumulate samples
	//jsTypedArray     js.Value  // Uint8Array or Float32Array in JS to pass data
	onSamplesReady   js.Value // JS callback function to send samples
	bufferSizeCycles int      // Buffer size in emulation cycles (approximately)
	cyclesSinceFlush int
}

// NewAudio creates and returns a new instance of the Audio structure with default initialization for playback management.
func NewAudio() *Audio {
	return &Audio{
		pos: 0,
	}
}

// Setup initializes the Audio instance with the given configuration and prepares it for further setup steps.
func (a *Audio) Setup(cfg *config.Config) error {
	// Initialize audioCtx and sampleRate from JavaScript in a later phase (e.g., Start or an exported func)
	a.cfg = cfg
	if err := a.initWasm(); err != nil {
		return err
	}
	return nil
}

// GetCurrentPosition returns the current playback position of the audio in cycles or buffer updates.
func (a *Audio) GetCurrentPosition() int {
	return a.pos
}

// Write processes input audio samples, normalizes them, and appends them to the internal buffer. Sends data to JS if buffer is full.
func (a *Audio) Write(buffer []uint32, _ int, samples int) {
	if a.audioCtx.IsUndefined() || !a.onSamplesReady.Truthy() {
		// AudioContext non ancora inizializzato da JS o callback mancante
		a.pos += samples // Avanza comunque la posizione
		return
	}

	// Normalize and add samples to goBuffer
	// This is a SIMPLIFICATION. Normalization depends on the SID's output range.
	// The C64 SID has a complex output, often a filter is used and then normalized.
	// Here we assume that the uint32 values are already in a range we can map to float32.
	for i := 0; i < samples; i++ {
		// Esempio di normalizzazione (DA ADATTARE ALL'OUTPUT REALE DEL SID!)
		// Se buffer[i] è un valore a 16bit (0-65535)
		// sampleFloat32 := (float32(buffer[i]) / 32768.0) - 1.0
		// Se buffer[i] è un valore a 8bit (0-255) dal SID (improbabile per qualità)
		// sampleFloat32 := (float32(buffer[i]) / 128.0) - 1.0
		// Per ora, usiamo un placeholder. Devi implementare la normalizzazione corretta.
		// Il buffer in input è []uint32, ma il SID produce tipicamente valori più piccoli.
		// Assumiamo che i valori in `buffer` siano già pre-processati o che questa sia
		// una rappresentazione intermedia.
		// Per questo esempio, convertiamo semplicemente e scaliamo (NON CORRETTO PER IL SID REALE):
		sampleFloat32 := float32(int32(buffer[i])) / float32(0x7FFFFFFF) // Placeholder
		a.goBuffer = append(a.goBuffer, sampleFloat32)
	}

	a.pos += samples
	a.cyclesSinceFlush += samples // Assumiamo 1 campione = 1 ciclo "audio" per semplicità qui

	// Se abbiamo accumulato abbastanza campioni o cicli, inviali a JavaScript
	if len(a.goBuffer) >= cap(a.goBuffer) || a.cyclesSinceFlush >= a.bufferSizeCycles {
		if len(a.goBuffer) > 0 {
			// Crea o riutilizza il Float32Array JavaScript
			// È più efficiente creare una volta e riempire se possibile, ma CopyBytesToJS è più semplice
			// per iniziare e gestisce la creazione di un nuovo ArrayBuffer in JS.
			// Per performance massime, si dovrebbe scrivere direttamente nella memoria WASM
			// un Float32Array e passare il puntatore/lunghezza a JS, come per il video.
			// Ma la Web Audio API preferisce ricevere Float32Array direttamente.

			jsBuffer := js.Global().Get("Float32Array").New(len(a.goBuffer))
			// Converti goBuffer (slice di float32) in []byte per CopyBytesToJS
			// Questo richiede di "reinterpretare" la slice di float32 come una slice di byte.
			// Ogni float32 è 4 byte.
			header := (*[3]uintptr)(unsafe.Pointer(&a.goBuffer))
			byteSlice := (*[1 << 30]byte)(unsafe.Pointer(header[0]))[: len(a.goBuffer)*4 : len(a.goBuffer)*4]

			js.CopyBytesToJS(jsBuffer, byteSlice)

			a.onSamplesReady.Invoke(jsBuffer)

			a.goBuffer = make([]float32, 0, cap(a.goBuffer)) // Oppure a.goBuffer = a.goBuffer[:0]
		}
		a.cyclesSinceFlush = 0
	}
}

// Play resumes audio playback if the AudioContext is suspended. Resets or starts AudioBufferSourceNode if needed.
func (a *Audio) Play() {
	if a.audioCtx.Truthy() && a.audioCtx.Get("state").String() == "suspended" {
		a.audioCtx.Call("resume")
		println("Go Audio: Resumed")
	}
}

// Pause suspends audio playback if the underlying AudioContext is currently in the "running" state.
func (a *Audio) Pause() {
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

		// Calculate a reasonable buffer size, e.g., for ~20-40ms of audio
		// This depends on the emulator frequency (e.g., 1MHz)
		// and how many samples the SID produces per cycle.
		// For now, we use a fixed size, but it should be calculated better.
		// If the SID produces one sample every N cycles, and we want to buffer M milliseconds:
		// bufferLenSamples = (sampleRate * M_ms) / 1000 // a.bufferSizeCycles = (N_cycles_per_sample * bufferLenSamples)
		// For example, if SID produces one sample every ~20 C64 cycles (@1MHz),
		// and JS sampleRate is 48000Hz, for 20ms:
		// bufferLenSamples = (48000 * 20) / 1000 = 960 samples
		// a.bufferSizeCycles = 20 * 960 = 19200 C64 cycles
		// This means Write will be called many times before reaching bufferSizeCycles.

		// La dimensione di goBuffer deve essere scelta per accumulare abbastanza campioni
		// prima di inviarli a JS, per evitare troppe chiamate JS.
		// Ad es. 1024 o 2048 campioni.
		a.goBuffer = make([]float32, 0, 2048) // Buffer di accumulo in Go
		a.bufferSizeCycles = 20000            // Valore indicativo, da affinare

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

	js.Global().Set("initAudioContextAndGetCallback", js.FuncOf(init))
	js.Global().Set("wasmAudioPlay", js.FuncOf(play))
	js.Global().Set("wasmAudioPause", js.FuncOf(pause))

	return nil
}

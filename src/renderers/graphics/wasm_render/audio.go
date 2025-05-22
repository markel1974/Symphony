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
	cfg               *config.Config
	pos               int
	audioCtx          js.Value  // Reference to the JavaScript AudioContext
	sampleRate        float64   // Sample rate of the AudioContext
	goBuffer          []float32 // Buffer in Go to accumulate samples
	goBufferTargetLen int
	onSamplesReady    js.Value // JS callback function to send samples
	bufferSizeCycles  int      // Buffer size in emulation cycles (approximately)
	cyclesSinceFlush  int
	//jsTypedArray     js.Value  // Uint8Array or Float32Array in JS to pass data
}

// NewAudio creates and returns a new instance of the Audio structure with default initialization for playback management.
func NewAudio() *Audio {
	return &Audio{
		pos: 0,
	}
}

// Setup initializes the Audio instance with the given configuration and prepares it for further setup steps.
func (a *Audio) Setup(cfg *config.Config) error {
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
func (a *Audio) Write(buffer []uint32, _ int, samplesOutputCount int) {
	if a.audioCtx.IsUndefined() {
		fmt.Println("Write Error: AudioContext not initialized")
		a.pos += samplesOutputCount
		return
	}
	if a.onSamplesReady.IsUndefined() || !a.onSamplesReady.Truthy() {
		fmt.Println("Write Error: OnSamplesReady not initialized")
		a.pos += samplesOutputCount
		return
	}
	numUint32sToProcess := 0
	if samplesOutputCount > 0 {
		if (samplesOutputCount % 2) != 0 {
			// Questo sarebbe un errore di logica da parte del chiamante, samplesOutputCount dovrebbe essere pari
			fmt.Printf("Write Warning: samplesOutputCount (%d) is odd, implies half a uint32.\n", samplesOutputCount)
			// Decidi come gestire: troncare, arrotondare, errore? Per ora, tronchiamo.
			samplesOutputCount--
		}
		numUint32sToProcess = samplesOutputCount / 2
	}
	if len(buffer) < numUint32sToProcess {
		fmt.Printf("Write Error: Input buffer too small. Has %d uint32s, expected %d for %d output samples.\n", len(buffer), numUint32sToProcess, samplesOutputCount)
		// Potresti voler troncare numUint32sToProcess a len(buffer) o ritornare un errore.
		// Se si tronca, anche samplesOutputCount dovrebbe essere aggiustato:
		// numUint32sToProcess = len(buffer)
		// samplesOutputCount = numUint32sToProcess * 2
		// Per ora, usciamo se c'è un mismatch grave per evidenziare il problema.
		// In produzione, potresti voler gestire diversamente.
		return
	}

	// Processa ogni uint32 per estrarre due campioni float32
	for i := 0; i < numUint32sToProcess; i++ {
		inputUint32 := buffer[i]

		// QUI DEVI IMPLEMENTARE L'ESTRAZIONE E LA NORMALIZZAZIONE CORRETTA
		// Esempio: se inputUint32 contiene due campioni a 16 bit (firmati, little-endian)
		sample1_raw_int16 := int16(inputUint32 & 0xFFFF)
		sample2_raw_int16 := int16((inputUint32 >> 16) & 0xFFFF)

		// Normalizza a float32 (da -1.0 a 1.0)
		sample1_float32 := float32(sample1_raw_int16) / 32768.0
		sample2_float32 := float32(sample2_raw_int16) / 32768.0

		a.goBuffer = append(a.goBuffer, sample1_float32)
		a.goBuffer = append(a.goBuffer, sample2_float32)
	}
	a.pos += samplesOutputCount
	a.cyclesSinceFlush += samplesOutputCount // Assumiamo 1 campione = 1 ciclo "audio" per semplicità qui

	if len(a.goBuffer) >= a.goBufferTargetLen || a.cyclesSinceFlush >= a.bufferSizeCycles {
		if len(a.goBuffer) > 0 {
			// Converti a.goBuffer (slice di float32) in []byte
			// Ogni float32 è 4 byte.
			// Questa parte per ottenere byteSlice è corretta:
			header := (*[3]uintptr)(unsafe.Pointer(&a.goBuffer))
			byteSlice := (*[1 << 30]byte)(unsafe.Pointer(header[0]))[: len(a.goBuffer)*4 : len(a.goBuffer)*4]

			// 1. Crea un ArrayBuffer in JavaScript
			jsArrayBuffer := js.Global().Get("ArrayBuffer").New(len(byteSlice))

			// 2. Crea una vista Uint8Array di questo ArrayBuffer
			jsUint8Array := js.Global().Get("Uint8Array").New(jsArrayBuffer)

			// 3. Copia i byte da Go alla Uint8Array JavaScript
			// js.CopyBytesToJS restituisce il numero di byte copiati.
			copiedCount := js.CopyBytesToJS(jsUint8Array, byteSlice)
			if copiedCount != len(byteSlice) {
				fmt.Println("Attenzione: js.CopyBytesToJS non ha copiato tutti i byte.")
				// Potresti voler gestire questo caso in modo più robusto
			}

			// 4. Crea una vista Float32Array dello stesso ArrayBuffer
			jsFloat32Array := js.Global().Get("Float32Array").New(jsArrayBuffer)

			// 5. Invoca la callback JavaScript con il Float32Array
			a.onSamplesReady.Invoke(jsFloat32Array)

			// Resetta goBuffer (a.goBuffer[:0] è più efficiente se la capacità è già adeguata)
			a.goBuffer = a.goBuffer[:0]
		}
		a.cyclesSinceFlush = 0
	}
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
		const targetSamplesPerFlush = 3528
		a.goBufferTargetLen = targetSamplesPerFlush
		a.goBuffer = make([]float32, 0, targetSamplesPerFlush+512) // Capacità leggermente maggiore per sicurezza
		a.bufferSizeCycles = targetSamplesPerFlush                 // Valore indicativo, da affinare

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

//go:build js && wasm

package wasm_render

import (
	"fmt"
	"syscall/js"
	"unsafe"

	"github.com/markel1974/c64emu/src/config"
)

const (
	targetSamplesPerFlush = 3528
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

// Audio.Write modificato per gestire l'output attuale di calcBuffer
func (a *Audio) Write(bufferInput []uint32, writePosHint int, samplesOutputCountExpectedByCaller int) {
	// samplesOutputCountExpectedByCaller è il valore (es. 1764) che il chiamante si aspetta
	// in termini di avanzamento o dimensione logica del blocco.

	// Controlli iniziali (audioCtx, onSamplesReady) - mantienili come sono
	if a.audioCtx.IsUndefined() {
		fmt.Println("Go Audio.Write: Errore - AudioContext non inizializzato")
		// Come gestire a.pos qui? Se il chiamante si aspetta un avanzamento,
		// potresti dover usare samplesOutputCountExpectedByCaller.
		// Per ora, lo lasciamo così, ma è da considerare.
		a.pos += samplesOutputCountExpectedByCaller // O forse 0 se non si processa nulla?
		return
	}
	if a.onSamplesReady.IsUndefined() || !a.onSamplesReady.Truthy() {
		fmt.Println("Go Audio.Write: Errore - OnSamplesReady non inizializzato")
		a.pos += samplesOutputCountExpectedByCaller // Come sopra
		return
	}

	// 1. Determina quanti uint32 nel bufferInput sono stati effettivamente riempiti da calcBuffer
	var numUint32sConDatiValidi int
	if len(bufferInput) > 0 {
		// calcBuffer riempie (len(bufferInput_originale_passato_a_calcBuffer) / 2) + 1 elementi.
		// Assumiamo che bufferInput sia la slice originale passata a calcBuffer.
		numUint32sConDatiValidi = (len(bufferInput) / 2) + 1
		// Mettiamo un limite di sicurezza per non leggere oltre la fine di bufferInput,
		// anche se la formula sopra dovrebbe essere corretta se len(bufferInput) è la stessa
		// dimensione usata da calcBuffer come riferimento per len(buf)/2.
		if numUint32sConDatiValidi > len(bufferInput) {
			numUint32sConDatiValidi = len(bufferInput)
		}
	} else {
		numUint32sConDatiValidi = 0
	}

	// 2. Processa solo gli uint32 validi, estraendo un singolo campione mono da ciascuno
	numFloat32CampioniProdotti := 0
	for i := 0; i < numUint32sConDatiValidi; i++ {
		valUint32 := bufferInput[i]

		// Estrai il singolo campione mono a 16 bit (assumiamo sia nei 16 bit inferiori)
		// I 16 bit superiori vengono ignorati perché calcBuffer scrive un solo valore.
		campioneInt16 := int16(valUint32 & 0xFFFF)

		// Normalizza a float32
		campioneFloat32 := float32(campioneInt16) / 32768.0 // Normalizzazione per int16

		a.goBuffer = append(a.goBuffer, campioneFloat32)
		numFloat32CampioniProdotti++
	}

	// 3. Aggiorna la posizione e i cicli
	// a.pos: Questa è la parte più incerta.
	// Se a.pos deve riflettere l'avanzamento come inteso dal chiamante (AudioBuilder),
	// allora dovrebbe usare samplesOutputCountExpectedByCaller (es. 1764).
	// AudioBuilder aggiorna dr.sbPos con quel valore.
	// Se invece a.pos traccia i campioni effettivamente gestiti da Audio.Write,
	// allora dovrebbe usare numFloat32CampioniProdotti.
	// Per coerenza con come AudioBuilder sembra gestire la sua posizione (sbPos),
	// potresti voler usare samplesOutputCountExpectedByCaller.
	// Scegli una delle due o chiarisci il ruolo di a.pos.
	// Opzione 1: Usa l'aspettativa del chiamante
	a.pos += samplesOutputCountExpectedByCaller
	// Opzione 2: Usa i campioni effettivi (meno probabile se il chiamante gestisce la posizione esterna)
	// a.pos += numFloat32CampioniProdotti

	// a.cyclesSinceFlush deve tracciare i campioni effettivamente aggiunti a goBuffer
	a.cyclesSinceFlush += numFloat32CampioniProdotti

	// 4. Invia i dati a JavaScript quando goBuffer è pieno (questa logica rimane simile)
	if len(a.goBuffer) >= cap(a.goBuffer) || a.cyclesSinceFlush >= a.bufferSizeCycles {
		if len(a.goBuffer) > 0 {
			// ... (codice per creare jsArrayBuffer, jsUint8Array, js.CopyBytesToJS, jsFloat32Array)
			// Questo codice che hai già dovrebbe funzionare,
			// jsFloat32Array verrà creato con len(a.goBuffer) che ora riflette
			// il numero corretto di campioni float32 singoli.

			header := (*[3]uintptr)(unsafe.Pointer(&a.goBuffer))
			byteSlice := (*[1 << 30]byte)(unsafe.Pointer(header[0]))[: len(a.goBuffer)*4 : len(a.goBuffer)*4]

			jsArrayBuffer := js.Global().Get("ArrayBuffer").New(len(byteSlice))
			jsUint8Array := js.Global().Get("Uint8Array").New(jsArrayBuffer)

			js.CopyBytesToJS(jsUint8Array, byteSlice)

			jsFloat32Array := js.Global().Get("Float32Array").New(jsArrayBuffer)

			a.onSamplesReady.Invoke(jsFloat32Array)

			a.goBuffer = a.goBuffer[:0] // Resetta la lunghezza mantenendo la capacità
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

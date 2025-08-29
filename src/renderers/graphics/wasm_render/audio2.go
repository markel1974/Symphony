//go:build js && wasm

package wasm_render

/*
import (
	"fmt"
	"sync/atomic"
	"syscall/js"
	//"time"
	"unsafe"
)

// Definiamo la struttura del nostro buffer di controllo condiviso
// Questi sono i primi byte della SharedArrayBuffer
const (
	controlBlockSize = 16 // Dimensione totale in byte (4 * sizeof(uint32))
	dataBufferSize   = 8  // Esempio: capacità del ring buffer in numero di chunk
)

const (
	bufferStateEmpty   = 0 // La coda è completamente vuota.
	bufferStateTooFast = 1 // Il consumatore è troppo veloce, la coda si sta svuotando.
	bufferStateNice    = 2 // Stato di equilibrio ideale.
	bufferStateGood    = 3
	bufferStateStable  = 4 // Stato ancora stabile, agisce come "dead zone" per evitare oscillazioni.
)

var (
	// Buffer di memoria condiviso tra Go e JS
	_sharedBuffer js.Value
	// Viste a 32-bit per leggere/scrivere gli indici in modo atomico
	_sharedControl *[controlBlockSize / 4]uint32
	// Vista a float32 per i dati audio
	_sharedData *[]float32
)

// ... (tutto il codice per writeFn, Format, _formats rimane invariato) ...

// ContinuousReader ora è più semplice, la sua logica è quasi interamente
// contenuta nelle funzioni esportate.
type ContinuousReader struct {
	// ... (le sue dipendenze come interpolator, doubleBuffer, etc.)
	interpolator *LinearInterpolation
	chunkSize    int
	doubleBuffer *[]float32
	writeFn      writeFn
	ring         *CircularQueueWasm
	// ... (potrebbe non servire più un ring buffer separato in Go)
}

var reader *ContinuousReader // Istanza globale per semplicità

// initWasm viene chiamato da JS per inizializzare il sistema.
// Riceve la SharedArrayBuffer creata da JS e imposta le viste.
func initWasm(this js.Value, args []js.Value) interface{} {
	fmt.Println("Symphony WASM Audio: Initializing...")

	// Riceve i parametri da JavaScript
	_sharedBuffer = args[0]
	sampleRate := args[1].Int()
	chunkPerSecond := args[2].Int()

	// Crea l'istanza del nostro renderer
	reader = &ContinuousReader{}
	reader.chunkSize = sampleRate / chunkPerSecond
	reader.interpolator = NewLinearInterpolation(reader.chunkSize * 2)
	doubleBuffer := make([]float32, reader.chunkSize*2)
	reader.doubleBuffer = &doubleBuffer
	reader.writeFn = writeFloat32LE // Esempio

	// --- "Magia" Corretta e Sicura per mappare la SharedArrayBuffer in Go ---

	// Ottiene l'array di byte sottostante dalla SharedArrayBuffer
	jsBuffer := js.Globals().Get("Uint8Array").New(_sharedBuffer)

	// Crea uno slice Go che punta alla stessa memoria (usa unsafe in un solo punto per ottenere il puntatore iniziale)
	dataSlice := make([]byte, _sharedBuffer.Get("byteLength").Int())
	js.CopyBytesToGo(dataSlice, jsBuffer)

	// 1. Otteniamo un puntatore unsafe al primo byte del nostro slice. Questo è il nostro punto di ancoaraggio.
	basePtr := unsafe.Pointer(&dataSlice[0])

	// 2. Creiamo le viste direttamente da questo puntatore base, senza passi intermedi.
	//    Questa conversione è permessa e sicura secondo le regole di `unsafe`.
	_sharedControl = (*[controlBlockSize / 4]uint32)(basePtr)

	// 3. Per i dati, calcoliamo il nuovo puntatore usando aritmetica su uintptr
	//    ma la conversione finale avviene sempre partendo da un unsafe.Pointer.
	dataPtr := unsafe.Pointer(uintptr(basePtr) + uintptr(controlBlockSize))
	dataSamples := (len(dataSlice) - controlBlockSize) / 4
	sharedData := (*[1 << 30]float32)(dataPtr)[:dataSamples:dataSamples]
	_sharedData = &sharedData
	// --- Fine "Magia" Corretta ---

	fmt.Printf("Symphony WASM Audio: Ready. Shared buffer size: %d bytes\n", len(dataSlice))
	return nil
}

func (r *ContinuousReader) AddChunk(chunk *[]float32) {
	// Delega il lavoro alla nostra nuova coda atomica
	r.ring.Push(chunk)
}

// addChunk è la funzione che l'emulatore Go chiama per aggiungere dati.
// Ora scrive direttamente nella memoria condivisa.
func addChunk(this js.Value, args []js.Value) interface{} {
	// args[0] dovrebbe essere un Float32Array dal mondo Go
	goChunk := args[0]

	head := _sharedControl[headControlIndex]
	tail := _sharedControl[tailControlIndex]

	nextHead := (head + 1) % uint32(dataBufferSize)
	if nextHead == tail {
		// Buffer pieno, il produttore scarta il frame.
		// In un sistema reale, questo indica che il throttle deve rallentare.
		fmt.Println("WASM Audio: Ring buffer full, dropping frame!")
		return nil
	}

	// Copia i dati nel prossimo slot libero del ring buffer
	offset := int(head) * reader.chunkSize
	js.CopyBytesToJS(js.Globals().Get("Float32Array").New(_sharedBuffer, controlBlockSize+(offset*4), reader.chunkSize), goChunk)

	// Aggiorna l'indice di testa
	_sharedControl[headControlIndex] = nextHead

	return nil
}

// La funzione `read` esportata a JS rimane quasi uguale, ma ora legge gli indici
// in modo atomico dal blocco di controllo condiviso.
func read(this js.Value, args []js.Value) interface{} {
	outputJS := args[0]
	//bufLen := outputJS.Get("length").Int()

	head := atomic.LoadUint32(&_sharedControl[headControlIndex])
	tail := atomic.LoadUint32(&_sharedControl[tailControlIndex])

	counter := (head - tail + uint32(dataBufferSize)) % uint32(dataBufferSize)

	//counter := (head - tail + uint32(dataBufferSize)) % uint32(dataBufferSize)

	var samplesToPlay *[]float32

	switch counter {
	case bufferStateEmpty:
		// Restituiamo silenzio, ma non cambiamo nessun indice
		outputJS.Call("fill", 0)
		return nil

	case bufferStateTooFast:
		// "Stretching"
		offset := int(tail) * reader.chunkSize
		chunkToStretch := (*_sharedData)[offset : offset+reader.chunkSize]
		stretchedChunk, _ := reader.interpolator.LinearInterpolationF32(&chunkToStretch, reader.chunkSize*2)

		// Scriviamo la prima metà nell'output JS
		outputJS.Call("set", *stretchedChunk)

		// Scriviamo la seconda metà nello stesso slot del ring buffer, sovrascrivendolo
		js.CopyBytesToJS(js.Globals().Get("Float32Array").New(_sharedBuffer, controlBlockSize+(offset*4), reader.chunkSize), (*stretchedChunk)[reader.chunkSize:])

		// NON avanziamo il tail. Alla prossima chiamata, il consumatore rileggerà
		// lo stesso slot, che ora contiene la seconda metà "stirata".
		// Questo è il modo corretto di "creare" un chunk.

		return nil

	case bufferStateGood, bufferStateStable:
		// "Good"
		offset := int(tail) * reader.chunkSize
		d := (*_sharedData)[offset : offset+reader.chunkSize]
		samplesToPlay = &d
		_sharedControl[tailControlIndex] = (tail + 1) % uint32(dataBufferSize) // Avanza il cursore di lettura

	default: // counter > 4, "Too Slow"
		// "Squishing"
		offset1 := int(tail) * reader.chunkSize
		offset2 := int((tail+1)%uint32(dataBufferSize)) * reader.chunkSize

		// Copia i due chunk nel buffer di lavoro
		copy(*reader.doubleBuffer, (*_sharedData)[offset1:offset1+reader.chunkSize])
		copy((*reader.doubleBuffer)[reader.chunkSize:], (*_sharedData)[offset2:offset2+reader.chunkSize])

		squishedChunk, _ := reader.interpolator.LinearInterpolationF32(reader.doubleBuffer, reader.chunkSize)
		samplesToPlay = squishedChunk

		// Abbiamo consumato due chunk, quindi avanziamo il cursore di lettura di due posizioni
		_sharedControl[tailControlIndex] = (tail + 2) % uint32(dataBufferSize)
	}

	// Copia i campioni finali nel buffer di output di JS
	outputJS.Call("set", *samplesToPlay)
	return nil
}
*/

/*
// main registra le funzioni Go per renderle chiamabili da JavaScript.
func main() {
	c := make(chan struct{}, 0)
	js.Globals().Set("wasmAudioInit", js.FuncOf(initWasm))
	js.Globals().Set("wasmAudioAddChunk", js.FuncOf(addChunk))
	js.Globals().Set("wasmAudioRead", js.FuncOf(read))
	<-c // Tiene il programma Go in esecuzione
}
*/

/*
// main.js

// Verifica il supporto per SharedArrayBuffer
if (typeof SharedArrayBuffer === 'undefined') {
  alert('Your browser does not support SharedArrayBuffer. Please use a modern browser like Chrome or Firefox and ensure cross-origin isolation headers are set.');
}

async function run() {
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
  go.run(result.instance);

  const audioCtx = new AudioContext();
  // L'utente deve interagire con la pagina per avviare l'audio
  document.getElementById('startButton').addEventListener('click', () => {
    audioCtx.resume();
    console.log('AudioContext resumed.');
  });

  // Aggiungiamo il nostro processore audio custom
  await audioCtx.audioWorklet.addModule('audio-processor.js');

  // --- Creazione della Memoria Condivisa ---
  const chunkSize = Math.floor(audioCtx.sampleRate / 50); // es. 882
  const ringBufferChunks = 8; // Capacità del nostro ring buffer in chunk
  const controlSizeBytes = 16;
  const dataSizeBytes = ringBufferChunks * chunkSize * 4; // 4 byte per float32

  // Il buffer condiviso
  const sharedBuffer = new SharedArrayBuffer(controlSizeBytes + dataSizeBytes);

  // Inizializza il lato Go, passando il buffer e i parametri
  wasmAudioInit(sharedBuffer, audioCtx.sampleRate, 50);

  // Crea il nodo dell'AudioWorklet, passando il buffer condiviso
  const wasmNode = new AudioWorkletNode(audioCtx, 'wasm-audio-processor', {
    processorOptions: {
      sab: sharedBuffer,
      chunkSize: chunkSize
    }
  });

  wasmNode.connect(audioCtx.destination);

  // --- Funzione di Test per il Produttore (l'emulatore) ---
  // In un'applicazione reale, questa verrebbe chiamata dal loop di emulazione
  function produceSample() {
      // Creiamo un chunk di dati di esempio (es. una sinusoide)
      const chunk = new Float32Array(chunkSize);
      // ... (logica per riempire il chunk) ...
      // Chiamiamo la funzione Go per aggiungere il chunk
      wasmAudioAddChunk(chunk);
  }

  // Simuliamo il produttore che gira a 50Hz
  setInterval(produceSample, 20); // 1000ms / 50 = 20ms
}

run();
*/

/*
// audio-processor.js

const controlBlockSize = 16;
const headIndex = 0; // In uint32
const tailIndex = 1; // In uint32

class WasmAudioProcessor extends AudioWorkletProcessor {
  constructor(options) {
    super();
    this.sab = options.processorOptions.sab;
    this.chunkSize = options.processorOptions.chunkSize;

    // Vista a 32-bit per leggere gli indici del ring buffer
    this.controlView = new Uint32Array(this.sab, 0, controlBlockSize / 4);
  }

  process(inputs, outputs, parameters) {
    // Prendiamo il buffer di output che il browser ci chiede di riempire
    const outputChannel = outputs[0][0]; // Assumiamo mono

    // Chiamiamo la funzione 'read' del WASM, passandogli il buffer da riempire.
    // Tutta la logica di resampling e state machine avviene in Go.
    if (globalThis.wasmAudioRead) {
        globalThis.wasmAudioRead(outputChannel);
    }

    // Manteniamo il processore attivo
    return true;
  }
}

registerProcessor('wasm-audio-processor', WasmAudioProcessor);
*/

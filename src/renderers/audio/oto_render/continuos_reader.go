package oto_render

import (
	"encoding/binary"
	"math"
	"sync"
)

type ContinuousReader struct {
	lock             sync.Mutex
	lastChunk        []uint32
	currentPos       int
	lastChunkSamples int
}

func NewContinuousReader() *ContinuousReader {
	return &ContinuousReader{}
}

func (r *ContinuousReader) AddChunk(chunk []uint32, samples int) {
	chunkLen := len(chunk)
	if chunkLen == 0 {
		return
	}
	lastChunkSamples := samples / 2
	r.lock.Lock()
	defer r.lock.Unlock()
	if chunkLen != len(r.lastChunk) {
		r.lastChunk = make([]uint32, len(chunk))
	}
	copy(r.lastChunk, chunk)
	r.currentPos = 0
	r.lastChunkSamples = lastChunkSamples
}

func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.lastChunk == nil || r.currentPos >= len(r.lastChunk) {
		return 0, nil
	}
	written := 0
	for x := 0; x < r.lastChunkSamples; x++ {
		start := written
		end := start + 4
		if end > len(buf) {
			break
		}
		curr := r.lastChunk[x]
		const divisor = float32(1 << 20)
		campioneFloat32 := float32(int32(curr)) / divisor
		binary.LittleEndian.PutUint32(buf[start:end], math.Float32bits(campioneFloat32))
		written += 4
		r.currentPos++
		/*
			start := written
			end := start + 4
			if end > len(buf) {
				break
			}
			val := float32(r.lastChunk[x]) / 32768.0
			binary.LittleEndian.PutUint32(buf[start:end], math.Float32bits(val))
			written += 4
			r.currentPos++

		*/
		/*
			start := written
			end := start + 4
			if end > len(buf) {
				break
			}
			// 2. Normalizza questo int32 a float32.
			// Il divisore deve essere il valore massimo teorico che i tuoi int32 raggiungono dopo lo shift >> 10.
			// Questo è 2^(31-10) = 2^21.
			const normalizationDivisor = float32(1 << 21) // Valore: 2097152.0
			campioneFloat32 := float32(r.lastChunk[x]) / normalizationDivisor

			// 3. **** (OPZIONALE) AMPLIFICA PER IL VOLUME ****
			// Questo è il punto per aumentare il volume se i tuoi campioni originali sono deboli.
			// Sperimenta con questo valore (es. 10.0, 50.0, 100.0, 200.0)
			//const amplificationFactor = float32(100.0) // Esempio: fattore di amplificazione
			//campioneFloat32 *= amplificationFactor

			// 4. (OPZIONALE) Clampa il valore per evitare clipping (distorsione)
			//if campioneFloat32 > 1.0 {
			//	campioneFloat32 = 1.0
			//} else if campioneFloat32 < -1.0 {
			//	campioneFloat32 = -1.0
			//}

			// 5. Converti il float32 normalizzato nei suoi bit uint32 per oto (FormatFloat32LE)
			binary.LittleEndian.PutUint32(buf[start:end], math.Float32bits(campioneFloat32))
			written += 4
			r.currentPos++

		*/
		/*
			start := written
			end := start + 4
			if end > len(buf) {
				break
			}
			curr := r.lastChunk[x]
			campioneInt16 := int16(curr & 0xFFFF)
			campioneFloat32 := float32(campioneInt16) / 32768.0
			val := math.Float32bits(campioneFloat32)
			binary.LittleEndian.PutUint32(buf[start:end], val)
			written += 4
			r.currentPos++

		*/
	}
	return written, nil
}

/*
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.lastChunk == nil || r.currentPos >= len(r.lastChunk) {
		return 0, nil
	}

	bytesWritten := 0
	numFloatsToCopy := 0

	remainingBufBytes := len(buf) - bytesWritten
	numFloatsCanWrite := remainingBufBytes / 4
	numFloatsInChunk := len(r.lastChunk) - r.currentPos

	numFloatsToCopy = int(math.Min(float64(numFloatsCanWrite), float64(numFloatsInChunk)))

	if numFloatsToCopy == 0 {
		return 0, nil
	}

	for i := 0; i < numFloatsToCopy; i++ {
		uint32bits := r.lastChunk[r.currentPos]
		campioneInt16 := int16(uint32bits & 0xFFFF)
		campioneFloat32 := float32(campioneInt16) / 32768.0
		start := bytesWritten + (i * 4)
		end := start + 4
		val := math.Float32bits(campioneFloat32)
		binary.LittleEndian.PutUint32(buf[start:end], val)
		r.currentPos++
	}
	bytesWritten += numFloatsToCopy * 4

	// Restituisci il numero di byte copiati.
	return bytesWritten, nil
}

*/

/*
import (
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"
)

const maxChunks = 1024 // Capacità massima della coda di chunk (es. 1024 chunk)

// ContinuousReader è un io.Reader che gestisce una coda di chunk audio in modo circolare.
type ContinuousReader struct {
	chunks [maxChunks][]uint32 // Array preallocato per i chunk
	head   int                 // Indice del prossimo chunk da leggere (consumare)
	tail   int                 // Indice del prossimo slot libero per scrivere (produrre)
	count  int                 // Numero di chunk attualmente presenti nella coda

	currentChunk []uint32 // Il chunk che stiamo leggendo attualmente (il pezzo del lotto combinato per oto)
	currentPos   int      // Posizione corrente nel currentChunk

	lock sync.Mutex // Protegge head, tail, count

	// Canale per segnalare a Read che ci sono nuovi dati disponibili.
	dataAvailable chan struct{}
}

// NewContinuousReader crea e inizializza un nuovo ContinuousReader.
func NewContinuousReader() *ContinuousReader {
	r := &ContinuousReader{
		head:          0,
		tail:          0,
		count:         0,
		dataAvailable: make(chan struct{}, 1), // Buffer 1 per non bloccare AddChunk se Read non è pronto a ricevere segnale
	}
	// Pre-allocare i singoli slice di chunk se la dimensione è fissa.
	// Usiamo una dimensione tipica (es. 1764 campioni) basata sui tuoi log.
	exampleChunkSize := 1764
	for i := 0; i < maxChunks; i++ {
		r.chunks[i] = make([]uint32, exampleChunkSize)
	}
	return r
}

// AddChunk aggiunge un nuovo chunk alla coda circolare.
func (r *ContinuousReader) AddChunk(chunk []uint32) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.count == maxChunks {
		log.Printf("ContinuousReader: Chunk buffer full (%d/%d), overwriting oldest chunk.", r.count, maxChunks)
		r.head = (r.head + 1) % maxChunks
	} else {
		r.count++
	}

	copy(r.chunks[r.tail], chunk)
	r.tail = (r.tail + 1) % maxChunks

	select {
	case r.dataAvailable <- struct{}{}:
	default:
		// Il segnale è già in coda, o Reader non sta aspettando.
	}
}

// Read legge i byte dal reader, prendendo i chunk dalla coda interna.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	// Questo è il numero massimo di byte che Read dovrebbe provare a restituire.
	// Non dobbiamo necessariamente riempire tutto 'buf', ma cerchiamo di prelevare
	// un intero chunk del modello, se disponibile.
	//maxBytesToRead := len(r.chunks[0]) * 4 // Dimensione in byte di un singolo chunk del modello (1764*4)

	bytesWritten := 0

	r.lock.Lock() // Proteggi l'accesso alla coda

	// Se il chunk corrente è esaurito O se non abbiamo un chunk valido
	// E non ci sono chunk disponibili nella coda (r.count == 0)
	if (r.currentChunk == nil || r.currentPos >= len(r.currentChunk)) && r.count == 0 {
		r.lock.Unlock() // Rilascia il lock prima di bloccare

		log.Println("NON CI SONO PIU CHUNKS.... WAITING")

		// Attendi attivamente che ci siano dati disponibili
		select {
		case <-r.dataAvailable: // Aspetta il segnale da AddChunk
			log.Println("CHUNKS ARRIVED.... RETRYING")
		case <-time.After(5 * time.Second): // Timeout
			log.Println("Reader: Timeout waiting for new chunks. Returning 0 bytes.")
			return 0, nil // Restituisci 0 byte se non ci sono dati dopo il timeout.
		}

		// Dopo aver aspettato, riprendi il loop per prendere un chunk.
		return 0, nil // Riprova la Read la prossima volta, dato che ora non abbiamo ancora un chunk.
	}

	// Se siamo qui, ci sono chunk disponibili nella coda (r.count > 0).
	// Se il chunk corrente è esaurito, prendine uno nuovo dalla coda.
	if r.currentChunk == nil || r.currentPos >= len(r.currentChunk) {
		r.currentChunk = r.chunks[r.head] // Prendi il chunk dalla head
		r.head = (r.head + 1) % maxChunks // Avanza la head in modo circolare
		r.count--                         // Decrementa il contatore di elementi
		r.currentPos = 0                  // Resetta la posizione all'inizio del nuovo chunk
		// log.Printf("READING CHUNK FROM QUEUE: Head: %d, Count: %d", r.head, r.count)
	}
	r.lock.Unlock() // Rilascia il lock mentre copiamo i dati

	// Ora copia i dati dal currentChunk al 'buf' di oto.
	// Copia al massimo la dimensione di un singolo chunk del modello, o meno se 'buf' è più piccolo.
	numFloatsToCopy := len(r.currentChunk) - r.currentPos // Quanti float32 rimangono nel currentChunk

	// Non copiare più di quanto 'buf' possa contenere (anche se il tuo chunk è più grande)
	if numFloatsToCopy*4 > len(buf) {
		numFloatsToCopy = len(buf) / 4
	}

	for i := 0; i < numFloatsToCopy; i++ {
		//uint32bits := r.currentChunk[r.currentPos]
		//float32Val := math.Float32frombits(uint32bits)
		//start := bytesWritten + (i * 4)
		//end := start + 4
		//val := math.Float32bits(float32Val)
		//binary.LittleEndian.PutUint32(buf[start:end], val)
		uint32bits := r.currentChunk[r.currentPos]
		campioneInt16 := int16(uint32bits & 0xFFFF)
		campioneFloat32 := float32(campioneInt16) / 32768.0
		start := bytesWritten + (i * 4)
		end := start + 4
		val := math.Float32bits(campioneFloat32)
		binary.LittleEndian.PutUint32(buf[start:end], val)
		r.currentPos++
	}
	bytesWritten += numFloatsToCopy * 4

	//log.Println("END READ", buf[:32])
	// Restituisci i byte scritti, che rappresentano una parte o tutto il singolo chunk prelevato.
	return bytesWritten, nil
}

*/

/*
import (
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"
)

const maxChunks = 8192 // Capacità massima della coda di chunk

// ContinuousReader è un io.Reader che gestisce una coda di chunk audio in modo circolare.
type ContinuousReader struct {
	chunks [maxChunks][]uint32 // Array preallocato per i chunk
	head   int                 // Indice del prossimo chunk da leggere (consumare)
	tail   int                 // Indice del prossimo slot libero per scrivere (produrre)
	count  int                 // Numero di chunk attualmente presenti nella coda

	currentChunk []uint32 // Il chunk che stiamo leggendo attualmente
	currentPos   int      // Posizione corrente nel currentChunk (all'interno di currentChunk)

	lock sync.Mutex // Protegge head, tail, count
}

// NewContinuousReader crea e inizializza un nuovo ContinuousReader.
// Inizializza l'array di chunk con una dimensione tipica, se fissa.
func NewContinuousReader() *ContinuousReader {
	r := &ContinuousReader{
		head:  0,
		tail:  0,
		count: 0,
	}
	// Pre-allocare i singoli slice di chunk se la dimensione è fissa.
	// Usiamo una dimensione tipica (es. 1764) basata sui tuoi log.
	exampleChunkSize := 1764
	for i := 0; i < maxChunks; i++ {
		r.chunks[i] = make([]uint32, exampleChunkSize)
	}
	return r
}

// AddChunk aggiunge un nuovo chunk alla coda circolare.
// Implementa la logica di "perdita dati" (sovrascrittura) se il buffer è pieno.
func (r *ContinuousReader) AddChunk(chunk []uint32) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Se la dimensione del chunk è cambiata, dobbiamo gestire la riallocazione,
	// ma per semplicità assumiamo che sia costante per ora (1764).

	// Se il buffer è pieno, sovrascrivi il chunk più vecchio (all'head).
	// Questo è il tuo requisito di "perdere il contenuto della coda".
	if r.count == maxChunks {
		log.Println("ContinuousReader: Chunk buffer full, overwriting oldest chunk.")
		// In questo caso, head avanza e count rimane lo stesso.
		r.head = (r.head + 1) % maxChunks
		// Non decrementiamo count, perché un nuovo chunk lo rimpiazza.
	} else {
		// Se non è pieno, incrementa il contatore.
		r.count++
	}

	// Copia il chunk nel tail (punto di scrittura)
	copy(r.chunks[r.tail], chunk)

	// Avanza il tail in modo circolare
	r.tail = (r.tail + 1) % maxChunks

	//log.Println("ADDING CHUNKS", chunk[:32]) // Metti i tuoi log qui se vuoi.
}

// Read legge i byte dal reader, prendendo i chunk dalla coda circolare.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	bytesWritten := 0
	counter := 0
	bufLen := len(buf)
	//log.Printf("STARTING READ")
	for bytesWritten < bufLen {
		r.lock.Lock()
		// Se il currentChunk è esaurito o non è stato ancora impostato
		if r.currentChunk == nil || r.currentPos >= len(r.currentChunk) {
			// Prendi un nuovo chunk dalla coda solo se disponibile
			if r.count > 0 { // Ci sono elementi nella coda
				r.currentChunk = r.chunks[r.head] // Prendi il chunk dalla head
				r.head = (r.head + 1) % maxChunks // Avanza la head in modo circolare
				r.count--                         // Decrementa il contatore di elementi
				r.currentPos = 0                  // Resetta la posizione nel nuovo chunk
				//log.Println("READING CHUNK FROM QUEUE", r.currentChunk[:32]) // Metti i tuoi log qui se vuoi
			} else {
				// Non ci sono chunk disponibili, sblocca e aspetta nuovi dati.
				r.lock.Unlock() // IMPORTANTE: Sblocca prima dell'attesa

				log.Println("NON CI SONO PIU CHUNKS.... WAITING")
				// Qui non c'è una vera attesa di canale.
				// oto si aspetta che Read non blocchi troppo a lungo o che non sia vuoto all'avvio.
				// Se non ci sono dati, oto si fermerà, come abbiamo visto.
				// Per il debug, puoi aggiungere una piccola pausa per non spammare la CPU
				time.Sleep(20 * time.Millisecond) // Piccola pausa

				// Se torniamo qui e non ci sono ancora dati, Read restituirà 0 e oto si ferma.
				return bytesWritten, nil // Restituisci quello che hai, senza errore.
			}
		}
		// Se siamo qui, abbiamo un currentChunk valido o ne abbiamo appena preso uno.
		r.lock.Unlock() // Sblocca mentre elaboriamo il chunk

		remainingBufBytes := len(buf) - bytesWritten
		numFloatsCanWrite := remainingBufBytes / 4
		numFloatsInChunk := len(r.currentChunk) - r.currentPos

		numFloatsToCopy := int(math.Min(float64(numFloatsCanWrite), float64(numFloatsInChunk)))

		if numFloatsToCopy == 0 {
			// Questo accade se remainingBufBytes < 4 o il currentChunk è esaurito.
			// La prossima iterazione del loop 'for bytesWritten < len(buf)' proverà a prendere un nuovo chunk.
			continue
		}

		for i := 0; i < numFloatsToCopy; i++ {
			uint32bits := r.currentChunk[r.currentPos]
			float32Val := math.Float32frombits(uint32bits)
			binary.LittleEndian.PutUint32(buf[bytesWritten+i*4:bytesWritten+i*4+4], math.Float32bits(float32Val))
			r.currentPos++
		}
		bytesWritten += numFloatsToCopy * 4
		counter++

		//log.Println("END READ", "REQUIRED CHUNK", counter, bytesWritten, buf[:32])
		r.currentChunk = nil
		r.currentPos = 0
		return bytesWritten, nil

	}
	log.Println("END READ", "REQUIRED CHUNK", counter, buf[:32])
	return bytesWritten, nil
}

*/

/*
// AddChunk aggiunge un nuovo chunk alla coda usando il canale.
// Questa operazione è intrinsecamente thread-safe per i canali.
func (r *ContinuousReader) AddChunk(chunk []uint32) {
	cloned := make([]uint32, len(chunk))
	copy(cloned, chunk)


	// Non stampare qui per i log, sposta il debug all'origine del modello
	// fmt.Println("ADDING CHUNKS", cloned[:32])

	// Manda il chunk sul canale. Si bloccherà se il canale è pieno.
	select {
	case r.chunkChan <- cloned:
		// Chunk aggiunto con successo
	default:
		// Canale pieno, non possiamo aggiungere il chunk per ora.
		// Questo significa che stiamo producendo chunk più velocemente di quanto oto li consumi.
		// È un segnale di "buffer over-run", e si scarta il chunk per mantenere il ritmo.
		log.Println("ContinuousReader: Chunk buffer full, dropping chunk.")
	}
}

// Read legge i byte dal reader, prendendo i chunk dal canale.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	bytesWritten := 0
	for bytesWritten < len(buf) {
		// Se il currentChunk è esaurito o non è stato ancora impostato, prendine uno nuovo.
		if r.currentChunk == nil || r.currentPos >= len(r.currentChunk) {
			// Tenta di leggere un nuovo chunk dal canale
			// Questo è il punto in cui il reader si blocca se non ci sono chunk
			select {
			case newChunk := <-r.chunkChan:
				// Chunk ricevuto, lo usiamo
				r.currentChunk = newChunk
				r.currentPos = 0
				// log.Println("READING CHUNK FROM QUEUE", r.currentChunk[:32]) // Se vuoi i log di lettura
			case <-time.After(5 * time.Second):
				// Timeout: non ci sono chunk per 5 secondi.
				// Restituiamo 0 byte e nessun errore, ma segnaliamo.
				log.Println("NON CI SONO PIU CHUNKS.... TIMEOUT WAITING")
				// In questo caso, oto potrebbe fermarsi se non riceve dati.
				// Puoi provare a restituire io.EOF qui se vuoi una chiusura controllata.
				return bytesWritten, nil // Restituisci quello che hai, senza errore.
			}
		}

		// Se siamo qui, abbiamo un currentChunk valido o ne abbiamo appena preso uno.
		// La logica di copia rimane la stessa, ma senza mutex, perché il canale ha garantito la sicurezza.
		remainingBufBytes := len(buf) - bytesWritten
		numFloatsCanWrite := remainingBufBytes / 4
		numFloatsInChunk := len(r.currentChunk) - r.currentPos

		numFloatsToCopy := int(math.Min(float64(numFloatsCanWrite), float64(numFloatsInChunk)))

		if numFloatsToCopy == 0 {
			// Questo accade se remainingBufBytes < 4 o il currentChunk è esaurito.
			// Il ciclo si ripeterà per prendere un nuovo chunk dal canale.
			continue
		}

		for i := 0; i < numFloatsToCopy; i++ {
			uint32bits := r.currentChunk[r.currentPos]
			float32Val := math.Float32frombits(uint32bits)
			binary.LittleEndian.PutUint32(buf[bytesWritten+i*4:bytesWritten+i*4+4], math.Float32bits(float32Val))
			r.currentPos++
		}
		bytesWritten += numFloatsToCopy * 4
	}
	return bytesWritten, nil
}


*/

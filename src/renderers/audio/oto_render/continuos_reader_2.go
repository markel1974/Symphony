package oto_render

import (
	"fmt"
	"github.com/hajimehoshi/oto/v2"
	"sync"
	"time"
)

// Definiamo delle costanti per gli stati del buffer, per rendere il codice più leggibile.
const (
	bufferStateEmpty   = 0 // La coda è completamente vuota.
	bufferStateTooFast = 1 // Il consumatore è troppo veloce, la coda si sta svuotando.
	bufferStateGood    = 2 // Stato di equilibrio ideale.
	bufferStateStable  = 3 // Stato ancora stabile, agisce come "dead zone" per evitare oscillazioni.
)

// ContinuousReader2 gestisce un flusso audio continuo verso un player 'oto'.
// Utilizza una coda circolare bufferizzata e un algoritmo di resampling a stati
// per sincronizzare dinamicamente un produttore di dati "free to run" con un
// consumatore a tempo reale (la scheda audio), garantendo un flusso audio
// stabile e di alta qualità.
type ContinuousReader2 struct {
	player oto.Player
	//bytesPerSample int
	writeFn      writeFn
	ring         *CircularQueue2      // La nostra coda circolare di chunk audio.
	interpolator *LinearInterpolation // Per il resampling (stretch/squish).
	lock         sync.Mutex
	chunkSize    int        // La dimensione standard di un chunk in campioni.
	doubleBuffer *[]float32 // Un buffer di lavoro per unire due chunk.
}

// NewContinuousReader2 crea una nuova istanza di ContinuousReader2.
func NewContinuousReader2() *ContinuousReader2 {
	return &ContinuousReader2{}
}

// Setup inizializza il ContinuousReader con i parametri audio specificati.
// Configura il contesto audio 'oto', la coda circolare, l'interpolatore
// e tutti i buffer di lavoro necessari.
func (r *ContinuousReader2) Setup(sampleRate int, chunkPerSecond int, channels int, fo string) error {
	format, ok := _formats[fo]
	if !ok {
		return fmt.Errorf("audio format not found")
	}

	r.chunkSize = sampleRate / chunkPerSecond

	options := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channels,
		Format:       format.Format,
		BufferSize:   time.Second / time.Duration(chunkPerSecond),
	}
	ctx, ready, err := oto.NewContextWithOptions(options)
	if err != nil {
		return fmt.Errorf("failed to create oto context: %w", err)
	}
	<-ready

	r.ring = NewCircularQueue2(chunkPerSecond, r.chunkSize)
	// L'interpolatore e il doubleBuffer devono poter contenere due chunk per la logica "TOO SLOW".
	r.interpolator = NewLinearInterpolation(r.chunkSize * 2)
	doubleBuffer := make([]float32, r.chunkSize*2)
	r.doubleBuffer = &doubleBuffer

	//r.bytesPerSample = format.Bytes
	r.writeFn = format.Func
	r.player = ctx.NewPlayer(r)
	r.player.SetVolume(1.0)

	return nil
}

// Play avvia la riproduzione audio.
func (r *ContinuousReader2) Play() {
	r.player.Play()
}

// Err restituisce l'ultimo errore incontrato dal player.
func (r *ContinuousReader2) Err() error {
	return r.player.Err()
}

// AddChunk aggiunge un nuovo chunk di dati audio alla coda, in modo thread-safe.
func (r *ContinuousReader2) AddChunk(chunk *[]float32, samples int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.ring.Push(chunk)
}

// Read è il cuore del renderer. Viene chiamato da 'oto' quando ha bisogno di dati.
// Implementa una macchina a stati per decidere come processare i dati audio
// in base al livello di riempimento del buffer.
func (r *ContinuousReader2) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	switch r.ring.Counter() {
	case bufferStateEmpty:
		return r.handleEmpty(buf)
	case bufferStateTooFast:
		return r.handleTooFast(buf)
	case bufferStateGood, bufferStateStable:
		return r.handleGood(buf)
	default:
		return r.handleTooSlow(buf)
	}
}

// handleEmpty viene chiamato quando il buffer è vuoto. Riempie il buffer di 'oto'
// con silenzio per mantenere lo stream audio attivo e prevenire interruzioni.
func (r *ContinuousReader2) handleEmpty(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}

// handleGood viene chiamato quando il buffer è in uno stato di equilibrio (2 o 3 chunk).
// Preleva un chunk e lo riproduce così com'è, senza resampling.
func (r *ContinuousReader2) handleGood(buf []byte) (int, error) {
	chunkToPlay, _ := r.ring.Pop()
	written := 0
	for _, sample := range *chunkToPlay {
		written += r.writeFn(buf, sample, written)
	}
	return written, nil
}

// handleTooFast viene chiamato quando il buffer si sta svuotando (1 chunk rimasto).
// Significa che il consumatore (oto) è più veloce del produttore.
// Per "creare tempo", prende un chunk, lo "stira" al doppio della sua lunghezza,
// ne suona la prima metà e rimette la seconda metà in coda.
func (r *ContinuousReader2) handleTooFast(buf []byte) (int, error) {
	fmt.Printf("[%s] CONSUMER TOO FAST: Stretching...\n", time.Now().Format(time.RFC3339Nano))

	fmt.Println(time.Now().Format(time.RFC3339Nano), "CONSUMER TOO FAST: Stretching...")
	chunkToStretch, _ := r.ring.Pop()
	stretchedChunk, _ := r.interpolator.LinearInterpolationF32(chunkToStretch, r.chunkSize*2)
	firstHalf := (*stretchedChunk)[:r.chunkSize]
	secondHalf := (*stretchedChunk)[r.chunkSize:]
	r.ring.Push(&secondHalf)
	written := 0
	for _, sample := range firstHalf {
		written += r.writeFn(buf, sample, written)
	}
	return written, nil
}

// handleTooSlow viene chiamato quando il buffer si sta riempiendo (più di 3 chunk).
// Significa che il produttore (emulatore) è più veloce del consumatore.
// Per "recuperare il ritardo", preleva due chunk, li "comprime" in uno solo
// e lo suona, consumando dati dal buffer al doppio della velocità.
func (r *ContinuousReader2) handleTooSlow(buf []byte) (int, error) {
	fmt.Printf("[%s] PRODUCER TOO SLOW (Lag > %d chunks): Squishing...\n", time.Now().Format(time.RFC3339Nano), r.ring.Counter())
	chunk1, _ := r.ring.Pop()
	chunk2, _ := r.ring.Pop()
	copy(*r.doubleBuffer, *chunk1)
	copy((*r.doubleBuffer)[len(*chunk1):], *chunk2)
	squishedChunk, squishedLen := r.interpolator.LinearInterpolationF32(r.doubleBuffer, r.chunkSize)
	written := 0
	for i := 0; i < squishedLen; i++ {
		written += r.writeFn(buf, (*squishedChunk)[i], written)
	}
	return written, nil
}

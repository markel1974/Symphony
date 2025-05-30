package oto_render

import (
	"github.com/markel1974/c64emu/src/config"
	"log"
	"time"
)

type Audio struct {
	audioSampleRate    int
	audioChannelCount  int
	audioNextStartTime time.Time
	audioReader        *ContinuousReader
	cfg                *config.Config
	pos                int
	debug              bool
}

func NewAudio() *Audio {
	return &Audio{
		audioSampleRate:   44100,
		audioChannelCount: 1,
		pos:               0,
		debug:             false,
	}
}

func (d *Audio) Setup(cfg *config.Config) error {
	//StartStub()
	d.cfg = cfg
	reader, err := NewContinuousReader(d.audioSampleRate, d.audioChannelCount)
	if err != nil {
		return err
	}
	d.audioReader = reader

	go func() {
		log.Println("Starting continuous audio player...")
		for {
			d.audioReader.Play()
			if playerErr := d.audioReader.Err(); playerErr != nil {
				log.Println("Error playing audio:", playerErr)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		log.Println("Continuous audio player stopped.", err)
	}()
	return nil
}

func (d *Audio) GetCurrentPosition() int {
	return d.pos
}

// SamplesReadyCallback ora invia i nuovi chunk al reader del player esistente.
func (d *Audio) Write(values []uint32, pos int, samples int) {
	d.pos += samples
	d.audioReader.AddChunk(values, samples)

	if d.debug {
		currentTime := time.Now()
		// bufferDuration è già calcolato correttamente ora
		durationMs := (float64(len(values)) / float64(d.audioSampleRate)) * 1000.0
		bufferDuration := time.Duration(durationMs) * time.Millisecond

		if d.audioNextStartTime.IsZero() {
			// Inizializza solo la prima volta, o dopo un reset completo
			d.audioNextStartTime = currentTime //.Add(10 * time.Millisecond) // Piccola latenza iniziale se vuoi
		} else {
			// Se siamo MOLTO in ritardo (es. > 100ms), resettiamo il tempo di partenza.
			// Questo allinea la riproduzione al "NOW" per evitare un lag infinito.
			const lagThreshold = 20 * time.Millisecond // Puoi sperimentare con questo valore
			if d.audioNextStartTime.Before(currentTime.Add(-lagThreshold)) {
				log.Printf("Go: Major lag detected. Resetting start time. Ideal: %s, Current: %s. Data still being buffered.\n",
					d.audioNextStartTime.Format("15:04:05.000"), currentTime.Format("15:04:05.000"))
				// Resetta audioNextStartTime al tempo attuale (o leggermente in avanti)
				d.audioNextStartTime = currentTime //.Add(10 * time.Millisecond)
			}
		}
		// Fa avanzare il tempo ideale per il prossimo buffer
		d.audioNextStartTime = d.audioNextStartTime.Add(bufferDuration)
		//log.Printf("Debug: audioNextStartTime before add: %s, bufferDuration: %v, after add: %s",
		//	d.audioNextStartTime.Add(-bufferDuration).Format("15:04:05.000.000"), bufferDuration, d.audioNextStartTime.Format("15:04:05.000.000"))
	}
}

func (d *Audio) Play() {
}

func (d *Audio) Pause() {
}

func (d *Audio) Resume() {
}

/*
func main() {
	if err := InitializeAudio(); err != nil {
		log.Fatalf("Error initializing audio: %v", err)
	}

	// Esempio d'uso: Simula la ricezione di chunk audio come []uint32
	// Ora inviamo i chunk a intervalli regolari, e il singolo player li riprodurrà.
	for i := 0; i < 20; i++ { // Simuliamo più chunk per vedere il flusso continuo
		// Crea un buffer dummy (ad es. 1 secondo di onda sinusoidale)
		dummyUint32Buffer := make([]uint32, audioSampleRate/2) // 0.5 secondi di audio per chunk
		for j := 0; j < len(dummyUint32Buffer); j++ {
			// Crea una sinusoide semplice per il test
			val := float32(math.Sin(float64(j)/float64(audioSampleRate)*2*math.Pi*440)) * 0.2 // Volume ridotto
			dummyUint32Buffer[j] = math.Float32bits(val)
		}
		samplesReadyCallback(dummyUint32Buffer)
		time.Sleep(450 * time.Millisecond) // Invia un nuovo chunk leggermente prima che il precedente finisca
	}

	// Lascia il main goroutine attivo per consentire la riproduzione audio
	log.Println("Finished sending chunks. Player should continue playing buffered audio.")
	time.Sleep(5 * time.Second) // Attendi ancora un po' per l'audio residuo
	fmt.Println("Example finished.")
}
*/

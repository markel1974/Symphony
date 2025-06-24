package oto_render

import (
	"github.com/markel1974/c64emu/src/config"
	"log"
	"time"
)

// Audio represents a structure for managing audio playback and configuration in a continuous stream setup.
// It includes options for sample rate, channel count, playback control, debugging, and current audio position tracking.
// The type integrates with ContinuousReader for audio streaming and playback functionality.
type Audio struct {
	audioSampleRate    int
	audioChannelCount  int
	audioChunks        int
	audioNextStartTime time.Time
	audioReader        *ContinuousReader
	cfg                *config.Config
	debug              bool
}

// NewAudio creates and returns a new instance of the Audio type with default configurations for sample rate and channel count.
func NewAudio() *Audio {
	//TODO DRIVER CONFIG!
	const sampleRate = 44100
	const channels = 1
	const chunks = 50
	a := &Audio{
		audioSampleRate:   sampleRate,
		audioChannelCount: channels,
		audioChunks:       chunks,
		debug:             false,
	}

	return a
}

// Setup initializes the Audio instance with the specified configuration and sets up continuous audio playback.
func (d *Audio) Setup(cfg *config.Config) error {
	//StartStub()
	d.cfg = cfg
	reader := NewContinuousReader()
	if err := reader.Setup(d.audioSampleRate, d.audioChunks, d.audioChannelCount, "FLOAT32LE"); err != nil {
		return err
	}
	d.audioReader = reader
	go func() {
		d.audioReader.Play()
		if err := d.audioReader.Err(); err != nil {
			log.Println("Error playing audio:", err)
		}
	}()
	return nil
}

// Write processes and buffers audio samples for playback, updating the current position and managing playback timing.
func (d *Audio) Write(values *[]float32, samples int) {
	d.audioReader.AddChunk(values, samples)

	if d.debug {
		currentTime := time.Now()
		// bufferDuration è già calcolato correttamente ora
		durationMs := (float64(len(*values)) / float64(d.audioSampleRate)) * 1000.0
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

// Play starts or resumes the audio playback for the Audio instance, leveraging the associated continuous reader.
func (d *Audio) Play() {
}

// Pause halts the current audio playback, maintaining the current playback position.
func (d *Audio) Pause() {
}

// Resume resumes audio playback by restarting the associated audio reader or player.
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

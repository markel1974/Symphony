package board

import (
	"fmt"
	"time"
)

/*
Dati tecnici originali:
Tempo di seek traccia-traccia: 3 ms (documentazione tecnica del 1541)
Modalità movimento: Full-step (1 traccia per step, non half-step)
Velocità massima: ~333 tracce/secondo (1 / 0.003)
*/

const (
	rpm          = 300
	stepsPerSec  = 1_000_000
	usPerStep    = 1_000_000 / stepsPerSec
	degreesPerUs = 360.0 * float64(rpm) / (60 * 1_000_000)
	stepDelay    = 3000 // 3ms per traccia (simula il motore passo-passo)
)

var _sectorsByZone = [...]int{21, 19, 18, 17} // Settori per zona (0-3)
var _tracksPerZone = [...]int{17, 24, 30, 35} // Limiti tracce per zona

type DiskSimulator struct {
}

type FloppySimulator struct {
	currentAngle       float64
	currentTrack       int
	targetTrack        int
	currentZone        int
	sectorAngle        float64
	stepCounter        int
	stepDelayRemaining int
}

func (fs *FloppySimulator) Emulate() {
	// 1. Gestione rotazione costante
	fs.currentAngle += degreesPerUs * usPerStep
	if fs.currentAngle >= 360 {
		fs.currentAngle -= 360
	}

	// 2. Gestione movimento testina (se c'è un seek in corso)
	if fs.currentTrack != fs.targetTrack {
		if fs.stepDelayRemaining <= 0 {
			// Muovi di una traccia verso il target
			if fs.targetTrack > fs.currentTrack {
				fs.currentTrack++
			} else {
				fs.currentTrack--
			}
			fs.updateZone(fs.currentTrack)
			fs.stepDelayRemaining = stepDelay // Reimposta il delay
		} else {
			fs.stepDelayRemaining -= usPerStep
		}
	}

	// 3. Calcola settore corrente (basato sulla zona attuale)
	currentSector := int(fs.currentAngle / fs.sectorAngle)

	fmt.Println(currentSector)
	// [...] Qui gestisci la lettura del settore
}

func (fs *FloppySimulator) updateZone(currentTrack int) {
	for zone, maxTrack := range _tracksPerZone {
		if currentTrack <= maxTrack {
			fs.currentZone = zone
			fs.sectorAngle = 360.0 / float64(_sectorsByZone[zone])
			return
		}
	}
}

func (fs *FloppySimulator) SeekTrack(track int) {
	if track < 1 || track > 35 {
		return // Traccia non valida
	}
	fs.targetTrack = track
}

// Esempio d'uso:
func main() {
	sim := FloppySimulator{currentTrack: 18}
	sim.SeekTrack(25) // Cerca la traccia 25 (zona 2)

	for {
		sim.Emulate()
		time.Sleep(time.Microsecond)
	}
}

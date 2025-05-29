package mos6581

import (
	"log"
	"math"
)

// FilterType represents the type of filter applied in the audio processing pipeline.
//type FilterType int

// FilterNone represents no filtering.
// FilterLp represents a low-pass filter.
// FilterBp represents a band-pass filter.
// FilterLpBp represents a combination of low-pass and band-pass filters.
// FilterHp represents a high-pass filter.
// FilterNotch represents a notch filter.
// FilterHpBp represents a combination of high-pass and band-pass filters.
// FilterAll represents applying all filters.

/*
const (
	FilterNone = FilterType(iota)
	FilterLp
	FilterBp
	FilterLpBp
	FilterHp
	FilterNotch
	FilterHpBp
	FilterAll
)

*/

type FilterType uint8

const (
	FilterNone FilterType = 0b000
	FilterLp   FilterType = 0b001 // Bit 0 (Low Pass)
	FilterBp   FilterType = 0b010 // Bit 1 (Band Pass)
	FilterHp   FilterType = 0b100 // Bit 2 (High Pass)
)

// calcResonanceLp computes the resonance low-pass filter value based on the given input x using a polynomial equation.
func calcResonanceLp(x float64) float64 {
	v := 227.755 - (1.7635 * x) - (0.0176385 * x * x) + (0.00333484 * x * x * x) - (9.05683e-6 * x * x * x * x)
	return v
}

// calcResonanceHp computes the high-pass filter resonance value based on the input parameter x.
func calcResonanceHp(x float64) float64 {
	v := 366.374 - (14.0052 * x) + (0.603212 * x * x) - (0.000880196 * x * x * x)
	return v
}

// const resonanceSize = 1 << 11
const resonanceSize = 1 << 8

// Filters represents the state and configuration of an audio filter, including coefficients, resonance, and frequencies.
type Filters struct {
	filterType         FilterType // Filter type
	filterFreqHigh     uint8      // SID filter frequency (upper 8 bits)
	filterFreqLow      uint8
	filterRes          uint8   // Filter resonance (0..15)
	filterAmpl         float64 // IIR filter input attenuation;
	d1, d2, g1, g2     float64 // IIR filter coefficients
	xn1, xn2, yn1, yn2 float64 // IIR filter previous input/output signal
	resonanceLP        [resonanceSize]float64
	resonanceHP        [resonanceSize]float64
}

// NewFilters initializes a new Filters instance with default parameters and precomputes resonance lookup tables.
func NewFilters() *Filters {
	f := &Filters{
		filterType:     FilterNone,
		filterFreqHigh: 0,
		filterFreqLow:  0,
		filterRes:      0,
		filterAmpl:     1.0,
		d1:             0.0,
		d2:             0.0,
		g1:             0.0,
		g2:             0.0,
		xn1:            0.0,
		xn2:            0.0,
		yn1:            0.0,
		yn2:            0.0,
	}
	for i := 0; i < resonanceSize; i++ {
		f.resonanceLP[i] = calcResonanceLp(float64(i))
		f.resonanceHP[i] = calcResonanceHp(float64(i))
	}
	return f
}

// Reset reinitializes all filter parameters in the Filters struct to their default values, effectively clearing any previous state.
func (f *Filters) Reset() {
	f.filterType = FilterNone
	f.filterFreqHigh = 0
	f.filterRes = 0
	f.filterAmpl = 1.0
	f.d1 = 0.0
	f.d2 = 0.0
	f.g1 = 0.0
	f.g2 = 0.0
	f.xn1 = 0.0
	f.xn2 = 0.0
	f.yn1 = 0.0
	f.yn2 = 0.0
}

// Compute applies the filter defined in the Filters struct to the output signal and returns the resulting value.
func (f *Filters) Compute(outputFilter int32) int32 {
	xn := float64(outputFilter) * f.filterAmpl
	yn := xn + (f.d1 * f.xn1) + (f.d2 * f.xn2) - (f.g1 * f.yn1) - (f.g2 * f.yn2)
	f.yn2 = f.yn1
	f.yn1 = yn
	f.xn2 = f.xn1
	f.xn1 = xn
	outputFilter = int32(yn)
	return outputFilter
}

// UpdateFreqLow sets the low part of the filter frequency and recalculates filter coefficients if the value is changed.
func (f *Filters) UpdateFreqLow(data uint8) {
	if data != f.filterFreqLow {
		f.filterFreqLow = data
		f.compute()
	}
}

// UpdateFreqHigh updates the high part of the filter frequency and recalculates filter coefficients if the value changes.
func (f *Filters) UpdateFreqHigh(data uint8) {
	if data != f.filterFreqHigh {
		f.filterFreqHigh = data
		f.compute()
	}
}

// UpdateRes updates the filter resonance value using the upper 4 bits of the input and recalculates filters if enabled.
func (f *Filters) UpdateRes(data uint8) {
	resonance := data >> 4
	if resonance != f.filterRes {
		f.filterRes = resonance
		f.compute()
	}
}

// UpdateType adjusts the filter type of the Filters instance based on the provided data and resets filter state if changed.
func (f *Filters) UpdateType(data uint8) {
	v := FilterType((data >> 4) & 7)
	if v != f.filterType {
		f.filterType = v
		f.xn1 = 0.0
		f.xn2 = 0.0
		f.yn1 = 0.0
		f.yn2 = 0.0
		f.compute()
	}
}

func (f *Filters) compute() {
	// Controlla casi speciali per ottimizzazione
	if f.filterType == FilterNone {
		f.d1, f.d2, f.g1, f.g2 = 0.0, 0.0, 0.0, 0.0
		f.filterAmpl = 0.0
		return
	}

	// Calcola frequenza di taglio a 11 bit
	cutoff11bit := (uint16(f.filterFreqHigh) << 8) | uint16(f.filterFreqLow)

	// Usa tutti gli 11 bit per l'indice (0-2047)
	//filterIndex := cutoff11bit
	// Usiamo gli 8 bit più significativi
	filterIndex := uint8(cutoff11bit >> 3)

	// Determina quali filtri sono attivi
	hasLP := (f.filterType & FilterLp) != 0
	hasBP := (f.filterType & FilterBp) != 0
	hasHP := (f.filterType & FilterHp) != 0

	// Seleziona tabella di risonanza
	var fr float64
	if hasLP || hasBP {
		fr = f.resonanceLP[filterIndex]
	} else {
		fr = f.resonanceHP[filterIndex]
	}

	// Calcola argomento normalizzato
	arg := fr / float64(SampleFreq>>1)
	if arg > 0.99 {
		arg = 0.99
	} else if arg < 0.01 {
		arg = 0.01
	}

	// Calcola poli del filtro
	f.g2 = 0.55 + 1.2*arg*arg - 1.2*arg + float64(f.filterRes)*0.0133333333
	f.g1 = -2.0 * math.Sqrt(f.g2) * math.Cos(math.Pi*arg)

	// Aumenta risonanza se Band Pass è combinato con altri filtri
	if hasBP {
		f.g2 += 0.1
	}

	// Stabilizzazione filtro
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}

	// Calcola coefficienti in base alla combinazione
	switch {
	case hasLP && hasHP: // Notch Filter
		f.d1 = -2.0 * math.Cos(math.Pi*arg)
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)

	case hasLP && hasBP: // Low Pass + Band Pass
		f.d1 = 2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)

	case hasHP && hasBP: // High Pass + Band Pass
		f.d1 = -2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2)

	case hasLP: // Solo Low Pass
		f.d1 = 2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)

	case hasHP: // Solo High Pass
		f.d1 = -2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2)

	case hasBP: // Solo Band Pass
		f.d1 = 0.0
		f.d2 = -1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)

	default: // Combinazioni non gestite
		log.Printf("Unsupported filter combination: %b", f.filterType)
		f.filterAmpl = 0.0
	}
}

/*
// compute adjusts internal filter coefficients and characteristics based on filter type, frequency, and resonance settings.
func (f *Filters) compute() {
	var fr float64
	// Check for some trivial cases
	if f.filterType == FilterAll {
		f.d1 = 0.0
		f.d2 = 0.0
		f.g1 = 0.0
		f.g2 = 0.0
		f.filterAmpl = 1.0
		return
	} else if f.filterType == FilterNone {
		f.d1 = 0.0
		f.d2 = 0.0
		f.g1 = 0.0
		f.g2 = 0.0
		f.filterAmpl = 0.0
		return
	}

	// Combina i byte di frequenza per ottenere il valore a 11 bit (0-2047)
	cutoff11bit := (uint16(f.filterFreqHigh) << 8) | uint16(f.filterFreqLow)

	// Mappa il valore a 11 bit all'indice a 8 bit (0-255) per le tabelle di lookup
	// Usiamo gli 8 bit più significativi
	//filterIndex := uint8(cutoff11bit >> 3)
	filterIndex := uint8(cutoff11bit)

	// Calculate resonance frequency
	if f.filterType == FilterLp || f.filterType == FilterLpBp {
		fr = f.resonanceLP[filterIndex]
	} else {
		fr = f.resonanceHP[filterIndex]
	}
	// Limit to <1/2 sample frequency, avoid div by 0 in case FilterBp below
	arg := fr / float64(SampleFreq>>1)
	if arg > 0.99 {
		arg = 0.99
	}
	if arg < 0.01 {
		arg = 0.01
	}
	// Calculate poles (resonance frequency and resonance)
	f.g2 = 0.55 + 1.2*arg*arg - 1.2*arg + float64(f.filterRes)*0.0133333333
	f.g1 = -2.0 * math.Sqrt(f.g2) * math.Cos(math.Pi*arg)
	// Increase resonance if LP/HP combined with BP
	if f.filterType == FilterLpBp || f.filterType == FilterHpBp {
		f.g2 += 0.1
	}
	// Stabilize filter
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		if f.g1 > 0.0 {
			f.g1 = f.g2 + 0.99
		} else {
			f.g1 = -(f.g2 + 0.99)
		}
	}
	// Calculate roots (filter characteristic) and input attenuation
	switch f.filterType {
	case FilterLpBp, FilterLp:
		f.d1 = 2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)
	case FilterHpBp, FilterHp:
		f.d1 = -2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2)
	case FilterBp:
		f.d1 = 0.0
		f.d2 = -1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)
	case FilterNotch:
		f.d1 = -2.0 * math.Cos(math.Pi*arg)
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)
	default:
		log.Printf("SID FILTER NOT IMPLEMENTED %d\n", f.filterType)
	}
}
*/

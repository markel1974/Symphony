package mos6581

import (
	"log"
	"math"
)

/*
type FilterConfig struct {
	d1, d2, g1, g2, ampl float64
}

var _filterPresets = map[FilterType]FilterConfig{
	0b101: { // LP+HP (Notch Filter)
		d1:   -2.0,
		d2:   1.0,
		g1:   0.6,
		g2:   0.3,
		ampl: 0.22,
	},
	0b110: { // BP+HP
		d1:   -2.1,
		d2:   1.05,
		g1:   0.7,
		g2:   0.4,
		ampl: 0.18,
	},
	0b111: { // LP+BP+HP
		d1:   -2.0,
		d2:   1.0,
		g1:   0.65,
		g2:   0.35,
		ampl: 0.25,
	},
}

*/

// FilterType represents the type of filter configuration as a bitmask.
type FilterType uint8

// FilterNone represents no filter applied (binary value 000).
// FilterLp represents a low-pass filter (binary value 001, Bit 0).
// FilterBp represents a band-pass filter (binary value 010, Bit 1).
// FilterHp represents a high-pass filter (binary value 100, Bit 2).
const (
	FilterNone FilterType = 0b000
	FilterLp   FilterType = 0b001 // Bit 0 (Low Pass)
	FilterBp   FilterType = 0b010 // Bit 1 (Band Pass)
	FilterHp   FilterType = 0b100 // Bit 2 (High Pass)
)

// calcResonanceLp calculates the low-pass filter resonance factor based on the input value scaled to an 8-bit range.
func calcResonanceLp(x float64) float64 {
	x = x / 8.0 // Scalatura a un range equivalente a 8 bit (0 - 255.875)
	// I coefficienti del polinomio rimangono gli stessi
	v := 227.755 - (1.7635 * x) - (0.0176385 * x * x) + (0.00333484 * x * x * x) - (9.05683e-6 * x * x * x * x)
	return v
}

// calcResonanceHp calculates the resonance high-pass filter value based on the input parameter x.
// The input x is scaled to a range equivalent to 8 bits before applying the polynomial function.
// The function uses a cubic polynomial to compute the high-pass resonance value.
func calcResonanceHp(x float64) float64 {
	x = x / 8.0 // Scalatura a un range equivalente a 8 bit (0 - 255.875)
	// I coefficienti del polinomio rimangono gli stessi
	v := 366.374 - (14.0052 * x) + (0.603212 * x * x) - (0.000880196 * x * x * x)
	return v
}

// resonanceSize represents the size of the resonance arrays used in filter calculations, defined as 2048 (2^11).
const resonanceSize = 1 << 11

// Filters represents an audio filter system with various configuration and IIR filter coefficients.
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

// NewFilters initializes and returns a new instance of Filters with default values and precomputed resonance data.
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

// Reset reinitializes all filter parameters to their default values, clearing any previous configuration or state.
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

// Compute applies an IIR filter to the given output signal and updates internal state variables for future computations.
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

// UpdateFreqLow updates the low byte of the filter frequency if it differs from the current value and triggers computation.
func (f *Filters) UpdateFreqLow(data uint8) {
	if data != f.filterFreqLow {
		f.filterFreqLow = data
		f.compute()
	}
}

// UpdateFreqHigh updates the high byte of the filter frequency if the provided value is different from the current value.
// If updated, this method triggers a recalculation of the filter's coefficients by invoking the compute method.
func (f *Filters) UpdateFreqHigh(data uint8) {
	if data != f.filterFreqHigh {
		f.filterFreqHigh = data
		f.compute()
	}
}

// UpdateRes updates the filter resonance by extracting the upper 4 bits of the provided data and recalculates filter settings.
func (f *Filters) UpdateRes(data uint8) {
	resonance := data >> 4
	if resonance != f.filterRes {
		f.filterRes = resonance
		f.compute()
	}
}

// UpdateType updates the filter type based on the provided data and resets filter states if a change occurs.
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

// compute recalculates filter coefficients and configurations based on the current filter type, frequency, and resonance.
func (f *Filters) compute() {
	if f.filterType == FilterNone {
		f.d1, f.d2, f.g1, f.g2 = 0.0, 0.0, 0.0, 0.0
		f.filterAmpl = 0.0
		return
	}

	// Calcola frequenza di taglio a 11 bit
	cutoff11bit := (uint16(f.filterFreqHigh) << 8) | uint16(f.filterFreqLow)

	// Usiamo gli 8 bit più significativi
	//filterIndex := uint8(cutoff11bit >> 3)

	// Usa tutti gli 11 bit per l'indice (0-2047)
	filterIndex := cutoff11bit

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

	//if preset, ok := _filterPresets[f.filterType]; ok {
	//	f.d1 = preset.d1
	//	f.d2 = preset.d2
	//	f.g1 = preset.g1
	//	f.g2 = preset.g2
	//	f.filterAmpl = preset.ampl
	//	return
	//}

	// Calcola coefficienti in base alla combinazione
	switch {
	case hasLP && hasBP && hasHP: // LP+BP+HP (raro)
		f.d1 = -2.0 * math.Cos(math.Pi*arg)
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 - math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)

	case hasBP && hasHP: // BP+HP
		f.d1 = -2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg))

	case hasLP && hasHP:
		// Utilizza una formula più precisa basata su documentazione MOS6581
		notchWidth := 0.1 // Regolabile in base al resonance
		f.d1 = -2.0 * math.Cos(math.Pi*(arg+notchWidth/2))
		f.d2 = 1.0
		f.filterAmpl = 0.5 * (1.0 + f.g1 + f.g2) / (1.0 + notchWidth)

	case hasLP && hasBP: // Low Pass + Band Pass
		f.d1 = 2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)

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

//case hasHP && hasBP: // High Pass + Band Pass
//	f.d1 = -2.0
//	f.d2 = 1.0
//	f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2)

//case hasLP && hasHP: // Notch Filter
//	f.d1 = -2.0 * math.Cos(math.Pi*arg)
//	f.d2 = 1.0
//	f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)

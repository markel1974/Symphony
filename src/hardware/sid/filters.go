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

	ResonanceScaleFactor = 8.0 // 8 bit [0-255.875]
	ResonanceSize        = 1 << 11

	NotchWidthMax = 0.3  // Maximum width with zero resonance
	NotchWidthMin = 0.02 // Minimum width with resonance at 15
)

// calcResonanceLp calculates the low-pass filter resonance factor based on the input value scaled to an 8-bit range.
func calcResonanceLp(x float64) float64 {
	x = x / ResonanceScaleFactor
	// Polynomial coefficients
	v := 227.755 - (1.7635 * x) - (0.0176385 * x * x) + (0.00333484 * x * x * x) - (9.05683e-6 * x * x * x * x)
	return v
}

// calcResonanceHp calculates the resonance high-pass filter value based on the input parameter x.
// The input x is scaled to a range equivalent to 8 bits before applying the polynomial function.
// The function uses a cubic polynomial to compute the high-pass resonance value.
func calcResonanceHp(x float64) float64 {
	x = x / ResonanceScaleFactor
	// Polynomial coefficients
	v := 366.374 - (14.0052 * x) + (0.603212 * x * x) - (0.000880196 * x * x * x)
	return v
}

// Filters represents an audio filter system with various configuration and IIR filter coefficients.
type Filters struct {
	filterType         FilterType // Filter type
	filterFreqHigh     uint8      // SID filter frequency (upper 8 bits)
	filterFreqLow      uint8
	filterRes          uint8   // Filter resonance (0..15)
	filterAmpl         float64 // IIR filter input attenuation;
	d1, d2, g1, g2     float64 // IIR filter coefficients
	xn1, xn2, yn1, yn2 float64 // IIR filter previous input/output signal
	resonanceLP        [ResonanceSize]float64
	resonanceHP        [ResonanceSize]float64
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
	for i := 0; i < ResonanceSize; i++ {
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
	// Calculate 11-bit cutoff frequency
	cutoff := (uint16(f.filterFreqHigh) << 8) | uint16(f.filterFreqLow)
	// Use all 11 bits for the index (0-2047)
	filterIndex := cutoff
	// Determine which filters are active
	hasLP := (f.filterType & FilterLp) != 0
	hasBP := (f.filterType & FilterBp) != 0
	hasHP := (f.filterType & FilterHp) != 0
	// Select resonance table
	var fr float64
	if hasLP || hasBP {
		fr = f.resonanceLP[filterIndex]
	} else {
		fr = f.resonanceHP[filterIndex]
	}
	// Calculate normalized argument
	arg := fr / float64(SampleFreq>>1)
	if arg > 0.99 {
		arg = 0.99
	} else if arg < 0.01 {
		arg = 0.01
	}
	// Calculate filter poles
	f.g2 = 0.55 + (1.2 * arg * arg) - (1.2 * arg) + (float64(f.filterRes) * 0.0133333333)
	f.g1 = -2.0 * math.Sqrt(f.g2) * math.Cos(math.Pi*arg)
	// Increase resonance if Band Pass is combined with other filters
	if hasBP {
		f.g2 += 0.1
	}
	// Filter stabilization
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

	// Calculate coefficients based on the filter combination
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
		// Uses a formula based on MOS6581 documentation
		// Linearly interpolate between maximum and minimum width values
		resonanceRatio := float64(f.filterRes) / 15.0
		notchWidth := NotchWidthMax - (resonanceRatio * (NotchWidthMax - NotchWidthMin))
		//notchWidth := 0.1 // Adjustable based on resonance
		f.d1 = -2.0 * math.Cos(math.Pi*(arg+notchWidth/2))
		f.d2 = 1.0
		f.filterAmpl = 0.5 * (1.0 + f.g1 + f.g2) / (1.0 + notchWidth)

	case hasLP && hasBP: // Low Pass + Band Pass
		f.d1 = 2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)

	case hasLP: // Low Pass Only
		f.d1 = 2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)

	case hasHP: // High Pass Only
		f.d1 = -2.0
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2)

	case hasBP: // Band Pass Only
		f.d1 = 0.0
		f.d2 = -1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / math.Sin(math.Pi*arg)

	default:
		log.Printf("Unsupported filter combination: %b", f.filterType)
		f.filterAmpl = 0.0
	}
}

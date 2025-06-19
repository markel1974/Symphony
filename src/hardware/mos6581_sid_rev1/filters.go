package mos6581

import (
	"math"
)

// FilterType represents the type of audio filter used, defined as a 3-bit unsigned integer.
type FilterType uint8

// FilterNone represents no filter applied.
// FilterLp represents a low-pass filter (Bit 0).
// FilterBp represents a band-pass filter (Bit 1).
// FilterHp represents a high-pass filter (Bit 2).
// ResonanceScaleFactor defines the scaling factor for resonance calculations.
// ResonanceSize specifies the resolution size for resonance.
// NotchWidthMax defines the maximum notch width with zero resonance.
// NotchWidthMin defines the minimum notch width with a resonance level of 15.
const (
	FilterNone FilterType = 0b000
	FilterLp   FilterType = 0b001 // Bit 0 (Low Pass)
	FilterBp   FilterType = 0b010 // Bit 1 (Band Pass)
	FilterHp   FilterType = 0b100 // Bit 2 (High Pass)

	ResonanceMax         = 0.99
	ResonanceMin         = 0.01
	ResonanceScaleFactor = 8.0 // 8 bit [0-255.875]
	ResonanceSize        = 1 << 11

	NotchWidthMax = 0.1  // Maximum width with zero resonance
	NotchWidthMin = 0.04 // Minimum width with resonance at 15
)

// filterCalculatorFn defines a function type for implementing specific filter calculations within a filters system.
type filterCalculatorFn func()

// calcResonanceLp calculates the resonance low-pass value for a given input x based on a polynomial equation.
func calcResonanceLp(x float64) float64 {
	x = x / ResonanceScaleFactor
	// Polynomial coefficients
	v := 227.755 - (1.7635 * x) - (0.0176385 * x * x) + (0.00333484 * x * x * x) - (9.05683e-6 * x * x * x * x)
	return v
}

// calcResonanceHp calculates the resonance high-pass value for a given input x using a polynomial equation.
func calcResonanceHp(x float64) float64 {
	x = x / ResonanceScaleFactor
	// Polynomial coefficients
	v := 366.374 - (14.0052 * x) + (0.603212 * x * x) - (0.000880196 * x * x * x)
	return v
}

const filterCalculatorMask = 0x7

// Filters represents a configurable structure for various filter types and associated parameters in signal processing.
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
	filtersCalculator  []filterCalculatorFn
}

// NewFilters creates and initializes a new Filters instance with default values and pre-calculated resonance settings.
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

	for idx := range f.resonanceLP {
		lp := calcResonanceLp(float64(idx)) / float64(SampleFreq>>1)
		if lp > ResonanceMax {
			lp = ResonanceMax
		} else if lp < ResonanceMin {
			lp = ResonanceMin
		}
		f.resonanceLP[idx] = lp
	}
	for idx := range f.resonanceHP {
		hp := calcResonanceHp(float64(idx)) / float64(SampleFreq>>1)
		if hp > ResonanceMax {
			hp = ResonanceMax
		} else if hp < ResonanceMin {
			hp = ResonanceMin
		}
		f.resonanceHP[idx] = hp
	}
	f.filtersCalculator = f.buildFilters()
	return f
}

// Reset reinitializes all filter parameters to their default values, effectively clearing any applied filter settings.
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

// Compute applies a filtering algorithm to the given input and updates internal state variables for subsequent calculations.
func (f *Filters) Compute(outputFilter int32) float64 {
	xn := float64(outputFilter) * f.filterAmpl
	yn := xn + (f.d1 * f.xn1) + (f.d2 * f.xn2) - (f.g1 * f.yn1) - (f.g2 * f.yn2)
	f.yn2 = f.yn1
	f.yn1 = yn
	f.xn2 = f.xn1
	f.xn1 = xn
	return yn
}

// UpdateFreqLow updates the low-frequency filter value if the new value differs from the current one and triggers recalculation.
func (f *Filters) UpdateFreqLow(data uint8) {
	if data != f.filterFreqLow {
		f.filterFreqLow = data
		f.compute()
	}
}

// UpdateFreqHigh updates the high-frequency filter value if it differs from the current value and recalculates the filter.
func (f *Filters) UpdateFreqHigh(data uint8) {
	if data != f.filterFreqHigh {
		f.filterFreqHigh = data
		f.compute()
	}
}

// UpdateRes updates the filter's resonance based on the provided data and recalculates the filter if a change is detected.
func (f *Filters) UpdateRes(data uint8) {
	resonance := data >> 4
	if resonance != f.filterRes {
		f.filterRes = resonance
		f.compute()
	}
}

// UpdateType updates the filter type based on the input data and resets internal state if the type has changed.
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

// FilterIndex returns the calculated filter index based on the high and low frequency filter values.
func (f *Filters) FilterIndex() uint16 {
	return (uint16(f.filterFreqHigh) << 8) | uint16(f.filterFreqLow)
}

// buildFilters initializes and returns an array of filterCalculatorFn for handling various filter combinations.
func (f *Filters) buildFilters() []filterCalculatorFn {
	filterCalculators := make([]filterCalculatorFn, filterCalculatorMask+1)
	for idx := range filterCalculators {
		filterCalculators[idx] = f.filterCalculatorNone
	}
	filterCalculators[FilterNone] = f.filterCalculatorNone                   // 000
	filterCalculators[FilterLp] = f.filterCalculatorLP                       // 001
	filterCalculators[FilterBp] = f.filterCalculatorBP                       // 010
	filterCalculators[FilterLp|FilterBp] = f.filterCalculatorLPBP            // 011
	filterCalculators[FilterHp] = f.filterCalculatorHP                       // 100
	filterCalculators[FilterLp|FilterHp] = f.filterCalculatorLPHP            // 101
	filterCalculators[FilterBp|FilterHp] = f.filterCalculatorHPBP            // 110
	filterCalculators[FilterLp|FilterBp|FilterHp] = f.filterCalculatorLPBPHP // 111
	return filterCalculators
}

// compute executes the filter calculation logic based on the provided filter type and computes the filter amplitude.
func (f *Filters) compute() {
	f.filtersCalculator[uint8(f.filterType)&filterCalculatorMask]()
}

// computePoles calculates and returns the first and second pole values (g1, g2) based on the provided argument and filter parameters.
func (f *Filters) computePoles(arg float64) (float64, float64) {
	g2 := 0.55 + (1.2 * arg * arg) - (1.2 * arg) + (float64(f.filterRes) * 0.0133333333)
	g1 := -2.0 * math.Sqrt(g2) * math.Cos(math.Pi*arg)
	return g1, g2
}

// filterCalculatorNone resets all filter-related variables to 0. It clears any stored data or adjustments in the filter state.
func (f *Filters) filterCalculatorNone() {
	f.d1, f.d2, f.g1, f.g2 = 0.0, 0.0, 0.0, 0.0
	f.filterAmpl = 0.0
}

// filterCalculatorLP computes and sets filter properties specifically for a low-pass (LP) filter using resonance parameters.
func (f *Filters) filterCalculatorLP() {
	//LP or BP filter using resonanceLP
	arg := f.resonanceLP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//NOT BP filter: not increment
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	f.d1 = 2.0
	f.d2 = 1.0
	f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)
}

// filterCalculatorHP computes high-pass filter parameters using resonanceHP and updates filter coefficients.
func (f *Filters) filterCalculatorHP() {
	//HP filter using resonanceHP
	arg := f.resonanceHP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//NOT BP filter: not increment
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	f.d1 = -2.0
	f.d2 = 1.0
	f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2)
}

// filterCalculatorBP calculates band-pass filter parameters based on resonance and updates internal states accordingly.
func (f *Filters) filterCalculatorBP() {
	//LP or BP filter using resonanceLP
	arg := f.resonanceLP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//BP filter must increment
	f.g2 += 0.1
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	f.d1 = 0.0
	f.d2 = -1.0
	sinArg := math.Sin(math.Pi * arg)
	if sinArg == 0 {
		f.filterAmpl = 0
		return
	}
	f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg)) / sinArg
}

// filterCalculatorLPHP adjusts the filter parameters for a high-pass filter using resonance values and pole computations.
func (f *Filters) filterCalculatorLPHP() {
	//LP or BP filter using resonanceLP
	arg := f.resonanceLP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//NOT BP filter: not increment
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	//TODO Definire Correttamente i parametri di NotchWidthMax NotchWidthMin
	resonanceRatio := float64(f.filterRes) / 15.0
	notchWidth := NotchWidthMax - (resonanceRatio * (NotchWidthMax - NotchWidthMin))
	//notchWidth := 0.1
	f.d1 = -2.0 * math.Cos(math.Pi*(arg+notchWidth/2))
	f.d2 = 1.0
	// Here too, 1.0 + notchWidth could be close to 0 if notchWidth = -1, but this shouldn't happen.
	f.filterAmpl = 0.5 * (1.0 + f.g1 + f.g2) / (1.0 + notchWidth)
}

// filterCalculatorLPBP calculates low-pass and band-pass filter parameters.
// It updates the filter's internal state variables for resonance and amplitude.
// Ensures stability of the filter poles by checking and adjusting computed values.
func (f *Filters) filterCalculatorLPBP() {
	//LP or BP filter using resonanceLP
	arg := f.resonanceLP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//BP filter must increment
	f.g2 += 0.1
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	f.d1 = 2.0
	f.d2 = 1.0
	f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2)
}

// filterCalculatorHPBP calculates and updates high-pass and band-pass filter parameters for the filter object.
// It utilizes resonanceHP and computes poles based on the filter index, updating internal filter attributes accordingly.
func (f *Filters) filterCalculatorHPBP() {
	//LP or BP filter using resonanceLP
	arg := f.resonanceLP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//BP filter must increment
	f.g2 += 0.1
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	f.d1 = -2.0
	f.d2 = 1.0
	f.filterAmpl = 0.25 * (1.0 - f.g1 + f.g2) * (1.0 + math.Cos(math.Pi*arg))
}

// filterCalculatorLPBPHP calculates filter coefficients and amplitude based on resonance and current filter configuration.
func (f *Filters) filterCalculatorLPBPHP() {
	//LP or BP filter using resonanceLP
	arg := f.resonanceLP[f.FilterIndex()]
	f.g1, f.g2 = f.computePoles(arg)
	//BP filter must increment
	f.g2 += 0.1
	if math.Abs(f.g1) >= (f.g2 + 1.0) {
		f.g1 = math.Copysign(f.g2+0.99, f.g1)
	}
	f.d1 = -2.0 * math.Cos(math.Pi*arg)
	f.d2 = 1.0
	sinArg := math.Sin(math.Pi * arg)
	if sinArg == 0 {
		f.filterAmpl = 0
		return
	}
	f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 - math.Cos(math.Pi*arg)) / sinArg
}

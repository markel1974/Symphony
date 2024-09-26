package mos6581

import (
	"log"
	"math"
)

type FilterType int

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

func calcResonanceLp(X float64) float64 {
	v := 227.755 - 1.7635*X - 0.0176385*X*X + 0.00333484*X*X*X - 9.05683e-6*X*X*X*X
	return v
}

func calcResonanceHp(X float64) float64 {
	v := 366.374 - 14.0052*X + 0.603212*X*X - 0.000880196*X*X*X
	return v
}

type Filters struct {
	filterType         FilterType   // Filter type
	filterFreq         uint8        // SID filter frequency (upper 8 bits)
	filterRes          uint8        // Filter resonance (0..15)
	filterAmpl         float64      // IIR filter input attenuation;
	d1, d2, g1, g2     float64      // IIR filter coefficients
	xn1, xn2, yn1, yn2 float64      // IIR filter previous input/output signal
	resonanceLP        [256]float64 // shortcut for calc_filter
	resonanceHP        [256]float64
	useFilters         bool
}

func NewFilters(useFilters bool) *Filters {
	f := &Filters{
		useFilters: useFilters,
		filterType: FilterNone,
		filterFreq: 0,
		filterRes:  0,
		filterAmpl: 1.0,
		d1:         0.0,
		d2:         0.0,
		g1:         0.0,
		g2:         0.0,
		xn1:        0.0,
		xn2:        0.0,
		yn1:        0.0,
		yn2:        0.0,
	}
	for i := 0; i < 256; i++ {
		f.resonanceLP[i] = calcResonanceLp(float64(i))
		f.resonanceHP[i] = calcResonanceHp(float64(i))
	}
	return f
}

func (f *Filters) Reset() {
	f.filterType = FilterNone
	f.filterFreq = 0
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

func (f *Filters) Compute(outputFilter int32) int32 {
	if f.useFilters {
		xn := float64(outputFilter) * f.filterAmpl
		yn := xn + (f.d1 * f.xn1) + (f.d2 * f.xn2) - (f.g1 * f.yn1) - (f.g2 * f.yn2)
		f.yn2 = f.yn1
		f.yn1 = yn
		f.xn2 = f.xn1
		f.xn1 = xn
		outputFilter = int32(yn)
	}
	return outputFilter
}

func (f *Filters) UpdateFreq(data uint8) {
	if data != f.filterFreq {
		f.filterFreq = data
		if f.useFilters {
			f.compute()
		}
	}
}

func (f *Filters) UpdateRes(data uint8) {
	if (data >> 4) != f.filterRes {
		f.filterRes = data >> 4
		if f.useFilters {
			f.compute()
		}
	}
}

func (f *Filters) UpdateType(data uint8) {
	if FilterType((data>>4)&7) != f.filterType {
		f.filterType = FilterType((data >> 4) & 7)
		f.xn1 = 0.0
		f.xn2 = 0.0
		f.yn1 = 0.0
		f.yn2 = 0.0
		if f.useFilters {
			f.compute()
		}
	}
}

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

	// Calculate resonance frequency
	if f.filterType == FilterLp || f.filterType == FilterLpBp {
		fr = f.resonanceLP[f.filterFreq]
	} else {
		fr = f.resonanceHP[f.filterFreq]
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
	f.g1 = -2.0 * math.Sqrt(f.g2) * math.Cos(MPi*arg)
	// Increase resonance if LP/HP combined with BP
	if f.filterType == FilterLpBp || f.filterType == FilterHpBp {
		f.g2 += 0.1
	}
	// Stabilize filter
	if math.Abs(f.g1) >= f.g2+1.0 {
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
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(MPi*arg)) / math.Sin(MPi*arg)
	case FilterNotch:
		f.d1 = -2.0 * math.Cos(MPi*arg)
		f.d2 = 1.0
		f.filterAmpl = 0.25 * (1.0 + f.g1 + f.g2) * (1.0 + math.Cos(MPi*arg)) / math.Sin(MPi*arg)
	default:
		log.Printf("SID FILTER NOT IMPLEMENTED %d\n", f.filterType)
	}
}

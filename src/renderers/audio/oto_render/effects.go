package oto_render

import "math"

func Float32ToPCM16(sample float32) uint16 {
	const target16 = 1 << 15
	sample = clamp(sample)
	return uint16(sample * target16)
}

func AmplifyCopy(samples []float32, gain float32) []float32 {
	amplified := make([]float32, len(samples))
	for i, sample := range samples {
		amplified[i] = clamp(sample * gain)
	}
	return amplified
}

func Amplify(samples []float32, gain float32) {
	for i := range samples {
		samples[i] *= gain
		samples[i] = clamp(samples[i])
	}
}

func AutoGain(samples []float32, targetDB float32) []float32 {
	peak := findPeak(samples)
	// Calcola il guadagno necessario
	currentDB := 20 * float32(math.Log10(float64(peak)))
	gain := dBToLinear(targetDB - currentDB)
	return AmplifyCopy(samples, gain)
}

func SoftClip(samples []float32) []float32 {
	amplified := make([]float32, len(samples))
	for i, sample := range samples {
		amplified[i] = softClip(sample)
	}
	return amplified
}

func clamp(value float32) float32 {
	if value > 1.0 {
		return 1.0
	}
	if value < -1.0 {
		return -1.0
	}
	return value
}

func dBToLinear(db float32) float32 {
	return float32(math.Pow(10, float64(db)/20))
}

func softClip(value float32) float32 {
	// Curva di compressione non lineare
	if math.Abs(float64(value)) > 0.9 {
		return float32(math.Tanh(float64(value * 0.8)))
	}
	return value
}

func findPeak(samples []float32) float32 {
	max := float32(0)
	for _, s := range samples {
		if abs := float32(math.Abs(float64(s))); abs > max {
			max = abs
		}
	}
	return max
}

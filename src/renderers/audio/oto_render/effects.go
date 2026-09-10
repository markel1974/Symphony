package oto_render

import "math"

// Float32ToPCM16 converts a float32 sample in the range [-1.0, 1.0] to a PCM16 int16 value in the range [-32768, 32767].
// Values exceeding the PCM16 range are clamped to the minimum or maximum int16 limits.
func Float32ToPCM16(sample float32) int16 {
	const scaleFactor float32 = math.MaxInt16 + 1 // target16
	// sample is now in [-1.0, 1.0]
	//clampedSample := clamp(sample)
	scaledFloat := sample * scaleFactor // Float result in range [-32768.0, 32768.0]
	if scaledFloat > math.MaxInt16 {
		return math.MaxInt16 // Maximum positive int16 value
	} else if scaledFloat < math.MinInt16 {
		return math.MinInt16
	}
	return int16(scaledFloat)
}

// AmplifyCopy applies a gain factor to each sample in the input slice and returns a new slice with the amplified values.
func AmplifyCopy(samples []float32, gain float32) []float32 {
	amplified := make([]float32, len(samples))
	for i, sample := range samples {
		amplified[i] = clamp(sample * gain)
	}
	return amplified
}

// Amplify modifies an array of audio samples by multiplying each sample with the specified gain and clamping the result.
func Amplify(samples []float32, gain float32) {
	for i := range samples {
		samples[i] *= gain
		samples[i] = clamp(samples[i])
	}
}

// AutoGain adjusts the amplitude of the audio samples to match the specified target decibel (dB) level.
func AutoGain(samples []float32, targetDB float32) []float32 {
	peak := findPeak(samples)
	// Calculate the necessary gain
	currentDB := 20 * float32(math.Log10(float64(peak)))
	gain := dBToLinear(targetDB - currentDB)
	return AmplifyCopy(samples, gain)
}

// SoftClip applies a non-linear soft clipping function to a slice of audio samples to prevent hard clipping distortion.
func SoftClip(samples []float32) []float32 {
	amplified := make([]float32, len(samples))
	for i, sample := range samples {
		amplified[i] = softClip(sample)
	}
	return amplified
}

// clamp ensures the given value stays within the range of -1.0 to 1.0, inclusive.
func clamp(value float32) float32 {
	if value > 1.0 {
		return 1.0
	}
	if value < -1.0 {
		return -1.0
	}
	return value
}

// dBToLinear converts a decibel (dB) value to its linear scale equivalent.
func dBToLinear(db float32) float32 {
	return float32(math.Pow(10, float64(db)/20))
}

// softClip applies a non-linear compression curve to the input value, preserving dynamics while preventing harsh clipping.
func softClip(value float32) float32 {
	// Non-linear compression curve
	if math.Abs(float64(value)) > 0.9 {
		return float32(math.Tanh(float64(value * 0.8)))
	}
	return value
}

// findPeak calculates and returns the maximum absolute value in the given slice of float32 samples.
func findPeak(samples []float32) float32 {
	maxV := float32(0)
	for _, s := range samples {
		if abs := float32(math.Abs(float64(s))); abs > maxV {
			maxV = abs
		}
	}
	return maxV
}

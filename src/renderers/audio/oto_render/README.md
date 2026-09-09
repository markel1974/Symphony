# Dynamic Audio Sync & Resampling in Pure Go

This module (`oto_render`) solves one of the most notoriously difficult problems in real-time emulation and audio streaming: **Clock Drift Synchronization**.

## The Problem: Audio Clock Drift
When building an emulator or a real-time audio synthesizer, you have two independent clocks that will **never** perfectly align:
1.  **The Producer (The Emulator):** Generates audio samples at a rigid, mathematically perfect rate (e.g., the C64 SID chip running at 1.02 MHz).
2.  **The Consumer (The Host OS):** The physical sound card (or WebAudio in WASM) consuming samples at its own hardware clock (e.g., 44100Hz).

Because these clocks drift, a naive audio buffer will eventually fail:
*   If the producer is slightly faster, the buffer fills up, causing unacceptable audio lag.
*   If the consumer is slightly faster, the buffer empties, causing audio crackling and underruns.

## The Solution: A Dynamic State Machine
This module implements a dynamic, stateful resampler using a Circular Queue. Instead of dropping frames or injecting silence, it actively **squishes** or **stretches** the audio chunks in real-time using Linear Interpolation to keep the buffer at an optimal equilibrium.

### O(1) Branchless State Resolution
To maintain extreme performance (essential for sub-millisecond audio threads), the state machine avoids `if/else` chains. Instead, it maps the current fill level of the `CircularQueue` directly to an array of function pointers:

```go
// r.states maps the buffer fill count to the specific resolution strategy in O(1)
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.states[r.ring.Counter()&0xf](buf)
}
```

### The 3 States of Synchronization

#### 1. Equilibrium (`handleGood`)
If the buffer contains a safe amount of data (2-4 chunks), the consumer and producer are perfectly in sync. The chunk is popped and played identically without any processing overhead.

#### 2. Consumer is too fast (`handleTooFast`)
When the buffer runs low (1 chunk left), the sound card is consuming audio faster than the emulator generates it. To prevent a buffer underrun (crackling), we "buy time":
*   Pop the single chunk.
*   Use Linear Interpolation to **stretch** it to double its length.
*   Play the first half, and push the second half back into the queue.
*   *Result:* We generate more audio data from the existing waveform, giving the producer time to catch up, without the human ear noticing the pitch shift.

#### 3. Producer is too fast (`handleTooSlow`)
When the buffer is overfilling (lag is increasing), the emulator is producing audio faster than the sound card can play it. To catch up without brutally dropping frames:
*   Pop **two** chunks from the queue.
*   Use Linear Interpolation to **squish** (compress) both chunks down into a single chunk.
*   Play the resulting chunk.
*   *Result:* We consume data at twice the speed, safely lowering the buffer fill level and eliminating lag.

## Why this matters
By treating audio synchronization as a dynamic resampling problem rather than a hard-blocking buffer problem, this renderer ensures that the audio stream remains perfectly smooth, glitch-free, and tightly synchronized with the video frames, regardless of the host operating system's exact clock drift. All implemented in pure Go with zero CGo dependencies.

# MOS6581 SID Chip Emulator in Go

This project is a Go-based emulation of the MOS Technology 6581 Sound Interface Device (SID), the iconic sound chip used in the Commodore 64 and other 8-bit computers. It aims to replicate the sound generation capabilities of the SID, including its three distinct voices, various waveform types, envelope generators, and its characteristic resonant filter.

## Features

* **Three Independent Voices:** Each voice can be programmed with its own frequency, waveform, envelope, and modulation parameters.
* **Multiple Waveforms:**
    * Triangle
    * Sawtooth
    * Pulse (with variable pulse width)
    * Noise (pseudo-random)
    * Combined waveforms (Triangle+Sawtooth, Triangle+Pulse, Sawtooth+Pulse, Triangle+Sawtooth+Pulse - though combined waveforms often result in specific behaviors rather than simple mixing in a real SID, and the implementation reflects this through lookup tables or specific logic for mixed modes).
* **ADSR Envelope Generators:** Each voice has an Attack, Decay, Sustain, Release envelope generator for dynamic volume control.
* **Synchronization & Ring Modulation:** Voices can synchronize their phase or apply ring modulation effects with preceding voices.
* **Programmable Filter:**
    * Low-Pass
    * Band-Pass
    * High-Pass
    * Notch (Low-Pass + High-Pass)
    * Configurable cutoff frequency and resonance.
    * Ability to route individual voices through the filter.
* **Master Volume Control.**
* **External Audio Input (Conceptual):** While not fully detailed in processing, register bits for external input filtering are present.
* **Potentiometer Registers:** Read-only registers for emulating paddle inputs (`POTX`, `POTY`).
* **Oscillator/Envelope Read-Back:** Registers for reading the current state of oscillator 3 and envelope 3.
* **TEST Bit Functionality:** Implements the TEST bit behavior which can reset/freeze oscillators and affect waveform generation.
* **Reflection API:** Provides getter/setter methods for all SID registers, facilitating debugging and external control (via `mos6581_reflect.go`).
* **Component-Based Design:** Designed to be integrated as a component within a larger emulation system (e.g., a C64 emulator), as suggested by `factory.go` and `component.BaseComponent`.

## File Structure

The emulator is organized into several Go files:

* `mos6581.go`: The core SID emulation logic, including register handling, voice mixing, filter processing, and audio buffer generation.
* `voice.go`: Implements the logic for a single SID voice, including waveform generation, envelope processing, and modulation.
* `filters.go`: Implements the SID's resonant filter, including different filter types and coefficient calculations.
* `tables.go`: Contains precomputed lookup tables for waveforms, envelope rates, and various constants used in the emulation.
* `factory.go`: Provides a factory pattern for creating SID component instances, likely for use in a larger emulator framework.
* `mos6581_reflect.go`: Defines constants for register indices and provides a reflection interface for easy register access by name.

## Core Concepts

### Voices (`voice.go`)

Each of the three voices in the SID is an independent sound generator. Key aspects include:

* **Waveform Generation:** Selected via control registers. The `count` variable acts as a phase accumulator, incremented by `add` (derived from the frequency registers). Different waveforms (Triangle, Sawtooth, Pulse, Noise) are generated based on this accumulator. Combined waveforms often use lookup tables or specific logic derived from analyzing SID behavior. The `TEST` bit significantly alters waveform generation, often producing fixed or specific test outputs.
* **Envelope Generator (ADSR):** Controls the amplitude of the voice over time.
    * `EGState`: Tracks current phase (Attack, Decay, Sustain, Release, Idle).
    * `aAdd`, `dSub`, `rSub`: Rates for Attack, Decay, and Release phases, derived from lookup tables (`_eGTable`).
    * `sLevel`: Sustain level.
    * `gate`: Bit that triggers the Attack phase and holds the envelope in Sustain/Decay, or starts Release when cleared.
    * The `_eGDRShiftTable` is used to implement the non-linear decay/release rates characteristic of the SID.
* **Modulation:**
    * **Synchronization (`sync`):** The phase accumulator of a voice can be reset when the preceding voice's accumulator overflows.
    * **Ring Modulation (`ring`):** Modifies the phase of the triangle waveform based on the output of the preceding voice.
* **Pulse Width (`pw`):** For the pulse waveform, this 12-bit value determines the duty cycle.

### Filters (`filters.go`)

The SID features a versatile analog filter that can be applied to voice outputs.

* **Filter Types:** Low-Pass, Band-Pass, High-Pass, and combinations (e.g., Notch from LP+HP).
* **Cutoff Frequency:** A 11-bit value set by `fcLO` (3 bits) and `fcHI` (8 bits) registers. The implementation uses these bits to index into precomputed resonance tables.
* **Resonance (`filterRes`):** A 4-bit value that controls the 'peakiness' of the filter.
* **IIR Filter:** The filter is implemented as an Infinite Impulse Response (IIR) filter. The `compute()` method updates coefficients (`d1`, `d2`, `g1`, `g2`) based on the filter type, cutoff, and resonance.
    * `resonanceLP` and `resonanceHP` arrays store pre-calculated values based on polynomial functions to determine the filter's response.
    * The `arg` variable, derived from the cutoff frequency and resonance table, is crucial for calculating filter pole positions.
    * Different coefficient formulas are used for various combinations of LP, BP, and HP active modes.
    * Filter stabilization logic is included to prevent extreme behavior.

### Main SID Logic (`mos6581.go`)

This file ties everything together:

* **Register Handling:**
    * `registers`: A slice holding the current state of all 29 SID registers.
    * `WriteRegister(addr, data)`: Handles writes to SID registers. It updates the `registers` slice and then calls a specific write handler function for that register (from `writes` array). These handlers often trigger updates in voices or filters.
    * `ReadRegister(addr)`: Handles reads from SID registers. For most registers, it returns the stored value. Special read handlers (from `reads` array) exist for `OSC3` (reads voice 2's waveform output) and `ENV3` (reads voice 2's envelope level).
* **Sound Generation (`calcSoundBuffer`)**: This is the heart of the audio output.
    1.  **Volume Interpolation:** The master volume (from register `$D418`) is sampled periodically (via `Prepare()`) into `sampleBuf`. `calcSoundBuffer` interpolates these volume changes across the audio samples being generated in the current block to provide smoother volume transitions.
    2.  **Voice Processing Loop:** For each voice:
        * `ComputeEnvelopeGenerators()`: Updates the voice's envelope.
        * `UpdateCount()`: Advances the voice's phase accumulator.
        * `ComputeWaveForm()`: Generates the raw waveform output for the voice.
        * The output is scaled by the current envelope level.
        * Voice outputs are summed into either a `sumOutputFiltered` or `sumOutputNonFiltered` path, depending on whether the voice is routed to the filter.
    3.  **Filter Application:** `sumOutputFiltered` is processed by `filters.Compute()`.
    4.  **Mixing:** The filtered and non-filtered signals are added together.
    5.  **Master Volume:** The interpolated master volume is applied.
    6.  **Scaling:** The final signal is scaled down to fit into the output `soundBuffer`.
* **Audio Output:** The generated `soundBuffer` is passed to an `IAudioRender` interface, which handles the actual playback.
* **Initialization (`NewSID`, `Bind`, `Setup`):** Sets up voices, filters, register maps, and prepares for audio generation.

### Lookup Tables (`tables.go`)

To improve performance and emulate certain SID characteristics more accurately, several lookup tables are used:

* `_triTable`, `_triSawTable`, `_triRectTable`, `_sawRectTable`, `_triSawRectTable`: Used for generating combined or complex waveforms.
* `_eGTable`: Contains pre-calculated rates for the ADSR envelope generator phases.
* `_eGDRShiftTable`: Used for the non-linear decay/release characteristic of the SID envelope.
* Constants like `SampleFreq`, SID `Frequency`, `Cycles` per sample, `RegisterCount`, etc.

## How to Use (Integration)

This SID emulator is designed as a component to be integrated into a larger system, such as a Commodore 64 emulator.

1.  **Creation:**
    * Use `NewFactory().Create(...)` if using the provided factory pattern.
    * Or directly call `mos6581.NewSID(parent, factory, label, instance)`.
2.  **Setup & Binding:**
    * Call `sid.Setup()` for initial configuration.
    * Call `sid.Bind(sidSocket, fragFreq, rasters)` to initialize audio parameters like fragment frequency, buffer sizes, and associate with an audio rendering backend.
3.  **Register Access:**
    * **Writing:** Use `sid.WriteRegister(address, value)` to write to SID registers (e.g., `0xD400` to `0xD41C`). This will trigger internal state changes in voices and filters.
    * **Reading:** Use `sid.ReadRegister(address)` to read from SID registers.
4.  **Sound Generation Cycle:**
    * **`Prepare()`:** Call this method regularly (e.g., once per emulated raster line or a fixed number of times per frame) to sample the current master volume setting. This builds up `sampleBuf` which is used by `calcSoundBuffer`.
    * **`Update()`:** Call this method at the desired audio fragment rate (e.g., 50Hz or 60Hz for PAL/NTSC). This triggers `calcSoundBuffer()` to generate the next block of audio samples and sends it to the audio player.
5.  **Reset:**
    * Call `sid.Reset()` to reset the SID to its power-on state.

## Dependencies

* `github.com/markel1974/c64emu/src/component`
* `github.com/markel1974/c64emu/src/config`
* `github.com/markel1974/c64emu/src/references`
* `github.com/markel1974/c64emu/src/registry`

(These suggest the SID emulator is part of a larger `c64emu` project.)

## Notes and Considerations

* **Accuracy:** The level of detail, especially in filter coefficient calculations, waveform generation logic (including the TEST bit), and envelope characteristics, suggests a focus on accurate emulation.
* **Filter Implementation Details:** The filter coefficients are dynamically calculated based on a normalized argument `arg` derived from the cutoff frequency and resonance tables. Specific formulas are applied for different filter mode combinations (e.g., LP+HP, BP+HP).
* **TEST Bit:** The TEST bit (bit 3 of the voice control register) has specific effects on each waveform type (e.g., resetting phase, outputting fixed values, or altering noise generation). It also freezes envelope progression.
* **`SampleFreq`:** The code is hardcoded for a `SampleFreq` of 44100 Hz. Other audio constants are derived from this.
* **Reflection API (`mos6581_reflect.go`):** This provides a convenient way to inspect and modify SID registers using human-readable names, which is very useful for debugging or building development tools.

This README provides a comprehensive overview of your MOS6581 SID emulator. You can adapt and expand it further as needed.
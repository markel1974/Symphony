# MOS 6581 SID Chip Emulator in Go

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony)

This project is a high-fidelity emulation of the **MOS Technology 6581 Sound Interface Device (SID)**, the iconic sound chip from the Commodore 64 home computer, written entirely in Go. It aims to accurately replicate the behavior and sonic nuances of the original chip, providing a robust and performant component for integration into larger emulation frameworks.

## Key Features

This emulator implements all key features of the SID 6581 with a strong focus on accuracy.

* **Three Independent Voices**: Each of the three voices is a complete and programmable oscillator with its own frequency, waveform, and envelope parameters.

* **Multiple, Accurate Waveforms**:
    * Triangle, Pulse (with variable pulse width), and pseudo-random Noise.
    * **Non-Linear Sawtooth**: The sawtooth waveform is implemented using a model that simulates the 6581's DAC non-linearity, replicating its characteristic raw sound.
    * **Combined Waveforms**: Waveform combinations are handled via lookup tables to emulate the specific behavior of the SID, which often does not correspond to a simple mathematical sum.

* **Per-Voice ADSR Envelope Generators**: Each voice features a complete Attack/Decay/Sustain/Release envelope generator, including a simulation of the SID's characteristic non-linear decay/release rates via a lookup table.

* **Inter-Voice Modulation**:
    * **Oscillator Synchronization**: A voice can reset the phase of another voice's oscillator to create complex, metallic, and harmonic-rich sounds.
    * **Ring Modulation**: Applies ring modulation to the triangle waveform for dissonant and bell-like effects.

* **Multi-Mode Resonant Filter**:
    * Implementation of Low-Pass, Band-Pass, and High-Pass filter modes.
    * Handles all filter combinations, including a **dynamic-width Notch filter (LP+HP)** whose bandwidth correctly varies based on the resonance setting.
    * Includes **filter stabilization logic** to prevent instability at high resonance settings.
    * Provides the ability to route each voice through the filter individually.

* **Accurate `TEST` Bit Handling**: The `TEST` bit correctly resets the oscillator and noise generator phase **without altering the ADSR envelope state**, ensuring compatibility with games and demos that use advanced programming techniques.

* **`float32` Audio Pipeline**: The entire mixing and final output generation process uses floating-point arithmetic (`float32`) for maximum precision and audio quality before being sent to the audio driver.

* **Interpolated Master Volume**: Changes to the master volume are interpolated across audio samples to produce smooth transitions and prevent audible "clicks".

* **Reflection API for Debugging**: A comprehensive reflection API provides programmatic getter/setter access to all SID registers via human-readable names, facilitating debugging and external control.

## Architecture and Design

The emulator was designed with modern software engineering principles to ensure clarity, maintainability, and performance.

* **High-Fidelity Behavioral Model**: The approach is to replicate the observable *behavior* of the real chip with the highest possible fidelity, implementing its known non-linearities and idiosyncrasies through mathematical models and lookup tables.
* **Dispatch Table Pattern**: Throughout the package, state- and type-based decisions (waveform selection, envelope state, filter type, register handling) are implemented using **dispatch tables (arrays of functions)** instead of large `switch` statements. This architectural pattern, applied consistently, reduces branching, improves performance, and makes the code highly modular and readable.
* **Component-Based ("Headless") Design**: As indicated by the `factory.go` file, this package is designed as a "headless" component for integration into the "Symphony" emulation framework. The core SID logic is completely decoupled from the audio output and user interface.

## File Structure

* `mos6581.go`: Contains the core SID logic, register handling, voice mixing, and final audio buffer generation.
* `voice.go`: Implements a single SID voice, with all logic for waveform generation and ADSR envelope management.
* `filters.go`: Implements the multi-mode resonant filter, including coefficient calculations and the dispatch logic.
* `tables.go`: Contains constants and pre-calculated lookup tables for performance and accuracy.
* `mos6581_reflect.go`: Provides the reflection API for programmatic register access.
* `factory.go`: Implements the factory pattern for integration into the Symphony framework.
* `interpolations.go`: Provides utility functions, such as linear interpolation.
* `ARCHITECTURE.md`: Describes the SID register map.

## How to Use (Integration)

This SID emulator is designed as a component. The typical lifecycle is:
1.  **Creation**: Use `NewSID()` to create an instance.
2.  **Binding**: Call `Bind()` to connect it to an audio backend and configure timing parameters (e.g., 50/60Hz refresh rate).
3.  **Register Access**: Use `WriteRegister(address, value)` and `ReadRegister(address)` to simulate the emulated CPU's access to the SID registers.
4.  **Audio Generation Cycle**:
    * Call `Prepare()` regularly (e.g., once per emulated raster line) to sample the master volume state.
    * Call `Update()` at the audio buffer refresh rate (e.g., 50Hz). This method triggers `calcSoundBuffer()` to generate the block of audio samples and sends it to the player.
5.  **Reset**: Call `Reset()` to return the SID to its power-on state.

## Dependencies

* `github.com/markel1974/symphony/src/component`
* `github.com/markel1974/symphony/src/config`
* `github.com/markel1974/symphony/src/references`
* `github.com/markel1974/symphony/src/registry`

(These suggest the SID emulator is part of the larger `symphony` project.)
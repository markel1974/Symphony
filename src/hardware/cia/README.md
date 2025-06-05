# MOS 6526 CIA Emulation Package

This package provides a high-fidelity, cycle-accurate emulation of the **MOS 6526 Complex Interface Adapter (CIA)**, a critical component of the Commodore 64 home computer. It is designed for use in C64 emulators and other projects requiring a precise software model of the CIA's behavior.

## Overview

The `mos6526` package delivers a robust software representation of the CIA chip, meticulously engineered to replicate the hardware's functionality and timing characteristics. The primary design goal is to achieve **cycle-accuracy** to ensure compatibility with software that relies on the subtle and often undocumented behaviors of the original silicon. This package is intended for emulator developers seeking a high-quality, modular, and accurate CIA component.

## Key Features

* **Two Programmable 16-bit Timers**: Full emulation of two independent 16-bit timers (Timer A and Timer B). This includes support for various clocking sources (system $\phi_2$ clock, external CNT pin events), one-shot and continuous modes, and timer output to Port B pins.

* **Time-of-Day (TOD) Clock**: A 24-hour clock that maintains time in Binary-Coded Decimal (BCD) format. It features a programmable alarm, hardware-accurate register latching upon read, and a cycle-accurate internal clocking mechanism.

* **Serial Data Register (SDR)**: A complete 8-bit shift register for serial I/O. The implementation is accurately clocked by Timer A underflows, correctly handling the bit-by-bit shifting mechanism for both data transmission and reception.

* **Parallel I/O Ports**: Two 8-bit bidirectional I/O ports (Port A and Port B) with their corresponding Data Direction Registers (DDRs), managed via a decoupled interface for easy integration with keyboard, joystick, and other peripheral logic.

* **Comprehensive Interrupt System**: A robust interrupt system that correctly handles all 5 hardware interrupt sources (Timer A/B Underflow, TOD Alarm, SDR transfer complete, FLAG pin). It features a fully emulated Interrupt Control Register (ICR) and Interrupt Mask Register (IMR).

## Emulation Accuracy

This implementation goes beyond basic functionality to replicate the specific, nuanced behaviors of the physical hardware.

* **Cycle-Accurate Timer IRQs**: Timer underflow interrupts are correctly delayed by one clock cycle. The interrupt flag in the ICR is set on the cycle of the underflow, but the IRQ line is asserted on the *following* cycle. This is handled by a flag system within the main `Emulate()` loop and is critical for compatibility with timing-sensitive software.

* **Cycle-Accurate TOD Clocking**: The TOD clock is driven by an internal frequency divider within the `CIA.Emulate()` method. This design is both highly accurate—ticking at the precise system cycle—and highly performant, as the full BCD update logic is only executed when the divider reaches zero.

* **Advanced Timer Logic**: The implementation correctly distinguishes between the two different uses of the `CNT` pin:
  * **Timer A** correctly counts **rising edges** (pulses) on the CNT pin.
  * **Timer B** correctly checks the **static level** (high/low state) of the CNT pin during a Timer A underflow.
    This distinction is essential for software that uses these advanced timer modes.

* **Hardware-Accurate State Machine**: The timers are implemented with a detailed state machine (`TimerState`) that correctly models subtle behaviors such as the one-cycle startup delay and the "strobe" nature of the Force Load bit, which is actioned on write but never stored.

* **Register Latching**:
  * **TOD**: Reading the TOD Hours register correctly latches the clock's value, preventing race conditions. The clock is subsequently unlatched upon reading the Tenths of a Second register, mirroring the hardware.
  * **Timers**: The 16-bit timer latches are loaded under the correct hardware conditions, including on a Force Load strobe, on underflow, or after a write to the high-byte while the timer is stopped.

## Design Philosophy

* **Decoupling**: The CIA component is decoupled from the rest of the emulator system via the `ICIASocket` interface. This allows it to be integrated cleanly, with the "socket" providing the necessary connections to the system bus, interrupt lines, and I/O pins.
* **Composition**: The `CIA` struct is composed of smaller, dedicated components (`Timer`, `TOD`), each handling its own logic. The `CIA` acts as an orchestrator, managing the complex interactions between these parts, just like the physical chip.

## Package Structure

The code is organized into modules with clear responsibilities:

-   `cia.go`: Contains the main `CIA` struct, which orchestrates all sub-components and handles interaction with the system bus and interrupts.
-   `timer.go`: Implements the internal logic and complex state machine for the two programmable timers.
-   `tod.go`: Manages the state and BCD logic for the Time-of-Day clock and its alarm.
-   `factory.go`: Provides the factory pattern for registering and instantiating the `mos6526` component within a larger system.
-   `timer_reflect.go`, `tod_reflect.go`: Support files that expose the internal state of components for debugging and introspection purposes.

## TODO / Future Work

-   Complete the implementation of the `tod_reflect.go` file for full introspection capabilities.
-   Expand the unit test suite to cover more edge cases, particularly for complex interactions.

## License

This project is released under the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
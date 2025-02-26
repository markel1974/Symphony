# Package mos6522

The `mos6522` package provides an emulation of the MOS 6522 Versatile Interface Adapter (VIA), a chip commonly used in 8-bit systems such as the Commodore VIC-20 and other computers of the era. The 6522 VIA handles multiple I/O functions, timers, and interrupt-driven operations to facilitate communication between the CPU and peripheral devices.

## Overview

The `mos6522` package accurately emulates the MOS 6522, offering:

- **Bidirectional I/O Ports**: Support for two fully configurable 8-bit ports for input/output.
- **Timers**: Two flexible timers (Timer 1 and Timer 2) for precise timing and interval measurements.
- **Shift Register**: Emulation of the Serial Peripheral Interface for simple serial communication.
- **Interrupt Handling**: Full support for interrupt-driven operations, including event triggers from ports, timers, or external signals.

## Features

- **Port Management**: Configurable Port A and Port B with support for Data Direction Registers (DDRs) for managing bit-level input/output.
- **Timer Functionality**:
    - Timer 1: Designed for single-shot or continuous intervals.
    - Timer 2: Supports counting external pulses or clock cycles.
- **Shift Register**: Implements serial data transfer functionality, enabling communication with external devices.
- **Interrupt Controller**:
    - Edge-detection, pulsed, or latched interrupts.
    - Control via Interrupt Enable Register (IER) and Interrupt Flag Register (IFR).
- **Cycle-Accurate Emulation**: Ensures precise timer and register operations by accounting for individual clock cycles.

## Use Cases

The MOS 6522 VIA is typically used for:

- Managing I/O operations for peripherals such as keyboards, joysticks, or printers.
- Precise timing or delay generation using the built-in timers.
- Serial communications via the shift register.
- Interrupt-driven workflows to increase efficiency in system operations.

## Package Structure

- `via.go`: Defines the `VIA` struct, which represents the core functionality of the MOS 6522.
- `socket.go`: 

## Key Components

### I/O Ports

The MOS 6522 has two 8-bit I/O ports (Port A and Port B). Each bit in these ports can be independently set as input or output using their respective Data Direction Registers (DDRs).

### Timers

- **Timer 1**:
    - Programmable for single-shot (one-time) or continuous modes.
    - Used for generating delays or periodic triggers.
- **Timer 2**:
    - Designed for event counting or timing with optional external input.
    - Supports cascaded counting (linked with other timers).

### Shift Register

The shift register allows for serial communication by sending or receiving data one bit at a time. It is typically used for applications requiring simple serial output or input.

### Interrupts

The package fully implements interrupt generation and handling as per the 6522’s original design. Interrupts can be triggered by:

- Timer underflow (Timer 1 or Timer 2).
- Positive or negative edge detection on the ports.
- Serial data transfer completion or external events.

## Example Usage

```go
package main

import "github.com/yourproject/mos6522"

func main() {
    // Create a new VIA instance
    via := mos6522.NewVIA()

    // Configure Port A as output
    via.SetPortDirection(mos6522.PortA, 0xFF) // All bits as output

    // Write a value to Port A
    via.WritePort(mos6522.PortA, 0xAA) // 10101010 in binary
    
    // Configure Timer 1 for a single-shot interval
    via.SetTimer1(0xFFFF, mos6522.SingleShotMode)

    // Perform operations based on interrupts
    if via.CheckInterrupt() {
        fmt.Println("Interrupt triggered by Timer 1!")
    }
}
```

## Dependencies

This package has no significant external dependencies but integrates easily with other emulation packages for custom hardware emulators.

## Notes

- All core functionalities of the MOS 6522 are supported and aim to achieve cycle-accurate emulation.
- The package is optimized for use in larger emulation environments for classic 8-bit computers or embedded systems.

## TODO

- Add additional test coverage for the shift register and edge-case scenarios for interrupts.
- Improve configuration options for different clock speeds and external inputs.
- Document specific behaviors for undocumented or rarely used modes.

## Contributing

Contributions are welcome! To participate:

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Submit a pull request with a clear description of the changes.

For questions or discussions, please open an issue.

## License

This project is released under the [Apache 2.0 license](https://opensource.org/licenses/Apache-2.0).
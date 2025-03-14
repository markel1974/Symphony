# Package mos6526

This package implements the emulation of the MOS 6526 Complex Interface Adapter (CIA), a key component in the Commodore 64 responsible for I/O operations, timers, and interrupt management.

## Overview

The `mos6526` package provides a software representation of the CIA, including:

- **Peripheral Interface**: Management of data ports for user interfaces.
- **Timers**: Emulation of Timer A and Timer B, both essential for time-based operations in the system.
- **Time of Day (TOD) Clock**: Emulates the TOD functionality, which acts as a 24-hour clock.
- **Serial I/O**: Emulation of serial communication for peripherals.
- **Interrupt Handling**: Full support for the CIA’s interrupt system, including interrupt registers and flags.
- **Cycles Management**: Emulates the clock and timing behavior of the chip, ensuring accurate synchronization.

## Features

- Emulates the important functions of the MOS 6526 for robust integration into Commodore 64 emulators.
- Implements **I/O Port** handling for bidirectional communication.
- Supports the **Timer** functionality, including one-shot and continuous modes.
- Includes the **Time of Day Clock**, which is configurable and accurate to the original.
- Manages **Interrupt Triggers**, enabling event-driven emulation.

## Package Structure

The package is organized in a modular and extensible way, separating functionality into clear sections:

- `cia.go`: Provides the main `CIA` struct, which represents the MOS 6526 chip, and its core behavior.
- `timer.go`: Implements the timer functionality (Timer A and Timer B), including count modes and underflow handling.
- `tod.go`: Handles the Time of Day clock, offering accuracy and configurability.
- `ports.go`: Manages I/O ports, providing read and write functionalities.
- `interrupts.go`: Implements the CIA’s interrupt management system.
- `utils.go`: Provides utility functions for internal operations.
- `cia_test.go`: Includes unit tests to verify the emulation accuracy of different CIA features (timers, I/O, interrupts, etc.).

## Timers

The MOS 6526 includes two timers (Timer A and Timer B) that are used for various time-based operations. The `mos6526` package accurately emulates these timers, including:

- **Count Modes**:
    - System clock
    - External signal
- **Operation Modes**:
    - One-shot mode
    - Continuous mode
- Underflow behavior and interaction with the interrupt system.

## I/O Ports

The package fully supports the two 8-bit ports of the MOS 6526:

- Port A and Port B, which can be configured for input or output.
- Support for bit-level direction control, including Data Direction Registers (DDRs).
- Handles interactions with external devices accurately.

## Time of Day (TOD) Clock

The TOD clock is implemented with the following features:

- Accurate 24-hour format emulation with support for AM/PM mode.
- Full support for the divider chain behavior and adjustment registers.
- Interrupt capability when the TOD matches a preset alarm value.

## Interrupts

The interrupt system is fully implemented, including:

- Support for timer underflows, TOD match, and external interrupts.
- Emulation of the Interrupt Control Register (ICR) and Interrupt Mask Register (IMR).
- Accurate behavior for setting and clearing interrupt flags.

## Dependencies

The following external modules are utilized in this package:

- **Memory Handling**: Integration with the memory subsystem for register mapping.
- **Clock Synchronization**: Interfaces with the central clock mechanism to maintain cycle accuracy.

## Notes

- This module ensures high compatibility with software using the MOS 6526 interface.
- It aims for **cycle-accurate emulation** to maintain proper timing and synchronization within the larger system.

## TODO

- Add further unit tests for all edge cases, particularly on timer and TOD behaviors.
- Add a configuration system for variable clock speeds.
- Improve inline documentation and code clarity.

## Contributing

We welcome contributions to improve and extend the `mos6526` package. To contribute, please:

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request with a clear description.

For questions or support, please open an issue.

## License

This project is released under the [Apache 2.0 license](https://opensource.org/licenses/Apache-2.0).
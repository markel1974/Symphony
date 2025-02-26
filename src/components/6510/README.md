# Package mos6510

This package implements the emulation of the MOS 6510 microprocessor (a variant of the 6502) used in the Commodore 64.

## Overview

The `mos6510` package provides a software representation of the 6510 CPU, including:

*   Registers (A, X, Y, PC, SP, SR).
*   Processor flags (N, V, B, D, I, Z, C).
*   Instruction execution cycle (fetch, decode, execute).
*   Interrupt handling (NMI, IRQ, Reset).
*   Instruction implementation for the 6510 (using micro-operations).
*   Stack management.
*   Lookup tables.

## Package Structure

The package is organized into the following files:

*   `cpu.go`: Defines the `CPU` struct and the main methods for emulation (execution cycle, register management, etc.).
*   `instructions.go`: Contains the *declarations* of the functions implementing the individual instructions of the 6510 (divided into micro-operations).
*   `inst_*.go`: Contain the *implementation* of the micro-operations for the instructions, grouped by category (load/store, arithmetic, logic, etc.).
*   `opcodes.go`: Defines the dispatch tables (`_modeTable` and `_opTable`) that map opcodes to functions for addressing mode handling and instruction execution.
*   `stack.go`: Implements operations on the 6510 stack.
*   `utils.go`: Contains utility functions.
*   `interrupts_test.go`: Contains tests for interrupt handling.
*   `opcodes_test.go`: Contains tests for opcode operations.

## Implemented Instructions

[**TODO:** List *all* implemented instructions, with a brief description of each, their addressing modes, modified flags, and clock cycles. This can be done in table format or as a list.]

**Example:**

| Instruction | Addressing Mode            | Description                                    | Flags Affected | Cycles |
| :---------- | :------------------------- | :--------------------------------------------- | :------------- | :----- |
| LDA         | Immediate                  | Loads an immediate value into the accumulator. | N, Z           | 2      |
| LDA         | Zero Page                  | Loads a value from a Zero Page address.        | N, Z           | 3      |
| ...         | ...                        | ...                                            | ...            | ...    |

## Addressing Modes

[**TODO:** Describe the addressing modes of the 6502/6510, with examples.]

## Interrupts

[**TODO:** Explain how interrupts (NMI, IRQ, Reset) are handled.]

## Dependencies

*   `github.com/markel1974/c64emu/src/memory` (for memory access)
*   `github.com/markel1974/c64emu/src/components/quartz` (for clock management)
*   Other interfaces

## Notes

*   This emulator implements *all* undocumented instructions of the 6502/6510.
*   This emulator *aims* for cycle-by-cycle accuracy.

## TODO

*   Add unit tests for *all* instructions, in *all* addressing modes.
*   Improve error handling.
*   Add detailed comments to the micro-operations.
*   Complete the implementation of any missing instructions (if any are still missing).

## Contributing

[**TODO:** If you accept contributions, describe how to do so.]

## License

This project is released under the [Apache 2.0 license](https://opensource.org/licenses/Apache-2.0).
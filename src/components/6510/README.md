# Package mos6510

This package (`src/components/mos6510`) provides an emulation of the MOS 6510 microprocessor, a variant of the 6502, used in the Commodore 64 home computer.  The emulator is designed for accuracy, aiming for cycle-accurate emulation where feasible. It's used as the CPU component within the `g64` Commodore 64 emulator.

**Note:** This emulator is part of a larger project (`g64`) and is not intended for standalone use. It relies on other components of `g64` (specifically the `memory.Memory` interface and `IBanks`, `IPic` interfaces) for its operation.

## Architecture

The 6510 emulation is based on a micro-operation approach.  Each 6502/6510 instruction is broken down into a sequence of smaller operations, each corresponding (roughly) to a CPU clock cycle. This allows for a high degree of timing accuracy, which is essential for accurate emulation of the VIC-II and SID chips.

**Key Components:**

*   **`cpu.go`:** Contains the core `CPU` struct and its methods. This includes:
    *   `CPU` struct: Represents the state of the 6510 CPU (registers, flags, program counter, stack pointer, etc.).
    *   `Emulate()`: The main execution loop. This function repeatedly fetches and executes instructions until `c.stop` is true.
    *   `read()`: Reads a byte from the memory.
    *   `Reset()`: Resets the cpu to initial state.
    *   `popFlags()`: Reads the status register from memory.
    *   `pushFlags()`: Create the status register to write to the memory.
    * `branch()`: Handles the logic for branching.
    * `doADC()`: Implement the ADC instruction.
    * `doSBC()`: Implement the SBC instruction.
    *   Helper functions for flag manipulation (e.g., `SetFlagNZ`).
    *   Helper functions for address calculation, different addressing mode.

*   **`instructions.go`:** Contains *declarations* of all the functions that implement the individual 6510 instructions.  These functions are grouped by category in separate files (see below).

*   **`inst_*.go`:**  These files contain the *implementations* of the 6502/6510 instructions, grouped by category:
    *   `inst_load_store.go`:  Load/Store instructions (LDA, STA, LDX, STX, LDY, STY, etc.).
    *   `inst_arithmetic.go`: Arithmetic instructions (ADC, SBC, INC, DEC, etc.).
    *   `inst_logic.go`: Logical instructions (AND, ORA, EOR, BIT, ecc.).
    *   `inst_shift_rotate.go`: Shift and rotate instructions (ASL, LSR, ROL, ROR, ecc.).
    *   `inst_branch.go`: Branch instructions (BCC, BCS, BEQ, BNE, JMP, JSR, RTS, etc.).
    *   `inst_flag.go`: Flag manipulation instructions (CLC, SEC, CLI, SEI, CLD, SED, CLV).
    *   `inst_stack.go`: Stack instructions (PHA, PLA, PHP, PLP, TSX, TXS).
    *    `inst_interrupt.go`: Interrupt instructions
    *   `inst_transfer.go`: Register transfer instructions (TAX, TXA, TAY, TYA, etc.).
    *  `inst_control.go`: Other instructions.
    *   `inst_undocumented.go`: Undocumented (illegal) opcodes.

*   **`tables.go`:** Contains the dispatch tables (`opcodeTable` and `addressingModeTable`) that map opcodes to the corresponding addressing mode and instruction execution functions.

*   **`stack.go`:**  Provides functions for managing the 6510 stack.

*   **`utils.go`:** Contains utility functions.

**Execution Flow:**

The `Emulate` method in `cpu.go` is the main execution loop.  It repeatedly:

1.  Executes the function `next`.
2. The function read the `opcodeTable` for the next instruction.
3. The function read the `addressingModeTable` for the next addressing mode.
4. Execute the instruction.
5. Handles interrupts.

**Addressing Modes:**

The 6510 supports various addressing modes.  The functions to handle address calculation are mostly located within `cpu.go`.

**Interrupts:**

The 6510 supports three types of interrupts:

*   **NMI (Non-Maskable Interrupt):**  A high-priority interrupt that cannot be ignored by the CPU.
*   **IRQ (Maskable Interrupt):**  A lower-priority interrupt that can be enabled or disabled by the CPU.
*   **Reset:**  Resets the CPU to its initial state.

The `pic.go` and the methods in `cpu.go` handle interrupt generation and processing. The `IPic` interface is used to communicate with a separate interrupt controller.

**Status Register (Flags):**

The 6510 has a status register (SR) that contains several flags that reflect the state of the CPU and the result of the last operation.

*   **N (Negative):** Set if the result of the last operation was negative (bit 7 set).
*   **V (Overflow):** Set if the last operation resulted in a signed overflow.
*   **- (Unused):** This bit is always 1.
* **B (Break):**
*   **D (Decimal Mode):**  Used for BCD arithmetic (not fully supported in the original 6502, and often has different, or no, behavior on the 6510).
*   **I (Interrupt Disable):**  Set to disable IRQ interrupts.
*   **Z (Zero):** Set if the result of the last operation was zero.
*   **C (Carry):** Set if the last operation resulted in a carry (addition) or borrow (subtraction).

**Supported Instructions:**

All the instructions are implemented, the documentated and the not documentated instructions.

**Undocumented Instructions:**

The instructions are implemented: NOP, LAX, SAX, SLO, RLA, SRE, RRA, DCP, ISB, ANC, ASR, ARR, ANE, LXA, SBX, LAS, SHS, SHY, SHX, SHA and JAM.

## Usage

[**TODO:** Provide examples of how to use the `mos6510` package.  This should include:]

*   **Creating a `CPU` instance:**

    ```go
    // Assuming you have an implementation of memory.Memory called `myMemory`
    cpu := mos6510.NewCPU(myMemory)
    ```

*   **Loading code into memory:**  (using the `memory.Memory` interface).

*   **Setting the program counter (PC):** (using the `SetPC` method).

*   **Running the emulator:**  (using the `Run` method).

*   **Accessing registers:**  (using the `GetA`, `GetX`, `GetY`, `GetPC`, `GetSP` methods).

*   **Reading and writing memory:** (using the `Read` and `Write` methods).
* **Handling interrupts:**

## Dependencies
* `github.com/markel1974/c64emu/src/memory`
* `github.com/markel1974/c64emu/src/components/quartz`
* `github.com/markel1974/c64emu/src/bits`
* `github.com/markel1974/c64emu/src/c64/banks`

## Limitations

*   **No Unit Tests:** This package currently lacks unit tests. This means that the correctness of the emulation cannot be guaranteed.
*   **Incomplete Documentation:**  The documentation is incomplete.

## Contributing

[**TODO:** If you accept contributions, explain how to do it.]

## License

This project is released under the [Apache 2.0 License](https://www.google.com/url?sa=E&source=gmail&q=https://www.google.com/url?sa=E%26source=gmail%26q=https://opensource.org/licenses/Apache-2.0).
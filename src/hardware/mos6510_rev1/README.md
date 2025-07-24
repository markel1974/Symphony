# Package mos6510

This package (`src/components/mos6510`) provides an emulation of the **MOS 6510** microprocessor, a variant of the 6502 used in the Commodore 64. The emulator is designed for high fidelity, aiming for cycle-accurate emulation.

**Note:** This emulator is part of a larger project and is not intended for standalone use. It relies on other interfaces from the project (such as `IMos6510Banks` and `IQuartz`) for its operation.

---

## Architecture and Advanced Design

The emulator is distinguished by two key architectural choices that ensure high performance and strong modularity.

### 1. Near-Zero Branch Execution via a State Machine

Unlike a traditional approach based on a large `switch` statement for opcode decoding, this implementation uses a **finite state machine** based on function pointers.

The core of the emulation loop is a single function pointer, `cpu.next`, which points directly to the next micro-operation to be executed. Each 6510 instruction is broken down into a sequence of these micro-operations, each roughly corresponding to a single CPU clock cycle.

This approach, sometimes called *threaded code*, offers significant advantages:
* **High Performance**: It nearly eliminates branching in the emulator's critical *hot-path*, drastically reducing `branch misprediction` on the host CPU and improving instruction cache usage.
* **Cycle Accuracy**: It allows for very fine-grained control over timing, which is essential for correctly emulating synchronization with complex chips like the VIC-II and SID.

### 2. Component-Based Architecture with a Pluggable `ControlUnit`

The design is strongly decoupled and based on components with well-defined responsibilities:
* **`CPU`**: The main orchestrator that maintains register state but delegates complex operations.
* **`Bus`**: Manages all memory access logic, abstracting away the complexities of the underlying hardware.
* **`Interrupts`**: Encapsulates the logic for handling IRQ, NMI, and RESET signals.
* **`ControlUnit`**: Acts as the "catalog" for all micro-operations and addressing modes.

The `CPU` does not depend on a concrete implementation of the `ControlUnit` but interacts with it solely through the `IControlUnit` interface. This makes the **`ControlUnit` completely pluggable**, allowing one to:
* Swap out the entire instruction set to emulate different chip revisions.
* Inject a debug `ControlUnit` for tracing.
* Mock the component during unit tests to isolate CPU behavior.

---

## Execution Flow

The `Emulate()` method in `cpu.go` is the entry point for the execution cycle. On each call, it simply executes the function pointed to by `cpu.next`.

The flow is as follows:
1.  The current micro-operation is executed (e.g., `InstOpINI`).
2.  This function reads the opcode, determines the correct sequence of micro-operations for the addressing mode and the instruction itself, and sets `cpu.next` to the next function in the chain.
3.  The cycle repeats.
4.  Interrupts are handled cleanly within this flow, checked only at the beginning of a new instruction cycle (`InstOpINI`).

---

## Supported Instructions

All instructions are implemented, both **documented** and **undocumented** (illegal), including: NOP, LAX, SAX, SLO, RLA, SRE, RRA, DCP, ISB, ANC, ASR, ARR, ANE, LXA, SBX, LAS, SHS, SHY, SHX, SHA, and JAM.

The organization of `inst_*.go` files by category (arithmetic, logic, branch, etc.) keeps the code organized and maintainable.

---

## Limitations and Future Development
* **Lack of Unit Tests:** The package needs a comprehensive test suite to ensure the correctness of every instruction and addressing mode.
* **Incomplete Documentation:** The code documentation can be further improved.

---

## License

This project is released under the [Apache 2.0 License](https://www.google.com/url?sa=E&source=gmail&q=https://www.google.com/url?sa=E%26source=gmail%26q=https://opensource.org/licenses/Apache-2.0).
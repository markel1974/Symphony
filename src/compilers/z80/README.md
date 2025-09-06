# Z80 Transpiler for the Go-Native VM

This document provides a detailed overview of the Z80 transpiler, a sophisticated compiler frontend designed to translate Z80 machine code into the high-performance bytecode of the Go-Native Virtual Machine (VM). This project is not a traditional interpreter-based emulator; instead, it is a true transpiler that re-targets Z80 instructions to a custom, high-level bytecode format.

This approach serves as a powerful demonstration of the VM's core architectural principle: the **"interchangeable instruction disk."** By swapping the compiler frontend, the VM is transformed from a runtime for a Go-like language into a high-performance execution engine for a classic 8-bit CPU architecture.

## Architectural Concept: Emulation via Transpilation

The core strategy is to translate an entire Z80 ROM into an equivalent bytecode program that can be executed by the VM. This is achieved through two key concepts:

1.  **CPU State Mapping**: The Z80's internal state—including 8-bit and 16-bit registers (A, F, BC, DE, HL, SP, PC) and the 64KB RAM—is mapped directly to global variables within the VM's scope. This allows the VM's bytecode to manipulate the Z80's state directly, as if they were native variables.

2.  **1-to-1 Instruction Translation**: The primary goal is to achieve a **1-to-1 mapping** where a single Z80 instruction is translated into a single, powerful VM bytecode instruction. This avoids the overhead of interpreting many small bytecode operations for one Z80 instruction, leading to significantly higher performance.

## The Three-Tiered 1-to-1 Translation Strategy

To achieve maximum efficiency, the transpiler employs a sophisticated, multi-tiered strategy for mapping Z80 opcodes to VM bytecode.

### Tier 1: Simple Instructions -> Direct Bytecode Mapping
Simple, register-to-register Z80 instructions are mapped to a single, equivalent VM bytecode opcode. This provides a direct and highly efficient translation.

* **Example**: The Z80 instruction `LD B, C` (copy the value from register C to register B) is translated directly into a single `OpGlobalCopy` instruction. The VM executor for `OpGlobalCopy` then performs a direct memory copy between the global slots representing the B and C registers.

### Tier 2: ALU Instructions -> Specialized "Fast Path" Opcodes
All standard 8-bit arithmetic and logical operations (`ADD`, `SUB`, `AND`, `OR`, `XOR`, `CP`) are handled by a dedicated set of specialized VM opcodes: `OpIntArithmetic` and `OpIntLogical`.

This is a critical "fast path" optimization:
* These opcodes bypass the VM's stack entirely for their operands.
* They take the indices of the source and destination registers as direct arguments within the instruction itself.
* The executor, written in native Go, performs the calculation and updates the destination register directly in the VM's global state.

* **Example**: A Z80 `ADD A, B` instruction is translated into a single `OpIntArithmetic` instruction. The executor for this opcode reads the values from the global slots for registers A and B, calculates the sum, and writes the result back to the slot for register A, all within a single, atomic operation.

### Tier 3: Complex Instructions -> SDK Bridge via `OpCallImportGlobal`
The most complex Z80 instructions, especially those involving intricate flag calculations or multi-byte stack operations, are mapped to a single `OpCallImportGlobal` instruction. This instruction calls a highly optimized function written in native Go, which is exposed to the VM via a dedicated Z80 SDK.

This approach offers the best of both worlds:
* **Maximum Performance**: The complex logic is executed as a single, fast native Go function, avoiding any interpretation overhead.
* **Clean Bytecode**: The resulting bytecode remains clean and high-level, with a single instruction representing a complex Z80 operation.

* **Example (Conceptual)**: A full ALU operation that correctly updates all six of the Z80's condition flags (S, Z, H, P/V, N, C) would be handled by a single call to a function like `z80.alu`. The transpiler would emit the bytecode to call this function, and the `alu` function in the SDK would perform the arithmetic and all the complex bitwise logic for the flags in native Go. Similarly, `PUSH AF` would be a single call to `z80.push16`.

## Implementation Highlights

* **Helper Abstraction**: The `z80/helper.go` module provides a clean abstraction layer that encapsulates the logic for emitting the correct bytecode for different Z80 instructions. This keeps the main transpiler loop in `z80/compiler.go` clean and focused on dispatching opcodes.

* **Control Flow Mapping**: Z80 control flow instructions like `JP`, `CALL`, and `RET` (both conditional and unconditional) are intelligently mapped to the VM's own powerful control flow opcodes, such as `OpJump`, `OpJumpTruthy`, and `OpJumpIndirect`.

## Conclusion

The Z80 transpiler is a powerful testament to the flexibility of the underlying VM. By leveraging a multi-tiered 1-to-1 translation strategy, it transforms the VM into a high-performance engine for executing Z80 machine code. This approach is more efficient and elegant than traditional emulation, resulting in a compact bytecode representation and a fast, modern runtime for classic 8-bit software.
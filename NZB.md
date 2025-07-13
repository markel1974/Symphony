# The NZB (Near Zero Branch) Pattern

> A software design pattern that replaces conditional branching in performance-critical loops with pre-calculated dispatch tables, shifting complexity from runtime execution to a setup-time "compilation" phase.

## 1. Overview

The **NZB (Near Zero Branch)** pattern is an architectural philosophy for building high-performance, state-driven systems, particularly in domains like hardware emulation, game engines, and interpreters. Its primary goal is to minimize or eliminate conditional branches (`if`, `switch`) from "hot loops"—code paths that are executed millions of times per second.

The core problem it solves is the performance penalty associated with **branch misprediction** in modern CPUs. In a critical loop, a conditional statement can stall the processor's pipeline if its outcome is not correctly predicted. The NZB pattern avoids this by transforming a decision-making process into a direct, predictable memory lookup.

## 2. The Core Principle: Pre-calculate to Solve

The fundamental idea behind NZB is to shift the "decision" logic out of the execution loop. Instead of evaluating conditions at runtime, the system pre-calculates the correct execution path based on the current state and stores it in a **dispatch table** (an array of function pointers).

#### Traditional Approach (with Branching)

A typical state machine in a hot loop might look like this:

```go
func execute_hot_loop(state int, data an_object) {
  // This switch statement introduces branching
  switch state {
  case STATE_A:
    handle_state_a(data)
  case STATE_B:
    handle_state_b(data)
  case STATE_C:
    handle_state_c(data)
  default:
    handle_default(data)
  }
}
The NZB Approach (Zero Branch)
The NZB pattern re-structures this by creating a dispatch table at setup time.

Setup / State Change:

Go

// This dispatch table is our "pre-calculated" execution path
var dispatch_table [STATE_COUNT]func(an_object)

func setup() {
  dispatch_table[STATE_A] = handle_state_a
  dispatch_table[STATE_B] = handle_state_b
  dispatch_table[STATE_C] = handle_state_c
  // ...
}
Hot Loop Execution:

Go

func execute_hot_loop(state int, data an_object) {
  // A single, predictable lookup and direct function call
  dispatch_table[state](data)
}
The result is a branch-free execution path that is both faster and, in many ways, architecturally cleaner.

3. Key Characteristics & Benefits
High Performance: By eliminating conditional branches from critical paths, it minimizes the risk of CPU pipeline stalls, leading to significant performance gains.
Architectural Elegance & Clarity: The main execution loop becomes trivial. The logic for each state is neatly encapsulated in its own dedicated function, aligning well with the Strategy design pattern.
Enhanced Maintainability: Adding a new state is a clean operation: create a new function and add it to the dispatch table during setup. The core execution logic remains untouched.
"Pseudo-JIT" Behavior: The pattern effectively acts as a "Just-In-Time" compiler for the emulator's own logic. It "compiles" a system state into a highly optimized, direct execution path.
4. Canonical Implementation: The Symphony Framework
The Symphony emulation framework serves as the reference implementation of the NZB pattern, applying it pervasively across every major component to achieve a unique combination of cycle-exact fidelity and high performance.

CPU (MOS 6510): The CPU uses a two-phase dispatch system. An initial table maps an opcode to its addressing mode logic. Once the address is resolved, a second table maps the opcode to its ALU operation logic. This avoids complex switch statements in the instruction execution cycle.
Video Chip (MOS 6569 VIC-II): The VIC-II emulation is the purest expression of NZB. The 63 clock cycles of a PAL scanline are modeled as a pre-linked chain of 63 distinct functions. The main emulation loop is a single line that executes the current cycle's function and advances the pointer to the next, creating a "film strip" of cycle-exact hardware events with zero branching.
Memory Controller (PLA): The applyMemoryConfig function acts as the NZB "compiler." It takes the current memory configuration and builds the bankRead and bankWrite dispatch tables. A memory access from the CPU becomes a single, direct lookup into these tables to find the correct handler (RAM, ROM, or I/O).
Audio Renderer: The adaptive audio renderer uses a dispatch table to handle buffer pressure. The number of chunks in the audio buffer is used as a direct index into a states array, which immediately calls the correct resampling strategy (handleGood, handleTooFast, handleTooSlow) without any if statements in the critical audio callback path.
5. Conclusion
The NZB (Near Zero Branch) pattern is a powerful architectural strategy for developing high-performance, state-driven software. By formalizing the principle of "pre-calculating to solve," it provides a robust and elegant method for building complex systems that are both exceptionally fast and architecturally sound. The Symphony framework stands as a testament to the power and versatility of this design philosophy.
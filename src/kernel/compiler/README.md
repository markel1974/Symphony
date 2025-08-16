
# Go-Based Compiler: An Architectural Overview

This document provides a detailed overview of the Go-based compiler designed to translate a subset of the Go language into bytecode for its companion Virtual Machine. The compiler has a pragmatic design, leveraging the Go ecosystem to create a robust, powerful, and maintainable compilation pipeline.

## Core Architectural Philosophy

The compiler's primary goal is to provide a seamless and efficient bridge between a high-level, Go-like syntax and the low-level bytecode understood by the VM. Its architecture is built on several key principles:

-   **Pragmatism over Reinvention:** The compiler intelligently reuses existing, high-quality components from the Go toolchain, focusing its efforts on the core logic of compilation rather than on parsing and syntax definition.
-   **Robustness through Structure:** It employs a multi-pass compilation strategy to correctly handle the complexities of a structured programming language, such as forward references and mutual recursion.
-   **Seamless Integration:** The compiler is not a standalone component but a perfectly integrated counterpart to the VM. It "thinks" in the same object model and produces bytecode tailored precisely for the execution engine.

---

## Key Architectural Components

### 1. The Parser: Leveraging the Go Toolchain

Instead of defining a custom language and writing a parser from scratch, the compiler makes a brilliant and pragmatic choice: it uses the official `go/parser` and `go/ast` packages.

-   **A Subset of Pure Go:** The language this compiler accepts is not merely "Go-like"; it is a **syntactically valid subset of the Go language**. This means any script written for the VM can be analyzed, formatted, and linted by standard Go tools.
-   **Immediate Robustness:** By using Go's own parser, the compiler inherits a production-grade, battle-tested component, eliminating an entire class of potential bugs and development effort.
-   **Developer Familiarity:** Developers can write scripts in a familiar syntax, significantly lowering the barrier to entry.

### 2. The Compilation Pipeline: A Multi-Pass Strategy

The compiler avoids the limitations of a single-pass approach by orchestrating a sophisticated multi-pass compilation process. This is most evident in the `doFile` function:

1.  **Pass 1: Declaration Separation:** The compiler first traverses the AST to separate function declarations from all other top-level statements (globals, imports).
2.  **Pass 2: Function Pre-definition:** It then processes only the function signatures to create "placeholder" objects in the constants table. This is a **critical step** that reserves a stable, known index for every function before their bodies are compiled.
3.  **Pass 3: Global Compilation:** With function indices locked in, the compiler processes all global variable declarations and other non-function code.
4.  **Pass 4: Function Body Compilation:** Finally, it revisits each function to compile its body. At this stage, it can correctly resolve all symbols—including calls to functions defined later in the file (forward references) or mutual recursion—because all global symbols and function signatures are already known.

This multi-pass architecture is what enables the compiler to correctly handle a structured language.

### 3. The Standard Library (`stdlib`): A Bridge to Go's Power

One of the compiler's most powerful features is its standard library, which exposes the richness of Go's native `stdlib` to the scripting environment.

-   **Elegant Wrapper Pattern:** The system uses a set of clever function wrappers (e.g., `FuncASRS`, `FuncAIIRE` defined in `objects/function.go`) to adapt native Go functions to the VM's `IObject` interface.
-   **Extensive Functionality:** This approach provides the scripting language with powerful, high-performance modules for file I/O (`os`), text manipulation (`text`, `regexp`), JSON processing (`json`), mathematics (`math`), and more, with minimal effort.
-   **Modular and Resolvable:** The `Loader` component manages these modules, allowing the compiler to resolve selector expressions like `fmt.Println` into a `OpReferences` instruction that the VM can use to look up the correct function at runtime.

### 4. Symbol and Scope Management

At the core of the compiler is a robust system for managing scopes and symbols (`SymbolTable` and `Scopes`).

-   **Nested Scopes:** The system correctly handles nested scopes, ensuring that variables are resolved to the correct context (global or local).
-   **Correct Closure Implementation:** Crucially, the symbol table correctly identifies "free variables" (variables referenced from an inner function but defined in an outer scope). It compiles them into `OpGetFree`/`OpSetFree` instructions, enabling the VM to implement closures properly.

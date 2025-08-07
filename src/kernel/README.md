# The Symphony MicroKernel Framework

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony)

---

## 1. Overview: A Framework for Introspective Systems

The **Symphony MicroKernel Module** is an advanced, Go-based framework for creating powerful, interactive, and stateful command-line applications. It is architected as a modular **microkernel** designed to serve as a real-time control, introspection, and live debugging interface for any complex system.

Its core philosophy is that a complex application must be **transparent and debuggable by design**. The framework provides the foundation to build systems where the internal state and behavior of any component can be inspected and controlled at runtime via a secure SSH console. This turns the command line into a powerful live diagnostics and recovery tool, enabling developers to perform interactive 'software surgery' on a running system.

---

## 2. Microkernel Architecture

The project fully embraces a **microkernel design philosophy**: a minimal, robust core provides essential services, while all other functionalities are delegated to isolated processes.

* **Central Kernel (`kernel/core/kernel.go`):** The heart of the system. It manages the entire lifecycle of processes, orchestrates I/O, and runs a central, asynchronous `eventLoop`. Its primary role is to act as a **message router**, decoupling the system's components.
* **Dynamic Message Handling:** The kernel uses a dynamic `handlers` map. Each peripheral `Server` (e.g., for rendering, filesystem) registers the message types it can handle. This eliminates hardcoded dependencies and allows for a truly modular system.
* **The Shell as a User Process (`xsh`):** As a testament to the purity of the design, the main shell (`xsh`) is not a privileged part of the kernel but is itself an application that runs as a regular "Process".

### Kernel-Space / User-Space Process Separation

A fundamental aspect of the architecture is the strict separation of concerns, mirroring a traditional operating system's design. Every process is composed of two distinct parts:

1.  **Kernel-Space Representation (`KernelProcess`):** This is the "Process Control Block" that the Kernel sees and manages. It is a lightweight wrapper containing only the metadata essential for the kernel's operation: PID, parent-child relationships, security flags (like `protected`), and a list of kernel-managed resources (like timers). The Kernel interacts exclusively with this representation, ensuring it cannot access or corrupt the application's internal state.
2.  **User-Space Implementation (`Process`):** This is where the application logic resides. It contains the command to be executed, the private message queue for its event loop, and the application-specific context (e.g., the state of a game). This space is completely isolated from the kernel.

This separation provides strong security and stability boundaries, preventing a bug in an application from destabilizing the kernel.

### True Asynchronous IPC via Message Queues

True to the microkernel philosophy, the system avoids direct function calls between components. Instead, it relies on a robust, asynchronous message-passing mechanism:

* **Dedicated Event Queues:** Every running process (`IProcess`) has its own dedicated Go channel (`messageChan`), serving as a private, buffered event queue.
* **Kernel as a Dispatcher:** The Kernel acts purely as a central message dispatcher. It forwards an `IMessage` to the target process's queue without directly accessing its internal state, ensuring components are truly decoupled.
* **Concurrent and Isolated Execution:** Each process runs its own event loop in a separate goroutine. This design ensures that a slow, busy, or blocked component will not affect the Kernel or any other part of the system, leading to a highly resilient and stable multitasking environment.

---

## 3. Key Features

### a. Multitasking & TUI Window Manager

The kernel provides a fully-featured multitasking environment with a graphical TUI, managed by the `Render` server.

* **Concurrent Tasks**: Launch and run multiple applications (e.g., system monitors, control scripts) simultaneously.
* **Windowed Interface**: Every graphical task is rendered within its own distinct "window" on the terminal, complete with borders and a title caption.
* **Process Management**: Standard OS-like commands are provided: `ps`, `kill`, `killall`, and `fg`.

### b. Advanced Rendering Architecture: A Retained-Mode TUI Engine

The framework implements a sophisticated graphics architecture known as **"retained mode"** or display list rendering. This approach decouples an application's logical state from its physical representation, leading to significant performance and responsiveness gains, especially for dynamic TUIs.

**How It Works:**

1.  **Abstract Scene Description**: When a process receives a `paint` event, it does not draw pixels to a framebuffer. Instead, it builds a `DescriptiveSurface`, which is a lightweight **list of drawing commands** (e.g., "draw this text at x,y," "trace this line with these colors"). This list is a resolution-independent, abstract description of *what* the process wants to display.

2.  **The Render Server as an Interpreter**: The `Render Server` receives and **retains** this command list in memory. Its job is to interpret these commands and execute them to render the final bitmap into its own internal cache.

**The Key Advantage: Zero-Cost Windowing Operations**

The true power of this design becomes evident during UI operations like **resizing or moving a window**:

> **A new `paint` event is not required from the application.**

When a window is resized, the application process is not interrupted and does not need to re-run its rendering logic. The `Render Server` autonomously handles the entire operation by:

1.  Creating a new destination bitmap of the correct size.
2.  **Re-executing the exact same command list it has already stored** onto the new, resized surface.

This turns complex UI operations into simple, high-speed re-executions of a display list, providing enormous benefits:

* **Maximum Efficiency**: Avoids expensive re-computation of application logic (`OnPaint`) just for a layout change.
* **Instant Responsiveness**: Window resizing and moving feel fluid and instantaneous, as they are handled entirely by the server optimized for this task.
* **Clean Architecture**: Reinforces the separation of concerns. The process handles the **"what,"** while the Render Server handles the **"how."**

In essence, the framework provides the TUI equivalent of a **vector graphics engine**. The scene is described abstractly and can be re-rendered at any size without information loss and without needing to involve the application's core logic again.

### c. Real-time Introspection & Go Runtime Profiling

The included `stats` module transforms the shell into a powerful, live diagnostics tool for the Go runtime.

* **Live Runtime Stats**: Get instant snapshots of memory usage (`rt`), CPU status, and goroutine counts (`cpu`).
* **Integrated `pprof` Profiling**: Start and stop CPU profiling (`startcpuprofile`, `stopcpuprofile`) and generate heap profiles (`memprofile`) on the fly for any application built on the framework.

### d. Filesystem-like Command Hierarchy

* **Hierarchical Structure**: Commands are organized in a virtual filesystem, allowing for intuitive navigation with `cd`, `ls`, and `pwd`.
* **Tab Completion**: Rich auto-completion for commands and paths.

### e. Remote & Secure Access

* **SSH Server**: A built-in, full-featured SSH server supports both password and public-key authentication for secure remote access.
* **Telnet Server**: A Telnet server is included for simpler, unencrypted connections in trusted networks.

---

## 4. Example Applications & Starter Toolkit

The framework's power is demonstrated through a suite of built-in applications that can serve as examples or as a starter toolkit for a new project:

* **`system/`**: Essential shell utilities that interact with the kernel's core services.
* **`stats/`**: The live profiling and introspection toolkit for the Go runtime.
* **`games/`**: A collection of fully implemented, interactive games (`invaders`, `tetris`, `snake`) that showcase the rendering engine, multitasking, and event handling capabilities.

---

## 5. Dependencies

* Go Standard Library
* `golang.org/x/crypto/ssh` (for the SSH server)

---
*Module Author: Marcello (born 1974)*
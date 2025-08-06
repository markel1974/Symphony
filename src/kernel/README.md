# Symphony MicroKernel Module

## 1. Overview

The **Symphony MicroKernel Module** is an advanced, Go-based framework for creating powerful, interactive command-line shells. It is architected as a modular **microkernel** designed to serve as a real-time control, introspection, and live debugging interface for complex systems, such as the Symphony emulator.

The framework simulates a complete **multitasking, windowed TUI environment** directly within the terminal. This allows users to manage concurrent processes, interact with a hierarchical command system, and dynamically inspect and modify the host application's state.

More than just a shell toolkit, Symphony is a framework for building **resilient, transparent, and debuggable-by-design systems**. Its core philosophy is that a complex application must be fully introspectable and controllable at runtime. This turns the shell into a powerful live diagnostics and recovery console, empowering developers to move beyond post-mortem log analysis and perform interactive 'software surgery'. Using this approach, complex issues that are invisible to traditional logging and telemetry can be **diagnosed, confirmed, and even temporarily patched on a live system in minutes.**

---

## 2. Microkernel Architecture

The project fully embraces a **microkernel design philosophy**: a minimal, robust core provides essential services, while all other functionalities are delegated to isolated, user-space processes called "Tasks".

* **Central Kernel (`kernel/core/kernel.go`)**: The heart of the system. It manages the entire lifecycle of tasks, orchestrates I/O, and runs a central, asynchronous `eventLoop`. Its primary role is to act as a **message router**, decoupling the system's components.
* **Dynamic Message Handling**: The kernel uses a dynamic, self-populating `handlers` map. Each peripheral `Server` (e.g., for rendering, filesystem) registers the message types it can handle by implementing the `IServer` interface. When a server is added, the kernel automatically updates its routing table to forward the correct messages, eliminating the need for hardcoded dependencies.
* **Task Management (`kernel/process/process.go`)**: Every command or application runs as a concurrent `Task` (implementing `IProcess`), complete with its own PID, isolated context, and state.
* **The Shell as a User Process (`xsh`)**: As a testament to the purity of the design, the main shell (`xsh`) is not part of the kernel but is itself an application that runs as a regular "Task". On boot, the kernel simply launches the `xsh` process to provide the user interface.
* **Interface-Driven Design (`kernel/interfaces/`)**: The entire system is built upon a foundation of interfaces (`IProcess`, `ICommand`, `IRender`, `IFileSystem`), ensuring a clean separation of concerns and promoting high cohesion and low coupling.

### True Asynchronous IPC via Message Queues

True to the microkernel philosophy, the system avoids direct function calls between the kernel and its processes. Instead, it relies on a robust, asynchronous message-passing mechanism that ensures true process isolation:

* **Dedicated Event Queues**: Every running process (`IProcess`) is equipped with its own dedicated Go channel (`messageChan`), which serves as a private, buffered event queue.
* **Kernel as a Dispatcher**: The Kernel acts purely as a central message dispatcher. When an event occurs (e.g., user input, a timer tick), the Kernel forwards an `IMessage` to the target process's specific message queue. It does not directly access the process's internal state or methods.
* **Concurrent and Isolated Execution**: Each process runs its own event loop in a separate goroutine, continuously reading from its message queue. This design ensures that processes are fully isolated and truly concurrent. A slow, busy, or blocked process will not affect the Kernel or any other process in the system, leading to a highly resilient and stable multitasking environment.

---

## 3. Key Features

### a. Multitasking & TUI Window Manager

The shell is a fully-featured multitasking environment with a graphical TUI, managed by the `Render` server.

* **Concurrent Tasks**: Launch and run multiple applications (e.g., games, system monitors) simultaneously. The kernel's `AdaptiveTicker` ensures efficient scheduling of timed events for smooth, concurrent execution.
* **Windowed Interface**: Every graphical task is rendered within its own distinct "window" on the terminal, complete with borders and a title caption managed by the `Surface` object.
* **Dynamic Window Management**: A `WindowSelector` enables a special mode where the user can:
    * **Cycle through windows** using `Tab` and `q`.
    * **Move selected windows** around the screen using `w,a,s,d` or arrow keys.
    * **Resize selected windows** using `+` and `-` to scale their content.
* **Process Management**: Standard OS-like commands are provided to manage tasks: `ps` to list active processes, `kill` to terminate a specific PID, `killall` to terminate by name, and `fg` to bring a task to the foreground.

### b. Real-time Introspection & Profiling

The `stats` module transforms the shell into a powerful, live diagnostics tool for the Go runtime.

* **Live Runtime Stats**: Get instant snapshots of memory usage (`rt`), CPU status, and goroutine counts (`cpu`).
* **Integrated `pprof` Profiling**:
    * **CPU Profiling**: Start and stop CPU profiling on the fly (`startcpuprofile`, `stopcpuprofile`) and save the output for analysis with `go tool pprof`.
    * **Memory Profiling**: Generate heap profiles (`memprofile`) to debug memory leaks and allocation patterns.
* **Graphical Real-time Monitoring**: The `rtplot` command launches an interactive, scrolling graph to monitor memory metrics (like heap allocation or GC cycles) over time, using the built-in `plotter`.

### c. Advanced Rendering Engine

* **Surface-Based Rendering**: An `ISurface` abstraction provides a drawable area for each task, handling scaling, offsets, and window decorations automatically.
* **2D Sprite & Matrix Engine**: A powerful engine (`matrix`) for complex TUI graphics, featuring:
    * `Entity` objects with physics properties (mass, velocity) and multi-frame, colored `Sprite` rendering.
    * An `AABBTree` for highly efficient 2D collision detection, showcased in the `invaders` game.

### d. Filesystem-like Command Hierarchy

* **Hierarchical Structure**: Commands are organized in a virtual filesystem, allowing for intuitive navigation with `cd`, `ls`, and `pwd`, all managed by the `FileSystem` server.
* **Tab Completion**: Rich auto-completion for commands and paths, using Levenshtein distance for fuzzy suggestions.
* **Persistent History**: Command history is managed by the `HistoryHandler` and can be configured to be saved across sessions.

### e. Remote & Secure Access

* **SSH Server**: Full-featured SSH server supporting both password and public-key authentication for secure remote access.
* **Telnet Server**: A Telnet server is included for simpler, unencrypted connections.

---

## 4. Built-in Applications Showcase

The framework's power is demonstrated through a suite of built-in applications:

* **`system/`**: Essential shell utilities that interact with the kernel's core services.
* **`stats/`**: The live profiling and introspection toolkit.
* **`games/`**: Fully implemented, interactive games that showcase the rendering engine, multitasking, and event handling capabilities:
    * **`invaders`**: A feature-rich clone with animated sprites, destructible barricades, and particle explosions.
    * **`tetris`**: A complete Tetris game with scoring, levels, and a next-piece preview.
    * **`snake`**: The classic Snake game.

---

## 5. Dependencies

* Go Standard Library
* `golang.org/x/crypto/ssh` (for the SSH server)

---
*Module Author: Marcello (born 1974)*
*Primary Context: Symphony Emulator*
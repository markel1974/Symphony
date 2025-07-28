# Symphony Shell Module

## Overview

This module implements an advanced interactive shell, written entirely in Go. It is designed to be highly modular, extensible, and decoupled, primarily serving as the **control, introspection, and live debugging interface** for the [Symphony](link-to-symphony-project-if-exists) emulation framework (or other complex Go systems).

The shell allows interaction with the host application as if navigating a filesystem, executing commands, managing internal processes (tasks), and inspecting/modifying the application's state in real-time.

## Key Features

* **Remote Access:** Integrated servers for **SSH** and **Telnet** connections.
* **Hierarchical Command System:** Defines and manages commands organized in a tree structure (`cli`, `interfaces.ICommand`).
* **Filesystem Navigation:** Simulates filesystem navigation for the command structure, supporting `cd`, absolute/relative paths (`interfaces.IFileSystem`).
* **Task/Process Management:** An internal "kernel" (`context/kernel`) manages the lifecycle of concurrent tasks (`interfaces.ITask`) with PIDs, foreground/background states, and basic inter-task communication.
* **Live Introspection:** Allows commands to inspect (and potentially modify) the internal state of other tasks/components of the host application.
* **Advanced Rendering:** Terminal rendering system (`render`, `interfaces.ISurface`) supporting colors, basic TUI (Text User Interface), plotting (`render/plotter`), and textual graphics with sprites/matrices (`render/matrix`).
* **Task Timers:** Timing mechanism (`adaptiveticker`) to execute code at regular intervals within tasks.
* **Command History:** Manages command history with navigation (up/down arrows) and optional persistence (`shell/history.go`).
* **Command Completion:** Supports suggestion generation and completion via Tab key.
* **Authentication:** Module for user authentication (`interfaces.IAuthenticator`, `authenticator/simple.go`).
* **Extensibility:** Easily extensible by adding new commands (`cli.Command`) or entire "applications" (`apps/`).

## Architecture

The module is designed following modern software development principles:

* **Interface-Driven:** The core design relies on the `interfaces/` package, defining contracts between components, ensuring low coupling and high testability.
* **Separation of Concerns (SoC):** Each package has a well-defined purpose (see below).
* **Modularity:** Components are largely independent and interact via defined interfaces.
* **Idiomatic Go:** Leverages Go features like interfaces, goroutines, and channels (e.g., in the `kernel` and `adaptiveticker`).

## Core Components (Packages)

* `interfaces/`: Definitions of key APIs/interfaces (ICommand, ITask, IRender, ITerminal, IFileSystem, IAuthenticator, etc.).
* `context/`: Contains the `Kernel` (main orchestrator, event and task manager), `IFileSystem` implementation, and the `Render` manager.
* `cli/`: System for defining, hierarchically organizing, and parsing commands.
* `render/`: Rendering engine onto `ISurface`, includes the `plotter` and `matrix` utilities (for TUI/games).
* `shell/`: Main logic for the interactive shell loop, user input handling, and history management.
* `terminal/`: Abstraction for interaction with the physical/virtual terminal (VT100 implementation included).
* `ssh/`, `telnet/`: Server implementations for remote access.
* `adaptiveticker/`: Custom handler for periodic timers associated with tasks.
* `authenticator/`: Logic for authentication (simple implementation included).
* `apps/`: Collections of example commands/applications using the framework:
    * `apps/core`: Basic shell commands (cd, ls, ps, kill, help, history, etc.).
    * `apps/stats`: Commands for displaying runtime statistics (memory, CPU, GC).
    * `apps/games`: Examples of text-based games (Snake, Tetris, Invaders) demonstrating rendering and event handling capabilities.

## Typical Integration

An application using this module typically:
1.  Creates a root command structure using the `cli` package.
2.  Instantiates an `IAuthenticator`.
3.  Instantiates a server (`shell.NewServer`) choosing between SSH or Telnet, passing the authenticator, port, and root command structure.
4.  Configures the desired prompt (`server.SetPrompt`).
5.  Starts the server (`server.Start` or `server.AsyncStart`).

## Dependencies

* Go Standard Library
* `golang.org/x/crypto/ssh` (for the SSH server)

---
*Module Author: Marcello (born 1974)*
*Primary Context: [Symphony](link-to-symphony-project-if-exists)*
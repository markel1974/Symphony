# The Symphony Framework: A Microkernel for "Transparent" Systems

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)

This is not just another kernel. It is an **advanced, Go-based microkernel framework for building complex systems that are transparent, inspectable, and manipulable in real-time by design.**

Its core philosophy is an answer to one of the most critical challenges in modern software engineering: what do you do when a complex service in production begins to misbehave, and your logs and telemetry fall silent?

Symphony's answer is to **stop building black boxes**. This framework provides the foundation to create applications equipped from day one with a built-in "control room and operating theater," accessible via a secure remote console. It transforms debugging from a post-mortem forensic analysis into **live, open-heart surgery** on a running system.

---

## 1. Microkernel Architecture: Stability and Modularity

The project fully embraces a microkernel design philosophy: a minimal, robust core orchestrates essential services, while all other functionalities are delegated to isolated processes and servers.

* **Central Kernel (`kernel/core/kernel.go`)**: The heart of the system. It acts as an asynchronous **message router**, managing the process lifecycle and orchestrating communication without knowing the implementation details of the components it connects.

* **Specialized Servers**: Complex functionalities like rendering (`servers/render/`) or the virtual filesystem (`servers/file_system/`) are implemented as long-running servers. Each server runs in its own event loop and communicates exclusively via messages, ensuring that a slow or blocked server cannot destabilize the kernel or other processes.

* **User-Space Processes**: Every application, including the shell itself (`xsh`), is an isolated user-space process with a clean separation from its kernel-space representation (`KernelProcess`), ensuring stability and security.

---

## 2. The Power of Live Introspection

Symphony's killer application is the ability to inspect and interact with the internal state of any component in real-time. This is made possible by a powerful reflection system and a code generator (`component/exporter.go`).

* **Properties**: Any field within a component's struct can be exposed as a "Property." This allows an engineer to read internal values (e.g., the contents of a cache, the state of a connection pool, the registers of an emulated CPU) with a simple `get` command from the shell.
* **Commands**: Any Go method can be exposed as a "Command," invokable from the shell with `exec`. This allows for activating internal functions, forcing resource cleanups, or altering the behavior of a service on the fly, without requiring a restart or redeployment.

The included **MOS 6581 SID chip emulator** is the perfect example of this philosophy in action, where every register and internal state of the sound chip is inspectable and manipulable in real-time.

---

## 3. An Educational Laboratory for Operating Systems 🎓

Beyond its professional applications, Symphony is an **outstanding educational tool**. Its clean, modern implementation in Go makes it the perfect bridge between textbook theory and the complexity of real-world kernels like MINIX.

* **Concepts Made Tangible**: Instead of just reading about "message passing," students can open `kernel/core/kernel.go` and see the `messageChan` and the `eventLoop` in action.
* **Interactive Experimentation**: The interactive shell transforms learning into an experiment. A student can launch multiple processes (`ps`), terminate them (`kill`), and observe component isolation in real-time, gaining a deep, practical understanding of microkernel principles.

---

## 4. Other Key Features

* **Multitasking TUI**: A responsive and efficient "retained-mode" TUI window manager.
* **Secure Remote Access**: A built-in SSH server with support for both password and public-key authentication.
* **Filesystem-like Command Hierarchy**: A virtual filesystem for commands, with navigation (`cd`, `ls`) and rich autocompletion.

---
*Module Author: Marcello (born 1974)*
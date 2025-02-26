# Package quartz

The `quartz` package provides a high-performance clock management system designed for the emulation of hardware components that rely on precise timing, such as CPUs, peripherals, and other cycle-dependent entities.

## Overview

The `quartz` package acts as a central timing mechanism for emulated environments, offering:

- **Clock Generation**: Provides clock signals for components requiring cycle-based operations.
- **Synchronization**: Ensures all emulated components operate in sync with a shared clock source.
- **Tick Management**: Tracks clock cycles and advances components based on configurable intervals.
- **Event Scheduling**: Allows components to register and handle events triggered after a specific number of cycles.

## Features

- **Flexible Clock Configuration**: Supports setting different clock frequencies to match the requirements of target hardware.
- **Cycle-Accurate Timing**: Ensures precise clock tick handling to maintain synchronization across multiple components.
- **Event Handling**: Provides an event queue system for scheduling tasks in specific ticks or cycles.
- **Pause/Resume Functionality**: Allows for stopping and restarting the clock without losing state.

## Use Cases

The `quartz` package is designed to be used in emulation systems where components such as CPUs, CIAs, VIC-II, or audio processors need to operate together under precise time constraints.

For example, in the Commodore 64, `quartz` can be used to provide clock signals for:

- **CPU Emulation**: Advancing the processor based on clock ticks.
- **Peripheral Updates**: Ensuring that devices like the 6526 (CIA) or the SID chip operate in sync with the system clock.
- **Interrupt Scheduling**: Triggering time-based interrupts accurately.

## Package Structure

The package consists of the following key components:

- `clock.go`: Implements the core `Clock` struct and its functionality for tick management and synchronization.
- `scheduler.go`: Manages the event queue and provides utilities for scheduling tasks based on future cycles.
- `quartz_test.go`: Includes unit tests to ensure timing accuracy and proper operation of the clock and event scheduler.

## Key Components

### Clock

The heart of the package is the `Clock` struct, which:

- Tracks the number of elapsed cycles (ticks).
- Allows components to attach and synchronize their cycle-based operations.
- Supports configurable clock frequencies for different emulated hardware.

### Scheduler

The `Scheduler` provides:

- An event-driven system where tasks can be scheduled to run after a specific number of clock ticks.
- A priority queue ensuring tasks are triggered in the correct order, even if multiple tasks are scheduled for the same cycle.

### Example Usage

Below is a high-level example of how the `quartz` package might be used:

```go
package main

import "github.com/yourproject/quartz"

func main() {
    // Create a new clock instance
    clock := quartz.NewClock(1000000) // 1 MHz clock
    
    // Register a callback for a scheduled event
    clock.ScheduleEvent(500, func() {
        fmt.Println("Event triggered after 500 cycles")
    })
    
    // Advance the clock
    clock.Tick(500)
}
```

## Dependencies

This package does not rely on significant external libraries, making it lightweight and reliable in core applications.

## Notes

- The `quartz` package is designed for systems requiring **high timing precision**.
- It is highly optimized for integration in hardware emulators and other systems dependent on synchronized simulations.

## TODO

- Add support for nested clocks (e.g., sub-clocks or derived clocks for specific components).
- Optimize event queue performance for systems with a high frequency of scheduled events.
- Expand unit tests to cover edge cases and stress tests.

## Contributing

To contribute to the `quartz` package:

1. Fork the repository.
2. Create a branch for your changes.
3. Submit a pull request with a detailed explanation of your changes or improvements.

For questions or support, feel free to open an issue.

## License

This project is released under the [Apache 2.0 license](https://opensource.org/licenses/Apache-2.0).
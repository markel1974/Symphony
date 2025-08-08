package interfaces

// IUserRouter defines the interface for posting messages *to* a user-space
// process or a server. It acts as the "inbox" for any running component,
// allowing the kernel to deliver events, responses, or notifications.
type IUserRouter interface {
	// PostUserMessage queues a message for the component to process in its
	// own event loop.
	PostUserMessage(msg IMessage)
}

// IKernelRouter defines the interface for a user-space process to post
// messages *to* the kernel. This is the primary channel for applications
// to request system services (e.g., creating a timer, writing to the screen,
// or exiting).
//
// The sender's PID is implicitly handled by the kernel-space wrapper
// (`KernelProcess`) that mediates the call, so it is not needed as a parameter.
type IKernelRouter interface {
	// PostKernelMessage queues a message for the kernel's central event loop
	// to process.
	PostKernelMessage(msg IMessage)
}

// IKernelServerRouter defines the interface used by kernel-space servers (and the
// kernel itself) to post messages to a specific destination process via the kernel.
//
// Unlike IKernelRouter, this interface requires a destination `pid` because
// servers often need to route messages to specific processes. This is used for two
// primary purposes:
//
//  1. **Replying to Requests:** To send a response back to the specific process
//     that originated a request.
//  2. **Asynchronous Notifications:** To proactively send notifications about
//     system-wide events, such as a change in the filesystem or an invalidated
//     surface, to interested processes.
type IKernelServerRouter interface {
	// PostKernelServerMessage queues a message for the kernel to deliver to the
	// process identified by `pid`.
	PostKernelServerMessage(pid int, msg IMessage)
}

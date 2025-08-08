package interfaces

type IUserRouter interface {
	PostUserMessage(msg IMessage)
}

type IKernelRouter interface {
	PostKernelMessage(msg IMessage)
}

type IKernelServerRouter interface {
	PostKernelServerMessage(pid int, msg IMessage)
}

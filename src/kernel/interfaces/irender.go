package interfaces

type IRender interface {
	IServer

	CallGetScreenSize(process IRouter) (int, int)
}

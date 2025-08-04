package interfaces

type IServer interface {
	Register() []MessageType

	Start()

	PostMessage(m IMessage)

	HandleMessage(msg IMessage)

	NotifyProcessCreation(desc *ProcessDescription)

	NotifyProcessTermination(desc *ProcessDescription)

	NotifyProcessForeground(desc *ProcessDescription)
}

type Server struct {
	IServer
	messageChan chan IMessage
}

// Start begins the process by setting its state to running and initiating its event loop asynchronously.
func (t *Server) Start() {
	c := make(chan bool)
	t.eventLoop(c)
	_ = <-c
}

func (t *Server) PostMessage(m IMessage) {
	t.messageChan <- m
}

// evenLoop continuously listens on the message channel and processes incoming messages until a quit message is received.
func (t *Server) eventLoop(r chan bool) {
	go func() {
		r <- true
		for {
			select {
			case m, ok := <-t.messageChan:
				if !ok {
					return
				}
				if m.GetType() == MessageTypeQuit {
					close(t.messageChan)
					return
				} else {
					t.HandleMessage(m)
				}
			}
		}
	}()
}

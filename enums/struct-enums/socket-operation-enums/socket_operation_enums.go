package socket_operation_enums

var SocketOperation = newSocketOperation()

func newSocketOperation() *socketOperation {
	return &socketOperation{
		CHAT:        "CHAT",
	}
}

type socketOperation struct {
	CHAT        string
}

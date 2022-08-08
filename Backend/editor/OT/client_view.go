package editor

import (
	"github.com/gorilla/websocket"
)

// clientView is the embodiment of all data relating to a clientView connection
// it is mostly managed by the documentServer
// quick sidenote:
//		since sendOp and sendAcknowledgement are bounded we can actually deadlock the system
//		by filling them up :P, the documentServer avoids this with a little hack, it spins a goroutine
//		to publish updates from the clientView within the pipe
type clientView struct {
	socket *websocket.Conn

	sendOp              chan op
	sendAcknowledgement chan empty
	sendTerminateSignal chan empty
}

// models an operation a client can propagte to the documentServer
type op struct {
	data     string
	location int
	// todo: add more metadata
}

type opRequest struct {
	requestType string
	operation   op
}

func newClient(socket *websocket.Conn) *clientView {
	return &clientView{
		socket:              socket,
		sendOp:              make(chan op),
		sendAcknowledgement: make(chan empty),
		sendTerminateSignal: make(chan empty),
	}
}

// run starts the client's spinny loop
// in this loop the clientView listens for updates and pushes
// them down the websocket, it also pulls stuff up the websocket
// the documentServer will use the appropriate channels to communicate
// updates to the client, namely: sendOp and sendAcknowledgement
func (c *clientView) run(serverPipe pipe, terminatePipe alertLeaving) {
	for {
		select {
		case <-c.sendOp:
			// push the operation down the websocket
			break

		case <-c.sendAcknowledgement:
			// push the acknowlegement down the websocket
			break

		case <-c.sendTerminateSignal:
			// looks like we've been told to terminate by the documentServer
			// propagate this to the client and close this connection
			c.socket.Close()
			return

		default:
			request := opRequest{}

			err := c.socket.ReadJSON(&request)
			if err != nil {
				// todo: push a terminate signal to the client, also tell the server we're leaving
				terminatePipe()
				c.socket.Close()
			}

			// push the update to the documentServer
			serverPipe(request.operation)
		}
	}
}

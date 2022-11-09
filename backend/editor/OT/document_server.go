package editor

import (
	"sync"

	"cms.csesoc.unsw.edu.au/editor/OT/operations"
	"cms.csesoc.unsw.edu.au/pkg/cmsjson"
	"github.com/google/uuid"
)

type documentServer struct {
	// todo: change to whatever data structure is being used for
	// todo: stop the clientView map from growing too large using some compaction
	// strategy or a more appropriate ds
	// state management
	ID        uuid.UUID
	state     cmsjson.AstNode
	stateLock sync.Mutex

	clients     map[int]*clientState
	clientsLock sync.Mutex

	operationHistory []operations.Operation
}

type clientState struct {
	*clientView
	canSendOps bool
}

func newDocumentServer() *documentServer {
	// ideally state shouldn't be a string due to its immutability
	// any update requires the allocation + copy of a new string in memory
	return &documentServer{
		state:       nil,
		stateLock:   sync.Mutex{},
		clients:     make(map[int]*clientState),
		clientsLock: sync.Mutex{},
	}
}

// a pipe is a closure that the clientView can use to communicate
// with the server, it wraps its internal clientView ID for security reasons
type pipe = func(op operations.Operation)

// alertLeaving is like a pipe except a client uses it to tell a document
// that it is leaving
type alertLeaving = func()

// connectClient connects a clientView to a documentServer and returns a one way pipe
// it can use for communication with the documentServer
// TODO: synchronise this properly
func (s *documentServer) connectClient(c *clientView) (pipe, alertLeaving) {
	// register this clientView
	s.clientsLock.Lock()
	clientID := len(s.clients)
	s.clients[clientID] = &clientState{
		clientView: c,
		canSendOps: true,
	}
	s.clientsLock.Unlock()

	// we need to create a new worker for this clientView too
	workerHandle := make(chan func())
	killHandle := make(chan empty)
	go createAndStartWorker(workerHandle, killHandle)

	// finally build a comm pipe for this clientView
	return s.buildClientPipe(clientID, workerHandle, killHandle), s.buildAlertLeavingSignal(clientID, killHandle)
}

// disconnectClient removes a client from a document server
func (s *documentServer) disconnectClient(clientID int) {
	s.clientsLock.Lock()
	if _, ok := s.clients[clientID]; !ok {
		panic("Trying to disconnect non-existent client")
	}

	delete(s.clients, clientID)
	s.clientsLock.Unlock()

	// if we have no more connected clients it may be time to terminate ourselves
	GetDocumentServerFactoryInstance().closeDocumentServer(s.ID)
}

// buildClientPipe is a function that returns the "pipe" for a clientView
// this pipe contains all the necessary code that the clientView needs to communicate with the documentServer
// when the clientView wishes to send data to the documentServer they simply just call this pipe with the operation
func (s *documentServer) buildClientPipe(clientID int, workerWorkHandle chan func(), workerKillHandle chan empty) func(operations.Operation) {
	return func(op operations.Operation) {
		// this could also just be captured from the outer func
		clientState := s.clients[clientID]
		thisClient := clientState.clientView
		if !clientState.canSendOps {
			// terminate this clientView
			// this is the only thing we can do in order to enforce
			// consistency across all clients
			s.disconnectClient(clientID)
			go func() { clientState.sendTerminateSignal <- empty{} }()
			go func() { workerKillHandle <- empty{} }()
			return
		}

		// to deal with this incoming operation we need to push
		// data to the worker assigned to this clientView
		workerWorkHandle <- func() {
			defer func() {
				clientState.canSendOps = true
				thisClient.sendAcknowledgement <- empty{}
			}()

			clientState.canSendOps = false

			// apply op to clientView states
			s.stateLock.Lock()

			// TODO: apply operation
			//	- Blockers:
			//		- Gary's TLB
			//		- Updated Traversal

			// apply the operation locally and log the new operation
			transformedOperation := s.transformOperation(op)
			s.operationHistory = append(s.operationHistory, transformedOperation)

			s.stateLock.Unlock()

			if transformedOperation.IsNoOp {
				return
			}

			// propagate updates to all connected clients except this one
			// if we send it to this clientView then we may deadlock the server and clientView
			s.clientsLock.Lock()
			for id, connectedClient := range s.clients {
				if id == clientID {
					continue
				}

				// push update
				connectedClient.sendOp <- transformedOperation
			}
			s.clientsLock.Unlock()
		}
	}
}

// transformOperation transforms an incoming client operation against the history of applied server operations
// 		note: the baseOpIndex indicates what operation to start applying against
func (s *documentServer) transformOperation(incomingOp operations.Operation) operations.Operation {
	for _, op := range s.operationHistory[incomingOp.AcknowledgedServerOps:] {
		_, incomingOp = operations.TransformPipeline(incomingOp, op)
	}

	return incomingOp
}

// buildAlertLeavingSignal builds a leaving signal for the client view
// to use when it wants to tell the document server that it is leaving
func (s *documentServer) buildAlertLeavingSignal(clientID int, workerKillHandle chan empty) func() {
	// go doesn't have currying :(
	return func() {
		workerKillHandle <- empty{}
		s.disconnectClient(clientID)
	}
}

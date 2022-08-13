package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cms.csesoc.unsw.edu.au/internal/logger"
	"cms.csesoc.unsw.edu.au/internal/session"
)

// Basic organization of a response we will receive from the API
type (
	empty struct{}

	// APIResponse is the public response type that is marshalled and presented to consumers of the API
	APIResponse[V any] struct {
		Status   int
		Message  string
		Response V
	}

	// handlerResponse is a special response type only returned by HTTP Handlers
	handlerResponse[V any] struct {
		Status   int
		Response V
	}
)

// This file contains a series of types defined to make writing http handlers
// a bit easier and less messy

type (
	handler[T, V any] struct {
		FormType    string
		Handler     func(form T, dependencyFactory DependencyFactory) (response handlerResponse[V])
		IsMultipart bool
	}

	// authenticatedHandler is basically a regular http handler the only difference is that
	// they can only be accessed by an authenticated client
	authenticatedHandler[T, V any] handler[T, V]
)

// ServeHTTP is an overloaded implementation of method on the http.HttpHandler interface
func (fn handler[T, V]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Determine what type of form parser to use first
	parser := getParser(fn)
	parsedForm := new(T)

	if parseStatus := parser(r, fn.FormType, parsedForm); parseStatus != http.StatusOK {
		writeResponse(w, handlerResponse[empty]{
			Status:   parseStatus,
			Response: empty{},
		})

		return
	}

	// construct a dependency factory for this request, which implies instantiating a logger
	logger := buildLogger(r.Method, r.URL.Path)
	dependencyFactory := DependencyProvider{Log: logger}
	response := fn.Handler(*parsedForm, dependencyFactory)

	// Record and write out any useful information
	writeResponse(w, response)
	logResponse(logger, response)
	logger.Close()
}

// ServeHTTP is an overloaded implementation of method on the http.HttpHandler interface, the constraint for the authenticateHandler
// is that it wraps the target handler up in an authentication check
func (fn authenticatedHandler[T, V]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok, err := session.IsAuthenticated(w, r); !ok || err != nil {
		writeResponse(w, handlerResponse[empty]{
			Status:   http.StatusUnauthorized,
			Response: empty{},
		})

		return
	}

	// parse request over to main handler
	handler[T, V](fn).ServeHTTP(w, r)
}

// getMessageFromStatus fetches the message corresponding to a given status code
func getMessageFromStatus(statusCode int) string {
	statusMappings := map[int]string{
		http.StatusBadRequest:          "missing parameters (check documentation)",
		http.StatusMethodNotAllowed:    "invalid method",
		http.StatusNotFound:            "unable to find requested object",
		http.StatusNotAcceptable:       "unable to preform requested operation",
		http.StatusInternalServerError: "somethings wrong I can feel it",
		http.StatusOK:                  "ok",
	}

	if message, ok := statusMappings[statusCode]; ok {
		return message
	}

	return "..."
}

// writeResponse is a small helper function to write out a received handler response to the response writer
func writeResponse[V any](dest http.ResponseWriter, response handlerResponse[V]) {
	out := APIResponse[V]{
		Status:   response.Status,
		Response: response.Response,
		Message:  getMessageFromStatus(response.Status),
	}

	dest.Header().Set("Content-Type", "application/json")
	re, _ := json.Marshal(out)
	dest.Write(re)
}

// buildLogger instantiates a logger instance given a method / endpoint of the handler
func buildLogger(method string, endpoint string) *logger.Log {
	return logger.OpenLog(fmt.Sprintf("Handling http %s request to %s", method, endpoint))
}

// logResponse just logs a handler response under the provided log
func logResponse[V any](logger *logger.Log, response handlerResponse[V]) {
	switch response.Status {
	case http.StatusOK:
		logger.Write([]byte("successfully handled request"))
	default:
		logger.Write([]byte(fmt.Sprintf("failed to handle request! status: %d \nresponse %v", response.Status, response.Response)))
	}
}

// formParser is a type indicating a valid form parser (see below)
type formParser = func(*http.Request, string, interface{}) int

// getParser fetches the required parser for a specific handler configuration
func getParser[T, V any](config handler[T, V]) formParser {
	if config.IsMultipart {
		return ParseMultiPartFormToSchema
	}

	return ParseParamsToSchema
}

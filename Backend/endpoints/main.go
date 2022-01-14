package endpoints

import (
	"encoding/json"
	"net/http"

	"cms.csesoc.unsw.edu.au/internal/session"
)

// Basic organisation of a response we will receieve from the API
type Response struct {
	Status  int
	Message string
	Reponse interface{}
}

// This file contains a series of types defined to make writing http handlers
// a bit easier and less messy
type handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

// authenticated handler is basically a regular http handler the only difference is that
// they can only be accessed by an authenticated client
type authenticatedHandler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

// impl of http Handler interface so that it can serve http requests
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if status, resp, err := fn(w, r); err == nil {
		if status != http.StatusMovedPermanently {
			w.WriteHeader(status)
		}
		out := Response{
			Status:  status,
			Reponse: resp,
		}

		// advanced machine learning model that predicts a http response
		// message given a request code :O
		switch status {
		case http.StatusBadRequest:
			out.Message = "missing parameters (check documentation)"
			break
		case http.StatusMethodNotAllowed:
			out.Message = "invalid method"
			break
		case http.StatusNotFound:
			out.Message = "unable to find requested object"
			break
		case http.StatusNotAcceptable:
			out.Message = "unable to preform requested operation"
			break
		case http.StatusInternalServerError:
			out.Message = "somethings wrong I can feel it"
			break
		case http.StatusOK:
			out.Message = "ok"
			break
		default:
			out.Message = "..."
			break
		}
		re, _ := json.Marshal(out)
		w.Write(re)
	} else {
		out, _ := json.Marshal(Response{
			Status:  http.StatusInternalServerError,
			Message: "something went wrong",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
	}
}

// authentiacated handler impl
func (fn authenticatedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check authentication
	if ok, err := session.IsAuthenticated(w, r); !ok || err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		out, _ := json.Marshal(Response{
			Status:  http.StatusInternalServerError,
			Message: "unauthorised",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	// parse request over to main handler
	handler(fn).ServeHTTP(w, r)
}

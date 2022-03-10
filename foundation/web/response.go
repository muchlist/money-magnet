package web

import (
	"fmt"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, status int, message interface{}) {
	env := Envelope{"error": message}
	err := WriteJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	ErrorResponse(w, http.StatusMethodNotAllowed, message)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, http.StatusInternalServerError, message)
}

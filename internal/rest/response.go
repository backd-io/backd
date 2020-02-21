package rest

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ResponseErr is the Response where a straight error cannot be assigned. Simplified for if err != nil
func ResponseErr(w http.ResponseWriter, err error) int {
	return Response(w, nil, err, http.StatusOK, "")
}

// Response is a default func to return data
func Response(w http.ResponseWriter, data interface{}, err error, desiredStatus int, location string) int {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Accept", "application/json")

	switch err {
	case nil:
		if location != "" {
			w.Header().Set("Location", location)
		}
		w.WriteHeader(desiredStatus)
		if data != nil {
			json.NewEncoder(w).Encode(data)
		}
		return desiredStatus
	case ErrConflict:
		return ErrorResponse(w, http.StatusConflict, "")
	case ErrBadRequest:
		return ErrorResponse(w, http.StatusBadRequest, "")
	case ErrUnauthorized:
		return ErrorResponse(w, http.StatusUnauthorized, "")
	}
	fmt.Printf("Unknown err: type: %T; value: %q\n", err, err)
	return ErrorResponse(w, http.StatusInternalServerError, err.Error())
}

// NotFound is a generic 404 response that will be returned if the router cannot
// match a route
func NotFound(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusNotFound, "")
}

// NotAllowed is a generic 405 response that will be returned if the router can match the method
func NotAllowed(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusMethodNotAllowed, "")
}

// BadRequest is a generic 400 response
func BadRequest(w http.ResponseWriter, r *http.Request, reason string) {
	ErrorResponse(w, http.StatusBadRequest, reason)
}

// Unauthorized is a generic 401 response
func Unauthorized(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusUnauthorized, "")
}

// ErrorResponse returns a formatted error
func ErrorResponse(w http.ResponseWriter, status int, reason string) int {
	w.WriteHeader(status)
	var err HTTPErrorStatus
	err.Code = status
	err.Message = http.StatusText(status)
	if reason != "" {
		err.Reason = reason
	}
	json.NewEncoder(w).Encode(err)
	return status
}

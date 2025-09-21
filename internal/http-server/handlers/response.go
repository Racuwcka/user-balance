package handlers

import (
	"bytes"
	"net/http"
)

func InvalidValidation(w http.ResponseWriter, err error) {
	GetErrorResponse(w, "invalid validation data", err, http.StatusUnprocessableEntity)
	return
}

func InvalidRequest(w http.ResponseWriter, err error) {
	GetErrorResponse(w, "invalid request body", err, http.StatusBadRequest)
	return
}

func InvalidResponse(w http.ResponseWriter, err error) {
	GetErrorResponse(w, "invalid create response", err, http.StatusInternalServerError)
	return
}

func GetErrorResponse(w http.ResponseWriter, errName string, err error, statusCode int) {
	w.WriteHeader(statusCode)
	buf := bytes.NewBufferString(errName)
	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteString("\n")
	_, _ = w.Write(buf.Bytes())
}

func GetSuccessResponse(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

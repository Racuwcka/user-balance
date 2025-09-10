package handlers

import (
	"bytes"
	"net/http"
)

func InvalidValidation(w http.ResponseWriter, err error, statusCode int) {
	GetErrorResponse(w, "invalid validation data", err, statusCode)
}

func InvalidRequest(w http.ResponseWriter, err error, statusCode int) {
	GetErrorResponse(w, "invalid request body", err, statusCode)
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

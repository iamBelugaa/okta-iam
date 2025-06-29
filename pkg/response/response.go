package response

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data,omitempty"`
	Error      any    `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
	Fields  any    `json:"fields,omitempty"`
}

func NewAPIError(msg string, fields any) *apiError {
	return &apiError{Message: msg, Fields: fields}
}

func Success(w http.ResponseWriter, code int, msg string, data any) {
	response := apiResponse{Success: true, Data: data, StatusCode: code, Message: msg, Error: nil}
	respond(w, code, response)
}

func Error(w http.ResponseWriter, code int, err *apiError) {
	response := apiResponse{Success: false, StatusCode: code, Error: err, Message: err.Message}
	respond(w, response.StatusCode, response)
}

func respond(w http.ResponseWriter, statusCode int, data apiResponse) {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

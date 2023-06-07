package helper

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

type Params struct {
	Error          error
	ResponseWriter http.ResponseWriter
	CustomMessage  string
}

// GetError : This is helper function to prepare error model.
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func GetError(p Params) {

	// log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: p.Error.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	if len(p.CustomMessage) > 0 {
		response.ErrorMessage = p.CustomMessage
	}

	message, _ := json.Marshal(response)

	// w.WriteHeader(response.StatusCode)
	p.ResponseWriter.Write(message)
}

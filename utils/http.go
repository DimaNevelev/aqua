package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type RouteProvider interface {
	RegisterRoutes(chi.Router)
}

type HTTPResponse struct {
	Code  int         `json:"code"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func SendHTTPResponse(writer http.ResponseWriter, data interface{}) {
	response := &HTTPResponse{
		Code: http.StatusOK,
		Data: data,
	}
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		SendHTTPError(writer, err, http.StatusInternalServerError)
	}
}

func SendHTTPError(writer http.ResponseWriter, err error, code int) {
	writer.WriteHeader(code)
	writer.Header().Set("X-Error", err.Error())

	response := &HTTPResponse{
		Code:  code,
		Error: err.Error(),
	}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

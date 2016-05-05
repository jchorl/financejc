package util

import (
	"encoding/json"
	"net/http"
)

type errorStruct struct {
	Message string `json:"message"`
}

func WriteJSONResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		WriteJSONError(w, err)
	}
}

func WriteJSONError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errStruct := errorStruct{
		Message: err.Error(),
	}
	WriteJSONResponse(w, errStruct)
}

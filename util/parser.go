package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
)

func ParseRequestBody[T interface{}](w http.ResponseWriter, r *http.Request) (res T) {
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	return
}

func EncodeSuccessResponse[T interface{}](w http.ResponseWriter, res T, status int) {
	w.WriteHeader(status)
	resp := dto.NewSuccessResponse[T](res)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to create response json: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func EncodeErrorResponse(w http.ResponseWriter, errorMessage string, status int) {
	w.WriteHeader(status)
	resp := dto.NewErrorResponse(errorMessage, status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to create error response json: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

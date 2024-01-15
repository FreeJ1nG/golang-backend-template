package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRequestBody[T interface{}](w http.ResponseWriter, r *http.Request) (res T) {
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	return
}

func EncodeResponse[T interface{}](w http.ResponseWriter, res T, status int) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to create response json: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

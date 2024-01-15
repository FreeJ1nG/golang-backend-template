package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRequestBody[T interface{}](r *http.Request) (res T, err error) {
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		err = fmt.Errorf("unable to parse request body: %s", err.Error())
		return
	}
	return
}

func EncodeResponse[T interface{}](w http.ResponseWriter, res T, status int) (err error) {
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		err = fmt.Errorf("unable to create response json: %s", err.Error())
		return
	}
	return
}

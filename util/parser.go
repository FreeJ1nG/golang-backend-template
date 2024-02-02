package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/pagination"
)

func ParseRequestBody[T interface{}](w http.ResponseWriter, r *http.Request) (res T) {
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse request body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	return
}

func EncodeSuccessResponse[T interface{}](w http.ResponseWriter, res T, status int, metadata *pagination.Metadata) {
	w.WriteHeader(status)
	resp := dto.NewSuccessResponse[T](res, metadata)
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

func ParseStringToInt(s string) (res int, status int, err error) {
	status = http.StatusOK
	res, err = strconv.Atoi(s)
	if err != nil {
		err = fmt.Errorf("unable to convert %s to integer: %s", s, err.Error())
		status = http.StatusBadRequest
		return
	}
	return
}

func ConvertMapKeys(m map[string]interface{}, converter func(s string) string) (res map[string]interface{}) {
	res = make(map[string]interface{})
	for key, value := range m {
		res[converter(key)] = value
	}
	return
}

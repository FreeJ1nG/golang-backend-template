package dto

import "github.com/FreeJ1nG/backend-template/app/pagination"

var EmptySuccessMessage = map[string]bool{
	"success": true,
}

type ErrorData struct {
	Message string `json:"message"`
}

type ErrorPayload struct {
	Status int       `json:"status"`
	Data   ErrorData `json:"data"`
}

type Response[T interface{}] struct {
	Data     *T                   `json:"data,omitempty"`
	Error    *ErrorPayload        `json:"error,omitempty"`
	Metadata *pagination.Metadata `json:"metadata,omitempty"`
}

func NewSuccessResponse[T interface{}](data T, metadata *pagination.Metadata) (res Response[T]) {
	return Response[T]{
		Data:     &data,
		Error:    nil,
		Metadata: metadata,
	}
}

func NewErrorResponse(errorMessage string, status int) (res Response[struct{}]) {
	return Response[struct{}]{
		Data: nil,
		Error: &ErrorPayload{
			Status: status,
			Data: ErrorData{
				Message: errorMessage,
			},
		},
		Metadata: nil,
	}
}

package dto

type ErrorData struct {
	Message string `json:"message"`
}

type ErrorPayload struct {
	Status int       `json:"status"`
	Data   ErrorData `json:"data"`
}

type Response[T interface{}] struct {
	Data  *T            `json:"data,omitempty"`
	Error *ErrorPayload `json:"error,omitempty"`
}

func NewSuccessResponse[T interface{}](data T) (res Response[T]) {
	return Response[T]{
		Data:  &data,
		Error: nil,
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
	}
}

package domain

type MessageWrapper struct {
	Message string `json:"message"`
}

func NewMessageWrapper(msg string) MessageWrapper {
	return MessageWrapper{Message: msg}
}

type IdWrapper struct {
	Id int64 `json:"id"`
}

func NewIdWrapper(id int64) IdWrapper {
	return IdWrapper{Id: id}
}

type Result[T any] struct {
	Success      bool   `json:"success"`
	Data         T      `json:"data,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	ShowType     int    `json:"showType,omitempty"`
}

func FailureResult[T any](message string) *Result[T] {
	return &Result[T]{
		Success:      false,
		ErrorMessage: message,
		ErrorCode:    "400",
		ShowType:     2,
	}
}

func SuccessResult[T any](data T) *Result[T] {
	return &Result[T]{
		Success: true,
		Data:    data,
	}
}

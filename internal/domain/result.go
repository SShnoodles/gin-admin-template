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

// Result 传统返回结构
type Result[T any] struct {
	Data         *T     `json:"data"`
	ErrorMessage string `json:"errorMessage"`
	ResultCode   string `json:"resultCode"`
}

func FailureResult[T any](message string) *Result[T] {
	return &Result[T]{
		ErrorMessage: message,
		ResultCode:   "400",
	}
}

func SuccessResult[T any](data *T) *Result[T] {
	return &Result[T]{
		Data:         data,
		ErrorMessage: "",
		ResultCode:   "200",
	}
}

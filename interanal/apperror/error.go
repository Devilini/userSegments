package apperror

import "encoding/json"

var (
	ErrNotFound = NewAppError(nil, "Not found")
)

type AppError struct {
	Err     error  `json:"-"`
	Message string `json:"message,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
	}
}

func SystemError(err error) *AppError {
	return NewAppError(err, "internal system error")
}

func NotFoundError(text string) *AppError {
	ErrNotFound = NewAppError(nil, text)
	return ErrNotFound
}

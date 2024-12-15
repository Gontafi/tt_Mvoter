package models

type Error struct {
	IsError bool   `json:"is_error"`
	Message string `json:"message"`
}

func (e *Error) SetError(msg string) {
	e.IsError = true
	e.Message = msg
}

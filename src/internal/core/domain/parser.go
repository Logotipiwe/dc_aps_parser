package domain

import "errors"

type ParserData struct {
	ChatID int64
	Link   string
}

type NotAllowedError struct {
	error
	RequestedNum int
	AllowedNum   int
}

func NewNotAllowedError(requested, allowed int) NotAllowedError {
	return NotAllowedError{
		error:        errors.New("requested aps num not allowed"),
		RequestedNum: requested,
		AllowedNum:   allowed,
	}
}

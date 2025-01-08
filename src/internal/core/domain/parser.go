package domain

import "errors"

type ParserData struct {
	ChatID   int64
	Link     string
	UserName string
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

type ParserParams struct {
	ChatID               int64
	ParseLink            string
	IsStartedFromStorage bool
	UserName             string
}

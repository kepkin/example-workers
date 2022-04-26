package ipm

import (
	"fmt"
)

type ErrorID int

const (
	ErrNoSuchTarget ErrorID = iota
	ErrDuplicateMonitor
)

var errorMsg = map[ErrorID]string{
	ErrNoSuchTarget:     "no such target %v",
	ErrDuplicateMonitor: "such monitor already exists %v",
}

type Error struct {
	id ErrorID
	m  string
}

func (e Error) Is(target error) bool {
	err, ok := target.(Error)
	if !ok {
		return false
	}
	return err.id == e.id
}

func (e Error) Error() string {
	return e.m
}

func NewError(id ErrorID, args ...interface{}) error {
	return &Error{
		id: id,
		m:  fmt.Sprintf(errorMsg[id], args...),
	}
}

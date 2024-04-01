package util

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidMacAddress = errors.New("invalid mac address")
)

type TimeInvalidArgument struct {
	Message string
}

func (e *TimeInvalidArgument) Error() string {
	return fmt.Sprintf("invalid time argument: %v", e.Message)
}

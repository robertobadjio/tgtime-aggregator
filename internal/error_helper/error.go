package error_helper

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidMacAddress ???
	ErrInvalidMacAddress = errors.New("invalid mac address")
)

// TimeInvalidArgument ???
type TimeInvalidArgument struct {
	Message string
}

// Error Возвращает текст ошибки
func (e *TimeInvalidArgument) Error() string {
	return fmt.Sprintf("invalid time argument: %v", e.Message)
}

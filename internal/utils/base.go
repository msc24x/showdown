package utils

import (
	"errors"
	"fmt"
)

func NewError(err error, reason string) error {
	if err == nil {
		return nil
	}
	reason = fmt.Sprintf("[%s] %s", reason, err.Error())
	return errors.New(reason)
}

package utils

import (
	"errors"
	"fmt"
	"log"
)

// Prints log as a warning
func LogWarn(msg string) {
	log.Printf("[WARNING] %s", msg)
}

// Prints log for worker methods
func LogWorker(msg string, a ...any) {
	log.Printf("[WORKER] %s", fmt.Sprintf(msg, a...))
}

// Creates a prefixed error object from string message
func NewError(err error, reason string) error {
	if err == nil {
		return nil
	}
	reason = fmt.Sprintf("[%s] %s", reason, err.Error())
	return errors.New(reason)
}

// Panic if error exists
func PanicIf(err error) {
	if err != nil {
		msg := NewError(err, "PANIC").Error()
		LogWarn(msg)
		panic(msg)
	}
}

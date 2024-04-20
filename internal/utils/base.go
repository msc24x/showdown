package utils

import (
	"errors"
	"fmt"
	"log"
)

// Prints log as a warning
func LogWarn(fmsg string, a ...any) {
	log.Printf("[WARNING] %s", fmt.Sprintf(fmsg, a...))
}

// Prints log for worker methods
func LogWorker(fmsg string, a ...any) {
	log.Printf("[WORKER] %s", fmt.Sprintf(fmsg, a...))
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

func BPanicIf(flag bool, freason string, fargs ...any) {
	if flag {
		msg := NewError(errors.New(fmt.Sprintf(freason, fargs...)), "PANIC").Error()
		LogWarn(msg)
		panic(msg)
	}
}

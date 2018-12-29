package utils

import "log"

// HandleError handles error by log the message and
// execute the error callback (without any parameters) if
// the error is not nil
func HandleError(err error, msg string, errCallback func()) {
	if err != nil {
		if msg == "" {
			msg = err.Error()
		}
		log.Println(msg)
		if errCallback != nil {
			errCallback()
		}
	}
}

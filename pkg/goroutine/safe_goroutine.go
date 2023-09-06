package goroutine

import (
	"log"
)

// Run executes Fn a function safely in a goroutine. In case of a panic, it recovers and logs the error.
func Run(fn func()) {
	go func() {
		defer RecoverPanicAndError()
		fn()
	}()
}

func RecoverPanicAndError() {
	if r := recover(); r != nil {
		log.Println("recovered from panic")
	}
}

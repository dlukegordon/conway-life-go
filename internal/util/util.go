package util

import (
	"log"
)

func FailIf(err error, msg string) {
	if err != nil {
		log.Fatalf("!!! Error %s: %v", msg, err)
	}
}

func Fail(msg string) {
	log.Fatalf("!!! Error %s", msg)
}

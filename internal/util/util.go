package util

import (
	"log"
	"regexp"
)

func FailIf(err error, msg string) {
	if err != nil {
		log.Fatalf("!!! Error %s: %v", msg, err)
	}
}

func Fail(msg string) {
	log.Fatalf("!!! Error %s", msg)
}

func NumEscapeChars(str string) uint {
	// ANSI escape sequences follow the pattern: ESC [ some characters ending with a letter or tilde
	escapeRegex := regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z~]`)

	escapeSequences := escapeRegex.FindAllString(str, -1)

	count := 0
	for _, seq := range escapeSequences {
		count += len(seq)
	}

	return uint(count)
}

func Abs(n int) uint {
	if n < 0 {
		n *= -1
	}
	return uint(n)
}

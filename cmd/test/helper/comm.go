package helper

import (
	"bytes"
	"log"
)

func CaptureOutput(f func()) string {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	f()

	return buf.String()
}

package test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/baijum/logger"
)

// BufLog takes a reader and log it and return another reader
// with the same content
func BufLog(t *testing.T, r io.Reader, msg string) io.Reader {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	b, err := ioutil.ReadAll(tee)
	if err != nil {
		t.Error("Cannot read from buffer")
	}
	if logger.Level <= logger.DEBUG {
		t.Log(msg, string(b))
	}
	return &buf
}

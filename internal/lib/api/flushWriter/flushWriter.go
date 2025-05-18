package flushWriter

import (
	"io"
	"net/http"
)

type FlushWriter struct {
	W       io.Writer
	Flusher http.Flusher
}

func (fw *FlushWriter) Write(p []byte) (int, error) {
	n, err := fw.W.Write(p)
	if fw.Flusher != nil {
		fw.Flusher.Flush()
	}
	return n, err
}

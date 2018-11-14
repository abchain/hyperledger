package util

import (
	"fmt"
	"hash"
	"io"
)

type ConWriter interface {
	Write([]byte) ConWriter
	Error() error
}

type failConWriter struct {
	Err error
}

func (w *failConWriter) Write([]byte) ConWriter {
	return w
}

func (w *failConWriter) Error() error {
	return w.Err
}

type conWriter struct {
	io.Writer
}

func (w conWriter) Write(p []byte) ConWriter {
	i, err := w.Writer.Write(p)

	if err != nil {
		return &failConWriter{err}
	} else if i < len(p) {
		return &failConWriter{fmt.Errorf("Write %d for %d bytes", i, len(p))}
	} else {
		return w
	}
}

func (w conWriter) Error() error {
	return nil
}

func NewConWriter(wr io.Writer) ConWriter {
	return &conWriter{wr}
}

func NewHashWriter(h hash.Hash) ConWriter {
	return NewConWriter(h)
}

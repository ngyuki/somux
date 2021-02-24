package main

import "io"

// ReaderWriterMix is
type ReaderWriterMix struct {
	r io.ReadCloser
	w io.WriteCloser
}

// Read is
func (mix *ReaderWriterMix) Read(p []byte) (n int, err error) {
	return mix.r.Read(p)
}

// Write is
func (mix *ReaderWriterMix) Write(p []byte) (n int, err error) {
	return mix.w.Write(p)
}

// Close is
func (mix *ReaderWriterMix) Close() error {
	err1 := mix.r.Close()
	err2 := mix.w.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

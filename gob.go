package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"unsafe"
)

func readGob(r io.Reader, data interface{}) (err error) {
	sizebuf := make([]byte, uint16(unsafe.Sizeof(uint16(0))))
	n, err := r.Read(sizebuf)
	if err != nil {
		return
	}
	if n != len(sizebuf) {
		err = fmt.Errorf("read gob too short: %v < %v", n, len(sizebuf))
		return
	}

	size := binary.BigEndian.Uint16(sizebuf)

	databuf := make([]byte, size)
	n, err = r.Read(databuf)
	if err != nil {
		return
	}
	if n != len(databuf) {
		err = fmt.Errorf("forward desc recv too short: %v < %v", n, len(databuf))
		return
	}

	err = gob.NewDecoder(bytes.NewBuffer(databuf)).Decode(data)
	if err != nil {
		return
	}
	return
}

func writeGob(w io.Writer, data interface{}) error {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(data)
	if err != nil {
		return err
	}

	databuf := buf.Bytes()

	sizebuf := make([]byte, uint16(unsafe.Sizeof(uint16(0))))
	binary.BigEndian.PutUint16(sizebuf, uint16(len(databuf)))
	n, err := w.Write(sizebuf)
	if err != nil {
		return err
	}
	if n != len(sizebuf) {
		err = fmt.Errorf("forward desc send too short: %v < %v", n, len(sizebuf))
		return err
	}

	n, err = w.Write(databuf)
	if err != nil {
		return err
	}
	if n != len(databuf) {
		err = fmt.Errorf("forward desc send too short: %v < %v", n, len(databuf))
		return err
	}

	return nil
}

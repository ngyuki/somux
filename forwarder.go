package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/hashicorp/yamux"
	"golang.org/x/sync/errgroup"
)

func startForwardListener(logger *Logger, session *yamux.Session, local string, remote string) error {

	logger.printf("listen forward: %v -> %v", local, remote)

	listener, err := net.Listen("tcp4", local)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			stream, err := session.OpenStream()
			if err != nil {
				conn.Close()
				logger.errorln(err)
				return
			}

			logger = logger.withStreamID(stream.StreamID())

			logger.printf("connect accept: %v", conn.RemoteAddr())
			defer func() {
				logger.printf("connect close: %v", conn.RemoteAddr())
				conn.Close()
			}()

			logger.printf("stream open: %06d", stream.StreamID())
			defer func() {
				logger.printf("stream close: %06d", stream.StreamID())
				stream.Close()
			}()

			logger.printf("forward: %v -> %v", local, remote)

			err = writeGob(stream, remote)
			if err != nil {
				logger.errorln(err)
				return
			}

			err = handleForward(logger, conn, stream)
			if err != nil {
				logger.errorln(err)
				return
			}
		}()
	}
}

func startSessionAccepter(logger *Logger, session *yamux.Session) error {
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			return fmt.Errorf("stream accept: %w", err)
		}

		go func() {
			logger := logger.withStreamID(stream.StreamID())

			logger.printf("stream accept: %06d", stream.StreamID())
			defer func() {
				logger.printf("stream close: %06d", stream.StreamID())
				stream.Close()
			}()

			var remote string
			err := readGob(stream, &remote)
			if err != nil {
				err = fmt.Errorf("recv desc: %w", err)
				logger.errorln(err)
				return
			}

			logger.printf("connect: stream -> %v", remote)
			conn, err := net.Dial("tcp", remote)
			if err != nil {
				err = fmt.Errorf("connect error: %w", err)
				logger.errorln(err)
				return
			}

			logger.printf("connect open: %v", remote)
			defer func() {
				logger.printf("connect close: %v", remote)
				conn.Close()
			}()

			err = handleForward(logger, conn, stream)
			if err != nil {
				logger.errorln(err)
			}
		}()
	}
}

func handleForward(logger *Logger, conn net.Conn, stream *yamux.Stream) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		defer stream.SetDeadline(time.Now())
		_, err := io.Copy(stream, conn)
		if err != nil {
			if !(errors.Is(err, io.EOF) || errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, yamux.ErrTimeout)) {
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		defer conn.SetDeadline(time.Now())
		_, err := io.Copy(conn, stream)
		if err != nil {
			if !(errors.Is(err, io.EOF) || errors.Is(err, os.ErrDeadlineExceeded) || errors.Is(err, yamux.ErrTimeout)) {
				return err
			}
		}
		return nil
	})

	return g.Wait()
}

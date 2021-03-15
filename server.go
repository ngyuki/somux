package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/yamux"
)

func runServer(ctx context.Context, logger *Logger) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger.println("start server session")
	conn := &ReaderWriterMix{r: os.Stdin, w: os.Stdout}
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return err
	}
	defer session.Close()

	setting := setting{}

	{
		stream, err := session.AcceptStream()
		if err != nil {
			return fmt.Errorf("accept initial stream: %w", err)
		}
		defer func() {
			logger.println("initial stream closing")
			stream.Close()
			logger.println("initial stream closed")
		}()

		err = readGob(stream, &setting)
		if err != nil {
			return fmt.Errorf("read initial stream: %w", err)
		}
	}

	for _, v := range setting.ReverseForwards {
		forward := v
		go func() {
			defer cancel()
			err := startForwardListener(logger, session, forward.Bind, forward.Connect)
			if err != nil {
				logger.errorln(err)
			}
		}()
	}

	go func() {
		defer cancel()
		err := startSessionAccepter(logger, session)
		if err != nil {
			logger.errorln(err)
		}
	}()

	<-ctx.Done()
	return nil
}

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/yamux"
)

func runServer(logger *Logger) error {

	conn := &ReaderWriterMix{r: os.Stdin, w: os.Stdout}

	logger.println("start server session")
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
		defer stream.Close()

		err = readGob(stream, &setting)
		if err != nil {
			return fmt.Errorf("read initial stream: %w", err)
		}
	}

	for _, v := range setting.ReverseForwards {
		forward := v
		go func() {
			err := startForwardListener(logger, session, forward.Bind, forward.Connect)
			if err != nil {
				logger.errorln(err)
			}
		}()
	}

	return startSessionAccepter(logger, session)
}

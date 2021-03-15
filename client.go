package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/hashicorp/yamux"
)

func runClient(ctx context.Context, logger *Logger, setting *setting) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd, conn, err := startServerProcess(ctx, logger, setting.Command)
	if err != nil {
		return err
	}

	logger.println("start client session")
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}
	defer session.Close()

	stream, err := session.OpenStream()
	if err != nil {
		return fmt.Errorf("open initial stream: %w", err)
	}
	defer stream.Close()

	err = writeGob(stream, setting)
	if err != nil {
		return fmt.Errorf("write initial data: %w", err)
	}

	for _, v := range setting.LocalForwards {
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

	cmd.Wait()
	logger.printf("exit server process: %v", cmd.ProcessState.ExitCode())
	return nil
}

func startServerProcess(ctx context.Context, logger *Logger, command []string) (*exec.Cmd, io.ReadWriteCloser, error) {
	cmd := exec.Command(
		command[0],
		command[1:]...,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	logger.println("spawn server process", command)
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text())
		}
	}()

	go func() {
		<-ctx.Done()
		logger.println("terminate server process")
		cmd.Process.Signal(syscall.SIGTERM)
	}()

	conn := &ReaderWriterMix{r: stdout, w: stdin}
	return cmd, conn, nil
}

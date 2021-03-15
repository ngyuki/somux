package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetOutput(os.Stderr)

	setting, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.SetFlags(0)
	if setting.Verbose {
		log.SetFlags(log.Lshortfile)
	}

	logger := newLogger()
	logger.verbose = setting.Verbose

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		s := <-sig
		logger.printf("recv signal %v", s)
		cancel()
	}()

	if len(setting.Command) > 0 {
		logger = logger.withName("C")
		err = runClient(ctx, logger, setting)
	} else {
		logger = logger.withName("S")
		err = runServer(ctx, logger)
	}
	if err != nil {
		logger.errorln(err)
		os.Exit(1)
	}
	logger.printf("exit: ok")
}

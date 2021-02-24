package main

import (
	"fmt"
	"log"
	"os"
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

	if len(setting.Command) > 0 {
		logger = logger.withName("C")
		err = runClient(logger, setting)
	} else {
		logger = logger.withName("S")
		err = runServer(logger)
	}
	if err != nil {
		logger.errorln(err)
		os.Exit(1)
	}
	logger.printf("exit: ok")
}

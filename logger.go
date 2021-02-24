package main

import (
	"fmt"
	"log"
)

// Logger is
type Logger struct {
	name    string
	id      uint32
	verbose bool
	logger  *log.Logger
}

func newLogger() *Logger {
	return &Logger{
		logger: log.New(log.Writer(), "", log.Flags()),
	}
}

func (logger *Logger) formatPrefix(name string, id uint32) string {
	if logger.verbose {
		if id == 0 {
			return fmt.Sprintf("(%1s) [%6s] ", name, "")
		}
		return fmt.Sprintf("(%1s) [%06d] ", name, id)
	} else if id == 0 {
		return fmt.Sprintf("(%1s) ", name)
	}
	return fmt.Sprintf("(%1s) ", name)
}

func (logger *Logger) withName(name string) *Logger {
	return &Logger{
		name:    name,
		id:      logger.id,
		verbose: logger.verbose,
		logger: log.New(
			logger.logger.Writer(),
			logger.formatPrefix(name, logger.id),
			logger.logger.Flags(),
		),
	}
}

func (logger *Logger) withStreamID(id uint32) *Logger {
	return &Logger{
		name:    logger.name,
		id:      id,
		verbose: logger.verbose,
		logger: log.New(
			logger.logger.Writer(),
			logger.formatPrefix(logger.name, id),
			logger.logger.Flags(),
		),
	}
}

func (logger *Logger) printf(fotmat string, v ...interface{}) {
	if logger.verbose {
		logger.logger.Output(2, fmt.Sprintf(fotmat, v...))
	}
}

func (logger *Logger) println(v ...interface{}) {
	if logger.verbose {
		logger.logger.Output(2, fmt.Sprintln(v...))
	}
}

func (logger *Logger) errorln(err error) {
	logger.logger.Output(2, fmt.Sprintln(err))
}

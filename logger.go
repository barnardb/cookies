package main

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	writer    io.Writer
	verbosity int
}

func LoggerWithVerbosity(verbosity int) *Logger {
	if verbosity > 0 {
		return &Logger{os.Stderr, verbosity}
	}
	return nil
}

func (l *Logger) RequireVerbosity(verbosity int) *Logger {
	if l != nil && l.verbosity >= verbosity {
		return l
	}
	return nil
}

func (l *Logger) Printf(format string, a ...interface{}) {
	if l == nil {
		return
	}
	fmt.Fprintf(l.writer, string(append([]byte(format), 0x0a)), a...)
}

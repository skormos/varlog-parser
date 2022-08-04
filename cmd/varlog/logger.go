package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func stdoutLoggerContext(module string) zerolog.Context {
	return loggerContext(os.Stdout, module)
}

func loggerContext(writer io.Writer, module string) zerolog.Context {
	return zerolog.New(writer).
		With().
		Timestamp().
		Stack().
		Caller().
		Str("module", module)
}

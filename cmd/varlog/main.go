package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/skormos/varlog-parser/internal/syslog"
)

type Parser interface {
	Parse(ctx context.Context, reader io.Reader) ([]*syslog.LogLine, error)
}

func main() {
	mainLogger := stdoutLoggerContext("main").Logger()

	mainLogger.Info().Msg("started")

	parser, err := syslog.NewSyslogParser()
	if err != nil {
		mainLogger.Err(err).Msg("while instantiating syslog parser. Exiting...")
		return
	}

	filename := "/var/log/system.log"
	lines, err := parseFile(context.Background(), filename, parser, func(closeErr error) {
		if closeErr != nil {
			mainLogger.Err(closeErr).Msgf("while closing file %s", filename)
		}
	})

	if err != nil {
		mainLogger.Err(err).Msgf("while parsing file %s", filename)
	} else {
		for _, line := range lines {
			fmt.Printf("%+v\n", line)
		}
	}

	mainLogger.Info().Msg("good-bye")
}

func parseFile(ctx context.Context, filename string, parser Parser, closeCallback func(error)) ([]*syslog.LogLine, error) {
	fileReader, err := os.Open(filename)
	defer func() {
		closeCallback(fileReader.Close())
	}()
	if err != nil {
		return nil, fmt.Errorf("while opening file: %w", err)
	}

	res, err := parser.Parse(ctx, fileReader)
	if err != nil {
		return nil, fmt.Errorf("while parsing file: %w", err)
	}

	return res, nil
}

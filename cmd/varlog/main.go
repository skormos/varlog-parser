package main

import (
	"context"
	"fmt"
	stdos "os"

	"github.com/skormos/varlog-parser/internal/logparser"
	"github.com/skormos/varlog-parser/internal/os"
)

func main() {
	mainLogger := stdoutLoggerContext("main").Logger()

	config := parseFlags()

	mainLogger.Info().Msg("started")

	fileHandler, err := os.NewFileHandler(config.dirPath)
	if err != nil {
		mainLogger.Err(err).Msgf("registering file handler")
		if pwd, err := stdos.Getwd(); err != nil {
			mainLogger.Err(err).Msg("could not retrieve present working directory.")
		} else {
			mainLogger.Info().Msgf("Present working directory: %s", pwd)
		}
		return
	}

	logFile, err := fileHandler.Open("benchmark.log")
	defer func() {
		if err := logFile.Close(); err != nil {
			mainLogger.Warn().Err(err).Msgf("closing requested log file")
		}
	}()
	if err != nil {
		mainLogger.Err(err).Msgf("retrieving requested log file")
		return
	}

	lines, err := logparser.ParseLastNLines(context.Background(), logFile, 7)
	if err != nil {
		mainLogger.Err(err).Msgf("parsing requested log file")
		return
	}

	for idx, line := range lines {
		fmt.Println(idx, line)
	}

	mainLogger.Info().Msg("good-bye")
}

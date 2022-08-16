package main

import (
	"context"
	"fmt"
	stdhttp "net/http"
	stdos "os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"

	"github.com/skormos/varlog-parser/cmd/varlog/http"
	"github.com/skormos/varlog-parser/internal/logparser"
	"github.com/skormos/varlog-parser/internal/os"
)

func main() {
	mainLogger := stdoutLoggerContext("main").Logger()
	defer recoverPanic(mainLogger)

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

	server := http.NewServerWrapper(stdoutLoggerContext("http"), stdhttp.HandlerFunc(func(writer stdhttp.ResponseWriter, _ *stdhttp.Request) {
		_, _ = writer.Write([]byte("this is the output"))
	}), http.WithPort(config.http.port))

	grp := new(errgroup.Group)
	grp.Go(onShutdown(mainLogger, server.Stop))
	grp.Go(server.Start)

	if err := grp.Wait(); err != nil {
		mainLogger.Err(err).Msgf("unexpected shutdown")
	} else {
		mainLogger.Info().Msg("Successfully completed shutdown")
	}
}

func onShutdown(logger zerolog.Logger, shutdownFn func()) func() error {
	return func() error {
		signalChan := signalShutdown()
		select { //nolint:gosimple // select shouldn't be used for a single channel, but this is more readable
		case sig := <-signalChan:
			if syscall.SIGSTOP == sig {
				logger.Info().Msg("integrity issue has invoked a shutdown...")
			} else {
				logger.Info().Msgf("%v : server shutdown requested...", sig)
			}

			shutdownFn()
		}
		return nil
	}
}

func signalShutdown() chan stdos.Signal {
	shutdown := make(chan stdos.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return shutdown
}

func recoverPanic(logger zerolog.Logger) {
	if r := recover(); r != nil {
		logger.Error().Msgf("Main Recovered from Panic: %v", r)
	}
}

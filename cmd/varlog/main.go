package main

import (
	stdos "os"
	"os/signal"
	"syscall"

	"github.com/skormos/varlog-parser/internal/handler/varlog"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"

	"github.com/skormos/varlog-parser/cmd/varlog/http"
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
	mainLogger.Info().Msgf("file handler registered for directory: %s", config.dirPath)

	httpLogContext := stdoutLoggerContext("http")

	httpHandler := rootHandler(apiHandler(varlog.NewHandler(httpLogContext, fileHandler)))
	server := http.NewServerWrapper(httpLogContext, httpHandler, http.WithPort(config.http.port))

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

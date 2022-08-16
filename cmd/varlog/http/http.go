// Package http provides convenience wrappers and options for creating a http.Server instance.
package http

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type (
	// ServerOption defines the function signature for helper methods to update values on the wrapper.
	ServerOption func(wrapper *ServerWrapper)

	// ServerWrapper is an abstraction to make it easier to start, shutdown and configure a http.Server instance.
	ServerWrapper struct {
		logger          zerolog.Logger
		server          *http.Server
		shutdownTimeout time.Duration
	}
)

// WithHostPort sets the host and port for the server to listen on.
func WithHostPort(host, port string) ServerOption {
	return func(wrapper *ServerWrapper) {
		wrapper.server.Addr = net.JoinHostPort(host, port)
	}
}

// WithPort sets the port for the server to listen on, and an empty host string.
func WithPort(port string) ServerOption {
	return WithHostPort("", port)
}

// WithReadTimeout sets the ReadTimeout value on the underlying http.Server instance.
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(wrapper *ServerWrapper) {
		wrapper.server.ReadTimeout = timeout
	}
}

// WithShutdownTimeout sets the time to wait to try and shutdown the underlying http.Server instance gracefully.
func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return func(wrapper *ServerWrapper) {
		wrapper.shutdownTimeout = timeout
	}
}

// WithWriteTimeout sets the WriteTimeout value on the underlying http.Server instance.
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(wrapper *ServerWrapper) {
		wrapper.server.WriteTimeout = timeout
	}
}

// NewServerWrapper creates a new instance of a http.Server and applies the provided option functions. The error
// callback function is used to handle errors returned from the http.Server.ListenAndServe() function.
func NewServerWrapper(logCtx zerolog.Context, rootHandler http.Handler, options ...ServerOption) *ServerWrapper {
	logger := logCtx.Logger()
	wrapper := &ServerWrapper{
		logger: logger,
		server: &http.Server{
			Addr:              net.JoinHostPort("", "8080"),
			Handler:           rootHandler,
			ReadTimeout:       120 * time.Second,
			WriteTimeout:      120 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			ErrorLog:          log.New(logger, "", log.LstdFlags),
		},
		shutdownTimeout: 60 * time.Second,
	}

	for _, optionFn := range options {
		optionFn(wrapper)
	}

	return wrapper
}

// Start uses the provided options and defaults to start the http server. It will return an error IFF the error is not
// http.ErrServerClosed.
func (w *ServerWrapper) Start() error {
	w.logger.Info().Msgf("starting http server at: %s", w.server.Addr)
	if err := w.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

// Stop attempts to gracefully stop the server, and if it passes the timeout will force close.
func (w *ServerWrapper) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), w.shutdownTimeout)
	defer cancel()

	err := w.server.Shutdown(ctx)
	if err != nil {
		w.logger.Info().Err(err).Msgf("Graceful shutdown did not complete in: %s. Attempting close...", w.shutdownTimeout.String())
		if err := w.server.Close(); err != nil {
			w.logger.Info().Err(err).Msg("while attempting server close")
		} else {
			w.logger.Info().Msg("server close completed successfully.")
		}
	} else {
		w.logger.Info().Msg("Graceful shutdown of server completed successfully.")
	}
}

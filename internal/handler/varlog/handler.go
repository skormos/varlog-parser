package varlog

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	v1 "github.com/skormos/varlog-parser/internal/api/rest/v1"
)

// NewHandler creates a new http.Handler which conforms to the varlog API Spec.
func NewHandler(logCtx zerolog.Context, opener FileOpener) http.Handler {
	logger := logCtx.Logger()

	return v1.HandlerWithOptions(NewLogParserHandler(logCtx, opener), v1.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Err(err).Msgf("while calling %s", r.RequestURI)

			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	})
}

func respond(writer http.ResponseWriter, input interface{}, status int) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("while marshalling %v for http response: %w", input, err)
	}

	writer.WriteHeader(status)
	if _, err := writer.Write(bytes); err != nil {
		return fmt.Errorf("while writing %v as bytes to Response: %w", input, err)
	}

	return nil
}

package handler

import (
	"net/http"
	stdos "os"

	"github.com/rs/zerolog"

	v1 "github.com/skormos/varlog-parser/internal/api/rest/v1"
	"github.com/skormos/varlog-parser/internal/logparser"
	"github.com/skormos/varlog-parser/internal/os"
)

type (
	// FileOpener defines the interface which is used to open file resources based on a single file name.
	FileOpener interface {
		Open(filename string) (*stdos.File, error)
	}

	// LogParserHandler implements the v1 ServerInterface to open files and read lines from the end of it.
	LogParserHandler struct {
		opener FileOpener
		logger zerolog.Logger
	}

	rawLines []string
)

// GetEntries uses the provided FileOpener implementation to retrieve the most recent log entries as part of the query
// parameters.
func (l *LogParserHandler) GetEntries(w http.ResponseWriter, r *http.Request, filename string, params v1.GetEntriesParams) {
	reader, err := l.opener.Open(filename)
	defer func() {
		if reader != nil {
			if closeErr := reader.Close(); closeErr != nil {
				l.logger.Err(closeErr).Msgf("could not close file %s", filename)
			}
		}
	}()
	if err != nil {
		if err == os.ErrNotExists {
			http.Error(w, "requested file with name could not be located", http.StatusNotFound)
			return
		}

		if err == os.ErrNoReadPerm {
			http.Error(w, "requested file does not have sufficient permissions to be read", http.StatusForbidden)
			return
		}

		l.logger.Err(err).Msgf("while requesting filename: %s", filename)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	numLines := 25
	if params.NumEntries != nil {
		val := params.NumEntries
		if *val > 1000 {
			numLines = 1000
		} else if *val > 0 {
			numLines = *val
		}
	}

	lines, err := logparser.ParseLastNLines(r.Context(), reader, numLines)
	if err != nil {
		l.logger.Err(err).Msgf("while parsing %d lines for file %s", numLines, filename)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := v1.GetEntriesResponse{
		Entries: rawLines(lines).toEntries(),
	}

	if err := respond(w, resp, http.StatusOK); err != nil {
		l.logger.Err(err).Msgf("attempting to send a response for file %s", filename)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// NewLogParserHandler returns a new instance of the LogParserHandler.
func NewLogParserHandler(logCtx zerolog.Context, opener FileOpener) *LogParserHandler {
	return &LogParserHandler{
		logger: logCtx.Str("handler", "logparser").Logger(),
		opener: opener,
	}
}

func (l rawLines) toEntries() []v1.LogEntry {
	out := make([]v1.LogEntry, 0, len(l))

	for _, line := range l {
		out = append(out, v1.LogEntry(line))
	}

	return out
}

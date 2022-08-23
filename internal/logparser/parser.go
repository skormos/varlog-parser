package logparser

import (
	"context"
	"fmt"
	"os"
	"strings"
)

const defaultLineSize = 120

// ParseLastNLinesSeek takes an open File, seeks to end, and attempts to read via chunks backward from the bottom of the
// file, returning at most of the number of requested lines. The results assume newer lines are appended to the end of
// the file, and since the results are returned in descending order of when they were appended, the last line will be
// the first item in the return slice.
//
// The specified filter is applied inline, so the results will only contain at most the nLines that pass the filter.
// if there is an error in file.Read or file.Seek operations, those will be wrapped and returned.
//
// Currently, the context parameter is not used, but is reserved for idiomatic usage later.
//
// It is up to the caller of this method to manager the file on return or on error.
func ParseLastNLinesSeek(_ context.Context, file *os.File, nLines int, filter Filterer) ([]string, error) {
	if nLines <= 0 {
		return []string{}, nil
	}

	chunkSize := int64(nLines * defaultLineSize)
	out := make([]string, 0, nLines)

	pos, err := file.Seek(0, 2)
	if err != nil {
		return nil, fmt.Errorf("while getting file size %w", err)
	}

	readSize := chunkSize

	for len(out) < nLines && pos > 0 {

		if pos < readSize {
			readSize = pos
		}

		// start at the end, and seek back a chunk
		_, err = file.Seek(int64(-1*readSize), 1)
		if err != nil {
			return nil, fmt.Errorf("while seeking %w", err)
		}

		// attempt to read the chunk from the seek spot
		b := make([]byte, readSize)
		bytesRead, err := file.Read(b)
		if err != nil {
			return nil, fmt.Errorf("while reading %w", err)
		}

		// strip \r if they exist
		raw := strings.ReplaceAll(string(b), "\r", "")

		// split on new lines
		lines := strings.Split(raw, "\n")

		if readSize == chunkSize {
			// remove the first line, as it's incomplete, and use the length to seek forward to the next new line
			cutLine := lines[0]
			lines = lines[1:]
			pos, err = file.Seek(int64(-1*(bytesRead-len(cutLine))), 1)
			if err != nil {
				return nil, fmt.Errorf("while moving forward on seek %w", err)
			}
		} else {
			// this should be the last read of the file. Just set the position to 0 to exit.
			pos = 0
		}

		// for the remaining lines, filter the lines, appending results.
		for i := len(lines) - 1; i >= 0; i-- {
			if filter.Filter(lines[i]) {
				out = append(out, lines[i])
			}

			if len(out) == nLines {
				break
			}
		}
	}

	return out, nil
}

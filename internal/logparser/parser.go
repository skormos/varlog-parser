package logparser

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// ParseLastNLines takes a Reader containing strings, splits them by lines, and returns a number of lines from the end
// of the reader. This is returned in reverse order, so the last line is first in the return string slice parameter.
//
// This uses bufio.Scanner internally, so if there's an error reading the file, it will be returned wrapped in the
// error parameter.
//
// Currently, the context parameter is not used, but is reserved for idiomatic usage later.
//
// It is up to the caller of this method to manager the reader.
func ParseLastNLines(_ context.Context, reader io.Reader, nLines int) ([]string, error) {
	scanner := bufio.NewScanner(reader)

	linesRead := 0
	lines := make([]string, nLines)

	for scanner.Scan() {
		lines[linesRead%nLines] = scanner.Text()
		linesRead++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("while scanning lines %w", err)
	}

	var out []string
	out = make([]string, 0, nLines)
	offset := linesRead % nLines

	for i := offset - 1; i >= 0; i-- {
		out = append(out, lines[i])
	}

	if linesRead >= nLines {
		for i := nLines - 1; i >= offset; i-- {
			out = append(out, lines[i])
		}
	}

	return out, nil
}

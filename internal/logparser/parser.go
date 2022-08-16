package logparser

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// ParseLastNLines takes a Reader containing strings, splits them by lines, and returns a number of lines from the end
// of the reader. The results assume newer lines are appended to the end of the file, and since the results are returned
// in descending order of when they were appended, the last line will be the first item in the return slice.
//
// The specified filter is applied inline, so the results will only contain at most the nLines that pass the filter.
//
// This uses bufio.Scanner internally, so if there's an error reading the file, it will be returned wrapped in the
// error parameter.
//
// Currently, the context parameter is not used, but is reserved for idiomatic usage later.
//
// It is up to the caller of this method to manager the reader.
func ParseLastNLines(_ context.Context, reader io.Reader, nLines int, filter Filterer) ([]string, error) {
	scanner := bufio.NewScanner(reader)

	accepted := 0
	lines := make([]string, nLines)

	var text string
	for scanner.Scan() {
		text = scanner.Text()
		if filter.Filter(text) {
			lines[accepted%nLines] = text
			accepted++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("while scanning lines %w", err)
	}

	var out []string
	out = make([]string, 0, nLines)
	offset := accepted % nLines

	for i := offset - 1; i >= 0; i-- {
		out = append(out, lines[i])
	}

	if accepted >= nLines {
		for i := nLines - 1; i >= offset; i-- {
			out = append(out, lines[i])
		}
	}

	return out, nil
}

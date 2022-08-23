package logparser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dummyDataSingle = "Single process wrote this message."
	// reminder the back-ticks contain formatting
	dummyDataMulti = `Multi-line message starts with this line.
	The next line is here, which explains a bunch of things that are related to the log entry.
	Finally, this will just be an additional line to add even more context.`
)

// Unskip this to create a large benchmark file
func TestCreateBenchmarkFile(t *testing.T) {
	t.SkipNow()

	filename := filepath.Join("./testdata", "benchmark2GB.log")
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		if datafile, err := os.Create(filename); err != nil {
			panic(err)
		} else {
			defer func() {
				if err := datafile.Close(); err != nil {
					fmt.Println("while closing benchmark data file", err)
				}
			}()
			if err := populateBenchmarkLog(t, datafile, 2_000_000_000); err != nil {
				panic(err)
			}
		}
	}
}

// quick and dirty way to populate a test file without it taking up long-term space
func populateBenchmarkLog(t *testing.T, file *os.File, size int) error {
	t.Helper()

	total := 0

	now := time.Now().UTC().Format(time.Stamp)
	entryPattern := "%s the-host-name thisprocess[4321] %s\n"

	for total < size {
		written, err := file.WriteString(fmt.Sprintf(entryPattern, now, dummyDataSingle))
		if err != nil {
			return err
		}
		total += written

		written, err = file.WriteString(fmt.Sprintf(entryPattern, now, dummyDataMulti))
		if err != nil {
			return err
		}
		total += written

		written, err = file.WriteString(fmt.Sprintf(entryPattern, now, dummyDataSingle))
		if err != nil {
			return err
		}
		total += written
	}

	return nil
}

func TestParseLastNLinesSeekSpeed(t *testing.T) {
	t.SkipNow()

	file, err := os.Open("./testdata/benchmark2GB.log")
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	require.NoError(t, err)

	start := time.Now()
	out, err := ParseLastNLinesSeek(context.Background(), file, 100, FilterNone())
	require.NoError(t, err)

	dur := time.Now().Sub(start)
	fmt.Println(dur)

	assert.Len(t, out, 100)
}

func TestParseLastNLinesSeek(t *testing.T) {
	file, err := os.Open("./testdata/benchmark-small.log")
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	require.NoError(t, err)

	// last 4 lines
	out, err := ParseLastNLinesSeek(context.TODO(), file, 4, FilterNone())
	require.NoError(t, err)
	assert.Len(t, out, 4)
	assert.True(t, strings.HasPrefix(out[0], "11"))
	assert.True(t, strings.HasPrefix(out[3], "08"))

	// filter for all existing lines with "thisprocess"
	out, err = ParseLastNLinesSeek(context.TODO(), file, 5, FilterOnSubstring("thisprocess"))
	require.NoError(t, err)
	assert.Len(t, out, 5)
	assert.True(t, strings.HasPrefix(out[0], "07"))
	assert.True(t, strings.HasPrefix(out[4], "01"))

	// filter for the first existing line with "thisprocess"
	out, err = ParseLastNLinesSeek(context.TODO(), file, 1, FilterOnSubstring("thisprocess"))
	require.NoError(t, err)
	assert.Len(t, out, 1)
	assert.True(t, strings.HasPrefix(out[0], "07"))

	// try to get more lines than exist
	out, err = ParseLastNLinesSeek(context.TODO(), file, 300, FilterNone())
	require.NoError(t, err)
	assert.Len(t, out, 11)
	assert.True(t, strings.HasPrefix(out[0], "11"))
	assert.True(t, strings.HasPrefix(out[10], "01"))

	// get no lines
	out, err = ParseLastNLinesSeek(context.TODO(), file, 0, FilterNone())
	require.NoError(t, err)
	assert.Empty(t, out)
}

//func TestParseLastNLines(t *testing.T) {
//	tests := map[string]struct {
//		input    string
//		nLines   int
//		expected []string
//	}{
//		"Number of input lines less than NLines returns successfully": {
//			input:    "Line 0\nLine 1\nLine 2\nLine 3\nLine 4\nLine 5",
//			nLines:   100,
//			expected: []string{"Line 5", "Line 4", "Line 3", "Line 2", "Line 1", "Line 0"},
//		},
//		"NLines equals the number of input lines returns successfully": {
//			input:    "Line 0\nLine 1\nLine 2\nLine 3\nLine 4",
//			nLines:   5,
//			expected: []string{"Line 4", "Line 3", "Line 2", "Line 1", "Line 0"},
//		},
//		"Empty input returns successfully": {
//			input:    "",
//			nLines:   10,
//			expected: []string{},
//		},
//		"Number of input lines greater than NLines returns successfully": {
//			input:    "Line 0\nLine 1\nLine 2\nLine 3\nLine 4\nLine 5",
//			nLines:   2,
//			expected: []string{"Line 5", "Line 4"},
//		},
//	}
//
//	for name, test := range tests {
//		t.Run(name, func(tt *testing.T) {
//			actual, actualErr := ParseLastNLinesReader(context.TODO(), strings.NewReader(test.input), test.nLines, FilterNone())
//			require.NoError(t, actualErr)
//			require.NotNil(t, actual)
//
//			assert.Equal(tt, test.expected, actual)
//		})
//	}
//}
//
//func TestParseLastNLines_WithError(t *testing.T) {
//
//	expectedError := errors.New("my custom error")
//
//	actual, actualErr := ParseLastNLinesReader(context.TODO(), iotest.ErrReader(expectedError), 100, FilterNone())
//	require.Nil(t, actual)
//	require.Error(t, actualErr)
//
//	assert.Equal(t, expectedError, errors.Unwrap(actualErr))
//}

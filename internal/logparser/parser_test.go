package logparser

import (
	"context"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLastNLines(t *testing.T) {
	tests := map[string]struct {
		input    string
		nLines   int
		expected []string
	}{
		"Number of input lines less than NLines returns successfully": {
			input:    "Line 0\nLine 1\nLine 2\nLine 3\nLine 4\nLine 5",
			nLines:   100,
			expected: []string{"Line 5", "Line 4", "Line 3", "Line 2", "Line 1", "Line 0"},
		},
		"NLines equals the number of input lines returns successfully": {
			input:    "Line 0\nLine 1\nLine 2\nLine 3\nLine 4",
			nLines:   5,
			expected: []string{"Line 4", "Line 3", "Line 2", "Line 1", "Line 0"},
		},
		"Empty input returns successfully": {
			input:    "",
			nLines:   10,
			expected: []string{},
		},
		"Number of input lines greater than NLines returns successfully": {
			input:    "Line 0\nLine 1\nLine 2\nLine 3\nLine 4\nLine 5",
			nLines:   2,
			expected: []string{"Line 5", "Line 4"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			actual, actualErr := ParseLastNLines(context.TODO(), strings.NewReader(test.input), test.nLines, FilterNone())
			require.NoError(t, actualErr)
			require.NotNil(t, actual)

			assert.Equal(tt, test.expected, actual)
		})
	}
}

func TestParseLastNLines_WithError(t *testing.T) {
	expectedError := errors.New("my custom error")

	actual, actualErr := ParseLastNLines(context.TODO(), iotest.ErrReader(expectedError), 100, FilterNone())
	require.Nil(t, actual)
	require.Error(t, actualErr)

	assert.Equal(t, expectedError, errors.Unwrap(actualErr))
}

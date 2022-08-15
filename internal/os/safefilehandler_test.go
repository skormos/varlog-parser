package os

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUnreadableDirPath  = "./testdata/unreadable"
	testUnreadableFilePath = "./testdata/unreadable.log"

	testPermNoRead   = 0355
	testPermOwnerAll = 0755
)

func TestMain(m *testing.M) {
	if err := os.Chmod(testUnreadableDirPath, testPermNoRead); err != nil {
		panic(fmt.Errorf("could not chmod for unreadable test dir: %w", err))
	}
	if err := os.Chmod(testUnreadableFilePath, testPermNoRead); err != nil {
		panic(fmt.Errorf("could not chmod for unreadable test file: %w", err))
	}

	m.Run()

	if err := os.Chmod(testUnreadableDirPath, testPermOwnerAll); err != nil {
		panic(fmt.Errorf("could reset chmod for unreadable test dir: %w", err))
	}
	if err := os.Chmod(testUnreadableFilePath, testPermOwnerAll); err != nil {
		panic(fmt.Errorf("could reset chmod for unreadable test file: %w", err))
	}
}

func TestNewFileHandler(t *testing.T) {
	tests := map[string]struct {
		path        string
		expectError bool
	}{
		"Non-existent directory should return an error": {
			path:        "./testdata/section31",
			expectError: true,
		},
		"Path to a file should return an error": {
			path:        "./testdata/single.log",
			expectError: true,
		},
		"Unreadable directory should return an error": {
			path:        "./testdata/unreadable",
			expectError: true,
		},
		"Existing readable directory should create a new handler": {
			path: "./testdata",
		},
	}

	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			actual, actualErr := NewFileHandler(test.path)

			if test.expectError {
				require.Nil(tt, actual)
				require.Error(tt, actualErr)
			} else {
				require.NotNil(tt, actual)
				require.NoError(tt, actualErr)
			}
		})
	}
}

func TestSafeFileHandler_Open(t *testing.T) {

	handler, err := NewFileHandler("./testdata")
	require.NoError(t, err)
	require.NotNil(t, handler)

	genericTests := map[string]struct {
		filename    string
		expectError bool
	}{
		"Non existing file returns an error": {
			filename:    "doesnotexist.log",
			expectError: true,
		},
		"Unreadable file should return an error": {
			filename:    "unreadable.log",
			expectError: true,
		},
		"Existing file returns successfully": {
			filename:    "empty.log",
			expectError: false,
		},
	}

	for name, test := range genericTests {
		t.Run(name, func(tt *testing.T) {
			actual, actualErr := handler.Open(test.filename)
			defer func() {
				if actual != nil {
					assert.NoError(tt, actualErr, actual.Close(), "closing test file")
				}
			}()

			if test.expectError {
				require.Nil(tt, actual)
				require.Error(tt, actualErr)
			} else {
				require.NotNil(tt, actual)
				require.NoError(tt, actualErr)
			}
		})
	}

	specificErrorTests := map[string]struct {
		filename      string
		expectedError error
	}{
		"Non existing file returns ErrNotExists": {
			filename:      "doesnotexist.log",
			expectedError: ErrNotExists,
		},
		"Unreadable file returns ErrNoReadPerm": {
			filename:      "unreadable.log",
			expectedError: ErrNoReadPerm,
		},
	}

	for name, test := range specificErrorTests {
		t.Run(name, func(tt *testing.T) {
			actual, actualErr := handler.Open(test.filename)
			require.Nil(tt, actual)
			require.Equal(tt, test.expectedError, actualErr)
		})
	}
}

package os

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// ErrNoReadPerm is returned if a file or directory is unreadable.
	ErrNoReadPerm = errors.New("file does not have read permissions")

	// ErrNotExists is returned if a file or directory does not exist.
	ErrNotExists = errors.New("file does not exist")
)

// SafeFileHandler is a convenience wrapper to perform file operations and abstracting some operations early in the process.
type SafeFileHandler struct {
	dirPath string
}

// NewFileHandler returns a new instance of FileHandler. Provided the given path, it will validate the path exists, and
// is readable.
func NewFileHandler(dirPath string) (*SafeFileHandler, error) {
	dirPath = filepath.Clean(dirPath)

	dirInfo, err := os.Lstat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("could not use directory [%s] %w", dirPath, err)
	}

	if !dirInfo.IsDir() {
		return nil, fmt.Errorf("provided path [%s] is not a directory", dirPath)
	}

	_, err = os.ReadDir(dirPath)
	if err != nil {
		return nil, ErrNoReadPerm
	}

	return &SafeFileHandler{
		dirPath: dirPath,
	}, nil
}

// Open is a simple wrapper around os.Open, but also joins the filename to the directory, cleans the path, checks the
// file exists.
func (h *SafeFileHandler) Open(filename string) (*os.File, error) {
	cleanPath := filepath.Clean(filepath.Join(h.dirPath, filename))

	info, err := os.Lstat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotExists
		}
		return nil, fmt.Errorf("could not get file info for file %s in directory %s: %w", filename, h.dirPath, err)
	}

	// check if owner readable
	if info.Mode().Perm()&0400 == 0 {
		return nil, ErrNoReadPerm
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s in directory %s: %w", filename, h.dirPath, err)
	}

	return file, nil
}

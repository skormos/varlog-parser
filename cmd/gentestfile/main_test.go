package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	tempDataDir  = "./testdata/temp"
	tempDataFile = "benchmark.log"
)

func TestMain(m *testing.M) {
	if isBench(os.Args) {
		if _, err := os.Stat(tempDataDir); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(tempDataDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else if err != nil {
			panic(err)
		}

		filename := filepath.Join(tempDataDir, tempDataFile)
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			if datafile, err := os.Create(filename); err != nil {
				panic(err)
			} else {
				defer func() {
					if err := datafile.Close(); err != nil {
						fmt.Println("while closing benchmark data file", err)
					}
				}()
				if err := populateBenchmarkLog(datafile, 1_000_000_000); err != nil {
					panic(err)
				}
			}
		}
	}

	os.Exit(m.Run())
}

func isBench(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-test.bench") {
			return true
		}
	}

	return false
}

const (
	dummyDataSingle = "Single process wrote this message."
	// reminder the back-ticks contain formatting
	dummyDataMulti = `Multi-line message starts with this line.
	The next line is here, which explains a bunch of things that are related to the log entry.
	Finally, this will just be an additional line to add even more context.`
)

// quick and dirty way to populate a test file without it taking up long-term space
func populateBenchmarkLog(file *os.File, size int) error {
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

package syslog

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    []*LogLine
		expectError bool
	}{
		"One single-line entry should return a valid record": {
			input: "Feb  1 10:11:12 the-host-name processname[23444]: This is the message written to the log",
			expected: []*LogLine{
				&LogLine{
					Timestamp:   time.Date(0, time.February, 1, 10, 11, 12, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         23444,
					ProcessName: "processname",
					Message:     "This is the message written to the log",
					Raw:         "Feb  1 10:11:12 the-host-name processname[23444]: This is the message written to the log",
				},
			},
		},
		"Multiple single-line entries should return multiple valid records": {
			input: `Feb  1 10:11:12 the-host-name processname[23444]: This is the message written to the log
Feb  1 10:11:14 the-host-name processname[23444]: This is the next message written from the same process`,
			expected: []*LogLine{
				&LogLine{
					Timestamp:   time.Date(0, time.February, 1, 10, 11, 12, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         23444,
					ProcessName: "processname",
					Message:     "This is the message written to the log",
					Raw:         "Feb  1 10:11:12 the-host-name processname[23444]: This is the message written to the log",
				},
				&LogLine{
					Timestamp:   time.Date(0, time.February, 1, 10, 11, 14, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         23444,
					ProcessName: "processname",
					Message:     "This is the next message written from the same process",
					Raw:         "Feb  1 10:11:14 the-host-name processname[23444]: This is the next message written from the same process",
				},
			},
		},
		"Single multi-line should return a single valid record": {
			input: `May 12 13:59:34 the-host-name processname[1098]: This is the beginning of the message
	followed by the second line of the message that doesn't have a date
	and ending with the rest of the message.`,
			expected: []*LogLine{
				&LogLine{
					Timestamp:   time.Date(0, time.May, 12, 13, 59, 34, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         1098,
					ProcessName: "processname",
					Message:     "This is the beginning of the message \tfollowed by the second line of the message that doesn't have a date \tand ending with the rest of the message.",
					Raw:         "May 12 13:59:34 the-host-name processname[1098]: This is the beginning of the message \tfollowed by the second line of the message that doesn't have a date \tand ending with the rest of the message.",
				},
			},
		},
		"Mixed single and multi-line entries should return valid records": {
			input: `Dec 30 22:59:58 the-host-name processname[2295]: This is the beginning of the message
	followed by the second line of the message that doesn't have a date
	and ending with the rest of the message.
Dec 30 23:00:02 the-host-name anotherprocess[39841]: A ping from another process called another process.
Dec 30 23:00:08 the-host-name processname[2295]: This starts a similar message
	followed by the second line of a similar message that doesn't
	have a date and ending with the rest of the message.`,
			expected: []*LogLine{
				&LogLine{
					Timestamp:   time.Date(0, time.December, 30, 22, 59, 58, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         2295,
					ProcessName: "processname",
					Message:     "This is the beginning of the message \tfollowed by the second line of the message that doesn't have a date \tand ending with the rest of the message.",
					Raw:         "Dec 30 22:59:58 the-host-name processname[2295]: This is the beginning of the message \tfollowed by the second line of the message that doesn't have a date \tand ending with the rest of the message.",
				},
				&LogLine{
					Timestamp:   time.Date(0, time.December, 30, 23, 0, 2, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         39841,
					ProcessName: "anotherprocess",
					Message:     "A ping from another process called another process.",
					Raw:         "Dec 30 23:00:02 the-host-name anotherprocess[39841]: A ping from another process called another process.",
				},
				&LogLine{
					Timestamp:   time.Date(0, time.December, 30, 23, 00, 8, 0, time.UTC),
					Host:        "the-host-name",
					Pid:         2295,
					ProcessName: "processname",
					Message:     "This starts a similar message \tfollowed by the second line of a similar message that doesn't \thave a date and ending with the rest of the message.",
					Raw:         "Dec 30 23:00:08 the-host-name processname[2295]: This starts a similar message \tfollowed by the second line of a similar message that doesn't \thave a date and ending with the rest of the message.",
				},
			},
		},
		"Incorrect timestamp should return error": {
			input:       "08-Aug 1:12:23.000 localhost processname[123]: Incorrect timestamp format.",
			expectError: true,
		},
		"Invalid PID format should return error": {
			input:       "Aug 08 13:13:13 localhost processname[PID]: PID is not a number.",
			expectError: true,
		},
		"Missing colon should return error": {
			input:       "Aug 08 13:13:13 localhost processname[123] Missing message separator.",
			expectError: true,
		},
		"No leading tab in multi-line should return error": {
			input: `Aug 08 13:13:13 localhost processname[123]: Missing message separator.
	this line has a tab, and should format correctly, but
this following line does not and should fail.`,
			expectError: true,
		},
	}

	target, err := NewSyslogParser()
	require.NoError(t, err)

	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			actual, actualErr := target.Parse(context.TODO(), strings.NewReader(test.input))
			if test.expectError {
				require.Nil(tt, actual)
				require.Error(tt, actualErr)
			} else {
				require.NoError(tt, actualErr)
				assert.Equal(tt, test.expected, actual)
			}
		})
	}
}

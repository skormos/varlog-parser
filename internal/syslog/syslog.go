package syslog

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	headerPattern = `(?P<Timestamp>[A-Z][a-z]{2}\s[\s|0-9][0-9]\s[0-9]{2}:[0-9]{2}:[0-9]{2})\s(?P<Host>.+)\s(?P<ProcessName>.+)\[(?P<PID>\d+)\]:\s+(?P<Message>.+)`
)

const (
	headerFieldRaw = iota
	headerFieldTimestamp
	headerFieldHost
	headerFieldProcessName
	headerFieldPID
	headerFieldMessage
)

type (
	LogLine struct {
		Timestamp   time.Time
		Host        string
		Pid         uint
		ProcessName string
		Message     string
		Raw         string
	}

	Parser struct {
		headerExp *regexp.Regexp
	}
)

func NewSyslogParser() (*Parser, error) {
	exp, err := regexp.Compile(headerPattern)

	if err != nil {
		return nil, fmt.Errorf("could not initialize Log Header RegExp pattern: %w", err)
	}

	return &Parser{
		headerExp: exp,
	}, nil
}

func (p *Parser) Parse(ctx context.Context, reader io.Reader) ([]*LogLine, error) {
	scanner := bufio.NewScanner(reader)

	lines := make([]*LogLine, 0)

	for scanner.Scan() {
		raw := scanner.Text()
		if strings.HasPrefix(raw, string('\t')) {
			line := lines[len(lines)-1]
			line.Raw += " " + raw
			line.Message += " " + raw
		} else {
			line, err := p.parseLine(raw)
			if err != nil {
				return nil, fmt.Errorf("could not parse line [%s]: %w", raw, err)
			}
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("while scanning lines: %w", err)
	}

	return lines, nil
}

func (p *Parser) parseLine(raw string) (*LogLine, error) {
	fields := p.headerExp.FindStringSubmatch(raw)
	if fields == nil || len(fields) != len(p.headerExp.SubexpNames()) {
		return nil, fmt.Errorf("incorrect header format for input: [%s]", raw)
	}

	pid, err := strconv.Atoi(fields[headerFieldPID])
	if err != nil {
		return nil, fmt.Errorf("while converting the PID [%s] to an int: %w", fields[headerFieldPID], err)
	}

	timestamp, err := time.Parse(time.Stamp, fields[headerFieldTimestamp])
	if err != nil {
		return nil, fmt.Errorf("while converting the timestamp [%s] to time: %w", fields[headerFieldTimestamp], err)
	}

	return &LogLine{
		Timestamp:   timestamp,
		Host:        fields[headerFieldHost],
		Pid:         uint(pid),
		ProcessName: fields[headerFieldProcessName],
		Message:     fields[headerFieldMessage],
		Raw:         fields[headerFieldRaw],
	}, nil
}

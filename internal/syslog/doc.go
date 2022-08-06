// Package syslog is used to parse logger file content as an opened Reader. It is up to callers to create and handle
// the reader instance correctly.
//
// Expected Usage:
//    reader := <instantiate io.Reader>
//    parser, err := NewSyslogParser()
//    <handle error>
//    outputSlice, err := parser.Parse(ctx, reader)
//    <handle error>
//    <iterate over outputSlice>
package syslog

// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package sanefmt

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

var (
	_ func(r io.Reader) ([]byte, error) = Format
	_ func() (string, error)            = Version
)

type Error struct {
	Text, Details string
}

func (e Error) Error() string {
	return e.Text
}

func parseStderr(s []byte, fallback error) error {
	var err Error
	var details []byte
	s = bytes.TrimPrefix(s, []byte("ERROR:"))
	first := true
	for line := range bytes.Lines(s) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if first {
			const sep = "file:///STDIN.ts:"
			if i := bytes.Index(line, []byte(sep)); i != -1 {
				line = append(line[:i], line[i+len(sep):]...)
			}
			original, originalSize := utf8.DecodeRune(line)
			lower := unicode.ToLower(original)
			lowerSize := utf8.RuneLen(lower)
			switch {
			case originalSize == lowerSize:
				utf8.EncodeRune(line, lower)
			case originalSize > lowerSize:
				utf8.EncodeRune(line, lower)
				line = append(line[:lowerSize], line[originalSize:]...)
			default:
				rest := line[originalSize:]
				line = make([]byte, 0, lowerSize+len(rest))
				line = utf8.AppendRune(line, lower)
				line = append(line, rest...)
			}
			err.Text = string(line)
		} else {
			details = append(details, line...)
			details = append(details, '\n')
		}
		first = false
	}
	if err.Text == "" {
		return fallback
	}
	err.Details = string(bytes.TrimSuffix(details, []byte{'\n'}))
	return &err
}

type buffer []byte

func (b *buffer) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

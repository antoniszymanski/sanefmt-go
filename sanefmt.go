// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package sanefmt

import (
	"bytes"
	"io"
	"math"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

func Format(r io.Reader) ([]byte, error) {
	return format(r)
}

func Version() (string, error) {
	return version()
}

type Error struct {
	Text, Details string
}

func (e Error) Error() string {
	return e.Text
}

func parseStderr(b []byte, fallback error) error {
	var err Error
	var detailsLines [][]byte
	minIndent := math.MaxInt
	isFirst := true
	b = bytes.TrimPrefix(b, []byte("ERROR:"))
	for line := range bytes.Lines(b) {
		if isFirst {
			line = bytes.TrimSpace(line)
		} else {
			line = bytes.TrimRightFunc(line, unicode.IsSpace)
		}
		switch {
		case len(line) == 0:
			continue
		case isFirst:
			const sep = "file:///STDIN.ts:"
			if i := bytes.Index(line, []byte(sep)); i != -1 {
				line = append(line[:i], line[i+len(sep):]...)
			}
			if len(line) == 0 {
				return fallback
			}
			err.Text = replaceFirstRune(line, unicode.ToLower)
			isFirst = false
		default:
			for i := range line {
				if line[i] == '\t' {
					line[i] = ' '
				}
			}
			minIndent = min(minIndent, bytes.IndexFunc(line, func(r rune) bool { return !unicode.IsSpace(r) }))
			detailsLines = append(detailsLines, line)
		}
	}
	if minIndent != math.MaxInt {
		for i := range detailsLines {
			detailsLines[i] = detailsLines[i][minIndent:]
		}
	}
	err.Details = string(bytes.Join(detailsLines, []byte("\n")))
	return &err
}

func replaceFirstRune(b []byte, f func(rune) rune) string {
	if len(b) == 0 {
		return ""
	}
	rune, size := utf8.DecodeRune(b)
	mapped := f(rune)
	mappedSize := utf8.RuneLen(mapped)
	if mappedSize == -1 {
		mappedSize = len(string(utf8.RuneError))
	}
	remaining := b[size:]
	out := make([]byte, 0, mappedSize+len(remaining))
	out = utf8.AppendRune(out, mapped)
	out = append(out, remaining...)
	return unsafe.String(unsafe.SliceData(out), len(out))
}

type buffer []byte

func (b *buffer) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

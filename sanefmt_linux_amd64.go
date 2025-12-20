// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package sanefmt

import (
	"bytes"
	_ "embed"
	"io"
	"strings"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed resources/sane-fmt-x86_64-unknown-linux-musl
var exe []byte

func Format(r io.Reader) ([]byte, error) {
	var stdout bytes.Buffer
	var stderr buffer
	cmd := exec.Command(exe, "--stdio")
	cmd.Stdin = r
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, parseStderr(stderr, err)
	}
	return stdout.Bytes(), nil
}

func Version() (string, error) {
	var stdout strings.Builder
	var stderr buffer
	cmd := exec.Command(exe, "--version")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", parseStderr(stderr, err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

//go:build linux && amd64

/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package sanefmt

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"strings"

	"codeberg.org/msantos/embedexe/exec"
)

//go:embed resources/sane-fmt-x86_64-unknown-linux-musl
var exe []byte

func format(r io.Reader) (*bytes.Buffer, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(exe, "--stdio")
	cmd.Stdin = r
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New(stderr.String())
	}

	return &stdout, nil
}

func version() (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(exe, "--version")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

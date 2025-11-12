// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package sanefmt

import (
	"errors"
	"io"
	"strings"
)

var (
	_ func(r io.Reader) ([]byte, error) = Format
	_ func() (string, error)            = Version
)

func newError(s string) error {
	return errors.New(strings.TrimSpace(s))
}

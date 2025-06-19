//go:build !(linux && amd64)

/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package sanefmt

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	wasip1 "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed resources/sane-fmt-wasm32-wasi.wasm
var binary []byte

var (
	ctx         = context.Background()
	rt          wazero.Runtime
	compiled    wazero.CompiledModule
	initRuntime = sync.OnceFunc(func() {
		rt = wazero.NewRuntime(ctx)
		var err error
		if _, err = wasip1.Instantiate(ctx, rt); err != nil {
			panic(err)
		}
		compiled, err = rt.CompileModule(ctx, binary)
		if err != nil {
			panic(err)
		}
	})
)

func Format(r io.Reader) ([]byte, error) {
	initRuntime()

	var stdout bytes.Buffer
	var stderr strings.Builder
	config := wazero.NewModuleConfig().
		WithStdin(r).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithArgs("sane-fmt", "--stdio")
	m, err := rt.InstantiateModule(ctx, compiled, config)
	if err != nil {
		return nil, errors.New(stderr.String())
	}
	defer m.Close(ctx) //nolint:errcheck

	return stdout.Bytes(), nil
}

func Version() (string, error) {
	initRuntime()

	var stdout, stderr strings.Builder
	config := wazero.NewModuleConfig().
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithArgs("sane-fmt", "--version")
	m, err := rt.InstantiateModule(ctx, compiled, config)
	if err != nil {
		return "", errors.New(stderr.String())
	}
	defer m.Close(ctx) //nolint:errcheck

	return strings.TrimSpace(stdout.String()), nil
}

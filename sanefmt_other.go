// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

//go:build !(linux && amd64)

package sanefmt

import (
	"bytes"
	"context"
	_ "embed"
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
	runtime     wazero.Runtime
	compiled    wazero.CompiledModule
	err         error
	initRuntime = sync.OnceFunc(func() {
		runtime = wazero.NewRuntime(ctx)
		if _, err = wasip1.Instantiate(ctx, runtime); err == nil {
			compiled, err = runtime.CompileModule(ctx, binary)
		}
	})
)

func Format(r io.Reader) ([]byte, error) {
	if initRuntime(); err != nil {
		return nil, err
	}
	var stdout bytes.Buffer
	var stderr strings.Builder
	config := wazero.NewModuleConfig().
		WithStdin(r).
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithArgs("sane-fmt", "--stdio")
	module, err := runtime.InstantiateModule(ctx, compiled, config)
	if err != nil {
		return nil, newError(stderr.String())
	}
	defer module.Close(ctx) //nolint:errcheck
	return stdout.Bytes(), nil
}

func Version() (string, error) {
	if initRuntime(); err != nil {
		return "", err
	}
	var stdout, stderr strings.Builder
	config := wazero.NewModuleConfig().
		WithStdout(&stdout).
		WithStderr(&stderr).
		WithArgs("sane-fmt", "--version")
	module, err := runtime.InstantiateModule(ctx, compiled, config)
	if err != nil {
		return "", newError(stderr.String())
	}
	defer module.Close(ctx) //nolint:errcheck
	return strings.TrimSpace(stdout.String()), nil
}

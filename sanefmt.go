package sanefmt

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed resources/sane-fmt-wasm32-wasi.wasm
var binary []byte

var (
	ctx         = context.Background()
	rt          wazero.Runtime
	compiled    wazero.CompiledModule
	initRuntime = sync.OnceFunc(func() {
		rt = wazero.NewRuntime(ctx)
		wasi_snapshot_preview1.MustInstantiate(ctx, rt)
		var err error
		compiled, err = rt.CompileModule(ctx, binary)
		if err != nil {
			panic(err)
		}
	})
)

func Format(r io.Reader) (*bytes.Buffer, error) {
	initRuntime()

	var stdout, stderr bytes.Buffer
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

	return &stdout, nil
}

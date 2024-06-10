package datatransform

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudquery/plugin-sdk/v4/glob"
	"github.com/rs/zerolog/log"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"golang.org/x/exp/maps"
)

type Wazero struct {
	wazero.Runtime
	modules map[string]*WazeroModule
}

type WazeroModule struct {
	api.Module
}

var _ DataTransformer = (*Wazero)(nil)

func NewWazeroDataTransformer(ctx context.Context) (*Wazero, error) {
	cfg := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)
	rt := wazero.NewRuntimeWithConfig(ctx, cfg)
	return &Wazero{Runtime: rt, modules: map[string]*WazeroModule{}}, nil
}

func (w *Wazero) Close() error      { return w.Runtime.Close(context.Background()) }
func (w *Wazero) Modules() []string { return maps.Keys(w.modules) }

func (w *Wazero) InitializeModule(ctx context.Context, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to load %s: %w", path, err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, w.Runtime)

	compiled, err := w.CompileModule(ctx, content)
	if err != nil {
		return err
	}
	host := w.NewHostModuleBuilder("env")
	builder := host.NewFunctionBuilder()

	builder.WithFunc(func(ctx context.Context, m api.Module, offset, byteCount uint32) {
		buf, ok := m.Memory().Read(offset, byteCount)
		if !ok {
			log.Panic().Msgf("Memory.Read(%d, %d) out of range", offset, byteCount)
		}

		log.Info().Str("wasm_transformer", m.Name()).Msg(string(buf))
	}).Export("log")

	if _, err := host.Instantiate(ctx); err != nil {
		return err
	}

	cm, err := w.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithArgs("wasm-transform").WithStderr(os.Stderr).WithStdout(os.Stdout))
	if err != nil {
		return err
	}
	mod := &WazeroModule{Module: cm}
	w.modules[path] = mod
	return nil
}

const cqFunctionPrefix = "_cqtransform_"

func (w *Wazero) ExecuteModule(ctx context.Context, path string, table string, rec []byte) ([]byte, error) {
	mod, ok := w.modules[path]
	if !ok {
		return nil, fmt.Errorf("requested unknown module %s", path)
	}

	matchingFunctionsWithFilter := map[string]string{}

	for _, wf := range mod.ExportedFunctionDefinitions() {
		if !strings.HasPrefix(wf.Name(), cqFunctionPrefix) {
			continue
		}

		table := wf.Name()[len(cqFunctionPrefix):]
		stop := 0
		for idx, ch := range table {
			if ch == '@' && len(table) > idx+1 && table[idx+1] == '@' {
				stop = idx
				break
			}
		}

		if stop == 0 {
			continue
		}
		filterEnd := table[:stop]

		matchingFunctionsWithFilter[wf.Name()] = strings.ReplaceAll(filterEnd, "\"", "")
	}

	if len(matchingFunctionsWithFilter) == 0 {
		return rec, nil
	}

	for name, filter := range matchingFunctionsWithFilter {
		if !glob.Glob(filter, table) {
			continue
		}

		out, err := doTransform(ctx, mod, name, rec)
		if err := checkedExitZero(err); err != nil {
			return nil, err
		}

		if out != nil {
			rec = out
		}
	}

	return rec, nil
}

func doTransform(ctx context.Context, mod *WazeroModule, function string, rec []byte) ([]byte, error) {
	rv, err := mod.ExportedFunction("allocate").Call(ctx, uint64(len(rec)))
	if err != nil {
		return nil, err
	}

	if !mod.Memory().Write(uint32(rv[0]), rec) {
		return nil, fmt.Errorf("couldn't write to memory offset: %d", rv[0])
	}
	ret, err := mod.ExportedFunction(function).Call(ctx,
		rv[0],
		uint64(len(rec)),
	)

	if err := checkedExitZero(err); err != nil {
		return nil, err
	}

	outPtr := uint32(ret[0] >> 32)
	outSize := uint32(ret[0])

	bo, ok := mod.Memory().Read(outPtr, outSize)
	if !ok {
		return nil, fmt.Errorf("failed to read output buffer %d (size: %d)", outPtr, outSize)
	}

	return bo, nil
}

func checkedExitZero(err error) error {
	if err == nil {
		return nil
	}

	if exitErr, ok := err.(*sys.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			return nil
		}
	}

	return err
}

package datatransform

import (
	"context"
	"fmt"
	"log"
	"os"

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
			log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
		}
		fmt.Println("[rust] " + string(buf))
	}).Export("log")

	host.Instantiate(ctx)

	cm, err := w.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithArgs("wasm-transform").WithStderr(os.Stderr).WithStdout(os.Stdout))
	if err != nil {
		return err
	}
	mod := &WazeroModule{Module: cm}
	w.modules[path] = mod
	return nil
}

func (w *Wazero) ExecuteModule(ctx context.Context, path string, rec []byte) ([]byte, error) {
	mod, ok := w.modules[path]
	if !ok {
		return nil, fmt.Errorf("requested unknown module %s", path)
	}

	fmt.Printf("%s exported memory table:\n", path)
	for _, mdef := range mod.ExportedMemoryDefinitions() {
		fmt.Println(mdef.ExportNames()[0])
	}
	fmt.Printf("\n%s exported functions table:\n", path)
	for _, fdef := range mod.ExportedFunctionDefinitions() {
		fmt.Println(fdef.Name(), fdef.ParamTypes(), fdef.ResultTypes())
	}
	fmt.Println()

	rv, err := mod.ExportedFunction("allocate").Call(ctx, uint64(len(rec)))
	if err != nil {
		return nil, err
	}

	if !mod.Memory().Write(uint32(rv[0]), rec) {
		return nil, fmt.Errorf("couldn't write to memory offset: %d", rv[0])
	}
	ret, err := mod.ExportedFunction("cloudquery_transform").Call(ctx,
		rv[0],
		uint64(len(rec)),
	)

	fmt.Println(ret)

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

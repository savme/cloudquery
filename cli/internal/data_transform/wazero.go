package datatransform

// type Wazero struct {
// 	wazero.Runtime
// 	modules map[string]wazero.CompiledModule
// }

// var _ DataTransformer = (*Wazero)(nil)

// func NewWazeroDataTransform(ctx context.Context) (*Wazero, error) {
// 	cfg := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)
// 	rt := wazero.NewRuntimeWithConfig(ctx, cfg)

// 	_, err := wasi_snapshot_preview1.Instantiate(ctx, rt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Wazero{Runtime: rt, modules: map[string]wazero.CompiledModule{}}, nil
// }

// func (w *Wazero) Close() error      { return w.Close() }
// func (w *Wazero) Modules() []string { return maps.Keys(w.modules) }

// func (w *Wazero) InitializeModule(ctx context.Context, path string) error {
// 	content, err := os.ReadFile(path)
// 	if err != nil {
// 		return fmt.Errorf("unable to load %s: %w", path, err)
// 	}
// 	mod, err := w.CompileModule(ctx, content)
// 	if err != nil {
// 		return err
// 	}
// 	w.modules[path] = mod
// 	return nil
// }

// func (w *Wazero) ExecuteModule(path string) error {
// 	mod, ok := w.modules[path]
// 	if !ok {
// 		return fmt.Errorf("requested unknown module %s", path)
// 	}

// 	fn, ok := mod.ExportedFunctions()["main"]
// 	if !ok {
// 		return fmt.Errorf("requested unknown function %s in module %s", "main", path)
// 	}

// }

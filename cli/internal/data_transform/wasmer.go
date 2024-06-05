package datatransform

import (
	"errors"
	"fmt"
	"os"

	"github.com/wasmerio/wasmer-go/wasmer"
	"golang.org/x/exp/maps"
)

type Wasmer struct {
	store   *wasmer.Store
	modules map[string]*WasiModule
}

type WasiModule struct {
	*wasmer.Instance
	*wasmer.WasiEnvironment
}

var _ DataTransformer = (*Wasmer)(nil)

func getRuntimeEngineKind() (wasmer.EngineKind, error) {
	if wasmer.IsEngineAvailable(wasmer.UNIVERSAL) {
		return wasmer.UNIVERSAL, nil
	}

	if wasmer.IsEngineAvailable(wasmer.DYLIB) {
		return wasmer.DYLIB, nil
	}

	return wasmer.EngineKind(0), errors.New("no available wasm engines")
}

func getRuntimeCompilerKind() (wasmer.CompilerKind, error) {
	compilersByPriority := []wasmer.CompilerKind{
		wasmer.CRANELIFT,
		wasmer.LLVM,
		wasmer.SINGLEPASS,
	}

	for _, kind := range compilersByPriority {
		if !wasmer.IsCompilerAvailable(kind) {
			continue
		}

		return kind, nil
	}

	return wasmer.CompilerKind(0), errors.New("no available wasm compilers")
}

func NewWasmerDataTransformer() (*Wasmer, error) {
	engineKind, err := getRuntimeEngineKind()
	if err != nil {
		return nil, err
	}
	compilerKind, err := getRuntimeCompilerKind()
	if err != nil {
		return nil, err
	}

	engineConfig := wasmer.NewConfig()
	switch engineKind {
	case wasmer.DYLIB:
		engineConfig.UseDylibEngine()
	default:
		engineConfig.UseJITEngine()
	}

	switch compilerKind {
	case wasmer.CRANELIFT:
		engineConfig.UseCraneliftCompiler()
	case wasmer.LLVM:
		engineConfig.UseLLVMCompiler()
	default:
		engineConfig.UseSinglepassCompiler()
	}

	engine := wasmer.NewEngineWithConfig(engineConfig)
	store := wasmer.NewStore(engine)
	return &Wasmer{store: store, modules: map[string]*WasiModule{}}, nil
}

func (w *Wasmer) Close() error {
	w.store.Close()
	return nil
}

func (w *Wasmer) Modules() []string { return maps.Keys(w.modules) }

func (w *Wasmer) InitializeModule(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to load %s: %w", path, err)
	}

	mod, err := wasmer.NewModule(w.store, content)
	if err != nil {
		return err
	}

	env, err := wasmer.NewWasiStateBuilder(path).CaptureStderr().CaptureStdout().Finalize()
	if err != nil {
		return err
	}
	imports, err := env.GenerateImportObject(w.store, mod)
	if err != nil {
		return err
	}

	instance, err := wasmer.NewInstance(mod, imports)
	if err != nil {
		return err
	}

	w.modules[path] = &WasiModule{Instance: instance, WasiEnvironment: env}
	return nil
}

func (w *Wasmer) ExecuteModule(path string) error {
	mod, ok := w.modules[path]
	if !ok {
		return fmt.Errorf("requested unknown module %s", path)
	}

	start, err := mod.Exports.GetWasiStartFunction()
	if err != nil {
		return err
	}
	rv, err := start()
	fmt.Println(rv, string(mod.ReadStdout()), string(mod.ReadStderr()))

	// if _, ok := err.(*wasmer.TrapError); ok {
	// 	return nil
	// }
	return err
}

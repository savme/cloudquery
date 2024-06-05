package specs

import "errors"

type Transform struct {
	Name string `json:"name" jsonschema:"required,minLength=1"`
	Path string `json:"path"`
}

func (*Transform) GetWarnings() Warnings {
	warnings := make(map[string]string)
	return warnings
}

func (*Transform) SetDefaults() {}

func (t *Transform) Validate() error {
	if t.Name == "" {
		return errors.New("name is required")
	}

	if t.Path == "" {
		return errors.New("path to a wasm module is required")
	}

	return nil
}

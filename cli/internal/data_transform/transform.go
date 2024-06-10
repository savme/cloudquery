package datatransform

import (
	"context"
)

type DataTransformer interface {
	Modules() []string
	InitializeModule(context.Context, string) error
	ExecuteModule(context.Context, string, string, []byte) ([]byte, error)

	Close() error
}

type NoopDataTransformer struct{}

var _ DataTransformer = (*NoopDataTransformer)(nil)

func NewNoopDataTransformer() (*NoopDataTransformer, error) {
	return &NoopDataTransformer{}, nil
}

func (*NoopDataTransformer) Close() error { return nil }
func (*NoopDataTransformer) ExecuteModule(context.Context, string, string, []byte) ([]byte, error) {
	return nil, nil
}
func (*NoopDataTransformer) InitializeModule(context.Context, string) error { return nil }
func (*NoopDataTransformer) Modules() []string                              { return []string{} }

package datatransform

type DataTransformer interface {
	Modules() []string
	InitializeModule(string) error
	ExecuteModule(string) error

	Close() error
}

type NoopDataTransformer struct{}

var _ DataTransformer = (*NoopDataTransformer)(nil)

func NewNoopDataTransformer() (*NoopDataTransformer, error) {
	return &NoopDataTransformer{}, nil
}

func (n *NoopDataTransformer) Close() error                  { return nil }
func (n *NoopDataTransformer) ExecuteModule(string) error    { return nil }
func (n *NoopDataTransformer) InitializeModule(string) error { return nil }
func (n *NoopDataTransformer) Modules() []string             { return []string{} }

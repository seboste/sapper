package ports

type ParameterResolver interface {
	Resolve(name string) string
}

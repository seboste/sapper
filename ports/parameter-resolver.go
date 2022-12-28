package ports

type ParameterResolver interface {
	Resolve(name string, defaultValue string) string
}

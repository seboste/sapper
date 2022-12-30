package ports

type ServiceApi interface {
	Add(templateName string, parentDir string, parameterResolver ParameterResolver) error
	Describe(path string) (string, error)
	Upgrade(path string) error
	Build(path string) error
	Test()
	Deploy()
}

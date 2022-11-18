package ports

type ServiceApi interface {
	Add(templateName string, parentDir string, parameterResolver ParameterResolver) error
	Update()
	Build(path string) error
	Test()
	Deploy()
}

package ports

import "io"

type ServiceApi interface {
	Add(templateName string, parentDir string, parameterResolver ParameterResolver) error
	Describe(path string, writer io.Writer) error
	Upgrade(path string) error
	Build(path string) error
	Test()
	Deploy()
}

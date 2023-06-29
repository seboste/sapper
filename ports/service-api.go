package ports

import "io"

type ServiceApi interface {
	Add(templateName string, parentDir string, parameterResolver ParameterResolver) (Service, error)
	Describe(path string, writer io.Writer) error
	Upgrade(path string, keepMajorVersion bool) error
	Build(path string) (string, error)
	Test(path string) error
	Deploy(path string) error
	Run(path string) error
}

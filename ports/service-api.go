package ports

type ServiceApi interface {
	Add(name string, version string, templateName string, parentDir string) error
	Update()
	Build(path string) error
	Test()
	Deploy()
}

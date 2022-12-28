package ports

type BrickDependency struct {
	Id      string
	Version string
}

type Service struct {
	Id           string
	Path         string `yaml:"-"`
	BrickIds     []BrickDependency
	Dependencies []PackageDependency `yaml:"-"`
}

type ServicePersistence interface {
	Load(path string) (Service, error)
	Save(service Service) error
}

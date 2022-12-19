package ports

type ServiceDependency struct {
	Id      string
	Version string
}

type BrickDependency struct {
	Id      string
	Version string
}

type Service struct {
	Id           string
	Path         string `yaml:"-"`
	BrickIds     []BrickDependency
	Dependencies []ServiceDependency `yaml:"-"`
}

type ServicePersistence interface {
	Load(path string) (Service, error)
	Save(service Service) error
}

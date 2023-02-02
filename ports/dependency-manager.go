package ports

type PackageDependency struct {
	Id      string
	Version string
}

type PackageDependencyReader interface {
	ReadFromService(s Service) ([]PackageDependency, error)
	ReadFromBrick(b Brick) ([]PackageDependency, error)
}

type PackageDependencyWriter interface {
	WriteToService(s Service, dependency string, version string) error
}

type DependencyInfo interface {
	AvailableVersions(dependency string) ([]string, error)
}

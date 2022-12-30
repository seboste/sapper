package ports

type PackageDependency struct {
	Id      string
	Version string
}

type DependencyReader interface {
	Read(s Service) ([]PackageDependency, error)
}

type DependencyWriter interface {
	Write(s Service) error
}

type DependencyInfo interface {
	AvailableVersions(dependency string) ([]string, error)
}

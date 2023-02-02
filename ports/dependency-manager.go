package ports

type PackageDependency struct {
	Id      string
	Version string
}

type PacakgeDependencySectionPredicate func(line string, state string) (isActive bool, newState string)

type PackageDependencyReader interface {
	ReadFromService(s Service) ([]PackageDependency, error)
	ReadFromBrick(b Brick, p PacakgeDependencySectionPredicate) ([]PackageDependency, error)
}

type PackageDependencyWriter interface {
	WriteToService(s Service, dependency string, version string) error
}

type DependencyInfo interface {
	AvailableVersions(dependency string) ([]string, error)
}

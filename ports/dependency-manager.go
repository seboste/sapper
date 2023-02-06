package ports

type PackageDependency struct {
	Id      string
	Version string
}

type PackageDependencySectionPredicate func(line string, state string) (isActive bool, newState string)

type ServicePackageDependencyReader interface {
	ReadFromService(s Service) ([]PackageDependency, error)
}

type BrickPackageDependencyReader interface {
	ReadFromBrick(b Brick, p PackageDependencySectionPredicate) ([]PackageDependency, error)
}

type ServicePackageDependencyWriter interface {
	WriteToService(s Service, d PackageDependency) error
}

type BrickPackageDependencyWriter interface {
	WriteToBrick(b Brick, d PackageDependency, p PackageDependencySectionPredicate) error
}

type DependencyInfo interface {
	AvailableVersions(dependency string) ([]string, error)
}

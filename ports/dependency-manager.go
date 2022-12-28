package ports

type PackageDependency struct {
	Id      string
	Version string
}

type DependencyReader interface {
	Read(s Service) []PackageDependency
}

type DependencyWriter interface {
	Write(s Service)
}

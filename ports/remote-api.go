package ports

type Remote struct {
	Name string
	Path string
}

type RemoteApi interface {
	Add()
	List()
}

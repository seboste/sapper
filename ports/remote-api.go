package ports

type Remote struct {
	Path string
}

type RemoteApi interface {
	Add()
	List()
}

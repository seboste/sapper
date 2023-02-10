package ports

type RemoteKind int

const (
	FilesystemRemote RemoteKind = iota
	GitRemote
)

type Remote struct {
	Name string
	Kind RemoteKind
	Src  string //folder path for file system remotes and url for git repositories
}

type RemoteApi interface {
	Add(name string, src string, position int) error
	Remove(name string) error
	Update(name string) error
	Upgrade(name string) error
	List() []Remote
}

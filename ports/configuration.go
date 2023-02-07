package ports

type Configuration interface {
	ConfigurationDir() string
	Remotes() []Remote
	UpdateRemotes(remotes []Remote) error
}

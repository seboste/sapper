package ports

type Configuration interface {
	Save() error
	DefaultRemotesDir() string
	Remotes() []Remote
	UpdateRemotes(remotes []Remote)
}

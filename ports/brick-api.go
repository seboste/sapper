package ports

type BrickApi interface {
	Add()
	List() []Brick
	Search(term string) []Brick
}

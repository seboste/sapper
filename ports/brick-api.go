package ports

type BrickApi interface {
	Add(servicePath string, brickId string, parameterResolver ParameterResolver) error
	Upgrade(brickId string) error
	List() []Brick
	Search(term string) []Brick
}

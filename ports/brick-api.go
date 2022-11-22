package ports

type BrickApi interface {
	Add(servicePath string, brickId string, parameterResolver ParameterResolver) error
	List() []Brick
	Search(term string) []Brick
}

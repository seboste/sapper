package ports

type BrickUpgrader interface {
	UpgradeInDB(brickId string, db BrickDB) error
}

type BrickApi interface {
	Add(servicePath string, brickId string, parameterResolver ParameterResolver) error
	Upgrade(brickId string) error
	List() []Brick
	Search(term string) []Brick
}

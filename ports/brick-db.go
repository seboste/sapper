package ports

type BrickParameters struct {
	Name    string
	Default string
}

type BrickKind int

const (
	Template = iota
	Extension
)

type Brick interface {
	GetID() string
	GetDescription() string
	GetVersion() string
	GetKind() BrickKind
	GetParameters() BrickParameters
	GetDependencies() string
	GetFiles() []string
}

type BrickDB interface {
	Init(Path string)
	GetBricks(kind BrickKind) []Brick
	GetBrick(id string) Brick
}

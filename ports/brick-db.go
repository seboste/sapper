package ports

type BrickParameters struct {
	Name    string
	Default string
}

type BrickKind int

const (
	Template BrickKind = iota
	Extension
)

type Brick interface {
	GetId() string
	GetDescription() string
	GetVersion() string
	GetKind() BrickKind
	GetParameters() []BrickParameters
	GetDependencies() []string
	GetFiles() []string
}

type BrickDB interface {
	Init(Path string) error
	Bricks(kind BrickKind) []Brick
	Brick(id string) Brick
}

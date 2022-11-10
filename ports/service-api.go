package ports

type ServiceApi interface {
	Add()
	Update()
	Build()
	Test()
	Deploy()
}

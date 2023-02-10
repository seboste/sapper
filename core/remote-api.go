package core

import (
	"github.com/seboste/sapper/ports"
)

type RemoteApi struct {
}

func (r RemoteApi) Add(name string, src string, position int) error {
	return nil
}
func (r RemoteApi) Remove(name string) error {
	return nil
}
func (r RemoteApi) Update(name string) error {
	return nil
}
func (r RemoteApi) Upgrade(name string) error {
	return nil
}
func (r RemoteApi) List() []ports.Remote {
	return nil
}

var _ ports.RemoteApi = RemoteApi{}

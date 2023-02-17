package core

import (
	"fmt"
	"os"

	"github.com/seboste/sapper/ports"
)

type RemoteApi struct {
	Configuration  ports.Configuration
	BrickDBFactory ports.BrickDBFactory
	BrickUpgrader  ports.BrickUpgrader
}

func findRemote(remotes []ports.Remote, name string) (index int, remote ports.Remote, ok bool) {
	for index, remote = range remotes {
		if remote.Name == name {
			ok = true
			return
		}
	}
	remote = ports.Remote{}
	index = -1
	ok = false
	return
}

func inferKind(src string) (kind ports.RemoteKind, err error) {
	if src[len(src)-4:] == ".git" {
		return ports.GitRemote, nil
	}
	if fileInfo, err := os.Stat(src); err == nil && fileInfo.IsDir() {
		return ports.FilesystemRemote, nil
	}
	return -1, fmt.Errorf("invalid git url or path on filesystem")
}

func Add(remotes []ports.Remote, r ports.Remote, pos int) []ports.Remote {
	if pos < 0 || pos >= len(remotes) {
		return append(remotes, r)
	}

	remotes = append(remotes[:pos+1], remotes[pos:]...)
	remotes[pos] = r
	return remotes
}

func (r RemoteApi) Add(name string, src string, position int) error {
	remotes := r.Configuration.Remotes()

	if _, _, ok := findRemote(remotes, name); ok {
		return fmt.Errorf("remote with name %s does already exist", name)
	}

	kind, err := inferKind(src)
	if err != nil {
		return err
	}

	remote := ports.Remote{Name: name, Src: src, Kind: kind}
	if _, err := r.BrickDBFactory.MakeBrickDB(remote, r.Configuration.DefaultRemotesDir()); err != nil {
		return err
	}

	remotes = Add(remotes, remote, position)

	r.Configuration.UpdateRemotes(remotes)

	if err := r.Configuration.Save(); err != nil {
		return err
	}

	return nil
}

func (r RemoteApi) Remove(name string) error {
	remotes := r.Configuration.Remotes()
	if i, _, ok := findRemote(remotes, name); ok {
		remotes = append(remotes[:i], remotes[i+1:]...)
		r.Configuration.UpdateRemotes(remotes)

		if err := r.Configuration.Save(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("remote %s does not exist", name)
}
func (r RemoteApi) Update(name string) error {
	_, remote, ok := findRemote(r.Configuration.Remotes(), name)
	if !ok {
		return fmt.Errorf("remote %s does not exist", name)
	}

	brickDB, err := r.BrickDBFactory.MakeBrickDB(remote, r.Configuration.DefaultRemotesDir())
	if err != nil {
		return err
	}

	return brickDB.Update()
}
func (r RemoteApi) Upgrade(name string) error {
	_, remote, ok := findRemote(r.Configuration.Remotes(), name)
	if !ok {
		return fmt.Errorf("remote %s does not exist", name)
	}

	brickDB, err := r.BrickDBFactory.MakeBrickDB(remote, r.Configuration.DefaultRemotesDir())
	if err != nil {
		return err
	}

	errorCount := 0
	allBricks := []ports.Brick{}
	for _, k := range ports.BrickKinds {
		allBricks = append(allBricks, brickDB.Bricks(k)...)
	}

	for i, brick := range allBricks {
		fmt.Printf("upgrading %s (%v/%v)...\n", brick.Id, i+1, len(allBricks))
		err := r.BrickUpgrader.UpgradeInDB(brick.Id, brickDB)
		if err != nil {
			errorCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("%v of %v bricks failed to upgrade.", errorCount, len(allBricks))
	}
	return nil
}
func (r RemoteApi) List() []ports.Remote {
	return r.Configuration.Remotes()
}

var _ ports.RemoteApi = RemoteApi{}

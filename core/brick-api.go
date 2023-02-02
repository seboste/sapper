package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/seboste/sapper/ports"
	"github.com/seboste/sapper/utils"
)

type BrickApi struct {
	Db                      ports.BrickDB
	ServicePersistence      ports.ServicePersistence
	PackageDependencyReader ports.PackageDependencyReader
	DependencyInfo          ports.DependencyInfo
	ServiceApi              ports.ServiceApi
}

func removeBricks(bricks []ports.Brick, brickIdsToRemove []ports.BrickDependency) []ports.Brick {
	brickIdsToRemoveMap := make(map[string]bool)
	for _, b := range brickIdsToRemove {
		brickIdsToRemoveMap[b.Id] = true
	}

	filteredBricks := []ports.Brick{}
	for _, b := range bricks {
		if brickIdsToRemoveMap[b.Id] == false {
			filteredBricks = append(filteredBricks, b)
		}
	}
	return filteredBricks
}

func (b BrickApi) Add(servicePath string, brickId string, parameterResolver ports.ParameterResolver) error {

	bricks, err := GetBricksRecursive(brickId, b.Db, map[string]bool{})
	if err != nil {
		return err
	}

	service, err := b.ServicePersistence.Load(servicePath)
	if err != nil {
		return err
	}

	bricks = removeBricks(bricks, service.BrickIds) //remove all bricks that are already there

	if len(bricks) == 0 {
		return fmt.Errorf("brick %s has already been added.", brickId)
	}

	parameters, err := ResolveParameterSlice(bricks, parameterResolver)
	if err != nil {
		return err
	}

	for _, brick := range bricks {
		if err := AddSingleBrick(&service, brick, parameters); err != nil {
			return err
		}
	}

	if err := b.ServicePersistence.Save(service); err != nil {
		return err
	}
	return nil
}

func (b BrickApi) Upgrade(brickId string) error {
	brick, err := b.Db.Brick(brickId)
	if err != nil {
		return err
	}

	dependencies, err := b.PackageDependencyReader.ReadFromBrick(brick, func(line string, state string) (bool, string) {
		state = getCurrentSection(line, state)
		return state == "CONAN-DEPENDENCIES", state
	})

	if err != nil {
		return err
	}

	allUptodate := true
	for _, d := range dependencies {
		vus := VersionUpgradeSpec{previous: d.Version, latestWorking: d.Version}
		availableVersionStrings, err := b.DependencyInfo.AvailableVersions(d.Id)
		if err != nil {
			fmt.Printf("%s: unable to find any versions (%v)\n", d.Id, err)
		} else if len(availableVersionStrings) == 0 {
			fmt.Printf("%s: unable to find any versions\n", d.Id)
		} else {
			vus.latestAvailable = availableVersionStrings[len(availableVersionStrings)-1]
			vus.target = vus.latestAvailable
			if vus.UpgradeRequired() {
				fmt.Printf("%s: scheduled for upgrade from %s to %s \n", d.Id, vus.previous, vus.target)
				allUptodate = false
			} else {
				fmt.Printf("%s: current version %s is already now up to date. No upgrade required.\n", d.Id, vus.previous)
			}
		}
	}

	if allUptodate {
		fmt.Println("all dependencies are up to date. Nothing to do.")
		return nil
	}

	parentDir, err := ioutil.TempDir("", "sapper_upgrade_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(parentDir)

	fmt.Printf("creating temp service...")
	service, err := b.ServiceApi.Add(brickId, parentDir, utils.DummyParameterResolver{})
	if err != nil {
		fmt.Printf("failed\n")
		return err
	}
	fmt.Printf("success\n")

	fmt.Printf("building service...")
	buildLogFilename, err := b.ServiceApi.Build(service.Path)
	if err != nil {
		fmt.Printf("failed (see %s for details)\n", buildLogFilename)
		return err
	} else {
		fmt.Printf("success\n")
	}

	return nil
}

func (b BrickApi) List() []ports.Brick {
	return b.Db.Bricks(ports.Extension)
}

func (b BrickApi) Search(term string) []ports.Brick {
	filteredBricks := []ports.Brick{}
	for _, brick := range b.Db.Bricks(ports.Extension) {
		if strings.Contains(brick.Id, term) || strings.Contains(brick.Description, term) {
			filteredBricks = append(filteredBricks, brick)
		}
	}
	return filteredBricks
}

var _ ports.BrickApi = BrickApi{}

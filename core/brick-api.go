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
	Configuration           ports.Configuration
	BrickDBFactory          ports.BrickDBFactory
	ServicePersistence      ports.ServicePersistence
	PackageDependencyReader ports.BrickPackageDependencyReader
	PackageDependencyWriter ports.BrickPackageDependencyWriter
	DependencyInfo          ports.DependencyInfo
	ServiceApi              ServiceApi
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
	db, err := b.BrickDBFactory.MakeAggregatedBrickDB(b.Configuration.Remotes(), b.Configuration.DefaultRemotesDir())
	if err != nil {
		return err
	}

	bricks, err := GetBricksRecursive(brickId, db, map[string]bool{})
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

func isInDependencySection(line string, state string) (bool, string) {
	state = getCurrentSection(line, state)
	return state == "CONAN-DEPENDENCIES", state
}

func (b BrickApi) UpgradeInDB(brickId string, db ports.BrickDB) error {
	brick, err := db.Brick(brickId)
	if err != nil {
		return err
	}
	dependencies, err := b.PackageDependencyReader.ReadFromBrick(brick, isInDependencySection)
	if err != nil {
		return err
	}

	//2. check if update is required
	toBeUpdated := []ports.PackageDependency{}
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
				toBeUpdated = append(toBeUpdated, d)
			} else {
				fmt.Printf("%s: current version %s is already up to date. No upgrade required.\n", d.Id, vus.previous)
			}
		}
	}
	if len(toBeUpdated) == 0 {
		fmt.Println("all dependencies are up to date. Nothing to do.")
		return nil
	}

	//3. create & build service with just that brick
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

	//4. do the upgrade
	upgradeMap := map[string]VersionUpgradeSpec{}
	for _, d := range toBeUpdated {
		fmt.Printf("%s: ", d.Id)
		vus, err := b.ServiceApi.upgradeDependency(service, d, false)
		if err == nil {
			upgradeMap[d.Id] = vus
		} else {
			fmt.Printf("failed (%v)\n", err)
		}
	}

	//5. write out the new package dependencies
	pd := []ports.PackageDependency{}
	for dependency, vus := range upgradeMap {
		if vus.UpgradeRequired() && !vus.UpgradeCompletelyFailed() {
			pd = append(pd, ports.PackageDependency{Id: dependency, Version: vus.latestWorking})
		}
	}
	err = b.PackageDependencyWriter.WriteToBrick(brick, pd, isInDependencySection)
	if err != nil {
		return err
	}

	//6. print status
	errorCount := 0
	for dependency, vus := range upgradeMap {
		pd := ports.PackageDependency{Id: dependency, Version: vus.previous}
		vus.PrintStatus(os.Stdout, pd)
		if !vus.UpgradeToTargetSuccessful() {
			errorCount++
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("%s: failed to upgrade %v of %v dependencies", brickId, errorCount, len(upgradeMap))
	}

	return nil
}

func (b BrickApi) Upgrade(brickId string) error {
	//1. read package dependencies from brick
	db, err := b.BrickDBFactory.MakeAggregatedBrickDB(b.Configuration.Remotes(), b.Configuration.DefaultRemotesDir())
	if err != nil {
		return err
	}

	return b.UpgradeInDB(brickId, db)
}

func (b BrickApi) List() []ports.Brick {
	db, err := b.BrickDBFactory.MakeAggregatedBrickDB(b.Configuration.Remotes(), b.Configuration.DefaultRemotesDir())
	if err != nil {
		return nil
	}

	return db.Bricks(ports.Extension)
}

func (b BrickApi) Search(term string) []ports.Brick {
	db, err := b.BrickDBFactory.MakeAggregatedBrickDB(b.Configuration.Remotes(), b.Configuration.DefaultRemotesDir())
	if err != nil {
		return nil
	}

	filteredBricks := []ports.Brick{}
	for _, brick := range db.Bricks(ports.Extension) {
		if strings.Contains(brick.Id, term) || strings.Contains(brick.Description, term) {
			filteredBricks = append(filteredBricks, brick)
		}
	}
	return filteredBricks
}

var _ ports.BrickApi = BrickApi{}
var _ ports.BrickUpgrader = BrickApi{}

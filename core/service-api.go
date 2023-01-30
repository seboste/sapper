package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/seboste/sapper/ports"
	"github.com/seboste/sapper/utils"
)

type ServiceApi struct {
	Db                 ports.BrickDB
	ServicePersistence ports.ServicePersistence
	ServiceBuilder     ports.ServiceBuilder
	ParameterResolver  ports.ParameterResolver
	DependencyInfo     ports.DependencyInfo
	DependencyWriter   ports.DependencyWriter
}

func ResolveParameters(bp []ports.BrickParameters, pr ports.ParameterResolver) (map[string]string, error) {
	parameters := make(map[string]string)
	for _, p := range bp {
		value := pr.Resolve(p.Name, p.Default)
		if value == "" {
			return nil, fmt.Errorf("unable to resolve value for parameter %s", p.Name)
		}
		parameters[p.Name] = value
	}
	return parameters, nil
}

func ResolveParameterSlice(bricks []ports.Brick, pr ports.ParameterResolver) (map[string]string, error) {
	combinedParameters := map[string]string{}
	for _, brick := range bricks {
		p, err := ResolveParameters(brick.Parameters, pr)
		if err != nil {
			return nil, err
		}
		for k, v := range p {
			combinedParameters[k] = v
		}
	}
	return combinedParameters, nil
}

func AddSingleBrick(s *ports.Service, b ports.Brick, parameters map[string]string) error {
	for _, f := range b.Files {
		inputFilePath := filepath.Join(b.BasePath, f)
		if _, err := os.Stat(inputFilePath); err != nil {
			return err
		}

		outputFilePath := filepath.Join(s.Path, f)
		outputDir, _ := filepath.Split(outputFilePath)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}

		content, err := ioutil.ReadFile(inputFilePath)
		if err != nil {
			return err
		}
		contentStr := string(content)

		contentStr = replaceParameters(contentStr, parameters)

		if _, err := os.Stat(outputFilePath); errors.Is(err, os.ErrNotExist) { //file does not exit => just write out the content
			err = ioutil.WriteFile(outputFilePath, []byte(contentStr), 0644)
			if err != nil {
				return err
			}
		} else { //file exists => merge Sections
			inputSectionSlice, err := readSections(contentStr)
			if err != nil {
				return err
			}
			inputSections := toMap(inputSectionSlice)

			outputContent, err := ioutil.ReadFile(outputFilePath)
			if err != nil {
				return err
			}

			outputContentStr := string(outputContent)

			mergedOutputContentStr, err := mergeSections(outputContentStr, inputSections)

			err = ioutil.WriteFile(outputFilePath, []byte(mergedOutputContentStr), 0644)
			if err != nil {
				return err
			}
		}
	}

	s.BrickIds = append(s.BrickIds, ports.BrickDependency{b.Id, b.Version})

	return nil
}

func mergeSection(base section, incoming section) (string, error) {

	if base.name != incoming.name {
		return "", fmt.Errorf("Unable to merge section %s with section %s. Names must match.", base.name, incoming.name)
	}

	if base.verb != "" {
		return "", fmt.Errorf("Unable to merge section %s. Base operation must not be defined.", base.name)
	}

	if incoming.verb == "REPLACE" {
		return incoming.content, nil
	} else if incoming.verb == "PREPEND" {
		var content string
		content = content + incoming.content
		if base.content != "" && incoming.content != "" {
			content = content + fmt.Sprintln("")
		}
		content = content + base.content
		return content, nil
	} else if incoming.verb == "APPEND" {
		var content string
		content = content + base.content
		if base.content != "" && incoming.content != "" {
			content = content + fmt.Sprintln("")
		}
		content = content + incoming.content
		return content, nil
	} else {
		return "", fmt.Errorf("Unable to merge section %s. Invalid incoming operation %s.", base.name, incoming.verb)
	}
}

func mergeSections(content string, inputSections map[string]section) (string, error) {
	outputSections, err := readSections(content)
	if err != nil {
		return content, err
	}

	scanner := bufio.NewScanner(strings.NewReader(content))

	outputContent := ""
	lineNumber := 0
	for _, s := range outputSections {
		//advance to next section
		for lineNumber < s.lineBegin && scanner.Scan() {
			outputContent = outputContent + fmt.Sprintln(scanner.Text())
			lineNumber = lineNumber + 1
		}

		if incomingSection, ok := inputSections[s.name]; ok {
			mergedSectionContent, err := mergeSection(s, incomingSection)
			if err != nil {
				return "", err
			}
			if mergedSectionContent != "" {
				outputContent = outputContent + fmt.Sprintln(mergedSectionContent)
			}
		} else { // no incoming section => just use base content
			outputContent = outputContent + s.content
		}
		for lineNumber < s.lineEnd && scanner.Scan() { //skip section content
			lineNumber = lineNumber + 1
		}

		//take care of end tag
		scanner.Scan()
		outputContent = outputContent + fmt.Sprintln(scanner.Text())
		lineNumber = lineNumber + 1
	}

	//copy incoming content after last section
	for scanner.Scan() {
		outputContent = outputContent + fmt.Sprintln(scanner.Text())
		lineNumber = lineNumber + 1
	}

	return outputContent, nil
}

func replaceParameters(content string, parameters map[string]string) string {
	for name, value := range parameters {
		pattern := "<<<" + name + ">>>"
		content = strings.ReplaceAll(content, pattern, value)
	}
	return content
}

func GetBricksRecursive(brickId string, db ports.BrickDB, parentBrickIds map[string]bool) ([]ports.Brick, error) {

	brickIds := make(map[string]bool)
	for k, v := range parentBrickIds {
		brickIds[k] = v
	}

	bricks := []ports.Brick{}

	if brickIds[brickId] == true {
		return nil, fmt.Errorf("cyclic brick dependency")
	}

	brick, err := db.Brick(brickId)
	if err != nil {
		return bricks, fmt.Errorf("invalid brick %s", brickId)
	}

	bricks = append(bricks, brick)
	brickIds[brick.Id] = true

	//deep copy to identify cyclic dependencies
	baselineBrickIds := make(map[string]bool)
	for k, v := range brickIds {
		baselineBrickIds[k] = v
	}

	for _, dependencyId := range brick.Dependencies {
		dependencies, err := GetBricksRecursive(dependencyId, db, baselineBrickIds)
		if err != nil {
			return nil, err
		}

		for _, dependency := range dependencies {
			if brickIds[dependency.Id] == false {
				bricks = append(bricks, dependency)
				brickIds[dependency.Id] = true
			}
		}
	}

	return bricks, nil
}

func (s ServiceApi) Add(templateName string, parentDir string, parameterResolver ports.ParameterResolver) error {

	bricks, err := GetBricksRecursive(templateName, s.Db, map[string]bool{})
	if err != nil {
		return err
	}

	parameters, err := ResolveParameterSlice(bricks, parameterResolver)
	if err != nil {
		return err
	}

	if len(bricks) == 0 {
		return fmt.Errorf("invalid template %s", templateName)
	}

	serviceName := parameters["NAME"]
	if serviceName == "" {
		return fmt.Errorf("invalid service name %s", serviceName)
	}
	outputBasePath := filepath.Join(parentDir, serviceName)
	if err := os.MkdirAll(outputBasePath, os.ModePerm); err != nil {
		return err
	}

	service := ports.Service{Id: serviceName, Path: outputBasePath}

	for _, brick := range bricks {
		if err := AddSingleBrick(&service, brick, parameters); err != nil {
			return err
		}
	}

	if err := s.ServicePersistence.Save(service); err != nil {
		return err
	}
	return nil
}

func (s ServiceApi) upgradeDependencyToVersion(service ports.Service, d ports.PackageDependency, targetVersion string) (string, error) {
	err := s.DependencyWriter.Write(service, d.Id, targetVersion)
	if err != nil {
		return "", err
	}
	buildLogFileName, err := s.Build(service.Path)
	if err != nil {
		return buildLogFileName, err
	}

	return buildLogFileName, nil
}

// sortedVersions must range from current version to latest version to be considered (must at least have one entry)
// isWorkinbg is a predicate that checks if a specific version is working
// returns latest working version
func findLatestWorkingVersion(sortedVersions []SemanticVersion, isWorking func(v SemanticVersion) bool) SemanticVersion {

	latestWorkingVersion := sortedVersions[0]
	sortedVersions = sortedVersions[1:]

	i := len(sortedVersions) - 1
	for len(sortedVersions) >= 1 {
		if isWorking(sortedVersions[i]) {
			//this works => exclude all that are lower than current version
			latestWorkingVersion = sortedVersions[i]
			sortedVersions = sortedVersions[i+1:]
		} else {
			//this version does not work => exclude all that are higher or equal to the current version
			sortedVersions = sortedVersions[:i]
		}
		i = len(sortedVersions) / 2
	}

	return latestWorkingVersion
}

type VersionUpgradeSpec struct {
	previous, target, latestAvailable, latestWorking string
}

func filterSemvers(in []SemanticVersion, predicate func(SemanticVersion) bool) []SemanticVersion {
	out := []SemanticVersion{}
	for _, v := range in {
		if predicate(v) {
			out = append(out, v)
		}
	}
	return out
}

func (s ServiceApi) upgradeDependency(service ports.Service, d ports.PackageDependency, keepMajorVersion bool) (VersionUpgradeSpec, error) {
	vus := VersionUpgradeSpec{previous: d.Version, latestWorking: d.Version} //assume that the current version is working

	//1. determine all available versions
	availableVersionStrings, err := s.DependencyInfo.AvailableVersions(d.Id)
	if err != nil {
		return vus, err
	}
	if len(availableVersionStrings) == 0 {
		return vus, fmt.Errorf("unable to find any versions of %s", d.Id)
	}
	vus.latestAvailable = availableVersionStrings[len(availableVersionStrings)-1]

	//2. check if we can use semantic versions
	semvers := []SemanticVersion{}
	currentSemver, err := ParseSemanticVersion(vus.previous)
	if err == nil {
		semvers, err = ConvertToSemVer(availableVersionStrings)
	}

	if err == nil {
		//yes => use semantic versions

		//a) sort versions
		sort.Sort(ByVersion(semvers))
		vus.latestAvailable = semvers[len(semvers)-1].String()

		//a) exclude all old versions
		semvers = filterSemvers(semvers, func(v SemanticVersion) bool { return !Less(v, currentSemver) })
		//b) if wanted, exclude all versions with a different major version
		if keepMajorVersion {
			semvers = filterSemvers(semvers, func(v SemanticVersion) bool { return currentSemver.Major == v.Major })
		}

		if len(semvers) == 0 { //this should not happen because the current version should always be included
			vus.target = vus.previous
			return vus, fmt.Errorf("unable to find any versions of %s that can be considered", d.Id)
		}

		//early exit?
		vus.target = semvers[len(semvers)-1].String()
		if vus.previous == vus.target {
			return vus, nil
		}

		vus.latestWorking = findLatestWorkingVersion(semvers, func(v SemanticVersion) bool {
			fmt.Printf("trying to upgrade to %v...", v)
			buildLogFilename, err := s.upgradeDependencyToVersion(service, d, v.String())
			if err == nil {
				fmt.Printf("success\n")
				return true
			} else {
				fmt.Printf("failed (see %s for details)\n", buildLogFilename)
				return false
			}
		}).String()
	} else {
		//no => no semantic versioning, just upgrade to the latest
		vus.target = vus.latestAvailable
		if vus.previous == vus.target {
			return vus, nil
		}

		fmt.Printf("%s => simply trying to upgrade to the latest version %s...", err.Error(), vus.target)
		buildLogFilename, err := s.upgradeDependencyToVersion(service, d, vus.target)
		if err == nil {
			vus.latestWorking = vus.target
			fmt.Printf("success\n")
		} else {
			fmt.Printf("failed (see %s for details)", buildLogFilename)
		}

	}

	//4. set the latest working version
	err = s.DependencyWriter.Write(service, d.Id, vus.latestWorking)
	if err != nil {
		return vus, err
	}

	return vus, nil
}

func (s ServiceApi) Upgrade(path string, keepMajorVersion bool) error {
	service, err := s.ServicePersistence.Load(path)
	if err != nil {
		return err
	}

	fmt.Printf("building service...")
	buildLogFilename, err := s.Build(path)
	if err != nil {
		fmt.Printf("failed (see %s for details)\n", buildLogFilename)
		return err
	} else {
		fmt.Println("success")
	}

	hasError := false
	for _, d := range service.Dependencies {
		fmt.Printf("upgrading %s (current version %s)\n", d.Id, d.Version)
		vus, err := s.upgradeDependency(service, d, keepMajorVersion)
		if err != nil {
			fmt.Println(err.Error())
			hasError = true
		} else {
			if vus.previous == vus.latestAvailable {
				fmt.Printf("%s is already now up to date. No upgrade required.\n", d.Id)
			} else if vus.latestWorking == vus.latestAvailable {
				fmt.Printf("upgrade from %s to %s succeeded. %s is now up to date.\n", vus.previous, vus.latestAvailable, d.Id)
			} else if vus.latestWorking == vus.target {
				fmt.Printf("upgrade from %s to %s succeeded. However, there is a newer version %s available.\n", vus.previous, vus.target, vus.latestAvailable)
			} else if vus.latestWorking != vus.previous {
				if vus.target == vus.latestAvailable {
					fmt.Printf("upgrade from %s to %s failed => upgrade to latest working version %s instead\n", vus.previous, vus.target, vus.latestWorking)
				} else {
					fmt.Printf("upgrade from %s to %s failed => upgrade to latest working version %s instead. Note that there is an even newer version %s available.\n", vus.previous, vus.target, vus.latestWorking, vus.latestAvailable)
				}
			} else if vus.latestWorking == vus.previous {
				fmt.Printf("upgrade from %s to %s failed => keeping version %s\n", vus.previous, vus.target, vus.previous)
			}
		}
	}

	if hasError == true {
		return fmt.Errorf("Unable to upgrade all dependencies")
	}

	return nil
}

func (s ServiceApi) Build(path string) (string, error) {
	service, err := s.ServicePersistence.Load(path)
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile("", "sapper_build_log_*.log")
	if err != nil {
		return "", err
	}

	slw := utils.MakeSingleLineWriter(os.Stdout)
	defer slw.Cleanup()

	err = s.ServiceBuilder.Build(service, io.MultiWriter(slw, f))
	if err != nil {
		return f.Name(), err
	} else {
		os.Remove(f.Name())
		return "", nil
	}

}

func (s ServiceApi) Test() {
	fmt.Println("test")
}

func (s ServiceApi) Describe(path string, writer io.Writer) error {

	service, err := s.ServicePersistence.Load(path)
	if err != nil {
		return err
	}

	writer.Write([]byte(fmt.Sprintln("Id:", service.Id)))
	writer.Write([]byte(fmt.Sprintln("Path:", service.Path)))
	writer.Write([]byte(fmt.Sprintln("BrickIds:")))
	for _, brickId := range service.BrickIds {
		writer.Write([]byte(fmt.Sprintln("  - Id:", brickId.Id)))
		writer.Write([]byte(fmt.Sprintln("    Version:", brickId.Version)))
	}
	writer.Write([]byte(fmt.Sprintln("Dependencies:")))
	for _, dependency := range service.Dependencies {
		writer.Write([]byte(fmt.Sprintln("  - Id:", dependency.Id)))
		availableVersions, err := s.DependencyInfo.AvailableVersions(dependency.Id)
		if err != nil || len(availableVersions) == 0 || availableVersions[len(availableVersions)-1] == dependency.Version {
			writer.Write([]byte(fmt.Sprintln("    Version:", dependency.Version)))
		} else {
			writer.Write([]byte(fmt.Sprintln("    Version:", dependency.Version, ", newer version", availableVersions[len(availableVersions)-1], "available")))
		}
	}

	return nil
}

func (s ServiceApi) Deploy() {
	fmt.Println("deploy")
}

var _ ports.ServiceApi = ServiceApi{}

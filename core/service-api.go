package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/seboste/sapper/ports"
	"github.com/seboste/sapper/utils"
)

type ServiceApi struct {
	Db                 ports.BrickDB
	ServicePersistence ports.ServicePersistence
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

func (s ServiceApi) upgradeDependencyToVersion(service ports.Service, d ports.PackageDependency, targetVersion string) error {
	err := s.DependencyWriter.Write(service, d.Id, targetVersion)
	if err != nil {
		return err
	}
	err = s.Build(service.Path)
	if err != nil {
		return err
	}

	return nil
}

// returns highestVersion, highestWorkingVersion
func findLatestWorkingVersion(versions []SemanticVersion, isWorking func(v SemanticVersion) bool) (*SemanticVersion, *SemanticVersion) {
	if len(versions) < 1 {
		return nil, nil
	}

	sort.Sort(ByVersion(versions))

	i := len(versions) - 1 // start with highest version
	highestVersion := versions[i]
	var highestWorkingVersion *SemanticVersion = nil

	for len(versions) > 1 || (len(versions) == 1 && highestWorkingVersion == nil) {

		if isWorking(versions[i]) {
			//this works => exclude all that are lower than current version
			currentSemver := versions[i]
			highestWorkingVersion = &currentSemver
			versions = versions[i:]
		} else {

			//this version does not work => exclude all that are higher or equal to the current version
			versions = versions[:i]
		}
		i = len(versions) / 2
	}

	return &highestVersion, highestWorkingVersion
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

func (s ServiceApi) upgradeDependency(service ports.Service, d ports.PackageDependency, keepMajorVersion bool) (string, error) {
	availableVersionStrings, err := s.DependencyInfo.AvailableVersions(d.Id)
	if err != nil {
		return "", err
	}
	if len(availableVersionStrings) == 0 {
		return "", fmt.Errorf("unable to find any versions of %s", d.Id)
	}

	latest := availableVersionStrings[len(availableVersionStrings)-1]
	if d.Version == latest {
		fmt.Printf("%s is already up to date (%s)\n", d.Id, d.Version)
		return d.Version, nil
	}

	semvers := []SemanticVersion{}
	currentSemver, err := ParseSemanticVersion(d.Version)
	if err == nil {
		semvers, err = ConvertToSemVer(availableVersionStrings)
	}
	newVersion := ""
	if err == nil {
		//use semantic versions

		semvers = filterSemvers(semvers, func(v SemanticVersion) bool { return Less(currentSemver, v) })
		if keepMajorVersion {
			semvers = filterSemvers(semvers, func(v SemanticVersion) bool { return currentSemver.Major == v.Major })
		}

		if len(semvers) == 0 {
			return d.Version, fmt.Errorf("unable to find any versions of %s higher than current version %v", d.Id, currentSemver)
		}

		highestSemver, highestSuccessSemver := findLatestWorkingVersion(semvers, func(v SemanticVersion) bool {
			fmt.Printf("Trying to upgrade to %v\n", v)
			err := s.upgradeDependencyToVersion(service, d, v.String())
			return err == nil
		})

		if highestSuccessSemver != nil && highestSemver != nil {
			newVersion = highestSuccessSemver.String()
			err = s.DependencyWriter.Write(service, d.Id, highestSuccessSemver.String())
			if err != nil {
				return "", err
			}
			if *highestSemver == *highestSuccessSemver {
				fmt.Printf("Upgrade from %v to %v succeeded. %s is now up to date.\n", currentSemver, highestSemver, d.Id)
			} else {
				fmt.Printf("Upgrade from %v to %v failed => upgrade to highest working version %v instead\n", currentSemver, highestSemver, *highestSuccessSemver)
			}
		} else {
			newVersion = currentSemver.String()
			err = s.DependencyWriter.Write(service, d.Id, currentSemver.String())
			if err != nil {
				return "", err
			}
			fmt.Printf("Upgrade from %v to %v failed => rolling back to previous version %v instead\n", currentSemver, highestSemver, currentSemver)
		}
	} else {
		fmt.Printf("%s => simply trying to upgrade to the latest version %s\n", err.Error(), latest)
		err := s.upgradeDependencyToVersion(service, d, latest)
		if err != nil {
			fmt.Printf("Upgrade from %s to %s failed to build => rollback to %s\n", latest, d.Id, d.Version)
			s.DependencyWriter.Write(service, d.Id, d.Version) //roll back
			return d.Version, err
		}
		fmt.Printf("Upgrade from %s to %s succeeded\n", d.Version, latest)
		newVersion = latest
	}

	return newVersion, nil
}

func (s ServiceApi) Upgrade(path string, keepMajorVersion bool) error {
	service, err := s.ServicePersistence.Load(path)
	if err != nil {
		return err
	}

	fmt.Println("building service...")
	err = s.Build(path)
	if err != nil {
		return err
	}

	hasError := false
	for _, d := range service.Dependencies {
		fmt.Printf("upgrading %s (current version %s)\n", d.Id, d.Version)
		_, err := s.upgradeDependency(service, d, keepMajorVersion)
		if err != nil {
			fmt.Println(err.Error())
			hasError = true
		} else {
			//fmt.Printf("version of %s has been set from %s to %s\n", d.Id, d.Version, version)
		}
	}

	if hasError == true {
		return fmt.Errorf("Unable to upgrade all dependencies")
	}

	return nil
}

func (s ServiceApi) Build(path string) error {
	cmd := exec.Command("make", "build", "-B")
	cmd.Dir = path

	slw := utils.MakeSingleLineWriter(os.Stdout)
	cmd.Stdout = slw
	cmd.Stderr = slw
	err := cmd.Run()
	slw.Cleanup()
	return err
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

package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/seboste/sapper/ports"
)

type Service struct {
}

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func ReplaceInFile(path string, old string, new string) error {

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)

	if strings.Contains(content, old) {
		content = strings.ReplaceAll(content, old, new)
		data = []byte(content)
		err = ioutil.WriteFile(path, data, info.Mode())
	}
	return nil
}

func ReplaceInPath(path string, old string, new string) error {
	filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				ReplaceInFile(path, old, new)
			}

			return nil
		})

	return nil
}

func (s Service) Add(name string, version string, templateName string, parentDir string) error {

	inputDir := filepath.Join("./remote", templateName)
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return err
	}

	outputDir := filepath.Join(parentDir, name)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	if err := Dir(inputDir, outputDir); err != nil {
		return err
	}

	if err := ReplaceInPath(outputDir, "<<<VERSION>>>", version); err != nil {
		return err
	}

	if err := ReplaceInPath(outputDir, "<<<NAME>>>", name); err != nil {
		return err
	}

	return nil
}

func (s Service) Update() {
	fmt.Println("update")
}

func (s Service) Build(path string) error {
	cmd := exec.Command("make", "build", "-B")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (s Service) Test() {
	fmt.Println("test")
}

func (s Service) Deploy() {
	fmt.Println("deploy")
}

var _ ports.ServiceApi = Service{}

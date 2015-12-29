package yum

import (
	".."
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateRepo(t *testing.T) {
	dirname := "../rpms"

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		t.Fatalf("Error reading RPM files: %v", err)
	}

	// scan for rpm files to create work queue
	rpmFiles := make([]string, 0)
	for _, file := range files {
		path := filepath.Join(dirname, file.Name())
		if strings.HasSuffix(strings.ToLower(path), ".rpm") {
			rpmFiles = append(rpmFiles, path)
		}
	}

	// start producer routine
	fileChannel := make(chan string, 100)
	go func() {
		for _, rpmFile := range rpmFiles {
			fileChannel <- rpmFile
		}

		close(fileChannel)
	}()

	// start workers
	workerCount := 4
	packages := make(chan *rpm.Package, workerCount)
	for i := 0; i < workerCount; i++ {
		go OpenPackages(fileChannel, packages)
	}

	// fan in results
	rpms := make([]rpm.Package, len(rpmFiles))
	for i := 0; i < len(rpmFiles); i++ {
		p := <-packages
		if p != nil {
			rpms[i] = *p
			t.Logf("Added package: %v", p)
		}
	}
}

func OpenPackages(paths <-chan string, packages chan<- *rpm.Package) {
	//packages := make(chan *rpm.Package)
	//go func() {
	for path := range paths {
		p, err := rpm.OpenPackage(path)
		if err != nil {
			// todo
		} else {
			packages <- p
		}
	}
	//close(packages)
	//}()

	//return packages
}

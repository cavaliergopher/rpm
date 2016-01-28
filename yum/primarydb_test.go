package yum

import (
	"github.com/cavaliercoder/go-rpm"
	//"os"
	"testing"
)

func TestAddPackages(t *testing.T) {
	dbpath := "./primary_db.sqlite"
	//defer os.Remove(dbpath)

	// read packages from ../testdata
	packageFiles, err := rpm.OpenPackageFiles("../testdata")
	if err != nil {
		t.Fatalf("Error loading test packages: %v", err)
	}

	// create db
	pdb, err := CreatePrimaryDB(dbpath)
	if err != nil {
		t.Fatalf("Error creating primary_db: %v", err)
	}

	// convert PackageFiles to rpm.Package interfaces
	packages := make([]rpm.Package, len(packageFiles))
	for i := 0; i < len(packageFiles); i++ {
		packages[i] = rpm.Package(packageFiles[i])
	}

	// add package to database
	err = pdb.AddPackages(packages)
	if err != nil {
		t.Errorf("Error adding package to database: %v", err)
	}
}

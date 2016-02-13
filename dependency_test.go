package rpm

import (
	"testing"
)

func TestDependencies(t *testing.T) {
	// load package file paths
	files, err := packages(t)
	if err != nil {
		t.Fatalf("Error listing rpm packages: %v", err)
	}

	// load each package
	for _, path := range files {
		p, err := OpenPackageFile(path)
		if err != nil {
			t.Errorf("%v", err)
		}

		// all should have Requires
		if reqs := p.Requires(); len(reqs) == 0 {
			t.Errorf("No Require dependencies found for package %v", p)
		}

		// all should have Provides
		if provs := p.Provides(); len(provs) == 0 {
			t.Errorf("No Provides dependencies found for package %v", p)
		}

		// some will have Conflicts
		p.Conflicts()

		// some will have Obsoletes
		p.Obsoletes()
	}
}

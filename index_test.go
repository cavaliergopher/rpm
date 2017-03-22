package rpm

import (
	"os"
	"testing"
)

// global to defeat compiler optimization
var X interface{}

func BenchmarkIndexReads(b *testing.B) {
	path := os.Getenv("RPM_DIR")
	if path == "" {
		path = "testdata"
	}

	pkgs, err := OpenPackageFiles(path)
	if err != nil {
		panic(err)
	}

	// open and read the package list b.N times
	var V interface{} // defeat compiler op
	for n := 0; n < b.N; n++ {
		for _, p := range pkgs {
			V = p.String()
			V = p.Requires()
			V = p.Provides()
			V = p.Conflicts()
			V = p.Obsoletes()
			V = p.Files()
			V = p.Summary()
			V = p.Description()
		}
	}

	X = V
}

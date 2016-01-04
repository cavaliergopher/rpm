package yum_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm/yum"
	"net/http"
)

func ExampleReadRepoMetadata() {
	// base url for a public yum repository
	baseurl := "http://mirror.centos.org/centos/7/os/x86_64/"

	// repo metadata is always found at the repodata/repomd.xml subpath
	repomdurl := baseurl + "repodata/repomd.xml"

	// get repo metadata from url
	resp, err := http.Get(repomdurl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// decode http stream into repo metadata struct
	repomd, err := yum.ReadRepoMetadata(resp.Body)
	if err != nil {
		panic(err)
	}

	// profit
	fmt.Printf("Downloaded repository metadata revision %d\n", repomd.Revision)
}

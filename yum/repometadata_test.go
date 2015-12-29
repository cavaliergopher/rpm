package yum

import (
	"net/http"
	"testing"
)

func TestReadRepoMetadata(t *testing.T) {
	// open repo metadata from URL
	resp, err := http.Get("http://mirror.centos.org/centos/7/updates/x86_64/repodata/repomd.xml")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	// decode repo metadata into struct
	_, err = ReadRepoMetadata(resp.Body)
	if err != nil {
		t.Errorf("%v", err)
	}
}

package yum_test

import (
	"fmt"
	"github.com/cavaliercoder/go-rpm/yum"
)

func ExampleValidateFileChecksum() {
	file := "primary_db.sqlite.bz2"
	checksum := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	checksumType := "sha256"

	if err := yum.ValidateFileChecksum(file, checksum, checksumType); err == yum.ErrChecksumMismatch {
		fmt.Printf("File failed checksum validation.\n")
	} else if err == nil {
		fmt.Printf("File passed checksum validation.\n")
	} else {
		panic(err)
	}
}

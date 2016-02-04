// +build gofuzz

package rpm

import "bytes"

func Fuzz(data []byte) int {
	if _, err := ReadPackageFile(bytes.NewReader(data)); err != nil {
		// err expected on random input
		return 0
	}

	return 1
}

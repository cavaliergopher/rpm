package rpm

import (
	"testing"
)

func TestPrimaryDB(t *testing.T) {
	err := CreatePrimaryDB("./rpms")
	if err != nil {
		t.Errorf("Failed to create primary DB: %v", err)
	}
}

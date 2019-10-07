package util

import "testing"

func TestUUID(t *testing.T) {
	//6b1d3cc6-03bc-436e-bbd9-a292c3741822
	if !IsUUID(`6b1d3cc6-03bc-436e-bbd9-a292c3741822`) {
		t.Errorf("6b1d3cc6-03bc-436e-bbd9-a292c3741822 is not a uuid")
	}
}

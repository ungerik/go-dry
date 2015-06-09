package dry

import "testing"

func TestGetenvDefault(t *testing.T) {
	if GetenvDefault("GO_DRY_BOGUS_ENVIRONMENT_VARIABLE", "default") != "default" {
		t.Fail()
	}
}

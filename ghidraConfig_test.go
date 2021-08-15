package ghidraScriptRunner

import "testing"

func TestValidConfig(t *testing.T) {
	if _, got := NewConfiguration("headless", "projectLocation", "project"); got != nil {
		t.FailNow()
	}

	var nilString string
	if _, got := NewConfiguration(nilString, "", ""); got == nil {
		t.FailNow()
	}
}

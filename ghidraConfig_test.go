package ghidraScriptRunner

import "testing"

func testValidConfig(t *testing.T) {
	if _, err := NewConfiguration("headless", "projectLocation", "project", "script"); err != nil {
		t.FailNow()
	}

	var nilString string
	if _, err := NewConfiguration(nilString, "", "", ""); err == nil {
		t.FailNow()
	}
}

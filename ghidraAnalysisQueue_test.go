package ghidraScriptRunner

import (
	"testing"
)

func TestAddToQueue(t *testing.T) {
	config, _ := NewConfiguration("test", "test", "test", "test")
	ghidraScriptService := NewGhidraScriptService(config)
	binName := "testBinName"
	script := "testScript"

	ghidraScriptService.AddToQueue(&binName, &script)

	if linkedListElement := ghidraScriptService.findElement(&binName); linkedListElement == nil {
		t.FailNow()
	}

}

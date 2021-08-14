package ghidraScriptRunner

import (
	"testing"
)

var config, _ = NewConfiguration("test", "test", "test", "test")

func TestAddToQueue(t *testing.T) {
	ghidraScriptService := NewGhidraScriptService(config)
	binName := "testBinName"
	script := "testScript"

	ghidraScriptService.AddToQueue(&binName, &script)

	if linkedListElement := ghidraScriptService.findElement(&binName); linkedListElement == nil {
		t.FailNow()
	}

	invalidBinName := "test"
	if linkedListElement := ghidraScriptService.findElement(&invalidBinName); linkedListElement != nil {
		t.FailNow()
	}

}

func TestRemoveFromQueue(t *testing.T) {
	ghidraScriptService := NewGhidraScriptService(config)
	binName := "testBinName"
	script := "testScript"

	ghidraScriptService.AddToQueue(&binName, &script)

	if linkedListElement := ghidraScriptService.findElement(&binName); linkedListElement == nil {
		t.FailNow()
	}

	ghidraScriptService.RemoveFromQueue(&binName)
	if linkedListElement := ghidraScriptService.findElement(&binName); linkedListElement != nil {
		t.FailNow()
	}

}

package ghidraScriptRunner

import (
	"container/list"
	"sync"
	"testing"
	"time"
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

func TestUpdateTaskStatus(t *testing.T) {
	ghidraScriptService := NewGhidraScriptService(config)
	binName := "testBinName"
	script := "testScript"

	ghidraScriptService.AddToQueue(&binName, &script)
	ghidraScriptService.UpdateTaskStatus(&binName, runningStatus)
	if taskStatus := ghidraScriptService.GetStatus(&binName); *taskStatus != runningStatus {
		t.FailNow()
	}
}

func TestGetAllStatus(t *testing.T) {
	ghidraScriptService := NewGhidraScriptService(config)
	binName := "testBinName"
	script := "testScript"
	ghidraScriptService.AddToQueue(&binName, &script)

	binName2 := "testBinName2"
	script2 := "testScript2"
	ghidraScriptService.AddToQueue(&binName2, &script2)

	ghidraScriptService.UpdateTaskStatus(&binName, runningStatus)
	if tasksCount := len(ghidraScriptService.GetAllStatus()); tasksCount != 2 {
		t.FailNow()
	}
}

func TestWaitForQueuedItems(t *testing.T) {
	mutex := sync.Mutex{}
	ghidraScriptService := &GhidraScriptService{queue: list.New(), syncCondi: sync.NewCond(&mutex), ghidraConfig: config}
	binName := "testBinName"
	script := "testScript"
	ghidraScriptService.AddToQueue(&binName, &script)
	go ghidraScriptService.waitForQueuedItems()
	time.Sleep(time.Duration(1) * time.Second)
	ghidraScriptService.AddToQueue(&binName, &script)
	time.Sleep(time.Duration(1) * time.Second)
}

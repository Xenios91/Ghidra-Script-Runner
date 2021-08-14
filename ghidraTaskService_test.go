package ghidraScriptRunner

import (
	"container/list"
	"fmt"
	"sync"
	"testing"
	"time"
)

var config, _ = NewConfiguration("test", "test", "test", "test")

func TestAddToQueue(t *testing.T) {
	ghidraScriptService := NewGhidraTaskService(config)
	taskID := "testID"
	script := "testScript"

	ghidraScriptService.AddToQueue(&taskID, &script)

	if linkedListElement := ghidraScriptService.findElementByTaskID(&taskID); linkedListElement == nil {
		t.FailNow()
	}

	invalidBinName := "test"
	if linkedListElement := ghidraScriptService.findElementByTaskID(&invalidBinName); linkedListElement != nil {
		t.FailNow()
	}

}

func TestAddNewTaskToQueue(t *testing.T) {
	mutex := sync.Mutex{}
	ghidraScriptService := &GhidraTaskService{statusMap: make(map[string]*GhidraTaskStatus), queue: list.New(), syncCondi: sync.NewCond(&mutex)}
	taskID := "testID"
	script := "testScript"
	newGhidraTask := NewGhidraScriptTask(&taskID, &script)
	ghidraScriptService.AddNewTaskToQueue(newGhidraTask)
	linkedListElement := ghidraScriptService.findElementByTaskID(&taskID)
	if linkedListElement == nil {
		t.FailNow()
	}
}

func TestRemoveFromQueue(t *testing.T) {
	ghidraScriptService := NewGhidraTaskService(config)
	taskID := "testID"
	script := "testScript"

	ghidraScriptService.AddToQueue(&taskID, &script)

	if linkedListElement := ghidraScriptService.findElementByTaskID(&taskID); linkedListElement == nil {
		t.FailNow()
	}

	ghidraScriptService.RemoveFromQueueByTaskID(&taskID)
	if linkedListElement := ghidraScriptService.findElementByTaskID(&taskID); linkedListElement != nil {
		t.FailNow()
	}

}

func TestUpdateTaskStatus(t *testing.T) {
	ghidraScriptService := NewGhidraTaskService(config)
	taskID := "testID"
	script := "testScript"

	ghidraScriptService.AddToQueue(&taskID, &script)
	ghidraScriptService.UpdateTaskStatusByTaskID(&taskID, completeStatus)
	if taskStatus := ghidraScriptService.GetStatusByTaskID(&taskID); *taskStatus != completeStatus {
		t.FailNow()
	}
}

func TestGetAllStatus(t *testing.T) {
	mutex := sync.Mutex{}
	ghidraScriptService := &GhidraTaskService{statusMap: make(map[string]*GhidraTaskStatus), queue: list.New(), syncCondi: sync.NewCond(&mutex)}
	taskID := "testID"
	script := "testScript"
	ghidraScriptService.AddToQueue(&taskID, &script)

	taskID2 := "testID2"
	script2 := "testScript2"
	ghidraScriptService.AddToQueue(&taskID2, &script2)

	ghidraScriptService.UpdateTaskStatusByTaskID(&taskID, runningStatus)
	ghidraScriptService.UpdateTaskStatusByTaskID(&taskID2, waitingStatus)
	if tasksCount := len(ghidraScriptService.GetAllStatus()); tasksCount != 2 {
		t.FailNow()
	}

	statusMap := ghidraScriptService.GetAllStatus()
	if *statusMap[taskID] != runningStatus {
		t.FailNow()
	}

	if *statusMap[taskID2] != waitingStatus {
		t.FailNow()
	}

}

func TestWaitForQueuedItems(t *testing.T) {
	mutex := sync.Mutex{}
	ghidraScriptService := &GhidraTaskService{statusMap: make(map[string]*GhidraTaskStatus), queue: list.New(), syncCondi: sync.NewCond(&mutex), ghidraConfig: config}
	taskID := "testID"
	script := "testScript"
	ghidraScriptService.AddToQueue(&taskID, &script)
	go ghidraScriptService.waitForQueuedItems()
	time.Sleep(time.Duration(1) * time.Second)
	ghidraScriptService.AddToQueue(&taskID, &script)
	time.Sleep(time.Duration(1) * time.Second)
}

func TestIsQueueEmpty(t *testing.T) {
	mutex := sync.Mutex{}
	ghidraScriptService := &GhidraTaskService{statusMap: make(map[string]*GhidraTaskStatus), queue: list.New(), syncCondi: sync.NewCond(&mutex), ghidraConfig: config}
	if !ghidraScriptService.IsQueueEmpty() {
		t.FailNow()
	}
}

func TestGetNextTaskAndRemoveFromQueue(t *testing.T) {
	mutex := sync.Mutex{}
	ghidraScriptService := &GhidraTaskService{statusMap: make(map[string]*GhidraTaskStatus), queue: list.New(), syncCondi: sync.NewCond(&mutex), ghidraConfig: config}
	taskID := "testID"
	script := "testScript"
	ghidraScriptService.AddToQueue(&taskID, &script)
	task := ghidraScriptService.getNextTaskAndRemoveFromQueue()
	if *(*task).GetTaskID() != taskID {
		t.FailNow()
	}
}

func ExampleAddToQueue() {
	//create new confi
	config, err := NewConfiguration("ghidraHeadless", "ghidraProjectLocation", "ghidraProject", "ghidraScript")
	if err == nil {
		fmt.Println("Oh no! We have an invalid config!")
	}
	//create new task service
	ghidraScriptService := NewGhidraTaskService(config)

	taskID := "testID"
	script := "testScript"
	//add example to queue
	ghidraScriptService.AddToQueue(&taskID, &script)
}

func ExampleAddNewTaskToQueue() {
	//create new confi
	config, err := NewConfiguration("ghidraHeadless", "ghidraProjectLocation", "ghidraProject", "ghidraScript")
	if err == nil {
		fmt.Println("Oh no! We have an invalid config!")
	}
	//create new task service
	ghidraScriptService := NewGhidraTaskService(config)

	//create a new task, this task can be any struct that implements the public interface 'GhidraTask'
	taskID := "testID"
	script := "testScript"
	newTask := NewGhidraScriptTask(&taskID, &script)

	//add example to queue
	ghidraScriptService.AddNewTaskToQueue(newTask)
}

func ExampleGetStatusByTaskID() {
	//create new confi
	config, err := NewConfiguration("ghidraHeadless", "ghidraProjectLocation", "ghidraProject", "ghidraScript")
	if err == nil {
		fmt.Println("Oh no! We have an invalid config!")
	}
	//create new task service
	ghidraScriptService := NewGhidraTaskService(config)

	//create a new task, this task can be any struct that implements the public interface 'GhidraTask'
	taskID := "testID"
	script := "testScript"
	newTask := NewGhidraScriptTask(&taskID, &script)

	//add example to queue
	ghidraScriptService.AddNewTaskToQueue(newTask)

	//check status
	status := ghidraScriptService.GetStatusByTaskID(&taskID)
	fmt.Println(string(*status))
}

package ghidraScriptRunner

import (
	"container/list"
	"sync"
)

type GhidraScriptService struct {
	queue        *list.List
	syncCondi    *sync.Cond
	ghidraConfig *Configuration
}

func (ghidraScriptService *GhidraScriptService) waitForQueuedItems() {
	linkedListElement := ghidraScriptService.queue.Front()
	syncCondi := ghidraScriptService.syncCondi

	for {
		syncCondi.L.Lock()
		if linkedListElement == nil {
			syncCondi.Wait()
		} else {
			task := linkedListElement.Value.(ghidraQueueTask)
			ghidraScriptService.UpdateTaskStatus(task.fileName, runningStatus)
			success := task.runScript(ghidraScriptService.ghidraConfig)
			if success {
				ghidraScriptService.UpdateTaskStatus(task.fileName, completeStatus)
			} else {
				ghidraScriptService.UpdateTaskStatus(task.fileName, errorStatus)
			}
		}
		linkedListElement = linkedListElement.Next()
	}
}

//NewGhidraScriptService Used to create a new analysis queue. Required to start polling,
func NewGhidraScriptService(config *Configuration) *GhidraScriptService {
	if config == nil {
		panic("Ghidra Config was nil!")
	}
	mutex := sync.Mutex{}
	newGhidraScriptService := &GhidraScriptService{list.New(), sync.NewCond(&mutex), config}
	go newGhidraScriptService.waitForQueuedItems()
	return newGhidraScriptService
}

//AddToQueue Adds a new analysis task to the queue
func (ghidraScriptService *GhidraScriptService) AddToQueue(binaryName, script *string) {
	queueValue := newGhidraQueueItem(binaryName, script)
	ghidraScriptService.queue.PushBack(queueValue)
}

func (ghidraScriptService *GhidraScriptService) findElement(binaryName *string) *list.Element {
	linkedListElement := ghidraScriptService.queue.Front()
	for {
		if linkedListElement != nil {
			task := linkedListElement.Value.(*ghidraQueueTask)
			if task.fileName == binaryName {
				return linkedListElement
			} else {
				linkedListElement = linkedListElement.Next()
			}
		} else {
			break
		}
	}
	return nil
}

//RemoveFromQueue removes a binary name from the queue of binaries being processed by ghidra.
func (ghidraScriptService *GhidraScriptService) RemoveFromQueue(binaryName *string) {
	linkedListElement := ghidraScriptService.findElement(binaryName)
	if linkedListElement != nil {
		ghidraScriptService.queue.Remove(linkedListElement)
	}
}

//UpdateTaskStatus updates the status of a binary currently in the queue.
func (ghidraScriptService *GhidraScriptService) UpdateTaskStatus(binaryName *string, statusUpdate queueStatus) {
	if linkedListElement := ghidraScriptService.findElement(binaryName); linkedListElement != nil {
		task := linkedListElement.Value.(*ghidraQueueTask)
		task.status = statusUpdate
	}
}

//GetStatus Returns the current status of a binary being processed by Ghidra.
func (ghidraScriptService *GhidraScriptService) GetStatus(binaryName *string) *queueStatus {
	if linkedListElement := ghidraScriptService.findElement(binaryName); linkedListElement != nil {
		status := linkedListElement.Value.(*ghidraQueueTask).status
		return &status
	}
	return nil
}

//GetAllStatus Returns a map with the status of all binaries being processed by Ghidra.
func (ghidraScriptService *GhidraScriptService) GetAllStatus() map[string]*queueStatus {
	statusMap := make(map[string]*queueStatus, ghidraScriptService.queue.Len())
	linkedListElement := ghidraScriptService.queue.Front()

	for {
		if linkedListElement != nil {
			task := linkedListElement.Value.(*ghidraQueueTask)
			statusMap[*task.fileName] = &task.status
			linkedListElement = linkedListElement.Next()
		} else {
			break
		}
	}
	return statusMap
}

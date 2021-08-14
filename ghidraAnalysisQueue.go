package ghidraScriptRunner

import (
	"container/list"
	"sync"
)

//GhidraScriptService the service utilized to run ghidra scripts and manage the task queue.
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
			linkedListElement = ghidraScriptService.queue.Front()
		} else {
			task := linkedListElement.Value.(*queueScriptTask)
			ghidraScriptService.UpdateTaskStatus(task.fileName, runningStatus)
			success := task.runTask(ghidraScriptService.ghidraConfig)
			if success {
				ghidraScriptService.UpdateTaskStatus(task.fileName, completeStatus)
			} else {
				ghidraScriptService.UpdateTaskStatus(task.fileName, errorStatus)
			}
		}
		linkedListElement = linkedListElement.Next()
		syncCondi.L.Unlock()
	}
}

//NewGhidraScriptService used to create a new ghidra script service. Required to start polling task
func NewGhidraScriptService(config *Configuration) *GhidraScriptService {
	if config == nil {
		panic("Ghidra Config was nil!")
	}
	mutex := sync.Mutex{}
	newGhidraScriptService := &GhidraScriptService{list.New(), sync.NewCond(&mutex), config}
	go newGhidraScriptService.waitForQueuedItems()
	return newGhidraScriptService
}

//AddToQueue adds a new task to the queue
func (ghidraScriptService *GhidraScriptService) AddToQueue(binaryName, script *string) {
	queueValue := newGhidraQueueItem(binaryName, script)
	ghidraScriptService.queue.PushBack(queueValue)
	ghidraScriptService.syncCondi.Signal()
}

func (ghidraScriptService *GhidraScriptService) findElement(binaryName *string) *list.Element {
	linkedListElement := ghidraScriptService.queue.Front()
	for {
		if linkedListElement != nil {
			task := linkedListElement.Value.(*queueScriptTask)
			if task.fileName == binaryName {
				return linkedListElement
			}

			linkedListElement = linkedListElement.Next()
		} else {
			break
		}
	}
	return nil
}

//RemoveFromQueue removes a task from the queue
func (ghidraScriptService *GhidraScriptService) RemoveFromQueue(binaryName *string) {
	linkedListElement := ghidraScriptService.findElement(binaryName)
	if linkedListElement != nil {
		ghidraScriptService.queue.Remove(linkedListElement)
	}
}

//UpdateTaskStatus updates the status of a task currently in the queue.
func (ghidraScriptService *GhidraScriptService) UpdateTaskStatus(binaryName *string, statusUpdate QueueStatus) {
	if linkedListElement := ghidraScriptService.findElement(binaryName); linkedListElement != nil {
		task := linkedListElement.Value.(*queueScriptTask)
		task.status = statusUpdate
	}
}

//GetStatus returns the current status of a task currently in the queue.
func (ghidraScriptService *GhidraScriptService) GetStatus(binaryName *string) *QueueStatus {
	if linkedListElement := ghidraScriptService.findElement(binaryName); linkedListElement != nil {
		status := linkedListElement.Value.(*queueScriptTask).status
		return &status
	}
	return nil
}

//GetAllStatus returns a map with the status of all task in the queue.
func (ghidraScriptService *GhidraScriptService) GetAllStatus() map[string]*QueueStatus {
	statusMap := make(map[string]*QueueStatus, ghidraScriptService.queue.Len())
	linkedListElement := ghidraScriptService.queue.Front()

	for {
		if linkedListElement != nil {
			task := linkedListElement.Value.(*queueScriptTask)
			statusMap[*task.fileName] = &task.status
			linkedListElement = linkedListElement.Next()
		} else {
			break
		}
	}
	return statusMap
}

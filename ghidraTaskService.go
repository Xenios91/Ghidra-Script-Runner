package ghidraScriptRunner

import (
	"container/list"
	"sync"
)

//GhidraTaskService the service utilized to run ghidra scripts and manage the task queue.
type GhidraTaskService struct {
	statusMap    map[string]*GhidraTaskStatus
	queue        *list.List
	syncCondi    *sync.Cond
	ghidraConfig *Configuration
}

func (ghidraTaskService *GhidraTaskService) waitForQueuedItems() {
	syncCondi := ghidraTaskService.syncCondi

	for {
		syncCondi.L.Lock()
		if ghidraTaskService.IsQueueEmpty() {
			syncCondi.Wait()
		} else {
			task := ghidraTaskService.getNextTaskAndRemoveFromQueue()
			ghidraTaskService.UpdateTaskStatusByTaskID((*task).GetTaskID(), runningStatus)
			success := (*task).RunTask(ghidraTaskService.ghidraConfig)
			if success {
				ghidraTaskService.UpdateTaskStatusByTaskID((*task).GetTaskID(), completeStatus)
			} else {
				ghidraTaskService.UpdateTaskStatusByTaskID((*task).GetTaskID(), errorStatus)
			}
		}
		syncCondi.L.Unlock()
	}
}

//NewGhidraTaskService used to create a new ghidra task service. Required to start polling task
func NewGhidraTaskService(config *Configuration) *GhidraTaskService {
	if config == nil {
		panic("Ghidra Config was nil!")
	}
	mutex := sync.Mutex{}
	newGhidraTaskService := &GhidraTaskService{make(map[string]*GhidraTaskStatus), list.New(), sync.NewCond(&mutex), config}
	go newGhidraTaskService.waitForQueuedItems()
	return newGhidraTaskService
}

//AddNewTaskToQueue takes a GhidraQueueTask and adds it to the task queue
func (ghidraTaskService *GhidraTaskService) AddNewTaskToQueue(task GhidraTask) {
	ghidraTaskService.statusMap[*task.GetTaskID()] = task.GetTaskStatus()
	ghidraTaskService.queue.PushBack(task)
	ghidraTaskService.syncCondi.Signal()
}

//AddToQueue adds a new task to the queue
func (ghidraTaskService *GhidraTaskService) AddToQueue(taskID, script *string) {
	task := NewGhidraScriptTask(taskID, script)
	ghidraTaskService.statusMap[*task.GetTaskID()] = task.GetTaskStatus()
	ghidraTaskService.queue.PushBack(task)
	ghidraTaskService.syncCondi.Signal()
}

func (ghidraTaskService *GhidraTaskService) findElementByTaskID(taskID *string) *list.Element {
	linkedListElement := ghidraTaskService.queue.Front()
	for {
		if linkedListElement != nil {
			task := linkedListElement.Value.(GhidraTask)
			if task.GetTaskID() == taskID {
				return linkedListElement
			}

			linkedListElement = linkedListElement.Next()
		} else {
			break
		}
	}
	return nil
}

//RemoveFromQueueByTaskID removes a task from the queue
func (ghidraTaskService *GhidraTaskService) RemoveFromQueueByTaskID(taskID *string) {
	linkedListElement := ghidraTaskService.findElementByTaskID(taskID)
	if linkedListElement != nil {
		ghidraTaskService.queue.Remove(linkedListElement)
		delete(ghidraTaskService.statusMap, *taskID)
	}
}

//UpdateTaskStatusByTaskID updates the status of a task currently in the queue.
func (ghidraTaskService *GhidraTaskService) UpdateTaskStatusByTaskID(taskID *string, statusUpdate GhidraTaskStatus) {
	if linkedListElement := ghidraTaskService.findElementByTaskID(taskID); linkedListElement != nil {
		task := linkedListElement.Value.(GhidraTask)
		task.SetTaskStatus(&statusUpdate)
	}
}

//GetStatusByTaskID returns the current status of a task currently in the queue.
func (ghidraTaskService *GhidraTaskService) GetStatusByTaskID(taskID *string) *GhidraTaskStatus {
	return ghidraTaskService.statusMap[*taskID]
}

//GetAllStatus returns a map with the status of all task in the queue.
func (ghidraTaskService *GhidraTaskService) GetAllStatus() map[string]*GhidraTaskStatus {
	return ghidraTaskService.statusMap
}

//IsQueueEmpty returns a bool indicating whether the queue is empty
func (ghidraTaskService *GhidraTaskService) IsQueueEmpty() bool {
	return ghidraTaskService.queue.Front() == nil
}

func (ghidraTaskService *GhidraTaskService) getNextTaskAndRemoveFromQueue() *GhidraTask {
	var task GhidraTask
	if frontElement := ghidraTaskService.queue.Front(); frontElement != nil {
		task = ghidraTaskService.queue.Remove(frontElement).(GhidraTask)
	}
	return &task
}

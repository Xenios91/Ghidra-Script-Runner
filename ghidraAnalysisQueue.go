package ghidraScriptRunner

type GhidraAnalysisQueue struct {
	queue        map[string]*ghidraQueueTask
	wait         chan bool
	ghidraConfig *Configuration
}

func (ghidraQueue *GhidraAnalysisQueue) waitForQueuedItems() {
	for {
		if len(ghidraQueue.queue) == 0 {
			<-ghidraQueue.wait
		} else {
			for key, value := range ghidraQueue.queue {
				ghidraQueue.UpdateTaskStatus(&key, runningStatus)
				success := value.runScript(ghidraQueue.ghidraConfig)
				if success {
					ghidraQueue.UpdateTaskStatus(&key, completeStatus)
				} else {
					ghidraQueue.UpdateTaskStatus(&key, errorStatus)
				}
			}
		}

	}
}

//NewGhidraAnalysisQueue Used to create a new analysis queue. Required to start polling,
func NewGhidraAnalysisQueue(config *Configuration) *GhidraAnalysisQueue {
	newAnalysisQueue := &GhidraAnalysisQueue{make(map[string]*ghidraQueueTask), make(chan bool), config}
	go newAnalysisQueue.waitForQueuedItems()
	return newAnalysisQueue
}

//AddToQueue Adds a new analysis task to the queue
func (ghidraQueue *GhidraAnalysisQueue) AddToQueue(binaryName, script *string) {
	queueValue := newGhidraQueueItem(binaryName, script)
	ghidraQueue.queue[*binaryName] = queueValue
	ghidraQueue.wait <- true
}

//RemoveFromQueue removes a binary name from the queue of binaries being processed by ghidra.
func (ghidraQueue *GhidraAnalysisQueue) RemoveFromQueue(binaryName *string) {
	delete((*ghidraQueue).queue, *binaryName)
}

//UpdateTaskStatus updates the status of a binary currently in the queue.
func (ghidraQueue *GhidraAnalysisQueue) UpdateTaskStatus(binaryName *string, statusUpdate queueStatus) {
	queueValue := ghidraQueue.queue[*binaryName]
	if queueValue != nil {
		ghidraQueue.queue[*binaryName].status = statusUpdate
	}
}

//GetStatus Returns the current status of a binary being processed by Ghidra.
func (ghidraQueue *GhidraAnalysisQueue) GetStatus(binaryName *string) *queueStatus {
	var status queueStatus
	queueValue := ghidraQueue.queue[*binaryName]
	if queueValue != nil {
		status = (*queueValue).status
	}
	return &status
}

//GetAllStatus Returns a map with the status of all binaries being processed by Ghidra.
func (ghidraQueue *GhidraAnalysisQueue) GetAllStatus() map[string]*queueStatus {
	statusMap := make(map[string]*queueStatus, len(ghidraQueue.queue))

	for key, value := range ghidraQueue.queue {
		statusMap[key] = &value.status
	}

	return statusMap
}

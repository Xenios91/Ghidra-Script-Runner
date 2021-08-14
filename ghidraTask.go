package ghidraScriptRunner

//GhidraQueueTask a single task for the queue
type GhidraQueueTask interface {
	RunTask(config *Configuration) bool
	GetTaskID() *string
	GetTaskStatus() *QueueStatus
	SetTaskStatus(status *QueueStatus)
}

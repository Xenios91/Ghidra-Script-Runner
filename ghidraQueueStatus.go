package ghidraScriptRunner

//QueueStatus the status of a task in the queue
type QueueStatus string

const (
	beginningStatus QueueStatus = "Waiting on Ghidra"
	errorStatus     QueueStatus = "An error has occurred"
	runningStatus   QueueStatus = "Running"
	completeStatus  QueueStatus = "Task Complete"
)

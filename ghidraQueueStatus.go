package ghidraScriptRunner

//QueueStatus the status of a task in the queue
type QueueStatus string

const (
	beginningStatus QueueStatus = "Waiting on Ghidra"
	errorStatus     QueueStatus = "An error has occured"
	runningStatus   QueueStatus = "Running"
	completeStatus  QueueStatus = "Script Complete"
)

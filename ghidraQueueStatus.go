package ghidraScriptRunner

//GhidraTaskStatus the status of a task in the queue
type GhidraTaskStatus string

const (
	waitingStatus  GhidraTaskStatus = "Waiting on Ghidra"
	errorStatus    GhidraTaskStatus = "An error has occurred"
	runningStatus  GhidraTaskStatus = "Running"
	completeStatus GhidraTaskStatus = "Task Complete"
)

package ghidraScriptRunner

type queueStatus string

const (
	beginningStatus queueStatus = "Waiting on Ghidra"
	errorStatus     queueStatus = "An error has occured"
	runningStatus   queueStatus = "Running"
	completeStatus  queueStatus = "Script Complete"
)

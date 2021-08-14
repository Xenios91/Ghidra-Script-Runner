package ghidraScriptRunner

//GhidraTask a single task for the queue
type GhidraTask interface {
	RunTask(config *Configuration) bool
	GetTaskID() *string
	GetTaskStatus() *GhidraTaskStatus
	SetTaskStatus(status *GhidraTaskStatus)
}

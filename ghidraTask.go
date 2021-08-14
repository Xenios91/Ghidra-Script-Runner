package ghidraScriptRunner

//GhidraQueueTask a single task for the queue
type GhidraQueueTask interface {
	runTask()
}

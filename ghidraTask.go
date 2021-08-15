package ghidraScriptRunner

//GhidraTask a single task for the queue
type GhidraTask interface {
	Run(config *Configuration) error
	ID() string
	Status() *GhidraTaskStatus
	SetStatus(status *GhidraTaskStatus)
}

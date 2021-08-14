package ghidraScriptRunner

import "os/exec"

//GhidraScriptTask a single task for ghidra to execute
type GhidraScriptTask struct {
	fileName *string
	script   *string
	status   GhidraTaskStatus
}

//NewGhidraScriptTask returns a new GhidraScriptTask struct, should be used instead of calling new
func NewGhidraScriptTask(fileName, script *string) *GhidraScriptTask {
	return &GhidraScriptTask{fileName, script, waitingStatus}
}

//RunTask run the task assigned to this GhidraScriptTask struct
func (queuedItem *GhidraScriptTask) RunTask(ghidraConfig *Configuration) (successful bool) {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation,
		*ghidraConfig.ghidraProject, "-import", *queuedItem.fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	return err != nil
}

//GetTaskID returns the taskID associated with this GhidraScriptTask struct
func (queuedItem *GhidraScriptTask) GetTaskID() *string {
	return queuedItem.fileName
}

//GetTaskStatus returns the GhidraTaskStatus associated with this GhidraScriptTask struct
func (queuedItem *GhidraScriptTask) GetTaskStatus() *GhidraTaskStatus {
	return &queuedItem.status
}

//SetTaskStatus sets the GhidraTaskStatus of this GhidraScriptTask to the argument passed to this method
func (queuedItem *GhidraScriptTask) SetTaskStatus(queueStatus *GhidraTaskStatus) {
	queuedItem.status = *queueStatus
}

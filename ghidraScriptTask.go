package ghidraScriptRunner

import (
	"os/exec"
)

//GhidraScriptTask a single task for ghidra to execute
type GhidraScriptTask struct {
	fileName *string
	status   GhidraTaskStatus
	script   *string
}

//NewGhidraScriptTask returns a new GhidraScriptTask struct, should be used instead of calling new
func NewGhidraScriptTask(fileName, script *string) *GhidraScriptTask {
	return &GhidraScriptTask{fileName, waitingStatus, script}
}

//RunTask run the task assigned to this GhidraScriptTask struct
func (ghidraScriptTask *GhidraScriptTask) RunTask(ghidraConfig *Configuration) (successful bool) {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation,
		*ghidraConfig.ghidraProject, "-import", *ghidraScriptTask.fileName, "-postScript", *ghidraScriptTask.script, "-overwrite").Start()
	return err != nil
}

//GetTaskID returns the taskID associated with this GhidraScriptTask struct
func (ghidraScriptTask *GhidraScriptTask) GetTaskID() *string {
	return ghidraScriptTask.fileName
}

//GetTaskStatus returns the GhidraTaskStatus associated with this GhidraScriptTask struct
func (ghidraScriptTask *GhidraScriptTask) GetTaskStatus() *GhidraTaskStatus {
	return &ghidraScriptTask.status
}

//SetTaskStatus sets the GhidraTaskStatus of this GhidraScriptTask to the argument passed to this method
func (ghidraScriptTask *GhidraScriptTask) SetTaskStatus(queueStatus *GhidraTaskStatus) {
	ghidraScriptTask.status = *queueStatus
}

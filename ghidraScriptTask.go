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

func (queuedItem *GhidraScriptTask) RunTask(ghidraConfig *Configuration) (successful bool) {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation,
		*ghidraConfig.ghidraProject, "-import", *queuedItem.fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	return err != nil
}

func (queuedItem *GhidraScriptTask) GetTaskID() *string {
	return queuedItem.fileName
}

func (queuedItem *GhidraScriptTask) GetTaskStatus() *GhidraTaskStatus {
	return &queuedItem.status
}

func (queuedItem *GhidraScriptTask) SetTaskStatus(queueStatus *GhidraTaskStatus) {
	queuedItem.status = *queueStatus
}

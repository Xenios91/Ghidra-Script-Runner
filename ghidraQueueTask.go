package ghidraScriptRunner

import "os/exec"

type queueScriptTask struct {
	fileName *string
	script   *string
	status   GhidraTaskStatus
}

func NewGhidraScriptTask(fileName, script *string) *queueScriptTask {
	return &queueScriptTask{fileName, script, waitingStatus}
}

func (queuedItem *queueScriptTask) RunTask(ghidraConfig *Configuration) (successful bool) {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation,
		*ghidraConfig.ghidraProject, "-import", *queuedItem.fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	return err != nil
}

func (queuedItem *queueScriptTask) GetTaskID() *string {
	return queuedItem.fileName
}

func (queuedItem *queueScriptTask) GetTaskStatus() *GhidraTaskStatus {
	return &queuedItem.status
}

func (queuedItem *queueScriptTask) SetTaskStatus(queueStatus *GhidraTaskStatus) {
	queuedItem.status = *queueStatus
}

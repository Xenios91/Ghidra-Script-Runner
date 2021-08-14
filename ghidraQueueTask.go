package ghidraScriptRunner

import "os/exec"

type queueScriptTask struct {
	fileName *string
	script   *string
	status   QueueStatus
}

func newGhidraQueueItem(fileName, script *string) *queueScriptTask {
	return &queueScriptTask{fileName, script, beginningStatus}
}

func (queuedItem *queueScriptTask) runTask(ghidraConfig *Configuration) (successful bool) {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation,
		*ghidraConfig.ghidraProject, "-import", *queuedItem.fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	return err != nil
}

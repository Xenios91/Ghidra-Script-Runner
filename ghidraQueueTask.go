package ghidraScriptRunner

import "os/exec"

type ghidraQueueTask struct {
	fileName *string
	script   *string
	status   queueStatus
}

func newGhidraQueueItem(fileName, script *string) *ghidraQueueTask {
	return &ghidraQueueTask{fileName, script, beginningStatus}
}

func (queuedItem *ghidraQueueTask) runScript(ghidraConfig *Configuration) (successful bool) {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation,
		*ghidraConfig.ghidraProject, "-import", *queuedItem.fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	return err != nil
}

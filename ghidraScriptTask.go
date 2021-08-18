package ghidraScriptRunner

import (
	"os/exec"
)

//GhidraScriptTask a single task for ghidra to execute
type GhidraScriptTask struct {
	fileName string
	status   GhidraTaskStatus
	script   string
}

//NewGhidraScriptTask returns a new GhidraScriptTask struct, should be used instead of calling new
func NewGhidraScriptTask(fileName, script string) *GhidraScriptTask {
	return &GhidraScriptTask{fileName, waitingStatus, script}
}

//Run run the task assigned to this GhidraScriptTask struct
func (ghidraScriptTask *GhidraScriptTask) Run(ghidraConfig *Configuration) error {
	var err error
	if ghidraConfig.shouldOverWrite {
		err = exec.Command(ghidraConfig.ghidraHeadless, ghidraConfig.ghidraProjectLocation,
			ghidraConfig.ghidraProject, "-import", ghidraScriptTask.fileName, "-postScript", ghidraScriptTask.script, "-overwrite").Start()
	} else {
		err = exec.Command(ghidraConfig.ghidraHeadless, ghidraConfig.ghidraProjectLocation,
			ghidraConfig.ghidraProject, "-import", ghidraScriptTask.fileName, "-postScript", ghidraScriptTask.script).Start()
	}
	return err
}

//ID returns the ID associated with this GhidraScriptTask struct
func (ghidraScriptTask *GhidraScriptTask) ID() string {
	return ghidraScriptTask.fileName
}

//Status returns the Status associated with this GhidraScriptTask struct
func (ghidraScriptTask *GhidraScriptTask) Status() *GhidraTaskStatus {
	return &ghidraScriptTask.status
}

//SetStatus sets the Status of this GhidraScriptTask to the argument passed to this method
func (ghidraScriptTask *GhidraScriptTask) SetStatus(queueStatus *GhidraTaskStatus) {
	ghidraScriptTask.status = *queueStatus
}

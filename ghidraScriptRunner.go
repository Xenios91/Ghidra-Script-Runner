package ghidraScriptRunner

import (
	"log"
	"os/exec"
	"path/filepath"
	"sync"
)

type ghidraAnalysisQueue map[string]*ghidraQueueValue

type ghidraQueueValue struct {
	script *string
	status *string
}

var (
	once         sync.Once
	ghidraQueue  ghidraAnalysisQueue
	ghidraConfig *Configuration
)

//StartGhidraAnalysis Starts analysis on a supplied binary and returns a boolean value indicating if the analysis was successfully started.
func StartGhidraAnalysis(fileName *string, script *string) error {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation, *ghidraConfig.ghidraProject, "-import", *fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	if err != nil {
		return err
	}
	addToQueue(filepath.Base(*fileName), script)
	return nil
}

//LoadGhidraAnalysis sets all configuration information for ghidra based on arguments supplied.
func LoadGhidraAnalysis(configuration *Configuration) {
	once.Do(func() {
		log.Print("Loading Ghidra analysis queue... ")
		ghidraQueue = make(ghidraAnalysisQueue)
		ghidraConfig = new(Configuration)
		log.Println("Ghidra Analysis Queue successfully loaded!")
	})
}

func addToQueue(binaryName string, script *string) {
	queueValue := new(ghidraQueueValue)
	beginningStatus := "Waiting on Ghidra"
	queueValue.status = &beginningStatus
	queueValue.script = script
	ghidraQueue[binaryName] = queueValue
}

//RemoveFromQueue removes a binary name from the queue of binaries being processed by ghidra.
func RemoveFromQueue(binaryName string) {
	delete(ghidraQueue, binaryName)
}

//UpdateQueue updates the status of a binary currently in the queue.
func UpdateQueue(binaryName *string, statusUpdate *string) {
	queueValue := ghidraQueue[*binaryName]
	if queueValue != nil {
		ghidraQueue[*binaryName].status = statusUpdate
	}
}

//GetStatus Returns the current status of a binary being processed by Ghidra.
func GetStatus(binaryName *string) *string {
	var status string
	queueValue := ghidraQueue[*binaryName]
	if queueValue != nil {
		status = *queueValue.status
	}
	return &status
}

//GetAllStatus Returns a map with the status of all binaries being processed by Ghidra.
func GetAllStatus() map[string]*string {
	statusMap := make(map[string]*string, len(ghidraQueue))

	for key, element := range ghidraQueue {
		statusMap[key] = element.status
	}

	return statusMap
}

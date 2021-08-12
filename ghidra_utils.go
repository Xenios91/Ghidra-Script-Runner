package ghidrautils

import (
	config "glyph/glyph/config"
	"log"
	"os/exec"
	"path/filepath"
	"sync"
)

type ghidraAnalysisQueue map[string]*ghidraQueueValue

type ghidraQueueValue struct {
	isTrainingData bool
	status         *string
}

type ghidraAnalysisConfig struct {
	ghidraHeadless        *string
	ghidraProjectLocation *string
	ghidraProject         *string
	ghidraScript          *string
}

var (
	once         sync.Once
	ghidraQueue  ghidraAnalysisQueue
	ghidraConfig *ghidraAnalysisConfig
)

func main() {
	fmt.Println("well, hello there")
}

//StartGhidraAnalysis Starts analysis on a supplied binary and returns a boolean value indicating if the analysis was successfully started.
func StartGhidraAnalysis(fileName string, trainingData bool) bool {
	err := exec.Command(*ghidraConfig.ghidraHeadless, *ghidraConfig.ghidraProjectLocation, *ghidraConfig.ghidraProject, "-import", fileName, "-postScript", *ghidraConfig.ghidraScript, "-overwrite").Start()
	if err != nil {
		log.Println(err)
		return false
	}
	addToQueue(filepath.Base(fileName), trainingData)
	return true
}

//LoadGhidraAnalysis sets all configuration information for ghidra based on arguments supplied.
func LoadGhidraAnalysis(configuration *config.Configuration) {
	once.Do(func() {
		log.Print("Loading Ghidra analysis queue... ")
		ghidraQueue = make(ghidraAnalysisQueue)
		ghidraConfig = new(ghidraAnalysisConfig)
		ghidraConfig.ghidraHeadless = configuration.GhidraHeadless
		ghidraConfig.ghidraProjectLocation = configuration.GhidraProjectLocation
		ghidraConfig.ghidraProject = configuration.GhidraProject
		ghidraConfig.ghidraScript = configuration.GhidraScript
		log.Println("Ghidra Analysis Queue successfully loaded!")
	})
}

func addToQueue(binaryName string, trainingData bool) {
	queueValue := new(ghidraQueueValue)
	beginningStatus := "Waiting on Ghidra"
	queueValue.status = &beginningStatus
	queueValue.isTrainingData = trainingData
	ghidraQueue[binaryName] = queueValue
}

//RemoveFromQueue removes a binary name from the queue of binaries being processed by ghidra.
func RemoveFromQueue(binaryName string) {
	delete(ghidraQueue, binaryName)
}

func checkIfTraining(binaryName *string) bool {
	var isTraining bool
	queueValue := ghidraQueue[*binaryName]

	if queueValue != nil {
		isTraining = queueValue.isTrainingData
	} else {
		isTraining = false
	}
	return isTraining
}

//CheckIfTrainingAndRemove returns true/false if a binary being processed by ghidra currently is training data, and removes it from the queue.
func CheckIfTrainingAndRemove(binaryName string) bool {
	isTraining := checkIfTraining(&binaryName)
	RemoveFromQueue(binaryName)
	return isTraining
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

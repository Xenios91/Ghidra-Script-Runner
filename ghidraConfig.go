package ghidraScriptRunner

import "errors"

//Configuration The ghidra configuration to be utilized for running scripts
type Configuration struct {
	ghidraHeadless        string
	ghidraProjectLocation string
	ghidraProject         string
	shouldOverWrite       bool
}

//NewConfiguration returns a new ghidra configuration
func NewConfiguration(ghidraHeadless, ghidraProjectLocation, ghidraProject string, shouldOverWrite bool) (*Configuration, error) {
	config := &Configuration{ghidraHeadless, ghidraProjectLocation, ghidraProject, shouldOverWrite}
	if config.checkConfig() {
		return config, nil
	}
	return nil, errors.New("invalid configuration")
}

func (config *Configuration) checkConfig() (validConfig bool) {
	validConfig = len(config.ghidraHeadless) > 0 && len(config.ghidraProjectLocation) > 0 && len(config.ghidraProject) > 0
	return validConfig
}

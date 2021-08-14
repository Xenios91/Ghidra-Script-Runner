package ghidraScriptRunner

import "errors"

type Configuration struct {
	ghidraHeadless        *string
	ghidraProjectLocation *string
	ghidraProject         *string
	ghidraScript          *string
}

func NewConfiguration(ghidraHeadless, ghidraProjectLocation, ghidraProject, ghidraScript string) (*Configuration, error) {
	config := &Configuration{ghidraHeadless: &ghidraHeadless, ghidraProjectLocation: &ghidraProjectLocation,
		ghidraProject: &ghidraProject, ghidraScript: &ghidraScript}
	if config.checkConfig() {
		return config, nil
	}
	return nil, errors.New("invalid configuration")
}

func (config *Configuration) checkConfig() (validConfig bool) {
	validConfig = config.ghidraHeadless != nil && config.ghidraProjectLocation != nil && config.ghidraProject != nil && config.ghidraScript != nil
	return validConfig
}

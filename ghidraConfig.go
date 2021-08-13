package ghidraScriptRunner

type Configuration struct {
	ghidraHeadless        *string
	ghidraProjectLocation *string
	ghidraProject         *string
	ghidraScript          *string
}

func NewConfiguration(ghidraHeadless, ghidraProjectLocation, ghidraProject, ghidraScript string) *Configuration {
	return &Configuration{ghidraHeadless: &ghidraHeadless, ghidraProjectLocation: &ghidraProjectLocation,
		ghidraProject: &ghidraProject, ghidraScript: &ghidraScript}
}

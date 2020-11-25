package presets

// FakeParser implements all fake behaviors for using parser in tests.
type FakeParser struct {
	CalledExists, CalledLookUpFiles, CalledWriteFiles, CalledGetCreateCommand, CalledGetPresets, CalledGetLanguages, CalledLoadPresets bool

	MockExists        bool
	MockFoundFiles    []string
	MockFileError     string
	MockError         error
	MockCreateCommand string
	MockLanguages     []string
	MockPresets       []string
	MockAllPresets    map[string]map[string]string
}

// Exists check if preset exists
func (f *FakeParser) Exists(preset string) (exists bool) {
	f.CalledExists = true
	exists = f.MockExists
	return
}

// GetCreateCommand gets the command to create a new project
func (f *FakeParser) GetCreateCommand(preset string) (cmd string, err error) {
	f.CalledGetCreateCommand = true
	cmd = f.MockCreateCommand
	return
}

// GetLanguages get all presets languages
func (f *FakeParser) GetLanguages() (languages []string) {
	f.CalledGetLanguages = true
	languages = f.MockLanguages
	return
}

// GetPresets get all presets names
func (f *FakeParser) GetPresets(language string) (presets []string) {
	f.CalledGetPresets = true
	presets = f.MockPresets
	return
}

// LookUpFiles check if preset files exist
func (f *FakeParser) LookUpFiles(preset string) (foundFiles []string) {
	f.CalledLookUpFiles = true
	foundFiles = f.MockFoundFiles
	return
}

// WriteFiles write preset files
func (f *FakeParser) WriteFiles(preset string) (fileError string, err error) {
	f.CalledWriteFiles = true
	fileError = f.MockFileError
	err = f.MockError
	return
}

//LoadPresets loads all presets
func (f *FakeParser) LoadPresets(presets map[string]map[string]string) {
	f.CalledLoadPresets = true
	f.MockAllPresets = presets
}

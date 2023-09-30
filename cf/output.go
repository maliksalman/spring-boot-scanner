package cf

import (
	"encoding/json"
	"os"
)

type JavaApp struct {
	Org                 string
	Space               string
	Name                string
	Instances           int
	State               string
	Type                string
	JavaRuntimeVersion  string
	JavaCompilerVersion string
	SpringBootVersion   string
}

type OtherApp struct {
	Org       string
	Space     string
	Name      string
	Instances int
	State     string
	Type      string
}

type ScannerOutput struct {
	apps []any
}

func NewScannerOutput() *ScannerOutput {
	return &ScannerOutput{apps: make([]any, 0)}
}

func (o *ScannerOutput) AddJavaApp(org string, space string, name string, instances int, state string, runtimeVersion string, compilerVersion string, springBootVersion string) {
	o.apps = append(o.apps, JavaApp{
		Org:                 org,
		Space:               space,
		Name:                name,
		Instances:           instances,
		State:               state,
		Type:                "java",
		JavaRuntimeVersion:  runtimeVersion,
		JavaCompilerVersion: compilerVersion,
		SpringBootVersion:   springBootVersion,
	})
}

func (o *ScannerOutput) AddOtherApp(org string, space string, name string, instances int, state string, appType string) {
	o.apps = append(o.apps, OtherApp{
		Org:       org,
		Space:     space,
		Name:      name,
		Instances: instances,
		State:     state,
		Type:      appType,
	})
}

func (o *ScannerOutput) WriteAsJSON(filename string) {
	jsonBytes, _ := json.Marshal(o.apps)
	os.WriteFile(filename, jsonBytes, 0644)
}

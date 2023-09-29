package main

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
	apps     []any
	filename string
}

func NewScannerOutput(filename string) *ScannerOutput {
	return &ScannerOutput{apps: make([]any, 0), filename: filename}
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

func (o *ScannerOutput) writeAsJSON() {
	jsonBytes, _ := json.Marshal(o.apps)
	os.WriteFile(o.filename, jsonBytes, 0644)
}

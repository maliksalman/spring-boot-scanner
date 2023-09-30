package cf

import (
	"encoding/json"
	"github.com/anchore/clio"
	"github.com/anchore/syft/cmd/syft/cli"
	"log"
	"os"
)

func findJavaRuntimeAndBootVersions(dropletBytes []byte) (string, string) {

	dropletFile, _ := os.CreateTemp("", "app-droplet-*.tgz")
	os.WriteFile(dropletFile.Name(), dropletBytes, 0644)
	log.Printf("Wrote droplet: %s", dropletFile.Name())

	syftId := clio.Identification{
		Name:           "spring-boot-scanner",
		Version:        "v1.0.0",
		BuildDate:      "some-date",
		GitCommit:      "some-commit",
		GitDescription: "some-git",
	}

	syftFile, _ := os.CreateTemp("", "syft-*.json")

	command := cli.Command(syftId)
	command.SetArgs([]string{
		"packages",
		dropletFile.Name(),
		"--output",
		"syft-json",
		"--file",
		syftFile.Name(),
	})
	command.Execute()

	log.Printf("Analyzed droplet: File=%s\n", syftFile.Name())
	return findJavaRuntimeAndBootVersionsFromSyftJson(syftFile.Name())
}

type SyftOutput struct {
	Artifacts []SyftArtifact `json:"artifacts"`
}
type SyftArtifact struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

func findJavaRuntimeAndBootVersionsFromSyftJson(syftJsonFileName string) (string, string) {

	var result SyftOutput
	jsonBytes, _ := os.ReadFile(syftJsonFileName)
	json.Unmarshal(jsonBytes, &result)

	javaVersion := ""
	for _, artifact := range result.Artifacts {
		if artifact.Type == "binary" && artifact.Name == "java" {
			javaVersion = artifact.Version
			break
		}
	}

	bootVersion := ""
	for _, artifact := range result.Artifacts {
		if artifact.Type == "java-archive" && artifact.Name == "spring-boot" {
			bootVersion = artifact.Version
			break
		}
	}

	return javaVersion, bootVersion
}

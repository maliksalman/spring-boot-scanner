package scan

import (
	"encoding/json"
	"github.com/anchore/clio"
	"github.com/anchore/syft/cmd/syft/cli/commands"
	"io"
	"log"
	"os"
	"strings"
)

type RuntimeInfo struct {
	JavaRuntimeVersion   string
	BootVersion          string
	BootContentExtracted bool
	BootContentPrefix    string
	BootJar              string
}

func FindRuntimeInfoFromPath(contentPath string) RuntimeInfo {
	outputFile := runSyftCommand(contentPath)
	return findRuntimeInfoFromSyftJsonFile(outputFile)
}

func FindRuntimeInfoFromContent(dropletBytes []byte) RuntimeInfo {

	dropletFile, _ := os.CreateTemp("", "app-droplet-*.tgz")
	os.WriteFile(dropletFile.Name(), dropletBytes, 0644)
	log.Printf("Wrote droplet: %s", dropletFile.Name())

	outputFile := runSyftCommand(dropletFile.Name())
	return findRuntimeInfoFromSyftJsonFile(outputFile)
}

func runSyftCommand(contentPath string) *os.File {

	id := clio.Identification{
		Name:           "spring-boot-scanner",
		Version:        "v1.0.0",
		BuildDate:      "some-date",
		GitCommit:      "some-commit",
		GitDescription: "some-git",
	}
	cfg := clio.NewSetupConfig(id).
		WithGlobalConfigFlag().
		WithNoLogging()
	app := clio.New(*cfg)
	command := commands.Packages(app)

	outputFile, _ := os.CreateTemp("", "syft-*.json")
	command.SetArgs([]string{
		contentPath,
		"--output",
		"syft-json",
		"--file",
		outputFile.Name(),
	})
	command.Execute()

	log.Printf("Analyzed: Content=%s Output=%s\n", contentPath, outputFile.Name())
	return outputFile
}

type syftAtifactMetadata struct {
	VirtualPath string `json:"virtualPath"`
}

type syftArtifact struct {
	Name         string              `json:"name"`
	ArtifcatType string              `json:"type"`
	Version      string              `json:"version"`
	FoundBy      string              `json:"foundBy"`
	Metadata     syftAtifactMetadata `json:"metadata"`
}

type syftOutput struct {
	Artifacts []syftArtifact `json:"artifacts"`
}

func findRuntimeInfoFromSyftJsonFile(syftJsonFile *os.File) RuntimeInfo {

	var result syftOutput
	jsonBytes, _ := io.ReadAll(syftJsonFile)
	json.Unmarshal(jsonBytes, &result)

	javaVersion := ""
	bootVersion := ""
	bootJar := ""
	bootContentExtracted := false
	bootContentPrefix := ""

	for _, artifact := range result.Artifacts {
		if (artifact.ArtifcatType == "java-archive" && artifact.Name == "rt") ||
			(artifact.ArtifcatType == "java-archive" && artifact.Name == "jrt-fs") {
			javaVersion = artifact.Version
			break
		}
	}

	for _, artifact := range result.Artifacts {
		if artifact.ArtifcatType == "java-archive" && artifact.Name == "spring-boot" && artifact.FoundBy == "java-cataloger" {
			bootVersion = artifact.Version
			jar, springBootLocation, found := strings.Cut(artifact.Metadata.VirtualPath, ":")
			if found {
				bootJar = jar
				if strings.HasPrefix(bootJar, "/") {
					bootJar = bootJar[1:]
				}

				bootContentPrefix = findBootContentPrefix(springBootLocation)
			} else {
				bootContentExtracted = true
				bootContentPrefix = findBootContentPrefix(artifact.Metadata.VirtualPath)
			}
			break
		}
	}

	return RuntimeInfo{
		JavaRuntimeVersion:   javaVersion,
		BootVersion:          bootVersion,
		BootContentExtracted: bootContentExtracted,
		BootContentPrefix:    bootContentPrefix,
		BootJar:              bootJar,
	}
}

func findBootContentPrefix(springBootLocation string) string {

	if strings.HasPrefix(springBootLocation, "/") {
		springBootLocation = springBootLocation[1:]
	}

	bootIndex := strings.Index(springBootLocation, "BOOT-INF/lib/")
	if bootIndex != -1 {
		if bootIndex == 0 {
			return "BOOT-INF"
		} else {
			return springBootLocation[0:bootIndex] + "BOOT-INF"
		}
	}

	webIndex := strings.Index(springBootLocation, "WEB-INF/lib/")
	if webIndex != -1 {
		if webIndex == 0 {
			return "WEB-INF"
		} else {
			return springBootLocation[0:webIndex] + "WEB-INF"
		}
	}

	return ""
}

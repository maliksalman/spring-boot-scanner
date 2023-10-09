package scan

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setupTest(t *testing.T, rtName string, sbName string, sbPath string) *os.File {

	json := fmt.Sprintf(`
{
    "Artifacts": [
        {
            "Name": "%s",
            "Version": "17.0.6",
            "type": "java-archive",
            "FoundBy": "java-cataloger",
            "Metadata": {
                "VirtualPath": "/layers/paketo-buildpacks_bellsoft-liberica/jre/lib/jrt-fs.jar"
            }
        },
        {
            "Name": "%s",
            "Version": "3.0.3",
            "type": "java-archive",
            "FoundBy": "java-cataloger",
            "Metadata": {
                "VirtualPath": "%s"
            }
        }
    ]
}`, rtName, sbName, sbPath)

	jsonFile, err := os.CreateTemp("", "test-*.json")
	assert.NoError(t, err)

	err = os.WriteFile(jsonFile.Name(), []byte(json), 0644)
	assert.NoError(t, err)

	return jsonFile
}

func TestFindRuntimeInfoFromSyftJsonFile_BootAppExtractedInDockerImage(t *testing.T) {

	// GIVEN
	testFile := setupTest(t, "jrt-fs", "spring-boot", "/workspace/BOOT-INF/lib/spring-boot-3.0.3.jar")

	// WHEN
	runtimeInfo := findRuntimeInfoFromSyftJsonFile(testFile)

	// THEN
	assert.Equal(t, "3.0.3", runtimeInfo.BootVersion)
	assert.Equal(t, "17.0.6", runtimeInfo.JavaRuntimeVersion)
	assert.True(t, runtimeInfo.BootContentExtracted)
	assert.Equal(t, "workspace/BOOT-INF", runtimeInfo.BootContentPrefix)
}

func TestFindRuntimeInfoFromSyftJsonFile_BootAppExtractedInDroplet(t *testing.T) {

	// GIVEN
	testFile := setupTest(t, "jrt-fs", "spring-boot", "app/BOOT-INF/lib/spring-boot-3.0.3.jar")

	// WHEN
	runtimeInfo := findRuntimeInfoFromSyftJsonFile(testFile)

	// THEN
	assert.Equal(t, "3.0.3", runtimeInfo.BootVersion)
	assert.Equal(t, "17.0.6", runtimeInfo.JavaRuntimeVersion)
	assert.True(t, runtimeInfo.BootContentExtracted)
	assert.Equal(t, "app/BOOT-INF", runtimeInfo.BootContentPrefix)
}

func TestFindRuntimeInfoFromSyftJsonFile_BootAppInBootJarInDockerImage(t *testing.T) {

	// GIVEN
	testFile := setupTest(t, "jrt-fs", "spring-boot", "/boot-app.jar:BOOT-INF/lib/spring-boot-3.0.3.jar")

	// WHEN
	runtimeInfo := findRuntimeInfoFromSyftJsonFile(testFile)

	// THEN
	assert.Equal(t, "3.0.3", runtimeInfo.BootVersion)
	assert.Equal(t, "17.0.6", runtimeInfo.JavaRuntimeVersion)
	assert.False(t, runtimeInfo.BootContentExtracted)
	assert.Equal(t, "BOOT-INF", runtimeInfo.BootContentPrefix)
	assert.Equal(t, "boot-app.jar", runtimeInfo.BootJar)
}

func TestFindRuntimeInfoFromSyftJsonFile_BootAppInBootJarInDockerImageWithOlderRuntime(t *testing.T) {

	// GIVEN
	testFile := setupTest(t, "rt", "spring-boot", "/boot-app.jar:BOOT-INF/lib/spring-boot-3.0.3.jar")

	// WHEN
	runtimeInfo := findRuntimeInfoFromSyftJsonFile(testFile)

	// THEN
	assert.Equal(t, "3.0.3", runtimeInfo.BootVersion)
	assert.Equal(t, "17.0.6", runtimeInfo.JavaRuntimeVersion)
	assert.False(t, runtimeInfo.BootContentExtracted)
	assert.Equal(t, "BOOT-INF", runtimeInfo.BootContentPrefix)
	assert.Equal(t, "boot-app.jar", runtimeInfo.BootJar)
}

func TestFindRuntimeInfoFromSyftJsonFile_NotBootAppInDockerImage(t *testing.T) {

	// GIVEN
	testFile := setupTest(t, "jrt-fs", "not-spring-boot", "/not-boot-app.jar")

	// WHEN
	runtimeInfo := findRuntimeInfoFromSyftJsonFile(testFile)

	// THEN
	assert.Equal(t, "", runtimeInfo.BootVersion)
	assert.Equal(t, "17.0.6", runtimeInfo.JavaRuntimeVersion)
	assert.False(t, runtimeInfo.BootContentExtracted)
	assert.Equal(t, "", runtimeInfo.BootContentPrefix)
	assert.Equal(t, "", runtimeInfo.BootJar)
}

func TestFindRuntimeInfoFromSyftJsonFile_NotJavaAppInDockerImage(t *testing.T) {

	// GIVEN
	testFile := setupTest(t, "not-rt", "not-spring-boot", "/not-boot-app.exe")

	// WHEN
	runtimeInfo := findRuntimeInfoFromSyftJsonFile(testFile)

	// THEN
	assert.Equal(t, "", runtimeInfo.BootVersion)
	assert.Equal(t, "", runtimeInfo.JavaRuntimeVersion)
	assert.False(t, runtimeInfo.BootContentExtracted)
	assert.Equal(t, "", runtimeInfo.BootContentPrefix)
	assert.Equal(t, "", runtimeInfo.BootJar)
}

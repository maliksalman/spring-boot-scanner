package scan

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetClassFileVersion_Java8(t *testing.T) {
	testGetClassFileVersion(t, "8")
}

func TestGetClassFileVersion_Java11(t *testing.T) {
	testGetClassFileVersion(t, "11")
}

func TestGetClassFileVersion_Java17(t *testing.T) {
	testGetClassFileVersion(t, "17")
}

func testGetClassFileVersion(t *testing.T, expectedVer string) {

	// GIVEN
	classFile, err := os.Open(filepath.Join("test", "HelloWorld."+expectedVer+".class"))
	assert.NoError(t, err)

	// WHEN
	found, ver := getClassFileVersion(classFile, "HelloWorld.class")

	// THEN
	assert.Equal(t, expectedVer, ver)
	assert.True(t, found)
}

func TestFindJavaCompilerVersionFromBootJar(t *testing.T) {

	ver, err := FindJavaCompilerVersionFromBootJar(filepath.Join("test", "app.jar"), "BOOT-INF")

	assert.NoError(t, err)
	assert.Equal(t, "11", ver)
}

func testFindJavaCompilerVersionFromReader(t *testing.T, tgzFileName string, bootContentPrefix string) {

	tgzFile, err := os.Open(filepath.Join("test", tgzFileName))
	assert.NoError(t, err)

	found, ver := FindJavaCompilerVersionFromReader(tgzFile, bootContentPrefix)

	assert.True(t, found)
	assert.Equal(t, "11", ver)
}

func TestFindJavaCompilerVersionFromReader_Droplet(t *testing.T) {
	testFindJavaCompilerVersionFromReader(t, "droplet.tgz", "./workspace/BOOT-INF")
}

func TestFindJavaCompilerVersionFromReader_ImageLayerWithBootExtracted(t *testing.T) {
	testFindJavaCompilerVersionFromReader(t, "image-layer-boot-extracted.tgz", "/workspace/BOOT-INF")
}

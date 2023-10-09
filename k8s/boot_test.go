package k8s

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractBootJar_Found(t *testing.T) {

	tgzFile, err := os.Open(filepath.Join("test", "image-layer-boot-jar.tgz"))
	assert.NoError(t, err)

	found, extractedContentFileName, err := extractBootJarFromImageLayer(tgzFile, "/workspace/app.jar")
	assert.NoError(t, err)

	assert.True(t, found)

	expectedFileContent, err := os.ReadFile(filepath.Join("test", "app.jar"))
	assert.NoError(t, err)

	extractedFileContent, err := os.ReadFile(extractedContentFileName)
	assert.NoError(t, err)

	assert.Equal(t, expectedFileContent, extractedFileContent)
}

func TestExtractBootJar_NotFound(t *testing.T) {

	tgzFile, err := os.Open(filepath.Join("test", "image-layer-boot-jar.tgz"))
	assert.NoError(t, err)

	found, _, err := extractBootJarFromImageLayer(tgzFile, "/doesnt-exist.jar")
	assert.NoError(t, err)

	assert.False(t, found)
}

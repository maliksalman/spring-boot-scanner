package k8s_test

import (
	"github.com/maliksalman/spring-boot-scanner/k8s"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractBootJar_Found(t *testing.T) {

	tgzFile, err := os.Open(filepath.Join("testdata", "image-layer-boot-jar.tgz"))
	assert.NoError(t, err)

	found, extractedContentFileName, err := k8s.ExtractBootJarFromImageLayer(tgzFile, "/workspace/app.jar")
	assert.NoError(t, err)

	assert.True(t, found)

	expectedFileContent, err := os.ReadFile(filepath.Join("testdata", "app.jar"))
	assert.NoError(t, err)

	extractedFileContent, err := os.ReadFile(extractedContentFileName)
	assert.NoError(t, err)

	assert.Equal(t, expectedFileContent, extractedFileContent)
}

func TestExtractBootJar_NotFound(t *testing.T) {

	tgzFile, err := os.Open(filepath.Join("testdata", "image-layer-boot-jar.tgz"))
	assert.NoError(t, err)

	found, _, err := k8s.ExtractBootJarFromImageLayer(tgzFile, "/doesnt-exist.jar")
	assert.NoError(t, err)

	assert.False(t, found)
}

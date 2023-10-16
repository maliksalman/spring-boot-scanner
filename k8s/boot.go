package k8s

import (
	"archive/tar"
	"compress/gzip"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/maliksalman/spring-boot-scanner/scan"
	"io"
	"log"
	"os"
	"strings"
)

type JavaInfo struct {
	CompilerVersion   string
	RuntimeVersion    string
	SpringBootVersion string
}

func FindJavaInfoFromExtractedImage(extractedImagePath string) (bool, *JavaInfo, error) {

	// get the runtime information from the image
	runtimeInfo := scan.FindRuntimeInfoFromPath(extractedImagePath)

	// get out early if not spring-boot
	if runtimeInfo.BootVersion == "" {
		return false, nil, nil
	}

	// get the array of OCI image's layers
	ociLayers, err := getOciLayersFromExtractedImage(extractedImagePath)
	if err != nil {
		return false, nil, err
	}

	if runtimeInfo.BootContentExtracted {
		// go in reverse order of manifest and open each layer (tgz) to look for .class file
		for i := len(ociLayers); i > 0; i-- {
			layerTarGzFile, err := ociLayers[i-1].Compressed()
			if err != nil {
				return false, nil, err
			}
			// check for .class file in tgz file
			found, compilerVersion := scan.FindJavaCompilerVersionFromReader(layerTarGzFile, runtimeInfo.BootContentPrefix)

			if found {
				return true, &JavaInfo{
					CompilerVersion:   compilerVersion,
					RuntimeVersion:    runtimeInfo.JavaRuntimeVersion,
					SpringBootVersion: runtimeInfo.BootVersion,
				}, nil
			}
		}
	} else {
		// go in reverse order of manifest and open each layer (tgz) to look for boot-jar
		for i := len(ociLayers); i > 0; i-- {
			layerTarGzFile, err := ociLayers[i-1].Compressed()
			if err != nil {
				return false, nil, err
			}

			// some debug info
			layerDigest, _ := ociLayers[i-1].Digest()
			log.Printf("Looking for boot-jar: Jar=%s, LayerSha256=%s", runtimeInfo.BootJar, layerDigest.Hex)

			// find the boot-jar from the layer
			jarFileFound, jarFileName, err := ExtractBootJarFromImageLayer(layerTarGzFile, runtimeInfo.BootJar)
			if err != nil {
				return false, nil, err
			}

			// open the boot-jar and look for .class file
			if jarFileFound {
				compilerVersion, err := scan.FindJavaCompilerVersionFromBootJar(jarFileName, runtimeInfo.BootContentPrefix)
				if err != nil {
					return false, nil, err
				}
				deleteBootJar(jarFileName)

				return true, &JavaInfo{
					CompilerVersion:   compilerVersion,
					RuntimeVersion:    runtimeInfo.JavaRuntimeVersion,
					SpringBootVersion: runtimeInfo.BootVersion,
				}, nil
			}
		}
	}

	return true, nil, nil
}

func getOciLayersFromExtractedImage(extractedImagePath string) ([]v1.Layer, error) {

	path, err := layout.FromPath(extractedImagePath)
	if err != nil {
		return nil, err
	}

	index, err := path.ImageIndex()
	if err != nil {
		return nil, err
	}

	manifest, err := index.IndexManifest()
	if err != nil {
		return nil, err
	}

	image, err := path.Image(manifest.Manifests[0].Digest)
	if err != nil {
		return nil, err
	}

	return image.Layers()
}

func ExtractBootJarFromImageLayer(tarGzippedReader io.Reader, bootJarName string) (bool, string, error) {

	gzipReader, _ := gzip.NewReader(tarGzippedReader)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if hdr.Typeflag == tar.TypeReg && strings.HasSuffix(hdr.Name, bootJarName) {

			target, err := os.CreateTemp("", "boot-jar-*")
			if err != nil {
				return false, "", err
			}

			// copy over contents
			if _, err := io.Copy(target, tarReader); err != nil {
				return false, "", err
			}

			target.Close()
			return true, target.Name(), nil
		}
	}

	return false, "", nil
}

func deleteBootJar(bootJarPath string) error {
	return os.Remove(bootJarPath)
}
